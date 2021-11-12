// Package binancewrapper wraps the binance api client
package binancewrapper

import (
	"context"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/errors"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/bosdhill/golang-binance-service/libs/binancewrapper/retry"
	log "github.com/sirupsen/logrus"
)

var binanceOnce sync.Once

// binanceClient is a wrapper for the binance api.
type binanceClient struct {
	c *futures.Client
}

// NewClient returns a new binance client.
func NewClient(user *models.User) *binanceClient {
	client := futures.NewClient(user.APIKey, user.APISecret)
	b := binanceClient{client}
	binanceOnce.Do(func() {
		// Should store timeoffset somewhere for future use when new clients are
		// created since this function can be called concurrently
		err := retry.ServerTimeSync()
		if err != nil {
			log.Fatal(err, "could not get binance server time and set time offset")
		}
	})
	return &b
}

// GetAccount returns the User's USD-(s)M Futures Account.
func (b *binanceClient) GetAccount(ctx context.Context) (*futures.Account, error) {
	svc := b.c.NewGetAccountService()
	var res *futures.Account
	res, err := svc.Do(ctx)
	if err != nil {
		retryRes, err := retry.Do(err, func(opts ...futures.RequestOption) (interface{}, error) {
			log.WithField("recvWindow", opts).Info("Retrying GetAccount request")
			return svc.Do(ctx, opts...)
		})
		if err != nil {
			return nil, err
		}
		res = retryRes.(*futures.Account)
	}
	return res, nil
}

// getBalances returns the User's USD-(s)M Futures Balances.
func (b *binanceClient) getBalances(ctx context.Context) ([]*futures.Balance, error) {
	svc := b.c.NewGetBalanceService()
	var res []*futures.Balance
	res, err := svc.Do(ctx)
	if err != nil {
		retryRes, err := retry.Do(err, func(opts ...futures.RequestOption) (interface{}, error) {
			log.WithField("recvWindow", opts).Info("Retrying GetBalance request")
			return svc.Do(ctx, opts...)
		})
		if err != nil {
			return nil, err
		}
		res = retryRes.([]*futures.Balance)
	}
	return res, nil
}

// GetBalance returns the user's futures USDT Balance.
func (b *binanceClient) GetUSDTBalance(ctx context.Context) (*futures.Balance, error) {
	balances, err := b.getBalances(ctx)
	if err != nil {
		return nil, err
	}
	for _, balance := range balances {
		if balance.Asset == "USDT" {
			return balance, nil
		}
	}
	return nil, errors.NewNoUSDTBalance()
}
