package service

import (
	"time"
)

type fixturesWrapper struct {
	Data       []fixture `json:"data"`
	Page       int       `json:"page"`
	TotalPages int       `json:"total_pages"`
}

type fixture struct {
	ID        string    `json:"id"`
	StartDate time.Time `json:"start_date"`
	HomeTeam  string    `json:"home_team_display"`
	AwayTeam  string    `json:"away_team_display"`
	HasOdds   bool      `json:"has_odds"`
	Sport     struct {
		ID string `json:"id"`
	} `json:"sport"`
	League struct {
		ID string `json:"id"`
	} `json:"league"`
}

type fixtureOddsWrapper struct {
	Data []fixtureOdds `json:"data"`
}

type fixtureOdds struct {
	ID   string `json:"id"`
	Odds []odd  `json:"odds"`
}

type odd struct {
	ID                  string   `json:"id"`
	Sportsbook          string   `json:"sportsbook"`
	Market              string   `json:"market_id"`
	Selection           string   `json:"selection"`
	NormalizedSelection string   `json:"normalized_selection"`
	SelectionLine       string   `json:"selection_line"`
	Points              *float64 `json:"points"`
	Price               float64  `json:"price"`
	GroupingKey         string   `json:"grouping_key"`
	DeepLink            struct {
		Ios     string `json:"ios"`
		Desktop string `json:"desktop"`
	} `json:"deep_link"`
}
