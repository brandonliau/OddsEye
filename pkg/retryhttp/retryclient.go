package retryhttp

import (
	"fmt"
	"net/http"
	"time"

	"odds-eye/pkg/logger"
)

type RetryClient struct {
	client       *http.Client
	logger       logger.Logger
	RetryMax     int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	Backoff      Backoff
	RetryPolicy  RetryPolicy
}

func NewRetryClient(client *http.Client, logger logger.Logger) *RetryClient {
	return &RetryClient{
		client:       client,
		logger:       logger,
		RetryMax:     3,
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 5 * time.Second,
		Backoff:      DefaultBackoff,
		RetryPolicy:  DefaultRetryPolicy,
	}
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= c.RetryMax; attempt++ {
		c.logger.Debug("%s %s", req.Method, req.URL.String())
		resp, err = c.client.Do(req)

		retry, retryErr := c.RetryPolicy(resp, err)
		if !retry {
			// unrecoverable error, dont retry
			if retryErr != nil {
				return nil, retryErr
			}
			// return http response
			return resp, nil
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		if attempt < c.RetryMax {
			wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, attempt)
			c.logger.Debug("Retrying in %v (%d) %s %s", wait, c.RetryMax-attempt, req.Method, req.URL.String())
			time.Sleep(wait)
		}
	}
	return nil, fmt.Errorf("all retry attempts failed")
}

func (c *RetryClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
