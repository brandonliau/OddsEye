package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	"OddsEye/pkg/config"
	"OddsEye/pkg/retryhttp"
)

func fetchData(config *config.ServiceConfig, client *retryhttp.RetryClient, baseURL string, params url.Values) ([]byte, error) {
	params.Add("key", config.Token)
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

func unmarshalData(body []byte, result any) error {
	err := json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}

func (s *oddsService) Fixtures(sport string, league string) []fixture {
	baseURL := "https://api.opticodds.com/api/v3/fixtures/active"
	params := url.Values{}
	params.Add("sport", sport)
	params.Add("league", league)
	params.Add("start_date_after", time.Now().Format(time.RFC3339))
	params.Add("start_date_before", time.Now().Add(time.Duration(s.cfg.Window*24)*time.Hour).Format(time.RFC3339))

	var fixturesWrapper fixturesWrapper
	body, err := fetchData(s.cfg, s.client, baseURL, params)
	if err != nil {
		s.logger.Error("Failed to fetch fixtures: %v", err)
	}
	err = unmarshalData(body, &fixturesWrapper)
	if err != nil {
		s.logger.Error("Failed to unmarshal fixtures: %v", err)
	}
	fixtures := append([]fixture{}, fixturesWrapper.Data...)

	// handle pagination
	for fixturesWrapper.Page < fixturesWrapper.TotalPages {
		fixturesWrapper.Page += 1
		params.Set("page", strconv.Itoa(fixturesWrapper.Page))

		body, err := fetchData(s.cfg, s.client, baseURL, params)
		if err != nil {
			s.logger.Error("Failed to fetch fixtures: %v", err)
		}
		err = unmarshalData(body, &fixturesWrapper)
		if err != nil {
			s.logger.Error("Failed to unmarshal fixtures: %v", err)
		}
		fixtures = append(fixtures, fixturesWrapper.Data...)
	}
	return fixtures
}

func (s *oddsService) Odds(fixtureIDs []string, sportsbooks []string) []fixtureOdds {
	baseURL := "https://api.opticodds.com/api/v3/fixtures/odds"
	params := url.Values{}
	params.Add("odds_format", "decimal")
	for _, id := range fixtureIDs {
		params.Add("fixture_id", id)
	}
	for _, sportsbook := range sportsbooks {
		params.Add("sportsbook", sportsbook)
	}

	var fixtureOddsWrapper fixtureOddsWrapper
	body, err := fetchData(s.cfg, s.client, baseURL, params)
	if err != nil {
		s.logger.Error("Failed to fetch fixture odds: %v", err)
	}
	err = unmarshalData(body, &fixtureOddsWrapper)
	if err != nil {
		s.logger.Error("Failed to unmarshal fixture odds: %v", err)
	}
	return fixtureOddsWrapper.Data
}
