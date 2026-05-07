# ezBookkeeping 理财模块开发计划与设计建议

> 生成日期：2026-05-07
> 基于：核心需求.txt、FINANCIAL_PLANNING_MODULE_REQUIREMENTS.md、现有代码库上下文

---

## 一、现状分析

### 1.1 已具备的基础设施

| 层面 | 已完成 | 状态 |
|------|--------|------|
| **前端路由** | 6 条 investment 路由已配置（desktop.ts，lazy loaded） | ✅ |
| **前端导航** | MainLayout.vue 中理财模式侧边栏已完整设计 | ✅ |
| **模式切换** | `isInvestmentMode` ref + `toggleMode()` 已实现，点击 logo 切换 | ✅ |
| **前端占位页面** | OverviewPage / PortfolioPage / AnalysisPage / AssetsPage / TransactionsPage / StrategyPage 6 个 .vue 文件 | ✅ 仅有骨架 |
| **前端组件** | InvestmentReturnOverviewCard.vue | ✅ 仅有架子 |
| **Logo 资源** | ezinvestment-192.png 已引用 | ✅ |
| **国际化 Key** | `global.app.investmentTitle` 及部分投资相关 key 已定义 | ⚠️ 需补全 |

### 1.2 尚未实现的缺失项

| 层面 | 缺失项 |
|------|--------|
| **数据模型（Go）** | 无任何 InvestmentAsset / InvestmentTransaction / InvestmentPortfolio 等 Go struct |
| **数据库表** | 无对应 MySQL/SQLite 表 |
| **后端 API** | 无任何 `/api/v1/investment-*` 端点 |
| **前端 Store** | 无 investment 相关的 Pinia store |
| **前端 API 调用** | `src/lib/services.ts` 无投资相关 API 方法 |
| **移动端** | mobile router 无投资路由，无投资页面 |
| **数据兼容** | 现有 Transaction 模型与投资资产无关联机制 |

### 1.3 项目技术栈速览

| 层 | 技术 |
|----|------|
| 前端框架 | Vue 3 + TypeScript + Vite |
| 桌面 UI | Vuetify 3 |
| 移动 UI | Framework7 |
| 图表 | ECharts 6 |
| 状态管理 | Pinia |
| 后端 | Go 1.25 + Gin |
| ORM | XORM |
| 数据库 | SQLite3 / MySQL / PostgreSQL |
| 迁移方式 | 代码驱动（`SyncStructs`），无 SQL 文件 |

---

## 二、整体架构设计建议

### 2.1 核心理念：记账与理财解耦但共享基础设施

```
┌──────────────────────────────────────────────────┐
│              MainLayout.vue                       │
│  ┌──────────────┐    ┌──────────────────────────┐│
│  │ 记账模式      │    │ 理财模式                  ││
│  │ (isInvestment │    │ (isInvestment             ││
│  │  Mode=false)  │    │  Mode=true)               ││
│  │              │    │                          ││
│  │ Transaction  │    │ InvestmentAsset          ││
│  │ Account      │    │ InvestmentTransaction    ││
│  │ Category     │    │ InvestmentPortfolio      ││
│  │ Tag          │    │ InvestmentStrategy       ││
│  │              │    │ MarketData (历史数据)      ││
│  └──────┬───────┘    └───────────┬──────────────┘│
│         │                        │               │
│         └────────┬───────────────┘               │
│                  │                               │
│     ┌────────────▼──────────────┐               │
│     │   共享基础设施             │               │
│     │  - 用户认证 (JWT)          │               │
│     │  - 货币/汇率               │               │
│     │  - 国际化                  │               │
│     │  - 账户体系 (可复用)        │               │
│     │  - 存储服务                │               │
│     └───────────────────────────┘               │
└──────────────────────────────────────────────────┘
```

**建议：独立的投资交易表，而非扩展现有 Transaction 表。**

理由：
1. 投资交易（买入/卖出/分红/拆分）的字段结构与普通生活交易完全不同（数量、单价、手续费、基金代码等）
2. 查询模式不同——投资需按产品聚合、按时间序列分析收益率，普通交易按账户/分类聚合
3. 避免现有 Transaction 表被污染，影响记账功能性能
4. 通过 `investment_account_id` 关联到现有 Account 表（复用投资账户概念），实现资金层面的关联

### 2.2 数据兼容性策略

针对核心需求第 2、4 条：**采用 EAV（Entity-Attribute-Value）+ 类型表模式**。

```
investment_assets (主表)
├── id, uid, type (fund/stock/bond/...), name, code, currency
├── 通用字段：current_price, current_value, cost_basis, quantity
└── extra_info (JSON) — 不同类型特有字段
    ├── fund: { fund_manager, fund_type, ... }
    ├── stock: { exchange, sector, ... }
    └── bond: { maturity_date, coupon_rate, ... }

market_data (历史行情数据)
├── asset_id, date, price, volume
└── 解决 API 失效后仍可查询历史数据（核心需求第 4 条）

investment_transactions
├── asset_id, type(buy/sell/dividend/split), quantity, price, fee, account_id
└── 通过 account_id 关联到 accounts 表（复用现有账户体系）
```

**新增资产类型只需**：在 `type` 枚举中增加新值，在 `extra_info` 中存储新类型的特有字段即可。无需改表结构。

---

## 三、数据库设计（Go Models）

### 3.1 新增表清单

| 表名 | 用途 | 对应 Go Model |
|------|------|---------------|
| `investment_asset` | 投资资产主表（基金/股票/债券等） | `InvestmentAsset` |
| `investment_transaction` | 投资交易记录 | `InvestmentTransaction` |
| `market_data` | 历史行情数据（日线） | `MarketData` |
| `investment_portfolio` | 投资组合 | `InvestmentPortfolio` |
| `portfolio_asset` | 组合-资产关联（多对多） | `PortfolioAsset` |
| `investment_strategy` | 投资策略 | `InvestmentStrategy` |
| `investment_alert` | 策略提醒 | `InvestmentAlert` |

### 3.2 核心 Model 设计

#### 3.2.1 InvestmentAsset（投资资产）

```go
// pkg/models/investment_asset.go
type InvestmentAsset struct {
    AssetId         int64   `xorm:"PK"`
    Uid             int64   `xorm:"INDEX(IDX_asset_uid_deleted) NOT NULL"`
    Deleted         bool    `xorm:"INDEX(IDX_asset_uid_deleted) NOT NULL"`
    Type            string  `xorm:"VARCHAR(32) NOT NULL"`   // fund, stock, bond, etc.
    Name            string  `xorm:"VARCHAR(128) NOT NULL"`
    Code            string  `xorm:"VARCHAR(32)"`            // 基金代码/股票代码
    Currency        string  `xorm:"VARCHAR(3) NOT NULL"`
    CurrentPrice    int64   `xorm:"NOT NULL"`               // 最新单价(金额*10000存储)
    CostBasis       int64   `xorm:"NOT NULL"`               // 成本(金额*10000存储)
    Quantity        int64   `xorm:"NOT NULL"`               // 持有数量(数量*10000存储)
    IsActive        bool    `xorm:"NOT NULL"`
    ExtraInfo       string  `xorm:"TEXT"`                   // JSON, 不同类型特有字段
    Comment         string  `xorm:"VARCHAR(255)"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

> **关键设计决策**：`CurrentPrice` 和 `CostBasis` 使用 `int64` 存储（乘以 10000），与项目现有 `Transaction.Amount` 的精度处理方式一致。

#### 3.2.2 InvestmentTransaction（投资交易）

```go
// pkg/models/investment_transaction.go
type InvestmentTransaction struct {
    TransactionId      int64   `xorm:"PK"`
    Uid                int64   `xorm:"INDEX(IDX_inv_tx_uid_deleted) NOT NULL"`
    Deleted            bool    `xorm:"INDEX(IDX_inv_tx_uid_deleted) NOT NULL"`
    AssetId            int64   `xorm:"INDEX NOT NULL"`       // 关联 investment_asset
    AccountId          int64   `xorm:"INDEX NOT NULL"`       // 关联 accounts（资金账户）
    Type               string  `xorm:"VARCHAR(16) NOT NULL"` // buy/sell/dividend/split/interest
    Quantity           int64   `xorm:"NOT NULL"`             // 数量*10000
    Price              int64   `xorm:"NOT NULL"`             // 单价*10000
    Amount             int64   `xorm:"NOT NULL"`             // 交易总金额*10000
    Fee                int64   `xorm:"NOT NULL"`             // 手续费*10000
    TransactionTime    int64   `xorm:"INDEX NOT NULL"`
    TimezoneUtcOffset  int16   `xorm:"NOT NULL"`
    Comment            string  `xorm:"VARCHAR(255)"`
    CreatedUnixTime    int64
    UpdatedUnixTime    int64
    DeletedUnixTime    int64
}
```

#### 3.2.3 MarketData（历史行情数据）

```go
// pkg/models/market_data.go
type MarketData struct {
    DataId     int64   `xorm:"PK"`
    AssetId    int64   `xorm:"UNIQUE(UQE_market_data_asset_date) INDEX NOT NULL"`
    Date       int64   `xorm:"UNIQUE(UQE_market_data_asset_date) NOT NULL"` // 日期 unix(取0点)
    Price      int64   `xorm:"NOT NULL"`  // 当日净值/收盘价*10000
    Volume     int64   `xorm:"NOT NULL"`  // 成交量(可选)
    CreatedUnixTime int64
}
```

> **解决核心需求第 4 条**：每日 cron 任务拉取第三方 API 数据写入此表。即使 API 失效，历史数据仍在本地，不影响查询。

---

## 四、后端 API 设计

### 4.1 API 路由规划

在 `pkg/api/` 下新建 `investment.go`，遵循现有的 API 注册模式（参考 `transactions.go`、`accounts.go`）。

| 方法 | 路径 | 功能 |
|------|------|------|
| GET | `/api/v1/investment/assets/list.json` | 获取资产列表 |
| GET | `/api/v1/investment/assets/:id/get.json` | 获取资产详情 |
| POST | `/api/v1/investment/assets/add.json` | 创建资产 |
| POST | `/api/v1/investment/assets/:id/update.json` | 更新资产 |
| POST | `/api/v1/investment/assets/:id/delete.json` | 删除资产 |
| GET | `/api/v1/investment/assets/:id/history.json` | 获取历史行情 |
| GET | `/api/v1/investment/transactions/list.json` | 获取投资交易列表 |
| POST | `/api/v1/investment/transactions/add.json` | 创建投资交易 |
| POST | `/api/v1/investment/transactions/:id/update.json` | 更新投资交易 |
| POST | `/api/v1/investment/transactions/:id/delete.json` | 删除投资交易 |
| GET | `/api/v1/investment/analysis/overview.json` | 投资概览数据 |
| GET | `/api/v1/investment/analysis/performance.json` | 收益表现分析 |
| GET | `/api/v1/investment/analysis/allocation.json` | 资产配置分析 |
| GET | `/api/v1/investment/analysis/report.json` | 月报数据 |

### 4.2 遵循现有 API 模式

参考 `pkg/api/transactions.go` 中的模式：
- 每个 API 模块使用单例模式（`GetInvestmentApi()`）
- Handler 签名：`func (a *InvestmentApi) HandlerName(c *gin.Context)`
- 使用 `binding.Bind(c, &req)` 进行参数绑定
- 使用 `response.SendSuccess(c, result)` 返回成功
- 使用 `response.SendError(c, errCode, errMsg)` 返回错误

---

## 五、前端开发指南

### 5.1 目录结构规划

```
src/
├── views/desktop/investment/
│   ├── OverviewPage.vue          # 资产概览（总资产、收益、饼图）
│   ├── PortfolioPage.vue         # 投资组合管理
│   ├── AnalysisPage.vue          # 收益分析（日/周/月/年）
│   ├── AssetsPage.vue            # 资产管理（CRUD）
│   ├── TransactionsPage.vue      # 交易记录
│   ├── StrategyPage.vue          # 策略配置 & 提醒
│   └── components/
│       ├── InvestmentReturnOverviewCard.vue  # 收益概览卡片（已有骨架）
│       ├── AssetAllocationChart.vue          # 资产配置饼图
│       ├── PerformanceTrendChart.vue         # 收益趋势折线图
│       ├── AssetFormDialog.vue               # 资产编辑弹窗
│       ├── TransactionFormDialog.vue         # 交易编辑弹窗
│       └── MarketDataChart.vue              # 历史行情图
├── models/
│   ├── investment_asset.ts       # 前端投资资产类型定义
│   └── investment_transaction.ts # 前端投资交易类型定义
├── stores/
│   ├── investmentAsset.ts        # Pinia store for assets
│   ├── investmentTransaction.ts  # Pinia store for transactions
│   └── investmentAnalysis.ts     # Pinia store for analysis data
└── lib/
    └── services.ts               # 新增投资相关 API 方法
```

### 5.2 遵循现有前端模式

参考现有 `src/stores/account.ts` 的 Pinia store 模式：

```typescript
// src/stores/investmentAsset.ts (示例结构)
export const useInvestmentAssetsStore = defineStore('investmentAssets', () => {
    const allAssets = ref<InvestmentAsset[]>([]);
    const allAssetsMap = ref<Record<string, InvestmentAsset>>({});
    const assetListStateInvalid = ref<boolean>(true);
    const loading = ref<boolean>(false);

    const allActiveAssets = computed(() => allAssets.value.filter(a => a.isActive));

    function loadAllAssets(): Promise<void> { ... }
    function saveAsset(asset: InvestmentAssetCreateRequest): Promise<void> { ... }
    function updateAsset(assetId: number, asset: InvestmentAssetModifyRequest): Promise<void> { ... }
    function deleteAsset(assetId: number): Promise<void> { ... }
    function resetAssets(): void { ... }

    return { allAssets, allAssetsMap, assetListStateInvalid, loading, allActiveAssets,
             loadAllAssets, saveAsset, updateAsset, deleteAsset, resetAssets };
});
```

### 5.3 国际化

在 `src/locales/zh_Hans.json` 和 `en.json` 中新增约 50-80 个投资相关 key。Key 命名遵循现有风格：
- `investment.overview.title`
- `investment.asset.type.fund`
- `investment.transaction.type.buy`
- `investment.analysis.performance.daily`

### 5.4 ECharts 图表

Investment OverviewPage 已有 `v-chart` 骨架，遵循该项目使用 ECharts 的模式（参考 `StatisticsTransactionPage.vue`）：
- 使用 `vue-echarts` 的 `<v-chart>` 组件
- 图表配置通过 computed 属性生成
- 支持深色/浅色主题自适应（复用 `useTheme()` 机制）

---

## 六、基金数据源集成方案（akshare）

### 6.1 背景与需求

- 需要免费的基金/股票行情数据 API
- [akshare](https://github.com/akfamily/akshare) 是最成熟的 Python 金融数据库，封装了东方财富、新浪、腾讯等免费 API
- 但 akshare 是 Python 库，本项目后端是 Go
- 核心需求第 4 条：**每日行情数据需存储到本地，防止免费 API 失效**

### 6.2 三种集成方案对比

| 方案 | 复杂度 | 可靠性 | 维护性 | 需要 Python | Docker 额外容器 | 推荐场景 |
|------|--------|--------|--------|-------------|----------------|----------|
| **A: AKTools HTTP Sidecar** | ⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ | ✅ | ✅ | 长期生产部署 |
| **B: exec.Command 子进程** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ | ✅ | ❌ | 个人项目 MVP |
| **C: Go 原生库替代** | ⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ❌ | ❌ | 最简单部署 |

---

### 6.3 方案 A：AKTools HTTP Sidecar（推荐生产使用）

#### 原理

akshare 官方提供了 [AKTools](https://aktools.akfamily.xyz/)，一条命令将 akshare 包装成 REST HTTP 服务：

```bash
pip install aktools
python -m aktools --host 0.0.0.0 --port 8888
```

Go 后端通过 HTTP 调用获取数据。

#### 架构图

```
┌─────────────────────────────────────┐
│         Docker / 宿主机               │
│                                      │
│  ┌──────────────┐  ┌─────────────┐ │
│  │ Go 后端       │  │ AKTools      │ │
│  │ :8080        │──│ :8888        │ │
│  │              │  │ (Python)     │ │
│  └──────────────┘  └──────┬──────┘ │
│                           │         │
│                    HTTP GET         │
│                    /api/public/     │
│                    fund_open_fund_  │
│                    info_em          │
└─────────────────────────────────────┘
```

#### akshare 核心 API 能力

| API 函数 | 数据 | 更新频率 |
|----------|------|----------|
| `fund_name_em()` | 所有基金基本信息（代码、名称、类型） | 实时 |
| `fund_open_fund_daily_em()` | 所有开放式基金净值 | 每日 16:00-23:00 |
| `fund_open_fund_info_em(symbol="110022")` | 指定基金历史净值 | 全部历史 |
| `fund_etf_category_sina()` | ETF 列表及实时报价 | 实时 |
| `stock_zh_a_hist(symbol="000001")` | A 股历史 K 线 | 日/周/月 |

> **底层原理**：akshare 本质是封装东方财富/新浪/腾讯的免费公开 HTTP API，非爬虫。

#### 在 ezBookkeeping 中的集成方式

参考现有汇率数据源架构（`pkg/exchangerates/`），创建 `pkg/marketdata/`：

```
pkg/marketdata/
├── market_data_provider.go            # 定义 MarketDataProvider interface
├── market_data_provider_container.go  # 单例容器，根据配置选择数据源
├── aktools_data_source.go             # AKTools HTTP 数据源实现
├── native_data_source.go              # Go 原生数据源实现（方案 C）
└── subprocess_data_source.go          # exec.Command 数据源实现（方案 B）
```

**接口定义**（参考 `pkg/exchangerates/exchange_rates_data_provider.go` 的设计模式）：

```go
// pkg/marketdata/market_data_provider.go
package marketdata

import (
    "github.com/mayswind/ezbookkeeping/pkg/core"
    "github.com/mayswind/ezbookkeeping/pkg/models"
)

// MarketDataProvider 定义行情数据源接口
type MarketDataProvider interface {
    // GetFundList 获取所有基金基本信息
    GetFundList(c core.Context) ([]*models.FundBasicInfo, error)
    // GetFundNAV 获取指定基金的最新净值
    GetFundNAV(c core.Context, fundCode string) (*models.FundNAVData, error)
    // GetFundHistory 获取指定基金的历史净值
    GetFundHistory(c core.Context, fundCode string, startDate, endDate int64) ([]*models.FundNAVData, error)
    // GetStockHistory 获取股票历史行情
    GetStockHistory(c core.Context, stockCode string, startDate, endDate int64) ([]*models.StockDailyData, error)
}
```

**AKTools 实现**：

```go
// pkg/marketdata/aktools_data_source.go
package marketdata

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"

    "github.com/mayswind/ezbookkeeping/pkg/core"
    "github.com/mayswind/ezbookkeeping/pkg/log"
)

type AKToolsDataSource struct {
    baseURL    string // http://localhost:8888
    httpClient *http.Client
}

func NewAKToolsDataSource(baseURL string) *AKToolsDataSource {
    return &AKToolsDataSource{
        baseURL:    baseURL,
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

func (a *AKToolsDataSource) GetFundNAV(c core.Context, fundCode string) (*models.FundNAVData, error) {
    url := fmt.Sprintf("%s/api/public/fund_open_fund_info_em?symbol=%s&indicator=单位净值走势", a.baseURL, fundCode)

    resp, err := a.httpClient.Get(url)
    if err != nil {
        log.Errorf(c, "[aktools.GetFundNAV] request failed: %s", err.Error())
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // 解析 AKTools 返回的 JSON 数据
    var result []map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    // 转换为内部模型...
    return navData, nil
}
```

**容器初始化**（参考 `pkg/exchangerates/exchange_rates_data_provider_container.go`）：

```go
// pkg/marketdata/market_data_provider_container.go
package marketdata

import "github.com/mayswind/ezbookkeeping/pkg/settings"

var Container = &MarketDataProviderContainer{}

func InitializeMarketDataSource(config *settings.Config) error {
    switch config.MarketDataSource {
    case settings.AKToolsMarketDataSource:
        Container.current = NewAKToolsDataSource(config.AKToolsURL)
        return nil
    case settings.NativeMarketDataSource:
        Container.current = NewNativeDataSource()
        return nil
    case settings.SubprocessMarketDataSource:
        Container.current = NewSubprocessDataSource()
        return nil
    }
    return errs.ErrInvalidMarketDataSource
}
```

**注册 Cron 任务**（参考 `pkg/cron/cron_jobs.go` 和 `pkg/cron/cron_container.go`）：

```go
// 在 pkg/cron/cron_jobs.go 中新增
var UpdateFundMarketDataJob = &CronJob{
    Name:        "UpdateFundMarketData",
    Description: "Periodically update fund NAV data from market data source.",
    Period: CronJobFixedHourPeriod{
        Hour: 18,  // 每天 18:00 执行（基金净值一般在 16:00-23:00 发布）
    },
    Run: func(c *core.CronContext) error {
        return services.MarketData.UpdateAllActiveFundsNAV(c)
    },
}

// 在 pkg/cron/cron_container.go 的 registerAllJobs 中新增：
if config.EnableUpdateMarketData {
    Container.registerIntervalJob(ctx, UpdateFundMarketDataJob)
}
```

**Docker 部署**：

```yaml
# docker-compose.yml
services:
  ezbookkeeping:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MARKET_DATA_SOURCE=aktools
      - AKTOOLS_URL=http://aktools:8888
    depends_on:
      - aktools

  aktools:
    image: python:3.14-slim
    command: >
      sh -c "pip install aktools && python -m aktools --host 0.0.0.0 --port 8888"
    ports:
      - "8888:8888"
    restart: unless-stopped
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| akshare 生态完整，84+ 基金 API | 需要额外 Python 容器 |
| AKTools 一条命令启动，零代码 | Docker 部署多一个容器 |
| Go 和 Python 完全解耦 | 网络调用有 ~0.35ms 延迟 |
| akshare 更新后 `pip install` 即可 | 需要维护 Python 环境 |

---

### 6.4 方案 B：exec.Command 子进程（推荐 MVP）

#### 原理

Go 后端通过 `os/exec` 直接调用 Python 脚本，脚本内部使用 akshare 获取数据并输出 JSON。

#### 实现方式

**Python 脚本**（放在项目根目录 `scripts/fetch_fund_nav.py`）：

```python
#!/usr/bin/env python3
"""通过 akshare 获取基金净值数据，输出 JSON 到 stdout"""

import sys
import json
import akshare as ak

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "missing fund_code"}))
        sys.exit(1)

    fund_code = sys.argv[1]
    mode = sys.argv[2] if len(sys.argv) > 2 else "latest"

    if mode == "latest":
        # 获取最新净值
        data = ak.fund_open_fund_info_em(symbol=fund_code, indicator="单位净值走势")
        # 取最后一行（最新数据）
        latest = data.iloc[-1]
        result = {
            "fund_code": fund_code,
            "date": str(latest.iloc[0]),
            "nav": float(latest.iloc[1]),
            "accumulated_nav": float(latest.iloc[2]),
            "daily_growth_rate": float(latest.iloc[3]) if len(latest) > 3 else 0,
        }
    elif mode == "history":
        # 获取全部历史净值
        data = ak.fund_open_fund_info_em(symbol=fund_code, indicator="单位净值走势")
        records = []
        for _, row in data.iterrows():
            records.append({
                "date": str(row.iloc[0]),
                "nav": float(row.iloc[1]),
                "accumulated_nav": float(row.iloc[2]),
            })
        result = {"fund_code": fund_code, "history": records}

    print(json.dumps(result, ensure_ascii=False))

if __name__ == "__main__":
    main()
```

**Go 数据源实现**：

```go
// pkg/marketdata/subprocess_data_source.go
package marketdata

import (
    "encoding/json"
    "fmt"
    "os/exec"
    "path/filepath"

    "github.com/mayswind/ezbookkeeping/pkg/core"
    "github.com/mayswind/ezbookkeeping/pkg/log"
)

type SubprocessDataSource struct {
    pythonPath string // "python3" 或 "python"
    scriptPath  string // 项目根目录/scripts/fetch_fund_nav.py
}

func NewSubprocessDataSource() *SubprocessDataSource {
    return &SubprocessDataSource{
        pythonPath: "python3",
        scriptPath:  filepath.Join("scripts", "fetch_fund_nav.py"),
    }
}

func (s *SubprocessDataSource) GetFundNAV(c core.Context, fundCode string) (*models.FundNAVData, error) {
    cmd := exec.Command(s.pythonPath, s.scriptPath, fundCode, "latest")

    output, err := cmd.Output()
    if err != nil {
        log.Errorf(c, "[subprocess.GetFundNAV] python script failed for fund %s: %s", fundCode, err.Error())
        return nil, fmt.Errorf("python script failed: %w", err)
    }

    var result struct {
        FundCode       string  `json:"fund_code"`
        Date           string  `json:"date"`
        Nav            float64 `json:"nav"`
        AccumulatedNav float64 `json:"accumulated_nav"`
        DailyGrowth    float64 `json:"daily_growth_rate"`
    }

    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("parse python output failed: %w", err)
    }

    // 转换为内部模型（金额 * 10000 存储）
    return &models.FundNAVData{
        FundCode:       result.FundCode,
        Date:           result.Date,
        Nav:            int64(result.Nav * 10000),
        AccumulatedNav: int64(result.AccumulatedNav * 10000),
        DailyGrowth:    int64(result.DailyGrowth * 10000),
    }, nil
}

func (s *SubprocessDataSource) GetFundHistory(c core.Context, fundCode string, startDate, endDate int64) ([]*models.FundNAVData, error) {
    cmd := exec.Command(s.pythonPath, s.scriptPath, fundCode, "history")

    output, err := cmd.Output()
    if err != nil {
        log.Errorf(c, "[subprocess.GetFundHistory] python script failed: %s", err.Error())
        return nil, err
    }

    var result struct {
        FundCode string                `json:"fund_code"`
        History  []map[string]float64 `json:"history"`
    }

    if err := json.Unmarshal(output, &result); err != nil {
        return nil, err
    }

    var navData []*models.FundNAVData
    for _, item := range result.History {
        navData = append(navData, &models.FundNAVData{
            FundCode: result.FundCode,
            Date:     fmt.Sprintf("%.0f", item["date"]),
            Nav:      int64(item["nav"] * 10000),
        })
    }

    return navData, nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| **零基础设施**，写几行就能跑 | 每次调用启动 Python 解释器（~100-200ms 开销） |
| 调试简单，手动跑脚本即可验证 | 无连接复用，不适合高频调用 |
| Docker 部署无需额外容器 | 需要在 Go 容器中安装 Python |
| 代码最少，实现最快 | Python 进程崩溃会影响 Go 进程 |

> **对于个人项目每天定时拉一次数据，200ms 的启动开销完全可以接受。**

---

### 6.5 方案 C：Go 原生库（推荐最简部署）

#### 原理

直接使用 Go 编写的金融数据库，这些库和 akshare 一样封装了东方财富/新浪的免费 HTTP API，无需 Python。

#### 可用的 Go 原生库

| 库 | 功能覆盖 | 维护状态 | GitHub |
|----|---------|----------|--------|
| **Hanson/adata-golang** | ETF（4 个 API）、股票行情、基金基础信息 | 2026 年活跃维护 | https://github.com/Hanson/adata-golang |
| **BlakeLiAFK/akshare-go** | 1154+ 接口，包含 84 个基金 API | 2026 年活跃维护 | https://github.com/BlakeLiAFK/akshare-go |

#### 实现方式

```go
// pkg/marketdata/native_data_source.go
package marketdata

import (
    adata "github.com/Hanson/adata-golang/fund"
    // 或
    aksharego "github.com/BlakeLiAFK/akshare-go"
)

type NativeDataSource struct{}

func NewNativeDataSource() *NativeDataSource {
    return &NativeDataSource{}
}

func (n *NativeDataSource) GetFundNAV(c core.Context, fundCode string) (*models.FundNAVData, error) {
    // 使用 adata-golang
    service := adata.NewFundService()
    data, err := service.GetFundNAV(fundCode, "2024-01-01", "2024-12-31")

    // 或使用 akshare-go
    // data, err := aksharego.FundOpenFundInfoEm(fundCode, "单位净值走势")

    if err != nil {
        return nil, err
    }

    // 转换为内部模型...
    return navData, nil
}
```

#### 优缺点

| 优点 | 缺点 |
|------|------|
| **不需要 Python 运行时** | API 覆盖可能不如 akshare 完整 |
| **不需要额外 Docker 容器** | Go 库的维护活跃度不如 akshare |
| **部署最简单**，单一二进制 | 如果底层 API 变化，需要等 Go 库更新 |
| **延迟最低**，直接 HTTP 调用 | 部分 API 尚未实现 |
| 类型安全，无需 JSON 解析 | 社区规模较小 |

---

### 6.6 三种方案对比与选择建议

```
你想要什么？
│
├─ 最简单的部署（不要 Python）
│  └─ 方案 C：Go 原生库（adata-golang）
│     - 单二进制部署，零依赖
│     - 适合长期维护
│
├─ 最快实现 MVP（能用就行）
│  └─ 方案 B：exec.Command
│     - 10 行 Go 代码 + 1 个 Python 脚本
│     - 适合快速验证
│
└─ 最完整的数据覆盖
   └─ 方案 A：AKTools Sidecar
      - akshare 84+ 基金 API 全部可用
      - 需要多一个 Docker 容器
```

**对个人项目的建议**：

1. **先方案 B 快速起步**：写一个 Python 脚本 + `exec.Command`，在 cron job 中每天拉一次数据。开发阶段最快。
2. **成熟后迁移到方案 C**：如果 Go 原生库的 API 覆盖满足需求，替换为 `NativeDataSource`。只需实现新的数据源类，接口不变，改动最小。
3. **如果需要完整 API**：使用方案 A 的 AKTools Sidecar。

**三种方案共享同一个接口（`MarketDataProvider`）**，切换只需改配置，无需改动业务逻辑。

---

## 七、实施计划

### 阶段 1：数据层（后端） — 预估 1-2 周

| 任务 | 详细 |
|------|------|
| 1.1 | 创建 `pkg/models/investment_asset.go`、`investment_transaction.go`、`market_data.go` 等 7 个 model |
| 1.2 | 在 `cmd/database.go` 的 `updateAllDatabaseTablesStructure()` 中注册新表 |
| 1.3 | 创建 `pkg/marketdata/` 包，实现 `MarketDataProvider` 接口 + 至少一种数据源 |
| 1.4 | 创建 `pkg/services/investment_asset.go`、`investment_transaction.go` 业务服务层 |
| 1.5 | 创建 `pkg/api/investment.go`，实现全部 REST 端点 |
| 1.6 | 在 `pkg/server/` 中注册投资路由组 |
| 1.7 | 在 `pkg/cron/cron_jobs.go` 中注册 `UpdateFundMarketDataJob` |
| 1.8 | 编写单元测试 |

### 阶段 2：前端 Store + API 层 — 预估 1 周

| 任务 | 详细 |
|------|------|
| 2.1 | 创建 `src/models/investment_asset.ts`、`investment_transaction.ts` 类型定义 |
| 2.2 | 创建 `src/stores/investmentAsset.ts`、`investmentTransaction.ts`、`investmentAnalysis.ts` |
| 2.3 | 在 `src/lib/services.ts` 中新增投资 API 方法 |
| 2.4 | 在 `src/stores/index.ts`（rootStore）中集成新 store |

### 阶段 3：理财前端界面 — 预估 2-3 周

| 任务 | 详细 |
|------|------|
| 3.1 | **OverviewPage**：总资产卡片、日/周/月/年收益概览、资产配置饼图、收益趋势折线图 |
| 3.2 | **AssetsPage**：资产列表、CRUD 弹窗、搜索/筛选 |
| 3.3 | **TransactionsPage**：交易列表、买入/卖出/分红记录表单 |
| 3.4 | **PortfolioPage**：组合管理、组合内资产权重调整 |
| 3.5 | **AnalysisPage**：日/周/月/年收益分析、收益率计算、波动率、最大回撤 |
| 3.6 | **StrategyPage**：策略预设、条件配置、提醒列表 |

### 阶段 4：数据导入导出 & 第三方 API — 预估 2 周

| 任务 | 详细 |
|------|------|
| 4.1 | **导入增强**：导入时自动创建缺失的账户/标签/分类（核心需求第 1 条） |
| 4.2 | **投资数据导入**：支持 CSV/JSON 格式导入基金持仓、交易记录 |
| 4.3 | **第三方 API 集成**：接入天天基金/东方财富等免费 API 拉取基金净值 |
| 4.4 | **Cron 任务**：每日自动拉取行情数据存入 `market_data` 表 |

### 阶段 5：移动端 & 打磨 — 预估 1 周

| 任务 | 详细 |
|------|------|
| 5.1 | 移动端投资路由注册 |
| 5.2 | 移动端 Overview 页面（卡片布局） |
| 5.3 | 国际化补全 |
| 5.4 | 模式切换动画优化 |
| 5.5 | 整体测试 |

---

## 八、设计建议与风险提醒

### 8.1 金额精度

项目现有 `Transaction.Amount` 使用 `int64` 存储（金额 × 10000）。新表必须保持一致。前端显示时除以 10000。

### 8.2 数据兼容性（核心需求第 2 条）

```
投资交易 (investment_transactions)  ←关联→  投资资产 (investment_assets)
       │                                        │
       │  account_id                            │ type=fund/stock/...
       ↓                                        │ extra_info (JSON)
账户 (accounts)                                  ↓
  type = INVESTMENT                             市场行情 (market_data)
                                      (历史净值/价格，本地存储)
```

**不需要修改现有 Transaction 表结构**。普通交易和投资交易通过 `account_id` 在账户层面关联（资金流向），但在业务层面完全独立。

### 8.3 API 失效应对（核心需求第 4 条）

- 每次成功拉取第三方数据时，写入 `market_data` 表
- 计算收益率、趋势图等优先使用本地数据
- 第三方 API 仅用于更新最新价格，不是实时依赖
- 如果 API 连续失败 3 天，通过策略提醒通知用户

### 8.4 移动端考量

目前 mobile router 完全没有投资路由。建议：
- 阶段 1-3 先完整实现桌面端
- 阶段 5 再移植到移动端（Framework7 组件风格不同，不能直接复用 Vuetify 组件）

### 8.5 模式切换

MainLayout.vue 中 `toggleMode()` 已实现基本切换，建议优化：
- 添加 CSS transition 动画（logo 旋转 + 页面 crossfade）
- `isInvestmentMode` 状态持久化到 `localStorage`，刷新页面后保持
- 建议：将 `isInvestmentMode` 存入 Pinia store（如 `desktopPage.ts`），而非组件内 ref

### 8.6 性能

- `market_data` 表会随时间快速增长（每个资产每天一条记录）
- 建议添加复合索引：`(asset_id, date DESC)`
- 历史数据查询始终限制时间范围（默认展示最近 1 年）
- 列表页使用分页，每页 20-50 条

---

## 九、验收标准

| 检查项 | 标准 |
|--------|------|
| 资产 CRUD | 创建、编辑、删除投资资产正常 |
| 交易记录 | 买入/卖出/分红/拆分四种交易类型均可记录 |
| 收益计算 | 日/周/月/年收益率计算正确 |
| 图表渲染 | 饼图、折线图数据正确，深色模式正常 |
| API 鉴权 | 所有投资 API 需要登录 Token |
| 数据隔离 | 用户 A 看不到用户 B 的投资数据 |
| 导入增强 | 导入时自动创建缺失账户/分类/标签 |
| 历史数据 | 每日行情数据正常存储和查询 |
| 数据源切换 | 三种数据源方案均可通过配置切换 |
| 移动端 | 投资概览在手机上正常显示 |
| 国际化 | 中英文界面完整无缺漏 |

---

## 十、后续可扩展方向

1. **更多资产类型**：加密货币、P2P、房产等——只需在 `type` 枚举 + `extra_info` JSON 中扩展
2. **智能定投**：基于预设策略自动计算定投金额和时机
3. **税务计算**：按国家/地区规则计算资本利得税
4. **家庭共享**：多个用户共享投资组合，权限控制（核心需求中提及）
5. **PDF 月报导出**：生成专业月度投资报告
6. **数据源热切换**：运行时切换 akshare/Go 原生库，无需重启

---

*本文档基于项目现有代码库上下文生成，所有设计遵循 ezBookkeeping 的现有编码规范和架构模式。*
