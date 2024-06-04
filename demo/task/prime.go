package task

import (
	"context"
	"runtime/pprof"
	"sync"
)

const primeBound = 10000

func selectPrime(ctx context.Context, c chan int, wg *sync.WaitGroup, enableProfile bool) {
	if enableProfile {
		pprof.SetGoroutineLabels(ctx)
	}
	defer wg.Done()
	prime, ok := <-c
	if !ok {
		return
	}
	// fmt.Println(prime)
	newChan := make(chan int)
	newWg := sync.WaitGroup{}
	newWg.Add(1)
	go selectPrime(ctx, newChan, &newWg, enableProfile)
	for n := range c {
		if n%prime != 0 {
			newChan <- n
		}
	}
	close(newChan)
	newWg.Wait()
}

func PrintPrime(ctx context.Context, label string, enableProfile bool, externWg *sync.WaitGroup) {
	// 筛法求素数
	defer externWg.Done()
	if enableProfile {
		defer pprof.SetGoroutineLabels(ctx)
		ctx = pprof.WithLabels(ctx, pprof.Labels("task", label))
		pprof.SetGoroutineLabels(ctx)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	c := make(chan int)
	n := primeBound
	go selectPrime(ctx, c, &wg, enableProfile)
	for i := 2; i <= n; i++ {
		c <- i
	}
	close(c)
	wg.Wait()
}
