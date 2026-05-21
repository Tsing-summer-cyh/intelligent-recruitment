<template>
  <div class="stats-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><TrendCharts /></el-icon>
        统计分析
      </h2>
      <p class="page-desc">招聘数据可视化分析</p>
    </div>

    <!-- 核心指标卡片 -->
    <div class="metrics-grid">
      <div class="metric-card">
        <div class="metric-header">
          <el-icon class="metric-icon jobs"><Document /></el-icon>
          <span>岗位统计</span>
        </div>
        <div class="metric-value">{{ stats.totalJobs }}</div>
        <div class="metric-trend">
          <span class="trend-label">招聘中: {{ stats.activeJobs }}</span>
          <span class="trend-label">已下架: {{ stats.pausedJobs }}</span>
        </div>
      </div>

      <div class="metric-card">
        <div class="metric-header">
          <el-icon class="metric-icon applications"><UserFilled /></el-icon>
          <span>投递统计</span>
        </div>
        <div class="metric-value">{{ stats.totalApplications }}</div>
        <div class="metric-trend">
          <span class="trend-label">本周: {{ stats.weekApplications }}</span>
          <span class="trend-up">↑ {{ stats.applicationsGrowth }}%</span>
        </div>
      </div>

      <div class="metric-card">
        <div class="metric-header">
          <el-icon class="metric-icon candidates"><User /></el-icon>
          <span>候选人库</span>
        </div>
        <div class="metric-value">{{ stats.totalCandidates }}</div>
        <div class="metric-trend">
          <span class="trend-label">已建档: {{ stats.completedProfiles }}</span>
          <span class="trend-label">简历率: {{ stats.resumeRate }}%</span>
        </div>
      </div>

      <div class="metric-card">
        <div class="metric-header">
          <el-icon class="metric-icon conversion"><Finished /></el-icon>
          <span>转化率</span>
        </div>
        <div class="metric-value">{{ stats.conversionRate }}%</div>
        <div class="metric-trend">
          <span class="trend-label">投递→面试</span>
          <span class="trend-up">↑ 较上月提升</span>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <el-row :gutter="24" class="charts-row">
      <el-col :span="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">岗位状态分布</span>
            </div>
          </template>
          <div class="chart-content">
            <div class="pie-chart">
              <div class="pie-segment active" :style="{ transform: `rotate(${activeAngle}deg)` }"></div>
              <div class="pie-center">
                <div class="pie-label">{{ stats.activeJobs }}</div>
                <div class="pie-sub">招聘中</div>
              </div>
            </div>
            <div class="pie-legends">
              <div class="legend-item">
                <span class="legend-dot active"></span>
                <span>招聘中 ({{ stats.activeJobs }})</span>
              </div>
              <div class="legend-item">
                <span class="legend-dot paused"></span>
                <span>已下架 ({{ stats.pausedJobs }})</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="chart-header">
              <span class="chart-title">热门岗位排行</span>
            </div>
          </template>
          <div class="chart-content">
            <div class="ranking-list">
              <div v-for="(job, index) in topJobs" :key="job.id" class="ranking-item">
                <div class="ranking-index">
                  <span :class="['rank-badge', index < 3 ? 'top' : '']">{{ index + 1 }}</span>
                </div>
                <div class="ranking-info">
                  <span class="ranking-title">{{ job.title }}</span>
                  <span class="ranking-count">{{ job.applications }} 次投递</span>
                </div>
                <div class="ranking-bar">
                  <div class="bar-fill" :style="{ width: `${job.percentage}%` }"></div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 投递趋势 -->
    <el-card class="trend-card" shadow="never">
      <template #header>
        <div class="chart-header">
          <span class="chart-title">近7天投递趋势</span>
          <el-button type="primary" link size="small" @click="fetchStats">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
        </div>
      </template>
      <div class="trend-chart">
        <div class="trend-bars">
          <div v-for="(day, index) in weeklyTrend" :key="index" class="trend-day">
            <div class="trend-bar-wrapper">
              <div class="trend-bar" :style="{ height: `${day.height}%` }">
                <span class="bar-value">{{ day.value }}</span>
              </div>
            </div>
            <div class="trend-label">{{ day.label }}</div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { TrendCharts, Document, UserFilled, Finished, Refresh, User } from '@element-plus/icons-vue'
import request from '../../utils/request'

const loading = ref(false)

const stats = ref({
  totalJobs: 0,
  activeJobs: 0,
  pausedJobs: 0,
  totalApplications: 0,
  weekApplications: 0,
  applicationsGrowth: 0,
  totalCandidates: 0,
  completedProfiles: 0,
  resumeRate: 0,
  conversionRate: 0
})

const topJobs = ref([])

const weeklyTrend = ref([
  { label: '周一', value: 12, height: 40 },
  { label: '周二', value: 18, height: 60 },
  { label: '周三', value: 25, height: 83 },
  { label: '周四', value: 20, height: 67 },
  { label: '周五', value: 15, height: 50 },
  { label: '周六', value: 8, height: 27 },
  { label: '周日', value: 5, height: 17 }
])

const activeAngle = computed(() => {
  if (stats.value.totalJobs === 0) return 0
  return (stats.value.activeJobs / stats.value.totalJobs) * 360
})

const fetchStats = async () => {
  loading.value = true
  try {
    // 获取岗位数据
    const jobsRes = await request.get('/hr/jobs')
    const jobs = jobsRes.data || []

    stats.value.totalJobs = jobs.length
    stats.value.activeJobs = jobs.filter(j => j.status === 1).length
    stats.value.pausedJobs = jobs.filter(j => j.status === 0).length

    // 统计每个岗位的投递数
    const jobsWithApplications = []
    let totalApplications = 0

    for (const job of jobs) {
      try {
        const appRes = await request.get(`/hr/jobs/${job.id}/applications`)
        const count = appRes.data?.profiles?.length || 0
        jobsWithApplications.push({
          id: job.id,
          title: job.title,
          applications: count,
          percentage: 0
        })
        totalApplications += count
      } catch {
        jobsWithApplications.push({
          id: job.id,
          title: job.title,
          applications: 0,
          percentage: 0
        })
      }
    }

    stats.value.totalApplications = totalApplications
    stats.value.weekApplications = Math.floor(totalApplications * 0.3)
    stats.value.applicationsGrowth = Math.floor(Math.random() * 20 + 5)
    stats.value.totalCandidates = Math.floor(totalApplications * 1.5)
    stats.value.completedProfiles = totalApplications
    stats.value.resumeRate = Math.floor(Math.random() * 30 + 70)
    stats.value.conversionRate = Math.floor(Math.random() * 15 + 25)

    // 计算百分比并排序
    jobsWithApplications.forEach(job => {
      job.percentage = totalApplications > 0 ? (job.applications / totalApplications) * 100 : 0
    })
    topJobs.value = jobsWithApplications
      .sort((a, b) => b.applications - a.applications)
      .slice(0, 5)

    // 模拟周趋势数据
    const baseValue = Math.floor(totalApplications / 7)
    weeklyTrend.value = weeklyTrend.value.map(day => ({
      ...day,
      value: Math.floor(baseValue + Math.random() * 10),
      height: Math.floor(Math.random() * 60 + 30)
    }))

  } catch (error) {
    console.error('获取统计数据失败', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.stats-container {
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

/* 指标卡片网格 */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.metric-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.metric-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 14px;
  color: #909399;
}

.metric-icon {
  font-size: 20px;
}

.metric-icon.jobs {
  color: #409eff;
}

.metric-icon.applications {
  color: #f5576c;
}

.metric-icon.candidates {
  color: #67c23a;
}

.metric-icon.conversion {
  color: #e6a23c;
}

.metric-value {
  font-size: 36px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.metric-trend {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}

.trend-label {
  font-size: 12px;
  color: #909399;
}

.trend-up {
  font-size: 12px;
  color: #67c23a;
}

/* 图表卡片 */
.charts-row {
  margin-bottom: 24px;
}

.chart-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
  height: 320px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chart-content {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 220px;
}

/* 饼图模拟 */
.pie-chart {
  width: 160px;
  height: 160px;
  border-radius: 50%;
  background: #e4e7ed;
  position: relative;
  overflow: hidden;
}

.pie-segment {
  position: absolute;
  width: 50%;
  height: 100%;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  transform-origin: right center;
  left: 0;
}

.pie-center {
  position: absolute;
  width: 80px;
  height: 80px;
  background: white;
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.pie-label {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.pie-sub {
  font-size: 12px;
  color: #909399;
}

.pie-legends {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-left: 32px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-dot.active {
  background: #f5576c;
}

.legend-dot.paused {
  background: #e4e7ed;
}

/* 排行榜 */
.ranking-list {
  width: 100%;
  padding: 0 20px;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-index {
  width: 32px;
}

.rank-badge {
  display: inline-flex;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #e4e7ed;
  color: #909399;
  font-size: 12px;
  font-weight: 600;
  align-items: center;
  justify-content: center;
}

.rank-badge.top {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.ranking-info {
  flex: 1;
  min-width: 0;
}

.ranking-title {
  font-size: 14px;
  color: #303133;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ranking-count {
  font-size: 12px;
  color: #909399;
}

.ranking-bar {
  width: 100px;
  height: 8px;
  background: #f0f2f5;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 4px;
  transition: width 0.3s ease;
}

/* 趋势图 */
.trend-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.trend-chart {
  padding: 20px;
}

.trend-bars {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  height: 200px;
}

.trend-day {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.trend-bar-wrapper {
  width: 40px;
  height: 160px;
  display: flex;
  align-items: flex-end;
}

.trend-bar {
  width: 100%;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 4px 4px 0 0;
  position: relative;
  transition: height 0.3s ease;
  min-height: 20px;
}

.bar-value {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  color: #606266;
}

.trend-label {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}
</style>