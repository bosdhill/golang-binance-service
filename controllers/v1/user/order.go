package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	binance "github.com/bosdhill/golang-binance-service/libs/binancewrapper"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var timeout = 1 * time.Minute

// CancelAllOrders closes all futures orders for the user
func CancelAllOrders(c *gin.Context) {
	// TODO
}

// CreateOrder creates the futures order for the user. The order types are:
// MARKET, LIMIT, and STOP_MARKET
func CreateOrder(c *gin.Context) {
	var bot models.Bot

	err := c.BindJSON(&bot)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		log.Error(err)
		return
	}

	log.WithFields(log.Fields{
		"Side":  bot.Order.Side,
		"Order": fmt.Sprintf("%#v\n", bot.Order),
	}).Info("New order")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	client := binance.NewClient(&bot.User)
	defer cancel()

	orderResp, err := client.CreateOrder(ctx, bot.Order)
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
		"Side":          bot.Order.Side,
		"Symbol":        orderResp.Symbol,
		"ClientOrderID": orderResp.ClientOrderID,
		"OrigQuantity":  orderResp.OrigQuantity,
	}).Info("Created order")

	c.JSON(http.StatusOK, orderResp)
}
