package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"stock-backtest/backend/models"
	"stock-backtest/backend/services"
	"stock-backtest/backend/strategies"
)

// Handler HTTP处理器
type Handler struct {
	dataService *services.DataService
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{
		dataService: services.NewDataService(),
	}
}

// SearchStocks 搜索股票
func (h *Handler) SearchStocks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "查询参数不能为空"})
		return
	}

	stocks := h.dataService.SearchStocks(query)
	c.JSON(http.StatusOK, gin.H{
		"data": stocks,
	})
}

// GetStrategies 获取所有策略
func (h *Handler) GetStrategies(c *gin.Context) {
	strategies := strategies.GetAllStrategies()
	c.JSON(http.StatusOK, gin.H{
		"data": strategies,
	})
}

// RunBacktest 运行回测
func (h *Handler) RunBacktest(c *gin.Context) {
	var req models.BacktestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认初始资金
	if req.InitialCap <= 0 {
		req.InitialCap = 100000
	}

	// 获取股票数据
	prices, err := h.dataService.GetStockData(req.StockCode, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取股票数据失败: " + err.Error()})
		return
	}

	if len(prices) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到股票数据"})
		return
	}

	// 获取策略并执行
	strategy := strategies.GetStrategy(req.Strategy)
	result := strategy.Execute(prices, req.InitialCap, req.Parameters)

	// 补充股票信息
	result.StockCode = req.StockCode
	result.StockName = h.dataService.GetStockName(req.StockCode)

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// GetStockData 获取股票数据
func (h *Handler) GetStockData(c *gin.Context) {
	stockCode := c.Query("code")
	startDate := c.Query("start")
	endDate := c.Query("end")

	if stockCode == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不完整"})
		return
	}

	prices, err := h.dataService.GetStockData(stockCode, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": prices,
	})
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "股票回测系统运行正常",
	})
}
