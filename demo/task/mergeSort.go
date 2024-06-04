package task

import (
	"context"
	"math/rand"
	"runtime/pprof"
	"sort"
	"sync"
)

const minBound = 10
const arrayLen = 500000

func TaskMergeSort(ctx context.Context, label string, enableProfile bool, externWg *sync.WaitGroup) {
	defer externWg.Done()
	if enableProfile {
		defer pprof.SetGoroutineLabels(ctx)
		ctx = pprof.WithLabels(ctx, pprof.Labels("task", label))
		pprof.SetGoroutineLabels(ctx)

	}

	array := make([]int, arrayLen)
	for i := range array {
		array[i] = rand.Int()
	}
	ParallelMergeSort(ctx, array, enableProfile)
}

func ParallelMergeSort(ctx context.Context, array []int, enableProfile bool) {
	if enableProfile {
		pprof.SetGoroutineLabels(ctx)
	}
	if len(array) <= minBound {
		sort.Ints(array)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	leftArray := make([]int, len(array)/2)
	rightArray := make([]int, len(array)-len(array)/2)
	copy(leftArray[0:], array[0:len(array)/2])
	copy(rightArray[0:], array[len(array)/2:])
	go func() {
		ParallelMergeSort(ctx, leftArray, enableProfile)
		wg.Done()
	}()
	go func() {
		ParallelMergeSort(ctx, rightArray, enableProfile)
		wg.Done()
	}()
	wg.Wait()
	cnt := 0
	i := 0
	j := 0
	for i < len(array)/2 && j < len(array)-len(array)/2 {
		if leftArray[i] < rightArray[j] {
			array[cnt] = leftArray[i]
			i++
		} else {
			array[cnt] = rightArray[j]
			j++
		}
		cnt++
	}
	if i < len(array)/2 {
		copy(array[cnt:], leftArray[i:])
	} else {
		copy(array[cnt:], rightArray[j:])
	}
}
