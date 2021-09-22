package user

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/gin-gonic/gin"
)

// GetBalanceHelper returns the User's USD-(s)M Futures Balance
func GetBalanceHelper(user models.User) (*futures.Balance, error) {
	// 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// USD-(s)M Futures
	client := futures.NewClient(user.APIKey, user.APISecret)
	defer client.HTTPClient.CloseIdleConnections()
	balance, err := client.NewGetBalanceService().Do(ctx)
	if err != nil {
		return nil, err
	}

	// Only return the USDT balance
	for _, b := range balance {
		if b.Asset == "USDT" {
			return b, nil
		}
	}
	return nil, errors.New("no USDT balance")
}

// GetBalance returns the users balance based on the User's APIKey and APISecret
func GetBalance(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	res, err := GetBalanceHelper(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		// c.JSON(http.StatusOK, gin.H{"balances": res})
		c.JSON(http.StatusOK, res)
	}
}
