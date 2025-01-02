package wagermath

import (
	"math"
)

const (
	MaxIterations = 100
	ConvergenceThreshold = 1e-12
)

/* REFERENCES */
// Adjusting Bookmaker’s Odds to Allow for Overround: https://www.sciencepublishinggroup.com/article/10.11648/j.ajss.20170506.12
// Prices of State Contingent Claims with Insider Traders, and the Favourite-Longshot Bias: https://doi.org/10.2307/2234526
// Measuring the Incidence of Insider Trading in a Market for State-Contingent Claims: https://doi.org/10.2307/2234240
// Beating the market with a bad predictive model: https://arxiv.org/pdf/2010.12508
// A Comment on the Bias of Probabilities Derived From Betting Odds and Their Use in Measuring Outcome Uncertainty: https://doi.org/10.1177/1527002513519329
// A Family of Solutions Related to Shin’s Model For Probability Forecasts: https://www.cambridge.org/engage/coe/article-details/666672b2e7ccf7753a661ce7

func RemoveVigMultiplicative(prices ...float64) []float64 {
	novig := make([]float64, len(prices))
	booksum := TotalImpliedProbability(prices...)
	for i, price := range prices {
		novig[i] = 1 / (ImpliedProbability(price) / booksum)
	}
	return novig
}

func RemoveVigAdditive(prices ...float64) []float64 {
	novig := make([]float64, len(prices))
	overround := TotalImpliedProbability(prices...) - 1
	count := float64(len(prices))
	for i, price := range prices {
		novig[i] = 1 / (ImpliedProbability(price) - (overround / count))
	}
	return novig
}

func RemoveVigPower(prices ...float64) []float64 {
	novig := append([]float64{}, prices...)
	count := float64(len(prices))
	booksum := 0.0
	delta := math.Inf(1)
	iterations := 0

	for delta > ConvergenceThreshold && iterations < MaxIterations {
		prevBooksum := booksum
		booksum = TotalImpliedProbability(novig...)
		exponent := math.Log(count) / math.Log(count / booksum)
		for i := range(int(count)) {
			prob := ImpliedProbability(novig[i])
			novig[i] = 1 / (math.Pow(prob, exponent))
		}
		delta = math.Abs(booksum - prevBooksum)
		iterations++
	}
	return novig
}

func RemoveVigShin(prices ...float64) []float64 {
	if len(prices) == 2 {
		return RemoveVigAdditive(prices...)
	}

	novig := append([]float64{}, prices...)
	booksum := TotalImpliedProbability(prices...)
	delta := math.Inf(1)
	iterations := 0
	z := 0.0
	count := float64(len(prices))

	for delta > ConvergenceThreshold && iterations < MaxIterations {
		prevZ := z
		zSum := 0.0
		for _, price := range novig {
			prob := ImpliedProbability(price)
			term := math.Sqrt(math.Pow(z, 2) + 4.0 * (1.0 - z) * (math.Pow(prob, 2) / booksum))
			zSum += term
		}
		z = (zSum - 2.0) / (count - 2)
		delta = math.Abs(prevZ - z)
		iterations++
	}

	for i, price := range novig {
		prob := ImpliedProbability(price)
		term := math.Sqrt(math.Pow(z, 2) + 4.0 * (1.0 - z) * (math.Pow(prob, 2) / booksum)) - z
		novig[i] = 1 / (term / (2.0 * (1.0 - z)))
	}

	return novig
}

func RemoveVigWorstCase(prices ...float64) []float64 {
	multiplicative := RemoveVigMultiplicative(prices...)
	additive := RemoveVigAdditive(prices...)
	power := RemoveVigPower(prices...)
	shin := RemoveVigShin(prices...)

	var multiplicativeSum, additiveSum, powerSum, shinSum float64
	for i := range(prices) {
		multiplicativeSum += multiplicative[i]
		additiveSum += additive[i]
		powerSum += power[i]
		shinSum += shin[i]
	}
	
	worstCase := multiplicative
	minv := multiplicativeSum
	if additiveSum < minv {
		worstCase = additive
		minv = additiveSum
	}
	if powerSum < minv {
		worstCase = power
		minv = powerSum
	}
	if shinSum < minv {
		worstCase = shin
	}

	return worstCase
}
