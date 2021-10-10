package test

import (
	"github.com/adshao/go-binance/v2/futures"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// InitializeBinanceTests configures the binance tests
func InitializeBinanceTests() {
	gin.SetMode(gin.TestMode)

	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	futures.UseTestnet = true

	// Report the caller method in the logs
	log.SetReportCaller(true)

	log.SetFormatter(&log.JSONFormatter{})
}

// InitializeStoreTests configures the store tests
func IntializeStoreTests() {
	futures.UseTestnet = true

	// Report the caller method in the logs
	log.SetReportCaller(true)

	// // Only log the Debug severity or above
	// log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

// // NewTextCtx returns a context, futures client, and cancel function used in tests
// func NewTestCtx(user *models.User) (context.Context, *futures.Client, context.CancelFunc) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
// 	client := futures.NewClient(user.APIKey, user.APISecret)
// 	return ctx, client, cancel
// }
