# 股票回测系统 (Stock Backtest System)

一个基于 Go 后端和 JavaScript 前端的股票回测系统，支持多种策略回测。

## 功能特性

- 支持股票代码搜索和选择
- 多种回测策略（均线策略、MACD策略、RSI策略等）
- 可视化回测结果展示
- 实时计算收益率和交易记录

## 技术栈

- **后端**: Go + Gin 框架
- **前端**: HTML5 + JavaScript + Chart.js
- **部署**: Vercel

## 项目结构

```
.
├── api/              # Vercel API 入口
│   └── index.go      # 主入口文件
├── backend/          # Go 后端代码
│   ├── main.go       # 主程序
│   ├── handlers/     # HTTP 处理器
│   ├── strategies/   # 回测策略
│   ├── models/       # 数据模型
│   └── services/     # 业务逻辑
├── frontend/         # 前端代码
│   ├── index.html
│   ├── app.js
│   └── style.css
├── go.mod
├── go.sum
└── vercel.json
```

## 本地开发

```bash
# 启动后端
cd backend
go run main.go

# 访问前端
open frontend/index.html
```

## 部署

项目已配置 Vercel 自动部署，推送代码到 GitHub 后会自动部署。
