<template>
  <div class="login-container">
    <!-- 左侧品牌宣传区域 -->
    <div class="brand-section">
      <div class="brand-content">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor" opacity="0.9"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
          <h1 class="brand-name">智能招聘平台</h1>
        </div>
        <h2 class="brand-tagline">HR 专属管理后台</h2>
        <p class="brand-desc">
          集岗位发布、简历管理、AI智能问答于一体，
          让招聘工作更加高效便捷，轻松找到最佳候选人。
        </p>

        <div class="features-list">
          <div class="feature-item">
            <div class="feature-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <span>一站式岗位管理</span>
          </div>
          <div class="feature-item">
            <div class="feature-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M9 5H7C5.89543 5 5 5.89543 5 7V19C5 20.1046 5.89543 21 7 21H17C18.1046 21 19 20.1046 19 19V7C19 5.89543 18.1046 5 17 5H15M9 5C9 6.10457 9.89543 7 11 7H13C14.1046 7 15 6.10457 15 5M9 5C9 3.89543 9.89543 3 11 3H13C14.1046 3 15 3.89543 15 5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <span>简历智能筛选</span>
          </div>
          <div class="feature-item">
            <div class="feature-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M8 12H8.01M12 12H12.01M16 12H16.01M21 12C21 16.4183 16.9706 20 12 20C7.02944 20 3 16.4183 3 12C3 7.58172 7.02944 4 12 4C16.9706 4 21 7.58172 21 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <span>AI 智能问答助手</span>
          </div>
        </div>
      </div>

      <!-- 底部装饰波浪 -->
      <div class="wave-decoration">
        <svg viewBox="0 0 1440 320" preserveAspectRatio="none">
          <path fill="rgba(255,255,255,0.1)" d="M0,96L48,112C96,128,192,160,288,160C384,160,480,128,576,122.7C672,117,768,139,864,154.7C960,171,1056,181,1152,165.3C1248,149,1344,107,1392,85.3L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"/>
        </svg>
      </div>
    </div>

    <!-- 右侧登录区域 -->
    <div class="login-section">
      <div class="login-wrapper">
        <div class="login-header">
          <h2 class="login-title">HR 工作台登录</h2>
          <p class="login-subtitle">登录后台，高效管理招聘流程</p>
        </div>

        <el-card class="login-card" shadow="never">
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="rules"
            label-width="0"
            class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
                size="large"
                clearable
              />
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                class="login-btn"
                size="large"
                :loading="loading"
                @click="handleLogin"
              >
                <span v-if="!loading">登 录</span>
                <span v-else>登录中...</span>
              </el-button>
            </el-form-item>
          </el-form>

          <div class="login-footer">
            <p class="demo-tip">
              <el-icon><InfoFilled /></el-icon>
              测试账号: <code>test_hr</code> / 密码: <code>123456</code>
            </p>
          </div>
        </el-card>

        <div class="section-footer">
          <p>&copy; 2024 智能招聘平台 · HR 管理后台</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, InfoFilled } from '@element-plus/icons-vue'
import request from '../utils/request'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = reactive({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await request.post('/login', loginForm)

        ElMessage.success('登录成功！欢迎回来')

        localStorage.setItem('token', res.data.token)
        localStorage.setItem('role', res.data.role)
        localStorage.setItem('userId', res.data.user_id)

        if (res.data.role === 'hr') {
          router.push('/hr/dashboard')
        } else {
          router.push('/candidate/jobs')
        }
      } catch (error) {
        console.error('登录异常', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

/* 左侧品牌区域 */
.brand-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px;
  position: relative;
  overflow: hidden;
}

.brand-content {
  max-width: 480px;
  color: white;
  z-index: 2;
}

.logo-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  color: white;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.brand-name {
  font-size: 32px;
  font-weight: 700;
  margin: 0;
  letter-spacing: 2px;
}

.brand-tagline {
  font-size: 24px;
  font-weight: 500;
  margin: 0 0 16px 0;
  opacity: 0.95;
}

.brand-desc {
  font-size: 16px;
  line-height: 1.8;
  margin: 0 0 40px 0;
  opacity: 0.85;
}

.features-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 15px;
  opacity: 0.9;
}

.feature-icon {
  width: 24px;
  height: 24px;
  color: white;
}

.feature-icon svg {
  width: 100%;
  height: 100%;
}

.wave-decoration {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 200px;
}

.wave-decoration svg {
  width: 100%;
  height: 100%;
}

/* 右侧登录区域 */
.login-section {
  width: 420px;
  background: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px 50px;
}

.login-wrapper {
  width: 100%;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.login-subtitle {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.login-card {
  border: none;
  border-radius: 16px;
  background: transparent;
}

.login-form {
  margin-top: 0;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(240, 147, 251, 0.4);
}

.login-footer {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.demo-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #909399;
  margin: 0;
}

.demo-tip code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  color: #f5576c;
}

.section-footer {
  text-align: center;
  margin-top: 32px;
}

.section-footer p {
  font-size: 13px;
  color: #c0c4cc;
  margin: 0;
}

/* Element Plus 样式覆盖 */
:deep(.el-input__wrapper) {
  border-radius: 12px;
  padding: 4px 12px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
  transition: all 0.2s ease;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #f5576c inset, 0 0 0 3px rgba(245, 87, 108, 0.1);
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}
</style>