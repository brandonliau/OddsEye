package wagermath

import (
	"math"
)

const (
	MaxIterations        = 100
	ConvergenceThreshold = 1e-12
)

/* REFERENCES */
// Adjusting Bookmaker’s Odds to Allow for Overround: https://www.sciencepublishinggroup.com/article/10.11648/j.ajss.20170506.12

/* IMPLEMENTATION REFERENCES */
// Python/Rust implementation of Shin's method: https://github.com/mberk/shin
// R implementation of Shin's method: https://github.com/opisthokonta/implied

func RemoveVigMultiplicative(prices ...float64) []float64 {
	novig := make([]float64, len(prices))
	booksum := TotalImpliedProbability(prices...)
	for i, price := range prices {
		novig[i] = price * booksum
	}
	return novig
}

func RemoveVigAdditive(prices ...float64) []float64 {
	novig := make([]float64, len(prices))
	share := (TotalImpliedProbability(prices...) - 1.0) / float64(len(prices))
	for i, price := range prices {
		novig[i] = price / (1 - share*price)
	}
	return novig
}

func RemoveVigPower(prices ...float64) []float64 {
	novig := make([]float64, len(prices))
	copy(novig, prices)
	count := float64(len(novig))
	booksum := 0.0
	delta := math.Inf(1)
	iterations := 0
	for delta > ConvergenceThreshold && iterations < MaxIterations {
		prevBooksum := booksum
		booksum = TotalImpliedProbability(novig...)
		k := math.Log(count) / math.Log(count/booksum)
		for i := range int(count) {
			prob := ImpliedProbability(novig[i])
			novig[i] = 1 / (math.Pow(prob, k))
		}
		delta = math.Abs(booksum - prevBooksum)
		iterations++
	}
	return novig
}

func RemoveVigShin(prices ...float64) []float64 {
	novig := make([]float64, len(prices))
	copy(novig, prices)
	count := float64(len(novig))
	booksum := TotalImpliedProbability(novig...)
	delta := math.Inf(1)
	iterations := 0
	z := 0.0
	if len(novig) == 2 {
		diff := ImpliedProbability(novig[0]) - ImpliedProbability(novig[1])
		z = ((booksum - 1) * (diff*diff - booksum)) / (booksum * (diff*diff - 1))
	} else {
		for delta > ConvergenceThreshold && iterations < MaxIterations {
			prevZ := z
			zSum := 0.0
			for _, price := range novig {
				prob := ImpliedProbability(price)
				term := math.Sqrt((z * z) + 4.0*(1.0-z)*prob*prob/booksum)
				zSum += term
			}
			z = (zSum - 2.0) / (count - 2)
			delta = math.Abs(prevZ - z)
			iterations++
		}
	}
	for i, price := range novig {
		prob := ImpliedProbability(price)
		term := math.Sqrt((z*z)+4.0*(1.0-z)*prob*prob/booksum) - z
		novig[i] = 2.0 * (1.0 - z) / term
	}
	return novig
}

func RemoveVigWorstCase(methods ...[]float64) []float64 {
	minv := make([]float64, len(methods[0]))
	copy(minv, methods[0])

	for _, vals := range methods {
		for i, val := range vals {
			minv[i] = min(minv[i], val)
		}
	}
	return minv
}
