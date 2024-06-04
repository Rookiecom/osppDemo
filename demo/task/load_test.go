package task

import (
	"context"
	"fmt"
	cpuprofile "main/cpu_profile"
	"strconv"
	"sync"
	"testing"
	"time"
)

func parallelStartNprime(ctx context.Context, number int, wg *sync.WaitGroup) {
	for i := 0; i < number; i++ {
		wg.Add(1)
		go PrintPrime(ctx, "prime"+" "+strconv.Itoa(number)+"xload", true, wg)
	}
}

func startNprime(ctx context.Context, number int) {
	wg := sync.WaitGroup{}
	for i := 0; i < number; i++ {
		wg.Add(1)
		PrintPrime(ctx, "prime"+" "+strconv.Itoa(number)+"xload", true, &wg)
		wg.Wait()
	}
}

type testConsumer struct {
	mergeData []*TagsProfile
	splitDate [][]*TagsProfile
	consumer  *cpuprofile.Consumer
}

func (t *testConsumer) testDataHandle(profileData *cpuprofile.ProfileData) error {
	if profileData.Error != nil {
		return profileData.Error
	}
	profiles, err := Analyse(profileData.Data)
	if err != nil {
		return err
	}
	t.splitDate = append(t.splitDate, profiles)
	for i := 0; i < len(profiles); i++ {
		j := 0
		for ; j < len(t.mergeData); j++ {
			if profiles[i].Key == t.mergeData[j].Key {
				t.mergeData[j].Value += profiles[i].Value
				break
			}
		}
		if j == len(t.mergeData) {
			t.mergeData = append(t.mergeData, profiles[i])
		}
	}
	return nil
}

func TestLoad(t *testing.T) {
	fmt.Println("并行筛法求素数 20倍负载 和 200倍负载测试开始")
	fmt.Println("N倍负载同时启动：------")
	cpuprofile.StartCPUProfiler(time.Duration(1000)*time.Millisecond, time.Duration(500)*time.Millisecond)
	ParallelConsumer := testConsumer{}
	ParallelConsumer.consumer = cpuprofile.NewConsumer(ParallelConsumer.testDataHandle)
	ParallelConsumer.consumer.StartConsume()
	ctx := context.Background()
	wg := sync.WaitGroup{}
	parallelStartNprime(ctx, 20, &wg)
	parallelStartNprime(ctx, 200, &wg)
	wg.Wait()
	ParallelConsumer.consumer.StopConsume()
	for i := range ParallelConsumer.mergeData {
		if ParallelConsumer.mergeData[i].Key == "" {
			continue
		}
		fmt.Printf("任务：%s\n  CPU使用量：%d\n", ParallelConsumer.mergeData[i].Key, ParallelConsumer.mergeData[i].Value)
	}
	fmt.Println("N倍负载串行启动测试完成")
	// SerialConsumer := testConsumer{}
	// SerialConsumer.consumer = cpuprofile.NewConsumer(SerialConsumer.testDataHandle)
	// SerialConsumer.consumer.StartConsume()
	// startNprime(ctx, 20)
	// startNprime(ctx, 200)
	// SerialConsumer.consumer.StopConsume()
	// for i := range SerialConsumer.mergeData {
	// 	if SerialConsumer.mergeData[i].Key == "" {
	// 		continue
	// 	}
	// 	fmt.Printf("任务：%s\n  CPU使用量：%d\n", SerialConsumer.mergeData[i].Key, SerialConsumer.mergeData[i].Value)
	// }

}
