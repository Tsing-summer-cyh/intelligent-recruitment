<template>
  <div class="candidate-layout">
    <!-- 顶部导航 -->
    <el-header class="nav-header">
      <div class="header-content">
        <div class="logo">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2"/>
            </svg>
          </div>
          <span>智能招聘 · 职位大厅</span>
        </div>
        <div class="nav-actions">
          <template v-if="isLoggedIn">
            <el-button type="primary" plain @click="$router.push('/profile')">
              <el-icon><UserFilled /></el-icon>
              我的档案
            </el-button>
            <el-button type="danger" plain @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>
              退出
            </el-button>
          </template>
          <template v-else>
            <el-button type="primary" @click="$router.push('/login')">
              <el-icon><User /></el-icon>
              登录 / 注册
            </el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 主内容区 -->
    <el-main class="main-content">
      <!-- 统计 Banner -->
      <div class="stats-banner">
        <div class="stats-content">
          <div class="stats-item">
            <div class="stats-number">{{ stats.totalJobs }}</div>
            <div class="stats-label">在线职位</div>
          </div>
          <div class="stats-divider"></div>
          <div class="stats-item">
            <div class="stats-number">{{ stats.activeHrs }}</div>
            <div class="stats-label">活跃企业</div>
          </div>
          <div class="stats-divider"></div>
          <div class="stats-item">
            <div class="stats-number">{{ stats.todayApplications }}</div>
            <div class="stats-label">今日投递</div>
          </div>
        </div>
      </div>

      <!-- 搜索筛选区 -->
      <div class="search-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索岗位名称、关键词..."
          size="large"
          clearable
          class="search-input"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <div class="filter-tags">
          <el-tag
            v-for="tag in hotTags"
            :key="tag"
            :type="selectedTag === tag ? '' : 'info'"
            :effect="selectedTag === tag ? 'dark' : 'plain'"
            class="filter-tag"
            @click="toggleTag(tag)"
          >
            {{ tag }}
          </el-tag>
        </div>
      </div>

      <!-- 职位列表 -->
      <div class="job-list-container">
        <div class="list-header">
          <h2 class="page-title">
            <el-icon><Document /></el-icon>
            最新招聘职位
          </h2>
          <el-button type="primary" link @click="fetchJobs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>

        <!-- 空状态 -->
        <div v-if="filteredJobs.length === 0 && !loading" class="empty-state">
          <div class="empty-icon">
            <svg viewBox="0 0 24 24" fill="none">
              <path d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <h3>暂无匹配的职位</h3>
          <p>尝试更换搜索关键词或筛选条件</p>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <el-skeleton :rows="3" animated />
          <el-skeleton :rows="3" animated style="margin-top: 20px;"/>
          <el-skeleton :rows="3" animated style="margin-top: 20px;"/>
        </div>

        <!-- 职位卡片网格 -->
        <transition-group name="job-list" tag="div" class="job-grid" v-if="!loading">
          <div v-for="job in filteredJobs" :key="job.id" class="job-card-wrapper">
            <el-card class="job-card" shadow="hover">
              <div class="job-card-content">
                <!-- 卡片头部 -->
                <div class="job-header">
                  <div class="job-title-row">
                    <h3 class="job-title">{{ job.title }}</h3>
                    <el-tag type="success" effect="dark" size="small" class="status-tag">
                      <el-icon><Clock /></el-icon>
                      热招中
                    </el-tag>
                  </div>
                  <div class="job-meta">
                    <span class="meta-item">
                      <el-icon><OfficeBuilding /></el-icon>
                      HR ID: {{ job.hr_id }}
                    </span>
                  </div>
                </div>

                <!-- 岗位描述 -->
                <div class="job-desc">
                  {{ job.description }}
                </div>

                <!-- 技能标签 -->
                <div class="skill-tags">
                  <el-tag
                    v-for="(skill, index) in extractSkills(job.description)"
                    :key="index"
                    size="small"
                    type="info"
                    effect="plain"
                    class="skill-tag"
                  >
                    {{ skill }}
                  </el-tag>
                </div>

                <!-- 卡片底部 -->
                <div class="job-footer">
                  <div class="publish-time">
                    <el-icon><Calendar /></el-icon>
                    <span>刚刚发布</span>
                  </div>
                  <div class="action-buttons">
                    <el-button
                      v-if="!isLoggedIn"
                      type="warning"
                      plain
                      size="small"
                      @click="$router.push('/login')"
                    >
                      登录后投递
                    </el-button>
                    <el-button
                      v-else
                      type="primary"
                      size="small"
                      :loading="applyingId === job.id"
                      @click="checkAndApply(job.id)"
                    >
                      <el-icon><Promotion /></el-icon>
                      立即投递
                    </el-button>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </transition-group>
      </div>
    </el-main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  UserFilled, SwitchButton, User, Search, Refresh, Document,
  Clock, OfficeBuilding, Calendar, Promotion
} from '@element-plus/icons-vue'
import request from '../utils/request'

const router = useRouter()
const jobList = ref([])
const applyingId = ref(null)
const loading = ref(false)
const isLoggedIn = ref(!!localStorage.getItem('token'))

// 搜索和筛选
const searchKeyword = ref('')
const selectedTag = ref('')
const hotTags = ['Go', 'Vue', 'React', 'Java', 'Python', '前端', '后端', '全栈']

// 统计数据
const stats = ref({
  totalJobs: 0,
  activeHrs: 0,
  todayApplications: 0
})

// 过滤后的职位列表
const filteredJobs = computed(() => {
  let jobs = jobList.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    jobs = jobs.filter(job =>
      job.title.toLowerCase().includes(keyword) ||
      job.description.toLowerCase().includes(keyword)
    )
  }

  if (selectedTag.value) {
    const tag = selectedTag.value.toLowerCase()
    jobs = jobs.filter(job =>
      job.title.toLowerCase().includes(tag) ||
      job.description.toLowerCase().includes(tag)
    )
  }

  return jobs
})

const handleSearch = () => {
  // 搜索时触发重新计算
}

const toggleTag = (tag) => {
  if (selectedTag.value === tag) {
    selectedTag.value = ''
  } else {
    selectedTag.value = tag
  }
}

// 从描述中提取技能关键词
const extractSkills = (desc) => {
  const skillKeywords = ['Go', 'Vue', 'React', 'Java', 'Python', 'MySQL', 'Redis', 'Docker', 'K8s', '微服务', '前端', '后端']
  const found = skillKeywords.filter(skill => desc?.includes(skill))
  return found.slice(0, 4) //最多显示4个
}

const fetchJobs = async () => {
  loading.value = true
  try {
    const res = await request.get('/public/jobs')
    jobList.value = res.data || []

    // 更新统计
    stats.value.totalJobs = jobList.value.length
    stats.value.activeHrs = new Set(jobList.value.map(j => j.hr_id)).size
    stats.value.todayApplications = Math.floor(jobList.value.length * 0.3) // 模拟数据
  } catch (error) {
    console.error('获取岗位失败', error)
  } finally {
    loading.value = false
  }
}

const checkAndApply = async (jobId) => {
  const isCompleted = localStorage.getItem('profile_completed')

  if (!isCompleted) {
    ElMessageBox.confirm(
      '系统检测到您的结构化档案尚未完善或未上传合规简历，无法进行投递！',
      '投递拦截提示',
      {
        confirmButtonText: '去完善档案',
        cancelButtonText: '稍后再说',
        type: 'warning',
        customClass: 'apply-dialog'
      }
    ).then(() => {
      router.push('/profile')
    }).catch(() => {})
    return
  }

  applyJob(jobId)
}

const applyJob = async (jobId) => {
  applyingId.value = jobId
  try {
    await request.post(`/candidate/jobs/${jobId}/apply`)
    ElMessage.success('投递成功！简历已安全送达。')
  } catch (error) {
    console.error(error)
  } finally {
    applyingId.value = null
  }
}

const handleLogout = () => {
  localStorage.clear()
  isLoggedIn.value = false
  ElMessage.success('已切换为游客模式')
  router.push('/jobs').then(() => {
    fetchJobs()
  })
}

onMounted(() => { fetchJobs() })
</script>

<style scoped>
.candidate-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #f0f5ff 0%, #f7f9fa 100%);
}

/* 顶部导航 */
.nav-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.15);
  padding: 0;
  height: 72px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 40px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
  font-size: 20px;
  font-weight: 600;
}

.logo-icon {
  width: 32px;
  height: 32px;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.nav-actions {
  display: flex;
  gap: 12px;
}

/* 统计 Banner */
.stats-banner {
  background: white;
  margin: 24px auto;
  padding: 24px 40px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  max-width: 1200px;
}

.stats-content {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40px;
}

.stats-item {
  text-align: center;
}

.stats-number {
  font-size: 36px;
  font-weight: 700;
  color: #667eea;
  line-height: 1;
}

.stats-label {
  font-size: 14px;
  color: #909399;
  margin-top: 8px;
}

.stats-divider {
  width: 1px;
  height: 48px;
  background: #ebeef5;
}

/* 搜索筛选区 */
.search-section {
  max-width: 1200px;
  margin: 0 auto 24px;
  padding: 0 40px;
}

.search-input {
  margin-bottom: 16px;
}

:deep(.search-input .el-input__wrapper) {
  border-radius: 12px;
  padding: 8px 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.filter-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-tag {
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-tag:hover {
  transform: translateY(-2px);
}

/* 职位列表区 */
.job-list-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 40px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 20px;
  color: #303133;
  margin: 0;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 80px 20px;
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-icon svg {
  width: 100%;
  height: 100%;
}

.empty-state h3 {
  font-size: 18px;
  color: #606266;
  margin: 0 0 8px 0;
}

.empty-state p {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

/* 加载状态 */
.loading-state {
  padding: 20px;
}

/* 职位网格 */
.job-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 24px;
}

.job-card-wrapper {
  transition: all 0.3s ease;
}

.job-card {
  border-radius: 16px;
  border: 1px solid #ebeef5;
  transition: all 0.3s ease;
  overflow: hidden;
}

.job-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.15);
  border-color: #667eea;
}

.job-card-content {
  padding: 4px;
}

.job-header {
  margin-bottom: 16px;
}

.job-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.job-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  line-height: 1.4;
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

.job-meta {
  margin-top: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
}

.job-desc {
  color: #606266;
  font-size: 14px;
  line-height: 1.7;
  height: 68px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  margin-bottom: 16px;
}

.skill-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 16px;
}

.skill-tag {
  border-radius: 6px;
}

.job-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.publish-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

/* 列表动画 */
.job-list-enter-active,
.job-list-leave-active {
  transition: all 0.3s ease;
}

.job-list-enter-from,
.job-list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>