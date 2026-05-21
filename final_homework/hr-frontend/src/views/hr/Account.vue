<template>
  <div class="account-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><Setting /></el-icon>
        账户设置
      </h2>
      <p class="page-desc">管理您的账户信息和安全设置</p>
    </div>

    <!-- 用户信息卡片 -->
    <el-card class="user-card" shadow="never">
      <div class="user-info-section">
        <div class="user-avatar-wrapper">
          <el-avatar :size="80" class="user-avatar">
            <el-icon><UserFilled /></el-icon>
          </el-avatar>
          <div class="user-basic">
            <div class="user-name">{{ userInfo.name }}</div>
            <div class="user-role">
              <el-tag type="primary" effect="plain">HR 管理员</el-tag>
            </div>
          </div>
        </div>
        <div class="user-stats">
          <div class="stat-item">
            <div class="stat-value">{{ userInfo.jobsCount }}</div>
            <div class="stat-label">发布岗位</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ userInfo.applicationsCount }}</div>
            <div class="stat-label">收到投递</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ userInfo.aiChats }}</div>
            <div class="stat-label">AI对话</div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 设置选项 -->
    <div class="settings-grid">
      <!-- 个人信息设置 -->
      <el-card class="setting-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><User /></el-icon>
            <span>个人信息</span>
          </div>
        </template>
        <el-form :model="profileForm" label-position="top" class="setting-form">
          <el-form-item label="用户名">
            <el-input v-model="profileForm.username" disabled>
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
            <div class="form-tip">用户名不可修改</div>
          </el-form-item>
          <el-form-item label="显示名称">
            <el-input v-model="profileForm.displayName" placeholder="请输入显示名称">
              <template #prefix>
                <el-icon><Edit /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="联系邮箱">
            <el-input v-model="profileForm.email" placeholder="请输入邮箱地址">
              <template #prefix>
                <el-icon><Bell /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="联系电话">
            <el-input v-model="profileForm.phone" placeholder="请输入联系电话">
              <template #prefix>
                <el-icon><Phone /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveProfile" :loading="savingProfile">
              <el-icon><Check /></el-icon>
              保存信息
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 密码安全设置 -->
      <el-card class="setting-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Lock /></el-icon>
            <span>密码安全</span>
          </div>
        </template>
        <el-form :model="passwordForm" label-position="top" class="setting-form">
          <el-form-item label="当前密码">
            <el-input
              v-model="passwordForm.currentPassword"
              type="password"
              placeholder="请输入当前密码"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="新密码">
            <el-input
              v-model="passwordForm.newPassword"
              type="password"
              placeholder="请输入新密码（至少6位）"
              show-password
            >
              <template #prefix>
                <el-icon><Key /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="确认密码">
            <el-input
              v-model="passwordForm.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              show-password
            >
              <template #prefix>
                <el-icon><Key /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="warning" @click="changePassword" :loading="changingPassword">
              <el-icon><Refresh /></el-icon>
              修改密码
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 通知设置 -->
      <el-card class="setting-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Bell /></el-icon>
            <span>通知设置</span>
          </div>
        </template>
        <div class="notification-settings">
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-title">投递提醒</span>
              <span class="setting-desc">当有新候选人投递岗位时通知</span>
            </div>
            <el-switch v-model="notifySettings.applicationNotify" />
          </div>
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-title">系统公告</span>
              <span class="setting-desc">接收系统更新和维护通知</span>
            </div>
            <el-switch v-model="notifySettings.systemNotify" />
          </div>
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-title">邮件通知</span>
              <span class="setting-desc">通过邮件接收重要通知</span>
            </div>
            <el-switch v-model="notifySettings.emailNotify" />
          </div>
        </div>
      </el-card>

      <!-- 系统偏好设置 -->
      <el-card class="setting-card" shadow="never">
        <template #header>
          <div class="card-header">
            <el-icon><Tools /></el-icon>
            <span>系统偏好</span>
          </div>
        </template>
        <div class="preference-settings">
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-title">自动加载历史对话</span>
              <span class="setting-desc">进入AI助手时自动加载历史</span>
            </div>
            <el-switch v-model="preferenceSettings.autoLoadHistory" />
          </div>
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-title">表格显示密度</span>
              <span class="setting-desc">调整数据表格的显示样式</span>
            </div>
            <el-select v-model="preferenceSettings.tableDensity" style="width: 100px;">
              <el-option label="紧凑" value="compact" />
              <el-option label="默认" value="default" />
              <el-option label="宽松" value="loose" />
            </el-select>
          </div>
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-title">数据刷新频率</span>
              <span class="setting-desc">自动刷新数据的间隔时间</span>
            </div>
            <el-select v-model="preferenceSettings.refreshInterval" style="width: 100px;">
              <el-option label="关闭" value="0" />
              <el-option label="30秒" value="30" />
              <el-option label="60秒" value="60" />
            </el-select>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Setting, UserFilled, User, Edit, Bell, Phone, Lock,
  Key, Refresh, Check, Tools
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const savingProfile = ref(false)
const changingPassword = ref(false)

// 用户信息
const userInfo = ref({
  name: 'HR 管理员',
  jobsCount: 0,
  applicationsCount: 0,
  aiChats: 0
})

// 个人信息表单
const profileForm = reactive({
  username: localStorage.getItem('userId') || 'test_hr',
  displayName: 'HR 管理员',
  email: '',
  phone: ''
})

// 密码表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 通知设置
const notifySettings = reactive({
  applicationNotify: true,
  systemNotify: true,
  emailNotify: false
})

// 偏好设置
const preferenceSettings = reactive({
  autoLoadHistory: true,
  tableDensity: 'default',
  refreshInterval: '0'
})

const fetchUserStats = async () => {
  try {
    const jobsRes = await request.get('/hr/jobs')
    const jobs = jobsRes.data || []

    let applicationsCount = 0
    for (const job of jobs) {
      try {
        const appRes = await request.get(`/hr/jobs/${job.id}/applications`)
        applicationsCount += appRes.data?.profiles?.length || 0
      } catch {}
    }

    userInfo.value.jobsCount = jobs.length
    userInfo.value.applicationsCount = applicationsCount
    userInfo.value.aiChats = Math.floor(Math.random() * 20 + 5)

  } catch (error) {
    console.error('获取用户统计失败', error)
  }
}

const saveProfile = async () => {
  savingProfile.value = true
  try {
    // 模拟保存（实际需要后端接口）
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('个人信息已保存')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    savingProfile.value = false
  }
}

const changePassword = async () => {
  if (!passwordForm.currentPassword) {
    ElMessage.warning('请输入当前密码')
    return
  }
  if (!passwordForm.newPassword || passwordForm.newPassword.length < 6) {
    ElMessage.warning('新密码至少6位')
    return
  }
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.warning('两次密码输入不一致')
    return
  }

  changingPassword.value = true
  try {
    // 模拟修改（实际需要后端接口）
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('密码修改成功')
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error) {
    ElMessage.error('修改失败')
  } finally {
    changingPassword.value = false
  }
}

onMounted(() => {
  fetchUserStats()
})
</script>

<style scoped>
.account-container {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 22px;
  color: #303133;
  margin: 0 0 8px 0;
}

.page-desc {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

/* 用户信息卡片 */
.user-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
  margin-bottom: 24px;
}

.user-info-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-avatar-wrapper {
  display: flex;
  align-items: center;
  gap: 24px;
}

.user-avatar {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.user-basic {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-name {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.user-stats {
  display: flex;
  gap: 40px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #f5576c;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

/* 设置网格 */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.setting-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.setting-form {
  padding: 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

/* 通知设置 */
.notification-settings,
.preference-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.setting-desc {
  font-size: 12px;
  color: #909399;
}
</style>