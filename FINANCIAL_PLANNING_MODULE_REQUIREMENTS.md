# 理财模块需求文档

## 1. 文档概述

### 1.1 文档目的

本文档描述了基于ezBookkeeping现有架构实现理财模块的详细需求，包括功能需求、非功能需求、数据模型设计、API设计、前端界面设计和集成方案。

### 1.2 术语定义

| 术语 | 解释 |
|------|------|
| 理财计划 | 用户制定的长期或短期财务目标和规划 |
| 预算 | 针对特定周期或类别设定的支出上限 |
| 投资组合 | 用户持有的各类投资资产的集合 |
| 净值 | 投资组合的总价值 |
| 收益率 | 投资组合的收益与成本的比率 |
| 风险评估 | 对投资组合或财务状况的风险分析 |

## 2. 功能需求

### 2.1 预算管理

#### 2.1.1 预算创建与编辑
- 支持按周期（日、周、月、季度、年）创建预算
- 支持按账户、分类、标签创建预算
- 支持设置预算金额和预警阈值
- 支持编辑和删除预算
- 支持预算模板功能

#### 2.1.2 预算执行跟踪
- 实时显示预算执行情况
- 支持预算执行趋势分析
- 支持预算超支预警（邮件、应用内通知）
- 支持预算执行报告生成

### 2.2 财务目标管理

#### 2.2.1 目标创建与编辑
- 支持创建短期和长期财务目标
- 支持设置目标金额、截止日期、当前进度
- 支持目标分类（如：购房、教育、退休、旅行等）
- 支持编辑和删除目标

#### 2.2.2 目标进度跟踪
- 实时显示目标达成进度
- 支持目标达成趋势分析
- 支持目标达成预测
- 支持目标调整建议

### 2.3 投资组合管理

#### 2.3.1 投资产品管理
- 支持手动添加投资产品（股票、基金、债券、房产等）
- 支持投资产品分类和标签
- 支持设置投资产品的当前价值、成本、持有数量
- 支持投资产品价格自动更新（通过API集成）

#### 2.3.2 投资组合分析
- 实时计算投资组合净值
- 支持投资组合收益率分析（日、周、月、年）
- 支持投资组合资产配置分析
- 支持投资组合风险评估

### 2.4 财务报表与分析

#### 2.4.1 标准财务报表
- 资产负债表生成
- 收支表生成
- 现金流量表生成

#### 2.4.2 财务分析
- 财务健康度评估
- 收支结构分析
- 储蓄率分析
- 债务收入比分析

### 2.5 智能建议

- 基于消费习惯的预算调整建议
- 基于财务目标的储蓄建议
- 基于投资组合的资产配置建议
- 基于收支情况的理财规划建议

## 3. 非功能需求

### 3.1 性能需求

- 页面加载时间 < 2秒
- 数据更新延迟 < 1秒
- 支持同时处理1000+用户请求

### 3.2 安全需求

- 所有敏感数据加密存储
- API请求鉴权和授权
- 支持数据备份和恢复

### 3.3 可用性需求

- 系统可用性 ≥ 99.9%
- 支持多语言（与现有系统一致）
- 支持多设备访问（桌面、移动）

### 3.4 可扩展性需求

- 支持插件式架构，便于扩展新的投资产品类型
- 支持第三方API集成（如：证券交易所API、银行API等）

## 4. 数据模型设计

### 4.1 现有数据模型扩展

#### 4.1.1 交易模型扩展

在现有交易模型中添加以下字段：
- `is_investment`: 布尔值，标识是否为投资交易
- `investment_type`: 字符串，投资类型（如：股票、基金、债券等）
- `investment_product_id`: 整数，关联投资产品ID

#### 4.1.2 分类模型扩展

在现有分类模型中添加以下字段：
- `is_investment_category`: 布尔值，标识是否为投资分类

### 4.2 新增数据模型

#### 4.2.1 预算模型 (budgets)

| 字段名 | 数据类型 | 约束 | 描述 |
|--------|----------|------|------|
| `id` | `int64` | `PRIMARY KEY` | 预算ID |
| `user_id` | `int64` | `NOT NULL` | 用户ID |
| `name` | `string` | `NOT NULL` | 预算名称 |
| `description` | `string` | | 预算描述 |
| `period` | `string` | `NOT NULL` | 预算周期（day, week, month, quarter, year） |
| `amount` | `decimal` | `NOT NULL` | 预算金额 |
| `warning_threshold` | `decimal` | | 预警阈值（百分比） |
| `account_ids` | `json` | | 关联账户ID列表 |
| `category_ids` | `json` | | 关联分类ID列表 |
| `tag_ids` | `json` | | 关联标签ID列表 |
| `start_date` | `datetime` | `NOT NULL` | 预算开始日期 |
| `end_date` | `datetime` | `NOT NULL` | 预算结束日期 |
| `is_active` | `bool` | `NOT NULL` | 是否激活 |
| `created_at` | `datetime` | `NOT NULL` | 创建时间 |
| `updated_at` | `datetime` | `NOT NULL` | 更新时间 |

#### 4.2.2 财务目标模型 (financial_goals)

| 字段名 | 数据类型 | 约束 | 描述 |
|--------|----------|------|------|
| `id` | `int64` | `PRIMARY KEY` | 目标ID |
| `user_id` | `int64` | `NOT NULL` | 用户ID |
| `name` | `string` | `NOT NULL` | 目标名称 |
| `description` | `string` | | 目标描述 |
| `category` | `string` | `NOT NULL` | 目标分类 |
| `target_amount` | `decimal` | `NOT NULL` | 目标金额 |
| `current_amount` | `decimal` | `NOT NULL` | 当前金额 |
| `deadline` | `datetime` | `NOT NULL` | 截止日期 |
| `priority` | `int` | `NOT NULL` | 优先级（1-5） |
| `is_completed` | `bool` | `NOT NULL` | 是否完成 |
| `created_at` | `datetime` | `NOT NULL` | 创建时间 |
| `updated_at` | `datetime` | `NOT NULL` | 更新时间 |

#### 4.2.3 投资产品模型 (investment_products)

| 字段名 | 数据类型 | 约束 | 描述 |
|--------|----------|------|------|
| `id` | `int64` | `PRIMARY KEY` | 产品ID |
| `user_id` | `int64` | `NOT NULL` | 用户ID |
| `name` | `string` | `NOT NULL` | 产品名称 |
| `type` | `string` | `NOT NULL` | 产品类型（stock, fund, bond, real_estate, etc.） |
| `code` | `string` | | 产品代码（如股票代码） |
| `current_value` | `decimal` | `NOT NULL` | 当前价值 |
| `cost_basis` | `decimal` | `NOT NULL` | 成本 |
| `quantity` | `decimal` | `NOT NULL` | 持有数量 |
| `current_price` | `decimal` | | 当前单价 |
| `currency` | `string` | `NOT NULL` | 货币类型 |
| `is_active` | `bool` | `NOT NULL` | 是否激活 |
| `created_at` | `datetime` | `NOT NULL` | 创建时间 |
| `updated_at` | `datetime` | `NOT NULL` | 更新时间 |
| `last_price_update` | `datetime` | | 最后价格更新时间 |

#### 4.2.4 投资交易模型 (investment_transactions)

| 字段名 | 数据类型 | 约束 | 描述 |
|--------|----------|------|------|
| `id` | `int64` | `PRIMARY KEY` | 交易ID |
| `user_id` | `int64` | `NOT NULL` | 用户ID |
| `product_id` | `int64` | `NOT NULL` | 投资产品ID |
| `transaction_id` | `int64` | `NOT NULL` | 关联主交易ID |
| `type` | `string` | `NOT NULL` | 交易类型（buy, sell, dividend, interest, etc.） |
| `quantity` | `decimal` | `NOT NULL` | 交易数量 |
| `price` | `decimal` | `NOT NULL` | 交易单价 |
| `fee` | `decimal` | | 交易费用 |
| `created_at` | `datetime` | `NOT NULL` | 创建时间 |

#### 4.2.5 投资组合模型 (investment_portfolios)

| 字段名 | 数据类型 | 约束 | 描述 |
|--------|----------|------|------|
| `id` | `int64` | `PRIMARY KEY` | 组合ID |
| `user_id` | `int64` | `NOT NULL` | 用户ID |
| `name` | `string` | `NOT NULL` | 组合名称 |
| `description` | `string` | | 组合描述 |
| `is_default` | `bool` | `NOT NULL` | 是否为默认组合 |
| `created_at` | `datetime` | `NOT NULL` | 创建时间 |
| `updated_at` | `datetime` | `NOT NULL` | 更新时间 |

#### 4.2.6 投资组合产品关联模型 (portfolio_products)

| 字段名 | 数据类型 | 约束 | 描述 |
|--------|----------|------|------|
| `id` | `int64` | `PRIMARY KEY` | 关联ID |
| `portfolio_id` | `int64` | `NOT NULL` | 投资组合ID |
| `product_id` | `int64` | `NOT NULL` | 投资产品ID |
| `created_at` | `datetime` | `NOT NULL` | 创建时间 |

## 5. API设计

### 5.1 预算管理API

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/budgets` | `GET` | 获取预算列表 | 登录用户 |
| `/api/v1/budgets/:id` | `GET` | 获取单个预算详情 | 登录用户 |
| `/api/v1/budgets` | `POST` | 创建预算 | 登录用户 |
| `/api/v1/budgets/:id` | `PUT` | 更新预算 | 登录用户 |
| `/api/v1/budgets/:id` | `DELETE` | 删除预算 | 登录用户 |
| `/api/v1/budgets/:id/statistics` | `GET` | 获取预算执行统计 | 登录用户 |

### 5.2 财务目标API

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/financial-goals` | `GET` | 获取财务目标列表 | 登录用户 |
| `/api/v1/financial-goals/:id` | `GET` | 获取单个财务目标详情 | 登录用户 |
| `/api/v1/financial-goals` | `POST` | 创建财务目标 | 登录用户 |
| `/api/v1/financial-goals/:id` | `PUT` | 更新财务目标 | 登录用户 |
| `/api/v1/financial-goals/:id` | `DELETE` | 删除财务目标 | 登录用户 |
| `/api/v1/financial-goals/:id/progress` | `GET` | 获取财务目标进度 | 登录用户 |

### 5.3 投资管理API

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/investment-products` | `GET` | 获取投资产品列表 | 登录用户 |
| `/api/v1/investment-products/:id` | `GET` | 获取单个投资产品详情 | 登录用户 |
| `/api/v1/investment-products` | `POST` | 创建投资产品 | 登录用户 |
| `/api/v1/investment-products/:id` | `PUT` | 更新投资产品 | 登录用户 |
| `/api/v1/investment-products/:id` | `DELETE` | 删除投资产品 | 登录用户 |
| `/api/v1/investment-products/:id/transactions` | `GET` | 获取投资产品交易记录 | 登录用户 |
| `/api/v1/investment-portfolios` | `GET` | 获取投资组合列表 | 登录用户 |
| `/api/v1/investment-portfolios/:id` | `GET` | 获取单个投资组合详情 | 登录用户 |
| `/api/v1/investment-portfolios` | `POST` | 创建投资组合 | 登录用户 |
| `/api/v1/investment-portfolios/:id` | `PUT` | 更新投资组合 | 登录用户 |
| `/api/v1/investment-portfolios/:id` | `DELETE` | 删除投资组合 | 登录用户 |
| `/api/v1/investment-portfolios/:id/statistics` | `GET` | 获取投资组合统计 | 登录用户 |
| `/api/v1/investment-transactions` | `POST` | 创建投资交易 | 登录用户 |
| `/api/v1/investment-transactions/:id` | `DELETE` | 删除投资交易 | 登录用户 |

### 5.4 财务分析API

| 路径 | 方法 | 功能 | 权限 |
|------|------|------|------|
| `/api/v1/financial-analysis/balance-sheet` | `GET` | 获取资产负债表 | 登录用户 |
| `/api/v1/financial-analysis/income-statement` | `GET` | 获取收支表 | 登录用户 |
| `/api/v1/financial-analysis/cash-flow` | `GET` | 获取现金流量表 | 登录用户 |
| `/api/v1/financial-analysis/net-worth` | `GET` | 获取净值趋势 | 登录用户 |
| `/api/v1/financial-analysis/financial-health` | `GET` | 获取财务健康度评估 | 登录用户 |
| `/api/v1/financial-analysis/recommendations` | `GET` | 获取理财建议 | 登录用户 |

## 6. 前端界面设计

### 6.1 桌面端界面

#### 6.1.1 理财模块导航

- 在现有导航菜单中添加"理财"菜单项
- 子菜单包括：预算管理、财务目标、投资组合、财务分析

#### 6.1.2 预算管理界面

- 预算列表视图，显示所有预算及其执行情况
- 预算创建/编辑表单
- 预算执行图表（柱状图、折线图）
- 预算预警设置

#### 6.1.3 财务目标界面

- 财务目标列表视图，显示所有目标及其进度
- 目标创建/编辑表单
- 目标进度可视化（进度条、仪表盘）
- 目标达成预测图表

#### 6.1.4 投资组合界面

- 投资组合列表视图
- 投资产品列表及详情
- 投资组合净值趋势图表
- 资产配置饼图
- 收益率分析图表

#### 6.1.5 财务分析界面

- 财务报表视图（资产负债表、收支表、现金流量表）
- 财务健康度评分
- 智能建议卡片
- 多种财务指标图表

### 6.2 移动端界面

- 采用Framework7框架实现移动端界面
- 底部导航栏添加"理财"图标
- 采用卡片式布局展示各项功能
- 支持手势操作和下拉刷新
- 图表采用响应式设计，适配不同屏幕尺寸

## 7. 集成方案

### 7.1 与现有系统集成

#### 7.1.1 认证与授权

- 复用现有JWT认证机制
- 为理财模块API添加相应的权限控制

#### 7.1.2 数据存储

- 复用现有数据库连接
- 按照现有模型设计规范扩展数据模型

#### 7.1.3 国际化

- 复用现有国际化框架
- 为理财模块添加相应的语言文件

#### 7.1.4 通知系统

- 复用现有通知机制
- 为预算超支、目标达成等事件添加通知

### 7.2 第三方API集成

#### 7.2.1 投资产品价格更新

- 支持集成第三方金融数据API（如：Alpha Vantage、Yahoo Finance、新浪财经等）
- 实现定时更新投资产品价格

#### 7.2.2 财务数据导入

- 支持从银行、券商等金融机构导入投资数据
- 支持CSV、OFX等格式的投资数据导入

## 8. 实施计划

### 8.1 阶段划分

| 阶段 | 时间 | 主要工作 |
|------|------|----------|
| 需求分析与设计 | 2周 | 需求细化、数据模型设计、API设计、界面设计 |
| 后端开发 | 4周 | 实现数据模型、API接口、业务逻辑 |
| 前端开发 | 4周 | 实现桌面端和移动端界面、交互逻辑 |
| 测试 | 2周 | 单元测试、集成测试、用户测试 |
| 部署与上线 | 1周 | 部署到测试环境、性能测试、正式上线 |

### 8.2 资源需求

- 后端开发人员：1-2名
- 前端开发人员：1-2名
- 测试人员：1名
- 产品经理：1名

## 9. 风险评估

### 9.1 技术风险

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| 第三方API稳定性 | 投资产品价格更新失败 | 实现API重试机制、添加手动更新功能 |
| 数据量增长 | 查询性能下降 | 优化数据库查询、添加缓存机制 |
| 复杂计算性能 | 财务分析计算耗时 | 实现异步计算、优化算法 |

### 9.2 业务风险

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| 用户接受度 | 功能使用率低 | 进行用户调研、优化界面设计、提供使用引导 |
| 功能复杂度 | 用户学习成本高 | 实现渐进式功能展示、提供帮助文档和教程 |

## 10. 验收标准

### 10.1 功能验收

- 所有功能需求均已实现
- 功能符合设计规格
- 系统运行稳定，无重大bug

### 10.2 性能验收

- 页面加载时间 < 2秒
- API响应时间 < 500ms
- 支持1000+并发用户

### 10.3 用户体验验收

- 界面设计符合现有系统风格
- 操作流程简洁直观
- 支持多设备访问
- 提供充分的帮助和引导

## 11. 后续规划

- 支持更多投资产品类型
- 实现投资组合回测功能
- 支持税务计算和申报
- 集成AI智能理财顾问
- 支持家庭共享理财功能

---

**文档版本**：v1.0
**编制日期**：2025-12-29
**编制人**：AI Assistant
