package binancewrapper

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/store"
	"github.com/bosdhill/golang-binance-service/libs/test"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.InitializeBinanceTests()
}

func TestCalculatePositionSize(t *testing.T) {
	user := models.User{
		APIKey:    os.Getenv("FUTURES_API_KEY"),
		APISecret: os.Getenv("FUTURES_API_SECRET"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	client := NewClient(&user)
	defer cancel()

	b, err := client.GetUSDTBalance(ctx)
	if err != nil {
		t.Fatal(err)
	}

	usdtBalance, err := strconv.ParseFloat(b.Balance, 64)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		symbol     string
		percentage float64
		expected   float64
	}{
		{
			name:       "Size: 20%",
			symbol:     "BTCUSDT",
			percentage: 0.20,
			expected:   0.20 * usdtBalance * float64(defaultLeverage),
		},
		{
			name:       "Size: 50%",
			symbol:     "TRXUSDT",
			percentage: 0.50,
			expected:   0.50 * usdtBalance * float64(defaultLeverage),
		},
		{
			name:       "Size: 90%",
			symbol:     "ETHUSDT",
			percentage: 0.90,
			expected:   0.90 * usdtBalance * float64(defaultLeverage),
		},
		{
			name:       "Size: 100%",
			symbol:     "BTCUSDT",
			percentage: 1.0,
			expected:   1.0 * usdtBalance * float64(defaultLeverage),
		},
		// edge cases
		{
			name:       "Size: 0%",
			symbol:     "XRPUSDT",
			percentage: 0.0,
			expected:   0.0 * usdtBalance * float64(defaultLeverage),
		},
		{
			name:       "Size: 110%",
			symbol:     "ETHUSDT",
			percentage: 1.10,
			expected:   1.10 * usdtBalance * float64(defaultLeverage),
		},
	}

	for _, tc := range tests {
		actual, err := client.calculatePositionSize(ctx, tc.symbol, tc.percentage)

		if err != nil {
			assert.EqualError(t, err, errors.NewPositionSizeInvalid().Error(), tc.name)
		}

		assert.Equal(t, tc.expected, actual, tc.name)
	}
}

func TestCalculateLimitQuantity(t *testing.T) {
	user := &models.User{
		APIKey:    os.Getenv("FUTURES_API_KEY"),
		APISecret: os.Getenv("FUTURES_API_SECRET"),
	}

	ctx := context.Background()
	client := NewClient(user)

	b, err := client.GetUSDTBalance(ctx)
	if err != nil {
		t.Fatal(err)
	}

	usdtBalance, err := strconv.ParseFloat(b.Balance, 64)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		expected float64
		order    *models.Order
	}{
		{
			name: "calculate limit order quantity for BTCUSDT at 60000",
			order: &models.Order{
				Type:       futures.OrderTypeLimit,
				Symbol:     "BTCUSDT",
				Percentage: 0.10,
				Price:      "60000.0",
			},
			expected: (0.10 * usdtBalance * float64(defaultLeverage)) / 60000.0,
		},
		{
			name: "calculate limit order quantity for ETHUSDT at 4000.0",
			order: &models.Order{
				Type:       futures.OrderTypeLimit,
				Symbol:     "ETHUSDT",
				Percentage: 0.10,
				Price:      "4000.0",
			},
			expected: (0.10 * usdtBalance * float64(defaultLeverage)) / 4000.0,
		},
		{
			name: "calculate limit order quantity for DOTUSDT at 75.0",
			order: &models.Order{
				Type:       futures.OrderTypeLimit,
				Symbol:     "DOTUSDT",
				Percentage: 0.10,
				Price:      "75.0",
			},
			expected: (0.10 * usdtBalance * float64(defaultLeverage)) / 75.0,
		},
	}

	for _, tc := range tests {
		quantity, err := client.calculateLimitQuantity(ctx, tc.order)
		if err != nil {
			t.Fatal(err)
		}

		precision := store.NewInfo().GetQuantityPrecision(tc.order.Symbol)
		expected := strconv.FormatFloat(tc.expected, 'f', precision, 32)
		assert.Equal(t, expected, quantity, tc.name)
	}
}

func TestCalculateMarketQuantity(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		expected  float64
		size      float64
		lastPrice string
	}{
		{
			name:      "calculate market order quantity for BTCUSDT at lastPrice 54000",
			size:      100000,
			lastPrice: "54000",
			symbol:    "BTCUSDT",
			expected:  100000 / 54000.0,
		},
		{
			name:      "calculate market order quantity for ETHUSDT at lastPrice 3200.0",
			size:      100000,
			lastPrice: "3200.0",
			symbol:    "ETHUSDT",
			expected:  100000 / 3200.0,
		},
		{
			name:      "calculate market order quantity for DOTUSDT at lastPrice 33.0",
			size:      100000,
			lastPrice: "33",
			symbol:    "DOTUSDT",
			expected:  100000 / 33.0,
		},
	}

	for _, tc := range tests {
		quantity, err := calculateQuantity(tc.size, tc.symbol, tc.lastPrice)
		if err != nil {
			t.Fatal(err)
		}

		precision := store.NewInfo().GetQuantityPrecision(tc.symbol)
		expected := strconv.FormatFloat(tc.expected, 'f', precision, 32)
		assert.Equal(t, expected, quantity, tc.name)
	}
}

func TestCreateMarketOrder(t *testing.T) {
	user := &models.User{
		APIKey:    os.Getenv("FUTURES_API_KEY"),
		APISecret: os.Getenv("FUTURES_API_SECRET"),
	}

	ctx := context.Background()
	client := NewClient(user)

	tests := []struct {
		name  string
		order *models.Order
	}{
		{
			name: "create market order of Size 0.01 for BTCUSDT",
			order: &models.Order{
				Type:       futures.OrderTypeMarket,
				Symbol:     "BTCUSDT",
				Side:       futures.SideTypeBuy,
				Percentage: 0.01,
			},
		},
		{
			name: "create market order of Size 0.01 for ETHUSDT",
			order: &models.Order{
				Type:       futures.OrderTypeMarket,
				Symbol:     "ETHUSDT",
				Side:       futures.SideTypeBuy,
				Percentage: 0.01,
			},
		},
	}

	for _, tc := range tests {
		res, err := client.CreateOrder(ctx, *tc.order)
		if err != nil {
			t.Fatal(err)
		}

		got, err := json.MarshalIndent(&res, "", " ")
		if err != nil {
			t.Fatal(err, tc.name)
		}
		t.Logf(string(got))

		// Close all long positions by creating a STOP_MARKET sell order at a
		// very low stop price to guarantee that it will be triggered.
		// Similar to https://github.com/sammchardy/python-binance/issues/536#issuecomment-643743964
		res, err = client.CloseAllPositions(ctx, tc.order.Symbol, futures.SideTypeSell, "10.0")
		if err != nil {
			t.Fatal(err, tc.name)
		}

		got, err = json.MarshalIndent(&res, "", " ")
		if err != nil {
			t.Fatal(err, tc.name)
		}
		t.Logf(string(got))
	}
}
