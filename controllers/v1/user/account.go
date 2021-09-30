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

// getAccount returns the User's USD-(s)M Futures Account
func getAccount(
	ctx *context.Context,
	c *futures.Client,
	user models.User,
) (*futures.Account, error) {
	account, err := c.NewGetAccountService().Do(*ctx)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetAccount returns the users futures account based on the User's APIKey and APISecret
func GetAccount(c *gin.Context) {
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

	res, err := getAccount(&ctx, client, user)
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
