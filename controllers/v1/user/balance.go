package user

import (
	"context"
	"net/http"

	"github.com/adshao/go-binance/v2/common"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	binance "github.com/bosdhill/golang-binance-service/libs/binancewrapper"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetBalance returns the users balance based on the User's APIKey and APISecret
func GetBalance(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		log.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	client := binance.NewClient(&user)
	defer cancel()

	res, err := client.GetUSDTBalance(ctx)
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
