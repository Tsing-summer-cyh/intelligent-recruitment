# 智能招聘系统 - API 接口说明文档

## 全局规范
* **Base URL**: `http://localhost:8080/api/v1`
* **数据交互格式**: `application/json`
* **统一返回结构 (JSON)**:
  ```json
  {
    "code": 0,          // 0 表示成功，非 0 表示业务错误
    "msg": "success",   // 提示信息
    "data": {}          // 具体返回数据，失败时为空或null
  }
  ```
* **鉴权方式**: 除公开接口（如游客浏览岗位、登录）外，所有请求头必须携带 JWT 令牌。
  * 格式：`Authorization: Bearer <your_jwt_token>`

---

## 1. 账号与通用模块

### 1.1 用户登录与注册
* **路径**: `/login`
* **方法**: `POST`
* **说明**: 同一个用户仅有一套账号密码，登录后根据角色区分权限。
* **请求体**:
  ```json
  {
    "username": "hr_01",
    "password": "password123"
  }
  ```
* **成功响应 `data`**:
  ```json
  {
    "token": "eyJhbG...", // JWT 令牌
    "role": "hr",         // hr 或 candidate
    "user_id": 1
  }
  ```

---

## 2. HR 管理端专用接口 (需校验角色为 HR)

### 2.1 发布新岗位
* **路径**: `/hr/jobs`
* **方法**: `POST`
* **请求体**:
  ```json
  {
    "title": "Go后端开发工程师",
    "description": "熟悉 Gin 和 gRPC，有微服务开发经验..."
  }
  ```

### 2.2 编辑或下架岗位
* **路径**: `/hr/jobs/:id` (URL 参数携带岗位 ID)
* **方法**: `PUT`
* **请求体**:
  ```json
  {
    "title": "Go/前端全栈开发",
    "description": "修改后的描述...",
    "status": 0  // 0-下架, 1-上架
  }
  ```

### 2.3 分页查看个人发布的岗位
* **路径**: `/hr/jobs?page=1&size=10`
* **方法**: `GET`
* **成功响应 `data`**:
  ```json
  {
    "jobs": [
      {
        "id": 1, 
        "title": "Go后端开发", 
        "status": 1, 
        "created_at": "2026-05-11T10:00:00Z"
      }
    ],
    "total": 15
  }
  ```

### 2.4 查看岗位投递候选人详情
* **路径**: `/hr/jobs/:id/applications`
* **方法**: `GET`
* **成功响应 `data`**:
  ```json
  {
    "profiles": [
      {
        "candidate_id": 101,
        "real_name": "张三",
        "phone": "13800000000",
        "university": "武汉理工大学",
        "highest_education": "本科",
        "skills": "Go, Vue3",
        "resume_url": "https://oss-bucket... (私有签名限时URL)"
      }
    ]
  }
  ```

### 2.5 AI 智能对话 (对接 Eino 框架)
* **路径**: `/hr/ai/chat`
* **方法**: `POST`
* **说明**: 自动加载上下文，并将本次对话存入 MySQL 历史记录表。
* **请求体**:
  ```json
  {
    "message": "帮我统计一下Go后端开发岗位的投递总人数"
  }
  ```
* **成功响应 `data`**:
  ```json
  {
    "reply": "根据数据库记录，Go后端开发岗位目前共有 15 人投递。其中本科学历 10 人，硕士学历 5 人。"
  }
  ```

---

## 3. 候选人用户端接口

### 3.1 游客浏览公开岗位列表
* **路径**: `/public/jobs?page=1&size=10`
* **方法**: `GET`
* **权限**: **无任何权限门槛** (无需 JWT 即可访问)
* **说明**: 仅返回状态为 1（已上架）的岗位列表。

### 3.2 完善结构化个人档案
* **路径**: `/candidate/profile`
* **方法**: `PUT`
* **权限**: 需候选人 JWT
* **请求体**:
  ```json
  {
    "real_name": "李四",
    "phone": "13900000000",
    "highest_education": "本科",
    "university": "武汉理工大学",
    "experience": "参与AI对话系统开发...",
    "skills": "Go, React, MySQL",
    "resume_oss_url": "oss://bucket/path/resume.pdf" // 仅存 OSS 相对路径
  }
  ```

### 3.3 一键投递岗位
* **路径**: `/candidate/jobs/:id/apply`
* **方法**: `POST`
* **权限**: 需候选人 JWT 
* **说明**: 后端网关需强制校验该候选人是否已完善档案及上传规范简历，未完善则拦截并返回对应错误码。
* **成功响应**: 
  ```json
  {
    "code": 0,
    "msg": "投递成功",
    "data": null
  }
  ```