package model

import (
	"time"
)

type FixturesWrapper struct {
	Data       []Fixture `json:"data"`
	Page       int       `json:"page"`
	TotalPages int       `json:"total_pages"`
}

type Fixture struct {
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
