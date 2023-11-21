package main

import (
	"fmt"
	"runtime"
	"sync"
)

func parallelSearchForkJoin(arr []int, target int) bool {
	length := len(arr)
	resultCh := make(chan bool, runtime.NumCPU())
	var wg sync.WaitGroup

	searchInSubarray := func(start, end int) {
		defer wg.Done()

		for i := start; i < end; i++ {
			if arr[i] == target {
				resultCh <- true
				return
			}
		}
		resultCh <- false
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		start := i * (length / runtime.NumCPU())
		end := (i + 1) * (length / runtime.NumCPU())
		if i == runtime.NumCPU()-1 {
			end = length
		}
		go searchInSubarray(start, end)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		if result {
			return true
		}
	}
	return false
}

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	target := 8

	found := parallelSearchForkJoin(arr, target)
	if found {
		fmt.Printf("Элемент %d найден в массиве.\n", target)
	} else {
		fmt.Printf("Элемент %d не найден в массиве.\n", target)
	}
}
