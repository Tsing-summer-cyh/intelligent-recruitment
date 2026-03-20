const express = require("express");
const cors = require("cors");
const app = express();

// 允许跨域请求
app.use(cors());

// 定义多组事件数据，演示“同一接口承载不同功能”的最佳实践
const EVENT_TYPES = {
  未读消息: { event: "未读消息", data: "你有一条未读消息 (实践版)" },
  订单更新: { event: "订单更新", data: "订单 #10086 已发货 (实践版)" },
  系统通知: { event: "系统通知", data: "系统将于今晚 22:00 维护 (实践版)" },
};

app.get("/api/sse", (req, res) => {
  // 1. 核心步骤：设置 SSE 必需的 HTTP 响应头
  res.writeHead(200, {
    "Content-Type": "text/event-stream; charset=utf-8", // 声明这是事件流
    "Cache-Control": "no-cache",                        // 禁用缓存
    "Connection": "keep-alive",                         // 保持连接不断开
  });

  let tick = 0;
  const eventKeys = Object.keys(EVENT_TYPES);

  // 2. 模拟服务端持续不断地向客户端推送数据
  const intervalId = setInterval(() => {
    // 轮流获取事件类型
    const key = eventKeys[tick % eventKeys.length];
    const { event, data } = EVENT_TYPES[key];
    
    // 3. 核心规范：先写 event 类型，再写 data，最后用 \n\n 表示本条消息结束
    res.write(`event: ${event}\n`);
    res.write(`data: ${data}\n\n`);
    
    tick++;
  }, 2000);

  // 4. 监听客户端断开连接（如浏览器关闭或主动调用 close）
  res.on("close", () => {
    console.log("检测到客户端断开连接，清理定时器资源");
    clearInterval(intervalId);
    res.end();
  });
});

app.listen(3000, () => {
  console.log("SSE 实践服务已启动，访问地址: http://localhost:3000");
});