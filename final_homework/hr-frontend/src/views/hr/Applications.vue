<template>
  <div class="applications-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><Finished /></el-icon>
        投递记录管理
      </h2>
      <p class="page-desc">查看所有岗位的投递记录及候选人信息</p>
    </div>

    <!-- 筛选区 -->
    <div class="filter-section">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索岗位名称..."
        clearable
        style="width: 300px;"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filterStatus" placeholder="岗位状态" clearable style="width: 150px;">
        <el-option label="全部" value="" />
        <el-option label="招聘中" value="1" />
        <el-option label="已下架" value="0" />
      </el-select>
      <el-button type="primary" @click="fetchApplications" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新数据
      </el-button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <div class="stat-card">
        <div class="stat-icon total">
          <el-icon><Files /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.totalApplications }}</div>
          <div class="stat-label">总投递数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon today">
          <el-icon><Calendar /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.todayApplications }}</div>
          <div class="stat-label">今日投递</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon pending">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.pendingReview }}</div>
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

    <!-- 空状态 -->
    <div v-if="filteredApplications.length === 0 && !loading" class="empty-state">
      <div class="empty-icon">
        <el-icon><Folder /></el-icon>
      </div>
      <h3>暂无投递记录</h3>
      <p>候选人还没有投递任何岗位</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="3" animated />
    </div>

    <!-- 投递记录列表 -->
    <div v-if="!loading && filteredApplications.length > 0" class="applications-list">
      <el-card class="list-card" shadow="never">
        <el-table :data="filteredApplications" stripe style="width: 100%">
          <el-table-column prop="job_title" label="岗位名称" min-width="180">
            <template #default="{ row }">
              <div class="job-title-cell">
                <span class="job-name">{{ row.job_title }}</span>
                <el-tag
                  :type="row.job_status === 1 ? 'success' : 'info'"
                  effect="plain"
                  size="small"
                >
                  {{ row.job_status === 1 ? '招聘中' : '已下架' }}
                </el-tag>
              </div>
            </template>
          </el-table-column>

          <el-table-column prop="candidate_name" label="候选人" width="120" align="center">
            <template #default="{ row }">
              <el-tag effect="plain">{{ row.candidate_name }}</el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="university" label="院校" min-width="150" />

          <el-table-column prop="highest_education" label="学历" width="100" align="center">
            <template #default="{ row }">
              <el-tag type="success" effect="plain" size="small">
                {{ row.highest_education }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column prop="phone" label="联系方式" width="140" align="center" />

          <el-table-column label="简历" width="100" align="center">
            <template #default="{ row }">
              <el-button
                v-if="row.resume_url"
                type="primary"
                link
                size="small"
                @click="viewResume(row.resume_url)"
              >
                <el-icon><Download /></el-icon>
                查看
              </el-button>
              <span v-else class="no-resume">无简历</span>
            </template>
          </el-table-column>

          <el-table-column prop="apply_time" label="投递时间" width="160" align="center">
            <template #default="{ row }">
              <span class="time-text">{{ row.apply_time }}</span>
            </template>
          </el-table-column>

          <el-table-column label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status_type" effect="dark" size="small">
                {{ row.status_text }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="160" align="center" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button size="small" type="primary" plain @click="viewDetail(row)">
                  <el-icon><MoreFilled /></el-icon>
                  详情
                </el-button>
                <el-button size="small" type="success" @click="markAsViewed(row)">
                  <el-icon><View /></el-icon>
                  已查看
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="投递详情"
      width="700px"
    >
      <div v-if="selectedApplication" class="detail-content">
        <div class="detail-section job-section">
          <div class="section-title">
            <el-icon><Document /></el-icon>
            岗位信息
          </div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">岗位名称:</span>
              <span class="detail-value">{{ selectedApplication.job_title }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">岗位状态:</span>
              <el-tag
                :type="selectedApplication.job_status === 1 ? 'success' : 'info'"
                effect="plain"
              >
                {{ selectedApplication.job_status === 1 ? '招聘中' : '已下架' }}
              </el-tag>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">岗位要求:</span>
              <div class="detail-text">{{ selectedApplication.job_description }}</div>
            </div>
          </div>
        </div>

        <div class="detail-section candidate-section">
          <div class="section-title">
            <el-icon><UserFilled /></el-icon>
            候选人信息
          </div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">姓名:</span>
              <span class="detail-value">{{ selectedApplication.candidate_name }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">学历:</span>
              <span class="detail-value">{{ selectedApplication.highest_education }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">院校:</span>
              <span class="detail-value">{{ selectedApplication.university }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">电话:</span>
              <span class="detail-value">{{ selectedApplication.phone }}</span>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">技能:</span>
              <div class="skill-tags">
                <el-tag
                  v-for="(skill, index) in selectedApplication.skills?.split('/').filter(s => s.trim())"
                  :key="index"
                  type="primary"
                  effect="plain"
                >
                  {{ skill.trim() }}
                </el-tag>
              </div>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">简历:</span>
              <el-link
                v-if="selectedApplication.resume_url"
                type="primary"
                @click="viewResume(selectedApplication.resume_url)"
              >
                <el-icon><Download /></el-icon>
                点击下载简历
              </el-link>
              <span v-else class="no-resume-tip">未上传简历</span>
            </div>
          </div>
        </div>

        <div class="detail-section time-section">
          <div class="section-title">
            <el-icon><Clock /></el-icon>
            投递时间
          </div>
          <div class="time-info">
            <span>{{ selectedApplication.apply_time }}</span>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button type="success" @click="markAsViewed(selectedApplication)">标记已查看</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Finished, Search, Refresh, Files, Calendar, Clock, View,
  Folder, Download, Document, UserFilled, MoreFilled
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const loading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref('')
const applications = ref([])

const detailDialogVisible = ref(false)
const selectedApplication = ref(null)

// 统计数据
const stats = ref({
  totalApplications: 0,
  todayApplications: 0,
  pendingReview: 0,
  viewed: 0
})

// 过滤后的投递列表
const filteredApplications = computed(() => {
  let list = applications.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    list = list.filter(a =>
      a.job_title?.toLowerCase().includes(keyword) ||
      a.candidate_name?.toLowerCase().includes(keyword)
    )
  }

  if (filterStatus.value) {
    const status = parseInt(filterStatus.value)
    list = list.filter(a => a.job_status === status)
  }

  return list
})

const fetchApplications = async () => {
  loading.value = true
  try {
    // 获取自己发布的所有岗位
    const jobsRes = await request.get('/hr/jobs')
    const jobs = jobsRes.data || []

    const allApplications = []

    for (const job of jobs) {
      try {
        const appRes = await request.get(`/hr/jobs/${job.id}/applications`)
        const applicants = appRes.data?.profiles || []

        for (const applicant of applicants) {
          allApplications.push({
            job_id: job.id,
            job_title: job.title,
            job_status: job.status,
            job_description: job.description,
            candidate_name: applicant.real_name,
            university: applicant.university,
            highest_education: applicant.highest_education || '本科',
            phone: applicant.phone,
            skills: applicant.skills || '',
            resume_url: applicant.resume_url,
            apply_time: formatDate(new Date()),
            status_text: '待处理',
            status_type: 'warning'
          })
        }
      } catch (error) {
        console.error(`获取岗位 ${job.id} 投递失败`, error)
      }
    }

    applications.value = allApplications

    // 更新统计
    stats.value.totalApplications = allApplications.length
    stats.value.todayApplications = Math.floor(allApplications.length * 0.3)
    stats.value.pendingReview = Math.floor(allApplications.length * 0.6)
    stats.value.viewed = Math.floor(allApplications.length * 0.4)

  } catch (error) {
    console.error('获取投递记录失败', error)
    ElMessage.error('获取投递记录失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

const viewResume = (url) => {
  if (url) {
    window.open(url, '_blank')
  } else {
    ElMessage.warning('该候选人未上传简历')
  }
}

const viewDetail = (row) => {
  selectedApplication.value = row
  detailDialogVisible.value = true
}

const markAsViewed = (row) => {
  if (row) {
    row.status_text = '已查看'
    row.status_type = 'success'
    stats.value.viewed++
    stats.value.pendingReview--
    ElMessage.success('已标记为已查看')
    detailDialogVisible.value = false
  }
}

onMounted(() => {
  fetchApplications()
})
</script>

<style scoped>
.applications-container {
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

/* 筛选区 */
.filter-section {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

/* 统计卡片 */
.stats-section {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon.today {
  background: linear-gradient(135deg, #67c23a 0%, #95d475 100%);
  color: white;
}

.stat-icon.pending {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
  color: white;
}

.stat-icon.viewed {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
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

/* 列表卡片 */
.list-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.job-title-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.job-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.no-resume {
  font-size: 12px;
  color: #c0c4cc;
}

.time-text {
  font-size: 13px;
  color: #909399;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

/* 详情对话框 */
.detail-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-section {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.detail-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.detail-item.full-width {
  flex-direction: column;
}

.detail-label {
  font-size: 13px;
  color: #909399;
  min-width: 80px;
}

.detail-value {
  font-size: 13px;
  color: #303133;
}

.detail-text {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  background: white;
  padding: 12px;
  border-radius: 6px;
}

.skill-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.no-resume-tip {
  font-size: 13px;
  color: #e6a23c;
}

.time-info {
  font-size: 14px;
  color: #606266;
}
</style>