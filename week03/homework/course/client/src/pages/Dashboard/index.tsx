import React, { useEffect, useState } from 'react';
import { Card, Row, Col, Statistic } from 'antd';
import { ReadOutlined, TeamOutlined, RiseOutlined, FireOutlined } from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import request from '../../utils/request';

export default function Dashboard() {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    request.get('/api/dashboard').then((res) => setData(res));
  }, []);

  if (!data) return null;

  const { stats, charts } = data;

  return (
    <div className="space-y-4">
      <Row gutter={16}>
        <Col span={6}>
          {/* 加上 className="h-full" */}
          <Card className="h-full">
            <Statistic title="课程总数" value={stats.totalCourses} prefix={<ReadOutlined />} />
            <div className="text-gray-400 text-sm mt-2">/ 已发布 {stats.publishedCourses}</div>
          </Card>
        </Col>
        <Col span={6}>
          <Card className="h-full">
            <Statistic title="学生总数" value={stats.totalStudents} prefix={<TeamOutlined />} />
            <div className="text-gray-400 text-sm mt-2">/ 活跃 {stats.activeStudents}</div>
          </Card>
        </Col>
        <Col span={6}>
          <Card className="h-full">
            <Statistic title="课程发布率" value={Math.round((stats.publishedCourses / stats.totalCourses) * 100)} suffix="%" prefix={<RiseOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card className="h-full">
            <Statistic title="学生活跃率" value={Math.round((stats.activeStudents / stats.totalStudents) * 100)} suffix="%" prefix={<FireOutlined />} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Card title="课程选课人数 TOP 8">
            <ReactECharts option={{
              tooltip: {},
              xAxis: { type: 'category', data: charts.enrollment.map((item: any) => item.name), axisLabel: { interval: 0, rotate: 30 } },
              yAxis: { type: 'value' },
              series: [{ type: 'bar', data: charts.enrollment.map((item: any) => item.value), itemStyle: { color: '#91cc75' } }]
            }} style={{ height: 300 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="近 7 天学习活跃度">
            <ReactECharts option={{
              tooltip: { trigger: 'axis' },
              legend: { data: ['学习人数', '总时长(h)'] },
              xAxis: { type: 'category', data: charts.activity.map((item: any) => item.label) },
              yAxis: { type: 'value' },
              series: [
                { name: '学习人数', type: 'line', data: charts.activity.map((item: any) => item.students) },
                { name: '总时长(h)', type: 'line', data: charts.activity.map((item: any) => item.duration) }
              ]
            }} style={{ height: 300 }} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Card title="学生状态分布">
            <ReactECharts option={{
              tooltip: { trigger: 'item' },
              series: [{ type: 'pie', radius: '60%', data: charts.statusDist }]
            }} style={{ height: 300 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="课程分类分布">
            <ReactECharts option={{
              tooltip: { trigger: 'item' },
              series: [{ type: 'pie', radius: '60%', data: charts.categoryDist }]
            }} style={{ height: 300 }} />
          </Card>
        </Col>
      </Row>
    </div>
  );
}