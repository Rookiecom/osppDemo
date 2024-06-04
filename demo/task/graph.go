package task

import (
	"io"
	"log"
	"net/http"
	"sync"
	"text/template"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Monitor struct {
	graphData [][]*TagsProfile
	MaxRecord int
	Lock      sync.Mutex
}

func NewMonitor(max int) *Monitor {
	return &Monitor{
		MaxRecord: max,
	}
}

var monitor = NewMonitor(50)

func (m *Monitor) WriteTo(w io.Writer) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	keyTOlineIndex := make(map[string]int)
	indexTOkey := make(map[int]string)
	lineData := make([]plotter.XYs, 0)
	cnt := 0
	base := max(len(m.graphData)-m.MaxRecord, 0)
	for i := base; i < len(m.graphData); i++ {
		for j := range m.graphData[i] {
			if _, ok := keyTOlineIndex[m.graphData[i][j].Key]; !ok {
				keyTOlineIndex[m.graphData[i][j].Key] = cnt
				indexTOkey[cnt] = m.graphData[i][j].Key
				cnt++
				lineData = append(lineData, make(plotter.XYs, len(m.graphData)-base))
			}
			index := keyTOlineIndex[m.graphData[i][j].Key]
			lineData[index][i-base].X = float64(i + 1 - base)
			lineData[index][i-base].Y = float64(m.graphData[i][j].Percent)
		}
	}

	p := plot.New()
	p.Title.Text = "Change chart of CPU proportion of each tag"
	colorMap := make(map[string]int)
	colorMap["task=mergeSort"] = 0
	colorMap["task=prime"] = 1
	colorMap[""] = 2
	for i := 0; i < len(lineData); i++ {
		newLine, _ := plotter.NewLine(lineData[i])
		newLine.Color = plotutil.Color(colorMap[indexTOkey[i]])
		p.Add(newLine)
		if indexTOkey[i] == "" {
			indexTOkey[i] = "default"
		}
		p.Legend.Add(indexTOkey[i], newLine)
	}

	p.X.Min = 0
	p.X.Max = float64(m.MaxRecord)
	p.Y.Min = 0
	p.Y.Max = 1

	wc, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
	if err != nil {
		log.Fatal(err)
	}
	wc.WriteTo(w)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("task/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t.Execute(w, nil)
}

func image(w http.ResponseWriter, r *http.Request) {
	monitor.WriteTo(w)
}

func graphServe() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/image", image)

	s := &http.Server{
		Addr:    "localhost:8090",
		Handler: mux,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
