package v1

import (
	user "github.com/bosdhill/golang-binance-service/controllers/v1/user"
	"github.com/gin-gonic/gin"
)

func SetUserRoutes(rg *gin.RouterGroup) {
	rg.GET("user/ping", user.Ping)
}
