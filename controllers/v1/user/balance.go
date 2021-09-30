package user

import (
	"context"
	"net/http"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetUSDTBalanceHelper returns the User's USD-(s)M Futures Balance
func getBalance(
	ctx context.Context,
	c *futures.Client,
	user *models.User,
) (*futures.Balance, error) {
	balance, err := c.NewGetBalanceService().Do(ctx)
	if err != nil {
		return nil, err
	}

	// Only return the USDT balance
	for _, b := range balance {
		if b.Asset == "USDT" {
			return b, nil
		}
	}
	return nil, errors.NewNoUSDTBalance()
}

// GetBalance returns the users balance based on the User's APIKey and APISecret
func GetBalance(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		log.Error(err)
		return
	}

	// 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// USD-(s)M Futures
	client := futures.NewClient(user.APIKey, user.APISecret)
	defer client.HTTPClient.CloseIdleConnections()

	res, err := getBalance(ctx, client, &user)
	if err != nil {
		if common.IsAPIError(err) {
			apiErr := errors.NewAPIError(err)
			c.JSON(int(apiErr.Code), errors.NewAPIError(err))
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
		log.Error(err)
		return
	}

	log.WithFields(log.Fields{
		"AccountAlias":       res.AccountAlias,
		"Asset":              res.Asset,
		"Balance":            res.Balance,
		"CrossWalletBalance": res.CrossWalletBalance,
		"CrossUnPnl":         res.CrossUnPnl,
		"AvailableBalance":   res.AvailableBalance,
		"MaxWithdrawAmount":  res.MaxWithdrawAmount,
	}).Info("Got Balance")

	c.JSON(http.StatusOK, res)
}
