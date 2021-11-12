// Package store implements in an memory store for binance exchange price stats
package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/bosdhill/golang-binance-service/libs/test"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.IntializeStoreTests()
}

func TestNewStats(t *testing.T) {
	stats := NewStore()
	assert.NotEqual(t, nil, stats, "stats is not nil")
}

func TestGetSymbols(t *testing.T) {
	stats := NewStore()
	s := stats.GetSymbols()
	assert.NotEqual(t, 0, len(s), "number of symbols is not zero")
}

func TestGetLastPrice(t *testing.T) {
	stats := NewStore()
	tests := []struct {
		symbol    string
		lastPrice string
	}{
		{
			symbol:    "ETHUSDT",
			lastPrice: stats.GetLastPrice("ETHUSDT"),
		},
		{
			symbol:    "BTCUSDT",
			lastPrice: stats.GetLastPrice("BTCUSDT"),
		},
	}

	time.Sleep(15 * time.Second)
	for _, tc := range tests {
		newLastPrice := stats.GetLastPrice(tc.symbol)
		assert.NotEqual(
			t,
			tc.lastPrice,
			newLastPrice,
			fmt.Sprintf("last price not updated for %s", tc.symbol),
		)
	}
}
