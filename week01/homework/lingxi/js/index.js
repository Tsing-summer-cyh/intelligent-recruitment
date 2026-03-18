// ====================
// 1. 获取所有的 DOM 元素
// ====================
const themeToggleBtn = document.getElementById('theme-toggle');
const body = document.body;
const messageInput = document.getElementById('message-input');
const sendBtn = document.getElementById('send-btn');
const welcomeArea = document.getElementById('welcome-area');
const chatArea = document.getElementById('chat-area');
const clearChatBtn = document.getElementById('clear-chat');
const suggestionCards = document.querySelectorAll('.card');
const uploadBtn = document.getElementById('upload-btn');
const fileInput = document.getElementById('file-input');
const imagePreview = document.getElementById('image-preview');

// 用于存储当前选择的图片 Base64 数据
let currentImageBase64 = null;

// ====================
// 2. 主题切换逻辑 
// ====================
// 页面加载时，检查本地存储的主题
const savedTheme = localStorage.getItem('lingxi-theme') || 'light-theme';
body.className = savedTheme;
updateThemeButtonText();

themeToggleBtn.addEventListener('click', () => {
    if (body.classList.contains('light-theme')) {
        body.className = 'dark-theme';
        localStorage.setItem('lingxi-theme', 'dark-theme');
    } else {
        body.className = 'light-theme';
        localStorage.setItem('lingxi-theme', 'light-theme');
    }
    updateThemeButtonText();
});

function updateThemeButtonText() {
    themeToggleBtn.textContent = body.classList.contains('light-theme') ? '🌓 切换深色' : '☀️ 切换浅色';
}

// ====================
// 3. 底部输入区交互 [cite: 23, 24, 25, 31]
// ====================
// 监听输入框内容变化，控制发送按钮显隐，并自适应高度
messageInput.addEventListener('input', function() {
    // 自动调整输入框高度
    this.style.height = 'auto';
    this.style.height = (this.scrollHeight) + 'px';
    
    // 控制发送按钮显示/隐藏 (有文本或有图片时显示)
    if (this.value.trim().length > 0 || currentImageBase64) {
        sendBtn.classList.remove('hidden');
    } else {
        sendBtn.classList.add('hidden');
    }
});

// 监听键盘按键 (Enter 发送，Shift+Enter 换行) 
messageInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault(); // 阻止默认的回车换行
        if (messageInput.value.trim() || currentImageBase64) {
            handleSendMessage(messageInput.value.trim());
        }
    }
});

// 点击发送按钮
sendBtn.addEventListener('click', () => {
    const text = messageInput.value.trim();
    if (text || currentImageBase64) {
        handleSendMessage(text);
    }
});

// ====================
// 4. 图片上传与预览 
// ====================
uploadBtn.addEventListener('click', () => {
    fileInput.click(); // 触发隐藏的 file input
});

fileInput.addEventListener('change', (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
        currentImageBase64 = event.target.result;
        // 显示预览图
        imagePreview.innerHTML = `
            <div style="position: relative; display: inline-block;">
                <img src="${currentImageBase64}" style="max-height: 60px; border-radius: 8px; border: 1px solid var(--border-color);">
                <button onclick="clearImage()" style="position: absolute; top: -5px; right: -5px; background: red; color: white; border: none; border-radius: 50%; width: 20px; height: 20px; cursor: pointer; font-size: 12px;">×</button>
            </div>
        `;
        imagePreview.classList.remove('hidden');
        sendBtn.classList.remove('hidden'); // 有图片也可以发送
    };
    reader.readAsDataURL(file);
});

// 清除图片的全局函数
window.clearImage = function() {
    currentImageBase64 = null;
    imagePreview.innerHTML = '';
    imagePreview.classList.add('hidden');
    fileInput.value = ''; // 清空 input
    if (messageInput.value.trim().length === 0) {
        sendBtn.classList.add('hidden');
    }
}

// ====================
// 5. 点击快捷卡片发送 [cite: 12, 13]
// ====================
suggestionCards.forEach(card => {
    card.addEventListener('click', () => {
        // 获取卡片内的纯文本内容（去除图标）
        const text = card.textContent.replace(/[\uD800-\uDBFF][\uDC00-\uDFFF]|\s/g, '').trim(); 
        // 简单处理：提取文本并发送。你可以根据实际卡片内容稍微优化这里的文本提取逻辑
        const cleanText = card.innerText.split(' ').slice(1).join(' ') || card.innerText;
        handleSendMessage(cleanText);
    });
});

// ====================
// 6. 核心：处理消息发送与 UI 状态切换
// ====================
function handleSendMessage(text) {
    // 1. 隐藏欢迎区，显示对话区 [cite: 18]
    welcomeArea.classList.add('hidden');
    chatArea.classList.remove('hidden');

    // 2. 清空输入框和图片预览
    messageInput.value = '';
    messageInput.style.height = 'auto';
    sendBtn.classList.add('hidden');
    
    let imageToSend = currentImageBase64;
    if (imageToSend) {
        clearImage();
    }

    // 3. 在界面上渲染用户的消息 [cite: 16]
    appendUserMessage(text, imageToSend);

    // 4. TODO: 这里将调用大模型 API 获取回复 (第四步的内容)
    console.log('准备发送给 AI 的文本:', text);
    if(imageToSend) console.log('包含图片数据');
    
    // 模拟一下 AI 思考状态
    fetchAIResponse(text, imageToSend); 
}

// 渲染用户消息到界面的辅助函数
function appendUserMessage(text, imageBase64) {
    const msgDiv = document.createElement('div');
    msgDiv.style.display = 'flex';
    msgDiv.style.justifyContent = 'flex-end';
    msgDiv.style.marginBottom = '20px';

    let contentHtml = '';
    if (imageBase64) {
        contentHtml += `<img src="${imageBase64}" style="max-width: 200px; border-radius: 8px; display: block; margin-bottom: 8px;">`;
    }
    if (text) {
        contentHtml += `<div>${text}</div>`;
    }

    msgDiv.innerHTML = `
        <div style="background-color: #e3f2fd; color: #000; padding: 12px 16px; border-radius: 12px 0 12px 12px; max-width: 80%;">
            ${contentHtml}
        </div>
    `;
    chatArea.appendChild(msgDiv);
    // 滚动到底部
    chatArea.scrollTop = chatArea.scrollHeight;
}

// ====================
// 7. 清除对话功能 [cite: 28]
// ====================
clearChatBtn.addEventListener('click', () => {
    chatArea.innerHTML = '';
    chatArea.classList.add('hidden');
    welcomeArea.classList.remove('hidden');
    chatHistory = []; // 新增这行：清空对话历史数组
});

// ====================
// 8. Markdown 与代码高亮配置
// ====================
marked.setOptions({
    highlight: function(code, lang) {
        const language = hljs.getLanguage(lang) ? lang : 'plaintext';
        return hljs.highlight(code, { language }).value;
    },
    langPrefix: 'hljs language-' // 配合 highlight.js 的 CSS 类名
});

// ====================
// 9. 核心：大模型 API 接入与流式响应
// ====================
const LINGXI_API_KEY = '         请填写你的 Lingxi API Key      '; 

let currentAbortController = null; // 用于控制“停止生成”
let chatHistory = []; // 维护上下文对话历史

async function fetchAIResponse(text, imageBase64) {
    // 1. 创建 AI 回复的 UI 气泡
    const aiMsgDiv = document.createElement('div');
    aiMsgDiv.style.display = 'flex';
    aiMsgDiv.style.justifyContent = 'flex-start';
    aiMsgDiv.style.marginBottom = '20px';
    
    // 气泡内部结构：头像 + 内容区 + 停止按钮
    aiMsgDiv.innerHTML = `
        <div style="margin-right: 15px; font-size: 28px;">🤖</div>
        <div style="background-color: var(--card-bg); border: 1px solid var(--border-color); padding: 12px 16px; border-radius: 0 12px 12px 12px; max-width: 80%; width: 100%; position: relative;">
            <div class="markdown-body" style="line-height: 1.6; overflow-x: auto;">思考中...</div>
            <button class="stop-btn" style="margin-top: 10px; background: #ff5252; color: white; border: none; padding: 5px 10px; border-radius: 4px; cursor: pointer; font-size: 12px;">⏹️ 停止生成</button>
        </div>
    `;
    chatArea.appendChild(aiMsgDiv);
    chatArea.scrollTop = chatArea.scrollHeight;

    const contentDiv = aiMsgDiv.querySelector('.markdown-body');
    const stopBtn = aiMsgDiv.querySelector('.stop-btn');
    
    // 2. 绑定“停止生成”事件
    currentAbortController = new AbortController();
    stopBtn.addEventListener('click', () => {
        if (currentAbortController) {
            currentAbortController.abort(); // 中断 fetch 请求
            stopBtn.style.display = 'none';
        }
    });

    // 3. 构建请求体 (兼容图文和纯文本)
    // 使用阿里云百炼提供的 OpenAI 兼容接口，方便解析标准的 SSE 流
    const apiUrl = 'https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions';
    
    let messageContent;
    let modelName = 'qwen-plus'; // 默认纯文本模型
    
    if (imageBase64) {
        modelName = 'qwen-vl-plus'; // 如果有图片，必须切换到视觉模型
        messageContent = [
            { type: 'image_url', image_url: { url: imageBase64 } },
            { type: 'text', text: text || '请描述这张图片' }
        ];
    } else {
        messageContent = text;
    }

    // 将用户新消息加入历史记录，保持上下文连贯
    chatHistory.push({ role: 'user', content: messageContent });

    const requestBody = {
        model: modelName,
        messages: chatHistory,
        stream: true // 关键：开启流式输出
    };

    let fullAiResponse = '';

    try {
        const response = await fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${LINGXI_API_KEY}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(requestBody),
            signal: currentAbortController.signal // 传入 signal 以支持中断
        });

        if (!response.ok) {
            throw new Error(`HTTP 错误: ${response.status}`);
        }

        // 4. 处理 ReadableStream 流式响应
        const reader = response.body.getReader();
        const decoder = new TextDecoder('utf-8');

        while (true) {
            const { done, value } = await reader.read();
            if (done) break;

            const chunk = decoder.decode(value, { stream: true });
            const lines = chunk.split('\n');
            
            for (const line of lines) {
                // 解析 Server-Sent Events (SSE) 格式的数据
                if (line.startsWith('data: ') && line !== 'data: [DONE]') {
                    try {
                        const data = JSON.parse(line.substring(6));
                        if (data.choices && data.choices[0].delta && data.choices[0].delta.content) {
                            fullAiResponse += data.choices[0].delta.content;
                            // 实时渲染 Markdown 并更新 DOM
                            contentDiv.innerHTML = marked.parse(fullAiResponse);
                            chatArea.scrollTop = chatArea.scrollHeight;
                        }
                    } catch (e) {
                        console.error('JSON 解析跳过:', line);
                    }
                }
            }
        }
        
        // 生成正常结束
        stopBtn.style.display = 'none';
        chatHistory.push({ role: 'assistant', content: fullAiResponse }); // 保存 AI 回复到历史
        addCopyButtons(contentDiv); // 渲染完毕后，为代码块添加复制按钮

    } catch (error) {
        if (error.name === 'AbortError') {
            contentDiv.innerHTML += '<br><br><em style="color: #ff5252;">[用户已手动停止生成]</em>';
        } else {
            contentDiv.innerHTML = `<span style="color: red;">请求出错: ${error.message}</span>`;
        }
        stopBtn.style.display = 'none';
    } finally {
        currentAbortController = null;
    }
}

// ====================
// 10. 附加功能：代码块一键复制
// ====================
function addCopyButtons(container) {
    // 找到所有 Markdown 渲染出的 pre 标签（代码块）
    const preElements = container.querySelectorAll('pre');
    preElements.forEach(pre => {
        if (pre.querySelector('.copy-btn')) return; // 避免重复添加

        pre.style.position = 'relative';
        const copyBtn = document.createElement('button');
        copyBtn.className = 'copy-btn';
        copyBtn.innerText = '复制';
        // 简单的内联样式控制复制按钮
        copyBtn.style.cssText = 'position: absolute; top: 5px; right: 5px; background: rgba(128,128,128,0.2); border: 1px solid var(--border-color); border-radius: 4px; padding: 4px 8px; cursor: pointer; font-size: 12px; color: var(--text-color);';
        
        copyBtn.addEventListener('click', () => {
            const code = pre.querySelector('code').innerText;
            navigator.clipboard.writeText(code).then(() => {
                copyBtn.innerText = '已复制!';
                setTimeout(() => copyBtn.innerText = '复制', 2000); // 2秒后恢复文字
            });
        });
        pre.appendChild(copyBtn);
    });
}