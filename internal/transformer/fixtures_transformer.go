package transformer

import (
	"encoding/json"
	"sync"

	"odds-eye/internal/model"

	"odds-eye/pkg/config"
	"odds-eye/pkg/database"
	"odds-eye/pkg/logger"
)

type FixturesTransformer struct {
	config *config.FixturesConfig
	db     database.Database
	logger logger.Logger
}

func NewFixturesTransformer(config *config.FixturesConfig, db database.Database, logger logger.Logger) *FixturesTransformer {
	return &FixturesTransformer{
		config: config,
		db:     db,
		logger: logger,
	}
}

func (t *FixturesTransformer) Transform(wg *sync.WaitGroup, jobs chan []byte, results chan int) {
	defer wg.Done()

	stmt, _ := t.db.PrepareExec(t.config.Transformer.Insertion)

	for data := range jobs {
		var fixturesWrapper model.FixturesWrapper
		err := json.Unmarshal(data, &fixturesWrapper)
		if err != nil {
			t.logger.Error("Failed to unmarshal fixtures: %v", err)
		}
		fixtures := fixturesWrapper.Data

		for _, fixture := range fixtures {
			if !fixture.HasOdds {
				continue
			}
			_, err := stmt.Exec(fixture.ID, fixture.StartDate, fixture.HomeTeam, fixture.AwayTeam, fixture.Sport.ID, fixture.League.ID)
			if err != nil {
				t.logger.Error("Failed to execute fixtures insertion statement: %v", err)
				continue
			}

			results <- 0
		}
	}
}
