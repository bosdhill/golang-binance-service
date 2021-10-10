package binancewrapper

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/store"
	log "github.com/sirupsen/logrus"
)

// CloseAllPositions close all current long position if side is SELL, otherwise
// closes all current short positions if side is BUY.
func (b *binanceClient) CloseAllPositions(
	ctx context.Context,
	symbol string,
	side futures.SideType,
) (*futures.CreateOrderResponse, error) {
	res, err := b.c.NewCreateOrderService().
		Type(futures.OrderTypeStopMarket).
		Symbol(symbol).
		Side(side).
		ClosePosition(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, err
}

// CancelMultipleOrders cancels multiple open orders for a specified symbol.
func (b *binanceClient) CancelMultipleOrders(
	ctx context.Context,
	symbol string,
	orderIDs []int64,
	clientOrderIDs []string,
) ([]*futures.CancelOrderResponse, error) {
	res, err := b.c.NewCancelMultipleOrdersService().
		OrderIDList(orderIDs).
		OrigClientOrderIDList(clientOrderIDs).
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelAllOrders cancels all open futures orders for a specified symbol.
func (b *binanceClient) CancelAllOrders(ctx context.Context, symbol string) error {
	return b.c.NewCancelAllOpenOrdersService().Symbol(symbol).Do(ctx)
}

// CreateOrder creates a futures order.
func (b *binanceClient) CreateOrder(
	ctx context.Context,
	order models.Order,
) (*futures.CreateOrderResponse, error) {
	orderSvc := b.c.NewCreateOrderService()

	orderSvc.Type(order.Type).
		Symbol(order.Symbol).
		Side(order.Side)

	var quantity string
	var err error
	switch order.Type {
	case futures.OrderTypeMarket:
		quantity, err = b.calculateMarketQuantity(ctx, &order)
		if err != nil {
			return nil, err
		}

		orderSvc.Quantity(quantity)

		log.WithFields(log.Fields{
			"Symbol":     order.Symbol,
			"Side":       order.Side,
			"Quantity":   quantity,
			"Percentage": order.Percentage,
		}).Info("New Market Order")
	case futures.OrderTypeLimit:
		quantity, err = b.calculateLimitQuantity(ctx, &order)
		if err != nil {
			return nil, err
		}

		orderSvc.Quantity(quantity).
			Price(order.Price).
			TimeInForce(order.TimeInForce)

		log.WithFields(log.Fields{
			"Symbol":     order.Symbol,
			"Side":       order.Side,
			"Quantity":   quantity,
			"Price":      order.Price,
			"Percentage": order.Percentage,
		}).Info("New Limit Order")
	case futures.OrderTypeStopMarket:
		orderSvc.StopPrice(order.StopPrice)

		log.WithFields(log.Fields{
			"Symbol":    order.Symbol,
			"Side":      order.Side,
			"StopPrice": order.StopPrice,
		}).Info("New Stop Market Order")
	}

	orderResp, err := orderSvc.Do(ctx)
	return orderResp, err
}

// calculateQuantity returns the quantity for a given size, symbol, and price.
func calculateQuantity(size float64, symbol,
	orderPrice string) (string, error) {
	price, err := strconv.ParseFloat(orderPrice, 64)
	if err != nil {
		return "", err
	}

	quantity := size / price
	precision := store.NewInfo().GetQuantityPrecision(symbol)
	return strconv.FormatFloat(quantity, 'f', precision, 64), nil
}

// calculateMarketQuantity returns the quantity of a market order.
func (b *binanceClient) calculateMarketQuantity(ctx context.Context,
	order *models.Order) (string, error) {
	size, err := b.calculatePositionSize(ctx, order.Symbol, order.Percentage)
	if err != nil {
		return "", err
	}

	lastPrice := store.NewStats().GetLastPrice(order.Symbol)

	log.WithFields(log.Fields{
		"Last Price": lastPrice,
	}).Info("Last market price")

	quantity, err := calculateQuantity(size, order.Symbol, lastPrice)
	if err != nil {
		return "", err
	}

	return quantity, nil
}

// calculateLimitQuantity returns the quantity of a limit order.
func (b *binanceClient) calculateLimitQuantity(ctx context.Context,
	order *models.Order) (string, error) {
	size, err := b.calculatePositionSize(ctx, order.Symbol, order.Percentage)
	if err != nil {
		return "", err
	}

	quantity, err := calculateQuantity(size, order.Symbol, order.Price)
	if err != nil {
		return "", err
	}

	return quantity, nil
}

// calculatePositionSize returns the user's position size. The position size
// is calculated using Order.Size * usdtBalance * 10 (at 10x leverage). So if
// Order.Size is 0.10, then 0.10 * 10 * usdtBalance = usdtBalance position is
// opened for the user, with a margin cost of 0.10 * usdtBalance. The risk would
// be 1/leverage or 1/10 in this case.
func (b *binanceClient) calculatePositionSize(ctx context.Context,
	symbol string, percentage float64) (float64, error) {
	account, err := b.GetAccount(ctx)
	if err != nil {
		return 0.0, err
	}

	var usdtBalance float64
	for _, asset := range account.Assets {
		if asset.Asset == "USDT" {
			usdtBalance, err = strconv.ParseFloat(asset.WalletBalance, 64)
			if err != nil {
				return 0.0, err
			}
		}
	}

	// Since the default leverage for symbols is 20x, we might need to update
	// the symbol leverage
	_, err = b.checkSymbolLeverage(ctx, symbol, account.Positions)
	if err != nil {
		return 0.0, err
	}

	positionSize := percentage * usdtBalance * float64(defaultLeverage)

	if positionSize == 0.0 || positionSize > usdtBalance*float64(defaultLeverage) {
		return positionSize, errors.NewPositionSizeInvalid()
	}

	log.WithFields(log.Fields{
		"Balance":      usdtBalance,
		"PositionSize": positionSize,
	}).Info("Calculated position size")

	return positionSize, nil
}
