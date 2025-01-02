package service

import (
	// "regexp"
	// "strings"
	"sync"
)

// var (
// 	normalizeTeamRe = regexp.MustCompile(`[\s&]+`)
// )

type batchedJob struct {
	fixtureID  []string
	sportsbook []string
}

func createBatchJobs(fixtures []string, sportsbooks []string, batchSize int) (int, []batchedJob) {
	// create sportsbook batches
	numSportsbookBatches := (len(sportsbooks) + batchSize - 1) / batchSize
	batchSportsbooks := make([][]string, 0, numSportsbookBatches)
	for i := 0; i < len(sportsbooks); i += batchSize {
		end := i + batchSize
		if end > len(sportsbooks) {
			end = len(sportsbooks)
		}
		batchSportsbooks = append(batchSportsbooks, sportsbooks[i:end])
	}

	// create fixture batches
	numJobBatches := (len(fixtures) + batchSize - 1) / batchSize
	batchJobs := make([]batchedJob, 0, numJobBatches*numSportsbookBatches)
	for i := 0; i < len(fixtures); i += batchSize {
		end := i + batchSize
		if end > len(fixtures) {
			end = len(fixtures)
		}
		for _, sportsbooks := range batchSportsbooks {
			job := batchedJob{
				fixtureID:  fixtures[i:end],
				sportsbook: sportsbooks,
			}
			batchJobs = append(batchJobs, job)
		}
	}

	return len(fixtures) * len(sportsbooks), batchJobs
}

func launchWorkers[T any, K any](numWorkers int, jobs chan T, results chan K, process func(*sync.WaitGroup, chan T, chan K)) {
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go process(&wg, jobs, results)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

func distributeJobs[T any](tasks []T, jobs chan T) {
	go func() {
		for _, job := range tasks {
			jobs <- job
		}
		close(jobs)
	}()
}

// func normalizeTeam(team string) string {
// 	lower := strings.ToLower(team)
// 	return normalizeTeamRe.ReplaceAllString(lower, "_")
// }
