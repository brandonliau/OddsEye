package retryhttp

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	redirectsErrorRe     = regexp.MustCompile(`stopped after \d+ redirects\z`)
	schemeErrorRe        = regexp.MustCompile(`unsupported protocol scheme`)
	invalidHeaderErrorRe = regexp.MustCompile(`invalid header`)
	notTrustedErrorRe    = regexp.MustCompile(`certificate is not trusted`)
)

type Backoff func(retryMin time.Duration, retryMax time.Duration, attemptNum int) time.Duration

func DefaultBackoff(retryMin time.Duration, retryMax time.Duration, attemptNum int) time.Duration {
	min := float64(retryMin.Milliseconds())
	max := float64(retryMax.Milliseconds())
	mult := math.Min(min*math.Pow(2, float64(attemptNum)), max)
	jitter := rand.Float64() * min
	return time.Duration(mult+jitter) * time.Millisecond
}

type RetryPolicy func(resp *http.Response, err error) (bool, error)

func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects
			if redirectsErrorRe.MatchString(v.Error()) {
				return false, v
			}
			// Don't retry if the error was due to an invalid protocol scheme
			if schemeErrorRe.MatchString(v.Error()) {
				return false, v
			}
			// Don't retry if the error was due to an invalid header
			if invalidHeaderErrorRe.MatchString(v.Error()) {
				return false, v
			}
			// Don't retry if the error was due to TLS cert verification failure
			if notTrustedErrorRe.MatchString(v.Error()) {
				return false, v
			}
		}
		return true, nil
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return true, nil
	}
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
		return true, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}
	return false, nil
}
