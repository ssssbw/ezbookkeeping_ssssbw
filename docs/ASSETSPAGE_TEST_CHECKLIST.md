# AssetsPage 功能测试清单

> 最后更新：2026-06-05

---

## 测试入口

- 前端开发服务器：`npm run serve`
- 后端热重载：`air`（macOS）/ `air -c .air.windows.toml`（Windows）
- 路由：`/#/investment/assets`（需先登录，切换到理财模式）

---

## 需要关心的文件

### 前端（主要）

| 文件 | 说明 |
|------|------|
| `src/views/desktop/investment/AssetsPage.vue` | **主文件**，搜索、持仓/自选表格、买卖对话框、详情对话框、管理弹窗 CRUD |
| `src/models/investment.ts` | 数据类型定义（AssetModifyRequest、AssetDeleteRequest、AssetListRequest 等） |
| `src/lib/services.ts` | API 调用（listGlobalAssets、modifyGlobalAsset、deleteGlobalAsset 等，约 L1008-1030） |
| `src/locales/en.json` | 英文多语言（搜索 `asset.` 嵌套对象） |
| `src/locales/zh_Hans.json` | 中文多语言（搜索 `asset.` 嵌套对象） |
| `src/components/desktop/PaginationButtons.vue` | 分页组件（持仓/自选共用） |
| `src/views/desktop/MainLayout.vue` | 导航栏标题同步（watch currentRoutePath） |

### 后端（主要）

| 文件 | 说明 |
|------|------|
| `pkg/models/asset.go` | Asset model + Request/Response（AssetDeleteRequest、AssetListResponse） |
| `pkg/services/asset.go` | Asset 服务（DeleteAsset、GetAllAssetsCount、GetAllAssets keyword 过滤） |
| `pkg/api/investment.go` | Handler（GlobalAssetListHandler L497、ModifyHandler L536、DeleteHandler L568） |
| `cmd/webserver.go` | 路由注册（L483-488，global_assets 6 条路由） |

### 间接相关

| 文件 | 说明 |
|------|------|
| `src/stores/investment.ts` | Pinia store，isInvestmentMode 状态 |
| `pkg/models/user_asset.go` | 用户持仓 model |
| `pkg/services/user_asset.go` | 用户持仓服务（自选/持仓列表） |
| `pkg/services/investment_transaction.go` | 买卖交易服务 |
| `docs/AI_CONTEXT_HANDOFF.md` | 开发交接文档 |

---

## 测试项

### 一、搜索功能
- [ ] 搜索框输入代码/名称，下拉浮层显示匹配结果
- [ ] 清除搜索内容，浮层消失
- [ ] 点击搜索结果，弹出行情详情对话框
- [ ] 搜索结果中点击「+自选」按钮，添加到自选 Tab
- [ ] 搜索结果中点击「买入」按钮，弹出买入对话框
- [ ] 已在自选的资产不显示「+自选」按钮
- [ ] 搜索框左側有 mdiMagnify 图标，点击清除按钮可清空

### 二、持仓 Tab
- [ ] 表格列完整：代码、名称、市场、行业（▾）、数量、现价、市值、成本、浮盈亏、收益率、操作
- [ ] 表头无排序箭头（sortable: false）
- [ ] 点击行弹出详情对话框
- [ ] 操作列：买入（绿）、卖出（红）、更多（⋯）
- [ ] 持仓汇总栏（表格下方）：总市值、总成本、浮盈亏、收益率
- [ ] 浮盈亏/收益率颜色随正负变化（正=绿，负=红）
- [ ] Tab chip 数字正确显示持仓条数
- [ ] 无数据时显示空状态提示
- [ ] 分页组件正常翻页，页码 ≥ 1（即使无数据也显示「1」）
- [ ] 数据为 null 时显示「--」而非空白

### 三、自选 Tab
- [ ] 表格列完整：代码、名称、市场、行业（▾）、现价、操作
- [ ] 操作列：买入（绿）、移除（红）、更多（⋯）
- [ ] 「移除」成功后从自选 Tab 消失，chip 数字更新
- [ ] Tab chip 数字正确显示自选条数
- [ ] 分页正常
- [ ] 无数据时显示空状态

### 四、买入/卖出对话框
- [ ] 买入对话框：正确预填资产名称、代码
- [ ] 卖出对话框：正确预填资产名称、代码
- [ ] 提交后交易记录出现，持仓数量/成本更新
- [ ] 关闭对话框不提交
- [ ] 提交中按钮 loading 状态

### 五、详情对话框
- [ ] 显示字段：名称、代码、市场、分类（chip）、行业、现价、数量、成本、浮盈亏、收益率
- [ ] 浮盈亏/收益率颜色（正绿负红）

### 六、全局资产管理 — 管理按钮
- [ ] 非管理员登录，搜索框右侧不显示「全局资产管理」按钮
- [ ] 管理员登录，搜索框右侧显示「全局资产管理」按钮

### 七、全局资产管理弹窗 — 列表
- [ ] 弹窗打开后表格加载全局资产（代码、名称、市场、行业、操作）
- [ ] 搜索框按代码/名称过滤，回车或输入后触发搜索
- [ ] 行业下拉筛选
- [ ] 分页正常

### 八、全局资产管理弹窗 — 新增
- [ ] 点击「新增」按钮弹出表单
- [ ] 表单字段：代码、名称、市场（下拉）、分类（下拉）、行业（下拉）
- [ ] 提交成功后列表刷新，新资产出现
- [ ] 必填校验（空代码/名称提交时应报错）

### 九、全局资产管理弹窗 — 编辑
- [ ] 点击行编辑按钮弹出预填表单
- [ ] 修改后保存成功，列表刷新

### 十、全局资产管理弹窗 — 删除
- [ ] 点击行删除按钮弹出确认框
- [ ] 确认后资产被删除，列表刷新
- [ ] 取消后无变化

### 十一、导航 & 路由
- [ ] 进入 `/#/investment/assets`，导航栏标题显示「资产管理」
- [ ] 切换到其他页面再回来，标题同步更新
- [ ] isInvestmentMode 由 URL 决定，不从 localStorage 读取

### 十二、多语言
- [ ] 中文界面所有文本正确，无 key 裸露
- [ ] 英文界面所有文本正确，无 key 裸露
- [ ] 行业名称（Technology、Healthcare 等）中英文都正确

---

## 常见问题排查

| 问题 | 排查方向 |
|------|----------|
| 搜索无结果 | 检查全局资产表是否有数据（先通过管理弹窗新增） |
| API 返回错误 | 查看后端日志，检查路由是否注册（`cmd/webserver.go` L483-488） |
| 管理按钮不显示 | 确认当前用户是管理员（`checkAdmin` 调用 `/investment/admin/check.json`） |
| 分页组件不显示页码 | 检查 `totalPageCount` 是否为 0，代码已做 `Math.max(1, ...)` 兜底 |
| 行业显示为原始值 | 检查 `formatIndustry()` 函数是否覆盖了该行业枚举 |
| i18n key 裸露 | 检查 `en.json` / `zh_Hans.json` 末尾的 `asset` 嵌套对象 |
