package v1

import (
	user "github.com/bosdhill/golang-binance-service/controllers/v1/user"
	"github.com/bosdhill/golang-binance-service/middleware"
	"github.com/gin-gonic/gin"
)

func SetUserRoutes(rg *gin.RouterGroup) {
	// TODO: Add validators
	rg.GET("user/ping", user.Ping, gin.Logger())
	rg.GET("user/balance", user.GetBalance, gin.Logger(), middleware.Validator)
	rg.POST("user/trade", user.CreateOrder, gin.Logger(), middleware.Validator)
}
