package task

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	cpuprofile "main/cpu_profile"
	"sort"
	"strings"

	"github.com/google/pprof/profile"
)

type TagsProfile struct {
	Key     string
	Tags    []string
	Value   int64   // pprof cpu times
	Percent float64 // <= 1.0
}

func labelToTags(label map[string][]string) []string {
	tags := make([]string, 0, len(label)*2)
	for k, v := range label {
		tags = append(tags, k, strings.Join(v, ","))
	}
	return tags
}

func tagsToKey(tags []string) string {
	if len(tags)%2 != 0 {
		return ""
	}
	tagsPair := make([]string, 0, len(tags)/2)
	for i := 0; i < len(tags); i += 2 {
		tagsPair = append(tagsPair, fmt.Sprintf("%s=%s", tags[i], tags[i+1]))
	}
	// sort tags to make it a unique key
	sort.Strings(tagsPair)
	return strings.Join(tagsPair, "|")
}

func Analyse(data *bytes.Buffer) ([]*TagsProfile, error) {
	// parse protobuf data
	pf, err := profile.ParseData(data.Bytes())
	if err != nil {
		return nil, err
	}

	// filter cpu value index
	sampleIdx := -1
	for idx, st := range pf.SampleType {
		if st.Type == "cpu" {
			sampleIdx = idx
			break
		}
	}
	if sampleIdx < 0 {
		return nil, errors.New("profiler: sample type not found")
	}

	// calculate every sample expense
	counter := map[string]*TagsProfile{} // map[tagsKey]funcProfile
	var total int64
	for _, sm := range pf.Sample {
		value := sm.Value[sampleIdx]
		tags := labelToTags(sm.Label)
		tagsKey := tagsToKey(tags)
		tp, ok := counter[tagsKey]
		if !ok {
			tp = &TagsProfile{}
			counter[tagsKey] = tp
			tp.Key = tagsKey
			tp.Tags = tags
		}
		tp.Value += value
		total += value
	}

	profiles := make([]*TagsProfile, 0, len(counter)) // flat to array
	for _, l := range counter {
		l.Percent = float64(l.Value) / float64(total)
		profiles = append(profiles, l)
	}
	return profiles, nil
}

func prinfProfilerInfo(profiles []*TagsProfile) {
	if len(profiles) == 0 {
		return
	}
	log.Printf("profiler collect %d records", len(profiles))
	for _, p := range profiles {
		if p.Key != "" {
			log.Printf("profiler - %s %.2f%%", p.Key, p.Percent*100)
		} else {
			log.Printf("profiler - type=default %.2f%%", p.Percent*100)
		}
	}
	log.Println("---------------------------------")
	monitor.graphData = append(monitor.graphData, profiles)
}

func handleTaskProfile(profileData *cpuprofile.ProfileData) error {
	if profileData.Error != nil {
		return profileData.Error
	}
	profiles, err := Analyse(profileData.Data)
	if err != nil {
		return err
	}
	prinfProfilerInfo(profiles)
	return nil
}


func StartCPUProfile(enableTaskWebVisual bool) {
	consumer := cpuprofile.NewConsumer(handleTaskProfile)
	consumer.StartConsume()
	go graphServe()
}
