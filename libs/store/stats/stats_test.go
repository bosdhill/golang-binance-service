// Package store implements in an memory store for binance exchange price stats
package stats

import (
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
	symbols := []string{"ETHUSDT", "BTCUSDT"}
	stats := NewStore()

	var lastPrice string
	var newLastPrice string
	for _, symbol := range symbols {
		lastPrice = stats.GetLastPrice(symbol)
		time.Sleep(30 * time.Second)
		newLastPrice = stats.GetLastPrice(symbol)
		assert.NotEqual(t, lastPrice, newLastPrice, "last price is not updated")
	}
}
