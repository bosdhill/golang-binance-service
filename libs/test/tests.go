package test

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/core/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func IntializeControllerTests() {
	gin.SetMode(gin.TestMode)

	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	futures.UseTestnet = true

	// Report the caller method in the logs
	log.SetReportCaller(true)
}

func NewTestCtx(user *models.User) (context.Context, *futures.Client, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	client := futures.NewClient(user.APIKey, user.APISecret)
	return ctx, client, cancel
}
