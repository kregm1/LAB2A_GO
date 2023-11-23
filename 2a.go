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
	var mu sync.Mutex

	searchInSubarray := func(start, end int) {
		defer wg.Done()

		for i := start; i < end; i++ {
			if arr[i] == target {
				mu.Lock()
				select {
				case resultCh <- true:
				default:
				}
				mu.Unlock()
				return
			}
		}
		mu.Lock()
		select {
		case resultCh <- false:
		default:
		}
		mu.Unlock()
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

func parallelSearchThreadPool(arr []int, target int, numThreads int) bool {
	length := len(arr)
	resultCh := make(chan bool, 1)
	var wg sync.WaitGroup
	var mu sync.Mutex

	searchInSubarray := func(start, end int) {
		defer wg.Done()

		for i := start; i < end; i++ {
			if arr[i] == target {
				mu.Lock()
				select {
				case resultCh <- true:
				default:
				}
				mu.Unlock()
				return
			}
		}
		mu.Lock()
		select {
		case resultCh <- false:
		default:
		}
		mu.Unlock()
	}

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		start := i * (length / numThreads)
		end := (i + 1) * (length / numThreads)
		if i == numThreads-1 {
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
	numThreads := runtime.NumCPU()

	found := parallelSearchForkJoin(arr, target)
	if found {
		fmt.Printf("Элемент %d найден с использованием Fork-Join Pool.\n", target)
	} else {
		fmt.Printf("Элемент %d не найден с использованием Fork-Join Pool.\n", target)
	}

	foundTP := parallelSearchThreadPool(arr, target, numThreads)
	if foundTP {
		fmt.Printf("Элемент %d найден с использованием произвольного пула потоков.\n", target)
	} else {
		fmt.Printf("Элемент %d не найден с использованием произвольного пула потоков.\n", target)
	}
}
