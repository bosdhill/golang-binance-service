package store

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
	stats := NewStats()
	assert.NotEqual(t, nil, stats, "stats is not nil")
}

func TestGetSymbols(t *testing.T) {
	stats := NewStats()

	s := stats.GetSymbols()
	assert.NotEqual(t, 0, len(s), "number of symbols is not zero")
}

func TestGetLastPrice(t *testing.T) {
	symbols := []string{"ETHUSDT", "BTCUSDT"}
	stats := NewStats()

	var lastPrice string
	var newLastPrice string
	for _, symbol := range symbols {
		lastPrice = stats.GetLastPrice(symbol)
		time.Sleep(30 * time.Second)
		newLastPrice = stats.GetLastPrice(symbol)
		assert.NotEqual(t, lastPrice, newLastPrice, "last price is updated")
	}
}
