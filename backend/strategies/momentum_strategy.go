package strategies

import (
	"fmt"
	"math"
	"stock-backtest/backend/models"
)

// MomentumStrategy 动量策略
type MomentumStrategy struct{}

// Name 策略名称
func (s *MomentumStrategy) Name() string {
	return "动量策略"
}

// Description 策略描述
func (s *MomentumStrategy) Description() string {
	return "基于价格动量的变化进行交易，动量突破阈值时买入，低于阈值时卖出"
}

// Execute 执行策略
func (s *MomentumStrategy) Execute(prices []models.StockPrice, initialCap float64, params map[string]interface{}) *models.BacktestResult {
	// 获取参数
	period := 10
	threshold := 0.02

	if v, ok := params["period"].(float64); ok {
		period = int(v)
	}
	if v, ok := params["threshold"].(float64); ok {
		threshold = v
	}

	if len(prices) < period+1 {
		return &models.BacktestResult{
			Strategy:       s.Name(),
			InitialCapital: initialCap,
			FinalCapital:   initialCap,
			TotalReturn:    0,
		}
	}

	// 计算动量
	momentum := calculateMomentum(prices, period)

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
		currentMomentum := momentum[i]

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

		// 动量突破阈值买入
		if currentMomentum > threshold && position == 0 {
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
					Reason: fmt.Sprintf("动量突破: 动量=%.4f > 阈值=%.4f", currentMomentum, threshold),
				})
			}
		}

		// 动量低于阈值卖出
		if currentMomentum < -threshold && position > 0 {
			revenue := float64(position) * currentPrice
			trades = append(trades, models.Trade{
				Date:   currentDate,
				Type:   "SELL",
				Price:  currentPrice,
				Shares: position,
				Amount: revenue,
				Reason: fmt.Sprintf("动量回落: 动量=%.4f < -阈值=%.4f", currentMomentum, -threshold),
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

// calculateMomentum 计算动量指标
func calculateMomentum(prices []models.StockPrice, period int) []float64 {
	momentum := make([]float64, len(prices))

	for i := period; i < len(prices); i++ {
		if prices[i-period].Close != 0 {
			momentum[i] = (prices[i].Close - prices[i-period].Close) / prices[i-period].Close
		}
	}

	return momentum
}
