package service

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

func (s *oddsService) SeedFixtures() {
	query := "INSERT INTO all_fixtures (id, start_date, home_team, away_team, sport, league) VALUES (?, ?, ?, ?, ?, ?)"
	s.seedFixtures(query)
}

func (s *oddsService) SeedFixtureOdds() {
	numOdds, batchedJobs := createBatchJobs(s.repo.Fixtures(), s.cfg.Sportsbooks, 5)
	query := ` INSERT INTO all_odds (id, market, selection, selection_line, points, sportsbook, price, url, grouping_key)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	s.seedOdds(150, batchedJobs, numOdds, query)
	s.clean("all_odds")
}

func (s *oddsService) SeedFairOdds() {
	numOdds, batchedJobs := createBatchJobs(s.repo.Fixtures(), s.cfg.Sharpbooks, 5)
	query := ` INSERT INTO fair_odds (id, market, selection, selection_line, points, sportsbook, price, url, grouping_key)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	s.seedOdds(150, batchedJobs, numOdds, query)
	s.clean("fair_odds")
}

func (s *oddsService) seedFixtures(query string) {
	start := time.Now()
	processed := 0

	// prepare insertion statement
	stmt, _ := s.db.PrepareExec(query)

	// seed fixtures for each sport and league combination
	s.db.Begin()
	for sport, leagues := range s.cfg.Sports {
		for _, league := range leagues {
			for _, fixture := range s.Fixtures(sport, league) {
				if !fixture.HasOdds {
					continue
				}
				_, err := stmt.Exec(fixture.ID, fixture.StartDate, fixture.HomeTeam, fixture.AwayTeam, fixture.Sport.ID, fixture.League.ID)
				if err != nil {
					s.logger.Error("Failed to execute statement: %v", err)
					continue
				}
				s.repo.AddFixture(fixture.ID, fixture.HomeTeam, fixture.AwayTeam)
				processed++
			}
		}
	}
	s.db.Commit()
	s.logger.Info("Seeded %d fixtures in %v", processed, time.Since(start))
}

func (s *oddsService) seedOdds(numWorkers int, batchedJobs []batchedJob, numOdds int, query string) {
	start := time.Now()
	processed := 0

	// prepare insertion statement
	stmt, _ := s.db.PrepareExec(query)

	// create jobs and results channel
	jobs := make(chan batchedJob, len(batchedJobs))
	results := make(chan fixtureOdds, numOdds)

	// launch worker goroutines
	launchWorkers(numWorkers, jobs, results, s.seedOddsProcess)

	// distribute jobs to workers
	distributeJobs(batchedJobs, jobs)

	s.db.Begin()
	for fixtureOdds := range results {
		home, away := s.repo.Teams(fixtureOdds.ID)
		for _, odds := range fixtureOdds.Odds {

			// create grouping key
			groupingKey := "default"
			if odds.NormalizedSelection != "" && odds.SelectionLine != "" && odds.Points != nil {
				groupingKey = fmt.Sprintf("%s:%.1f", odds.NormalizedSelection, math.Abs(*odds.Points))
			} else if odds.NormalizedSelection != "" && odds.SelectionLine != "" && (!strings.Contains(odds.Market, "correct_score") && !strings.Contains(odds.Market, "set_betting")) {
				groupingKey = odds.NormalizedSelection
			} else if odds.Points != nil {
				if odds.Selection == home {
					groupingKey = fmt.Sprintf("default:%.2f", *odds.Points)
				} else if odds.Selection == away {
					groupingKey = fmt.Sprintf("default:%.2f", -*odds.Points)
				} else if odds.Selection == "Draw" {
					groupingKey = fmt.Sprintf("default:%.2f", -*odds.Points)
				} else {
					groupingKey = fmt.Sprintf("default:%.2f", *odds.Points)
				}
			}

			// insert entry into db
			_, err := stmt.Exec(
				fixtureOdds.ID,
				odds.Market,
				strings.Split(odds.ID, ":")[3],
				odds.SelectionLine,
				odds.Points,
				odds.Sportsbook,
				odds.Price,
				odds.DeepLink.Desktop,
				groupingKey,
			)
			if err != nil {
				s.logger.Error("Failed to execute statement: %v", err)
			}
			processed++
		}
	}

	s.db.Commit()
	s.logger.Info("Processed %d odds in %v", processed, time.Since(start))
}

func (s *oddsService) seedOddsProcess(wg *sync.WaitGroup, jobs chan batchedJob, results chan fixtureOdds) {
	defer wg.Done()
	for job := range jobs {
		fixtureOdds := s.Odds(job.fixtureID, job.sportsbook)
		for _, fixtureOdd := range fixtureOdds {
			results <- fixtureOdd
		}
	}
}

func (s *oddsService) clean(table string) {
	// replace -0.0
	if table == "all_odds" {
		query := `
			UPDATE all_odds
			SET grouping_key = REPLACE(grouping_key, '-0.00', '0.00')
			WHERE grouping_key LIKE '%-0.00%';
		`
		s.db.Exec(query)
	} else if table == "fair_odds" {
		query := `
			UPDATE fair_odds
			SET grouping_key = REPLACE(grouping_key, '-0.00', '0.00')
			WHERE grouping_key LIKE '%-0.00%';
		`
		s.db.Exec(query)
	}
}
