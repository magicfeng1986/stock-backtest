package services

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"stock-backtest/backend/models"
)

// DataService 数据服务
type DataService struct {
	// 模拟数据存储
	mockData map[string][]models.StockPrice
}

// NewDataService 创建数据服务
func NewDataService() *DataService {
	service := &DataService{
		mockData: make(map[string][]models.StockPrice),
	}
	service.initMockData()
	return service
}

// SearchStocks 搜索股票
func (s *DataService) SearchStocks(query string) []models.Stock {
	// 模拟股票数据库
	stocks := []models.Stock{
		{Code: "000001.SZ", Name: "平安银行", Market: "SZ"},
		{Code: "000002.SZ", Name: "万科A", Market: "SZ"},
		{Code: "000063.SZ", Name: "中兴通讯", Market: "SZ"},
		{Code: "000100.SZ", Name: "TCL科技", Market: "SZ"},
		{Code: "000333.SZ", Name: "美的集团", Market: "SZ"},
		{Code: "000568.SZ", Name: "泸州老窖", Market: "SZ"},
		{Code: "000651.SZ", Name: "格力电器", Market: "SZ"},
		{Code: "000725.SZ", Name: "京东方A", Market: "SZ"},
		{Code: "000858.SZ", Name: "五粮液", Market: "SZ"},
		{Code: "002001.SZ", Name: "新和成", Market: "SZ"},
		{Code: "002007.SZ", Name: "华兰生物", Market: "SZ"},
		{Code: "002024.SZ", Name: "苏宁易购", Market: "SZ"},
		{Code: "002027.SZ", Name: "分众传媒", Market: "SZ"},
		{Code: "002142.SZ", Name: "宁波银行", Market: "SZ"},
		{Code: "002230.SZ", Name: "科大讯飞", Market: "SZ"},
		{Code: "002236.SZ", Name: "大华股份", Market: "SZ"},
		{Code: "002271.SZ", Name: "东方雨虹", Market: "SZ"},
		{Code: "002304.SZ", Name: "洋河股份", Market: "SZ"},
		{Code: "002352.SZ", Name: "顺丰控股", Market: "SZ"},
		{Code: "002415.SZ", Name: "海康威视", Market: "SZ"},
		{Code: "002460.SZ", Name: "赣锋锂业", Market: "SZ"},
		{Code: "002475.SZ", Name: "立讯精密", Market: "SZ"},
		{Code: "002594.SZ", Name: "比亚迪", Market: "SZ"},
		{Code: "002714.SZ", Name: "牧原股份", Market: "SZ"},
		{Code: "002812.SZ", Name: "恩捷股份", Market: "SZ"},
		{Code: "300001.SZ", Name: "特锐德", Market: "SZ"},
		{Code: "300003.SZ", Name: "乐普医疗", Market: "SZ"},
		{Code: "300014.SZ", Name: "亿纬锂能", Market: "SZ"},
		{Code: "300015.SZ", Name: "爱尔眼科", Market: "SZ"},
		{Code: "300033.SZ", Name: "同花顺", Market: "SZ"},
		{Code: "300059.SZ", Name: "东方财富", Market: "SZ"},
		{Code: "300122.SZ", Name: "智飞生物", Market: "SZ"},
		{Code: "300124.SZ", Name: "汇川技术", Market: "SZ"},
		{Code: "300142.SZ", Name: "沃森生物", Market: "SZ"},
		{Code: "300274.SZ", Name: "阳光电源", Market: "SZ"},
		{Code: "300408.SZ", Name: "三环集团", Market: "SZ"},
		{Code: "300413.SZ", Name: "芒果超媒", Market: "SZ"},
		{Code: "300433.SZ", Name: "蓝思科技", Market: "SZ"},
		{Code: "300498.SZ", Name: "温氏股份", Market: "SZ"},
		{Code: "300750.SZ", Name: "宁德时代", Market: "SZ"},
		{Code: "600000.SH", Name: "浦发银行", Market: "SH"},
		{Code: "600009.SH", Name: "上海机场", Market: "SH"},
		{Code: "600016.SH", Name: "民生银行", Market: "SH"},
		{Code: "600028.SH", Name: "中国石化", Market: "SH"},
		{Code: "600030.SH", Name: "中信证券", Market: "SH"},
		{Code: "600031.SH", Name: "三一重工", Market: "SH"},
		{Code: "600036.SH", Name: "招商银行", Market: "SH"},
		{Code: "600048.SH", Name: "保利发展", Market: "SH"},
		{Code: "600050.SH", Name: "中国联通", Market: "SH"},
		{Code: "600104.SH", Name: "上汽集团", Market: "SH"},
		{Code: "600196.SH", Name: "复星医药", Market: "SH"},
		{Code: "600276.SH", Name: "恒瑞医药", Market: "SH"},
		{Code: "600309.SH", Name: "万华化学", Market: "SH"},
		{Code: "600332.SH", Name: "白云山", Market: "SH"},
		{Code: "600340.SH", Name: "华夏幸福", Market: "SH"},
		{Code: "600352.SH", Name: "浙江龙盛", Market: "SH"},
		{Code: "600362.SH", Name: "江西铜业", Market: "SH"},
		{Code: "600383.SH", Name: "金地集团", Market: "SH"},
		{Code: "600406.SH", Name: "国电南瑞", Market: "SH"},
		{Code: "600436.SH", Name: "片仔癀", Market: "SH"},
		{Code: "600438.SH", Name: "通威股份", Market: "SH"},
		{Code: "600519.SH", Name: "贵州茅台", Market: "SH"},
		{Code: "600585.SH", Name: "海螺水泥", Market: "SH"},
		{Code: "600588.SH", Name: "用友网络", Market: "SH"},
		{Code: "600600.SH", Name: "青岛啤酒", Market: "SH"},
		{Code: "600606.SH", Name: "绿地控股", Market: "SH"},
		{Code: "600637.SH", Name: "东方明珠", Market: "SH"},
		{Code: "600660.SH", Name: "福耀玻璃", Market: "SH"},
		{Code: "600690.SH", Name: "海尔智家", Market: "SH"},
		{Code: "600703.SH", Name: "三安光电", Market: "SH"},
		{Code: "600745.SH", Name: "闻泰科技", Market: "SH"},
		{Code: "600809.SH", Name: "山西汾酒", Market: "SH"},
		{Code: "600837.SH", Name: "海通证券", Market: "SH"},
		{Code: "600887.SH", Name: "伊利股份", Market: "SH"},
		{Code: "600893.SH", Name: "航发动力", Market: "SH"},
		{Code: "600900.SH", Name: "长江电力", Market: "SH"},
		{Code: "601012.SH", Name: "隆基绿能", Market: "SH"},
		{Code: "601066.SH", Name: "中信建投", Market: "SH"},
		{Code: "601088.SH", Name: "中国神华", Market: "SH"},
		{Code: "601111.SH", Name: "中国国航", Market: "SH"},
		{Code: "601138.SH", Name: "工业富联", Market: "SH"},
		{Code: "601166.SH", Name: "兴业银行", Market: "SH"},
		{Code: "601186.SH", Name: "中国铁建", Market: "SH"},
		{Code: "601288.SH", Name: "农业银行", Market: "SH"},
		{Code: "601318.SH", Name: "中国平安", Market: "SH"},
		{Code: "601319.SH", Name: "中国人保", Market: "SH"},
		{Code: "601328.SH", Name: "交通银行", Market: "SH"},
		{Code: "601336.SH", Name: "新华保险", Market: "SH"},
		{Code: "601390.SH", Name: "中国中铁", Market: "SH"},
		{Code: "601398.SH", Name: "工商银行", Market: "SH"},
		{Code: "601601.SH", Name: "中国太保", Market: "SH"},
		{Code: "601628.SH", Name: "中国人寿", Market: "SH"},
		{Code: "601668.SH", Name: "中国建筑", Market: "SH"},
		{Code: "601688.SH", Name: "华泰证券", Market: "SH"},
		{Code: "601766.SH", Name: "中国中车", Market: "SH"},
		{Code: "601818.SH", Name: "光大银行", Market: "SH"},
		{Code: "601857.SH", Name: "中国石油", Market: "SH"},
		{Code: "601888.SH", Name: "中国中免", Market: "SH"},
		{Code: "601899.SH", Name: "紫金矿业", Market: "SH"},
		{Code: "601933.SH", Name: "永辉超市", Market: "SH"},
		{Code: "601988.SH", Name: "中国银行", Market: "SH"},
		{Code: "601989.SH", Name: "中国重工", Market: "SH"},
		{Code: "603019.SH", Name: "中科曙光", Market: "SH"},
		{Code: "603288.SH", Name: "海天味业", Market: "SH"},
		{Code: "603501.SH", Name: "韦尔股份", Market: "SH"},
		{Code: "603659.SH", Name: "璞泰来", Market: "SH"},
		{Code: "603799.SH", Name: "华友钴业", Market: "SH"},
		{Code: "603986.SH", Name: "兆易创新", Market: "SH"},
		{Code: "603993.SH", Name: "洛阳钼业", Market: "SH"},
		{Code: "AAPL", Name: "苹果公司", Market: "US"},
		{Code: "MSFT", Name: "微软", Market: "US"},
		{Code: "GOOGL", Name: "谷歌", Market: "US"},
		{Code: "AMZN", Name: "亚马逊", Market: "US"},
		{Code: "TSLA", Name: "特斯拉", Market: "US"},
		{Code: "META", Name: "Meta", Market: "US"},
		{Code: "NVDA", Name: "英伟达", Market: "US"},
		{Code: "NFLX", Name: "奈飞", Market: "US"},
		{Code: "BABA", Name: "阿里巴巴", Market: "US"},
		{Code: "JD", Name: "京东", Market: "US"},
	}

	var results []models.Stock
	query = strings.ToUpper(strings.TrimSpace(query))

	for _, stock := range stocks {
		if strings.Contains(strings.ToUpper(stock.Code), query) ||
			strings.Contains(strings.ToUpper(stock.Name), query) {
			results = append(results, stock)
		}
	}

	return results
}

// GetStockData 获取股票历史数据
func (s *DataService) GetStockData(stockCode string, startDate, endDate string) ([]models.StockPrice, error) {
	// 首先尝试获取模拟数据
	if data, ok := s.mockData[stockCode]; ok {
		var filtered []models.StockPrice
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)

		for _, price := range data {
			if !price.Date.Before(start) && !price.Date.After(end) {
				filtered = append(filtered, price)
			}
		}
		return filtered, nil
	}

	// 如果没有模拟数据，生成随机数据
	return s.generateMockData(stockCode, startDate, endDate), nil
}

// GetStockName 获取股票名称
func (s *DataService) GetStockName(stockCode string) string {
	stocks := s.SearchStocks(stockCode)
	for _, stock := range stocks {
		if stock.Code == stockCode {
			return stock.Name
		}
	}
	return stockCode
}

// initMockData 初始化模拟数据
func (s *DataService) initMockData() {
	// 为一些热门股票生成模拟数据
	hotStocks := []string{"000001.SZ", "000002.SZ", "000333.SZ", "000651.SZ", "000858.SZ",
		"002594.SZ", "300750.SZ", "600000.SH", "600036.SH", "600519.SH",
		"601012.SH", "601318.SH", "603288.SH", "AAPL", "MSFT", "TSLA", "NVDA", "BABA"}

	for _, code := range hotStocks {
		s.mockData[code] = s.generateMockData(code, "2020-01-01", "2024-12-31")
	}
}

// generateMockData 生成模拟股票数据
func (s *DataService) generateMockData(stockCode string, startDate, endDate string) []models.StockPrice {
	var prices []models.StockPrice
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	// 基于股票代码生成一个固定的初始价格
	basePrice := 50.0 + float64(len(stockCode))*10
	if stockCode == "600519.SH" { // 贵州茅台
		basePrice = 1500
	} else if stockCode == "300750.SZ" { // 宁德时代
		basePrice = 200
	}

	currentPrice := basePrice
	date := start

	rand.Seed(int64(len(stockCode)))

	for !date.After(end) {
		// 跳过周末
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			// 生成随机波动 (-3% 到 +3%)
			change := (rand.Float64() - 0.5) * 0.06
			currentPrice = currentPrice * (1 + change)

			// 确保价格不会太低
			if currentPrice < basePrice*0.3 {
				currentPrice = basePrice * 0.3
			}

			open := currentPrice * (1 + (rand.Float64()-0.5)*0.01)
			high := currentPrice * (1 + rand.Float64()*0.02)
			low := currentPrice * (1 - rand.Float64()*0.02)
			close := currentPrice
			volume := int64(1000000 + rand.Int63n(9000000))

			prices = append(prices, models.StockPrice{
				Date:     date,
				Open:     open,
				High:     high,
				Low:      low,
				Close:    close,
				Volume:   volume,
				AdjClose: close,
			})
		}
		date = date.AddDate(0, 0, 1)
	}

	return prices
}

// FetchYahooFinance 从 Yahoo Finance 获取数据（备用）
func (s *DataService) FetchYahooFinance(symbol string) ([]models.StockPrice, error) {
	// Yahoo Finance API URL
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=0&period2=9999999999&interval=1d&events=history", symbol)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var prices []models.StockPrice
	// 跳过标题行
	for i, record := range records {
		if i == 0 {
			continue
		}
		date, _ := time.Parse("2006-01-02", record[0])
		open, _ := strconv.ParseFloat(record[1], 64)
		high, _ := strconv.ParseFloat(record[2], 64)
		low, _ := strconv.ParseFloat(record[3], 64)
		close, _ := strconv.ParseFloat(record[4], 64)
		adjClose, _ := strconv.ParseFloat(record[5], 64)
		volume, _ := strconv.ParseInt(record[6], 10, 64)

		prices = append(prices, models.StockPrice{
			Date:     date,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			Volume:   volume,
			AdjClose: adjClose,
		})
	}

	return prices, nil
}
