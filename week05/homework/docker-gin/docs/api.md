# AI 智能单词本 API 接口文档

**基础说明：**
* 生产环境下，所有 API 请求通过 Nginx 代理，基础路径为 `/api` (本地测试地址为 `http://localhost:8088/api`)。
* 鉴权方式：除注册/登录外，所有接口均需要在 HTTP Header 中携带 JWT Token，格式为 `Authorization: Bearer <token>`。
* 成功状态码通常为 200，失败状态码包括 400（参数错误）、401（未授权/Token失效）、500（服务器/AI接口错误）。
---

### 1. 用户注册
* **接口路径**: `/register`
* **请求方法**: `POST`
* **鉴权说明**: 无需鉴权
* **请求参数 (Body - JSON)**:
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```
* **成功响应 (200)**:
  ```json
  { "message": "注册成功" }
  ```
* **失败说明**: `400` 用户名已存在或参数为空。

### 2. 用户登录
* **接口路径**: `/login`
* **请求方法**: `POST`
* **鉴权说明**: 无需鉴权
* **请求参数 (Body - JSON)**:
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```
* **成功响应 (200)**:
  ```json
  {
    "message": "登录成功",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
* **失败说明**: `401` 用户不存在或密码错误。

### 3. 智能查询单词
* **接口路径**: `/query`
* **请求方法**: `GET`
* **鉴权说明**: 需要 JWT
* **请求参数 (Query)**:
  * `word` (必填): 需要查询的单词，如 "apple"
  * `ai_provider` (可选): 选择的模型，如 "DeepSeek"
* **成功响应 (200)**:
  ```json
  {
    "source": "ai", 
    "data": {
      "word": "apple",
      "meaning": "苹果",
      "sentences": [
        "I ate an apple for breakfast.",
        "She picked a red apple from the tree."
      ],
      "ai_provider": "DeepSeek"
    }
  }
  ```
* **备注**: 若数据库中已存在该用户的查询记录，`source` 将返回 `"database"`，且不消耗 AI 调用次数。

### 4. 手动保存单词
* **接口路径**: `/save`
* **请求方法**: `POST`
* **鉴权说明**: 需要 JWT
* **请求参数 (Body - JSON)**:
  ```json
  {
    "word": "apple",
    "meaning": "苹果",
    "sentences": ["例句1", "例句2", "例句3"],
    "ai_provider": "DeepSeek"
  }
  ```
* **成功响应 (200)**:
  ```json
  {
    "message": "保存成功",
    "data": {
      "id": 1,
      "user_id": 1,
      "word": "apple",
      "meaning": "苹果",
      "sentences": ["例句1", "例句2", "例句3"],
      "ai_provider": "DeepSeek"
    }
  }
  ```

### 5. 获取单词本列表 (分页)
* **接口路径**: `/words`
* **请求方法**: `GET`
* **鉴权说明**: 需要 JWT
* **请求参数 (Query)**:
  * `page` (可选): 当前页码，默认 1
  * `page_size` (可选): 每页条数，默认 10
* **成功响应 (200)**:
  ```json
  {
    "total": 25,
    "page": 1,
    "page_size": 10,
    "data": [
      {
        "id": 1,
        "word": "apple",
        "meaning": "苹果",
        "sentences": ["例句1", "例句2", "例句3"],
        "ai_provider": "DeepSeek",
        "created_at": "2026-04-27T20:00:00Z"
      }
    ]
  }
  ```

### 6. 删除单词
* **接口路径**: `/words/:id`
* **请求方法**: `DELETE`
* **鉴权说明**: 需要 JWT
* **请求参数 (Path)**:
  * `id`: 单词记录的 ID
* **成功响应 (200)**:
  ```json
  { "message": "删除成功" }
  ```