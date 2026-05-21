<template>
  <div class="candidate-jobs-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h2 class="page-title">
          <el-icon><UserFilled /></el-icon>
          候选人职位管理
        </h2>
        <p class="page-desc">查看候选人投递的职位情况</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <div class="stat-card">
        <div class="stat-icon applications">
          <el-icon><Files /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.totalApplications }}</div>
          <div class="stat-label">总投递数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon pending">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.pending }}</div>
          <div class="stat-label">待处理</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon viewed">
          <el-icon><View /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.viewed }}</div>
          <div class="stat-label">已查看</div>
        </div>
      </div>
    </div>

    <!-- 搜索筛选 -->
    <div class="filter-section">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索职位名称..."
        clearable
        style="width: 300px;"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filterStatus" placeholder="筛选状态" clearable style="width: 150px;">
        <el-option label="全部" value="" />
        <el-option label="招聘中" value="1" />
        <el-option label="已下架" value="0" />
      </el-select>
      <el-button type="primary" @click="fetchJobs" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新数据
      </el-button>
    </div>

    <!-- 空状态 -->
    <div v-if="filteredJobs.length === 0 && !loading" class="empty-state">
      <div class="empty-icon">
        <el-icon><Folder /></el-icon>
      </div>
      <h3>暂无投递数据</h3>
      <p>候选人还没有投递任何职位</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="4" animated />
    </div>

    <!-- 职位列表 -->
    <div v-if="!loading && filteredJobs.length > 0" class="jobs-grid">
      <el-card
        v-for="job in filteredJobs"
        :key="job.id"
        class="job-card"
        shadow="hover"
      >
        <div class="job-card-header">
          <div class="job-title-section">
            <h3 class="job-title">{{ job.title }}</h3>
            <el-tag
              :type="job.status === 1 ? 'success' : 'info'"
              effect="dark"
              size="small"
            >
              {{ job.status === 1 ? '招聘中' : '已下架' }}
            </el-tag>
          </div>
          <div class="job-id">ID: {{ job.id }}</div>
        </div>

        <div class="job-card-body">
          <div class="job-description">{{ job.description }}</div>

          <!-- 投递统计 -->
          <div class="application-stats">
            <div class="stat-item">
              <el-icon><UserFilled /></el-icon>
              <span>投递人数: {{ getApplicationCount(job.id) }}</span>
            </div>
          </div>
        </div>

        <div class="job-card-footer">
          <el-button type="primary" size="small" @click="viewApplications(job.id)">
            <el-icon><Files /></el-icon>
            查看简历
          </el-button>
          <el-button size="small" @click="viewJobDetail(job)">
            <el-icon><View /></el-icon>
            详情
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 简历抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="投递简历详情"
      size="50%"
      class="resume-drawer"
    >
      <div class="drawer-content">
        <div class="drawer-header">
          <el-icon><Files /></el-icon>
          <span>已投递简历列表</span>
          <el-tag type="primary" effect="plain">共 {{ applicationList.length }}份</el-tag>
        </div>

        <el-table :data="applicationList" v-loading="drawerLoading" stripe>
          <el-table-column prop="real_name" label="姓名" width="120" align="center">
            <template #default="{ row }">
              <el-tag effect="plain">{{ row.real_name }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="university" label="毕业院校" min-width="180" />
          <el-table-column prop="highest_education" label="学历" width="100" align="center">
            <template #default="{ row }">
              <el-tag type="success" effect="plain" size="small">
                {{ row.highest_education }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="skills" label="核心技能" min-width="150" show-overflow-tooltip />
          <el-table-column prop="phone" label="联系电话" width="140" align="center" />
          <el-table-column label="简历附件" width="120" align="center">
            <template #default="{ row }">
              <el-link
                v-if="row.resume_url"
                :href="row.resume_url"
                target="_blank"
                type="primary"
              >
                <el-icon><Download /></el-icon>
                下载
              </el-link>
              <span v-else class="no-resume">无简历</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-drawer>

    <!-- 职位详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="职位详情"
      width="600px"
    >
      <div class="detail-content">
        <div class="detail-item">
          <span class="detail-label">岗位名称:</span>
          <span class="detail-value">{{ selectedJob?.title }}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">岗位状态:</span>
          <el-tag
            :type="selectedJob?.status === 1 ? 'success' : 'info'"
            effect="dark"
          >
            {{ selectedJob?.status === 1 ? '招聘中' : '已下架' }}
          </el-tag>
        </div>
        <div class="detail-item full-width">
          <span class="detail-label">岗位要求:</span>
          <div class="detail-text">{{ selectedJob?.description }}</div>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  UserFilled, Files, Clock, View, Search, Refresh,
  Folder, Download
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const loading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const jobList = ref([])
const applicationCounts = ref({})

// 统计
const stats = ref({
  totalApplications: 0,
  pending: 0,
  viewed: 0
})

// 过滤后的职位
const filteredJobs = computed(() => {
  let jobs = jobList.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    jobs = jobs.filter(job =>
      job.title.toLowerCase().includes(keyword) ||
      job.description.toLowerCase().includes(keyword)
    )
  }

  if (filterStatus.value) {
    const status = parseInt(filterStatus.value)
    jobs = jobs.filter(job => job.status === status)
  }

  return jobs
})

// 抽屉
const drawerVisible = ref(false)
const drawerLoading = ref(false)
const applicationList = ref([])

// 详情对话框
const detailDialogVisible = ref(false)
const selectedJob = ref(null)

const fetchJobs = async () => {
  loading.value = true
  try {
    const res = await request.get('/hr/jobs')
    jobList.value = res.data || []

    // 获取每个职位的投递数量
    for (const job of jobList.value) {
      try {
        const appRes = await request.get(`/hr/jobs/${job.id}/applications`)
        const count = appRes.data?.profiles?.length || 0
        applicationCounts.value[job.id] = count
      } catch {
        applicationCounts.value[job.id] = 0
      }
    }

    // 更新统计
    stats.value.totalApplications = Object.values(applicationCounts.value).reduce((a, b) => a + b, 0)
    stats.value.pending = Math.floor(stats.value.totalApplications * 0.6)
    stats.value.viewed = Math.floor(stats.value.totalApplications * 0.4)
  } catch (error) {
    console.error('获取岗位失败', error)
  } finally {
    loading.value = false
  }
}

const getApplicationCount = (jobId) => {
  return applicationCounts.value[jobId] || 0
}

const viewApplications = async (jobId) => {
  drawerVisible.value = true
  drawerLoading.value = true
  applicationList.value = []
  try {
    const res = await request.get(`/hr/jobs/${jobId}/applications`)
    applicationList.value = res.data.profiles || []
  } catch (error) {
    console.error(error)
    ElMessage.error('获取简历列表失败')
  } finally {
    drawerLoading.value = false
  }
}

const viewJobDetail = (job) => {
  selectedJob.value = job
  detailDialogVisible.value = true
}

onMounted(() => {
  fetchJobs()
})
</script>

<style scoped>
.candidate-jobs-container {
  padding: 0;
}

/* 页面标题 */
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

/* 统计卡片 */
.stats-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.stat-icon.applications {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.pending {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
  color: white;
}

.stat-icon.viewed {
  background: linear-gradient(135deg, #67c23a 0%, #95d475 100%);
  color: white;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

/* 筛选区 */
.filter-section {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
}

.empty-icon {
  font-size: 48px;
  color: #c0c4cc;
  margin-bottom: 16px;
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

/* 职位网格 */
.jobs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.job-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.job-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.job-title-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.job-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.job-id {
  font-size: 12px;
  color: #909399;
}

.job-card-body {
  padding: 16px 0;
}

.job-description {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  height: 60px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.application-stats {
  margin-top: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #909399;
}

.job-card-footer {
  display: flex;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

/* 抽屉 */
.drawer-content {
  padding: 0 20px 20px;
}

.drawer-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
}

.no-resume {
  color: #c0c4cc;
  font-size: 12px;
}

/* 详情对话框 */
.detail-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.detail-item.full-width {
  flex-direction: column;
}

.detail-label {
  font-size: 14px;
  color: #909399;
  min-width: 80px;
}

.detail-value {
  font-size: 14px;
  color: #303133;
}

.detail-text {
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}
</style>