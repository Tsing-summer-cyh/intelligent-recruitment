<template>
  <el-container class="hr-layout">
    <!-- 顶部导航栏 -->
    <el-header class="header">
      <div class="header-content">
        <div class="header-left">
          <div class="logo">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor"/>
                <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2"/>
                <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2"/>
              </svg>
            </div>
            <span>智能招聘 · HR工作台</span>
          </div>
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item>首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ currentPageTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
        </div>
        <div class="header-right">
          <div class="header-actions">
            <!-- 系统通知按钮 -->
            <el-popover
              placement="bottom"
              :width="320"
              trigger="click"
              v-model:visible="notificationVisible"
            >
              <template #reference>
                <el-badge :value="unreadCount" class="action-badge" :hidden="unreadCount === 0">
                  <el-button circle @click="toggleNotifications">
                    <el-icon><Bell /></el-icon>
                  </el-button>
                </el-badge>
              </template>

              <!-- 通知内容 -->
              <div class="notification-panel">
                <div class="notification-header">
                  <span class="notification-title">系统通知</span>
                  <el-button type="primary" link size="small" @click="markAllAsRead">
                    全部已读
                  </el-button>
                </div>

                <div class="notification-list">
                  <div
                    v-for="(notification, index) in notifications"
                    :key="index"
                    :class="['notification-item', { unread: !notification.read }]"
                    @click="handleNotificationClick(notification)"
                  >
                    <div class="notification-icon">
                      <el-icon v-if="notification.type === 'application'"><UserFilled /></el-icon>
                      <el-icon v-else-if="notification.type === 'system'"><Bell /></el-icon>
                      <el-icon v-else><ChatDotRound /></el-icon>
                    </div>
                    <div class="notification-content">
                      <div class="notification-text">{{ notification.title }}</div>
                      <div class="notification-desc">{{ notification.desc }}</div>
                      <div class="notification-time">{{ notification.time }}</div>
                    </div>
                    <el-button
                      v-if="!notification.read"
                      type="primary"
                      link
                      size="small"
                      @click.stop="markAsRead(index)"
                    >
                      标记已读
                    </el-button>
                  </div>

                  <div v-if="notifications.length === 0" class="notification-empty">
                    <el-icon><Bell /></el-icon>
                    <span>暂无新通知</span>
                  </div>
                </div>

                <div class="notification-footer">
                  <el-button type="primary" link @click="notificationVisible = false">
                    关闭
                  </el-button>
                </div>
              </div>
            </el-popover>

            <el-tooltip content="系统设置" placement="bottom">
              <el-button circle @click="router.push('/hr/account')">
                <el-icon><Setting /></el-icon>
              </el-button>
            </el-tooltip>
          </div>
          <div class="user-info">
            <el-avatar :size="32" class="user-avatar">
              <el-icon><UserFilled /></el-icon>
            </el-avatar>
            <div class="user-detail">
              <span class="user-name">HR 管理员</span>
              <span class="user-role">招聘经理</span>
            </div>
            <el-dropdown trigger="click" @command="handleCommand">
              <el-button type="primary" link>
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="account">
                    <el-icon><User /></el-icon>
                    账户设置
                  </el-dropdown-item>
                  <el-dropdown-item command="stats">
                    <el-icon><TrendCharts /></el-icon>
                    统计分析
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </el-header>

    <el-container class="main-container">
      <!-- 侧边栏 -->
      <el-aside width="240px" class="aside">
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          background-color="#1a1f2e"
          text-color="#a3a7b2"
          active-text-color="#ffffff"
          @select="handleMenuSelect"
        >
          <!-- 工作台概览 -->
          <el-sub-menu index="dashboard">
            <template #title>
              <el-icon><DataBoard /></el-icon>
              <span>工作台概览</span>
            </template>
            <el-menu-item index="/hr/dashboard">
              <el-icon><Monitor /></el-icon>
              <span>数据看板</span>
            </el-menu-item>
            <el-menu-item index="/hr/stats">
              <el-icon><TrendCharts /></el-icon>
              <span>统计分析</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- 岗位管理 -->
          <el-sub-menu index="jobs">
            <template #title>
              <el-icon><Document /></el-icon>
              <span>岗位管理</span>
            </template>
            <el-menu-item index="/hr/dashboard">
              <el-icon><List /></el-icon>
              <span>岗位列表</span>
            </el-menu-item>
            <el-menu-item index="create-job">
              <el-icon><Plus /></el-icon>
              <span>发布新岗位</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- 候选人管理 -->
          <el-sub-menu index="candidates">
            <template #title>
              <el-icon><UserFilled /></el-icon>
              <span>候选人管理</span>
            </template>
            <el-menu-item index="/hr/profiles">
              <el-icon><User /></el-icon>
              <span>候选人档案</span>
            </el-menu-item>
            <el-menu-item index="/hr/candidate/jobs">
              <el-icon><Files /></el-icon>
              <span>投递职位</span>
            </el-menu-item>
            <el-menu-item index="/hr/applications">
              <el-icon><Finished /></el-icon>
              <span>投递记录</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- AI助手 -->
          <el-sub-menu index="ai">
            <template #title>
              <el-icon><ChatDotRound /></el-icon>
              <span>AI智能助手</span>
            </template>
            <el-menu-item index="/hr/dashboard">
              <el-icon><Service /></el-icon>
              <span>问答对话</span>
            </el-menu-item>
            <el-menu-item index="/hr/ai-history">
              <el-icon><Clock /></el-icon>
              <span>历史记录</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- 系统设置 -->
          <el-sub-menu index="settings">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </template>
            <el-menu-item index="/hr/account">
              <el-icon><User /></el-icon>
              <span>账户设置</span>
            </el-menu-item>
            <el-menu-item index="/hr/account">
              <el-icon><Lock /></el-icon>
              <span>安全设置</span>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>

      <!-- 主内容区 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Bell, Setting, UserFilled, ArrowDown, User, SwitchButton,
  DataBoard, Monitor, TrendCharts, Document, List, Plus, Files,
  Finished, ChatDotRound, Service, Clock, Lock
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const activeMenu = ref(route.path)
const notificationVisible = ref(false)

// 通知列表
const notifications = reactive([
  {
    id: 1,
    type: 'application',
    title: '新的投递通知',
    desc: '候选人张三投递了"Go开发工程师"岗位',
    time: '2分钟前',
    read: false
  },
  {
    id: 2,
    type: 'system',
    title: '系统维护通知',
    desc: '系统将于今晚22:00进行维护升级',
    time: '1小时前',
    read: false
  },
  {
    id: 3,
    type: 'ai',
    title: 'AI助手提醒',
    desc: '您的AI对话历史已保存，可随时查看',
    time: '3小时前',
    read: true
  }
])

// 未读数量
const unreadCount = computed(() => {
  return notifications.filter(n => !n.read).length
})

const currentPageTitle = computed(() => {
  const titleMap = {
    '/hr/dashboard': '数据看板',
    '/hr/stats': '统计分析',
    '/hr/profiles': '候选人档案',
    '/hr/applications': '投递记录',
    '/hr/account': '账户设置',
    '/hr/ai-history': 'AI对话历史',
    '/hr/candidate/jobs': '投递职位'
  }
  return titleMap[route.path] || '工作台'
})

const toggleNotifications = () => {
  notificationVisible.value = !notificationVisible.value
}

const handleNotificationClick = (notification) => {
  if (notification.type === 'application') {
    router.push('/hr/applications')
    notificationVisible.value = false
  } else if (notification.type === 'ai') {
    router.push('/hr/ai-history')
    notificationVisible.value = false
  }
  markAsRead(notifications.indexOf(notification))
}

const markAsRead = (index) => {
  notifications[index].read = true
}

const markAllAsRead = () => {
  notifications.forEach(n => n.read = true)
  ElMessage.success('已全部标记为已读')
}

const handleMenuSelect = (index) => {
  if (index.startsWith('/')) {
    router.push(index)
  } else if (index === 'create-job') {
    // 发布新岗位，跳转到dashboard并触发发布弹窗
    router.push('/hr/dashboard')
  }
}

const handleCommand = (command) => {
  if (command === 'logout') {
    localStorage.clear()
    ElMessage.success('已安全退出系统')
    router.push('/login')
  } else if (command === 'account') {
    router.push('/hr/account')
  } else if (command === 'stats') {
    router.push('/hr/stats')
  }
}
</script>

<style scoped>
.hr-layout {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background: #f0f2f5;
}

/* 顶部导航栏 */
.header {
  background: linear-gradient(135deg, #1a1f2e 0%, #2d3548 100%);
  height: 64px;
  padding: 0;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 100;
}

.header-content {
  width: 100%;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
}

.logo-icon {
  width: 28px;
  height: 28px;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo span {
  font-size: 18px;
  font-weight: 600;
}

.breadcrumb {
  display: flex;
  align-items: center;
}

:deep(.el-breadcrumb__inner) {
  color: #a3a7b2;
}

:deep(.el-breadcrumb__inner:last-child) {
  color: #ffffff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.header-actions .el-button {
  background: transparent;
  border: 1px solid #3a4156;
  color: #a3a7b2;
}

.header-actions .el-button:hover {
  border-color: #ffffff;
  color: #ffffff;
}

.action-badge :deep(.el-badge__content) {
  background: #f5576c;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.user-detail {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #ffffff;
}

.user-role {
  font-size: 12px;
  color: #a3a7b2;
}

.user-info .el-button {
  color: #a3a7b2;
}

/* 通知面板 */
.notification-panel {
  padding: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #ebeef5;
}

.notification-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.notification-list {
  max-height: 280px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.notification-item:hover {
  background: #f5f7fa;
}

.notification-item.unread {
  background: #f0f9ff;
}

.notification-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
}

.notification-content {
  flex: 1;
}

.notification-text {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.notification-desc {
  font-size: 12px;
  color: #606266;
  margin-top: 4px;
}

.notification-time {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.notification-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
  color: #c0c4cc;
}

.notification-footer {
  padding: 12px 16px;
  border-top: 1px solid #ebeef5;
  text-align: center;
}

/* 侧边栏 */
.main-container {
  height: calc(100vh - 64px);
}

.aside {
  background-color: #1a1f2e;
  overflow: hidden;
}

.el-menu-vertical {
  border-right: none;
  height: 100%;
}

:deep(.el-menu-vertical .el-sub-menu__title) {
  height: 56px;
  line-height: 56px;
}

:deep(.el-menu-vertical .el-sub-menu__title:hover) {
  background-color: #2d3548;
}

:deep(.el-menu-vertical .el-menu-item) {
  height: 48px;
  line-height: 48px;
}

:deep(.el-menu-vertical .el-menu-item:hover) {
  background-color: #2d3548;
}

:deep(.el-menu-vertical .el-menu-item.is-active) {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

/* 主内容区 */
.main-content {
  background-color: #f0f2f5;
  padding: 24px;
  overflow-y: auto;
}
</style>