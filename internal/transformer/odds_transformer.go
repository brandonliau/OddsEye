package transformer

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"

	"OddsEye/internal/model"

	"OddsEye/pkg/config"
	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
)

type OddsTransformer struct {
	config *config.OddsConfig
	db     database.Database
	logger logger.Logger
}

func NewOddsTransformer(config *config.OddsConfig, db database.Database, logger logger.Logger) *OddsTransformer {
	return &OddsTransformer{
		config: config,
		db:     db,
		logger: logger,
	}
}

func groupingKey(normalizedSelection, selection, selectionLine, home, away string, points *float64) string {
	groupingKey := "default"
	if normalizedSelection != "" && selectionLine != "" && points != nil {
		groupingKey = fmt.Sprintf("%s:%.2f", normalizedSelection, math.Abs(*points))
	} else if normalizedSelection != "" && selectionLine != "" {
		groupingKey = normalizedSelection
	} else if points != nil {
		if selection == home {
			groupingKey = fmt.Sprintf("default:%.2f", *points)
		} else if selection == away {
			groupingKey = fmt.Sprintf("default:%.2f", -*points)
		} else if selection == "Draw" {
			groupingKey = fmt.Sprintf("default:%.2f", -*points)
		} else {
			groupingKey = fmt.Sprintf("default:%.2f", *points)
		}
	}
	return groupingKey
}

func (t *OddsTransformer) teams(fixtureID string) (string, string) {
	query := "SELECT home_team, away_team FROM all_fixtures WHERE id = ?"
	rows, err := t.db.Query(query, fixtureID)
	if err != nil {
		t.logger.Error("Failed to retrieve teams: %v", err)
	}
	defer rows.Close()

	var home, away string
	for rows.Next() {
		err := rows.Scan(&home, &away)
		if err != nil {
			t.logger.Error("Failed to scan row: %v", err)
		}
	}
	return home, away
}

func (t *OddsTransformer) normalizeURL(url string) string {
	for k, v := range t.config.Transformer.Normalization {
		url = strings.ReplaceAll(url, k, v)
	}
	return url
}

func (t *OddsTransformer) Transform(wg *sync.WaitGroup, jobs chan []byte, results chan int) {
	defer wg.Done()

	stmt, _ := t.db.PrepareExec(t.config.Transformer.Insertion)

	for data := range jobs {
		var fixtureOddsWrapper model.FixtureOddsWrapper
		err := json.Unmarshal(data, &fixtureOddsWrapper)
		if err != nil {
			t.logger.Error("Failed to unmarshal fixture odds: %v", err)
		}

		fixtureOdds := fixtureOddsWrapper.Data
		for _, fixtureOdd := range fixtureOdds {
			id := fixtureOdd.ID
			home, away := t.teams(id)

			for _, odd := range fixtureOdd.Odds {
				_, err := stmt.Exec(
					id,
					odd.Market,
					strings.Split(odd.ID, ":")[3],
					odd.Sportsbook,
					odd.Price,
					t.normalizeURL(odd.DeepLink.Desktop),
					groupingKey(odd.NormalizedSelection, odd.Selection, odd.SelectionLine, home, away, odd.Points),
				)
				if err != nil {
					t.logger.Error("Failed to execute odds insertion statement: %v", err)
					continue
				}
				results <- 0
			}
		}
	}
}
