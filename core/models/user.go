package models

import "github.com/adshao/go-binance/v2/futures"

// User represents the telegram bot user and is identified by their binance
// futures api credentials
type User struct {
	// APIKey is the user's futures api key
	APIKey string `json:"api_key"`

	// APISecret is the user's futures api secret
	APISecret string `json:"api_secret"`
}

// Order represents the Limit/Take Profit, Market, or Stop Loss orders
// presented in the trading signal
type Order struct {
	// Type of order:
	// 	MARKET requires percentage
	//	LIMIT
	//	STOP_LOSS (SL) (or in the Binance API STOP_MARKET)
	Type futures.OrderType `json:"type"`

	// Symbol of the asset
	Symbol string `json:"symbol"`

	// Side or either buy or sell
	Side futures.SideType `json:"side"`

	// Used by LIMIT and MARKET
	// Percentage of futures balance to trade. Used by LIMIT and MARKET orders.
	Percentage string `json:"percentage"`

	// Used by LIMIT
	// Price to buy underlying asset. Used by LIMIT orders.
	Price string `json:"price"`

	// TimeInForce
	// 	GTC - Good Till Cancel
	// 	IOC - Immediate or Cancel
	// 	FOK - Fill or Kill
	// 	GTX - Good Till Crossing (Post Only)
	TimeInForce futures.TimeInForceType `json:"timeInForce"`

	// Used by STOP_MARKET
	// StopPrice closes the position at the market price
	StopPrice string `json:"stopPrice"`
}

// Bot represents a binance futures order
type Bot struct {
	// User's api key and secret
	User User
	// User's Order
	Order Order
}
