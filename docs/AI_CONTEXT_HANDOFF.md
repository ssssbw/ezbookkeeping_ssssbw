# AI 会话上下文交接文档

> 每次会话结束后更新此文件，确保下一个 AI 会话能无缝接续。
> 最后更新：2026-06-04

---

## 一、项目概况

| 项 | 值 |
|----|-----|
| 项目名 | ezBookkeeping 投资理财模块 |
| 项目路径 | /Users/xw.zhu/ezbookkeeping_ssssbw |
| 分支 | feature/add_invest_analysis |
| 后端 | Go 1.25 + Gin + XORM |
| 前端 | Vue 3 + TypeScript + Vuetify 3 + Pinia |
| 运行配置 | package/conf/ezbookkeeping.ini（MySQL, 127.0.0.1:3306, user=root, passwd=245680） |
| 开发配置 | conf/ezbookkeeping.ini（SQLite3, data/ezbookkeeping.db） |
| 数据库迁移 | 代码驱动 SyncStructs，无 SQL 文件 |
| 路由注册 | cmd/webserver.go |
| 表注册 | cmd/database.go |

---

## 二、已完成的工作

### 2.1 文档（docs/）

| 文件 | 状态 | 说明 |
|------|------|------|
| INVESTMENT_MODULE_DEV_PLAN.md | ✅ 已完成 | 开发计划（含阶段划分、技术栈、学习指南） |
| ER_DIAGRAM.md | ✅ 已完成 | 现有 16 张表 + 新增投资模块 ER 图 |
| INVESTMENT_MODULE_DESIGN_NOTES.md | ✅ 已完成 | 设计决策笔记（三层资金流、自洽性验证、交易类型定义） |
| AI_CONTEXT_HANDOFF.md | ✅ 已完成 | AI 会话上下文交接文档 |

### 2.2 代码

| 文件 | 状态 | 说明 |
|------|------|------|
| pkg/models/asset.go | ✅ 已完成 | Asset 全局资产表 struct + Request/Response（Code+Market UNIQUE） |
| pkg/models/user_asset.go | ✅ 已完成 | UserAsset 用户持仓表 struct + Request/Response（含软删除） |
| pkg/models/investment_asset.go | ❌ 已删除 | 旧 InvestmentAsset model（已清理） |
| pkg/models/investment_transaction.go | ✅ 已完成 | InvestmentTransaction struct + Request/Response |
| pkg/models/market_data.go | ✅ 已完成 | MarketData struct + Request/Response |
| pkg/models/investment_consts.go | ✅ 已存在 | InvestmentAssetType/TransactionType/Market 常量定义 |
| cmd/database.go | ✅ 已修改 | 注册 4 张新表（InvestmentTransaction/MarketData/Asset/UserAsset） |
| pkg/errs/investment_asset.go | ✅ 已完成 | 投资资产相关错误定义 |
| pkg/errs/investment_transaction.go | ✅ 已完成 | 投资交易相关错误定义 |
| pkg/errs/market_data.go | ✅ 已完成 | 行情数据相关错误定义 |
| pkg/errs/setting.go | ✅ 已修改 | 新增 ErrInvalidMarketDataSource |
| pkg/uuid/uuid_type.go | ✅ 已修改 | 新增 UUID_TYPE_ASSET (13) 和 UUID_TYPE_USER_ASSET (14) |
| pkg/services/asset.go | ✅ 已完成 | AssetService 全局资产 CRUD |
| pkg/services/user_asset.go | ✅ 已完成 | UserAssetService 用户持仓管理（软删除） |
| pkg/services/investment_asset.go | ❌ 已删除 | 旧 InvestmentAssetService（已清理） |
| pkg/services/investment_transaction.go | ✅ 已完成 | InvestmentTransactionService CRUD + Balance 自动维护 |
| pkg/services/market_data.go | ✅ 已完成 | MarketDataService CRUD + FetchAllActiveAssetsMarketData + InitAssetMarketData |
| pkg/api/investment.go | ✅ 已完成 | InvestmentApi Handler（全局资产/用户持仓/交易/行情，旧 Asset Handler 已清理） |
| cmd/webserver.go | ✅ 已修改 | 注册 /investment/ 路由组（20 个端点，旧 /assets/ 路由已清理） |
| pkg/marketdata/market_data_provider.go | ✅ 已完成 | MarketDataProvider 接口定义 |
| pkg/marketdata/akshare_market_data_provider.go | ✅ 已完成 | akshare 数据源实现（备选） |
| pkg/marketdata/eastmoney_market_data_provider.go | ✅ 已完成 | 东方财富数据源实现（主力） |
| pkg/marketdata/market_data_provider_container.go | ✅ 已完成 | 容器 + 配置切换 + fallback |
| pkg/settings/setting.go | ✅ 已修改 | 新增 MarketDataSource 配置项 |
| pkg/cron/cron_jobs.go | ✅ 已修改 | 新增 FetchMarketDataJob（每日 18:00） |
| pkg/cron/cron_container.go | ✅ 已修改 | 注册 FetchMarketDataJob |
| cmd/initializer.go | ✅ 已修改 | 初始化 MarketDataSource |
| .air.toml | ✅ 已完成 | 热重载配置（macOS/Linux） |
| .air.windows.toml | ✅ 已完成 | Windows 热重载配置（使用 entrypoint + dev.bat） |
| dev.bat | ✅ 已完成 | Windows 热重载启动脚本 |
| go.mod | ✅ 已修改 | 引入 akshare 依赖 |
| scripts/migrate_investment_asset.sql | ✅ 已完成 | 数据迁移脚本（InvestmentAsset → Asset + UserAsset） |
| src/models/investment.ts | ✅ 已完成 | 前端 TS 类型定义（4 个 Model 类 + 20+ Request/Response 接口） |
| src/lib/services.ts | ✅ 已修改 | 前端 API 方法（新增 18 个投资模块方法） |
| src/stores/investment.ts | ✅ 已完成 | 前端 Pinia Store（持仓/交易/行情 CRUD） |
| src/stores/index.ts | ✅ 已修改 | 注册 investmentStore + resetAllStates |
| pkg/models/investment_analysis.go | ✅ 已完成 | 持仓聚合/概览 Response 类型（InvestmentHoldingInfo/InvestmentOverviewResponse/InvestmentAllocationItem） |
| pkg/services/investment_analysis.go | ✅ 已完成 | InvestmentAnalysisService（GetHoldings 持仓计算 + GetOverview 概览聚合） |
| pkg/api/investment.go | ✅ 已增强 | 新增 HoldingsHandler/OverviewHandler + TransactionListHandler 嵌入资产/账户名称 |
| cmd/webserver.go | ✅ 已增强 | 新增 2 个分析端点（analysis/holdings.json + analysis/overview.json，共 22 个端点） |
| src/models/investment.ts | ✅ 已增强 | 新增 InvestmentHolding/InvestmentOverview Model 类 + 持仓/概览 Response 接口 |
| src/lib/services.ts | ✅ 已增强 | 新增 getInvestmentHoldings/getInvestmentOverview API 方法（共 20 个投资方法） |
| src/stores/investment.ts | ✅ 已增强 | 新增 holdings/overview 状态 + loadHoldings/loadOverview 方法 |
| pkg/settings/setting.go | ✅ 已增强 | 新增 InvestmentAdminUid 配置项 |
| pkg/api/investment.go | ✅ 已增强 | 新增 AdminCheckHandler |
| cmd/webserver.go | ✅ 已增强 | 新增 admin/check.json 路由（共 23 个端点） |
| src/models/investment.ts | ✅ 已增强 | 新增 AdminCheckResponse 接口 |
| src/lib/services.ts | ✅ 已增强 | 新增 checkInvestmentAdmin API 方法（共 21 个投资方法） |
| src/views/desktop/investment/AssetsPage.vue | ✅ 已重构 | 搜索下拉、分类/市场 chips 筛选、买卖弹窗、管理按钮 |
| src/views/desktop/investment/AssetsPage.vue | ✅ 已增强 | Tab+表格合并单卡片、行业列（▾筛选icon）、分页兜底、mdiMagnify SVG 图标 |
| src/views/desktop/MainLayout.vue | ✅ 已增强 | watch currentRoutePath 自动同步 isInvestmentMode（/investment/* → 理财模式标题） |
| pkg/models/asset.go | ✅ 已增强 | 新增 AssetDeleteRequest、AssetListResponse、AssetListRequest 加 keyword/page/pageSize |
| pkg/services/asset.go | ✅ 已增强 | 新增 DeleteAsset、GetAllAssetsCount、GetAllAssets 加 keyword 参数 |
| pkg/api/investment.go | ✅ 已增强 | 新增 GlobalAssetListHandler/ModifyHandler/DeleteHandler（共 26 个端点） |
| cmd/webserver.go | ✅ 已增强 | 注册 list/modify/delete 路由 |
| src/models/investment.ts | ✅ 已增强 | 新增 AssetModifyRequest/DeleteRequest/ListRequest/ListResponse |
| src/lib/services.ts | ✅ 已增强 | 新增 listGlobalAssets/modifyGlobalAsset/deleteGlobalAsset（共 24 个投资 API） |
| src/locales/en.json + zh_Hans.json | ✅ 已增强 | 新增 Industry 相关 11 key + asset.AssetCode/NoAssetsFound/ConfirmDeleteAsset |

- 已有的前端代码：投资 Overview 页面骨架（已有，非本次新增）
- 已有的路由切换：点击 logo 切换记账/理财模式

---

## 三、核心设计决策（已确定，不要改）

### 3.1 全局资产表 + 用户持仓表

```
Asset（全局资产表）
├── 所有基金/股票信息只存储一份
├── 按 category 分类：equity/fixed_income/commodity/digital
├── 行业分类 industry（用于持仓分析）
└── 标签 tags（JSON 数组，用于搜索）

UserAsset（用户持仓表）
├── uid + asset_id 唯一
├── 只存储用户和资产的关联关系
└── is_active 控制是否活跃
```

### 3.2 三层资金流动模型

```
Layer 1：Transaction 表 → 只管出入金（银行卡↔投资池），不记录具体买了什么
Layer 2：InvestmentTransaction 表 → 只管投资交易（买/卖/分红/转换），不碰 Account.Balance 以外的记账逻辑
Layer 3：MarketData 表 → 每日行情，计算浮动盈亏
```

### 3.3 投资池 = 复用 Account 表

- Category = 7 (INVESTION)，已有枚举值
- Type = MultiSubAccounts，策略池作为子账户
- Balance = 池内可用现金，InvestmentTransaction CRUD 时自动更新
- 前端账户列表隐藏（Category=7 过滤），转账 API 不受影响

### 3.4 投资交易不写入 Transaction 表

- 独立 InvestmentTransaction 表，字段结构完全不同
- 资金层面只关心：银行卡→投资池（入金），投资池→银行卡（出金）

### 3.5 行情数据源

- **主力**：东方财富免费 HTTP API（实时估值 + 历史净值）
- **备选**：BlakeLiAFK/akshare Go 库
- 数据源可切换（策略模式），改配置 `market_data_source` 即可
- 更新策略：Cron 定时（每日 18:00）+ 手动刷新 API
- **QDII 基金**：东方财富无实时估值，自动 fallback 到 akshare

### 3.6 资产类别分类

| 类别 | 代码 | 包含类型 |
|------|------|----------|
| 权益类 | equity | 股票、股票基金、混合基金、ETF |
| 固定收益类 | fixed_income | 债券、债券基金、货币基金 |
| 商品类 | commodity | 黄金、商品基金 |
| 数字资产类 | digital | 加密货币 |

### 3.7 精度约定

| 字段类型 | 精度 | 说明 |
|---------|------|------|
| 金额 Amount/Fee | ×10000 | 4 位小数，与现有 Transaction.Amount 一致 |
| 份额 Quantity | ×10000 | 4 位小数，基金净值标准精度 |
| 净值/单价 Price | ×10000 | 4 位小数，如 1.2345 |
| 时间 | int64 Unix 时间戳 | 与现有表一致 |

---

## 四、数据库设计（已实现）

### 新增五张表

#### Asset（全局资产表）

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| AssetId | int64 | PK | UUID |
| Code | string | INDEX(IDX_asset_code_market) | 资产代码，如 005827 |
| Market | InvestmentMarket | INDEX(IDX_asset_code_market) | 市场：1=中国, 2=香港, 3=美国 |
| Name | string | VARCHAR(64) | 资产名称 |
| Category | AssetCategory | INDEX | equity/fixed_income/commodity/digital |
| Currency | string | VARCHAR(3) | 计价货币：CNY/USD/HKD |
| Industry | string | INDEX | 行业分类：technology/healthcare/consumer/... |
| Tags | string | TEXT | 标签 JSON 数组，用于搜索 |
| ExtraInfo | string | TEXT | 扩展信息 JSON |
| CreatedUnixTime | int64 | | 创建时间 |
| UpdatedUnixTime | int64 | | 更新时间 |

#### UserAsset（用户持仓表）

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| Id | int64 | PK | UUID |
| Uid | int64 | INDEX(IDX_user_asset_uid_asset_id) NOT NULL | 用户 ID |
| AssetId | int64 | INDEX(IDX_user_asset_uid_asset_id) NOT NULL | 资产 ID |
| Deleted | bool | NOT NULL | 软删除 |
| IsActive | bool | NOT NULL | 是否活跃 |
| IsWatchlist | bool | NOT NULL | 是否自选关注: true=自选, false=持仓 |
| Comment | string | VARCHAR(255) NOT NULL | 备注 |
| CreatedUnixTime | int64 | | 创建时间 |
| UpdatedUnixTime | int64 | | 更新时间 |
| DeletedUnixTime | int64 | | 删除时间 |

#### InvestmentAsset（待废弃）

> 原投资资产表，数据已迁移到 Asset + UserAsset

#### InvestmentTransaction

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| TransactionId | int64 | PK | UUID |
| Uid | int64 | INDEX | 用户 ID |
| Deleted | bool | INDEX | 软删除 |
| AssetId | int64 | INDEX | FK → Asset |
| AccountId | int64 | INDEX | FK → Account（策略池子账户） |
| Type | string | VARCHAR(30) | buy/sell/dividend_cash/dividend_reinvest/split/conversion_out/conversion_in |
| TradeTime | int64 | | 下单时间 |
| ConfirmTime | int64 | | 确认时间（T+N） |
| Quantity | int64 | | 份额 ×10000 |
| Price | int64 | | 单价 ×10000 |
| Amount | int64 | | 金额 ×10000 |
| Fee | int64 | | 手续费 ×10000 |
| RelatedTransactionId | int64 | | 配对交易ID（conversion 互指） |
| TimezoneUtcOffset | int16 | | 时区偏移 |
| Comment | string | VARCHAR(255) | 备注 |
| CreatedUnixTime | int64 | | 创建时间 |
| UpdatedUnixTime | int64 | | 更新时间 |
| DeletedUnixTime | int64 | | 删除时间 |

#### MarketData

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| DataId | int64 | PK autoincr | 自增主键 |
| AssetId | int64 | UNIQUE(UQE_market_data_asset_id_date) | FK → Asset |
| Date | int64 | UNIQUE(UQE_market_data_asset_id_date) | 日期 Unix 0点 |
| Price | int64 | NOT NULL | 净值 ×10000 |
| Volume | int64 | | 成交量（可选） |
| CreatedUnixTime | int64 | | 创建时间 |
| UpdatedUnixTime | int64 | | 更新时间 |

---

## 五、开发计划（按顺序）

### 阶段 1：数据层（后端）+ 行情基础设施 — 预估 2 周

| # | 任务 | 状态 | 参考文件 |
|---|------|------|---------|
| 1.1 | 创建 Model struct（Asset/UserAsset/InvestmentTransaction/MarketData） | ✅ 已完成 | pkg/models/ |
| 1.2 | 在 cmd/database.go 注册新表 | ✅ 已完成 | cmd/database.go |
| 1.3 | 创建 Service 层 CRUD | ✅ 已完成 | pkg/services/ |
| 1.4 | 创建 API Handler | ✅ 已完成 | pkg/api/investment.go |
| 1.5 | 注册投资路由组（25+ 端点） | ✅ 已完成 | cmd/webserver.go |
| 1.6 | 行情数据 Provider（东方财富 + akshare） | ✅ 已完成 | pkg/marketdata/ |
| 1.7 | Cron 任务（每日 18:00 拉净值） | ✅ 已完成 | pkg/cron/cron_jobs.go |
| 1.8 | 数据迁移脚本 | ✅ 已完成 | scripts/migrate_investment_asset.sql |
| 1.9 | 单元测试 | ⬜ 未开始 | |

### 阶段 1 补充：数据库设计修复

| # | 任务 | 状态 | 说明 |
|---|------|------|------|
| 1.R.1 | Asset.Code+Market UNIQUE 约束 | ✅ 已完成 | 防止重复创建同一基金 |
| 1.R.2 | UserAsset 加软删除 | ✅ 已完成 | 与项目惯例一致 |
| 1.R.3 | InvestmentTransaction 加 TradeTime 索引 | ✅ 已完成 | 优化日期范围查询 |
| 1.R.4 | 清理 InvestmentAsset 旧体系 | ✅ 已完成 | 删除旧 model/service/API/路由 |
| 1.R.5 | 更新 ER_DIAGRAM.md | ✅ 已完成 | 反映 Asset+UserAsset 新设计 |

### 阶段 2：前端 Store + API 层

| # | 任务 | 状态 | 参考文件 |
|---|------|------|---------|
| 2.1 | TS 类型定义 | ✅ 已完成 | src/models/investment.ts |
| 2.2 | Pinia Store | ✅ 已完成 | src/stores/investment.ts |
| 2.3 | services.ts API 方法 | ✅ 已完成 | src/lib/services.ts |
| 2.4 | rootStore 集成 | ✅ 已完成 | src/stores/index.ts |

### 阶段 2.5：API 契约检查

| # | 任务 | 状态 |
|---|------|------|
| 2.5.1 | 逐页列出数据需求 | ✅ 已完成 |
| 2.5.2 | 对照后端 API 找差距 | ✅ 已完成 |
| 2.5.3 | 一次性补齐 | ✅ 已完成 |

### 阶段 3：前端界面 — 预估 2-3 周

| # | 任务 | 状态 |
|---|------|------|
| 3.1 | OverviewPage | ⬜ |
| 3.2 | AssetsPage | 🔧 进行中（搜索+持仓/自选+行业列+管理CRUD，细节待完善） |
| 3.3 | TransactionsPage | ⬜ |
| 3.4 | PortfolioPage | ⬜ |
| 3.5 | AnalysisPage | ⬜ |
| 3.6 | StrategyPage | ⬜ |

### 阶段 4：导入增强 & 高级数据 — 预估 1-2 周

| # | 任务 | 状态 |
|---|------|------|
| 4.1 | 导入时自动创建缺失账户/分类/标签 | ⬜ |
| 4.2 | 投资数据导入 CSV/JSON | ⬜ |
| 4.3 | 引入 BlakeLiAFK/akshare Go 库 | ✅ 已完成 |

### 阶段 5：打磨 & 测试 — 预估 1 周

| # | 任务 | 状态 |
|---|------|------|
| 5.1 | 国际化补全 | ⬜ |
| 5.2 | 动画优化 | ⬜ |
| 5.3 | 整体测试 | ⬜ |
| 5.4 | 性能优化 | ⬜ |

---

## 六、开发前技术参考

详见 INVESTMENT_MODULE_DEV_PLAN.md 第十节，核心参考文件：

| 知识点 | 参考文件 | 行号 |
|--------|---------|------|
| Model struct + XORM | pkg/models/transaction_category.go | 17-32 |
| 表注册 | cmd/database.go | 85-107 |
| Service CRUD | pkg/services/transaction_categories.go | 36-100 |
| API Handler | pkg/api/transaction_categories.go | 42-60 |
| 路由注册 | cmd/webserver.go | 379-425 |
| Provider 模式 | pkg/exchangerates/exchange_rates_data_provider.go | 10-12 |
| MarketData Provider | pkg/marketdata/market_data_provider.go | — |
| Cron Job | pkg/cron/cron_jobs.go | — |
| Pinia Store | src/stores/account.ts | 26-34 |
| API 调用 | src/lib/services.ts | 490-510 |
| TS 类型 | src/models/account.ts | — |

---

## 七、当前状态

- 无阻塞问题
- 下一步：阶段 3 继续（OverviewPage → TransactionsPage → PortfolioPage → AnalysisPage）
- AssetsPage 基本功能已实现，细节待完善
- 36 个投资模块专属 i18n key 待迁移到 `asset.` 嵌套对象（用户说后面再改）
- 构建验证方式：`.\build.bat backend --no-lint --no-test`（Windows）/ `bash build.sh backend --no-lint --no-test`（macOS/Linux）（不要用 `go build ./...`）
- 前端验证：`npm run lint`（vue-tsc + eslint）
- 热重载：`air`（macOS/Linux 配置文件 `.air.toml`）/ `air -c .air.windows.toml`（Windows）

---

## 八、AssetsPage 待办事项（2026-06-02 新增）

### 设计决策

| 项目 | 决策 |
|------|------|
| 搜索框位置 | 独立在顶部 |
| [管理] 按钮 | 指定用户可见（InvestmentAdminUid 配置） |
| Tab 状态 | 持仓(is_watchlist=false) / 自选(is_watchlist=true) |
| 交易操作 | 弹出表单（类似记账模式） |
| 全局资产管理 | 放在数据管理菜单，低频操作 |
| Asset 表初始化 | 首次部署全量拉取，每季度增量更新 |

### 待完成任务

| 优先级 | 任务 | 状态 | 说明 |
|--------|------|------|------|
| P0 | 重构 AssetsPage 布局 | ✅ 已完成 | 搜索框独立、Tab 切换、管理按钮 |
| P0 | 实现搜索功能 | ✅ 已完成 | 搜索本地 Asset 表 |
| P0 | 实现买入/卖出弹窗 | ✅ 已完成 | 交易表单 |
| P1 | 添加配置项 InvestmentAdminUid | ✅ 已完成 | 控制管理按钮可见性 |
| P1 | 全局资产管理页面 | ✅ 已完成 | 管理员弹窗内 CRUD（搜索+行业筛选+新增/编辑/删除） |
| P2 | Asset 表初始化脚本 | ⬜ 待做 | 全量拉取基金/股票数据 |
| P2 | Cron 定时同步 | ⬜ 待做 | 每季度增量更新 |

### 页面功能定义

#### AssetsPage（资产管理）

| Tab | 数据来源 | 操作 |
|-----|----------|------|
| 搜索结果 | 本地 Asset 表 | 添加自选、买入 |
| 自选列表 | UserAsset (is_watchlist=true) | 买入、移除 |
| 持仓列表 | UserAsset (is_watchlist=false) | 卖出、查看详情 |

#### TransactionsPage（交易记录）

| 功能 | 说明 |
|------|------|
| 交易列表 | 所有投资交易历史 |
| 添加交易 | 支持所有类型（买入/卖出/分红/转换等） |
| 编辑/删除 | 修改或删除交易 |

### UserAsset 字段状态

| is_active | is_watchlist | 状态 |
|-----------|--------------|------|
| true | false | 持仓 |
| true | true | 自选 |
| false | * | 不活跃 |

---

## 九、Git 提交记录

```
045d3d40 refactor: 数据库设计修复 + 清理旧 InvestmentAsset 体系
3f1c226b fix: 修复 air 在 Windows 上异常退出的问题
067ee587 feat: 完善行情数据 Provider（实时估值 + 历史初始化）
0067d4f4 fix: 修复 market_data 复合索引重复创建问题 + 添加 air 热重载配置
73d3d316 docs: add AI session context handoff document
9d4bcf77 docs: add investment module design notes with data flow verification
05678984 docs: add ER diagrams for existing and new investment module tables
4051a995 docs: optimize investment module development plan
```

---

## 九、会话历史摘要

### 会话 1（2026-05-08 ~ 2026-05-09）

讨论内容：
1. 移除移动端网页开发计划（后期独立 App）
2. 基金数据方案：选定 A+B 组合（东方财富 HTTP + BlakeLiAFK/akshare Go 库）
3. 行情数据采用策略模式（仿 exchange_rates Provider）
4. 开发顺序调整：行情基础设施提前到阶段 1
5. 新增阶段 2.5 API 契约检查
6. 生成现有数据库 ER 图 + 新增投资模块 ER 图
7. 技术学习清单（Go 后端 6 项 + Vue 前端 5 项）
8. 数据库设计深度讨论：
   - 投资池 = 复用 Account 表 Category=7
   - 三层资金流动模型（Transaction / InvestmentTransaction / MarketData）
   - 7 种交易类型定义
   - Balance 自洽性验证
   - 成本计算（加权平均法）
   - 三种收益率（TWR/MWR/累计）
9. 确认投资模块完全独立于记账系统，出入金走现有转账

### 会话 2（2026-05-14）

完成内容：
1. 阶段 1.1：创建 3 个 Model struct
   - `pkg/models/investment_asset.go`
   - `pkg/models/investment_transaction.go`
   - `pkg/models/market_data.go`
2. 阶段 1.2：在 `cmd/database.go` 注册 3 张表（182-206 行）
3. 构建验证：`.\build.bat backend --no-lint --no-test` 通过

### 会话 3（2026-05-22）

完成内容：
1. 阶段 1.3：创建 Service 层 CRUD
   - `pkg/services/investment_asset.go`（InvestmentAssetService）
   - `pkg/services/investment_transaction.go`（InvestmentTransactionService + Balance 自动维护）
   - `pkg/services/market_data.go`（MarketDataService）
2. 阶段 1.4：创建 API Handler + 注册路由
   - `pkg/api/investment.go`（InvestmentApi，15 个 Handler）
   - `cmd/webserver.go`（注册 /investment/ 路由组）
3. 新增错误定义
   - `pkg/errs/investment_asset.go`
   - `pkg/errs/investment_transaction.go`
   - `pkg/errs/market_data.go`
4. 新增 UUID 类型
   - `pkg/uuid/uuid_type.go`（UUID_TYPE_INVESTMENT_ASSET=11, UUID_TYPE_INVESTMENT_TRANS=12）
5. 构建验证：`.\build.bat backend --no-lint --no-test` 通过

### 会话 4（2026-05-29）

完成内容：
1. 阶段 1.6：行情数据 Provider
   - `pkg/marketdata/market_data_provider.go`（接口定义）
   - `pkg/marketdata/akshare_market_data_provider.go`（akshare 数据源）
   - `pkg/marketdata/eastmoney_market_data_provider.go`（东方财富数据源）
   - `pkg/marketdata/market_data_provider_container.go`（容器 + 配置切换）
2. 配置项添加
   - `pkg/settings/setting.go`（新增 MarketDataSource 配置）
   - `pkg/errs/setting.go`（新增 ErrInvalidMarketDataSource）
3. Cron 任务
   - `pkg/cron/cron_jobs.go`（新增 FetchMarketDataJob）
   - `pkg/cron/cron_container.go`（注册任务）
4. 手动刷新 API
   - `pkg/api/investment.go`（新增 MarketDataRefreshHandler）
   - `cmd/webserver.go`（注册 /investment/market_data/refresh.json）
5. 依赖引入
   - `go.mod`（引入 BlakeLiAFK/akshare）
6. 开发工具
   - `.air.toml`（热重载配置）
7. 构建验证：`bash build.sh backend --no-lint --no-test` 通过
8. 接口测试：通过 Apifox 测试资产/交易/行情 CRUD 接口
9. 数据导入：使用 `scripts/convert_investment_data.py` 导入测试数据

### 会话 5（2026-05-30）

完成内容：
1. 全局资产表重构
   - `pkg/models/asset.go`（Asset 全局资产表）
   - `pkg/models/user_asset.go`（UserAsset 用户持仓表）
   - `pkg/services/asset.go`（AssetService）
   - `pkg/services/user_asset.go`（UserAssetService）
2. 新增 UUID 类型
   - `pkg/uuid/uuid_type.go`（UUID_TYPE_ASSET=13, UUID_TYPE_USER_ASSET=14）
3. 新增 API 端点
   - `/investment/global_assets/search.json`（搜索全局资产）
   - `/investment/global_assets/get.json`（获取资产详情）
   - `/investment/global_assets/add.json`（创建资产）
   - `/investment/user_assets/list.json`（用户持仓列表）
   - `/investment/user_assets/add.json`（添加持仓）
   - `/investment/user_assets/remove.json`（移除持仓）
   - `/investment/market_data/init.json`（初始化历史净值）
   - `/investment/market_data/estimate.json`（实时估值）
4. 行情 Provider 完善
   - akshare 日期解析修复（int64 时间戳）
   - 东方财富空数据处理（QDII 基金）
   - fallback 机制：东方财富失败自动切换 akshare
5. 数据迁移脚本
   - `scripts/migrate_investment_asset.sql`（InvestmentAsset → Asset + UserAsset）
6. 构建验证：`bash build.sh backend --no-lint --no-test` 通过
  8. 接口测试：通过 Apifox 测试资产/交易/行情 CRUD 接口
  9. 数据导入：使用 `scripts/convert_investment_data.py` 导入测试数据

### 会话 6（2026-05-31）

完成内容：
1. 修复 air 在 Windows 上异常退出
    - `.air.windows.toml`：使用 `entrypoint` + `dev.bat` 包装
    - 根因：air v1.65.3 通过 PowerShell 执行命令，PS 5.1 不支持 `&&`；`full_bin` 字段会自动拼接 `.exe` 后缀
2. 数据库设计审查 + 修复
    - Asset.Code+Market: INDEX → UNIQUE
    - UserAsset: 加 Deleted/DeletedUnixTime/UpdatedUnixTime/Comment
    - InvestmentTransaction: 加 (Uid,Deleted,TradeTime) 复合索引
3. 清理 InvestmentAsset 旧体系
    - 删除 `pkg/models/investment_asset.go`、`pkg/services/investment_asset.go`
    - 删除 5 个旧 Asset API Handler + 5 条旧路由
    - `market_data.go` FetchAllActiveAssetsMarketData 改查 UserAsset+Asset
    - `cmd/database.go` 移除 InvestmentAsset 表注册
4. 更新 ER_DIAGRAM.md 反映当前设计
5. 阶段 2：前端 Store + API 层
    - `src/models/investment.ts`：4 个 Model 类 + 20+ Request/Response 接口 + 3 个枚举
    - `src/lib/services.ts`：新增 18 个 API 方法（全局资产/用户持仓/交易/行情）
    - `src/stores/investment.ts`：Pinia Store（持仓/交易/行情 CRUD + 脏标记 + reset）
    - `src/stores/index.ts`：注册 investmentStore + resetAllStates
6. 构建验证：`.\build.bat backend --no-lint --no-test` + `npm run build` 均通过

### 会话 7（2026-05-31）

完成内容：
1. 阶段 2.5 API 契约检查（2.5.1 + 2.5.2 + 2.5.3 全部完成）
2. 发现的缺口：
   - ❌ 无持仓聚合端点（每资产持有份额、成本、市值、收益）
   - ❌ 无概览聚合端点（总投入/总市值/收益率/配置比例）
   - ❌ 交易响应缺少资产名称/账户名称
   - ❌ 前端缺少持仓/概览类型和方法
3. 补齐内容（2.5.3）：
   - `pkg/models/investment_analysis.go`：InvestmentHoldingInfo / InvestmentOverviewResponse / InvestmentAllocationItem
   - `pkg/services/investment_analysis.go`：InvestmentAnalysisService
     - GetHoldings()：遍历交易记录计算每资产持仓（加权平均成本法），获取最新行情，计算市值/浮动盈亏/收益率
     - GetOverview()：汇总所有持仓，按 category 分组计算配置比例
   - `pkg/api/investment.go`：新增 HoldingsHandler + OverviewHandler
   - `pkg/models/investment_transaction.go`：Response 新增 assetName/assetCode/accountName 字段 + ToInvestmentTransactionInfoResponseWithInfo 方法
   - `pkg/api/investment.go`：TransactionListHandler 增强 — 批量加载资产/账户名称嵌入响应
   - `cmd/webserver.go`：注册 2 个新路由（analysis/holdings.json + analysis/overview.json，共 22 个端点）
   - `src/models/investment.ts`：新增 InvestmentHolding/InvestmentOverview Model 类 + 3 个新接口
   - `src/lib/services.ts`：新增 getInvestmentHoldings/getInvestmentOverview 方法（共 20 个投资 API）
   - `src/stores/investment.ts`：新增 holdings/overview 状态 + loadHoldings/loadOverview 方法
4. 构建验证：`.\build.bat backend --no-lint --no-test` + `npm run build` 均通过

### 会话 8（2026-06-02）

完成内容：
1. AssetsPage 重构（基于用户在公司设计的 UI 方案）
   - 搜索框独立在顶部，输入时弹出下拉列表（类似搜索引擎）
   - 分类筛选用 chips（全部/权益类/固收类/商品类/数字资产）
   - 市场筛选用 chips（全部市场/中国/香港/美国）
   - 整个页面只有一个卡片（tabs + 表格）
   - 搜索结果显示状态标签（已持有/自选中）
   - 买入/卖出弹窗（选择资金池、交易类型、金额等）
   - 管理按钮（仅 InvestmentAdminUid 可见）
2. 后端配置
   - `pkg/settings/setting.go`：新增 InvestmentAdminUid 配置项
   - `pkg/api/investment.go`：新增 AdminCheckHandler
   - `cmd/webserver.go`：注册 /investment/admin/check.json 路由
3. 前端实现
   - `src/views/desktop/investment/AssetsPage.vue`：完整重构
   - `src/lib/services.ts`：新增 checkInvestmentAdmin API 方法
   - `src/models/investment.ts`：新增 AdminCheckResponse 接口
4. i18n 补全
   - 新增 Search Results、All Markets、Search assets、Add to Watchlist、Buy 等 key
5. 构建验证：`.\build.bat backend --no-lint --no-test` + `npm run build` 均通过

下一个 AI 应该做什么：
- 读取此文档了解完整上下文
- 读取 docs/ 下三个文档了解详细设计
- 继续阶段 3：其他页面开发（OverviewPage → TransactionsPage → PortfolioPage → AnalysisPage）
- 参考文件：`src/stores/investment.ts`、`src/lib/services.ts`、`src/models/investment.ts`
- 构建验证方式：`.\build.bat backend --no-lint --no-test`（Windows）/ `bash build.sh backend --no-lint --no-test`（macOS/Linux）
- 前端构建验证：`npm run lint`

### 会话 9（2026-06-03 ~ 2026-06-04）

完成内容：
1. AssetsPage 修复合并冲突后持续完善
   - 修复合并冲突导致的全部 TS 编译错误（重复声明、缺失函数/变量）
   - 补全缺失函数：checkAdmin、showManageButton、submitTransaction 等
   - 对照设计稿实现：搜索栏+结果浮层、持仓/自选双 Tab 表格、汇总栏、买卖/详情对话框
2. UI 调整（根据用户反馈逐项优化）
   - Tab 栏+表格合并为单个 v-card（参考记账模式）
   - 分页器改为记账模式样式（只保留 PaginationButtons，去掉 v-select）
   - 表头去掉排序箭头（sortable: false）
   - 市场列后加行业列（title 带 ▾ 筛选 icon）
   - 搜索图标改为 mdiMagnify SVG（@mdi/js 导入，density=compact）
   - 分页 totalPageCount 加 Math.max(1, ...) 兜底
   - 代码列去掉头像
   - 导航栏标题自动同步路由（MainLayout.vue watch currentRoutePath）
   - Tab 数字用 v-chip size="small" variant="tonal"
3. i18n 补全
   - Industry 及 10 个行业分类中英文 key
   - asset.SearchAssetsByNameOrCode、asset.AssetName、asset.TotalMarketValue 等
4. 技术决策
   - isInvestmentMode 不持久化到 localStorage，完全由 URL 路由决定
   - i18n 投资专属 key 用 `asset.` 嵌套对象包裹
   - 36 个投资模块专属 key 用户说后面再迁移到 asset. 下

### 会话 10（2026-06-04）

完成内容：
1. 全局资产管理 CRUD（前后端一起实现）
   - 后端：AssetDeleteRequest model、DeleteAsset/GetAllAssetsCount service 方法
   - 后端：GlobalAssetListHandler/ModifyHandler/DeleteHandler（asset.go + investment.go）
   - 后端：注册 list.json/modify.json/delete.json 路由（共 26 个端点）
   - 前端：AssetModifyRequest/DeleteRequest/ListRequest/ListResponse 类型
   - 前端：listGlobalAssets/modifyGlobalAsset/deleteGlobalAsset API（共 24 个投资 API）
   - 前端：管理弹窗完整 CRUD（搜索+行业筛选表格、新增/编辑表单、删除确认）
   - i18n：新增 AssetCode、NoAssetsFound、ConfirmDeleteAsset
2. 后端 API 端点统计（共 26 个）
   - global_assets: search/get/list/add/modify/delete（6 个）
   - user_assets: list/add/remove（3 个）
   - transactions: list/get/add/modify/delete（5 个）
   - market_data: latest/list/add/modify/refresh/init/estimate（7 个）
   - analysis: holdings/overview（2 个）
   - admin: check（1 个）
   - overview: 策略页（2 个，非本次）
3. 验证：后端 `go build` + `go test` 通过，前端 `npm run lint` 零错误
