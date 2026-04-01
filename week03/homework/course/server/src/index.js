import Koa from 'koa';
import Router from '@koa/router';
import cors from '@koa/cors';
import bodyParser from 'koa-bodyparser';
// --- 新增下面三个 import 用于静态托管 ---
import serve from 'koa-static';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import fs from 'fs';

import { initDatabase } from './database/init.js';
import authRoutes from './routes/auth.js';
import dashboardRoutes from './routes/dashboard.js';
import courseRoutes from './routes/courses.js';
import studentRoutes from './routes/students.js';
import summaryRoutes from './routes/summary.js';
import staticRoutes from './routes/static.js';

// --- 新增：获取当前文件的绝对路径 ---
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const app = new Koa();
const router = new Router();

initDatabase();

app.use(cors({ credentials: true }));
app.use(bodyParser());

app.use(async (ctx, next) => {
  try {
    await next();
  } catch (err) {
    const status = err.status || 500;
    ctx.status = status;
    ctx.body = { code: status, msg: err.message || '服务器内部错误', data: null };
    console.error(`[${new Date().toISOString()}] ${err.message}`);
  }
});

router.use('/api/auth', authRoutes.routes());
router.use('/api/dashboard', dashboardRoutes.routes());
router.use('/api/courses', courseRoutes.routes());
router.use('/api/students', studentRoutes.routes());
router.use('/api/summary', summaryRoutes.routes());
router.use('/api/static', staticRoutes.routes());

app.use(router.routes());
app.use(router.allowedMethods());

// ==========================================
// --- 新增核心逻辑：静态资源托管与单页应用兜底 ---
// ==========================================
const distPath = join(__dirname, '../../client/dist');

// 1. 使用 koa-static 托管前端 build 出来的 dist 目录
app.use(serve(distPath));

// 2. 前端路由兜底：因为 React Router 使用了 History 模式
// 当用户直接刷新页面（例如访问 /dashboard），Koa 找不到对应的静态文件会报错
// 我们需要把所有非 /api 的请求，统统指向 index.html 交给前端路由去解析
app.use(async (ctx, next) => {
  if (!ctx.path.startsWith('/api')) {
    const indexPath = join(distPath, 'index.html');
    if (fs.existsSync(indexPath)) {
      ctx.type = 'html';
      ctx.body = fs.readFileSync(indexPath);
      return;
    }
  }
  await next();
});
// ==========================================

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`服务端已启动: http://localhost:${PORT}`);
});