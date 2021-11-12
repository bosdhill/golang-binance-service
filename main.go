package main

import (
	"fmt"
	"os"
	"strconv"

	binance "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/bosdhill/golang-binance-service/libs/store/info"
	"github.com/bosdhill/golang-binance-service/libs/store/stats"
	v1 "github.com/bosdhill/golang-binance-service/routers/v1"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type ServerCtx struct {
	Port       string
	UseTestnet bool
	Debug      bool
}

var (
	router      = gin.Default()
	defaultPort = "4200"
)

func loadServerCtx() *ServerCtx {
	s := &ServerCtx{defaultPort, false, false}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port != "" {
		s.Port = port
	}

	debugMode := os.Getenv("DEBUG")
	s.Debug, err = strconv.ParseBool(debugMode)
	if err != nil {
		log.Fatal(err)
	}

	useTestnet, err := strconv.ParseBool(os.Getenv("USE_TESTNET"))
	if err != nil {
		log.Fatal(err)
	}

	s.UseTestnet = useTestnet
	binance.UseTestnet = useTestnet
	futures.UseTestnet = useTestnet
	delivery.UseTestnet = useTestnet

	log.WithFields(log.Fields{
		"Port":       s.Port,
		"UseTestnet": s.UseTestnet,
		"Debug":      s.Debug,
	}).Info("Server configuration loaded")

	return s
}

func init() {
	version1 := router.Group("/v1")
	v1.InitRoutes(version1)

	// Report the caller method in the logs
	log.SetReportCaller(true)

	// Output to stdout instead of stderr
	log.SetOutput(os.Stdout)

	// Only log the Info severity or above
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	s := loadServerCtx()
	if s.Debug {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	}

	// Create in memory store to maintain price stats
	stats.NewStore().StartUpdates()

	// Create in memory store for exchange info
	info.NewStore().StartUpdates()

	router.Run(fmt.Sprintf(":%v", s.Port))
}
