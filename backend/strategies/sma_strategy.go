package strategies

import (
	"fmt"
	"math"
	"stock-backtest/backend/models"
)

// SMAStrategy 简单移动平均线策略
type SMAStrategy struct{}

// Name 策略名称
func (s *SMAStrategy) Name() string {
	return "均线策略 (SMA Crossover)"
}

// Description 策略描述
func (s *SMAStrategy) Description() string {
	return "当短期均线上穿长期均线时买入，下穿时卖出"
}

// Execute 执行策略
func (s *SMAStrategy) Execute(prices []models.StockPrice, initialCap float64, params map[string]interface{}) *models.BacktestResult {
	// 获取参数
	shortPeriod := 10
	longPeriod := 30

	if v, ok := params["shortPeriod"].(float64); ok {
		shortPeriod = int(v)
	}
	if v, ok := params["longPeriod"].(float64); ok {
		longPeriod = int(v)
	}

	if len(prices) < longPeriod+1 {
		return &models.BacktestResult{
			Strategy:       s.Name(),
			InitialCapital: initialCap,
			FinalCapital:   initialCap,
			TotalReturn:    0,
		}
	}

	// 计算移动平均线
	shortMA := calculateSMA(prices, shortPeriod)
	longMA := calculateSMA(prices, longPeriod)

	// 回测逻辑
	capital := initialCap
	position := 0 // 持仓数量
	var trades []models.Trade
	var equityCurve []models.EquityPoint
	var dailyReturns []models.DailyReturn

	maxCapital := initialCap
	maxDrawdown := 0.0

	for i := longPeriod; i < len(prices); i++ {
		currentPrice := prices[i].Close
		currentDate := prices[i].Date.Format("2006-01-02")

		// 计算当前权益
		currentEquity := capital + float64(position)*currentPrice

		// 记录权益曲线
		equityCurve = append(equityCurve, models.EquityPoint{
			Date:  currentDate,
			Value: currentEquity,
		})

		// 计算每日收益
		if i > longPeriod {
			prevEquity := equityCurve[len(equityCurve)-2].Value
			dailyReturn := (currentEquity - prevEquity) / prevEquity
			dailyReturns = append(dailyReturns, models.DailyReturn{
				Date:   currentDate,
				Return: dailyReturn,
				Value:  currentEquity,
			})
		}

		// 更新最大回撤
		if currentEquity > maxCapital {
			maxCapital = currentEquity
		}
		drawdown := (maxCapital - currentEquity) / maxCapital
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

		// 交易信号
		if i > longPeriod {
			prevShortMA := shortMA[i-1]
			prevLongMA := longMA[i-1]
			currShortMA := shortMA[i]
			currLongMA := longMA[i]

			// 金叉买入信号
			if prevShortMA <= prevLongMA && currShortMA > currLongMA && position == 0 {
				// 买入
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
						Reason: fmt.Sprintf("金叉: 短期均线(%.2f)上穿长期均线(%.2f)", currShortMA, currLongMA),
					})
				}
			}

			// 死叉卖出信号
			if prevShortMA >= prevLongMA && currShortMA < currLongMA && position > 0 {
				// 卖出
				revenue := float64(position) * currentPrice
				trades = append(trades, models.Trade{
					Date:   currentDate,
					Type:   "SELL",
					Price:  currentPrice,
					Shares: position,
					Amount: revenue,
					Reason: fmt.Sprintf("死叉: 短期均线(%.2f)下穿长期均线(%.2f)", currShortMA, currLongMA),
				})
				capital += revenue
				position = 0
			}
		}
	}

	// 计算最终收益
	finalCapital := capital + float64(position)*prices[len(prices)-1].Close
	totalReturn := (finalCapital - initialCap) / initialCap

	// 计算年化收益
	days := len(prices) - longPeriod
	years := float64(days) / 252.0
	annualizedReturn := 0.0
	if years > 0 {
		annualizedReturn = math.Pow(1+totalReturn, 1/years) - 1
	}

	// 计算夏普比率
	sharpeRatio := calculateSharpeRatio(dailyReturns)

	// 计算胜率
	winRate := calculateWinRate(trades)

	return &models.BacktestResult{
		Strategy:         s.Name(),
		StartDate:        prices[longPeriod].Date.Format("2006-01-02"),
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

// calculateSMA 计算简单移动平均线
func calculateSMA(prices []models.StockPrice, period int) []float64 {
	ma := make([]float64, len(prices))
	for i := 0; i < len(prices); i++ {
		if i < period-1 {
			ma[i] = 0
			continue
		}
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i-j].Close
		}
		ma[i] = sum / float64(period)
	}
	return ma
}

// calculateSharpeRatio 计算夏普比率
func calculateSharpeRatio(returns []models.DailyReturn) float64 {
	if len(returns) == 0 {
		return 0
	}

	// 计算平均收益
	var sumReturn float64
	for _, r := range returns {
		sumReturn += r.Return
	}
	avgReturn := sumReturn / float64(len(returns))

	// 计算标准差
	var sumSquaredDiff float64
	for _, r := range returns {
		diff := r.Return - avgReturn
		sumSquaredDiff += diff * diff
	}
	stdDev := math.Sqrt(sumSquaredDiff / float64(len(returns)))

	if stdDev == 0 {
		return 0
	}

	// 假设无风险利率为2%年化，日利率约为 0.02/252
	riskFreeRate := 0.02 / 252
	return (avgReturn - riskFreeRate) / stdDev * math.Sqrt(252)
}

// calculateWinRate 计算胜率
func calculateWinRate(trades []models.Trade) float64 {
	if len(trades) < 2 {
		return 0
	}

	wins := 0
	tradesCount := 0

	for i := 0; i < len(trades)-1; i++ {
		if trades[i].Type == "BUY" {
			tradesCount++
			// 找到对应的卖出交易
			for j := i + 1; j < len(trades); j++ {
				if trades[j].Type == "SELL" {
					if trades[j].Price > trades[i].Price {
						wins++
					}
					break
				}
			}
		}
	}

	if tradesCount == 0 {
		return 0
	}

	return float64(wins) / float64(tradesCount)
}
