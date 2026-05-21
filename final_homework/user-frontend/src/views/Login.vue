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
        <h2 class="brand-tagline">让求职更智能，让招聘更高效</h2>
        <p class="brand-desc">
          基于 AI 技术的新一代招聘平台，智能匹配岗位与人才，
          为每一位求职者找到最适合的发展机会。
        </p>

        <div class="features-list">
          <div class="feature-item">
            <div class="feature-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M9 12L11 14L15 10M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <span>智能岗位推荐</span>
          </div>
          <div class="feature-item">
            <div class="feature-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12C22 6.48 17.52 2 12 2ZM12 5C15.31 5 18 7.69 18 11C18 14.31 15.31 17 12 17C8.69 17 6 14.31 6 11" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <span>简历智能优化</span>
          </div>
          <div class="feature-item">
            <div class="feature-icon">
              <svg viewBox="0 0 24 24" fill="none">
                <path d="M17 21V19C17 16.79 15.21 15 13 15H5C2.79 15 1 16.79 1 19V21M9 11C11.21 11 13 9.21 13 7C13 4.79 11.21 3 9 3C6.79 3 5 4.79 5 7C5 9.21 6.79 11 9 11ZM23 21V19C22.9999 17.1771 21.7652 15.5385 20 15V15M16 3C17.6569 4.91154 17.6569 7.08846 16 9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <span>海量优质岗位</span>
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
          <h2 class="login-title">欢迎回来</h2>
          <p class="login-subtitle">登录您的账号，开启职业新旅程</p>
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
              测试账号: <code>test_candidate</code> / 密码: <code>123456</code>
            </p>
          </div>
        </el-card>

        <div class="section-footer">
          <p>&copy; 2024 智能招聘平台 · 让梦想起飞</p>
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
          ElMessage.warning('HR 请使用后台管理系统登录！')
          localStorage.clear()
          loading.value = false
          return
        } else {
          router.push('/jobs')
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
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
  color: #409EFF;
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
  box-shadow: 0 0 0 1px #667eea inset, 0 0 0 3px rgba(102, 126, 234, 0.1);
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}
</style>