package wagermath

func ImpliedProbability(price float64) float64 {
	return 1 / price
}

func TotalImpliedProbability(prices ...float64) float64 {
	var totalImpliedProb float64 = 0
	for _, price := range prices {
		totalImpliedProb += ImpliedProbability(price)
	}
	return totalImpliedProb
}

func ExpectedValue(fairPrice float64, price float64) float64 {
	prob := ImpliedProbability(fairPrice)
	return prob*(price-1) - (1 - prob)
}
