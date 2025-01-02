package wagermath

import (
	"math"
	"testing"
)

// Test when fair price equals market price
func TestExpectedValue_FairPriceEqualsMarketPrice(t *testing.T) {
	fairPrice := 2.0
	price := 2.0
	expected := 0.0

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}

// Test when market price is higher than fair price
func TestExpectedValue_MarketPriceHigher(t *testing.T) {
	fairPrice := 2.0
	price := 3.0
	expected := 0.5

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}

// Test when market price is lower than fair price
func TestExpectedValue_MarketPriceLower(t *testing.T) {
	fairPrice := 2.0
	price := 1.5
	expected := -0.25

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}

// Test edge case when fair price is 1 (certainty)
func TestExpectedValue_FairPriceIsOne(t *testing.T) {
	fairPrice := 1.0
	price := 2.0
	expected := 1.0

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}

// Test edge case when market price equals 1
func TestExpectedValue_MarketPriceIsOne(t *testing.T) {
	fairPrice := 2.0
	price := 1.0
	expected := -0.5

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}

// Test edge case with high fair price
func TestExpectedValue_HighFairPrice(t *testing.T) {
	fairPrice := 100.0
	price := 200.0
	expected := 1.0

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}

// Test edge case with very low fair price
func TestExpectedValue_LowFairPrice(t *testing.T) {
	fairPrice := 0.1
	price := 10.0
	expected := 99.0

	result := ExpectedValue(fairPrice, price)
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("ExpectedValue(%v, %v) = %v; want %v", fairPrice, price, result, expected)
	}
}
