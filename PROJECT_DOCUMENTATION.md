# ezBookkeeping 项目文档

## 1. 项目概述

ezBookkeeping 是一款轻量级、自托管的个人财务管理应用，具有用户友好的界面和强大的记账功能。它易于部署，只需一个 Docker 命令即可启动，设计资源高效且高度可扩展，可在树莓派等小型设备上平稳运行，也可扩展到 NAS、微服务器甚至大型集群环境。

### 核心优势
- **开源且自托管**：注重隐私和控制权
- **轻量级且快速**：性能优化，即使在低资源环境下也能流畅运行
- **易于安装**：支持 Docker、多数据库（SQLite、MySQL、PostgreSQL）
- **跨平台**：支持 Windows、macOS、Linux，适用于 x86、amd64、ARM 架构
- **用户友好界面**：针对移动和桌面设备的定制化界面
- **PWA 支持**：可添加到移动主屏幕，像原生应用一样使用

## 2. 技术栈

### 后端
- **语言**：Go 1.25
- **Web 框架**：Gin
- **ORM**：XORM
- **数据库支持**：SQLite3、MySQL、PostgreSQL
- **存储**：本地文件系统、MinIO、WebDAV
- **认证**：JWT、OAuth2、双因素认证

### 前端
- **框架**：Vue 3 + TypeScript
- **构建工具**：Vite
- **UI 库**：Vuetify（桌面）、Framework7（移动）
- **图表**：ECharts
- **状态管理**：Pinia
- **路由**：Vue Router

## 3. 架构设计

### 系统架构

ezBookkeeping 采用前后端分离架构：

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    前端应用     │     │    后端 API     │     │    数据存储     │
│  (Vue 3 + TS)   │────▶│  (Go + Gin)     │────▶│  (MySQL/PostgreSQL) │
└─────────────────┘     └─────────────────┘     └─────────────────┘
           │                         │     ┌─────────────────┐
           │                         └────▶│  对象存储       │
           │                               │  (本地/MinIO/WebDAV) │
           └───────────────────────────────▶│                 │
                                           └─────────────────┘
```

### 核心模块

1. **认证与授权模块**：处理用户登录、注册、JWT 令牌生成与验证
2. **账户管理模块**：管理用户账户和子账户
3. **交易管理模块**：处理交易记录、分类、标签等
4. **数据导入导出模块**：支持多种格式的数据导入导出
5. **AI 集成模块**：通过 LLM 实现收据图像识别
6. **MCP 服务器模块**：支持 AI/LLM 访问
7. **定时任务模块**：处理定期交易创建等

## 4. 目录结构

```
├── .gitea/              # Gitea CI/CD 配置
├── .github/             # GitHub CI/CD 配置
├── cmd/                 # 命令行工具实现
├── conf/                # 配置文件目录
├── docker/              # Docker 相关配置
├── etc/                 # 系统服务配置
├── pkg/                 # 后端核心包
├── public/              # 静态资源
├── src/                 # 前端源码
├── .editorconfig        # 编辑器配置
├── .gitignore           # Git 忽略文件
├── Dockerfile           # Docker 构建文件
├── LICENSE              # 许可证
├── README.md            # 项目说明
├── build.bat            # Windows 构建脚本
├── build.ps1            # PowerShell 构建脚本
├── build.sh             # Linux/macOS 构建脚本
├── docker-bake.hcl      # Docker 构建配置
├── ezbookkeeping.go     # 后端主入口
├── go.mod               # Go 依赖管理
├── go.sum               # Go 依赖校验
├── package-lock.json    # NPM 依赖锁定
├── package.json         # NPM 依赖管理
└── postcss.config.js    # PostCSS 配置
```

## 5. 主要文件功能说明

### 后端核心文件

| 文件路径 | 主要功能 |
|---------|--------|
| `ezbookkeeping.go` | 后端主入口文件，处理命令行参数和启动服务 |
| `cmd/webserver.go` | Web 服务器启动和配置 |
| `cmd/database.go` | 数据库初始化和迁移 |
| `cmd/security.go` | 安全相关命令（如生成密钥） |
| `conf/ezbookkeeping.ini` | 主配置文件，包含所有系统配置项 |
| `pkg/api/` | REST API 实现，包含所有 API 端点 |
| `pkg/core/application.go` | 应用核心逻辑，处理应用生命周期 |
| `pkg/models/` | 数据模型定义，包含所有数据库表结构 |
| `pkg/services/` | 业务逻辑层，处理核心业务逻辑 |
| `pkg/storage/` | 存储服务实现，支持多种存储类型 |
| `pkg/llm/` | LLM 集成，处理 AI 图像识别 |

### 前端核心文件

| 文件路径 | 主要功能 |
|---------|--------|
| `src/DesktopApp.vue` | 桌面端主应用组件 |
| `src/MobileApp.vue` | 移动端主应用组件 |
| `src/router/` | 路由配置，处理页面导航 |
| `src/stores/` | Pinia 状态管理，处理全局状态 |
| `src/components/` | 可复用组件，分为基础、通用、桌面和移动组件 |
| `src/lib/` | 核心工具库，包含各种工具函数 |
| `src/locales/` | 多语言支持，包含各种语言的翻译文件 |

## 6. 核心功能模块

### 6.1 账户管理

- 支持多级账户和子账户
- 账户余额实时计算
- 账户分类和图标自定义
- 支持多种货币

### 6.2 交易管理

- 交易记录的创建、编辑、删除
- 支持交易分类和标签
- 交易图片上传和管理
- 交易模板和定期交易
- 交易位置跟踪与地图集成
- 重复交易检测

### 6.3 数据导入导出

- 支持多种格式导入：CSV、TSV、Beancount、CAMT、QIF、OFX 等
- 支持导出为 CSV 和 TSV
- 智能字段映射和数据验证

### 6.4 AI 集成

- 支持多种 LLM 提供商：OpenAI、Google AI、Ollama 等
- 收据图像识别，自动提取交易信息
- MCP（Model Context Protocol）支持，允许 AI/LLM 访问数据

### 6.5 认证与安全

- JWT 令牌认证
- 双因素认证
- 邮件验证
- OAuth2 集成（OIDC、Nextcloud、Gitea、GitHub）
- 密码重置功能
- 登录失败次数限制

### 6.6 国际化

- 支持 15+ 种语言
- 本地化货币和日期格式
- 自动汇率更新（来自各国央行）

## 7. 数据结构详细解析

### 7.1 账户相关数据结构

#### 账户分类 (AccountCategory)
```go
// AccountCategory represents account category
type AccountCategory byte

// Account categories
const (
    ACCOUNT_CATEGORY_CASH                   AccountCategory = 1
    ACCOUNT_CATEGORY_CHECKING_ACCOUNT       AccountCategory = 2
    ACCOUNT_CATEGORY_CREDIT_CARD            AccountCategory = 3
    ACCOUNT_CATEGORY_VIRTUAL                AccountCategory = 4
    ACCOUNT_CATEGORY_DEBT                   AccountCategory = 5
    ACCOUNT_CATEGORY_RECEIVABLES            AccountCategory = 6
    ACCOUNT_CATEGORY_INVESTMENT             AccountCategory = 7
    ACCOUNT_CATEGORY_SAVINGS_ACCOUNT        AccountCategory = 8
    ACCOUNT_CATEGORY_CERTIFICATE_OF_DEPOSIT AccountCategory = 9
)
```

#### 账户类型 (AccountType)
```go
// AccountType represents account type
type AccountType byte

// Account types
const (
    ACCOUNT_TYPE_SINGLE_ACCOUNT     AccountType = 1
    ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS AccountType = 2
)
```

#### 账户模型 (Account)
```go
// Account represents account data stored in database
type Account struct {
    AccountId       int64           `xorm:"PK"`
    Uid             int64           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
    Deleted         bool            `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
    Category        AccountCategory `xorm:"NOT NULL"`
    Type            AccountType     `xorm:"NOT NULL"`
    ParentAccountId int64           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
    Name            string          `xorm:"VARCHAR(64) NOT NULL"`
    DisplayOrder    int32           `xorm:"INDEX(IDX_account_uid_deleted_parent_account_id_order) NOT NULL"`
    Icon            int64           `xorm:"NOT NULL"`
    Color           string          `xorm:"VARCHAR(6) NOT NULL"`
    Currency        string          `xorm:"VARCHAR(3) NOT NULL"`
    Balance         int64           `xorm:"NOT NULL"`
    Comment         string          `xorm:"VARCHAR(255) NOT NULL"`
    Extend          *AccountExtend  `xorm:"BLOB"`
    Hidden          bool            `xorm:"NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

#### 账户扩展信息 (AccountExtend)
```go
// AccountExtend represents account extend data stored in database
type AccountExtend struct {
    CreditCardStatementDate *int `json:"creditCardStatementDate"`
}
```

### 7.2 交易相关数据结构

#### 交易类型 (TransactionType)
```go
// TransactionType represents transaction type
type TransactionType byte

// Transaction types
const (
    TRANSACTION_TYPE_MODIFY_BALANCE TransactionType = 1
    TRANSACTION_TYPE_INCOME         TransactionType = 2
    TRANSACTION_TYPE_EXPENSE        TransactionType = 3
    TRANSACTION_TYPE_TRANSFER       TransactionType = 4
)
```

#### 交易数据库类型 (TransactionDbType)
```go
// TransactionDbType represents transaction type in database
type TransactionDbType byte

// Transaction db types
const (
    TRANSACTION_DB_TYPE_MODIFY_BALANCE TransactionDbType = 1
    TRANSACTION_DB_TYPE_INCOME         TransactionDbType = 2
    TRANSACTION_DB_TYPE_EXPENSE        TransactionDbType = 3
    TRANSACTION_DB_TYPE_TRANSFER_OUT   TransactionDbType = 4
    TRANSACTION_DB_TYPE_TRANSFER_IN    TransactionDbType = 5
)
```

#### 交易模型 (Transaction)
```go
// Transaction represents transaction data stored in database
type Transaction struct {
    TransactionId        int64             `xorm:"PK"`
    Uid                  int64             `xorm:"UNIQUE(UQE_transaction_uid_time) INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) INDEX(IDX_transaction_uid_deleted_time_longitude_latitude) NOT NULL"`
    Deleted              bool              `xorm:"INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) INDEX(IDX_transaction_uid_deleted_time_longitude_latitude) NOT NULL"`
    Type                 TransactionDbType `xorm:"INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) NOT NULL"`
    CategoryId           int64             `xorm:"INDEX(IDX_transaction_uid_deleted_category_id_time) NOT NULL"`
    AccountId            int64             `xorm:"INDEX(IDX_transaction_uid_deleted_account_id_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) NOT NULL"`
    TransactionTime      int64             `xorm:"UNIQUE(UQE_transaction_uid_time) INDEX(IDX_transaction_uid_deleted_time) INDEX(IDX_transaction_uid_deleted_type_time) INDEX(IDX_transaction_uid_deleted_type_account_id_time) INDEX(IDX_transaction_uid_deleted_category_id_time) INDEX(IDX_transaction_uid_deleted_account_id_time) NOT NULL"`
    TimezoneUtcOffset    int16             `xorm:"NOT NULL"`
    Amount               int64             `xorm:"NOT NULL"`
    RelatedId            int64             `xorm:"NOT NULL"`
    RelatedAccountId     int64             `xorm:"NOT NULL"`
    RelatedAccountAmount int64             `xorm:"NOT NULL"`
    HideAmount           bool              `xorm:"NOT NULL"`
    Comment              string            `xorm:"VARCHAR(255) NOT NULL"`
    GeoLongitude         float64           `xorm:"INDEX(IDX_transaction_uid_deleted_time_longitude_latitude)"`
    GeoLatitude          float64           `xorm:"INDEX(IDX_transaction_uid_deleted_time_longitude_latitude)"`
    CreatedIp            string            `xorm:"VARCHAR(39)"`
    ScheduledCreated     bool
    CreatedUnixTime      int64
    UpdatedUnixTime      int64
    DeletedUnixTime      int64
}
```

### 7.3 分类相关数据结构

#### 交易分类类型 (TransactionCategoryType)
```go
// TransactionCategoryType represents transaction category type
type TransactionCategoryType byte

// Transaction category types
const (
    CATEGORY_TYPE_INCOME   TransactionCategoryType = 1
    CATEGORY_TYPE_EXPENSE  TransactionCategoryType = 2
    CATEGORY_TYPE_TRANSFER TransactionCategoryType = 3
)
```

#### 交易分类模型 (TransactionCategory)
```go
// TransactionCategory represents transaction category data stored in database
type TransactionCategory struct {
    CategoryId       int64                   `xorm:"PK"`
    Uid              int64                   `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
    Deleted          bool                    `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
    Type             TransactionCategoryType `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
    ParentCategoryId int64                   `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
    Name             string                  `xorm:"VARCHAR(64) NOT NULL"`
    DisplayOrder     int32                   `xorm:"INDEX(IDX_category_uid_deleted_type_parent_category_id_order) NOT NULL"`
    Icon             int64                   `xorm:"NOT NULL"`
    Color            string                  `xorm:"VARCHAR(6) NOT NULL"`
    Hidden           bool                    `xorm:"NOT NULL"`
    Comment          string                  `xorm:"VARCHAR(255) NOT NULL"`
    CreatedUnixTime  int64
    UpdatedUnixTime  int64
    DeletedUnixTime  int64
}
```

### 7.4 标签相关数据结构

#### 交易标签模型 (TransactionTag)
```go
// TransactionTag represents transaction tag data stored in database
type TransactionTag struct {
    TagId        int64  `xorm:"PK"`
    Uid          int64  `xorm:"INDEX(IDX_tag_uid_deleted_group_id_order) NOT NULL"`
    Deleted      bool   `xorm:"INDEX(IDX_tag_uid_deleted_group_id_order) NOT NULL"`
    GroupId      int64  `xorm:"INDEX(IDX_tag_uid_deleted_group_id_order) NOT NULL"`
    Name         string `xorm:"VARCHAR(32) NOT NULL"`
    DisplayOrder int32  `xorm:"INDEX(IDX_tag_uid_deleted_group_id_order) NOT NULL"`
    Color        string `xorm:"VARCHAR(6) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

#### 交易标签组模型 (TransactionTagGroup)
```go
// TransactionTagGroup represents transaction tag group data stored in database
type TransactionTagGroup struct {
    GroupId      int64  `xorm:"PK"`
    Uid          int64  `xorm:"INDEX(IDX_tag_group_uid_deleted_order) NOT NULL"`
    Deleted      bool   `xorm:"INDEX(IDX_tag_group_uid_deleted_order) NOT NULL"`
    Name         string `xorm:"VARCHAR(32) NOT NULL"`
    DisplayOrder int32  `xorm:"INDEX(IDX_tag_group_uid_deleted_order) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

### 7.5 用户相关数据结构

#### 用户模型 (User)
```go
// User represents user data stored in database
type User struct {
    Uid                 int64  `xorm:"PK"`
    Username            string `xorm:"UNIQUE(UQE_user_username) VARCHAR(32) NOT NULL"`
    Email               string `xorm:"UNIQUE(UQE_user_email) VARCHAR(128) NOT NULL"`
    PasswordHash        string `xorm:"VARCHAR(128) NOT NULL"`
    Salt                string `xorm:"VARCHAR(16) NOT NULL"`
    Language            string `xorm:"VARCHAR(10) NOT NULL"`
    TimeZone            string `xorm:"VARCHAR(32) NOT NULL"`
    Currency            string `xorm:"VARCHAR(3) NOT NULL"`
    ExpenseAmountColor  string `xorm:"VARCHAR(6) NOT NULL"`
    IncomeAmountColor   string `xorm:"VARCHAR(6) NOT NULL"`
    InitialLogin        bool   `xorm:"NOT NULL"`
    EmailVerified       bool   `xorm:"NOT NULL"`
    TwoFactorAuthEnabled bool  `xorm:"NOT NULL"`
    CreatedUnixTime     int64
    UpdatedUnixTime     int64
}
```

### 7.6 交易图片相关数据结构

#### 交易图片信息模型 (TransactionPictureInfo)
```go
// TransactionPictureInfo represents transaction picture info data stored in database
type TransactionPictureInfo struct {
    PictureId        int64  `xorm:"PK"`
    Uid              int64  `xorm:"INDEX(IDX_picture_uid) NOT NULL"`
    TransactionId    int64  `xorm:"INDEX(IDX_picture_transaction_id) NOT NULL"`
    StorageKey       string `xorm:"VARCHAR(255) NOT NULL"`
    OriginalFileName string `xorm:"VARCHAR(255) NOT NULL"`
    Width            int    `xorm:"NOT NULL"`
    Height           int    `xorm:"NOT NULL"`
    Size             int64  `xorm:"NOT NULL"`
    UploadedUnixTime int64
}
```

### 7.7 交易模板相关数据结构

#### 交易模板模型 (TransactionTemplate)
```go
// TransactionTemplate represents transaction template data stored in database
type TransactionTemplate struct {
    TemplateId      int64  `xorm:"PK"`
    Uid             int64  `xorm:"INDEX(IDX_template_uid_deleted) NOT NULL"`
    Deleted         bool   `xorm:"INDEX(IDX_template_uid_deleted) NOT NULL"`
    Name            string `xorm:"VARCHAR(64) NOT NULL"`
    Type            TransactionType `xorm:"NOT NULL"`
    CategoryId      int64  `xorm:"NOT NULL"`
    SourceAccountId int64  `xorm:"NOT NULL"`
    DestinationAccountId int64 `xorm:"NOT NULL"`
    Amount          int64  `xorm:"NOT NULL"`
    Comment         string `xorm:"VARCHAR(255) NOT NULL"`
    Icon            int64  `xorm:"NOT NULL"`
    Color           string `xorm:"VARCHAR(6) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
    DeletedUnixTime int64
}
```

### 7.8 汇率相关数据结构

#### 汇率模型 (ExchangeRate)
```go
// ExchangeRate represents exchange rate data stored in database
type ExchangeRate struct {
    FromCurrency string  `xorm:"PK(VARCHAR(3))"`
    ToCurrency   string  `xorm:"PK(VARCHAR(3))"`
    Rate         float64 `xorm:"NOT NULL"`
    UpdatedUnixTime int64
}
```

### 7.9 认证相关数据结构

#### 令牌记录模型 (TokenRecord)
```go
// TokenRecord represents token record data stored in database
type TokenRecord struct {
    TokenId      string `xorm:"PK(VARCHAR(64))"`
    Uid          int64  `xorm:"INDEX(IDX_token_uid) NOT NULL"`
    Type         int    `xorm:"NOT NULL"`
    ExpiredUnixTime int64
    CreatedUnixTime int64
}
```

#### 双因素认证模型 (TwoFactor)
```go
// TwoFactor represents two-factor authentication data stored in database
type TwoFactor struct {
    Uid        int64  `xorm:"PK"`
    Secret     string `xorm:"VARCHAR(128) NOT NULL"`
    CreatedUnixTime int64
    UpdatedUnixTime int64
}
```

## 8. 配置说明

### 8.1 主要配置文件

配置文件位于 `conf/ezbookkeeping.ini`，包含以下主要部分：

- **[global]**：全局配置，如运行模式
- **[server]**：服务器配置，如协议、端口、域名等
- **[database]**：数据库配置，如类型、连接信息等
- **[storage]**：存储配置，如存储类型、路径等
- **[llm]**：LLM 配置，如 AI 提供商、API 密钥等
- **[security]**：安全配置，如密钥、令牌过期时间等
- **[auth]**：认证配置，如内部认证、OAuth2 等
- **[user]**：用户配置，如注册、邮件验证等

### 8.2 关键配置项

| 配置项 | 说明 | 默认值 |
|-------|------|--------|
| `server.protocol` | 服务器协议（http, https, socket） | http |
| `server.http_port` | HTTP 端口 | 8080 |
| `database.type` | 数据库类型（mysql, postgres, sqlite3） | sqlite3 |
| `storage.type` | 存储类型（local_filesystem, minio, webdav） | local_filesystem |
| `security.secret_key` | 签名密钥，必须修改以保证安全 | - |
| `user.enable_register` | 是否允许用户注册 | true |
| `user.enable_email_verify` | 是否启用邮件验证 | false |

## 9. 部署方式

### 9.1 Docker 部署

**最新稳定版**：
```bash
docker run -p8080:8080 mayswind/ezbookkeeping
```

**最新每日构建**：
```bash
docker run -p8080:8080 mayswind/ezbookkeeping:latest-snapshot
```

### 9.2 二进制部署

1. 从 [GitHub Releases](https://github.com/mayswind/ezbookkeeping/releases) 下载最新版本
2. 解压并运行：
   - Linux/macOS: `./ezbookkeeping server run`
   - Windows: `.ezbookkeeping.exe server run`

### 9.3 源码构建

1. 安装依赖：Go 1.25+、GCC、Node.js、NPM
2. 克隆代码：`git clone https://github.com/mayswind/ezbookkeeping.git`
3. 构建：
   - Linux/macOS: `./build.sh package -o ezbookkeeping.tar.gz`
   - Windows: `./build.bat package -o ezbookkeeping.zip`

## 10. 开发流程

### 10.1 后端开发

1. 安装 Go 1.25+ 和依赖管理工具
2. 安装依赖：`go mod download`
3. 运行开发服务器：`go run ezbookkeeping.go server run`
4. 运行测试：`go test ./pkg/...`

### 10.2 前端开发

1. 安装 Node.js 和 NPM
2. 安装依赖：`npm install`
3. 运行开发服务器：`npm run dev`
4. 构建生产版本：`npm run build`
5. 运行测试：`npm test`

## 11. 安全特性

### 11.1 数据安全

- 数据库连接加密（支持 TLS）
- 敏感数据加密存储
- JWT 令牌签名验证
- 密码哈希存储（bcrypt）

### 11.2 访问控制

- 基于角色的访问控制
- API 限流
- IP 访问限制
- 登录失败次数限制

### 11.3 安全更新

- 定期更新依赖库
- 及时修复安全漏洞
- 提供安全更新通知

## 12. 贡献指南

### 12.1 代码贡献

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/your-feature`
3. 提交更改：`git commit -m "Add your feature"`
4. 推送到分支：`git push origin feature/your-feature`
5. 创建 Pull Request

### 12.2 翻译贡献

1. 编辑或添加语言文件：`src/locales/xx.json`
2. 提交更改并创建 Pull Request

### 12.3 报告问题

- 使用 GitHub Issues 报告 Bug
- 提供详细的问题描述和复现步骤
- 包含相关日志和截图

## 13. 许可证

本项目采用 [MIT 许可证](LICENSE)，详见 LICENSE 文件。

## 14. 联系方式

- 项目主页：https://github.com/mayswind/ezbookkeeping
- 官方文档：https://ezbookkeeping.mayswind.net
- 演示地址：https://ezbookkeeping-demo.mayswind.net

## 15. 更新日志

### 最新版本

**v1.0.0** (2025-12-29)

- 初始版本发布
- 支持基本的记账功能
- 支持多语言和多货币
- 支持 Docker 部署
- 支持 PWA

### 后续计划

- 增强 AI 功能
- 支持更多数据源
- 改进报告和分析功能
- 增强移动体验
- 支持更多集成

---

**ezBookkeeping** - 轻量级、自托管的个人财务管理应用