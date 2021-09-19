package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/gin-gonic/gin"
)

func computeQuantity(order *models.TakeProfitOrder) (string, error) {
	// Get the user's USDT balance
	usdtBalance, err := GetBalanceHelper(order.User)
	if err != nil {
		return "", err
	}

	balance, err := strconv.ParseFloat(usdtBalance.Balance, 64)
	if err != nil {
		return "", err
	}

	price, err := strconv.ParseFloat(order.Price, 64)
	if err != nil {
		return "", err
	}

	percentage, err := strconv.ParseFloat(order.Percentage, 64)
	if err != nil {
		return "", err
	}

	// Return error if it is an impossible quantity
	quantity := percentage * balance * price

	return strconv.FormatFloat(quantity, 'f', 5, 32), nil
}

// CreateOrder creates the futures order for the user
func CreateOrder(c *gin.Context) {
	var order models.TakeProfitOrder

	err := c.BindJSON(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	quantity, err := computeQuantity(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	client := futures.NewClient(order.User.APIKey, order.User.APISecret)
	res, err := client.NewCreateOrderService().
		Price(order.Price).
		StopPrice(order.StopPrice).
		Quantity(quantity).
		Side(futures.SideTypeBuy).
		Type(futures.OrderTypeTakeProfit).
		Symbol(order.Symbol).
		Do(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, res)
}
