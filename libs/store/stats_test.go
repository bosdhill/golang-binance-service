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

func TestStats(t *testing.T) {
	symbols := []string{"ETHUSDT", "BTCUSDT"}
	stats := NewStats()

	var lastPrice string
	var newLastPrice string
	for _, symbol := range symbols {
		lastPrice = stats.GetLastPrice(symbol)
		time.Sleep(30 * time.Second)
		newLastPrice = stats.GetLastPrice(symbol)
		assert.NotEqual(t, lastPrice, newLastPrice)
	}
}
