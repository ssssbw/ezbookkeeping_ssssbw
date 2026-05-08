# ezBookkeeping 理财模块开发计划与设计建议

> 生成日期：2026-05-06
> 最后更新：2026-05-08
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
| **数据兼容** | 现有 Transaction 模型与投资资产无关联机制 |
| **移动端 App** | 后期独立开发移动端 App（React Native/Flutter），本项目不再维护移动端网页 |

### 1.3 项目技术栈速览

| 层 | 技术 |
|----|------|
| 前端框架 | Vue 3 + TypeScript + Vite |
| 桌面 UI | Vuetify 3 |
| 图表 | ECharts 6 |
| 状态管理 | Pinia |
| 后端 | Go 1.25 + Gin |
| ORM | XORM |
| 数据库 | SQLite3 / MySQL / PostgreSQL |
| 迁移方式 | 代码驱动（`SyncStructs`），无 SQL 文件 |
| 注 | 不再支持移动端网页；移动端将独立开发 App |

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

### 4.3 行情数据拉取架构（方案 A：仿汇率拉取模式）

项目 `pkg/exchangerates/` 已有成熟的 Provider 模式，行情数据拉取完全遵循此架构：

```
pkg/marketdata/
├── market_data_provider.go              # MarketDataProvider 接口
├── common_http_market_data_provider.go  # HTTP 通用实现
├── eastmoney_fund_datasource.go         # 东方财富基金数据源
└── market_data_provider_container.go    # 单例容器
```

**免费 API 端点**：

| 接口 | URL |
|------|-----|
| 实时估值 | `https://fundgz.1234567.com.cn/js/{code}.js` |
| 基金详情（含全量历史净值） | `http://fund.eastmoney.com/pingzhongdata/{code}.js` |
| 全部基金列表 | `http://fund.eastmoney.com/js/fundcode_search.js` |
| 历史净值（分页） | `http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code={code}` |

### 4.4 高级数据拉取（方案 B：Go 原生 akshare）

引入 `github.com/BlakeLiAFK/akshare`（Python akshare 的 Go 移植，84 个基金接口）：

```go
import "github.com/BlakeLiAFK/akshare/fund"

// 基金基本信息
fund.FundNameEm()
fund.FundInfoEm("007345")

// 净值数据
fund.FundOpenFundInfoEm("007345", "单位净值走势", "成立来")

// 持仓数据
fund.FundPortfolioHoldEm("007345")
```

Go 原生库，无需 Python 环境，`go get` 即可集成。

### 4.5 Cron 定时任务

在 `pkg/cron/cron_jobs.go` 新增：

```go
var FetchFundMarketDataJob = &CronJob{
    Name:        "FetchFundMarketData",
    Description: "Daily fetch fund NAV data from Eastmoney and write to market_data table.",
    Period: CronJobFixedHourPeriod{Hour: 18}, // 每日 18:00（收盘后）
    Run: func(c *core.CronContext) error {
        return services.MarketData.FetchAllActiveFunds(c)
    },
}
```

在 `pkg/cron/cron_container.go` 的 `registerAllJobs` 中注册。

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

## 六、实施计划

> **开发顺序说明**：阶段 4 中的行情数据基础设施（market_data model、provider 架构、cron job）提前到阶段 1 完成，因为前端图表（阶段 3）依赖历史净值数据。阶段 4 仅保留导入增强和 akshare 高级数据。

### 阶段 1：数据层（后端）+ 行情基础设施 — 预估 2 周

| 任务 | 详细 |
|------|------|
| 1.1 | 创建 `pkg/models/investment_asset.go`、`investment_transaction.go`、`market_data.go` 等 model |
| 1.2 | 在 `cmd/database.go` 的 `updateAllDatabaseTablesStructure()` 中注册新表 |
| 1.3 | 创建 `pkg/services/investment_asset.go`、`investment_transaction.go` 业务服务层 |
| 1.4 | 创建 `pkg/api/investment.go`，实现全部 REST 端点 |
| 1.5 | 在 `pkg/server/` 中注册投资路由组 |
| 1.6 | **行情数据 provider**：仿照 `pkg/exchangerates/` 建立 `pkg/marketdata/`，实现 `MarketDataProvider` 接口 + 东方财富 HTTP 数据源 |
| 1.7 | **Cron 任务**：新增 `FetchFundMarketDataJob`，每日 18:00 自动拉取活跃基金净值写入 `market_data` 表 |
| 1.8 | 编写单元测试 |

### 阶段 2：前端 Store + API 层 — 预估 1 周

| 任务 | 详细 |
|------|------|
| 2.1 | 创建 `src/models/investment_asset.ts`、`investment_transaction.ts` 类型定义 |
| 2.2 | 创建 `src/stores/investmentAsset.ts`、`investmentTransaction.ts`、`investmentAnalysis.ts` |
| 2.3 | 在 `src/lib/services.ts` 中新增投资 API 方法 |
| 2.4 | 在 `src/stores/index.ts`（rootStore）中集成新 store |

### 阶段 2.5：API 契约检查（前后对接验证）

> **在阶段 2 完成后、阶段 3 开始前执行。** 目的是提前发现后端数据缺口，避免前端开发过程中反复回改。

| 步骤 | 详细 |
|------|------|
| 2.5.1 | 逐页列出每个组件的数据需求（OverviewPage 需要什么字段、AssetsPage 需要什么查询维度……） |
| 2.5.2 | 对照阶段 1 的 API 响应结构，标记已满足 / 缺字段 / 缺端点 |
| 2.5.3 | 一次性补齐所有缺口（新增字段、新增端点、调整聚合逻辑），然后才开始阶段 3 |

### 阶段 3：理财前端界面 — 预估 2-3 周

| 任务 | 详细 |
|------|------|
| 3.1 | **OverviewPage**：总资产卡片、日/周/月/年收益概览、资产配置饼图、收益趋势折线图 |
| 3.2 | **AssetsPage**：资产列表、CRUD 弹窗、搜索/筛选 |
| 3.3 | **TransactionsPage**：交易列表、买入/卖出/分红记录表单 |
| 3.4 | **PortfolioPage**：组合管理、组合内资产权重调整 |
| 3.5 | **AnalysisPage**：日/周/月/年收益分析、收益率计算、波动率、最大回撤 |
| 3.6 | **StrategyPage**：策略预设、条件配置、提醒列表 |

### 阶段 4：导入增强 & 高级数据 — 预估 1-2 周

| 任务 | 详细 |
|------|------|
| 4.1 | **导入增强**：导入时自动创建缺失的账户/标签/分类（核心需求第 1 条） |
| 4.2 | **投资数据导入**：支持 CSV/JSON 格式导入基金持仓、交易记录 |
| 4.3 | **引入 BlakeLiAFK/akshare**：Go 原生 akshare 移植库（84 个基金接口），用于高级数据（持仓、分红、基本信息等） |

### 阶段 5：打磨 & 测试 — 预估 1 周

| 任务 | 详细 |
|------|------|
| 5.1 | 国际化补全 |
| 5.2 | 模式切换动画优化 |
| 5.3 | 整体测试 & bug 修复 |
| 5.4 | 性能优化（market_data 查询、分页、缓存） |

> **移动端 App**：独立项目开发（React Native 或 Flutter），通过本项目 API 交互。不在本项目中开发移动端网页。

---

## 七、设计建议与风险提醒

### 7.1 金额精度

项目现有 `Transaction.Amount` 使用 `int64` 存储（金额 × 10000）。新表必须保持一致。前端显示时除以 10000。

### 7.2 数据兼容性（核心需求第 2 条）

```
投资交易 (investment_transactions)  ←关联→  投资资产 (investment_assets)
       │                                        │
       │  account_id                            │ type=fund/stock/...
       ↓                                        │ extra_info (JSON)
账户 (accounts)                                  │
  type = INVESTMENT                             ↓
                                      市场行情 (market_data)
                                      (历史净值/价格，本地存储)
```

**不需要修改现有 Transaction 表结构**。普通交易和投资交易通过 `account_id` 在账户层面关联（资金流向），但在业务层面完全独立。

### 7.3 数据源可切换（策略模式）

行情数据拉取采用与汇率系统相同的**策略模式**，数据源失效时只需新增一个 Provider 实现并改配置，无需修改任何 API 端点或 cron job：

```
MarketDataProvider (接口)
├── EastMoneyHttpProvider    ← 方案 A：直接 HTTP 调东方财富（默认）
├── AkshareGoProvider        ← 方案 B：BlakeLiAFK/akshare
└── (未来可扩展)              ← Sina、Tencent、自定义 CSV...
```

切换方式：修改配置 `market_data_source = "eastmoney"` → `"akshare"`，重启即生效。

### 7.4 API 失效应对（核心需求第 4 条）

- 每次成功拉取第三方数据时，写入 `market_data` 表
- 计算收益率、趋势图等优先使用本地数据
- 第三方 API 仅用于更新最新价格，不是实时依赖
- 如果 API 连续失败 3 天，通过策略提醒通知用户

### 7.4 移动端考量

本项目不再支持移动端网页。移动端将作为独立 App 开发（React Native 或 Flutter），通过后端 API 与本项目交互。因此本文档仅覆盖桌面端（Web）的投资模块开发。

### 7.5 模式切换

MainLayout.vue 中 `toggleMode()` 已实现基本切换，建议优化：
- 添加 CSS transition 动画（logo 旋转 + 页面 crossfade）
- `isInvestmentMode` 状态持久化到 `localStorage`，刷新页面后保持
- 建议：将 `isInvestmentMode` 存入 Pinia store（如 `desktopPage.ts`），而非组件内 ref

### 7.6 性能

- `market_data` 表会随时间快速增长（每个资产每天一条记录）
- 建议添加复合索引：`(asset_id, date DESC)`
- 历史数据查询始终限制时间范围（默认展示最近 1 年）
- 列表页使用分页，每页 20-50 条

---

## 八、验收标准

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
| 国际化 | 中英文界面完整无缺漏 |

---

## 九、后续可扩展方向

1. **更多资产类型**：加密货币、P2P、房产等——只需在 `type` 枚举 + `extra_info` JSON 中扩展
2. **智能定投**：基于预设策略自动计算定投金额和时机
3. **税务计算**：按国家/地区规则计算资本利得税
4. **家庭共享**：多个用户共享投资组合，权限控制（核心需求中提及）
5. **PDF 月报导出**：生成专业月度投资报告

---

## 十、开发前技术准备

> 不需要系统学习 Go 和 Vue 全家桶，只需掌握本项目用到的具体模式。**直接对着参考文件写，遇到不会的再查。**

### 10.1 Go 后端（阶段 1 必备）

#### 🔴 必学：Go struct + XORM 标签

写 Model 时唯一需要掌握的语法。

**参考文件**：`pkg/models/transaction_category.go:17-32`

```go
type TransactionCategory struct {
    CategoryId       int64  `xorm:"PK"`                    // 主键
    Uid              int64  `xorm:"INDEX(...) NOT NULL"`    // 带索引
    Deleted          bool   `xorm:"INDEX(...) NOT NULL"`    // 软删除标记
    Type             byte   `xorm:"NOT NULL"`               // 枚举用 byte
    Name             string `xorm:"VARCHAR(64) NOT NULL"`   // 字符串 + 长度
    CreatedUnixTime  int64                                   // 创建时间
    UpdatedUnixTime  int64                                   // 更新时间
    DeletedUnixTime  int64                                   // 删除时间
}
```

**你要做的**：照着这个模式写 `InvestmentAsset`、`InvestmentTransaction`、`MarketData`，改字段名和类型就行。

#### 🔴 必学：Gin Handler 三步走

每个 API handler 固定三步流程。

**参考文件**：`pkg/api/transaction_categories.go:42-60`

```go
func (a *Api) Handler(c *core.WebContext) (any, *errs.Error) {
    // 1. 绑定请求参数
    var req models.SomeRequest
    err := c.ShouldBindQuery(&req)   // GET 用 ShouldBindQuery
    // err := c.ShouldBindJSON(&req) // POST 用 ShouldBindJSON

    // 2. 调 service 拿数据
    uid := c.GetCurrentUid()
    data, err := a.service.SomeMethod(c, uid, req)

    // 3. 返回结果（框架自动包装成 JSON）
    return data, nil
}
```

#### 🔴 必学：路由注册

**参考文件**：`cmd/webserver.go:379-385`

```go
apiV1Route.GET("/accounts/list.json", bindApi(api.Accounts.AccountListHandler))
apiV1Route.POST("/accounts/add.json", bindApi(api.Accounts.AccountCreateHandler))
```

**你要做的**：在旁边加一组 `/investment/` 路由。

#### 🟡 需学：Service 层 CRUD

**参考文件**：`pkg/services/transaction_categories.go`

```go
// 查询
s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=?", uid, false).Find(&categories)

// 插入
session.Insert(model)

// 更新
session.ID(id).Update(cols, model)

// 软删除
// 更新 deleted=true + deleted_unix_time
```

#### 🟡 需学：表注册

**参考文件**：`cmd/database.go:85-107`

```go
err = datastore.Container.UserDataStore.SyncStructs(new(models.InvestmentAsset))
err = datastore.Container.UserDataStore.SyncStructs(new(models.MarketData))
```

一行注册一张表，复制改 model 名即可。

#### 🟢 按需：Go 接口（Provider 模式）

阶段 1.6 行情数据架构需要。

**参考文件**：`pkg/exchangerates/exchange_rates_data_provider.go`、`exchange_rates_data_provider_container.go`

```go
type ExchangeRatesDataProvider interface {
    GetLatestExchangeRates(...) (*Response, error)
}
// 多个实现 → Container 按配置选择
```

### 10.2 Vue + TypeScript 前端（阶段 2-3 必备）

#### 🔴 必学：Pinia Store 模式

**参考文件**：`src/stores/account.ts:26-34`

```typescript
export const useXxxStore = defineStore('xxx', () => {
    // 状态
    const allItems = ref<Item[]>([]);
    const loading = ref<boolean>(false);

    // 计算属性
    const activeItems = computed(() => allItems.value.filter(i => i.isActive));

    // 加载数据
    function loadAll() { services.getXxx().then(...) }

    return { allItems, loading, activeItems, loadAll };
});
```

#### 🔴 必学：API 调用

**参考文件**：`src/lib/services.ts:490-510`

```typescript
// GET
getAllAccounts: ({ visibleOnly }) =>
    axios.get<ApiResponse<AccountInfoResponse[]>>('v1/accounts/list.json?visible_only=' + visibleOnly),

// POST
addAccount: (req) =>
    axios.post<ApiResponse<AccountInfoResponse>>('v1/accounts/add.json', req),
```

**你要做的**：照着加一组 `investment/` 的 API 方法。

#### 🔴 必学：TypeScript 类型定义

**参考文件**：`src/models/account.ts`、`src/models/transaction.ts`

```typescript
export interface AccountCreateRequest {
    name: string;
    type: number;
    currency: string;
    // ...
}
```

#### 🟢 按需：Vuetify 组件

用到什么查什么，最常用的：
- `v-card`、`v-list`、`v-data-table` — 列表页
- `v-dialog` — 弹窗
- `v-text-field`、`v-select`、`v-btn` — 表单

#### 🟢 按需：ECharts 图表

阶段 3 画图表时再看。

**参考文件**：`StatisticsTransactionPage.vue`（`v-chart` 组件用法）

### 10.3 学习优先级总览

| 优先级 | 知识点 | 参考文件 | 对应阶段 |
|--------|--------|----------|----------|
| 🔴 必学 | Go struct + XORM 标签 | `pkg/models/transaction_category.go` | 1.1 |
| 🔴 必学 | Gin Handler 三步走 | `pkg/api/transaction_categories.go` | 1.4 |
| 🔴 必学 | 路由注册 | `cmd/webserver.go:379` | 1.5 |
| 🔴 必学 | Pinia Store 模式 | `src/stores/account.ts` | 2.2 |
| 🔴 必学 | services.ts API 调用 | `src/lib/services.ts:490` | 2.3 |
| 🟡 需学 | XORM CRUD | `pkg/services/transaction_categories.go` | 1.3 |
| 🟡 需学 | 表注册 | `cmd/database.go:85` | 1.2 |
| 🟡 需学 | TS 类型定义 | `src/models/` | 2.1 |
| 🟢 按需 | Go 接口（Provider） | `pkg/exchangerates/` | 1.6 |
| 🟢 按需 | Vuetify 组件 | 官方文档 | 3.1-3.6 |
| 🟢 按需 | ECharts | `StatisticsTransactionPage.vue` | 3.1, 3.5 |

---

*本文档基于项目现有代码库上下文生成，所有设计遵循 ezBookkeeping 的现有编码规范和架构模式。*