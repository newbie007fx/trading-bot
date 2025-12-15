package indicator_test

import (
	"testing"

	"github.com/newbie007fx/trading-bot/internal/indicator"
)

func TestEmaSeries(t *testing.T) {
	prices := []float64{100, 102, 101, 105, 107, 110}
	result, err := indicator.EMASeries(prices, 5)
	if err != nil {
		t.Fatal("error should be nill")
	}

	t.Log(result)
}
