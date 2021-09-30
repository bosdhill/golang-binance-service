package user

import (
	"os"
	"strconv"
	"testing"

	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/test"
	"github.com/stretchr/testify/assert"
)

func init() {
	test.IntializeControllerTests()
}

func TestComputePercentBalance(t *testing.T) {
	percentage := "0.20"
	user := &models.User{
		APIKey:    os.Getenv("FUTURES_API_KEY"),
		APISecret: os.Getenv("FUTURES_API_SECRET"),
	}
	bot := &models.Bot{
		User: *user,
		Order: models.Order{
			Percentage: percentage,
		},
	}

	ctx, client, cancel := test.NewTestCtx(&bot.User)
	defer cancel()

	balance, err := getBalance(ctx, client, user)
	if err != nil {
		t.Fatal(err)
	}

	balFloat, err := strconv.ParseFloat(balance.Balance, 64)
	if err != nil {
		t.Fatal(err)
	}

	percentFloat, err := strconv.ParseFloat(percentage, 64)
	if err != nil {
		t.Fatal(err)
	}

	expectedBal := percentFloat * balFloat * leverageMultiplier
	actual, err := computePercentBalance(ctx, client, bot)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, actual, expectedBal)
}

func TestComputeLimitQuantity(t *testing.T) {
	// user := &models.User{
	// 	APIKey:    os.Getenv("FUTURES_API_KEY"),
	// 	APISecret: os.Getenv("FUTURES_API_SECRET"),
	// }
	// bot := &models.Bot{
	// 	User: *user,
	// 	Order: models.Order{
	// 		Percentage: "0.20",
	// 	},
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	// defer cancel()

	// client := futures.NewClient(user.APIKey, user.APISecret)
	// defer client.HTTPClient.CloseIdleConnections()

	// quantity, computeLimitQuantity(ctx, client, bot)

}

func TestComputeMarketQuantity(t *testing.T) {

}

func TestCreateMarketOrder(t *testing.T) {
	bot := &models.Bot{
		User: models.User{
			APIKey:    os.Getenv("FUTURES_API_KEY"),
			APISecret: os.Getenv("FUTURES_API_SECRET"),
		},
		Order: models.Order{
			Type:        "",
			Symbol:      "",
			Side:        "",
			Percentage:  "",
			Price:       "",
			TimeInForce: "",
			StopPrice:   "",
		},
	}

	ctx, client, cancel := test.NewTestCtx(&bot.User)
	defer cancel()

	createOrder(ctx, client, *bot)
}
