<template>
  <div class="ai-history-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><ChatDotRound /></el-icon>
        AI对话历史
      </h2>
      <p class="page-desc">查看与Eino AI助手的历史对话记录</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="stats-info">
        <div class="stat-item">
          <el-icon><ChatLineRound /></el-icon>
          <span>共 <strong>{{ historyList.length }}</strong> 条对话</span>
        </div>
        <div class="stat-item">
          <el-icon><Calendar /></el-icon>
          <span>最近对话: <strong>{{ latestTime }}</strong></span>
        </div>
      </div>
      <div class="action-buttons">
        <el-button type="primary" @click="fetchHistory" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新记录
        </el-button>
        <el-button type="danger" @click="clearHistory" :disabled="historyList.length === 0">
          <el-icon><Delete /></el-icon>
          清空历史
        </el-button>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="historyList.length === 0 && !loading" class="empty-state">
      <div class="empty-icon">
        <el-icon><ChatDotRound /></el-icon>
      </div>
      <h3>暂无对话历史</h3>
      <p>前往工作台与AI助手开始对话吧</p>
      <el-button type="primary" @click="$router.push('/hr/dashboard')">
        <el-icon><Promotion /></el-icon>
        去对话
      </el-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="3" animated />
      <el-skeleton :rows="3" animated style="margin-top: 20px;" />
    </div>

    <!-- 对话历史列表 -->
    <div v-if="!loading && historyList.length > 0" class="history-list">
      <el-card
        v-for="(item, index) in historyList"
        :key="index"
        class="history-card"
        shadow="hover"
      >
        <!-- 用户提问 -->
        <div class="message-section user-section">
          <div class="message-header">
            <el-avatar :size="32" class="avatar-user">
              <el-icon><UserFilled /></el-icon>
            </el-avatar>
            <div class="message-meta">
              <span class="role-label">我的提问</span>
              <span class="time-label">{{ formatTime(item.created_at) }}</span>
            </div>
          </div>
          <div class="message-content user-message">{{ item.question }}</div>
        </div>

        <!-- AI回答 -->
        <div class="message-section ai-section">
          <div class="message-header">
            <el-avatar :size="32" class="avatar-ai">
              <el-icon><Service /></el-icon>
            </el-avatar>
            <div class="message-meta">
              <span class="role-label">Eino AI 助手</span>
            </div>
          </div>
          <div class="message-content ai-message">{{ item.answer }}</div>
        </div>

        <!-- 操作按钮 -->
        <div class="card-actions">
          <el-button type="primary" link size="small" @click="copyAnswer(item.answer)">
            <el-icon><CopyDocument /></el-icon>
            复制回答
          </el-button>
          <el-button type="danger" link size="small" @click="deleteItem(index)">
            <el-icon><Delete /></el-icon>
            删除此条
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 分页 -->
    <div v-if="!loading && historyList.length > pageSize" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="totalItems"
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ChatDotRound, ChatLineRound, Calendar, Refresh, Delete,
  UserFilled, Service, Promotion, CopyDocument
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const loading = ref(false)
const historyList = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)

const latestTime = computed(() => {
  if (historyList.value.length === 0) return '无'
  return formatTime(historyList.value[0].created_at)
})

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

const fetchHistory = async () => {
  loading.value = true
  try {
    const res = await request.get('/hr/ai/history')
    const messages = res.data || []

    // 将消息列表转换为问答对
    const pairs = []
    for (let i = 0; i < messages.length; i += 2) {
      if (messages[i] && messages[i + 1]) {
        pairs.push({
          question: messages[i].content,
          answer: messages[i + 1].content,
          created_at: messages[i].created_at
        })
      }
    }

    historyList.value = pairs
    totalItems.value = pairs.length

    if (pairs.length === 0) {
      ElMessage.info('暂无对话历史记录')
    }
  } catch (error) {
    console.error('获取历史失败', error)
    ElMessage.error('获取对话历史失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
}

const copyAnswer = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

const deleteItem = (index) => {
  ElMessageBox.confirm('确定删除这条对话记录吗？', '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    historyList.value.splice(index, 1)
    ElMessage.success('已删除')
  }).catch(() => {})
}

const clearHistory = () => {
  ElMessageBox.confirm('确定清空所有对话历史吗？此操作不可恢复！', '清空确认', {
    confirmButtonText: '确定清空',
    cancelButtonText: '取消',
    type: 'danger'
  }).then(() => {
    historyList.value = []
    ElMessage.success('已清空所有历史记录')
  }).catch(() => {})
}

onMounted(() => {
  fetchHistory()
})
</script>

<style scoped>
.ai-history-container {
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

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 16px 24px;
  border-radius: 12px;
  margin-bottom: 24px;
}

.stats-info {
  display: flex;
  gap: 24px;
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

.action-buttons {
  display: flex;
  gap: 12px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
}

.empty-icon {
  font-size: 64px;
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
  margin: 0 0 24px 0;
}

/* 对话历史列表 */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.history-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.message-section {
  padding: 16px;
}

.user-section {
  border-bottom: 1px solid #ebeef5;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.avatar-user {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.avatar-ai {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.role-label {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.time-label {
  font-size: 12px;
  color: #909399;
}

.message-content {
  font-size: 14px;
  line-height: 1.8;
  padding: 12px 16px;
  border-radius: 8px;
}

.user-message {
  background: #f5f7fa;
  color: #606266;
}

.ai-message {
  background: linear-gradient(135deg, #f0f9ff 0%, #e8f5e9 100%);
  color: #303133;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 0 0 12px 12px;
}

/* 分页 */
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 24px;
  background: white;
  border-radius: 12px;
  margin-top: 24px;
}
</style>