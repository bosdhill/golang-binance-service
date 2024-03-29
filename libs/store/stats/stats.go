// Package store implements in an memory store for binance exchange price stats
package stats

import (
	"context"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
)

var (
	s            *statsStore
	once         sync.Once
	defaultDelay = "0s"
)

// Stats stores various price stats for a futures symbol.
type Stats struct {
	PriceChange        string
	PriceChangePercent string
	WeightedAvgPrice   string
	LastPrice          string
	LastQuantity       string
}

// statsStore stores various price stats for all futures symbols. A single
// websocket receives price stats updates on the entire market.
//
// On each update, the entire stats map is updated, which happens every
// updateDelay + 1 seconds (websocket has a default of 1 update every second).
type statsStore struct {
	stats       map[string]Stats
	m           sync.RWMutex
	updateDelay time.Duration
	symbols     []string
}

// NewStore returns a reference to the in memory latest price store.
func NewStore() *statsStore {
	once.Do(func() {
		s = &statsStore{}
		s.init()
	})
	return s
}

func (s *statsStore) init() {
	s.updateDelay, _ = time.ParseDuration(defaultDelay)
	s.fetchSymbolsAndPriceStats()
	s.startUpdates()
}

// WithDelay is the last price update delay in duration string format.
func (s *statsStore) WithDelay(d string) {
	s.updateDelay, _ = time.ParseDuration(d)
	log.WithFields(log.Fields{"stats store update delay": d}).Info()
}

// GetLastPrice gets the last price for a futures symbol.
func (s *statsStore) GetLastPrice(symbol string) string {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.stats[symbol].LastPrice
}

// GetSymbols returns the list of futures symbols.
func (s *statsStore) GetSymbols() []string {
	return s.symbols
}

// fetchSymbolsAndPriceStats initializes the price stats map.
func (s *statsStore) fetchSymbolsAndPriceStats() {
	s.stats = make(map[string]Stats)

	priceStats, err := futures.NewClient("", "").
		NewListPriceChangeStatsService().
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return
	}

	for _, priceStat := range priceStats {
		log.Debug(log.Fields{"symbol": priceStat.Symbol,
			"last price": priceStat.LastPrice})

		s.symbols = append(s.symbols, priceStat.Symbol)
		s.stats[priceStat.Symbol] = Stats{
			PriceChange:        priceStat.PriceChange,
			PriceChangePercent: priceStat.PriceChangePercent,
			WeightedAvgPrice:   priceStat.WeightedAvgPrice,
			LastPrice:          priceStat.LastPrice,
			LastQuantity:       priceStat.LastQuantity,
		}
	}
}

// update will update the actual map used to store each symbol's price stats
func (s *statsStore) update(events futures.WsAllMarketTickerEvent) {
	for _, priceStat := range events {
		log.WithFields(log.Fields{"symbol": priceStat.Symbol,
			"last price": priceStat.ClosePrice}).
			Debug("updating last price")

		// We care about LastPrice (used in new order quantity calc)
		s.stats[priceStat.Symbol] = Stats{
			PriceChange:        priceStat.PriceChange,
			PriceChangePercent: priceStat.PriceChangePercent,
			WeightedAvgPrice:   priceStat.WeightedAvgPrice,
			LastPrice:          priceStat.ClosePrice,
			LastQuantity:       priceStat.CloseQty,
		}
	}
}

// startUpdates opens the websocket and will start updating the entire
// statsStore every updateInterval + 1 sec.
func (s *statsStore) startUpdates() {
	eventHandler := func(events futures.WsAllMarketTickerEvent) {
		time.Sleep(s.updateDelay)
		s.m.Lock()
		defer s.m.Unlock()
		s.update(events)
	}

	errHandler := func(err error) {
		log.Trace(err)
	}

	_, _, err := futures.WsAllMarketTickerServe(eventHandler, errHandler)
	if err != nil {
		log.Fatal(err)
		return
	}
}
