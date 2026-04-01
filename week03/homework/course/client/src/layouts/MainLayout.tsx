import React from 'react';
import { Layout, Menu, Button, theme, Dropdown } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { DashboardOutlined, BookOutlined, TeamOutlined, EditOutlined, UserOutlined, LogoutOutlined } from '@ant-design/icons';

const { Header, Sider, Content } = Layout;

export default function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { token: { colorBgContainer, borderRadiusLG } } = theme.useToken();

  // 从本地存储获取当前登录的用户信息
  const userStr = localStorage.getItem('user');
  const user = userStr ? JSON.parse(userStr) : { name: '管理员' };

  // 退出登录
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };

  // 侧边栏菜单配置
  const items = [
    { key: '/dashboard', icon: <DashboardOutlined />, label: '工作台' },
    { key: '/course', icon: <BookOutlined />, label: '课程管理' },
    { key: '/student', icon: <TeamOutlined />, label: '学生管理' },
    { key: '/summary', icon: <EditOutlined />, label: '学习总结' },
  ];

  return (
    // 1. 将 min-h-screen 改为 h-screen (100vh)，并隐藏外层滚动条
    <Layout className="h-screen overflow-hidden">
      
      {/* 2. 侧边栏本身也占满整个高度 */}
      <Sider theme="light" width={220} className="h-full border-r border-gray-100">
        <div className="h-16 flex items-center justify-center border-b border-gray-100">
          <h1 className="text-xl font-bold text-gray-800 m-0">学习管理平台</h1>
        </div>
        {/* 防止未来菜单过多，侧边栏内部也可以独立滚动 */}
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={items}
          onClick={({ key }) => navigate(key)}
          className="border-r-0 mt-2 overflow-y-auto"
          style={{ height: 'calc(100vh - 64px)' }}
        />
      </Sider>

      <Layout className="h-full">
        <Header style={{ padding: '0 24px', background: colorBgContainer }} className="flex justify-end items-center shadow-sm z-10 shrink-0">
          <Dropdown menu={{ items: [{ key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: handleLogout }] }}>
            <Button type="text" className="flex items-center gap-2">
              <UserOutlined /> {user.name}
            </Button>
          </Dropdown>
        </Header>
        
        {/* 3. 核心：只在这里开启 overflow: 'auto'，让内容独立滚动 */}
        <Content 
          style={{ 
            margin: '24px 16px', 
            padding: 24, 
            background: colorBgContainer, 
            borderRadius: borderRadiusLG, 
            overflow: 'auto' 
          }}
        >
          {/* 页面内容在此渲染 */}
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}