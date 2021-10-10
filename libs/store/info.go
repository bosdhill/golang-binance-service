package store

import (
	"context"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/libs/logwrapper"
	log "github.com/sirupsen/logrus"
)

var (
	e           *exchangeInfoStore
	infoOnce    sync.Once
	infoLogFile = "exchangeInfo.log"
)

type exchangeInfoStore struct {
	info map[string]futures.Symbol
	m    sync.RWMutex
	// updateInterval time.Duration
	symbols []string
	l       *logwrapper.StandardLogger
}

// NewInfo returns a reference to the in memory exchangeInfo store
func NewInfo() *exchangeInfoStore {
	infoOnce.Do(func() {
		e = &exchangeInfoStore{}
		e.init()
	})
	return e
}

func (e *exchangeInfoStore) init() {
	e.l = logwrapper.New().WithLogFile(infoLogFile)
	e.fetchExchangeInfo()
}

// GetBaseAssetPrecision returns the base asset precision for a futures symbol.
// Used in the quantity calculation in the futures order.
func (e *exchangeInfoStore) GetBaseAssetPrecision(symbol string) int {
	// e.m.RLock()
	// defer e.m.RUnlock()
	s := e.info[symbol]
	return s.BaseAssetPrecision
}

// GetQuantityPrecision returns the quantity precision for a futures symbol.
func (e *exchangeInfoStore) GetQuantityPrecision(symbol string) int {
	// e.m.RLock()
	// defer e.m.RUnlock()
	s := e.info[symbol]
	return s.QuantityPrecision
}

// GetQuotePrecision returns the quote precision for a futures symbol.
func (e *exchangeInfoStore) GetQuotePrecision(symbol string) int {
	// e.m.RLock()
	// defer e.m.RUnlock()
	s := e.info[symbol]
	return s.QuotePrecision
}

// GetPricePrecision returns the price precision for a futures symbol.
func (e *exchangeInfoStore) GetPricePrecision(symbol string) int {
	// e.m.RLock()
	// defer e.m.RUnlock()
	s := e.info[symbol]
	return s.PricePrecision
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
