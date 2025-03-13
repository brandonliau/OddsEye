package processor

import (
	"fmt"
	"io"
	"net/url"
	"slices"

	"OddsEye/pkg/retryhttp"
)

const (
	numWorkers = 150
	batchSize  = 5
)

type Processor interface {
	Execute()
}

func fetchData(baseURL string, params url.Values, client *retryhttp.RetryClient) ([]byte, error) {
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func batchJobs(fixtures []string, sportsbooks []string) ([][]string, [][]string) {
	// create fixture batches
	batchFixtures := make([][]string, 0)
	for c := range slices.Chunk(fixtures, batchSize) {
		batchFixtures = append(batchFixtures, c)
	}

	// create sportsbook batches
	batchSportsbooks := make([][]string, 0)
	for c := range slices.Chunk(sportsbooks, batchSize) {
		batchSportsbooks = append(batchSportsbooks, c)
	}

	return batchFixtures, batchSportsbooks
}
