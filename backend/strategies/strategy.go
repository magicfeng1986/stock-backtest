package strategies

import "stock-backtest/backend/models"

// Strategy 策略接口
type Strategy interface {
	Name() string
	Description() string
	Execute(prices []models.StockPrice, initialCap float64, params map[string]interface{}) *models.BacktestResult
}

// GetAllStrategies 获取所有可用策略
func GetAllStrategies() []models.StrategyInfo {
	return []models.StrategyInfo{
		{
			ID:          "sma",
			Name:        "均线策略 (SMA Crossover)",
			Description: "基于短期和长期移动平均线的交叉信号进行买卖",
			Parameters: []models.StrategyParameter{
				{Name: "shortPeriod", Type: "number", Default: 10, Min: 5, Max: 30, Description: "短期均线周期"},
				{Name: "longPeriod", Type: "number", Default: 30, Min: 20, Max: 100, Description: "长期均线周期"},
			},
		},
		{
			ID:          "macd",
			Name:        "MACD策略",
			Description: "基于MACD指标的买卖信号进行交易",
			Parameters: []models.StrategyParameter{
				{Name: "fastPeriod", Type: "number", Default: 12, Min: 5, Max: 20, Description: "快速EMA周期"},
				{Name: "slowPeriod", Type: "number", Default: 26, Min: 15, Max: 40, Description: "慢速EMA周期"},
				{Name: "signalPeriod", Type: "number", Default: 9, Min: 5, Max: 15, Description: "信号线周期"},
			},
		},
		{
			ID:          "rsi",
			Name:        "RSI策略",
			Description: "基于RSI超买超卖信号进行交易",
			Parameters: []models.StrategyParameter{
				{Name: "period", Type: "number", Default: 14, Min: 7, Max: 30, Description: "RSI计算周期"},
				{Name: "overbought", Type: "number", Default: 70, Min: 60, Max: 90, Description: "超买阈值"},
				{Name: "oversold", Type: "number", Default: 30, Min: 10, Max: 40, Description: "超卖阈值"},
			},
		},
		{
			ID:          "bollinger",
			Name:        "布林带策略",
			Description: "基于布林带上下轨的突破信号进行交易",
			Parameters: []models.StrategyParameter{
				{Name: "period", Type: "number", Default: 20, Min: 10, Max: 50, Description: "布林带周期"},
				{Name: "stdDev", Type: "number", Default: 2.0, Min: 1.0, Max: 3.0, Description: "标准差倍数"},
			},
		},
		{
			ID:          "momentum",
			Name:        "动量策略",
			Description: "基于价格动量的变化进行交易",
			Parameters: []models.StrategyParameter{
				{Name: "period", Type: "number", Default: 10, Min: 5, Max: 30, Description: "动量计算周期"},
				{Name: "threshold", Type: "number", Default: 0.02, Min: 0.01, Max: 0.1, Description: "动量阈值"},
			},
		},
	}
}

// GetStrategy 根据ID获取策略
func GetStrategy(id string) Strategy {
	switch id {
	case "sma":
		return &SMAStrategy{}
	case "macd":
		return &MACDStrategy{}
	case "rsi":
		return &RSIStrategy{}
	case "bollinger":
		return &BollingerStrategy{}
	case "momentum":
		return &MomentumStrategy{}
	default:
		return &SMAStrategy{} // 默认使用均线策略
	}
}
