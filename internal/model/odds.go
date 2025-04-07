package model

type FixtureOddsWrapper struct {
	Data []FixtureOdds `json:"data"`
}

type FixtureOdds struct {
	ID   string `json:"id"`
	Odds []Odd  `json:"odds"`
}

type Odd struct {
	ID                  string   `json:"id"`
	Sportsbook          string   `json:"sportsbook"`
	Market              string   `json:"market_id"`
	Selection           string   `json:"selection"`
	NormalizedSelection string   `json:"normalized_selection"`
	SelectionLine       string   `json:"selection_line"`
	Points              *float64 `json:"points"`
	Price               float64  `json:"price"`
	DeepLink            struct {
		Desktop string `json:"desktop"`
	} `json:"deep_link"`
}
