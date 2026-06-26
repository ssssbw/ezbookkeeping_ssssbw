# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

ezBookkeeping - 轻量级、自托管的个人财务应用。前后端分离架构：
- **后端**：Go + Gin + xorm ORM，REST API
- **前端**：Vue 3 + TypeScript + Vite，支持桌面和移动双模式 UI

## 构建与开发命令

```bash
# 构建（Linux/macOS）
./build.sh backend              # 构建后端二进制
./build.sh frontend             # 构建前端
./build.sh package -o out.tar.gz  # 打包
./build.sh docker               # Docker 镜像

# 构建（Windows）
.\build.bat backend              # 构建后端二进制
.\build.bat frontend             # 构建前端
.\build.bat package -o out.zip   # 打包

# 跳过检查（构建时默认执行 lint 和 test）
# Linux/macOS:
./build.sh backend --no-lint --no-test
# Windows:
.\build.bat backend --no-lint --no-test

# 前端开发
npm run serve                   # 开发服务器
npm run build                   # 生产构建
npm run lint                    # ESLint + vue-tsc 类型检查

# 测试
go test ./pkg/...               # 后端测试
npm run test                    # 前端测试（vitest）
go test ./pkg/services/ -run TestXxx  # 运行单个后端测试
npx vitest run -t "test name"   # 运行单个前端测试
```

## 代码架构

### 后端结构（Go）

- `cmd/` - CLI 命令实现（使用 urfave/cli/v3）
  - `webserver.go` - Web 服务器启动和 API 路由注册
  - `database.go` - 数据库迁移命令
- `pkg/` - 核心业务逻辑
  - `api/` - REST API 处理器，按资源分文件
  - `models/` - 数据模型（xorm struct tags）
  - `services/` - 业务逻辑服务层
  - `converters/` - 数据导入导出格式转换（CSV、OFX、QIF、Beancount 等）
  - `auth/` - JWT、2FA、OIDC 认证
  - `llm/` - 大语言模型集成
  - `mcp/` - Model Context Protocol 支持

### 前端结构（Vue 3）

- `src/views/desktop/` - 桌面版页面组件
- `src/views/mobile/` - 移动版页面组件
- `src/views/base/` - 通用基础组件
- `src/stores/` - Pinia 状态管理
- `src/models/` - TypeScript 数据模型
- `src/lib/` - 工具函数
- `src/locales/` - 国际化文件（18 种语言）

入口文件：
- `src/desktop-main.ts` - 桌面版入口
- `src/mobile-main.ts` - 移动版入口

### 投资模块（Investment）

投资模块采用前后端分离架构，前端页面位于 `src/views/desktop/investment/`，后端逻辑位于 `pkg/` 下多个目录。

**前端页面**（路由 `/investment/*`）：
- `OverviewPage.vue` - 投资组合总览（资产配置饼图、组合摘要、月度表现）
- `PortfolioPage.vue` - 投资组合详情（汇总卡片、资产配置、持仓明细）
- `AssetsPage.vue` - 资产管理主页（持仓列表 Tab + 自选列表 Tab）
- `AssetDetailPage.vue` - 资产详情页（基本信息、持仓明细、价格走势、交易记录）
- `TransactionsPage.vue` - 交易记录（列表、新增交易弹窗）
- `AnalysisPage.vue` - 收益分析（收益率、资产类别表现、持仓收益率柱状图）
- `StrategyPage.vue` - 策略配置（目标配置、偏离分析、再平衡建议）

**前端子组件**：
- `assets/HoldingsTab.vue` - 持仓列表（支持展开行查看账户明细）
- `assets/WatchlistTab.vue` - 自选列表
- `assets/AssetSearchBar.vue` - 资产搜索栏
- `assets/HoldingsSummaryBar.vue` - 持仓汇总栏
- `components/TransactionFormDialog.vue` - 交易表单弹窗
- `components/MarketDataEditDialog.vue` - 行情编辑弹窗
- `components/InvestmentReturnOverviewCard.vue` - 收益概览卡片
- `assets/dialogs/AssetAdminDialog.vue` - 资产管理弹窗

**前端工具**：
- `assets/assetUtils.ts` - 格式化工具（价格、收益率、市场、类别等）
- `assets/types.ts` - DisplayAsset 类型定义

**后端结构**：
- `pkg/api/investment.go` - 所有投资相关 API handler（全局资产、用户资产、交易、行情、分析）
- `pkg/models/investment_transaction.go` - 交易模型
- `pkg/models/investment_analysis.go` - 分析模型（持仓、总览、配置）
- `pkg/models/investment_consts.go` - 投资常量（类别、市场、交易类型枚举）
- `pkg/models/market_data.go` - 行情数据模型（含 IsManual 手动标记）
- `pkg/services/investment_transaction.go` - 交易服务
- `pkg/services/investment_analysis.go` - 分析服务
- `pkg/services/market_data.go` - 行情数据服务（自动拉取跳过手动记录）
- `pkg/marketdata/` - 行情数据提供者（东方财富、AKShare）

**前端 Store** (`src/stores/investment.ts`)：
- `useInvestmentStore` - 管理用户资产、持仓、交易、行情、总览数据
- 核心 actions：`loadHoldings`、`loadOverview`、`loadTransactions`、`loadUserAssets`
- 核心 computed：`aggregatedHoldings`（按资产聚合）、`holdingsSummary`（汇总统计）

**前端 API 服务** (`src/lib/services.ts`)：
- 全局资产 CRUD：`searchAssets`、`getGlobalAsset`、`addGlobalAsset`、`modifyGlobalAsset`、`deleteGlobalAsset`
- 用户资产：`getUserAssets`、`addUserAsset`、`removeUserAsset`
- 交易 CRUD：`getInvestmentTransactions`、`addInvestmentTransaction`、`modifyInvestmentTransaction`、`deleteInvestmentTransaction`
- 行情数据：`getLatestMarketData`、`getMarketDataList`、`addMarketData`、`modifyMarketData`、`refreshMarketData`、`initMarketData`、`estimateMarketData`
- 分析：`getInvestmentHoldings`、`getInvestmentOverview`

### 数据库

支持 SQLite、MySQL、PostgreSQL。配置文件 `conf/ezbookkeeping.ini`。

## 关键技术细节

- ORM 使用 xorm，模型定义在 `pkg/models/`，通过 struct tags 映射数据库字段
- API 路由在 `cmd/webserver.go` 中统一注册
- 前端状态通过 Pinia stores 管理，stores 位于 `src/stores/`
- 前端 UI 框架使用 Vuetify 3 + Framework7（移动版）
- 图表使用 ECharts（全局注册为 `<v-chart>`，使用 `autoresize` 和 `:option` prop）
- ECharts 不支持 CSS 变量，必须使用硬编码颜色值（如 `#c67e48`）

### 投资模块技术要点

- **数值精度**：所有金额/价格/数量存储时 ×10000（DIVISOR=10000），前端显示时 ÷10000
- **价格格式化**：<1000 用 `toFixed(4)`，≥1000 用 `toFixed(2)`
- **红涨绿跌**：正收益用 `text-profit`（红色 `#d43f3f`），负收益用 `text-loss`（绿色 `#009688`）
- **手动价格保护**：MarketData 有 `IsManual` 字段，手动录入的价格不会被自动拉取覆盖
- **持仓聚合**：`aggregatedHoldings` 按 assetId 聚合多个账户的持仓，计算加权收益率
- **行情数据源**：东方财富（`eastmoney_market_data_provider.go`）和 AKShare（`akshare_market_data_provider.go`）
- **主色调**：投资模块使用 `#c67e48`（rgb 198, 126, 72）作为图表主色
