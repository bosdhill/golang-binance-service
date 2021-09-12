package main

import (
	v1 "github.com/bosdhill/golang-binance-service/routers/v1"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func init() {
	version1 := router.Group("/v1")
	v1.InitRoutes(version1)
}

func main() {
	router.Run(":5000")
}
