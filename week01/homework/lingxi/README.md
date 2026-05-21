# 灵犀 AI 对话助手 - 全栈训练营第一次作业

---

##  开发任务索引
本项目已完全脱离现代前端框架，使用原生 HTML + CSS + JavaScript 完成了以下核心任务：

- [x] **基础 UI 布局**：包含欢迎页、渐变色标题、快捷建议卡片及底部输入区。
- [x] **主题切换功能**：支持深色/浅色模式无缝切换，并通过 `localStorage` 实现刷新后状态保持。
- [x] **交互与多媒体**：
  - [x] 输入框自适应高度，按键监听（`Enter` 发送，`Shift+Enter` 换行）。
  - [x] 支持本地图片上传、Base64 预览及取消上传。
- [x] **大模型 API 接入**：
  - [x] 成功接入阿里云百炼 API。
  - [x] 支持纯文本对话（`qwen-plus`）与图文多模态对话（`qwen-vl-plus`）的动态模型切换。
- [x] **流式响应与渲染**：
  - [x] 使用 Fetch API 处理流式输出 (SSE)，实现打字机效果。
  - [x] 支持中断生成功能。
- [x] **高级渲染特性**：
  - [x] 实时解析并渲染 Markdown 内容（标题、表格、列表等）。
  - [x] 代码块语法高亮，并原生实现右上角“一键复制”功能。

---

##  核心技术实现思路

1. **流式响应与打字机效果**：
   使用原生 `fetch` API 并配置 `stream: true`。利用 `response.body.getReader()` 和 `TextDecoder` 逐块读取阿里云服务器返回的 Server-Sent Events (SSE) 数据流。将提取到的增量文本不断拼接到已有字符串中，再实时渲染到页面，从而实现平滑的逐字打印效果。

2. **多模态图文对话支持**：
   通过隐藏的 `<input type="file">` 标签唤起本地文件选择，使用 `FileReader` 将图片转换为 Base64 格式并渲染预览图。在发起网络请求前，动态判断是否存在图片数据：若有，则将模型名称切换为 `qwen-vl-plus`，并按官方接口规范构造包含 `image_url` 的 messages 数组。

3. **Markdown 渲染与代码高亮复制**：
   引入了 `marked.js` 和 `highlight.js`。在接收流式数据的 `while` 循环中，实时调用 `marked.parse()` 将拼接好的完整 Markdown 文本转为 HTML。为了实现代码复制，在单次生成结束后，通过 DOM 操作遍历所有生成的 `<pre>` 标签，动态插入绝对定位的 `<button>`，并结合 `navigator.clipboard.writeText()` API 实现一键复制至剪贴板。

4. **深浅色主题切换机制**：
   利用 CSS 变量（`--bg-color`, `--text-color` 等）在 `:root` 和 `body.dark-theme` 下定义两套色彩体系。通过 JS 监听切换按钮点击事件，动态切换 `document.body` 的 `class` 属性，并同步将当前主题标识存入 `localStorage`，确保用户体验连贯。

---

##  其他说明
- **API Key 安全说明**：为遵守作业规范并保障账户安全，提交的代码中已将真实的阿里云 API Key 移除，统一替换为了 `LINGXI_API_KEY` 假名。老师在本地使用 Live Server 测试时，如需发送真实请求，请在 `js/index.js` 中填入有效的 API Key。
- **项目启动**：请在项目根目录使用 VS Code 的 Live Server 插件打开 `index.html`，以确保所有相对路径资源正常加载。