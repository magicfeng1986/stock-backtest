package strategies

import (
	"fmt"
	"math"
	"stock-backtest/backend/models"
)

// BollingerStrategy 布林带策略
type BollingerStrategy struct{}

// Name 策略名称
func (s *BollingerStrategy) Name() string {
	return "布林带策略"
}

// Description 策略描述
func (s *BollingerStrategy) Description() string {
	return "基于布林带上下轨的突破信号进行交易，价格触及下轨买入，触及上轨卖出"
}

// Execute 执行策略
func (s *BollingerStrategy) Execute(prices []models.StockPrice, initialCap float64, params map[string]interface{}) *models.BacktestResult {
	// 获取参数
	period := 20
	stdDev := 2.0

	if v, ok := params["period"].(float64); ok {
		period = int(v)
	}
	if v, ok := params["stdDev"].(float64); ok {
		stdDev = v
	}

	if len(prices) < period+1 {
		return &models.BacktestResult{
			Strategy:       s.Name(),
			InitialCapital: initialCap,
			FinalCapital:   initialCap,
			TotalReturn:    0,
		}
	}

	// 计算布林带
	middle, upper, lower := calculateBollingerBands(prices, period, stdDev)

	// 回测逻辑
	capital := initialCap
	position := 0
	var trades []models.Trade
	var equityCurve []models.EquityPoint
	var dailyReturns []models.DailyReturn

	maxCapital := initialCap
	maxDrawdown := 0.0

	for i := period; i < len(prices); i++ {
		currentPrice := prices[i].Close
		currentDate := prices[i].Date.Format("2006-01-02")

		currentEquity := capital + float64(position)*currentPrice

		equityCurve = append(equityCurve, models.EquityPoint{
			Date:  currentDate,
			Value: currentEquity,
		})

		if i > period {
			prevEquity := equityCurve[len(equityCurve)-2].Value
			dailyReturn := (currentEquity - prevEquity) / prevEquity
			dailyReturns = append(dailyReturns, models.DailyReturn{
				Date:   currentDate,
				Return: dailyReturn,
				Value:  currentEquity,
			})
		}

		if currentEquity > maxCapital {
			maxCapital = currentEquity
		}
		drawdown := (maxCapital - currentEquity) / maxCapital
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

		// 价格触及下轨买入
		if currentPrice <= lower[i] && position == 0 {
			shares := int(capital / currentPrice)
			if shares > 0 {
				cost := float64(shares) * currentPrice
				capital -= cost
				position = shares
				trades = append(trades, models.Trade{
					Date:   currentDate,
					Type:   "BUY",
					Price:  currentPrice,
					Shares: shares,
					Amount: cost,
					Reason: fmt.Sprintf("触及下轨: 价格=%.2f, 下轨=%.2f", currentPrice, lower[i]),
				})
			}
		}

		// 价格触及上轨卖出
		if currentPrice >= upper[i] && position > 0 {
			revenue := float64(position) * currentPrice
			trades = append(trades, models.Trade{
				Date:   currentDate,
				Type:   "SELL",
				Price:  currentPrice,
				Shares: position,
				Amount: revenue,
				Reason: fmt.Sprintf("触及上轨: 价格=%.2f, 上轨=%.2f", currentPrice, upper[i]),
			})
			capital += revenue
			position = 0
		}
	}

	finalCapital := capital + float64(position)*prices[len(prices)-1].Close
	totalReturn := (finalCapital - initialCap) / initialCap

	days := len(prices) - period
	years := float64(days) / 252.0
	annualizedReturn := 0.0
	if years > 0 {
		annualizedReturn = math.Pow(1+totalReturn, 1/years) - 1
	}

	sharpeRatio := calculateSharpeRatio(dailyReturns)
	winRate := calculateWinRate(trades)

	return &models.BacktestResult{
		Strategy:         s.Name(),
		StartDate:        prices[period].Date.Format("2006-01-02"),
		EndDate:          prices[len(prices)-1].Date.Format("2006-01-02"),
		InitialCapital:   initialCap,
		FinalCapital:     finalCapital,
		TotalReturn:      totalReturn,
		AnnualizedReturn: annualizedReturn,
		MaxDrawdown:      maxDrawdown,
		SharpeRatio:      sharpeRatio,
		WinRate:          winRate,
		TradeCount:       len(trades),
		Trades:           trades,
		DailyReturns:     dailyReturns,
		EquityCurve:      equityCurve,
	}
}

// calculateBollingerBands 计算布林带
func calculateBollingerBands(prices []models.StockPrice, period int, stdDev float64) (middle, upper, lower []float64) {
	middle = make([]float64, len(prices))
	upper = make([]float64, len(prices))
	lower = make([]float64, len(prices))

	for i := period - 1; i < len(prices); i++ {
		// 计算中轨 (SMA)
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i-j].Close
		}
		sma := sum / float64(period)
		middle[i] = sma

		// 计算标准差
		variance := 0.0
		for j := 0; j < period; j++ {
			diff := prices[i-j].Close - sma
			variance += diff * diff
		}
		std := math.Sqrt(variance / float64(period))

		// 计算上下轨
		upper[i] = sma + stdDev*std
		lower[i] = sma - stdDev*std
	}

	return middle, upper, lower
}
