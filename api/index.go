package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"stock-backtest/backend/handlers"
)

// Handler Vercel serverless handler
func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// 创建处理器
	h := handlers.NewHandler()

	// API路由
	api := router.Group("/api")
	{
		api.GET("/health", h.HealthCheck)
		api.GET("/stocks/search", h.SearchStocks)
		api.GET("/strategies", h.GetStrategies)
		api.POST("/backtest", h.RunBacktest)
		api.GET("/stock/data", h.GetStockData)
	}

	// 静态文件服务 - 在Vercel上通过配置处理
	router.Static("/static", "./frontend")
	router.StaticFile("/", "./frontend/index.html")

	router.ServeHTTP(w, r)
}
