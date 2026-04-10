package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"stock-backtest/backend/handlers"
)

func main() {
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 创建处理器
	h := handlers.NewHandler()

	// API路由
	api := r.Group("/api")
	{
		api.GET("/health", h.HealthCheck)
		api.GET("/stocks/search", h.SearchStocks)
		api.GET("/strategies", h.GetStrategies)
		api.POST("/backtest", h.RunBacktest)
		api.GET("/stock/data", h.GetStockData)
	}

	// 静态文件服务
	r.Static("/static", "./frontend")
	r.StaticFile("/", "./frontend/index.html")

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
