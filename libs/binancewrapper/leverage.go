package binancewrapper

import (
	"context"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
)

var (
	// defaultLeverage is 10x cross
	defaultLeverage = 10
)

// changeLeverage changes the symbol's initial leverage.
func (b *binanceClient) changeLeverage(
	ctx context.Context,
	symbol string,
	leverage int,
) error {
	res, err := b.c.NewChangeLeverageService().
		Leverage(leverage).
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"Symbol":       res.Symbol,
		"New Leverage": res.Leverage,
	}).Info("Changed symbol leverage")

	return nil
}

// checkSymbolLeverage will change the symbol's initial leverage if its not the
// same as the desired symbol leverage. Returns whether or not the leverage was
// changed.
func (b *binanceClient) checkSymbolLeverage(
	ctx context.Context,
	symbol string,
	positions []*futures.AccountPosition,
) (bool, error) {
	changed := false

	currentLeverage, err := getCurrentLeverage(symbol, positions)
	if err != nil {
		return changed, err
	}

	if currentLeverage != defaultLeverage {
		err := b.changeLeverage(ctx, symbol, defaultLeverage)
		if err != nil {
			return changed, err
		}
		changed = true
	}

	return changed, nil
}

// getCurrentLeverage returns the current leverage for a symbol.
func getCurrentLeverage(symbol string, positions []*futures.AccountPosition) (int, error) {
	var currentLeverage int
	var err error
	for _, position := range positions {
		if position.Symbol == symbol {
			currentLeverage, err = strconv.Atoi(position.Leverage)
			if err != nil {
				return 0, err
			}
		}
	}

	return currentLeverage, nil
}
