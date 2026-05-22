# 投资模块设计笔记

> 基于 2026-05-09 讨论精炼
> 核心问题：资金流向、数据自洽、边界划分

---

## 一、核心结论：三层资金流动模型

```
Layer 1：记账系统（Transaction 表）
  ┌──────────┐    转账     ┌──────────────┐    转账     ┌──────────┐
  │  银行卡A  │ ────────> │  投资池-策略A  │ ────────> │  生活账户  │
  │  Balance  │           │  Balance      │           │  Balance  │
  └──────────┘           └──────┬───────┘           └──────────┘
                             │
Layer 2：投资交易（InvestmentTransaction 表） 只关心：投资池的钱买了什么
  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │ 基金005827│  │ 基金161725│  │ 股票AAPL │
  └──────────┘  └──────────┘  └──────────┘

Layer 3：行情数据（MarketData 表）每日净值 → 计算浮动盈亏
```

| 层 | 表 | 关心什么 | 不关心什么 |
|----|-----|---------|-----------|
| 记账 | Transaction | 银行卡↔投资池，转了多少 | 买了什么基金，赚了多少 |
| 投资 | InvestmentTransaction | 买了多少份额，什么价格 | 钱从哪张银行卡来 |
| 行情 | MarketData | 每天净值 | — |

### 为什么不把投资交易混入记账系统

1. 买入基金 ≠ 消费支出，是资产形态转换
2. 投资交易字段结构完全不同（份额、单价、手续费、确认时间、配对交易）
3. 查询模式不同（投资按资产聚合，记账按分类聚合）
4. 独立表 = 独立迭代，互不干扰

---

## 二、投资池 = 复用 Account 表

利用现有 `Account.Type=MultiSubAccounts, Category=7(INVESTMENT)` 实现策略隔离：

```
Account(总):     Name="投资账户",  Type=MultiSubAccounts, Category=7
Account(子1):    Name="定投策略",  ParentAccountId=总账户,  Balance=5000
Account(子2):    Name="高频策略",  ParentAccountId=总账户,  Balance=3000
Account(子3):    Name="海外配置",  ParentAccountId=总账户,  Balance=2000
```

- **出入金**：走现有 Transaction 转账，零改动
- **Balance 含义**：池内可用现金
- **Balance 维护**：创建/修改/删除 InvestmentTransaction 时自动更新
- **前端可见性**：账户列表隐藏（Category=7 过滤），但不影响转账 API

---

## 三、数据库表结构

### 3.1 InvestmentAsset（资产信息）

| 字段 | 类型 | 说明 |
|------|------|------|
| AssetId | int64 PK | |
| Uid | int64 INDEX | |
| Deleted | bool INDEX | 软删除 |
| Type | string | fund / stock / ETF / bond / crypto / ... |
| Market | string | CN / HK / US / ... |
| Code | string | 代码，如 005827 |
| Name | string | |
| Currency | string | CNY / USD / HKD |
| IsActive | bool | |
| ExtraInfo | text(JSON) | 行业标签、基金公司、经理、费率、持仓分布（自动拉取+手动补充，手动优先） |
| Comment | string | |
| CreatedUnixTime | int64 | |
| UpdatedUnixTime | int64 | |
| DeletedUnixTime | int64 | |

### 3.2 InvestmentTransaction（投资交易）

| 字段 | 类型 | 说明 |
|------|------|------|
| TransactionId | int64 PK | |
| Uid | int64 INDEX | |
| Deleted | bool INDEX | 软删除 |
| AssetId | int64 INDEX | FK → InvestmentAsset |
| AccountId | int64 | FK → Account（哪个策略池子账户） |
| Type | string | buy / sell / dividend_cash / dividend_reinvest / split / conversion_out / conversion_in |
| TradeTime | int64 | 下单时间 |
| ConfirmTime | int64 | 确认时间（T+N，份额到账日） |
| Quantity | int64 | 数量（×10000，4位小数精度） |
| Price | int64 | 单价（×10000，4位小数精度） |
| Amount | int64 | 金额（×10000，与现有 Transaction.Amount 精度一致） |
| Fee | int64 | 手续费（×10000） |
| RelatedTransactionId | int64 | 配对交易ID（conversion_out ↔ conversion_in 互指） |
| TimezoneUtcOffset | int16 | |
| Comment | string | |
| CreatedUnixTime | int64 | |
| UpdatedUnixTime | int64 | |
| DeletedUnixTime | int64 | |

### 3.3 MarketData（行情数据）

| 字段 | 类型 | 说明 |
|------|------|------|
| DataId | int64 PK | |
| AssetId | int64 INDEX | FK → InvestmentAsset |
| Date | int64 UNIQUE(AssetId+Date) | 日期（取 0 点） |
| Price | int64 | 当日净值/收盘价（×10000，4位小数精度） |
| Volume | int64 | 成交量（可选） |
| CreatedUnixTime | int64 | |
| UpdatedUnixTime | int64 | |

### 3.4 复用现有表：Account

| 字段 | 投资池用途 |
|------|-----------|
| Category | 7 (INVESTMENT) |
| Type | MultiSubAccounts（策略池作为子账户） |
| Balance | 池内可用现金，InvestmentTransaction 变更时自动更新 |
| Hidden | true（前端账户列表不显示） |

---

## 四、资金流自洽性验证

### 场景：入金 → 买入 → 分红 → 卖出 → 出金

```
1. 入金 10000
   Transaction: 银行卡 → 定投策略池, 10000
   → 定投策略池.Balance = 10000

2. 买入 5000 元基金（+7.50 手续费）
   InvestmentTransaction(buy): Amount=5000, Fee=7.50
   → 定投策略池.Balance = 10000 - 5000 - 7.50 = 4992.50

3. 现金分红 50 元
   InvestmentTransaction(dividend_cash): Amount=50
   → 定投策略池.Balance = 4992.50 + 50 = 5042.50

4. 卖出，回款 447.75（450 - 2.25 手续费）
   InvestmentTransaction(sell): Amount=450, Fee=2.25
   → 定投策略池.Balance = 5042.50 + 450 - 2.25 = 5490.25

5. 提取 2000 到生活账户
   Transaction: 定投策略池 → 生活账户, 2000
   → 定投策略池.Balance = 5490.25 - 2000 = 3490.25
```

**Balance 公式**：
```
投资池 Balance = 
  + SUM(入金转账)
  + SUM(卖出回款-手续费)
  + SUM(分红)
  - SUM(出金转账)
  - SUM(买入金额+手续费)
```

---

## 五、交易类型定义

| Type | 含义 | 份额变化 | 投资池 Balance | 是否配对 |
|------|------|---------|---------------|---------|
| buy | 买入 | +Quantity | -Amount - Fee | 否 |
| sell | 卖出 | -Quantity | +Amount - Fee | 否 |
| dividend_cash | 现金分红 | 不变 | +Amount | 否 |
| dividend_reinvest | 红利再投资 | +Quantity | 不变 | 否 |
| split | 份额拆分/合并 | ±Quantity | 不变 | 否 |
| conversion_out | 转出（转换出旧基金） | -Quantity | 不变 | 是，关联 conversion_in |
| conversion_in | 转入（转换入新基金） | +Quantity | 不变 | 是，关联 conversion_out |

---

## 六、成本与盈亏计算

### 持仓成本（加权平均法）

```
平均成本 = 累计投入 / 当前持有份额

卖出时：
  卖出成本 = 卖出份额 × 卖出前平均成本
  已实现盈亏 = (卖出回款 - 手续费) - 卖出成本
  卖出后平均成本不变（加权平均法下）
```

### 浮动盈亏

```
浮动盈亏 = 当前市值(MarketData.Price × 持仓份额) - 持仓成本(平均成本 × 持仓份额)
```

### 三种收益率（后期实现）

| 指标 | 含义 | 场景 |
|------|------|------|
| 累计收益率 | (当前市值-累计投入)/累计投入 | 简单看赚了多少 |
| 时间加权 TWR | 每期收益率连乘，消除入金影响 | 评价基金/策略本身表现 |
| 资金加权 MWR(IRR) | 考虑每笔现金流时间点的内部收益率 | 评价你自己的择时能力 |

---

## 七、设计决策记录

| 决策 | 选择 | 原因 |
|------|------|------|
| 投资交易是否写入 Transaction 表 | 不写，独立 InvestTransaction 表 | 字段结构完全不同，查询模式不同 |
| 投资池用什么表 | 复用 Account 表，Category=7 | 现有 MultiSubAccounts 模式完全匹配 |
| Balance 维护方式 | InvestmentTransaction CRUD 时自动更新 | 保证池内现金余额实时准确 |
| 前端是否显示投资池 | 账户列表隐藏，转账 API 不受影响 | 出金入金功能正常，视觉干净 |
| 行情数据获取 | 东方财富 HTTP + BlakeLiAFK/akshare | 纯 Go，零外部依赖 |
| 行业/持仓数据 | 自动拉取 + 手动补充（手动优先） | 覆盖海外资产，允许纠错 |
| 分红方式 | 两种都支持 | 真实场景都需要 |
| 手续费存储 | 单独字段，不混入 Amount | 盈亏计算需要区分净金额和手续费 |
| 确认时间 | 单独字段 ConfirmTime | T+N 确认制下，份额到账日前不计入持仓 |

---

## 八、与现有系统集成点

| 集成点 | 方式 | 说明 |
|--------|------|------|
| 用户认证 | 复用 JWT 中间件 | 零改动 |
| 账户体系 | Account.Category=7 | 现有枚举已有 INVESTMENT 值 |
| 货币/汇率 | 复用现有 exchange_rates 模块 | 海外资产需要多币种 |
| 时区 | 复用 TimezoneUtcOffset | 海外市场跨时区 |
| 行情数据拉取 | 仿照 exchange_rates Provider 模式 | 策略模式，数据源可切换 |
| Cron 定时任务 | 仿照现有 gocron 模式 | 每日 18:00 拉取基金净值 |
| 数据迁移 | 复用 code-driven SyncStructs | 在 cmd/database.go 注册新表 |