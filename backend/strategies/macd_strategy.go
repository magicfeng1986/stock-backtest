package strategies

import (
	"fmt"
	"math"
	"stock-backtest/backend/models"
)

// MACDStrategy MACD策略
type MACDStrategy struct{}

// Name 策略名称
func (s *MACDStrategy) Name() string {
	return "MACD策略"
}

// Description 策略描述
func (s *MACDStrategy) Description() string {
	return "基于MACD指标的金叉死叉信号进行交易"
}

// Execute 执行策略
func (s *MACDStrategy) Execute(prices []models.StockPrice, initialCap float64, params map[string]interface{}) *models.BacktestResult {
	// 获取参数
	fastPeriod := 12
	slowPeriod := 26
	signalPeriod := 9

	if v, ok := params["fastPeriod"].(float64); ok {
		fastPeriod = int(v)
	}
	if v, ok := params["slowPeriod"].(float64); ok {
		slowPeriod = int(v)
	}
	if v, ok := params["signalPeriod"].(float64); ok {
		signalPeriod = int(v)
	}

	if len(prices) < slowPeriod+signalPeriod+1 {
		return &models.BacktestResult{
			Strategy:       s.Name(),
			InitialCapital: initialCap,
			FinalCapital:   initialCap,
			TotalReturn:    0,
		}
	}

	// 计算EMA
	fastEMA := calculateEMA(prices, fastPeriod)
	slowEMA := calculateEMA(prices, slowPeriod)

	// 计算MACD线
	macdLine := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		macdLine[i] = fastEMA[i] - slowEMA[i]
	}

	// 计算信号线 (MACD的EMA)
	signalLine := calculateEMAFromValues(macdLine, signalPeriod)

	// 计算MACD柱状图
	histogram := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		histogram[i] = macdLine[i] - signalLine[i]
	}

	// 回测逻辑
	capital := initialCap
	position := 0
	var trades []models.Trade
	var equityCurve []models.EquityPoint
	var dailyReturns []models.DailyReturn

	maxCapital := initialCap
	maxDrawdown := 0.0

	startIndex := slowPeriod + signalPeriod

	for i := startIndex; i < len(prices); i++ {
		currentPrice := prices[i].Close
		currentDate := prices[i].Date.Format("2006-01-02")

		currentEquity := capital + float64(position)*currentPrice

		equityCurve = append(equityCurve, models.EquityPoint{
			Date:  currentDate,
			Value: currentEquity,
		})

		if i > startIndex {
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

		// MACD金叉买入信号
		if i > startIndex {
			prevMACD := macdLine[i-1]
			prevSignal := signalLine[i-1]
			currMACD := macdLine[i]
			currSignal := signalLine[i]

			if prevMACD <= prevSignal && currMACD > currSignal && position == 0 {
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
						Reason: fmt.Sprintf("MACD金叉: MACD(%.3f)上穿信号线(%.3f)", currMACD, currSignal),
					})
				}
			}

			// MACD死叉卖出信号
			if prevMACD >= prevSignal && currMACD < currSignal && position > 0 {
				revenue := float64(position) * currentPrice
				trades = append(trades, models.Trade{
					Date:   currentDate,
					Type:   "SELL",
					Price:  currentPrice,
					Shares: position,
					Amount: revenue,
					Reason: fmt.Sprintf("MACD死叉: MACD(%.3f)下穿信号线(%.3f)", currMACD, currSignal),
				})
				capital += revenue
				position = 0
			}
		}
	}

	finalCapital := capital + float64(position)*prices[len(prices)-1].Close
	totalReturn := (finalCapital - initialCap) / initialCap

	days := len(prices) - startIndex
	years := float64(days) / 252.0
	annualizedReturn := 0.0
	if years > 0 {
		annualizedReturn = math.Pow(1+totalReturn, 1/years) - 1
	}

	sharpeRatio := calculateSharpeRatio(dailyReturns)
	winRate := calculateWinRate(trades)

	return &models.BacktestResult{
		Strategy:         s.Name(),
		StartDate:        prices[startIndex].Date.Format("2006-01-02"),
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

// calculateEMA 计算指数移动平均线
func calculateEMA(prices []models.StockPrice, period int) []float64 {
	ema := make([]float64, len(prices))
	multiplier := 2.0 / float64(period+1)

	// 第一个EMA使用SMA
	if len(prices) >= period {
		sum := 0.0
		for i := 0; i < period; i++ {
			sum += prices[i].Close
		}
		ema[period-1] = sum / float64(period)
	}

	// 计算后续EMA
	for i := period; i < len(prices); i++ {
		ema[i] = (prices[i].Close-ema[i-1])*multiplier + ema[i-1]
	}

	return ema
}

// calculateEMAFromValues 从数值数组计算EMA
func calculateEMAFromValues(values []float64, period int) []float64 {
	ema := make([]float64, len(values))
	multiplier := 2.0 / float64(period+1)

	// 找到第一个非零值
	startIdx := 0
	for i, v := range values {
		if v != 0 {
			startIdx = i
			break
		}
	}

	// 第一个EMA使用SMA
	if startIdx+period <= len(values) {
		sum := 0.0
		for i := startIdx; i < startIdx+period; i++ {
			sum += values[i]
		}
		ema[startIdx+period-1] = sum / float64(period)
	}

	// 计算后续EMA
	for i := startIdx + period; i < len(values); i++ {
		ema[i] = (values[i]-ema[i-1])*multiplier + ema[i-1]
	}

	return ema
}
