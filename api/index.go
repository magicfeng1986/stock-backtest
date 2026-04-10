package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"stock-backtest/backend/handlers"
)

// 嵌入的前端HTML
const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>股票回测系统</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 100vh; color: #333; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        header { text-align: center; color: white; padding: 40px 0; }
        header h1 { font-size: 2.5rem; margin-bottom: 10px; text-shadow: 2px 2px 4px rgba(0,0,0,0.2); }
        header p { font-size: 1.1rem; opacity: 0.9; }
        main { background: white; border-radius: 16px; box-shadow: 0 20px 60px rgba(0,0,0,0.3); overflow: hidden; }
        .control-panel { padding: 30px; background: #f8f9fa; border-bottom: 1px solid #e9ecef; }
        .form-group { margin-bottom: 20px; }
        .form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        label { display: block; margin-bottom: 8px; font-weight: 600; color: #495057; }
        input, select { width: 100%; padding: 12px 16px; border: 2px solid #dee2e6; border-radius: 8px; font-size: 1rem; transition: all 0.3s ease; }
        input:focus, select:focus { outline: none; border-color: #667eea; box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1); }
        .search-box { position: relative; }
        .search-results { position: absolute; top: 100%; left: 0; right: 0; background: white; border: 2px solid #dee2e6; border-top: none; border-radius: 0 0 8px 8px; max-height: 300px; overflow-y: auto; z-index: 100; display: none; }
        .search-results.active { display: block; }
        .search-item { padding: 12px 16px; cursor: pointer; border-bottom: 1px solid #f1f3f5; transition: background 0.2s; }
        .search-item:hover { background: #f8f9fa; }
        .search-item .code { font-weight: 600; color: #667eea; }
        .search-item .name { color: #6c757d; margin-left: 8px; }
        .selected-stock { margin-top: 10px; padding: 10px 16px; background: #e7f3ff; border-radius: 8px; display: none; }
        .selected-stock.active { display: block; }
        .selected-stock .label { color: #6c757d; font-size: 0.9rem; }
        .selected-stock .value { font-weight: 600; color: #495057; }
        .strategy-desc { margin-top: 8px; color: #6c757d; font-size: 0.9rem; }
        .strategy-params { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; margin-bottom: 20px; }
        .param-group { background: white; padding: 15px; border-radius: 8px; border: 1px solid #dee2e6; }
        .param-group label { font-size: 0.9rem; margin-bottom: 6px; }
        .param-group input { padding: 8px 12px; }
        .param-group .param-desc { font-size: 0.8rem; color: #6c757d; margin-top: 4px; }
        .btn-primary { width: 100%; padding: 16px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; border: none; border-radius: 8px; font-size: 1.1rem; font-weight: 600; cursor: pointer; transition: all 0.3s ease; }
        .btn-primary:hover:not(:disabled) { transform: translateY(-2px); box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4); }
        .btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
        .results-section { padding: 30px; }
        .results-section h2 { margin-bottom: 20px; color: #495057; }
        .metrics-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; margin-bottom: 30px; }
        .metric-card { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 12px; text-align: center; }
        .metric-card h3 { font-size: 0.9rem; opacity: 0.9; margin-bottom: 8px; }
        .metric-value { font-size: 1.5rem; font-weight: 700; }
        .metric-value.positive { color: #4ade80; }
        .metric-value.negative { color: #f87171; }
        .chart-container { background: white; padding: 20px; border-radius: 12px; border: 1px solid #e9ecef; margin-bottom: 30px; }
        .chart-container h3 { margin-bottom: 15px; color: #495057; }
        .trades-section { background: white; padding: 20px; border-radius: 12px; border: 1px solid #e9ecef; }
        .trades-section h3 { margin-bottom: 15px; color: #495057; }
        .table-container { overflow-x: auto; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #e9ecef; }
        th { background: #f8f9fa; font-weight: 600; color: #495057; }
        tr:hover { background: #f8f9fa; }
        .badge { display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; font-weight: 600; }
        .badge-buy { background: #d4edda; color: #155724; }
        .badge-sell { background: #f8d7da; color: #721c24; }
        .loading { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.7); display: flex; flex-direction: column; align-items: center; justify-content: center; z-index: 1000; }
        .loading p { color: white; margin-top: 20px; font-size: 1.1rem; }
        .spinner { width: 50px; height: 50px; border: 4px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite; }
        @keyframes spin { to { transform: rotate(360deg); } }
        .hidden { display: none !important; }
        footer { text-align: center; padding: 30px; color: rgba(255,255,255,0.8); font-size: 0.9rem; }
        @media (max-width: 768px) {
            .form-row { grid-template-columns: 1fr; }
            .metrics-grid { grid-template-columns: repeat(2, 1fr); }
            header h1 { font-size: 1.8rem; }
            .strategy-params { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>股票回测系统</h1>
            <p>基于历史数据的股票策略回测平台</p>
        </header>
        <main>
            <section class="control-panel">
                <div class="form-group">
                    <label for="stock-search">股票搜索</label>
                    <div class="search-box">
                        <input type="text" id="stock-search" placeholder="输入股票代码或名称 (如: 000001 或 平安银行)" autocomplete="off">
                        <div id="search-results" class="search-results"></div>
                    </div>
                    <div id="selected-stock" class="selected-stock"></div>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label for="start-date">开始日期</label>
                        <input type="date" id="start-date" value="2023-01-01">
                    </div>
                    <div class="form-group">
                        <label for="end-date">结束日期</label>
                        <input type="date" id="end-date" value="2024-12-31">
                    </div>
                </div>
                <div class="form-group">
                    <label for="strategy">选择策略</label>
                    <select id="strategy">
                        <option value="">请选择策略</option>
                    </select>
                    <p id="strategy-desc" class="strategy-desc"></p>
                </div>
                <div id="strategy-params" class="strategy-params"></div>
                <div class="form-group">
                    <label for="initial-cap">初始资金 (元)</label>
                    <input type="number" id="initial-cap" value="100000" min="10000" step="10000">
                </div>
                <button id="run-backtest" class="btn-primary" disabled>运行回测</button>
            </section>
            <section id="results-section" class="results-section hidden">
                <h2>回测结果</h2>
                <div class="metrics-grid">
                    <div class="metric-card">
                        <h3>总收益率</h3>
                        <p id="total-return" class="metric-value">-</p>
                    </div>
                    <div class="metric-card">
                        <h3>年化收益率</h3>
                        <p id="annualized-return" class="metric-value">-</p>
                    </div>
                    <div class="metric-card">
                        <h3>最大回撤</h3>
                        <p id="max-drawdown" class="metric-value">-</p>
                    </div>
                    <div class="metric-card">
                        <h3>夏普比率</h3>
                        <p id="sharpe-ratio" class="metric-value">-</p>
                    </div>
                    <div class="metric-card">
                        <h3>胜率</h3>
                        <p id="win-rate" class="metric-value">-</p>
                    </div>
                    <div class="metric-card">
                        <h3>交易次数</h3>
                        <p id="trade-count" class="metric-value">-</p>
                    </div>
                </div>
                <div class="chart-container">
                    <h3>权益曲线</h3>
                    <canvas id="equity-chart"></canvas>
                </div>
                <div class="trades-section">
                    <h3>交易记录</h3>
                    <div class="table-container">
                        <table id="trades-table">
                            <thead>
                                <tr>
                                    <th>日期</th>
                                    <th>类型</th>
                                    <th>价格</th>
                                    <th>数量</th>
                                    <th>金额</th>
                                    <th>理由</th>
                                </tr>
                            </thead>
                            <tbody></tbody>
                        </table>
                    </div>
                </div>
            </section>
        </main>
        <footer>
            <p>股票回测系统 - 仅供学习研究使用，不构成投资建议</p>
        </footer>
    </div>
    <div id="loading" class="loading hidden">
        <div class="spinner"></div>
        <p>正在运行回测...</p>
    </div>
    <script>
const API_BASE_URL = '/api';
let selectedStock = null;
let strategies = [];
let equityChart = null;
const stockSearch = document.getElementById('stock-search');
const searchResults = document.getElementById('search-results');
const selectedStockDiv = document.getElementById('selected-stock');
const strategySelect = document.getElementById('strategy');
const strategyDesc = document.getElementById('strategy-desc');
const strategyParams = document.getElementById('strategy-params');
const runBacktestBtn = document.getElementById('run-backtest');
const loading = document.getElementById('loading');
const resultsSection = document.getElementById('results-section');

async function init() {
    await loadStrategies();
    setupEventListeners();
    setDefaultDates();
}

async function loadStrategies() {
    try {
        const response = await fetch(API_BASE_URL + '/strategies');
        const data = await response.json();
        strategies = data.data;
        strategySelect.innerHTML = '<option value="">请选择策略</option>';
        strategies.forEach(strategy => {
            const option = document.createElement('option');
            option.value = strategy.id;
            option.textContent = strategy.name;
            strategySelect.appendChild(option);
        });
    } catch (error) {
        console.error('加载策略失败:', error);
    }
}

function setDefaultDates() {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setFullYear(endDate.getFullYear() - 1);
    document.getElementById('end-date').value = formatDate(endDate);
    document.getElementById('start-date').value = formatDate(startDate);
}

function formatDate(date) {
    return date.toISOString().split('T')[0];
}

function setupEventListeners() {
    let searchTimeout;
    stockSearch.addEventListener('input', (e) => {
        clearTimeout(searchTimeout);
        const query = e.target.value.trim();
        if (query.length === 0) {
            searchResults.classList.remove('active');
            return;
        }
        searchTimeout = setTimeout(() => searchStocks(query), 300);
    });
    document.addEventListener('click', (e) => {
        if (!e.target.closest('.search-box')) {
            searchResults.classList.remove('active');
        }
    });
    strategySelect.addEventListener('change', onStrategyChange);
    runBacktestBtn.addEventListener('click', runBacktest);
}

async function searchStocks(query) {
    try {
        const response = await fetch(API_BASE_URL + '/stocks/search?q=' + encodeURIComponent(query));
        const data = await response.json();
        displaySearchResults(data.data);
    } catch (error) {
        console.error('搜索股票失败:', error);
    }
}

function displaySearchResults(stocks) {
    searchResults.innerHTML = '';
    if (stocks.length === 0) {
        searchResults.innerHTML = '<div class="search-item">未找到匹配的股票</div>';
    } else {
        stocks.forEach(stock => {
            const item = document.createElement('div');
            item.className = 'search-item';
            item.innerHTML = '<span class="code">' + stock.code + '</span><span class="name">' + stock.name + '</span>';
            item.addEventListener('click', () => selectStock(stock));
            searchResults.appendChild(item);
        });
    }
    searchResults.classList.add('active');
}

function selectStock(stock) {
    selectedStock = stock;
    stockSearch.value = '';
    searchResults.classList.remove('active');
    selectedStockDiv.innerHTML = '<span class="label">已选择:</span><span class="value">' + stock.name + ' (' + stock.code + ')</span>';
    selectedStockDiv.classList.add('active');
    checkFormValidity();
}

function onStrategyChange() {
    const strategyId = strategySelect.value;
    const strategy = strategies.find(s => s.id === strategyId);
    if (strategy) {
        strategyDesc.textContent = strategy.description;
        renderStrategyParams(strategy);
    } else {
        strategyDesc.textContent = '';
        strategyParams.innerHTML = '';
    }
    checkFormValidity();
}

function renderStrategyParams(strategy) {
    strategyParams.innerHTML = '';
    if (!strategy.parameters || strategy.parameters.length === 0) return;
    strategy.parameters.forEach(param => {
        const paramGroup = document.createElement('div');
        paramGroup.className = 'param-group';
        const inputType = param.type === 'number' ? 'number' : 'text';
        const step = param.type === 'number' && param.default % 1 !== 0 ? '0.01' : '1';
        paramGroup.innerHTML = '<label for="param-' + param.name + '">' + param.name + '</label><input type="' + inputType + '" id="param-' + param.name + '" name="' + param.name + '" value="' + param.default + '" step="' + step + '"' + (param.min !== undefined ? ' min="' + param.min + '"' : '') + (param.max !== undefined ? ' max="' + param.max + '"' : '') + '><p class="param-desc">' + param.description + '</p>';
        strategyParams.appendChild(paramGroup);
    });
}

function checkFormValidity() {
    const isValid = selectedStock && strategySelect.value && document.getElementById('start-date').value && document.getElementById('end-date').value;
    runBacktestBtn.disabled = !isValid;
}

document.getElementById('start-date').addEventListener('change', checkFormValidity);
document.getElementById('end-date').addEventListener('change', checkFormValidity);

async function runBacktest() {
    if (!selectedStock || !strategySelect.value) {
        alert('请填写完整信息');
        return;
    }
    const params = {};
    const strategy = strategies.find(s => s.id === strategySelect.value);
    if (strategy && strategy.parameters) {
        strategy.parameters.forEach(param => {
            const input = document.getElementById('param-' + param.name);
            if (input) {
                params[param.name] = param.type === 'number' ? parseFloat(input.value) : input.value;
            }
        });
    }
    const requestData = {
        stockCode: selectedStock.code,
        startDate: document.getElementById('start-date').value,
        endDate: document.getElementById('end-date').value,
        strategy: strategySelect.value,
        initialCap: parseFloat(document.getElementById('initial-cap').value) || 100000,
        parameters: params
    };
    loading.classList.remove('hidden');
    try {
        const response = await fetch(API_BASE_URL + '/backtest', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(requestData)
        });
        const data = await response.json();
        if (response.ok) {
            displayResults(data.data);
        } else {
            alert(data.error || '回测失败');
        }
    } catch (error) {
        console.error('回测请求失败:', error);
        alert('网络错误，请稍后重试');
    } finally {
        loading.classList.add('hidden');
    }
}

function displayResults(result) {
    resultsSection.classList.remove('hidden');
    updateMetric('total-return', result.totalReturn, true);
    updateMetric('annualized-return', result.annualizedReturn, true);
    updateMetric('max-drawdown', result.maxDrawdown, false);
    updateMetric('sharpe-ratio', result.sharpeRatio, false);
    updateMetric('win-rate', result.winRate, false);
    document.getElementById('trade-count').textContent = result.tradeCount;
    drawEquityChart(result.equityCurve);
    displayTrades(result.trades);
    resultsSection.scrollIntoView({ behavior: 'smooth' });
}

function updateMetric(id, value, isPercentage) {
    const element = document.getElementById(id);
    const numValue = value * (isPercentage ? 100 : 1);
    const formatted = numValue.toFixed(2) + (isPercentage ? '%' : '');
    element.textContent = formatted;
    element.classList.remove('positive', 'negative');
    if (value > 0) element.classList.add('positive');
    else if (value < 0) element.classList.add('negative');
}

function drawEquityChart(equityCurve) {
    const ctx = document.getElementById('equity-chart').getContext('2d');
    if (equityChart) equityChart.destroy();
    const labels = equityCurve.map(p => p.date);
    const data = equityCurve.map(p => p.value);
    equityChart = new Chart(ctx, {
        type: 'line',
        data: { labels: labels, datasets: [{ label: '权益值', data: data, borderColor: '#667eea', backgroundColor: 'rgba(102, 126, 234, 0.1)', borderWidth: 2, fill: true, tension: 0.4, pointRadius: 0, pointHoverRadius: 6 }] },
        options: { responsive: true, maintainAspectRatio: false, interaction: { intersect: false, mode: 'index' }, plugins: { legend: { display: false }, tooltip: { callbacks: { label: function(context) { return '权益: ¥' + context.parsed.y.toFixed(2); } } } }, scales: { x: { grid: { display: false }, ticks: { maxTicksLimit: 10 } }, y: { beginAtZero: false, ticks: { callback: function(value) { return '¥' + value.toFixed(0); } } } } }
    });
}

function displayTrades(trades) {
    const tbody = document.querySelector('#trades-table tbody');
    tbody.innerHTML = '';
    if (!trades || trades.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" style="text-align: center;">暂无交易记录</td></tr>';
        return;
    }
    trades.forEach(trade => {
        const row = document.createElement('tr');
        const badgeClass = trade.type === 'BUY' ? 'badge-buy' : 'badge-sell';
        const typeText = trade.type === 'BUY' ? '买入' : '卖出';
        row.innerHTML = '<td>' + trade.date + '</td><td><span class="badge ' + badgeClass + '">' + typeText + '</span></td><td>¥' + trade.price.toFixed(2) + '</td><td>' + trade.shares + '</td><td>¥' + trade.amount.toFixed(2) + '</td><td>' + trade.reason + '</td>';
        tbody.appendChild(row);
    });
}

document.addEventListener('DOMContentLoaded', init);
    </script>
</body>
</html>`

// Handler Vercel serverless handler
func Handler(w http.ResponseWriter, r *http.Request) {
	// 如果是根路径或HTML请求，直接返回嵌入的HTML
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(indexHTML))
		return
	}

	// API请求使用Gin处理
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// 创建处理器
	h := handlers.NewHandler()

	// API路由
	api := router.Group("/api")
	{
		api.GET("/health", h.HealthCheck)
		api.GET("/stocks/search", h.SearchStocks)
		api.GET("/strategies", h.GetStrategies)
		api.POST("/backtest", h.RunBacktest)
		api.GET("/stock/data", h.GetStockData)
	}

	router.ServeHTTP(w, r)
}
