// Package store implements in memory stores for the binance exchange info and
// price stats
package store

import (
	"context"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/libs/logwrapper"
	log "github.com/sirupsen/logrus"
)

var (
	s            *statsStore
	statsOnce    sync.Once
	statsLogFile = "stats.log"
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

// NewStats returns a reference to the in memory latest price store.
func NewStats() *statsStore {
	statsOnce.Do(func() {
		s = &statsStore{}
		s.init()
	})
	return s
}

func (c *statsStore) init() {
	c.l = logwrapper.New().WithLogFile(statsLogFile)
	c.fetchSymbolsAndPriceStats()
	c.startUpdates()
}

// WithInterval is the last price update interval in duration string format.
func (c *statsStore) WithInterval(seconds string) *statsStore {
	d, err := time.ParseDuration(seconds)
	if err != nil {
		return nil
	}

	c.updateInterval = d

	c.l.WithFields(log.Fields{"price update interval": seconds}).Info()

	return c
}

// GetLastPrice gets the last price for a futures symbol.
func (c *statsStore) GetLastPrice(symbol string) string {
	c.m.RLock()
	defer c.m.RUnlock()
	s := c.stats[symbol]
	return s.LastPrice
}

// GetSymbols returns the list of futures symbols.
func (c *statsStore) GetSymbols() []string {
	return c.symbols
}

// fetchSymbolsAndPriceStats initializes the price stats map.
func (c *statsStore) fetchSymbolsAndPriceStats() {
	c.stats = make(map[string]Stats)

	priceStats, err := futures.NewClient("", "").
		NewListPriceChangeStatsService().
		Do(context.Background())

	if err != nil {
		c.l.Fatal(err)
		return
	}

	for _, priceStat := range priceStats {
		c.l.Debug(log.Fields{"symbol": priceStat.Symbol,
			"last price": priceStat.LastPrice})

		c.symbols = append(c.symbols, priceStat.Symbol)
		c.stats[priceStat.Symbol] = Stats{
			PriceChange:        priceStat.PriceChange,
			PriceChangePercent: priceStat.PriceChangePercent,
			WeightedAvgPrice:   priceStat.WeightedAvgPrice,
			LastPrice:          priceStat.LastPrice,
			LastQuantity:       priceStat.LastQuantity,
		}
	}
}

// update will update the actual map used to store each symbol's price stats
func (c *statsStore) update(events futures.WsAllMarketTickerEvent) {
	for _, priceStat := range events {
		c.l.WithFields(log.Fields{"symbol": priceStat.Symbol,
			"last price": priceStat.ClosePrice}).
			Debug("updating last price")

		// We care about LastPrice (used in new order quantity calc)
		c.stats[priceStat.Symbol] = Stats{
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
func (c *statsStore) startUpdates() {
	eventHandler := func(events futures.WsAllMarketTickerEvent) {
		time.Sleep(c.updateInterval)
		c.m.Lock()
		defer c.m.Unlock()
		c.update(events)
	}

	errHandler := func(err error) {
		c.l.Trace(err)
	}

	_, _, err := futures.WsAllMarketTickerServe(eventHandler, errHandler)
	if err != nil {
		c.l.Fatal(err)
		return
	}
}
