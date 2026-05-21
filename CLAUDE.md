# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

ezBookkeeping - 轻量级、自托管的个人财务应用。前后端分离架构：
- **后端**：Go + Gin + xorm ORM，REST API
- **前端**：Vue 3 + TypeScript + Vite，支持桌面和移动双模式 UI

## 构建与开发命令

```bash
# 构建
./build.sh backend              # 构建后端二进制
./build.sh frontend             # 构建前端
./build.sh package -o out.tar.gz  # 打包
./build.sh docker               # Docker 镜像

# 跳过检查（构建时默认执行 lint 和 test）
./build.sh backend --no-lint --no-test

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

### 数据库

支持 SQLite、MySQL、PostgreSQL。配置文件 `conf/ezbookkeeping.ini`。

## 关键技术细节

- ORM 使用 xorm，模型定义在 `pkg/models/`，通过 struct tags 映射数据库字段
- API 路由在 `cmd/webserver.go` 中统一注册
- 前端状态通过 Pinia stores 管理，stores 位于 `src/stores/`
- 前端 UI 框架使用 Vuetify 3 + Framework7（移动版）
- 图表使用 ECharts
