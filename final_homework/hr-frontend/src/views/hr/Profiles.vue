<template>
  <div class="profiles-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><User /></el-icon>
        候选人档案库
      </h2>
      <p class="page-desc">查看所有候选人的结构化档案信息</p>
    </div>

    <!-- 搜索筛选 -->
    <div class="filter-section">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索候选人姓名、院校、技能..."
        clearable
        style="width: 300px;"
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filterEducation" placeholder="学历筛选" clearable style="width: 150px;">
        <el-option label="全部" value="" />
        <el-option label="博士" value="博士" />
        <el-option label="硕士" value="硕士" />
        <el-option label="本科" value="本科" />
        <el-option label="大专" value="大专" />
      </el-select>
      <el-button type="primary" @click="fetchProfiles" :loading="loading">
        <el-icon><Refresh /></el-icon>
        刷新数据
      </el-button>
    </div>

    <!-- 统计信息 -->
    <div class="stats-bar">
      <div class="stat-item">
        <el-icon><UserFilled /></el-icon>
        <span>总档案数: <strong>{{ profiles.length }}</strong></span>
      </div>
      <div class="stat-item">
        <el-icon><Document /></el-icon>
        <span>简历完整率: <strong>{{ resumeRate }}%</strong></span>
      </div>
      <div class="stat-item">
        <el-icon><Finished /></el-icon>
        <span>已筛选: <strong>{{ filteredProfiles.length }}</strong></span>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="filteredProfiles.length === 0 && !loading" class="empty-state">
      <div class="empty-icon">
        <el-icon><Folder /></el-icon>
      </div>
      <h3>暂无候选人档案</h3>
      <p>还没有候选人完善档案信息</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="3" animated />
      <el-skeleton :rows="3" animated style="margin-top: 20px;" />
    </div>

    <!-- 档案卡片网格 -->
    <div v-if="!loading && filteredProfiles.length > 0" class="profiles-grid">
      <el-card
        v-for="profile in filteredProfiles"
        :key="profile.user_id"
        class="profile-card"
        shadow="hover"
      >
        <!-- 头部信息 -->
        <div class="profile-header">
          <div class="avatar-wrapper">
            <el-avatar :size="48" class="profile-avatar">
              <el-icon><UserFilled /></el-icon>
            </el-avatar>
          </div>
          <div class="header-info">
            <div class="name-row">
              <span class="profile-name">{{ profile.real_name }}</span>
              <el-tag type="success" effect="plain" size="small">
                {{ profile.highest_education }}
              </el-tag>
            </div>
            <div class="university-row">
              <el-icon><School /></el-icon>
              <span>{{ profile.university }}</span>
            </div>
          </div>
        </div>

        <!-- 基本信息 -->
        <div class="profile-body">
          <div class="info-row">
            <div class="info-item">
              <el-icon><Phone /></el-icon>
              <span class="info-label">联系电话:</span>
              <span class="info-value">{{ profile.phone }}</span>
            </div>
          </div>

          <div class="skills-section">
            <span class="skills-label">核心技能:</span>
            <div class="skill-tags">
              <el-tag
                v-for="(skill, index) in profile.skills?.split('/').filter(s => s.trim()).slice(0, 4)"
                :key="index"
                size="small"
                type="primary"
                effect="plain"
              >
                {{ skill.trim() }}
              </el-tag>
            </div>
          </div>

          <div class="experience-section">
            <span class="experience-label">项目经历:</span>
            <div class="experience-text">{{ profile.experience || '暂无' }}</div>
          </div>
        </div>

        <!-- 底部操作 -->
        <div class="profile-footer">
          <div class="resume-status">
            <el-icon v-if="profile.resume_url"><Check /></el-icon>
            <el-icon v-else><Warning /></el-icon>
            <span :class="['status-text', profile.resume_url ? 'has-resume' : 'no-resume']">
              {{ profile.resume_url ? '已上传简历' : '未上传简历' }}
            </span>
          </div>
          <div class="action-buttons">
            <el-button
              v-if="profile.resume_url"
              type="primary"
              size="small"
              @click="viewResume(profile.resume_url)"
            >
              <el-icon><View /></el-icon>
              查看简历
            </el-button>
            <el-button size="small" @click="viewDetail(profile)">
              <el-icon><MoreFilled /></el-icon>
              详情
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="候选人档案详情"
      width="600px"
    >
      <div v-if="selectedProfile" class="detail-content">
        <div class="detail-section">
          <div class="section-title">
            <el-icon><User /></el-icon>
            基本信息
          </div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">姓名:</span>
              <span class="detail-value">{{ selectedProfile.real_name }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">联系电话:</span>
              <span class="detail-value">{{ selectedProfile.phone }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <div class="section-title">
            <el-icon><School /></el-icon>
            教育背景
          </div>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">学历:</span>
              <span class="detail-value">{{ selectedProfile.highest_education }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">院校:</span>
              <span class="detail-value">{{ selectedProfile.university }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <div class="section-title">
            <el-icon><Star /></el-icon>
            专业技能
          </div>
          <div class="skill-tags-detail">
            <el-tag
              v-for="(skill, index) in selectedProfile.skills?.split('/').filter(s => s.trim())"
              :key="index"
              type="primary"
              effect="plain"
            >
              {{ skill.trim() }}
            </el-tag>
          </div>
        </div>

        <div class="detail-section">
          <div class="section-title">
            <el-icon><Document /></el-icon>
            项目经历
          </div>
          <div class="experience-detail">{{ selectedProfile.experience || '暂无填写' }}</div>
        </div>

        <div class="detail-section">
          <div class="section-title">
            <el-icon><Files /></el-icon>
            简历附件
          </div>
          <div v-if="selectedProfile.resume_url" class="resume-link">
            <el-link type="primary" @click="viewResume(selectedProfile.resume_url)">
              <el-icon><Download /></el-icon>
              点击下载简历
            </el-link>
          </div>
          <div v-else class="no-resume-tip">未上传简历附件</div>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  User, Search, Refresh, UserFilled, Document, Finished, Folder,
  School, Phone, Check, Warning, View, MoreFilled,
  Star, Files, Download
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const loading = ref(false)
const searchKeyword = ref('')
const filterEducation = ref('')
const profiles = ref([])

const detailDialogVisible = ref(false)
const selectedProfile = ref(null)

// 过滤后的档案列表
const filteredProfiles = computed(() => {
  let list = profiles.value

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    list = list.filter(p =>
      p.real_name?.toLowerCase().includes(keyword) ||
      p.university?.toLowerCase().includes(keyword) ||
      p.skills?.toLowerCase().includes(keyword)
    )
  }

  if (filterEducation.value) {
    list = list.filter(p => p.highest_education === filterEducation.value)
  }

  return list
})

// 简历完整率
const resumeRate = computed(() => {
  if (profiles.value.length === 0) return 0
  const hasResume = profiles.value.filter(p => p.resume_url).length
  return Math.round((hasResume / profiles.value.length) * 100)
})

const handleSearch = () => {}

const fetchProfiles = async () => {
  loading.value = true
  try {
    // 先获取自己发布的所有岗位
    const jobsRes = await request.get('/hr/jobs')
    const jobs = jobsRes.data || []

    // 收集所有投递该岗位的候选人档案
    const allProfiles = []
    const userIdSet = new Set()

    for (const job of jobs) {
      try {
        const appRes = await request.get(`/hr/jobs/${job.id}/applications`)
        const applicants = appRes.data?.profiles || []

        for (const applicant of applicants) {
          if (!userIdSet.has(applicant.real_name)) {
            userIdSet.add(applicant.real_name)
            allProfiles.push({
              user_id: applicant.real_name,
              real_name: applicant.real_name,
              phone: applicant.phone,
              university: applicant.university,
              highest_education: applicant.highest_education || '本科',
              skills: applicant.skills || 'Go/Vue',
              experience: applicant.experience || '',
              resume_url: applicant.resume_url
            })
          }
        }
      } catch (error) {
        console.error(`获取岗位 ${job.id} 投递失败`, error)
      }
    }

    profiles.value = allProfiles

    if (allProfiles.length === 0) {
      ElMessage.info('暂无候选人档案数据')
    }
  } catch (error) {
    console.error('获取档案失败', error)
    ElMessage.error('获取候选人档案失败')
  } finally {
    loading.value = false
  }
}

const viewResume = (url) => {
  if (url) {
    window.open(url, '_blank')
  } else {
    ElMessage.warning('该候选人未上传简历')
  }
}

const viewDetail = (profile) => {
  selectedProfile.value = profile
  detailDialogVisible.value = true
}

onMounted(() => {
  fetchProfiles()
})
</script>

<style scoped>
.profiles-container {
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

/* 统计栏 */
.stats-bar {
  display: flex;
  gap: 24px;
  background: white;
  padding: 16px 24px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.stat-item strong {
  color: #f5576c;
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

/* 档案网格 */
.profiles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.profile-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
}

.profile-avatar {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.header-info {
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.profile-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.university-row {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.profile-body {
  padding: 16px 0;
}

.info-row {
  margin-bottom: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.info-label {
  color: #909399;
}

.info-value {
  color: #303133;
}

.skills-section {
  margin-bottom: 12px;
}

.skills-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
  display: block;
}

.skill-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.experience-section {
  margin-bottom: 0;
}

.experience-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 4px;
  display: block;
}

.experience-text {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.profile-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #ebeef5;
}

.resume-status {
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-text {
  font-size: 13px;
}

.has-resume {
  color: #67c23a;
}

.no-resume {
  color: #e6a23c;
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
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.detail-item {
  display: flex;
  gap: 8px;
}

.detail-label {
  font-size: 13px;
  color: #909399;
}

.detail-value {
  font-size: 13px;
  color: #303133;
}

.skill-tags-detail {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.experience-detail {
  font-size: 14px;
  color: #606266;
  line-height: 1.8;
  background: white;
  padding: 12px;
  border-radius: 6px;
}

.resume-link {
  display: flex;
  align-items: center;
  gap: 8px;
}

.no-resume-tip {
  font-size: 13px;
  color: #e6a23c;
}
</style>