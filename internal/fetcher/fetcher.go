package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"odds-eye/pkg/retryhttp"
)

type Fetcher[T any] interface {
	Fetch(wg *sync.WaitGroup, jobs chan T, results chan []byte)
}

func transport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
}

func fetchData(client *retryhttp.RetryClient, baseURL string, params url.Values) ([]byte, error) {
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
