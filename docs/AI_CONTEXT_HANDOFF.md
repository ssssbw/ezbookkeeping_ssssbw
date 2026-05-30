# AI 会话上下文交接文档

> 每次会话结束后更新此文件，确保下一个 AI 会话能无缝接续。
> 最后更新：2026-05-30

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
| pkg/models/asset.go | ✅ 已完成 | Asset 全局资产表 struct + Request/Response |
| pkg/models/user_asset.go | ✅ 已完成 | UserAsset 用户持仓表 struct + Request/Response |
| pkg/models/investment_asset.go | ✅ 已完成 | InvestmentAsset struct + Request/Response（待废弃） |
| pkg/models/investment_transaction.go | ✅ 已完成 | InvestmentTransaction struct + Request/Response |
| pkg/models/market_data.go | ✅ 已完成 | MarketData struct + Request/Response |
| pkg/models/investment_consts.go | ✅ 已存在 | InvestmentAssetType/TransactionType/Market 常量定义 |
| cmd/database.go | ✅ 已修改 | 注册 5 张新表 |
| pkg/errs/investment_asset.go | ✅ 已完成 | 投资资产相关错误定义 |
| pkg/errs/investment_transaction.go | ✅ 已完成 | 投资交易相关错误定义 |
| pkg/errs/market_data.go | ✅ 已完成 | 行情数据相关错误定义 |
| pkg/errs/setting.go | ✅ 已修改 | 新增 ErrInvalidMarketDataSource |
| pkg/uuid/uuid_type.go | ✅ 已修改 | 新增 UUID_TYPE_ASSET (13) 和 UUID_TYPE_USER_ASSET (14) |
| pkg/services/asset.go | ✅ 已完成 | AssetService 全局资产 CRUD |
| pkg/services/user_asset.go | ✅ 已完成 | UserAssetService 用户持仓管理 |
| pkg/services/investment_asset.go | ✅ 已完成 | InvestmentAssetService CRUD（待废弃） |
| pkg/services/investment_transaction.go | ✅ 已完成 | InvestmentTransactionService CRUD + Balance 自动维护 |
| pkg/services/market_data.go | ✅ 已完成 | MarketDataService CRUD + FetchAllActiveAssetsMarketData + InitAssetMarketData |
| pkg/api/investment.go | ✅ 已完成 | InvestmentApi Handler（全局资产/用户持仓/交易/行情） |
| cmd/webserver.go | ✅ 已修改 | 注册 /investment/ 路由组（25+ 个端点） |
| pkg/marketdata/market_data_provider.go | ✅ 已完成 | MarketDataProvider 接口定义 |
| pkg/marketdata/akshare_market_data_provider.go | ✅ 已完成 | akshare 数据源实现（备选） |
| pkg/marketdata/eastmoney_market_data_provider.go | ✅ 已完成 | 东方财富数据源实现（主力） |
| pkg/marketdata/market_data_provider_container.go | ✅ 已完成 | 容器 + 配置切换 + fallback |
| pkg/settings/setting.go | ✅ 已修改 | 新增 MarketDataSource 配置项 |
| pkg/cron/cron_jobs.go | ✅ 已修改 | 新增 FetchMarketDataJob（每日 18:00） |
| pkg/cron/cron_container.go | ✅ 已修改 | 注册 FetchMarketDataJob |
| cmd/initializer.go | ✅ 已修改 | 初始化 MarketDataSource |
| .air.toml | ✅ 已完成 | 热重载配置 |
| go.mod | ✅ 已修改 | 引入 akshare 依赖 |
| scripts/migrate_investment_asset.sql | ✅ 已完成 | 数据迁移脚本（InvestmentAsset → Asset + UserAsset） |

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
| Uid | int64 | INDEX(IDX_user_asset_uid_asset_id) | 用户 ID |
| AssetId | int64 | INDEX(IDX_user_asset_uid_asset_id) | 资产 ID |
| IsActive | bool | | 是否活跃 |
| AddedUnixTime | int64 | | 添加时间 |

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

### 阶段 2：前端 Store + API 层 — 预估 1 周

| # | 任务 | 状态 |
|---|------|------|
| 2.1 | TS 类型定义 | ⬜ |
| 2.2 | Pinia Store | ⬜ |
| 2.3 | services.ts API 方法 | ⬜ |
| 2.4 | rootStore 集成 | ⬜ |

### 阶段 2.5：API 契约检查

| # | 任务 | 状态 |
|---|------|------|
| 2.5.1 | 逐页列出数据需求 | ⬜ |
| 2.5.2 | 对照后端 API 找差距 | ⬜ |
| 2.5.3 | 一次性补齐 | ⬜ |

### 阶段 3：前端界面 — 预估 2-3 周

| # | 任务 | 状态 |
|---|------|------|
| 3.1 | OverviewPage | ⬜ |
| 3.2 | AssetsPage | ⬜ |
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
- 下一步：阶段 2 前端 Store + API 层
- 用户会在两台电脑间切换开发，此文档是 AI 会话的上下文桥梁
- 构建验证方式：`.\build.bat backend --no-lint --no-test`（Windows）/ `bash build.sh backend --no-lint --no-test`（macOS/Linux）（不要用 `go build ./...`）
- 热重载：`air`（配置文件：`.air.toml`）

---

## 八、Git 提交记录

```
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

下一个 AI 应该做什么：
- 读取此文档了解完整上下文
- 读取 docs/ 下三个文档了解详细设计
- 从阶段 2 开始：前端 Store + API 层
- 参考文件：`src/stores/account.ts`、`src/lib/services.ts`、`src/models/account.ts`
- 构建验证方式：`.\build.bat backend --no-lint --no-test`（Windows）/ `bash build.sh backend --no-lint --no-test`（macOS/Linux）
