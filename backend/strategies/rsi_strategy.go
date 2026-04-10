package strategies

import (
	"fmt"
	"math"
	"stock-backtest/backend/models"
)

// RSIStrategy RSI策略
type RSIStrategy struct{}

// Name 策略名称
func (s *RSIStrategy) Name() string {
	return "RSI策略"
}

// Description 策略描述
func (s *RSIStrategy) Description() string {
	return "基于RSI超买超卖信号进行交易，RSI低于超卖阈值买入，高于超买阈值卖出"
}

// Execute 执行策略
func (s *RSIStrategy) Execute(prices []models.StockPrice, initialCap float64, params map[string]interface{}) *models.BacktestResult {
	// 获取参数
	period := 14
	overbought := 70.0
	oversold := 30.0

	if v, ok := params["period"].(float64); ok {
		period = int(v)
	}
	if v, ok := params["overbought"].(float64); ok {
		overbought = v
	}
	if v, ok := params["oversold"].(float64); ok {
		oversold = v
	}

	if len(prices) < period+1 {
		return &models.BacktestResult{
			Strategy:       s.Name(),
			InitialCapital: initialCap,
			FinalCapital:   initialCap,
			TotalReturn:    0,
		}
	}

	// 计算RSI
	rsi := calculateRSI(prices, period)

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
		currentRSI := rsi[i]

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

		// RSI超卖买入信号
		if currentRSI < oversold && position == 0 {
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
					Reason: fmt.Sprintf("RSI超卖: RSI=%.2f < %.2f", currentRSI, oversold),
				})
			}
		}

		// RSI超买卖出信号
		if currentRSI > overbought && position > 0 {
			revenue := float64(position) * currentPrice
			trades = append(trades, models.Trade{
				Date:   currentDate,
				Type:   "SELL",
				Price:  currentPrice,
				Shares: position,
				Amount: revenue,
				Reason: fmt.Sprintf("RSI超买: RSI=%.2f > %.2f", currentRSI, overbought),
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

// calculateRSI 计算RSI指标
func calculateRSI(prices []models.StockPrice, period int) []float64 {
	rsi := make([]float64, len(prices))

	if len(prices) < period+1 {
		return rsi
	}

	// 计算价格变化
	changes := make([]float64, len(prices))
	for i := 1; i < len(prices); i++ {
		changes[i] = prices[i].Close - prices[i-1].Close
	}

	// 计算初始平均涨跌
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		if changes[i] > 0 {
			avgGain += changes[i]
		} else {
			avgLoss -= changes[i]
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// 计算第一个RSI
	if avgLoss == 0 {
		rsi[period] = 100
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// 计算后续RSI
	for i := period + 1; i < len(prices); i++ {
		gain := 0.0
		loss := 0.0
		if changes[i] > 0 {
			gain = changes[i]
		} else {
			loss = -changes[i]
		}

		avgGain = (avgGain*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	return rsi
}
