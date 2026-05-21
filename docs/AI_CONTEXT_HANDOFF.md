# AI 会话上下文交接文档

> 每次会话结束后更新此文件，确保下一个 AI 会话能无缝接续。
> 最后更新：2026-05-09

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

### 2.2 代码

| 文件 | 状态 | 说明 |
|------|------|------|
| pkg/models/investment_asset.go | ✅ 已完成 | InvestmentAsset struct + Request/Response |
| pkg/models/investment_transaction.go | ✅ 已完成 | InvestmentTransaction struct + Request/Response |
| pkg/models/market_data.go | ✅ 已完成 | MarketData struct + Request/Response |
| cmd/database.go | ✅ 已修改 | 注册 3 张新表（182-206 行） |

- 已有的前端代码：投资 Overview 页面骨架（已有，非本次新增）
- 已有的路由切换：点击 logo 切换记账/理财模式

---

## 三、核心设计决策（已确定，不要改）

### 3.1 三层资金流动模型

```
Layer 1：Transaction 表 → 只管出入金（银行卡↔投资池），不记录具体买了什么
Layer 2：InvestmentTransaction 表 → 只管投资交易（买/卖/分红/转换），不碰 Account.Balance 以外的记账逻辑
Layer 3：MarketData 表 → 每日行情，计算浮动盈亏
```

### 3.2 投资池 = 复用 Account 表

- Category = 7 (INVESTION)，已有枚举值
- Type = MultiSubAccounts，策略池作为子账户
- Balance = 池内可用现金，InvestmentTransaction CRUD 时自动更新
- 前端账户列表隐藏（Category=7 过滤），转账 API 不受影响

### 3.3 投资交易不写入 Transaction 表

- 独立 InvestmentTransaction 表，字段结构完全不同
- 资金层面只关心：银行卡→投资池（入金），投资池→银行卡（出金）

### 3.4 行情数据源

- 方案 A：东方财富免费 HTTP API（主力，仿照现有 exchange_rates Provider 模式）
- 方案 B：BlakeLiAFK/akshare Go 库（高级数据，持仓/分红等）
- 数据源可切换（策略模式），加一个 Provider 实现改配置即可

### 3.5 精度约定

| 字段类型 | 精度 | 说明 |
|---------|------|------|
| 金额 Amount/Fee | ×10000 | 与现有 Transaction.Amount 一致 |
| 份额/单价 Quantity/Price | ×100000000 | 8 位小数精度 |
| 时间 | int64 Unix 时间戳 | 与现有表一致 |

---

## 四、数据库设计（待实现）

### 新增三张表

#### InvestmentAsset

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| AssetId | int64 | PK | |
| Uid | int64 | INDEX | |
| Deleted | bool | INDEX | 软删除 |
| Type | string | VARCHAR(20) NOT NULL | fund/stock/ETF/bond/crypto |
| Market | string | VARCHAR(10) NOT NULL | CN/HK/US |
| Code | string | VARCHAR(20) NOT NULL | 代码如 005827 |
| Name | string | VARCHAR(64) NOT NULL | |
| Currency | string | VARCHAR(3) NOT NULL | CNY/USD/HKD |
| IsActive | bool | NOT NULL | |
| ExtraInfo | text | BLOB | JSON: 行业/基金公司/经理/费率/持仓分布 |
| Comment | string | VARCHAR(255) NOT NULL | |
| CreatedUnixTime | int64 | | |
| UpdatedUnixTime | int64 | | |
| DeletedUnixTime | int64 | | |

#### InvestmentTransaction

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| TransactionId | int64 | PK | |
| Uid | int64 | INDEX | |
| Deleted | bool | INDEX | 软删除 |
| AssetId | int64 | INDEX | FK → InvestmentAsset |
| AccountId | int64 | NOT NULL | FK → Account（策略池子账户） |
| Type | string | VARCHAR(30) NOT NULL | buy/sell/dividend_cash/dividend_reinvest/split/conversion_out/conversion_in |
| TradeTime | int64 | NOT NULL | 下单时间 |
| ConfirmTime | int64 | | 确认时间（T+N） |
| Quantity | int64 | NOT NULL | 份额 ×100000000 |
| Price | int64 | NOT NULL | 单价 ×100000000 |
| Amount | int64 | NOT NULL | 金额 ×10000 |
| Fee | int64 | NOT NULL DEFAULT 0 | 手续费 ×10000 |
| RelatedTransactionId | int64 | | 配对交易ID（conversion 互指） |
| TimezoneUtcOffset | int16 | NOT NULL | |
| Comment | string | VARCHAR(255) NOT NULL | |
| CreatedUnixTime | int64 | | |
| UpdatedUnixTime | int64 | | |
| DeletedUnixTime | int64 | | |

#### MarketData

| 字段 | 类型 | XORM | 说明 |
|------|------|------|------|
| DataId | int64 | PK | |
| AssetId | int64 | INDEX UNIQUE(AssetId+Date) | FK → InvestmentAsset |
| Date | int64 | UNIQUE(AssetId+Date) | 日期 Unix 0点 |
| Price | int64 | NOT NULL | 净值 ×100000000 |
| Volume | int64 | | 成交量（可选） |
| CreatedUnixTime | int64 | | |
| UpdatedUnixTime | int64 | | |

---

## 五、开发计划（按顺序）

### 阶段 1：数据层（后端）+ 行情基础设施 — 预估 2 周

| # | 任务 | 状态 | 参考文件 |
|---|------|------|---------|
| 1.1 | 创建 3 个 Model struct | ✅ 已完成 | pkg/models/investment_asset.go, investment_transaction.go, market_data.go |
| 1.2 | 在 cmd/database.go 注册 3 张表 | ✅ 已完成 | cmd/database.go:182-206 |
| 1.3 | 创建 Service 层 CRUD | ⬜ 未开始 | pkg/services/transaction_categories.go |
| 1.4 | 创建 API Handler | ⬜ 未开始 | pkg/api/transaction_categories.go |
| 1.5 | 注册投资路由组 | ⬜ 未开始 | cmd/webserver.go:379 |
| 1.6 | 行情数据 Provider（仿 exchange_rates） | ⬜ 未开始 | pkg/exchangerates/ |
| 1.7 | Cron 任务（每日 18:00 拉净值） | ⬜ 未开始 | pkg/cron/cron_jobs.go |
| 1.8 | 单元测试 | ⬜ 未开始 | |

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
| 4.3 | 引入 BlakeLiAFK/akshare Go 库 | ⬜ |

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
| Cron Job | pkg/cron/cron_jobs.go | — |
| Pinia Store | src/stores/account.ts | 26-34 |
| API 调用 | src/lib/services.ts | 490-510 |
| TS 类型 | src/models/account.ts | — |

---

## 七、当前状态

- 无阻塞问题
- 下一步：阶段 1.3 创建 Service 层 CRUD
- 用户会在两台电脑间切换开发，此文档是 AI 会话的上下文桥梁
- 构建验证方式：`bash build.sh backend --no-lint --no-test`（不要用 `go build ./...`）

---

## 八、Git 提交记录

```
73d3d316 docs: add AI session context handoff document
9d4bcf77 docs: add investment module design notes with data flow verification
05678984 docs: add ER diagrams for existing and new investment module tables
4051a995 docs: optimize investment module development plan
```

四个提交已在 feature/add_invest_analysis 分支，已 push 到远端。

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
3. 构建验证：`bash build.sh backend --no-lint --no-test` 通过

下一个 AI 应该做什么：
- 读取此文档了解完整上下文
- 读取 docs/ 下三个文档了解详细设计
- 从阶段 1.3 开始：创建 Service 层 CRUD
- 参考文件：`pkg/services/transaction_categories.go`
- 构建验证方式：`bash build.sh backend --no-lint --no-test`
