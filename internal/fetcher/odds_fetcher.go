package fetcher

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"OddsEye/pkg/config"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/retryhttp"
)

type oddsFetcher struct {
	config *config.OddsConfig
	client *retryhttp.RetryClient
	logger logger.Logger
}

type OddsFetchJob struct {
	Fixtures    []string
	Markets     []string
	Sportsbooks []string
}

func NewOddsFetcher(config *config.OddsConfig, logger logger.Logger) *oddsFetcher {
	httpClient := &http.Client{
		Transport: transport(),
		Timeout:   10 * time.Second,
	}
	client := retryhttp.NewRetryClient(httpClient, logger)

	return &oddsFetcher{
		config: config,
		client: client,
		logger: logger,
	}
}

func (f *oddsFetcher) Fetch(wg *sync.WaitGroup, jobs chan OddsFetchJob, results chan []byte) {
	defer wg.Done()

	baseURL := f.config.API.BaseURL

	params := url.Values{}
	params.Set("key", f.config.API.Token)
	params.Add("odds_format", f.config.API.Format)

	for param := range jobs {
		for _, fixture := range param.Fixtures {
			params.Add("fixture_id", fixture)
		}
		for _, market := range param.Markets {
			params.Add("market", market)
		}
		for _, sportsbook := range param.Sportsbooks {
			params.Add("sportsbook", sportsbook)
		}

		body, err := fetchData(f.client, baseURL, params)
		if err != nil {
			f.logger.Error("Failed to fetch odds data: %v", err)
		}

		results <- body
	}
}
