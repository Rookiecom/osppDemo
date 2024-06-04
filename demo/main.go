package main

import (
	"context"
	"flag"
	cpuprofile "main/cpu_profile"
	"main/task"
	"sync"
	"time"
)

var (
	enableProfile       bool
	enableTaskWebVisual bool
	window              *int
	interval            *int
)

func init() {
	flag.BoolVar(&enableProfile, "p", true, "是否开启task的profile功能")
	flag.BoolVar(&enableTaskWebVisual, "w", true, "是否开启task的CPU占用量web可视化")
	window = flag.Int("pwindow", 1000, "profile采样的窗口大小，以毫秒为单位")
	interval = flag.Int("pinterval", 500, "profile每次采样的间隔时间，以毫秒为单位")
	flag.Parse()
}

func main() {
	go cpuprofile.StartCPUProfiler(time.Duration(*window*int(time.Millisecond)), time.Duration(*interval*int(time.Millisecond)))
	if enableProfile {
		task.StartCPUProfile(enableTaskWebVisual)
	}
	cpuprofile.WebProfile("localhost:8080")
	ctx := context.Background()
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(2)
		go task.PrintPrime(ctx, "prime", enableProfile, &wg)
		go task.TaskMergeSort(ctx, "mergeSort", enableProfile, &wg)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
}
