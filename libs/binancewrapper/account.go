package binancewrapper

import (
	"context"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
)

// binanceClient is a wrapper for the binance api
type binanceClient struct {
	c *futures.Client
}

// NewClient returns a new binance client
func NewClient(user *models.User) *binanceClient {
	client := futures.NewClient(user.APIKey, user.APISecret)
	return &binanceClient{client}
}

// GetAccount returns the User's USD-(s)M Futures Account
func (b *binanceClient) GetAccount(ctx context.Context) (*futures.Account, error) {
	account, err := b.c.NewGetAccountService().Do(ctx)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetBalance returns the User's USD-(s)M Futures Balance
func (b *binanceClient) GetBalance(ctx context.Context) (*futures.Balance, error) {
	balance, err := b.c.NewGetBalanceService().Do(ctx)
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
