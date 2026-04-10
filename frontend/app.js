// API 基础URL
const API_BASE_URL = window.location.hostname === 'localhost'
    ? 'http://localhost:8080/api'
    : '/api';

// 全局状态
let selectedStock = null;
let strategies = [];
let equityChart = null;

// DOM 元素
const stockSearch = document.getElementById('stock-search');
const searchResults = document.getElementById('search-results');
const selectedStockDiv = document.getElementById('selected-stock');
const strategySelect = document.getElementById('strategy');
const strategyDesc = document.getElementById('strategy-desc');
const strategyParams = document.getElementById('strategy-params');
const runBacktestBtn = document.getElementById('run-backtest');
const loading = document.getElementById('loading');
const resultsSection = document.getElementById('results-section');

// 初始化
async function init() {
    await loadStrategies();
    setupEventListeners();
    setDefaultDates();
}

// 加载策略列表
async function loadStrategies() {
    try {
        const response = await fetch(`${API_BASE_URL}/strategies`);
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
        showError('加载策略列表失败，请刷新页面重试');
    }
}

// 设置默认日期
function setDefaultDates() {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setFullYear(endDate.getFullYear() - 1);

    document.getElementById('end-date').value = formatDate(endDate);
    document.getElementById('start-date').value = formatDate(startDate);
}

// 格式化日期
function formatDate(date) {
    return date.toISOString().split('T')[0];
}

// 设置事件监听
function setupEventListeners() {
    // 股票搜索
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

    // 点击外部关闭搜索结果
    document.addEventListener('click', (e) => {
        if (!e.target.closest('.search-box')) {
            searchResults.classList.remove('active');
        }
    });

    // 策略选择
    strategySelect.addEventListener('change', onStrategyChange);

    // 运行回测
    runBacktestBtn.addEventListener('click', runBacktest);
}

// 搜索股票
async function searchStocks(query) {
    try {
        const response = await fetch(`${API_BASE_URL}/stocks/search?q=${encodeURIComponent(query)}`);
        const data = await response.json();

        displaySearchResults(data.data);
    } catch (error) {
        console.error('搜索股票失败:', error);
    }
}

// 显示搜索结果
function displaySearchResults(stocks) {
    searchResults.innerHTML = '';

    if (stocks.length === 0) {
        searchResults.innerHTML = '<div class="search-item">未找到匹配的股票</div>';
    } else {
        stocks.forEach(stock => {
            const item = document.createElement('div');
            item.className = 'search-item';
            item.innerHTML = `
                <span class="code">${stock.code}</span>
                <span class="name">${stock.name}</span>
            `;
            item.addEventListener('click', () => selectStock(stock));
            searchResults.appendChild(item);
        });
    }

    searchResults.classList.add('active');
}

// 选择股票
function selectStock(stock) {
    selectedStock = stock;
    stockSearch.value = '';
    searchResults.classList.remove('active');

    selectedStockDiv.innerHTML = `
        <span class="label">已选择:</span>
        <span class="value">${stock.name} (${stock.code})</span>
    `;
    selectedStockDiv.classList.add('active');

    checkFormValidity();
}

// 策略变更
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

// 渲染策略参数
function renderStrategyParams(strategy) {
    strategyParams.innerHTML = '';

    if (!strategy.parameters || strategy.parameters.length === 0) {
        return;
    }

    strategy.parameters.forEach(param => {
        const paramGroup = document.createElement('div');
        paramGroup.className = 'param-group';

        const inputType = param.type === 'number' ? 'number' : 'text';
        const step = param.type === 'number' && param.default % 1 !== 0 ? '0.01' : '1';

        paramGroup.innerHTML = `
            <label for="param-${param.name}">${param.name}</label>
            <input type="${inputType}"
                   id="param-${param.name}"
                   name="${param.name}"
                   value="${param.default}"
                   step="${step}"
                   ${param.min !== undefined ? `min="${param.min}"` : ''}
                   ${param.max !== undefined ? `max="${param.max}"` : ''}>
            <p class="param-desc">${param.description}</p>
        `;

        strategyParams.appendChild(paramGroup);
    });
}

// 检查表单有效性
function checkFormValidity() {
    const isValid = selectedStock &&
                   strategySelect.value &&
                   document.getElementById('start-date').value &&
                   document.getElementById('end-date').value;

    runBacktestBtn.disabled = !isValid;
}

// 监听日期变化
document.getElementById('start-date').addEventListener('change', checkFormValidity);
document.getElementById('end-date').addEventListener('change', checkFormValidity);

// 运行回测
async function runBacktest() {
    if (!selectedStock || !strategySelect.value) {
        showError('请填写完整信息');
        return;
    }

    // 收集参数
    const params = {};
    const strategy = strategies.find(s => s.id === strategySelect.value);
    if (strategy && strategy.parameters) {
        strategy.parameters.forEach(param => {
            const input = document.getElementById(`param-${param.name}`);
            if (input) {
                params[param.name] = param.type === 'number'
                    ? parseFloat(input.value)
                    : input.value;
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
        const response = await fetch(`${API_BASE_URL}/backtest`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestData)
        });

        const data = await response.json();

        if (response.ok) {
            displayResults(data.data);
        } else {
            showError(data.error || '回测失败');
        }
    } catch (error) {
        console.error('回测请求失败:', error);
        showError('网络错误，请稍后重试');
    } finally {
        loading.classList.add('hidden');
    }
}

// 显示回测结果
function displayResults(result) {
    resultsSection.classList.remove('hidden');

    // 更新指标卡片
    updateMetric('total-return', result.totalReturn, true);
    updateMetric('annualized-return', result.annualizedReturn, true);
    updateMetric('max-drawdown', result.maxDrawdown, false);
    updateMetric('sharpe-ratio', result.sharpeRatio, false);
    updateMetric('win-rate', result.winRate, false);
    document.getElementById('trade-count').textContent = result.tradeCount;

    // 绘制权益曲线
    drawEquityChart(result.equityCurve);

    // 显示交易记录
    displayTrades(result.trades);

    // 滚动到结果区域
    resultsSection.scrollIntoView({ behavior: 'smooth' });
}

// 更新指标显示
function updateMetric(id, value, isPercentage) {
    const element = document.getElementById(id);
    const numValue = value * (isPercentage ? 100 : 1);
    const formatted = numValue.toFixed(2) + (isPercentage ? '%' : '');

    element.textContent = formatted;
    element.classList.remove('positive', 'negative');

    if (value > 0) {
        element.classList.add('positive');
    } else if (value < 0) {
        element.classList.add('negative');
    }
}

// 绘制权益曲线
function drawEquityChart(equityCurve) {
    const ctx = document.getElementById('equity-chart').getContext('2d');

    if (equityChart) {
        equityChart.destroy();
    }

    const labels = equityCurve.map(p => p.date);
    const data = equityCurve.map(p => p.value);

    equityChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: '权益值',
                data: data,
                borderColor: '#667eea',
                backgroundColor: 'rgba(102, 126, 234, 0.1)',
                borderWidth: 2,
                fill: true,
                tension: 0.4,
                pointRadius: 0,
                pointHoverRadius: 6
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
                intersect: false,
                mode: 'index'
            },
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return '权益: ¥' + context.parsed.y.toFixed(2);
                        }
                    }
                }
            },
            scales: {
                x: {
                    grid: {
                        display: false
                    },
                    ticks: {
                        maxTicksLimit: 10
                    }
                },
                y: {
                    beginAtZero: false,
                    ticks: {
                        callback: function(value) {
                            return '¥' + value.toFixed(0);
                        }
                    }
                }
            }
        }
    });
}

// 显示交易记录
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

        row.innerHTML = `
            <td>${trade.date}</td>
            <td><span class="badge ${badgeClass}">${typeText}</span></td>
            <td>¥${trade.price.toFixed(2)}</td>
            <td>${trade.shares}</td>
            <td>¥${trade.amount.toFixed(2)}</td>
            <td>${trade.reason}</td>
        `;

        tbody.appendChild(row);
    });
}

// 显示错误信息
function showError(message) {
    alert(message);
}

// 启动应用
document.addEventListener('DOMContentLoaded', init);
