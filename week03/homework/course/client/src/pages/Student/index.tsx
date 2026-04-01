import React, { useEffect, useState } from 'react';
import { Card, Table, Button, Form, Input, Select, Space, Tag, Modal, Popconfirm, message, Checkbox, Row, Col } from 'antd';
import request from '../../utils/request';

export default function Student() {
  const [form] = Form.useForm();
  const [modalForm] = Form.useForm();
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [classes, setClasses] = useState([]);
  const [courses, setCourses] = useState([]);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10 });
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);

  const fetchStudents = async (page = 1, pageSize = 10, values = {}) => {
    setLoading(true);
    try {
      // 修复 TS 报错：加上 : any
      const res: any = await request.get('/api/students', { params: { page, pageSize, ...values } });
      setData(res.list);
      setTotal(res.total);
      setPagination({ current: page, pageSize });
    } finally {
      setLoading(false);
    }
  };

  const fetchOptions = async () => {
    // 修复 TS 报错：加上 : any
    const [clsRes, crsRes]: any = await Promise.all([
      request.get('/api/students/classes'),
      request.get('/api/courses')
    ]);
    setClasses(clsRes);
    setCourses(crsRes.list);
  };

  useEffect(() => {
    fetchStudents();
    fetchOptions();
  }, []);

  const handleSearch = (values: any) => {
    fetchStudents(1, pagination.pageSize, values);
  };

  const handleDelete = async (id: number) => {
    await request.delete(`/api/students/${id}`);
    message.success('删除成功');
    fetchStudents(pagination.current, pagination.pageSize, form.getFieldsValue());
  };

  const openModal = (record?: any) => {
    setEditingId(record ? record.id : null);
    if (record) {
      modalForm.setFieldsValue({ ...record });
    } else {
      modalForm.resetFields();
    }
    setIsModalOpen(true);
  };

  const handleModalOk = async () => {
    const values = await modalForm.validateFields();
    if (editingId) {
      await request.put(`/api/students/${editingId}`, values);
      message.success('更新成功');
    } else {
      await request.post('/api/students', values);
      message.success('创建成功');
    }
    setIsModalOpen(false);
    fetchStudents(pagination.current, pagination.pageSize, form.getFieldsValue());
  };

  const columns = [
    { title: '姓名', dataIndex: 'name', width: '10%' },
    { title: '学号', dataIndex: 'student_no', width: '12%' },
    { title: '班级', dataIndex: 'class_name', width: '12%', render: (txt: string) => <Tag color="purple">{txt}</Tag> },
    { title: '联系方式', width: '18%', render: (_: any, r: any) => <div>{r.phone}<br/><span className="text-gray-400 text-xs">{r.email}</span></div> },
    {
      title: '已选课程',
      dataIndex: 'course_ids',
      width: '22%',
      render: (courseIds: number[]) => {
        if (!courseIds || courseIds.length === 0) return <span className="text-gray-400">暂无选课</span>;
        
        // 将课程 ID 映射为课程名称
        const enrolledNames = courseIds
          .map(id => {
            const course: any = courses.find((c: any) => c.id === id);
            return course ? course.name : null;
          })
          .filter(Boolean); // 过滤掉找不到的课程
          
        // 使用 CSS 限制最大宽度并超出省略，同时增加 title 悬浮显示完整内容
        return (
          <div className="text-sm text-gray-600 truncate max-w-[200px]" title={enrolledNames.join('、')}>
            {enrolledNames.join('、')}
          </div>
        );
      }
    },
    { title: '状态', dataIndex: 'status', width: '10%', render: (status: string) => <Tag color={status === 'active' ? 'green' : 'default'}>{status === 'active' ? '活跃' : '非活跃'}</Tag> },
    {
      title: '操作',
      width: '16%',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link" size="small" onClick={() => openModal(record)}>编辑</Button>
          <Popconfirm title="确定删除该学生吗？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger size="small">删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <Card title="学生管理" extra={<Button type="primary" onClick={() => openModal()}>+ 新增学生</Button>}>
      <Form form={form} layout="inline" onFinish={handleSearch} className="mb-4">
        <Form.Item name="keyword"><Input placeholder="搜索姓名/学号" /></Form.Item>
        <Form.Item name="className">
          <Select placeholder="全部班级" allowClear style={{ width: 120 }}>
            {classes.map(c => <Select.Option key={c as string} value={c}>{c}</Select.Option>)}
          </Select>
        </Form.Item>
        <Form.Item name="status">
          <Select placeholder="全部状态" allowClear style={{ width: 120 }}>
            <Select.Option value="active">活跃</Select.Option>
            <Select.Option value="inactive">非活跃</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item><Button type="primary" htmlType="submit">搜索</Button></Form.Item>
      </Form>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{ ...pagination, total, onChange: (p, s) => fetchStudents(p, s, form.getFieldsValue()) }}
      />

      <Modal title={editingId ? '编辑学生' : '新增学生'} open={isModalOpen} onOk={handleModalOk} onCancel={() => setIsModalOpen(false)} width={600}>
        <Form form={modalForm} layout="vertical">
          {/* 修复 JSX 报错：补全了 <Input /> 和闭合标签 */}
          <Row gutter={16}>
            <Col span={12}><Form.Item name="name" label="姓名" rules={[{ required: true }]}><Input /></Form.Item></Col>
            <Col span={12}><Form.Item name="student_no" label="学号" rules={[{ required: true }]}><Input /></Form.Item></Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}><Form.Item name="class_name" label="班级"><Input /></Form.Item></Col>
            <Col span={12}>
              <Form.Item name="status" label="状态">
                <Select placeholder="请选择状态" allowClear>
                  <Select.Option value="active">活跃</Select.Option>
                  <Select.Option value="inactive">非活跃</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}><Form.Item name="phone" label="手机号"><Input /></Form.Item></Col>
            <Col span={12}><Form.Item name="email" label="邮箱"><Input /></Form.Item></Col>
          </Row>
          <Form.Item name="course_ids" label="选择课程">
            <Checkbox.Group>
              <Row>
                {courses.map((c: any) => (
                  <Col span={12} key={c.id}><Checkbox value={c.id}>{c.name}</Checkbox></Col>
                ))}
              </Row>
            </Checkbox.Group>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}