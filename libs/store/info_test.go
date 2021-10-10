package store

import (
	"testing"

	"github.com/bosdhill/golang-binance-service/libs/test"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.IntializeStoreTests()
}

func TestInfo(t *testing.T) {
	exchangeInfo := NewInfo()

	tests := []struct {
		name          string
		symbol        string
		expectedPrice int
		expectedQty   int
		expectedBase  int
		expectedQuote int
	}{
		// Note: These values may change. Based off the latest exchangeInfo
		// response
		{
			name:          "btcusdt exchange info",
			symbol:        "BTCUSDT",
			expectedPrice: 2,
			expectedQty:   3,
			expectedBase:  8,
			expectedQuote: 8,
		},
		{
			name:          "ethusdt exchange info",
			symbol:        "ETHUSDT",
			expectedPrice: 2,
			expectedQty:   3,
			expectedBase:  8,
			expectedQuote: 8,
		},
		{
			name:          "trxusdt exchange info",
			symbol:        "TRXUSDT",
			expectedPrice: 5,
			expectedQty:   0,
			expectedBase:  8,
			expectedQuote: 8,
		},
	}

	var actual int
	for _, tc := range tests {
		actual = exchangeInfo.GetPricePrecision(tc.symbol)
		assert.Equal(t, tc.expectedPrice, actual, tc.name)

		actual = exchangeInfo.GetQuantityPrecision(tc.symbol)
		assert.Equal(t, tc.expectedQty, actual, tc.name)

		actual = exchangeInfo.GetBaseAssetPrecision(tc.symbol)
		assert.Equal(t, tc.expectedBase, actual, tc.name)

		actual = exchangeInfo.GetQuotePrecision(tc.symbol)
		assert.Equal(t, tc.expectedQuote, actual, tc.name)
	}
}
