package binancewrapper

import (
	"context"
	"os"
	"testing"

	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/test"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.InitializeBinanceTests()
}

func TestChangeSymbolLeverage(t *testing.T) {
	user := models.User{
		APIKey:    os.Getenv("FUTURES_API_KEY"),
		APISecret: os.Getenv("FUTURES_API_SECRET"),
	}

	ctx := context.Background()
	client := NewClient(&user)

	tests := []struct {
		name     string
		symbol   string
		expected bool
	}{
		{
			name:     "change BTCUSDT leverage",
			symbol:   "BTCUSDT",
			expected: true,
		},
		{
			name:     "don't change BTCUSDT leverage",
			symbol:   "BTCUSDT",
			expected: false,
		},
		{
			name:     "change ETHUSDT leverage",
			symbol:   "ETHUSDT",
			expected: true,
		},
		{
			name:     "don't change ETHUSDT leverage",
			symbol:   "ETHUSDT",
			expected: false,
		},
	}

	// change all symbols to the binance default 20x leverage (i know this is
	// redundant)
	for _, tc := range tests {
		client.changeLeverage(ctx, tc.symbol, 20)
	}

	for _, tc := range tests {
		// get current positions
		account, err := client.GetAccount(ctx)
		if err != nil {
			t.Fatal(err)
		}

		changed, err := client.changeSymbolLeverage(ctx, tc.symbol, account.Positions)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.expected, changed, tc.name)
	}
}
