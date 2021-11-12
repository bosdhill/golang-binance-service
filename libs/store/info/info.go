// Package info implements an in memory store for binance exchange info
package info

import (
	"context"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
)

var (
	e            *exchangeInfoStore
	once         sync.Once
	defaultDelay = "30m"
)

type exchangeInfoStore struct {
	info        map[string]futures.Symbol
	m           sync.RWMutex
	updateDelay time.Duration
	symbols     []string
}

// NewStore returns a reference to the in memory exchangeInfo store
func NewStore() *exchangeInfoStore {
	once.Do(func() {
		e = &exchangeInfoStore{}
		e.init()
	})
	return e
}

func (e *exchangeInfoStore) init() {
	e.updateDelay, _ = time.ParseDuration(defaultDelay)
	e.fetchExchangeInfo()
	e.startUpdates()
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

// GetBaseAssetPrecision returns the base asset precision for a futures symbol.
// Used in the quantity calculation in the futures order.
func (e *exchangeInfoStore) GetBaseAssetPrecision(symbol string) int {
	e.m.RLock()
	defer e.m.RUnlock()
	s := e.info[symbol]
	return s.BaseAssetPrecision
}

// GetQuantityPrecision returns the quantity precision for a futures symbol.
func (e *exchangeInfoStore) GetQuantityPrecision(symbol string) int {
	e.m.RLock()
	defer e.m.RUnlock()
	s := e.info[symbol]
	return s.QuantityPrecision
}

// GetQuotePrecision returns the quote precision for a futures symbol.
func (e *exchangeInfoStore) GetQuotePrecision(symbol string) int {
	e.m.RLock()
	defer e.m.RUnlock()
	s := e.info[symbol]
	return s.QuotePrecision
}

// GetPricePrecision returns the price precision for a futures symbol.
func (e *exchangeInfoStore) GetPricePrecision(symbol string) int {
	e.m.RLock()
	defer e.m.RUnlock()
	s := e.info[symbol]
	return s.PricePrecision
}

// GetPriceFilter returns a price filter for a symbol
func (e *exchangeInfoStore) GetPriceFilter(symbol string) *futures.PriceFilter {
	e.m.RLock()
	defer e.m.RUnlock()
	s := e.info[symbol]
	return s.PriceFilter()
}

// WithDelay is the last price update delay in duration string format.
func (e *exchangeInfoStore) WithDelay(d string) {
	e.updateDelay, _ = time.ParseDuration(d)
	log.WithFields(log.Fields{"exchange info store update delay": d}).Info()
}

func (e *exchangeInfoStore) update() {
	exchangeInfo, err := futures.NewClient("", "").
		NewExchangeInfoService().
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return
	}

	for _, s := range exchangeInfo.Symbols {
		log.WithFields(log.Fields{"symbol": s.Symbol,
			"info": s}).
			Debug("Updating symbol's exchange info")

		e.info[s.Symbol] = s
	}
}

// startUpdates opens the websocket and will start updating the entire
// statsStore every updateInterval + 1 sec.
func (c *exchangeInfoStore) startUpdates() {
	go func() {
		time.Sleep(c.updateDelay)
		e.m.Lock()
		defer e.m.Unlock()
		c.update()
	}()
}
