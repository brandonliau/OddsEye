package service

import (
	"OddsEye/pkg/wagermath"
	"sync"
)

type arbitrageJob struct {
	id       string
	market   string
	grouping string
}

type oddsEntry struct {
	selection  string
	price      float64
	sportsbook string
}

type arbitrageEntry struct {
	id               string
	market           string
	odds             []oddsEntry
	totalImpliedProb float64
}

func (s *oddsService) filter() {
	// clear entries that cannot be used for arbitrage
	query := `
		DELETE FROM all_odds_arbitrage
		WHERE (id, market, grouping_key, sportsbook) IN (
			SELECT 
				id,
				market,
				grouping_key,
				sportsbook
			FROM 
				all_odds_arbitrage
			GROUP BY 
				id, 
				market, 
				sportsbook,
				grouping_key
			HAVING
				COUNT(*) NOT IN (2, 3)
		)
	`
	s.db.Exec(query)
}

func (s *oddsService) ScanArbitrage(numWorkers int) {
	// create jobs and results channel
	jobs := make(chan arbitrageJob, 10000)
	results := make(chan arbitrageEntry, 10000)

	// launch worker goroutines
	launchWorkers(numWorkers, jobs, results, s.seedArbitrageProcess)

	// distribute jobs to workers
	go func() {
		rows, _ := s.db.Query("SELECT DISTINCT id, market, grouping_key FROM all_odds_arbitrage")
		defer rows.Close()

		var id, market, grouping string
		for rows.Next() {
			rows.Scan(&id, &market, &grouping)
			jobs <- arbitrageJob{id, market, grouping}
		}

		close(jobs)
	}()

	query := `INSERT INTO arbitrage (id, market, selection_α, selection_β, price_α, price_β, sportsbook_α, sportsbook_β, total_implied_probability, vig)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	two_way_stmt, _ := s.db.PrepareExec(query)

	query = `INSERT INTO arbitrage (id, market, selection_α, selection_β, selection_γ, price_α, price_β, price_γ, sportsbook_α, sportsbook_β, sportsbook_γ, total_implied_probability, vig)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	three_way_stmt, _ := s.db.PrepareExec(query)

	for entries := range results {
		id := entries.id
		market := entries.market
		a := entries.odds[0]
		b := entries.odds[1]
		totalImpliedProb := entries.totalImpliedProb
		vig := entries.totalImpliedProb - 1
		if len(entries.odds) == 2 {
			two_way_stmt.Exec(id, market, a.selection, b.selection, a.price, b.price, a.sportsbook, b.sportsbook, totalImpliedProb, vig)
		} else if len(entries.odds) == 3 {
			c := entries.odds[2]
			three_way_stmt.Exec(id, market, a.selection, b.selection, c.selection, a.price, b.price, c.price, a.sportsbook, b.sportsbook, c.sportsbook, totalImpliedProb, vig)
		}
	}
}

func generateCombinations(data map[string][]oddsEntry) [][]oddsEntry {
	// Extract keys
	keys := make([]string, 0, len(data))
	totalCombinations := 1
	for k := range data {
		keys = append(keys, k)
		totalCombinations *= len(data[k]) // Compute product of lengths
	}

	// Preallocate the result
	result := make([][]oddsEntry, 0, totalCombinations)

	// Preallocate current combination slice
	current := make([]oddsEntry, len(keys))

	var combine func(depth int)
	combine = func(depth int) {
		if depth == len(keys) {
			// Reached a full combination
			combination := make([]oddsEntry, len(current))
			copy(combination, current)
			result = append(result, combination)
			return
		}

		key := keys[depth]
		values := data[key]
		for _, val := range values {
			// Assign directly by index, no append
			current[depth] = val
			combine(depth + 1)
		}
	}

	// Start the recursion
	combine(0)

	return result
}

func (s *oddsService) seedArbitrageProcess(wg *sync.WaitGroup, jobs chan arbitrageJob, results chan arbitrageEntry) {
	defer wg.Done()
	stmt, _ := s.db.PrepareQuery("SELECT selection, sportsbook, price FROM all_odds_arbitrage WHERE id = ? AND market = ? AND grouping_key = ?")
	for job := range jobs {

		data := make(map[string][]oddsEntry)
		rows, _ := stmt.Query(job.id, job.market, job.grouping)

		var selection, sportsbook string
		var price float64
		for rows.Next() {
			err := rows.Scan(&selection, &sportsbook, &price)
			if err != nil {
				s.logger.Error("%v", err)
			}
			if data[selection] == nil {
				data[selection] = make([]oddsEntry, 0)
			}
			data[selection] = append(data[selection], oddsEntry{selection, price, sportsbook})
		}
		rows.Close()

		// generate combinations
		combinations := generateCombinations(data)

		// check for arbitrage combinations
		var totalImpliedProb float64
		for _, combo := range combinations {
			if len(combo) == 2 {
				totalImpliedProb = wagermath.TotalImpliedProbability(combo[0].price, combo[1].price)
			} else if len(combo) == 3 {
				totalImpliedProb = wagermath.TotalImpliedProbability(combo[0].price, combo[1].price, combo[2].price)
			} else {
				s.logger.Warn("%v", combo)
				s.logger.Error("error combo of length %d", len(combo))
			}

			if totalImpliedProb < 1.0 {
				results <- arbitrageEntry{job.id, job.market, combo, totalImpliedProb}
			}
		}
	}
}
