package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/store"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Assuming that user's are trading at 10x cross
// returns 0.0 on and error
var leverageMultiplier = 10.0

// computePercentBalance computes the percentage of the user's max leverage.
func computePercentBalance(
	ctx context.Context,
	c *futures.Client,
	bot *models.Bot,
) (float64, error) {
	usdt, err := getBalance(ctx, c, &bot.User)
	if err != nil {
		return 0.0, err
	}

	balance, err := strconv.ParseFloat(usdt.Balance, 64)
	if err != nil {
		return 0.0, err
	}

	percentage, err := strconv.ParseFloat(bot.Order.Percentage, 64)
	if err != nil {
		return 0.0, err
	}

	log.WithFields(log.Fields{
		"AccountAlias": usdt.AccountAlias,
		"Balance":      balance,
	}).Info()

	return percentage * balance * leverageMultiplier, nil
}

// computeMarketQuantity computes the position size of the market order.
func computeMarketQuantity(ctx context.Context, c *futures.Client,
	bot *models.Bot) (string, error) {
	percentageBal, err := computePercentBalance(ctx, c, bot)
	if err != nil {
		return "", err
	}

	lastPrice, err := strconv.ParseFloat(store.NewStats().
		GetLastPrice(bot.Order.Symbol), 64)
	if err != nil {
		return "", err
	}

	quantity := percentageBal / lastPrice

	log.WithFields(log.Fields{
		"Market Quantity": quantity,
	}).Info()

	precision := store.NewInfo().GetBaseAssetPrecision(bot.Order.Symbol)
	return strconv.FormatFloat(quantity, 'f', precision, 32), nil
}

// compteLimitQuantity computes the position size a limit order.
func computeLimitQuantity(ctx context.Context, c *futures.Client,
	bot *models.Bot) (string, error) {
	percentageBal, err := computePercentBalance(ctx, c, bot)
	if err != nil {
		return "", err
	}

	lastPrice, err := strconv.ParseFloat(bot.Order.Price, 64)
	if err != nil {
		return "", err
	}

	quantity := percentageBal / lastPrice

	log.WithFields(log.Fields{
		"Limit Quantity": quantity,
	}).Info()

	precision := store.NewInfo().GetBaseAssetPrecision(bot.Order.Symbol)
	return strconv.FormatFloat(quantity, 'f', precision, 32), nil
}

// createOrder creates the futures order
func createOrder(
	ctx context.Context,
	c *futures.Client,
	bot models.Bot,
) (*futures.CreateOrderResponse, error) {
	orderSvc := c.NewCreateOrderService()

	orderSvc.Type(bot.Order.Type).
		Symbol(bot.Order.Symbol).
		Side(bot.Order.Side)

	var quantity string
	var err error

	switch bot.Order.Type {
	case futures.OrderTypeMarket:
		quantity, err = computeMarketQuantity(ctx, c, &bot)
		if err != nil {
			return nil, err
		}

		orderSvc.Quantity(quantity)
	case futures.OrderTypeLimit:
		quantity, err = computeLimitQuantity(ctx, c, &bot)
		if err != nil {
			return nil, err
		}

		orderSvc.Quantity(quantity).
			Price(bot.Order.Price).
			TimeInForce(bot.Order.TimeInForce)
	case futures.OrderTypeStopMarket:
		orderSvc.StopPrice(bot.Order.StopPrice)
	}

	orderResp, err := orderSvc.Do(ctx)
	return orderResp, err
}

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

	// 1 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	client := futures.NewClient(bot.User.APIKey, bot.User.APISecret)
	defer client.HTTPClient.CloseIdleConnections()

	orderResp, err := createOrder(ctx, client, bot)
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
