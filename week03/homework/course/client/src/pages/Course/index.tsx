import React, { useEffect, useState } from 'react';
import { Card, Table, Button, Form, Input, Select, Space, Tag, Modal, Popconfirm, message, InputNumber } from 'antd';
import request from '../../utils/request';

export default function Course() {
  const [form] = Form.useForm();
  const [modalForm] = Form.useForm();
  const [data, setData] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState([]);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10 });
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);

  const fetchCourses = async (page = 1, pageSize = 10, values = {}, sorter: any = {}) => {
    setLoading(true);
    try {
      const params: any = { page, pageSize, ...values };
      // 处理服务端排序
      if (sorter.field && sorter.order) {
        params.sortField = sorter.field;
        params.sortOrder = sorter.order;
      }
      const res: any = await request.get('/api/courses', { params });
      setData(res.list);
      setTotal(res.total);
      setPagination({ current: page, pageSize });
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    const res: any = await request.get('/api/courses/categories');
    setCategories(res);
  };

  useEffect(() => {
    fetchCourses();
    fetchCategories();
  }, []);

  const handleTableChange = (newPagination: any, filters: any, sorter: any) => {
    fetchCourses(newPagination.current, newPagination.pageSize, form.getFieldsValue(), sorter);
  };

  const handleSearch = (values: any) => {
    fetchCourses(1, pagination.pageSize, values);
  };

  const handleDelete = async (id: number) => {
    await request.delete(`/api/courses/${id}`);
    message.success('删除成功');
    fetchCourses(pagination.current, pagination.pageSize, form.getFieldsValue());
  };

  const handleToggleStatus = async (id: number) => {
    await request.patch(`/api/courses/${id}/status`);
    message.success('状态更新成功');
    fetchCourses(pagination.current, pagination.pageSize, form.getFieldsValue());
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
      await request.put(`/api/courses/${editingId}`, values);
      message.success('更新成功');
    } else {
      await request.post('/api/courses', values);
      message.success('新增成功');
    }
    setIsModalOpen(false);
    fetchCourses(pagination.current, pagination.pageSize, form.getFieldsValue());
  };

  const columns = [
    {
      title: '课程名称', dataIndex: 'name', width: '25%', render: (txt: string, r: any) => (
        <div><div className="font-bold">{txt}</div><div className="text-gray-400 text-xs truncate max-w-[200px]">{r.description}</div></div>
      )
    },
    { title: '讲师', dataIndex: 'instructor' },
    { title: '分类', dataIndex: 'category', render: (txt: string) => <Tag color="blue">{txt}</Tag> },
    { title: '课时', dataIndex: 'lesson_count' },
    { title: '选课人数', dataIndex: 'student_count', sorter: true }, // 开启服务端排序
    { title: '状态', dataIndex: 'status', render: (s: string) => <Tag color={s === 'published' ? 'success' : 'default'}>{s === 'published' ? '已发布' : '草稿'}</Tag> },
    {
      title: '操作',
      render: (_: any, r: any) => (
        <Space>
          <Button type="link" size="small" onClick={() => openModal(r)}>编辑</Button>
          <Popconfirm title={`确定要${r.status === 'published' ? '下架' : '发布'}该课程吗？`} onConfirm={() => handleToggleStatus(r.id)}>
            <Button type="link" size="small">{r.status === 'published' ? '下架' : '发布'}</Button>
          </Popconfirm>
          <Popconfirm title="确定删除该课程吗？" onConfirm={() => handleDelete(r.id)}>
            <Button type="link" danger size="small">删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <Card title="课程管理" extra={<Button type="primary" onClick={() => openModal()}>+ 新增课程</Button>}>
      <Form form={form} layout="inline" onFinish={handleSearch} className="mb-4">
        <Form.Item name="keyword"><Input placeholder="搜索课程名/讲师" allowClear /></Form.Item>
        <Form.Item name="status">
          <Select placeholder="全部状态" allowClear style={{ width: 120 }}>
            <Select.Option value="published">已发布</Select.Option>
            <Select.Option value="draft">草稿</Select.Option>
          </Select>
        </Form.Item>
        <Form.Item name="category">
          <Select placeholder="全部分类" allowClear style={{ width: 120 }}>
            {categories.map(c => <Select.Option key={c as string} value={c}>{c}</Select.Option>)}
          </Select>
        </Form.Item>
        <Form.Item><Button type="primary" htmlType="submit">搜索</Button></Form.Item>
      </Form>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{ ...pagination, total, showSizeChanger: true }}
        onChange={handleTableChange}
      />

      <Modal title={editingId ? '编辑课程' : '新增课程'} open={isModalOpen} onOk={handleModalOk} onCancel={() => setIsModalOpen(false)}>
        <Form form={modalForm} layout="vertical">
          <Form.Item name="name" label="课程名称" rules={[{ required: true }]}><Input placeholder="请输入课程名称" /></Form.Item>
          <Form.Item name="description" label="课程描述"><Input.TextArea placeholder="请输入课程描述" /></Form.Item>
          <div className="flex gap-4">
            <Form.Item name="instructor" label="讲师" className="flex-1"><Input placeholder="请输入讲师" /></Form.Item>
            <Form.Item name="category" label="分类" className="flex-1">
              <Select placeholder="请选择分类" allowClear>
                {categories.map((c: any) => (
                  <Select.Option key={c} value={c}>{c}</Select.Option>
                ))}
              </Select>
            </Form.Item>
          </div>
          <div className="flex gap-4">
            {/* 删掉了 initialValue={0}，改用 placeholder="0" 显示灰色占位符 */}
            <Form.Item name="lesson_count" label="课时数" className="flex-1">
              <InputNumber min={0} className="w-full" placeholder="0" />
            </Form.Item>
            
            <Form.Item name="status" label="状态" className="flex-1">
              <Select placeholder="请选择状态" allowClear>
                <Select.Option value="draft">草稿</Select.Option>
                <Select.Option value="published">已发布</Select.Option>
              </Select>
            </Form.Item>
          </div>
        </Form>
      </Modal>
    </Card>
  );
}