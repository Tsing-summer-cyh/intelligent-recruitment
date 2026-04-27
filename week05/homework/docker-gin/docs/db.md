# 数据库设计文档 (MySQL)

本项目数据库名为 `smart_vocab`，字符集统一使用 `utf8mb4` 以支持多语言及特殊符号。包含两张核心业务表：用户认证表 (`users`) 和 单词学习记录表 (`words`)。

## 关联关系设计说明
**1 对 N 关系**：一个用户 (`User`) 可以拥有多个单词记录 (`Word`)。
在 `words` 表中设计了 `user_id` 作为外键，关联至 `users` 表的 `id` 主键。同时设置了 `ON DELETE CASCADE` 级联删除约束，确保用户账号被注销（物理删除）时，其名下的单词本记录一并清除。

---

## 1. 用户认证表 (`users`)

负责存储用户的基本注册信息与鉴权凭证。

| 字段名 | 数据类型 | 约束条件 | 索引 | 业务含义 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | 主键索引 | 用户唯一标识，自增 ID |
| `username` | VARCHAR(50) | NOT NULL, UNIQUE | 唯一索引 | 登录用户名，不允许重复 |
| `password_hash`| VARCHAR(255) | NOT NULL | - | bcrypt 加密后的密码散列值，严禁明文 |
| `created_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP | - | 账号注册时间 |
| `updated_at` | DATETIME(3) | DEFAULT ON UPDATE | - | 账号最后更新时间 |
| `deleted_at` | DATETIME(3) | DEFAULT NULL | 普通索引 | 软删除标记时间 (GORM 规范) |

---

## 2. 单词学习记录表 (`words`)

核心业务表，用于存储用户查询并保存的生词、AI 释义以及生成的例句。

| 字段名 | 数据类型 | 约束条件 | 索引 | 业务含义 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | BIGINT UNSIGNED | PRIMARY KEY, AUTO_INCREMENT | 主键索引 | 单词记录唯一标识，自增 ID |
| `user_id` | BIGINT UNSIGNED | NOT NULL | 外键索引 | 关联 `users` 表的 `id`，标明单词归属人 |
| `word` | VARCHAR(100) | NOT NULL | - | 查询的目标单词 |
| `meaning` | TEXT | NOT NULL | - | AI 返回的详细中文精准释义 |
| `sentences` | JSON | NOT NULL | - | AI 生成的 3 条语境例句，以 JSON 数组格式持久化存储 |
| `ai_provider` | VARCHAR(50) | NOT NULL | - | 提供本次翻译服务的大模型来源 (如 DeepSeek) |
| `created_at` | DATETIME(3) | DEFAULT CURRENT_TIMESTAMP | - | 该单词收录进生词本的时间 |
| `updated_at` | DATETIME(3) | DEFAULT ON UPDATE | - | 记录最后更新时间 |
| `deleted_at` | DATETIME(3) | DEFAULT NULL | 普通索引 | 软删除标记时间，实现用户层面的“移出单词本”但数据可找回 |

**外键约束**: 
`CONSTRAINT fk_words_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`