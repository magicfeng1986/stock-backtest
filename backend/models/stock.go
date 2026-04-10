package models

import "time"

// Stock 股票信息
type Stock struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Market string `json:"market"`
}

// StockPrice 股票价格数据
type StockPrice struct {
	Date     time.Time `json:"date"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
	AdjClose float64   `json:"adjClose"`
}

// BacktestRequest 回测请求
type BacktestRequest struct {
	StockCode  string  `json:"stockCode" binding:"required"`
	StartDate  string  `json:"startDate" binding:"required"`
	EndDate    string  `json:"endDate" binding:"required"`
	Strategy   string  `json:"strategy" binding:"required"`
	InitialCap float64 `json:"initialCap"`
	Parameters map[string]interface{} `json:"parameters"`
}

// BacktestResult 回测结果
type BacktestResult struct {
	StockCode       string         `json:"stockCode"`
	StockName       string         `json:"stockName"`
	Strategy        string         `json:"strategy"`
	StartDate       string         `json:"startDate"`
	EndDate         string         `json:"endDate"`
	InitialCapital  float64        `json:"initialCapital"`
	FinalCapital    float64        `json:"finalCapital"`
	TotalReturn     float64        `json:"totalReturn"`
	AnnualizedReturn float64       `json:"annualizedReturn"`
	MaxDrawdown     float64        `json:"maxDrawdown"`
	SharpeRatio     float64        `json:"sharpeRatio"`
	WinRate         float64        `json:"winRate"`
	TradeCount      int            `json:"tradeCount"`
	Trades          []Trade        `json:"trades"`
	DailyReturns    []DailyReturn  `json:"dailyReturns"`
	EquityCurve     []EquityPoint  `json:"equityCurve"`
}

// Trade 交易记录
type Trade struct {
	Date      string  `json:"date"`
	Type      string  `json:"type"` // BUY or SELL
	Price     float64 `json:"price"`
	Shares    int     `json:"shares"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

// DailyReturn 每日收益
type DailyReturn struct {
	Date   string  `json:"date"`
	Return float64 `json:"return"`
	Value  float64 `json:"value"`
}

// EquityPoint 权益曲线点
type EquityPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// StrategyInfo 策略信息
type StrategyInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  []StrategyParameter    `json:"parameters"`
}

// StrategyParameter 策略参数
type StrategyParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Min         float64     `json:"min,omitempty"`
	Max         float64     `json:"max,omitempty"`
	Description string      `json:"description"`
}
