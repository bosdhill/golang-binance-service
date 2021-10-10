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

// GetAccount returns the users futures account based on the User's APIKey and APISecret
func GetAccount(c *gin.Context) {
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

	res, err := client.GetAccount(ctx)
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
		// "Assets":                      fmt.Sprintf("%v", res.Assets),
		"CanDeposit":        res.CanDeposit,
		"CanTrade":          res.CanTrade,
		"CanWithdraw":       res.CanWithdraw,
		"FeeTier":           res.FeeTier,
		"MaxWithdrawAmount": res.MaxWithdrawAmount,
		// "Positions":                   fmt.Sprintf("%v", res.Positions),
		"TotalInitialMargin":          res.TotalInitialMargin,
		"TotalMaintMargin":            res.TotalMaintMargin,
		"TotalMarginBalance":          res.TotalMarginBalance,
		"TotalOpenOrderInitialMargin": res.TotalOpenOrderInitialMargin,
		"TotalPositionInitialMargin":  res.TotalPositionInitialMargin,
		"TotalUnrealizedProfit":       res.TotalUnrealizedProfit,
		"TotalWalletBalance":          res.TotalWalletBalance,
		"UpdateTime":                  res.UpdateTime,
	}).Info("Got Account")

	c.JSON(http.StatusOK, res)
}
