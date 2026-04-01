import Router from '@koa/router';
import db from '../database/db.js';
import { authenticateToken } from '../middleware/auth.js';
import { success, fail } from '../utils/response.js';

const router = new Router();

router.get('/', authenticateToken, async (ctx) => {
  const { keyword = '', className = '', status = '', courseId = '', page = 1, pageSize = 10 } = ctx.query;
  const offset = (Number(page) - 1) * Number(pageSize);

  let where = 'WHERE 1=1';
  const params = [];

  if (keyword) {
    where += ' AND (name LIKE ? OR student_no LIKE ?)';
    params.push(`%${keyword}%`, `%${keyword}%`);
  }
  if (className) {
    where += ' AND class_name = ?';
    params.push(className);
  }
  if (status) {
    where += ' AND status = ?';
    params.push(status);
  }

  let rows = db.prepare(`SELECT * FROM students ${where} ORDER BY created_at DESC`).all(...params);

  if (courseId) {
    rows = rows.filter(s => {
      const ids = JSON.parse(s.course_ids || '[]');
      return ids.includes(Number(courseId));
    });
  }

  const total = rows.length;
  const list = rows.slice(offset, offset + Number(pageSize)).map(s => ({
    ...s,
    course_ids: JSON.parse(s.course_ids || '[]'),
  }));

  success(ctx, { list, total, page: Number(page), pageSize: Number(pageSize) });
});

router.get('/classes', authenticateToken, async (ctx) => {
  const classes = db.prepare("SELECT DISTINCT class_name FROM students WHERE class_name != '' ORDER BY class_name")
    .all()
    .map(r => r.class_name);
  success(ctx, classes);
});

router.get('/:id', authenticateToken, async (ctx) => {
  const student = db.prepare('SELECT * FROM students WHERE id = ?').get(ctx.params.id);
  if (!student) {
    return fail(ctx, 404, '学生不存在');
  }
  student.course_ids = JSON.parse(student.course_ids || '[]');

  const courses = db.prepare('SELECT id, name FROM courses').all();
  const enrolledCourses = courses.filter(c => student.course_ids.includes(c.id));

  success(ctx, { ...student, enrolledCourses });
});

// 实现创建学生接口 POST /
router.post('/', authenticateToken, async (ctx) => {
  const { name, student_no, class_name, phone, email, status, course_ids } = ctx.request.body;

  // 1. 参数校验
  if (!name || !student_no) {
    return fail(ctx, 400, '姓名和学号不能为空');
  }

  // 2. 学号唯一性检查
  const existingStudent = db.prepare('SELECT * FROM students WHERE student_no = ?').get(student_no);
  if (existingStudent) {
    return fail(ctx, 400, '学号已存在');
  }

  // 3. 插入数据库 (注意将 course_ids 数组序列化为字符串)
  const result = db.prepare(`
    INSERT INTO students (name, student_no, class_name, phone, email, status, course_ids)
    VALUES (?, ?, ?, ?, ?, ?, ?)
  `).run(
    name, 
    student_no, 
    class_name || '', 
    phone || '', 
    email || '', 
    status || 'active', 
    JSON.stringify(course_ids || [])
  );

  // 4. 更新课程的选课计数
  updateCourseCounts();

  // 5. 返回创建的学生数据，状态码设为 201
  const student = db.prepare('SELECT * FROM students WHERE id = ?').get(result.lastInsertRowid);
  student.course_ids = JSON.parse(student.course_ids || '[]');
  
  ctx.status = 201;
  success(ctx, student);
});

// 实现更新学生接口 PUT /:id
router.put('/:id', authenticateToken, async (ctx) => {
  const id = ctx.params.id;
  
  // 1. 检查学生是否存在
  const existing = db.prepare('SELECT * FROM students WHERE id = ?').get(id);
  if (!existing) {
    return fail(ctx, 404, '学生不存在');
  }

  const { name, student_no, class_name, phone, email, status, course_ids } = ctx.request.body;

  // 2. 学号唯一性校验 (确保新学号没有被其他学生占用)
  if (student_no && student_no !== existing.student_no) {
    const conflict = db.prepare('SELECT * FROM students WHERE student_no = ?').get(student_no);
    if (conflict) {
      return fail(ctx, 400, '学号已被其他学生占用');
    }
  }

  // 3. 更新数据库
  db.prepare(`
    UPDATE students 
    SET name = ?, student_no = ?, class_name = ?, phone = ?, email = ?, status = ?, course_ids = ?, updated_at = CURRENT_TIMESTAMP
    WHERE id = ?
  `).run(
    name ?? existing.name,
    student_no ?? existing.student_no,
    class_name ?? existing.class_name,
    phone ?? existing.phone,
    email ?? existing.email,
    status ?? existing.status,
    course_ids ? JSON.stringify(course_ids) : existing.course_ids,
    id
  );

  // 4. 更新课程的选课计数
  updateCourseCounts();

  // 5. 返回更新后的数据
  const student = db.prepare('SELECT * FROM students WHERE id = ?').get(id);
  student.course_ids = JSON.parse(student.course_ids || '[]');
  
  success(ctx, student);
});

// 实现删除学生接口 DELETE /:id
router.delete('/:id', authenticateToken, async (ctx) => {
  const id = ctx.params.id;
  
  // 1. 检查学生是否存在
  const existing = db.prepare('SELECT * FROM students WHERE id = ?').get(id);
  if (!existing) {
    return fail(ctx, 404, '学生不存在');
  }

  // 2. 删除数据库记录
  db.prepare('DELETE FROM students WHERE id = ?').run(id);

  // 3. 更新选课计数
  updateCourseCounts();

  success(ctx, null, '删除成功');
});
function updateCourseCounts() {
  const courses = db.prepare('SELECT id FROM courses').all();
  const students = db.prepare('SELECT course_ids FROM students').all();

  for (const course of courses) {
    const count = students.filter(s => {
      const ids = JSON.parse(s.course_ids || '[]');
      return ids.includes(course.id);
    }).length;
    db.prepare('UPDATE courses SET student_count = ? WHERE id = ?').run(count, course.id);
  }
}

export default router;
