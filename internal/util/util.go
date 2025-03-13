package util

import (
	"sync"
)

func LaunchWorkers[T any, K any](numWorkers int, jobs chan T, results chan K, work func(wg *sync.WaitGroup, jobs chan T, results chan K)) {
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go work(&wg, jobs, results)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

func DistributeJobs[T any](tasks []T, jobs chan T) {
	go func() {
		for _, job := range tasks {
			jobs <- job
		}
		close(jobs)
	}()
}
