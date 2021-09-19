package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	v1 "github.com/bosdhill/golang-binance-service/routers/v1"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	router = gin.Default()
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	useTestnet, err := strconv.ParseBool(os.Getenv("USE_TEST_NET"))
	if err != nil {
		log.Fatal(err.Error())
	}

	binance.UseTestnet = useTestnet
	futures.UseTestnet = useTestnet
	delivery.UseTestnet = useTestnet
	version1 := router.Group("/v1")
	v1.InitRoutes(version1)
}

func main() {
	port := os.Getenv("PORT")
	router.Run(fmt.Sprintf(":%v", port))
}
