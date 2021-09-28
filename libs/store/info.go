package store

import (
	"context"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
)

var (
	e           *exchangeInfoStore
	newInfoOnce sync.Once
)

type exchangeInfoStore struct {
	info map[string]futures.Symbol
	m    sync.RWMutex
	// updateInterval time.Duration
	symbols []string
}

// NewInfo returns a reference to the in memory latest price store
func NewInfo() *exchangeInfoStore {
	newInfoOnce.Do(func() {
		e = &exchangeInfoStore{}
		e.fetchExchangeInfo()
	})
	return e
}

// GetBaseAssetPrecision returns the base asset precision for a futures symbol.
// Used in the quantity calculation in the futures order.
func (e *exchangeInfoStore) GetBaseAssetPrecision(symbol string) int {
	e.m.RLock()
	defer e.m.RUnlock()
	s := e.info[symbol]
	return s.BaseAssetPrecision
}

func (e *exchangeInfoStore) fetchExchangeInfo() {
	exchangeInfo, err := futures.NewClient("", "").
		NewExchangeInfoService().
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return
	}

	e.info = make(map[string]futures.Symbol)
	for _, s := range exchangeInfo.Symbols {

		e.symbols = append(e.symbols, s.Symbol)
		// We only care about BaseAssetPrecision (used in new order quantity calc)
		// from binance:
		// base asset refers to the asset that is the quantity of a symbol.
		// quote asset refers to the asset that is the price of a symbol.
		e.info[s.Symbol] = s
	}
}
