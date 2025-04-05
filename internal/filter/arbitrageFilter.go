package filter

import (
	"sync"

	"OddsEye/internal/repository"
	"OddsEye/internal/util"

	"OddsEye/pkg/database"
	"OddsEye/pkg/logger"
	"OddsEye/pkg/wagermath"
)

type arbitrageFilter struct {
	db     database.Database
	repo   repository.Repository
	logger logger.Logger
}

func NewArbitrageFilter(db database.Database, repo repository.Repository, logger logger.Logger) *arbitrageFilter {
	return &arbitrageFilter{
		db:     db,
		repo:   repo,
		logger: logger,
	}
}

type groupingCombo struct {
	id          string
	market      string
	groupingKey string
}

type selectionCombo struct {
	group      groupingCombo
	selections map[string]float64
}

func (f *arbitrageFilter) groupingCombos() []groupingCombo {
	query := "SELECT DISTINCT id, market, grouping_key FROM all_odds"
	rows, err := f.db.Query(query)
	if err != nil {
		f.logger.Error("Failed to retrieve distinct fixture odds: %v", err)
	}
	defer rows.Close()

	var combos []groupingCombo
	var id, market, groupingKey string
	for rows.Next() {
		err := rows.Scan(&id, &market, &groupingKey)
		if err != nil {
			f.logger.Error("Failed to scan row: %v", err)
		}
		combos = append(combos, groupingCombo{id, market, groupingKey})
	}
	return combos
}

func (f *arbitrageFilter) selectionCombos(wg *sync.WaitGroup, jobs chan groupingCombo, results chan selectionCombo) {
	defer wg.Done()

	query := "SELECT selection, price FROM all_odds WHERE id = ? AND market = ? AND grouping_key = ?"
	var selection string
	var price float64

	for job := range jobs {
		combos := make(map[string]float64)

		rows, err := f.db.Query(query, job.id, job.market, job.groupingKey)
		if err != nil {
			f.logger.Error("Failed to retrieve distinct fixture odds: %v", err)
		}

		for rows.Next() {
			err := rows.Scan(&selection, &price)
			if err != nil {
				f.logger.Error("Failed to scan row: %v", err)
			}
			if maxv, ok := combos[selection]; !ok || price > maxv {
				combos[selection] = price
			}
		}
		results <- selectionCombo{job, util.MapCopy(combos)}
		rows.Close()
	}
}

func (f *arbitrageFilter) Execute() {
	jobs := make(chan groupingCombo, numWorkers)
	results := make(chan selectionCombo, numWorkers)

	util.LaunchWorkers(numWorkers, jobs, results, f.selectionCombos)

	groupingComobos := f.groupingCombos()
	util.DistributeJobs(groupingComobos, jobs)

	q1 := "INSERT INTO arbitrage (id, market, selection_α, selection_β, price_α, price_β, total_implied_probability, vig) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	q2 := "INSERT INTO arbitrage (id, market, selection_α, selection_β, selection_γ, price_α, price_β, price_γ, total_implied_probability, vig) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	for res := range results {
		f.logger.Info("%+v", res.selections)
		if len(res.selections) == 2 {
			keys := util.Keys(res.selections)
			x := res.selections[keys[0]]
			y := res.selections[keys[1]]
			totalImpliedProb := wagermath.TotalImpliedProbability(x, y)
			vig := totalImpliedProb - 1
			f.db.Exec(q1, res.group.id, res.group.market, keys[0], keys[1], x, y, totalImpliedProb, vig)
		} else if len(res.selections) == 3 {
			keys := util.Keys(res.selections)
			x := res.selections[keys[0]]
			y := res.selections[keys[1]]
			z := res.selections[keys[2]]
			totalImpliedProb := wagermath.TotalImpliedProbability(x, y, z)
			vig := totalImpliedProb - 1
			f.db.Exec(q2, res.group.id, res.group.market, keys[0], keys[1], keys[2], x, y, z, totalImpliedProb, vig)
		}
	}
}
