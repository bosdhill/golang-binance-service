package v1

import (
	user "github.com/bosdhill/golang-binance-service/controllers/v1/user"
	"github.com/bosdhill/golang-binance-service/middleware"
	"github.com/gin-gonic/gin"
)

func SetUserRoutes(rg *gin.RouterGroup) {
	rg.GET("user/ping", user.Ping, gin.Logger())
	rg.GET("user/balance", user.GetBalance, gin.Logger(), middleware.Validator)
	rg.GET("user/account", user.GetAccount, gin.Logger(), middleware.Validator)
	rg.POST("user/order", user.CreateOrder, gin.Logger(), middleware.Validator)
}
