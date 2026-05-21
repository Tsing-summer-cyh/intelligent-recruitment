# 智能招聘系统 - 数据库设计文档 (MySQL)

## 1. 账号表 (users)
为了满足“一个用户只有一种角色”的要求，我们将 HR 和候选人存放在同一张表，通过角色字段区分。

| 字段名 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint | PRIMARY KEY, AUTO_INCREMENT | 用户唯一标识 |
| username | varchar(50) | UNIQUE, NOT NULL | 登录账号 |
| password_hash | varchar(255) | NOT NULL | 密码哈希值 |
| role | varchar(20) | NOT NULL | 角色: `hr` 或 `candidate` |
| created_at | datetime | DEFAULT CURRENT_TIMESTAMP | 注册时间 |

## 2. 招聘岗位表 (jobs)
对应 HR 端“自主新增、编辑、下架个人发布的招聘岗位”。

| 字段名 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint | PRIMARY KEY, AUTO_INCREMENT | 岗位唯一标识 |
| hr_id | bigint | NOT NULL | 发布岗位的 HR 账号 ID (关联 users.id) |
| title | varchar(100) | NOT NULL | 岗位名称 |
| description | text | NOT NULL | 岗位描述 |
| status | tinyint | DEFAULT 1 | 状态: 1-上架, 0-下架 |
| created_at | datetime | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | datetime | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

## 3. 候选人结构化档案表 (candidate_profiles)
对应候选人端“必填填写姓名、联系电话、最高学历、毕业院校、工作/项目经历、核心技能标签”。

| 字段名 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint | PRIMARY KEY, AUTO_INCREMENT | 档案唯一标识 |
| user_id | bigint | UNIQUE, NOT NULL | 候选人账号 ID (关联 users.id) |
| real_name | varchar(50) | NOT NULL | 真实姓名 |
| phone | varchar(20) | NOT NULL | 联系电话 |
| highest_education| varchar(50) | NOT NULL | 最高学历 |
| university | varchar(100) | NOT NULL | 毕业院校 |
| experience | text | NOT NULL | 工作/项目经历 |
| skills | varchar(255) | NOT NULL | 核心技能标签 (逗号分隔或 JSON) |
| resume_oss_url | varchar(500)| DEFAULT NULL | 私有 OSS 简历签名访问链接 |
| created_at | datetime | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | datetime | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

## 4. 岗位投递关联表 (job_applications)
用于记录谁投递了什么岗位，并支撑 AI 查询“单岗位投递统计、岗位热度数据”。

| 字段名 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint | PRIMARY KEY, AUTO_INCREMENT | 投递记录唯一标识 |
| job_id | bigint | NOT NULL | 岗位 ID (关联 jobs.id) |
| candidate_id | bigint | NOT NULL | 候选人账号 ID (关联 users.id) |
| created_at | datetime | DEFAULT CURRENT_TIMESTAMP | 投递时间 |

## 5. AI 对话历史记录表 (ai_chat_histories)
对应 Eino 框架要求：“每一条 HR 提问、AI 回复自动成对存入 MySQL 对话历史表，绑定 HR 管理员账号 ID”。

| 字段名 | 类型 | 约束 | 说明 |
| --- | --- | --- | --- |
| id | bigint | PRIMARY KEY, AUTO_INCREMENT | 记录唯一标识 |
| hr_id | bigint | NOT NULL | 提问的 HR 账号 ID (关联 users.id) |
| question | text | NOT NULL | HR 输入的自然语言提问 |
| answer | text | NOT NULL | 大模型返回的统计回答 |
| created_at | datetime | DEFAULT CURRENT_TIMESTAMP | 对话发生时间 |