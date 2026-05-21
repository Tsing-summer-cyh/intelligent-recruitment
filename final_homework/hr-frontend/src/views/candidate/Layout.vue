<template>
  <el-container class="candidate-layout">
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
            <span>智能招聘 · 候选人管理</span>
          </div>
          <div class="breadcrumb">
            <el-breadcrumb separator="/">
              <el-breadcrumb-item>HR工作台</el-breadcrumb-item>
              <el-breadcrumb-item>{{ currentPageTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
        </div>
        <div class="header-right">
          <div class="header-actions">
            <el-tooltip content="系统通知" placement="bottom">
              <el-badge :value="2" class="action-badge">
                <el-button circle>
                  <el-icon><Bell /></el-icon>
                </el-button>
              </el-badge>
            </el-tooltip>
          </div>
          <div class="user-info">
            <el-avatar :size="32" class="user-avatar">
              <el-icon><UserFilled /></el-icon>
            </el-avatar>
            <div class="user-detail">
              <span class="user-name">HR 管理员</span>
              <span class="user-role">候选人管理</span>
            </div>
            <el-dropdown trigger="click" @command="handleCommand">
              <el-button type="primary" link>
                <el-icon><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="dashboard">
                    <el-icon><DataBoard /></el-icon>
                    返回工作台
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
      <el-aside width="220px" class="aside">
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          background-color="#1a1f2e"
          text-color="#a3a7b2"
          active-text-color="#ffffff"
          @select="handleMenuSelect"
        >
          <!-- 返回工作台 -->
          <el-menu-item index="/hr/dashboard">
            <el-icon><DataBoard /></el-icon>
            <span>返回工作台</span>
          </el-menu-item>

          <!-- 候选人管理 -->
          <el-sub-menu index="candidates">
            <template #title>
              <el-icon><UserFilled /></el-icon>
              <span>候选人管理</span>
            </template>
            <el-menu-item index="/hr/candidate/jobs">
              <el-icon><Files /></el-icon>
              <span>投递职位</span>
            </el-menu-item>
            <el-menu-item index="/hr/profiles">
              <el-icon><User /></el-icon>
              <span>候选人档案</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- 统计分析 -->
          <el-sub-menu index="stats">
            <template #title>
              <el-icon><TrendCharts /></el-icon>
              <span>统计分析</span>
            </template>
            <el-menu-item index="/hr/stats">
              <el-icon><DataBoard /></el-icon>
              <span>数据概览</span>
            </el-menu-item>
            <el-menu-item index="/hr/stats">
              <el-icon><TrendCharts /></el-icon>
              <span>趋势分析</span>
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
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Bell, UserFilled, ArrowDown, SwitchButton, DataBoard,
  Files, User, TrendCharts
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const activeMenu = ref(route.path)

const currentPageTitle = computed(() => {
  const titleMap = {
    '/hr/candidate/jobs': '投递职位'
  }
  return titleMap[route.path] || '候选人管理'
})

const handleMenuSelect = (index) => {
  if (index.startsWith('/')) {
    router.push(index)
  }
}

const handleCommand = (command) => {
  if (command === 'logout') {
    localStorage.clear()
    ElMessage.success('已安全退出系统')
    router.push('/login')
  } else if (command === 'dashboard') {
    router.push('/hr/dashboard')
  }
}
</script>

<style scoped>
.candidate-layout {
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