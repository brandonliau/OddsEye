package wagermath

import (
	"fmt"
	"testing"
)

func TestRemoveVigPower(t *testing.T) {
	fmt.Printf("Multiplicative: %-9f\n", RemoveVigMultiplicative(1.05, 10.0, 18.5, 19.0))
	fmt.Printf("Additive:       %-9f\n", RemoveVigAdditive(1.05, 10.0, 18.5, 19.0))
	fmt.Printf("Power:          %-9f\n", RemoveVigPower(1.05, 10.0, 18.5, 19.0))
	fmt.Printf("Shin:           %-9f\n", RemoveVigShin(1.05, 10.0, 18.53, 19.0))

	fmt.Println()

	fmt.Printf("Multiplicative: %-8f\n", RemoveVigMultiplicative(1.909, 1.833))
	fmt.Printf("Additive:       %-8f\n", RemoveVigAdditive(1.909, 1.833))
	fmt.Printf("Power:          %-8f\n", RemoveVigPower(1.909, 1.833))
	fmt.Printf("Shin:           %-8f\n", RemoveVigShin(1.909, 1.833))

	mult := RemoveVigMultiplicative(1.909, 1.833)
	add := RemoveVigAdditive(1.909, 1.833)
	pow := RemoveVigPower(1.909, 1.833)
	shin := RemoveVigShin(1.909, 1.833)
	fmt.Printf("Worst Case: 	%-8f\n", RemoveVigWorstCase(mult, add, pow, shin))
}
