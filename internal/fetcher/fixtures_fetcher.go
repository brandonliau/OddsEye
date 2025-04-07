package fetcher

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"OddsEye/internal/model"

	"OddsEye/pkg/config"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type fixturesFetcher struct {
	config *config.FixturesConfig
	client *retryhttp.RetryClient
	logger logger.Logger
}

type FixturesFetchJob struct {
	Sport  string
	League string
}

func NewFixturesFetcher(config *config.FixturesConfig, logger logger.Logger) *fixturesFetcher {
	httpClient := &http.Client{
		Transport: transport(),
		Timeout:   10 * time.Second,
	}
	client := retryhttp.NewRetryClient(httpClient, logger)

	return &fixturesFetcher{
		config: config,
		client: client,
		logger: logger,
	}
}

func (f *fixturesFetcher) Fetch(wg *sync.WaitGroup, jobs chan FixturesFetchJob, results chan []byte) {
	defer wg.Done()

	baseURL := f.config.API.BaseURL

	curr := time.Now()
	startDateAfter := curr.Format(time.RFC3339)
	startDateBefore := curr.AddDate(0, 0, f.config.API.Window).Format(time.RFC3339)

	params := url.Values{}
	params.Set("key", f.config.API.Token)
	params.Set("start_date_after", startDateAfter)
	params.Set("start_date_before", startDateBefore)

	for param := range jobs {
		params.Set("sport", param.Sport)
		params.Set("league", param.League)

		body, err := fetchData(f.client, baseURL, params)
		if err != nil {
			f.logger.Error("Failed to fetch fixture data: %v", err)
		}

		var fixturesWrapper model.FixturesWrapper
		err = json.Unmarshal(body, &fixturesWrapper)
		if err != nil {
			f.logger.Error("Failed to unmarshal fixtures: %v", err)
		}

		// handle pagination
		for fixturesWrapper.Page < fixturesWrapper.TotalPages {
			fixturesWrapper.Page += 1
			params.Set("page", strconv.Itoa(fixturesWrapper.Page))

			nextPage, err := fetchData(f.client, baseURL, params)
			if err != nil {
				f.logger.Error("Failed to fetch fixture data: %v", err)
			}

			body = append(body, nextPage...)
			err = json.Unmarshal(nextPage, &fixturesWrapper)
			if err != nil {
				f.logger.Error("Failed to unmarshal fixtures: %v", err)
			}
		}

		results <- body
	}
}
