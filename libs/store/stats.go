package store

import (
	"context"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
)

var (
	s         *statsStore
	statsOnce sync.Once
)

type Stats struct {
	// Symbol             string
	PriceChange        string
	PriceChangePercent string
	WeightedAvgPrice   string
	// PrevClosePrice     string
	LastPrice    string
	LastQuantity string
	// OpenPrice          string
	// HighPrice          string
	// LowPrice           string
	// Volume             string
	// QuoteVolume        string
	// OpenTime           int64
	// CloseTime          int64
	// FristID            int64
	// LastID             int64
	// Count              int64
	PricePrecision     int
	BaseAssetPrecision int
	QuantityPrecision  int
}

type statsStore struct {
	stats          map[string]Stats
	m              sync.RWMutex
	updateInterval time.Duration
	symbols        []string
}

// NewStats returns a reference to the in memory latest price store
func NewStats() *statsStore {
	statsOnce.Do(func() {
		s = &statsStore{}
		s.init()
	})
	return s
}

// WithInterval is the last price update interval in duration string format
func (c *statsStore) WithInterval(seconds string) *statsStore {
	d, err := time.ParseDuration(seconds)
	if err != nil {
		return nil
	}

	c.updateInterval = d

	log.WithFields(log.Fields{"price update interval": seconds}).Info()

	return c
}

// GetLastPrice gets the last price for a futures symbol
func (c *statsStore) GetLastPrice(symbol string) string {
	c.m.RLock()
	defer c.m.RUnlock()
	s := c.stats[symbol]
	return s.LastPrice
}

// GetSymbols returns the list of futures symbols
func (c *statsStore) GetSymbols() []string {
	return c.symbols
}

func (c *statsStore) init() {
	c.fetchSymbolsAndPriceStats()
	c.startUpdates()
}

func (c *statsStore) fetchSymbolsAndPriceStats() {
	priceStats, err := futures.NewClient("", "").
		NewListPriceChangeStatsService().
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return
	}

	c.stats = make(map[string]Stats)
	for _, priceStat := range priceStats {
		log.Debug(log.Fields{"symbol": priceStat.Symbol,
			"last price": priceStat.LastPrice})

		c.symbols = append(c.symbols, priceStat.Symbol)
		stats := c.stats[priceStat.Symbol]
		stats.LastPrice = priceStat.LastPrice
		stats.LastQuantity = priceStat.LastQuantity
		stats.PriceChange = priceStat.PriceChange
		stats.WeightedAvgPrice = priceStat.WeightedAvgPrice
	}
}

func (c *statsStore) update(events futures.WsAllMarketTickerEvent) {
	for _, event := range events {
		log.WithFields(log.Fields{"symbol": event.Symbol,
			"last price": event.ClosePrice}).
			Debug("updating last price")

		s := c.stats[event.Symbol]

		// We care about LastPrice (used in new order quantity calc)
		s.LastPrice = event.ClosePrice
		s.LastQuantity = event.CloseQty
		s.WeightedAvgPrice = event.WeightedAvgPrice
		s.PriceChange = event.PriceChange
		s.PriceChangePercent = event.PriceChangePercent
	}
}

// startUpdates will start updating last prices ~updateInterval period
func (c *statsStore) startUpdates() {
	eventHandler := func(events futures.WsAllMarketTickerEvent) {
		time.Sleep(c.updateInterval)
		c.m.Lock()
		defer c.m.Unlock()
		c.update(events)
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
