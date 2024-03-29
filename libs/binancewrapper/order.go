// Package binancewrapper wraps the binance api client
package binancewrapper

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/binancewrapper/retry"
	"github.com/bosdhill/golang-binance-service/libs/store/info"
	"github.com/bosdhill/golang-binance-service/libs/store/stats"
	log "github.com/sirupsen/logrus"
)

// CloseAllPositions will create a STOP_MARKET order that will be triggered when
// the stopPrice is met with closePosition=true. If triggered, it will close all
// open long (BUY) positions if the side is SELL, otherwise it will close all
// open short (SELL) positions if the side is BUY.
//
// TODO: This doesn't guarantee all positions would be closed. In order to close a
// position, a market order must be placed in the opposite direction for the
// same quantity. OR this can be managed with all the frontend server bookkeeping
// (api keys, order info, position info)
func (b *binanceClient) CloseAllPositions(
	ctx context.Context,
	symbol string,
	side futures.SideType,
	stopPrice string,
) (*futures.CreateOrderResponse, error) {
	svc := b.c.NewCreateOrderService().
		Type(futures.OrderTypeStopMarket).
		Symbol(symbol).
		Side(side).
		StopPrice(stopPrice).
		ClosePosition(true)
	var res *futures.CreateOrderResponse
	res, err := svc.Do(ctx)
	if err != nil {
		retryRes, err := retry.Do(err, func(opts ...futures.RequestOption) (interface{}, error) {
			log.WithField("recvWindow", opts).Info("Retrying CloseAllPositions request")
			return svc.Do(ctx, opts...)
		})
		if err != nil {
			return nil, err
		}
		res = retryRes.(*futures.CreateOrderResponse)
	}

	log.WithFields(log.Fields{
		"Symbol":    symbol,
		"Side":      side,
		"StopPrice": stopPrice,
	}).Info("New Stop Market Order")

	if err != nil {
		return nil, err
	}
	return res, err
}

// CancelMultipleOrders cancels multiple open orders for a specified symbol.
// Note that once an order is filled, it becomes an open position and can not
// be cancelled.
func (b *binanceClient) CancelMultipleOrders(
	ctx context.Context,
	symbol string,
	orderIDs []int64,
	clientOrderIDs []string,
) ([]*futures.CancelOrderResponse, error) {
	svc := b.c.NewCancelMultipleOrdersService().
		OrderIDList(orderIDs).
		OrigClientOrderIDList(clientOrderIDs).
		Symbol(symbol)
	var res []*futures.CancelOrderResponse
	res, err := svc.Do(ctx)
	if err != nil {
		retryRes, err := retry.Do(err, func(opts ...futures.RequestOption) (interface{}, error) {
			log.WithField("recvWindow", opts).Info("Retrying CancelMultipleOrders request")
			return svc.Do(ctx, opts...)
		})
		if err != nil {
			return nil, err
		}
		res = retryRes.([]*futures.CancelOrderResponse)
	}

	log.WithFields(log.Fields{
		"Symbol":         symbol,
		"OrderIDs":       orderIDs,
		"ClientOrderIDs": clientOrderIDs,
	}).Info("New Cancel Multiple Orders")

	if err != nil {
		return nil, err
	}
	return res, nil
}

// CancelAllOrders cancels all open futures orders for a specified symbol.
func (b *binanceClient) CancelAllOrders(ctx context.Context, symbol string) error {
	svc := b.c.NewCancelAllOpenOrdersService().Symbol(symbol)
	err := svc.Do(ctx)
	if err != nil {
		_, err := retry.Do(err, func(opts ...futures.RequestOption) (interface{}, error) {
			log.WithField("recvWindow", opts).Info("Retrying CancelAllOrders request")
			return nil, svc.Do(ctx, opts...)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateOrder creates a futures order.
func (b *binanceClient) CreateOrder(
	ctx context.Context,
	order *models.Order,
) (*futures.CreateOrderResponse, error) {
	svc := b.c.NewCreateOrderService()

	svc.Type(order.Type).
		Symbol(order.Symbol).
		Side(order.Side)

	var quantity string
	var err error
	switch order.Type {
	case futures.OrderTypeMarket:
		quantity, err = b.calculateMarketQuantity(ctx, order)
		if err != nil {
			return nil, err
		}

		svc.Quantity(quantity)

		log.WithFields(log.Fields{
			"Symbol":     order.Symbol,
			"Side":       order.Side,
			"Quantity":   quantity,
			"Percentage": order.Percentage,
		}).Info("New Market Order")
	case futures.OrderTypeLimit:
		quantity, err = b.calculateLimitQuantity(ctx, order)
		if err != nil {
			return nil, err
		}

		svc.Quantity(quantity).
			Price(order.Price).
			TimeInForce(order.TimeInForce)

		log.WithFields(log.Fields{
			"Symbol":      order.Symbol,
			"Side":        order.Side,
			"Quantity":    quantity,
			"Price":       order.Price,
			"Percentage":  order.Percentage,
			"TimeInForce": order.TimeInForce,
		}).Info("New Limit Order")
	case futures.OrderTypeStopMarket:
		quantity, err = b.calculateStopMarketQuantity(ctx, order)
		if err != nil {
			return nil, err
		}

		svc.StopPrice(order.StopPrice).
			Quantity(quantity).
			TimeInForce(order.TimeInForce)

		log.WithFields(log.Fields{
			"Symbol":      order.Symbol,
			"Side":        order.Side,
			"Quantity":    quantity,
			"StopPrice":   order.StopPrice,
			"Percentage":  order.Percentage,
			"TimeInForce": order.TimeInForce,
		}).Info("New Market Order")
	}

	var res *futures.CreateOrderResponse
	res, err = svc.Do(ctx)
	if err != nil {
		retryRes, err := retry.Do(err, func(opts ...futures.RequestOption) (interface{}, error) {
			log.WithField("recvWindow", opts).Info("Retrying CreateOrder request")
			return svc.Do(ctx, opts...)
		})
		if err != nil {
			return nil, err
		}
		res = retryRes.(*futures.CreateOrderResponse)
	}
	return res, err
}

// calculateMarketQuantity returns the quantity of a market order.
func (b *binanceClient) calculateMarketQuantity(ctx context.Context,
	order *models.Order) (string, error) {
	return b.calculate(
		ctx,
		order,
		func(size float64) (string, error) {
			lastPrice := stats.NewStore().GetLastPrice(order.Symbol)
			return calculateQuantity(size, order.Symbol, lastPrice)
		},
	)
}

// calculateStopMarketQuantity returns the quantity of a stop market order
// (similar to calculateLimitQuantity except it uses StopPrice).
func (b *binanceClient) calculateStopMarketQuantity(
	ctx context.Context,
	order *models.Order,
) (string, error) {
	return b.calculate(
		ctx,
		order,
		func(size float64) (string, error) {
			return calculateQuantity(size, order.Symbol, order.StopPrice)
		},
	)
}

// calculateLimitQuantity returns the quantity of a limit order.
func (b *binanceClient) calculateLimitQuantity(ctx context.Context,
	order *models.Order) (string, error) {
	return b.calculate(
		ctx,
		order,
		func(size float64) (string, error) {
			return calculateQuantity(size, order.Symbol, order.Price)
		},
	)
}

// calculateQuantity returns the quantity for a given size, symbol, and price.
func calculateQuantity(size float64, symbol, orderPrice string) (string, error) {
	price, err := strconv.ParseFloat(orderPrice, 64)
	if err != nil {
		return "", err
	}

	quantity := size / price
	precision := info.NewStore().GetQuantityPrecision(symbol)
	return strconv.FormatFloat(quantity, 'f', precision, 64), nil
}

// calcFunc is the function used for calculating the quantity given the size.
type calcFunc func(size float64) (string, error)

// calculate returns the quantity by first calculating the position size for the
// symbol and then calcFunc to calculate the order quantity.
func (b *binanceClient) calculate(
	ctx context.Context,
	order *models.Order,
	calcQuantity calcFunc,
) (string, error) {
	size, err := b.calculatePositionSize(ctx, order.Symbol, order.Percentage)
	if err != nil {
		return "", err
	}
	return calcQuantity(size)
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
	_, err = b.changeSymbolLeverage(ctx, symbol, account.Positions)
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
