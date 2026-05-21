# AI Agent 智能招聘系统

基于 Gin + gRPC + Vue3 + 阿里云OSS + Eino AI 框架的全栈智能招聘平台。

## 项目简介

本系统是一个完整的智能招聘解决方案，采用两层微服务架构（Gin Web网关 + gRPC Logic业务层），实现HR招聘管理与候选人求职投递的全流程数字化。系统集成私有OSS安全简历存储、Eino AI智能问答（联动MySQL真实数据）等核心能力。

## 技术架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                           前端层                                     │
│  ┌─────────────────────┐        ┌─────────────────────┐            │
│  │   HR 管理端前端      │        │   候选人用户端前端   │            │
│  │   (Vue3 + Vite)      │        │   (Vue3 + Vite)      │            │
│  │   Element Plus       │        │   Element Plus       │            │
│  └─────────────────────┘        └─────────────────────┘            │
└─────────────────────────────────────────────────────────────────────┘
                              │ HTTP REST API
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Web 服务层 (Gin 网关)                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  web-gin-service (Port: 8080)                                │   │
│  │  - CORS 跨域处理                                             │   │
│  │  - JWT 统一鉴权                                              │   │
│  │  - 请求参数校验                                              │   │
│  │  - gRPC 远程调用转发                                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              │ gRPC Protocol
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Logic 业务服务层 (gRPC)                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  logic-grpc-service (Port: 50051)                            │   │
│  │  - 用户权限管控                                              │   │
│  │  - 岗位业务处理                                              │   │
│  │  - 私有 OSS 签名调度                                         │   │
│  │  - MySQL 全量数据读写                                        │   │
│  │  - Eino AI 对话封装（联动数据库真实数据）                     │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│   MySQL 数据库    │ │  阿里云 OSS 存储  │ │  通义千问大模型   │
│   (业务数据)      │ │  (简历文件)       │ │  (Eino AI)       │
└──────────────────┘ └──────────────────┘ └──────────────────┘
```

## 功能模块

### HR 管理端功能

| 功能模块 | 功能描述 |
|----------|----------|
| 账号登录 | 独立账号密码登录，JWT 令牌鉴权，无令牌禁止访问 |
| 岗位管理 | 自主新增、编辑、下架个人发布的招聘岗位 |
| 投递管理 | 分页查看岗位投递候选人、结构化档案、OSS简历签名链接 |
| 候选人档案 | 查看所有候选人档案库，支持学历筛选、技能搜索 |
| 统计分析 | 数据概览、投递趋势分析、岗位热度统计 |
| AI 智能问答 | 基于 Eino 框架的自然语言问答，**联动 MySQL 数据库回答真实数据** |
| 对话历史 | AI 对话自动持久化，页面刷新自动加载历史上下文 |
| 系统通知 | 实时通知推送，新投递提醒 |

### 候选人用户端功能

| 功能模块 | 功能描述 |
|----------|----------|
| 岗位浏览 | 游客免登录浏览全平台公开岗位列表 |
| 账号注册 | 独立账号注册登录 |
| 个人档案 | 结构化档案编辑（姓名、电话、学历、院校、经历、技能） |
| 简历上传 | PDF/DOC/DOCX 格式简历直传私有 OSS，签名访问 |
| 岗位投递 | 档案+简历校验后一键投递 |

## 目录结构

```
final_homework/
├── hr-frontend/              # HR 管理端前端 (Vue3 + Vite)
│   ├── src/
│   │   ├── views/
│   │   │   ├── Login.vue           # 登录页面
│   │   │   └ hr/
│   │   │   │   ├── Layout.vue      # HR 布局框架
│   │   │   │   ├── Dashboard.vue   # HR 工作台（岗位管理+AI问答）
│   │   │   │   ├── Stats.vue       # 统计分析页面
│   │   │   │   ├── Profiles.vue    # 候选人档案库
│   │   │   │   ├── Applications.vue # 投递记录
│   │   │   │   ├── AIHistory.vue   # AI对话历史
│   │   │   │   └── Account.vue     # 账户设置
│   │   │   └ candidate/
│   │   │   │   ├── Layout.vue      # 候选人管理布局
│   │   │   │   └ Jobs.vue          # 投递职位
│   │   ├── router/           # 路由配置
│   │   ├── utils/            # Axios 请求工具
│   └── package.json
│
├── user-frontend/            # 候选人用户端前端 (Vue3 + Vite)
│   ├── src/
│   │   ├── views/
│   │   │   ├── Login.vue     # 登录/注册页面
│   │   │   ├── JobBoard.vue  # 岗位大厅
│   │   │   ├── Profile.vue   # 个人档案编辑
│   │   ├── router/
│   │   ├── utils/
│   └── package.json
│
├── web-gin-service/          # Gin Web 网关服务
│   ├── main.go               # 路由定义 + gRPC 调用转发
│   ├── utils/
│   │   ├── jwt.go            # JWT 生成 + 鉴权中间件
│   │   ├── eino.go           # Eino 初始化
│   ├── pb/                   # Protobuf 生成的 Go 文件
│   └── go.mod
│
├── logic-grpc-service/       # Logic 核心业务 gRPC 服务
│   ├── main.go               # gRPC 服务实现
│   ├── proto/
│   │   └ recruitment.proto   # Protobuf 定义文件
│   ├── db/
│   │   ├── models.go         # GORM 数据模型
│   │   ├── mysql.go          # MySQL 连接初始化
│   ├── utils/
│   │   ├── oss.go            # OSS 上传 + 签名 URL 生成
│   │   ├── eino.go           # Eino AI 对话封装（含数据库上下文）
│   ├── pb/                   # Protobuf 生成的 Go 文件
│   ├── .env.example          # 环境变量配置模板 ⭐
│   ├── config.yaml.example   # YAML 配置模板
│   └── go.mod
│
├── api.md                    # API 接口说明文档
├── db.md                     # 数据库设计文档
├── README.md                 # 项目说明文档
└── answer.md                 # 思考题答案文档
```

## 环境要求

- **Go**: 1.21+
- **Node.js**: 18+
- **MySQL**: 8.0+ 或 SQLite 3（开发环境）
- **阿里云 OSS**: 已创建私有 Bucket
- **通义千问 API**: 阿里云 DashScope API Key

## 配置说明

### 重要：环境变量配置

本项目使用 `.env` 文件管理敏感配置，**所有密钥均通过环境变量读取**，请勿将真实密钥上传到 GitHub。

#### 第一步：创建环境变量文件

克隆项目后，进入 `logic-grpc-service` 目录，复制模板文件：

```bash
cd logic-grpc-service
cp .env.example .env
```

#### 第二步：编辑 `.env` 文件

打开 `.env` 文件，填入你的真实密钥：

```bash
# OSS 阿里云配置
OSS_ENDPOINT=oss-cn-beijing.aliyuncs.com
OSS_ACCESS_KEY_ID=你的阿里云AccessKeyID
OSS_ACCESS_KEY_SECRET=你的阿里云AccessKeySecret
OSS_BUCKET_NAME=你的私有Bucket名称

# AI 大模型配置 (通义千问 DashScope)
DASHSCOPE_API_KEY=你的DashScopeAPIKey
DASHSCOPE_MODEL=qwen-plus

# MySQL 数据库配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=你的数据库密码
MYSQL_DATABASE=smart_recruitment
```

#### 配置项说明

| 变量名 | 说明 | 获取方式 |
|--------|------|----------|
| `OSS_ENDPOINT` | OSS 地域节点 | 阿里云 OSS 控制台 |
| `OSS_ACCESS_KEY_ID` | 阿里云访问密钥ID | 阿里云 RAM 控制台 |
| `OSS_ACCESS_KEY_SECRET` | 阿里云访问密钥Secret | 阿里云 RAM 控制台 |
| `OSS_BUCKET_NAME` | OSS Bucket名称 | 阿里云 OSS 控制台（需创建私有Bucket） |
| `DASHSCOPE_API_KEY` | 通义千问API Key | [阿里云 DashScope](https://dashscope.console.aliyun.com/) |
| `DASHSCOPE_MODEL` | 模型名称 | 默认 `qwen-plus`，可选 `qwen-turbo`、`qwen-max` |

#### 获取密钥指南

**1. 阿里云 OSS 密钥**
- 登录 [阿里云 OSS 控制台](https://oss.console.aliyun.com/)
- 创建私有 Bucket（读写权限选择「私有」）
- 进入 [RAM 控制台](https://ram.console.aliyun.com/) 创建 AccessKey
- 建议使用 RAM 子账号，仅授予 OSS 权限

**2. 通义千问 API Key**
- 登录 [DashScope 控制台](https://dashscope.console.aliyun.com/)
- 开通「通义千问」服务
- 在 API-KEY 管理页面创建新密钥

## 启动部署

### 服务启动顺序

```bash
# 1. 配置环境变量（首次运行必须）
cd logic-grpc-service
cp .env.example .env
# 编辑 .env 填入真实密钥

# 2. 启动 Logic gRPC 核心业务服务（端口 50051）
cd logic-grpc-service
go run main.go

# 3. 启动 Gin Web 网关服务（端口 8080）
cd web-gin-service
go run main.go

# 4. 启动 HR 管理端前端
cd hr-frontend
npm install
npm run dev

# 5. 启动候选人用户端前端
cd user-frontend
npm install
npm run dev
```

### 端口说明

| 服务 | 端口 |
|------|------|
| Logic gRPC 服务 | 50051 |
| Gin Web 网关 | 8080 |
| HR 前端开发服务器 | 5173 |
| 候选人前端开发服务器 | 5174 |

### 测试账号

系统自动创建以下测试账号：

| 角色 | 用户名 | 密码 |
|------|--------|------|
| HR | test_hr | 123456 |
| 候选人 | test_candidate | 123456 |

## API 接口说明

详细接口文档请参考 [api.md](api.md)。

### 核心接口一览

| 接口 | 方法 | 权限 | 描述 |
|------|------|------|------|
| `/api/v1/login` | POST | 公开 | 用户登录 |
| `/api/v1/public/jobs` | GET | 公开 | 游客浏览岗位 |
| `/api/v1/hr/jobs` | POST | HR | 发布岗位 |
| `/api/v1/hr/jobs` | GET | HR | 查看个人岗位 |
| `/api/v1/hr/jobs/:id` | PUT | HR | 编辑/下架岗位 |
| `/api/v1/hr/jobs/:id/applications` | GET | HR | 查看投递详情 |
| `/api/v1/hr/profiles` | GET | HR | 候选人档案库 |
| `/api/v1/hr/stats` | GET | HR | 统计数据 |
| `/api/v1/hr/ai/chat` | POST | HR | AI智能问答（联动数据库） |
| `/api/v1/hr/ai/history` | GET | HR | AI对话历史 |
| `/api/v1/candidate/jobs` | GET | 候选人 | 岗位列表 |
| `/api/v1/candidate/profile` | GET/PUT | 候选人 | 个人档案 |
| `/api/v1/candidate/upload/resume` | POST | 候选人 | 简历上传OSS |
| `/api/v1/candidate/jobs/:id/apply` | POST | 候选人 | 投递岗位 |

## 数据库设计

详细数据库设计请参考 [db.md](db.md)。

### 数据表一览

| 表名 | 描述 |
|------|------|
| `users` | 用户账号表（hr/candidate角色） |
| `jobs` | 招聘岗位表 |
| `candidate_profiles` | 候选人结构化档案表 |
| `job_applications` | 岗位投递关联表 |
| `ai_chat_histories` | AI对话历史记录表 |

## 项目亮点

### 1. 两层微服务架构
采用 Gin Web网关 + gRPC Logic业务层的轻量化架构，Web层仅负责请求转发和鉴权，所有业务逻辑集中在 Logic 层，符合微服务分层设计原则。

### 2. 私有 OSS 签名 URL 安全机制
- 简历直传私有 OSS Bucket，彻底关闭匿名访问
- 动态生成 10 分钟有效期的签名 URL，防盗链保护隐私
- 本地不缓存、不留存任何简历源文件

### 3. Eino AI 框架集成 + 数据库联动
- 基于 CloudWeGo Eino 框架调用通义千问大模型
- **AI 助手可查询 MySQL 真实数据回答问题**（岗位数、投递数、候选人信息等）
- AI 对话自动持久化到 MySQL，绑定 HR 账号
- 前端自动加载历史对话上下文，延续问答交互

### 4. JWT 角色+权限双重校验
- 统一 JWT 令牌签发，24 小时有效期
- 中间件校验令牌 + 角色，HR/候选人路由隔离

### 5. 投递强制校验机制
- 未完善档案 + 未上传简历直接拦截投递请求
- 简历格式严格限制（PDF/DOC/DOCX），后端双重校验文件后缀和大小

### 6. 安全密钥管理
- 所有敏感配置通过 `.env` 环境变量管理
- `.env` 文件已加入 `.gitignore`，防止密钥泄露
- 提供配置模板 `.env.example` 供用户参考

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue3、Vite、Element Plus、Axios |
| Web网关 | Go、Gin、gRPC Client、JWT |
| 业务服务 | Go、gRPC Server、GORM、Eino、godotenv |
| 数据存储 | MySQL / SQLite |
| 文件存储 | 阿里云 OSS（私有 Bucket） |
| AI 大模型 | 通义千问（Eino 框架封装） |

## 开发者

- 姓名：陈榆昊
- 学号：1023004500
- 学校：武汉理工大学
- 方向：全栈开发（偏服务端）

## License

MIT License