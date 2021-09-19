package models

// User represents the telegram bot user and is identified by their binance futures api credentials
type User struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// TakeProfitOrder represents a binance futures take profit order
type TakeProfitOrder struct {
	// User's api key and secret
	User User
	// Price to buy asset
	Price string `json:"price"`
	// StopPrice to trigger a market sell
	StopPrice string `json:"stopPrice"`
	// Symbol of the asset
	Symbol string `json:"symbol"`
	// Side or either buy or sell
	Side string `json:"side"`
	// Percentage of futures balance to trade
	Percentage string `json:"percentage"`
}
