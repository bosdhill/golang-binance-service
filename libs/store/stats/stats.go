// Package store implements in an memory store for binance exchange price stats
package stats

import (
	"context"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/libs/logwrapper"
	log "github.com/sirupsen/logrus"
)

var (
	s               *statsStore
	once            sync.Once
	statsLogFile    = "stats.log"
	defaultInterval = "3s"
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
// updateInterval + 1 seconds (websocket has a default of 1 update every second).
type statsStore struct {
	stats          map[string]Stats
	m              sync.RWMutex
	updateInterval time.Duration
	symbols        []string
	l              *logwrapper.StandardLogger
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
	s.l = logwrapper.New().WithLogFile(statsLogFile)
	s.updateInterval, _ = time.ParseDuration(defaultInterval)
	s.fetchSymbolsAndPriceStats()
}

// WithInterval is the last price update interval in duration string format.
func (s *statsStore) WithInterval(d string) {
	s.updateInterval, _ = time.ParseDuration(d)
	s.l.WithFields(log.Fields{"stats store update interval": d}).Info()
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
		s.l.Fatal(err)
		return
	}

	for _, priceStat := range priceStats {
		s.l.Debug(log.Fields{"symbol": priceStat.Symbol,
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
		s.l.WithFields(log.Fields{"symbol": priceStat.Symbol,
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
func (s *statsStore) StartUpdates() {
	eventHandler := func(events futures.WsAllMarketTickerEvent) {
		time.Sleep(s.updateInterval)
		s.m.Lock()
		defer s.m.Unlock()
		s.update(events)
	}

	errHandler := func(err error) {
		s.l.Trace(err)
	}

	_, _, err := futures.WsAllMarketTickerServe(eventHandler, errHandler)
	if err != nil {
		s.l.Fatal(err)
		return
	}
}
