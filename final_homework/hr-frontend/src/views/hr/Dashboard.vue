<template>
  <div class="dashboard-container">
    <!-- 数据统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon total">
          <el-icon><Document /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.totalJobs }}</div>
          <div class="stat-label">总岗位数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon active">
          <el-icon><Check /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.activeJobs }}</div>
          <div class="stat-label">招聘中</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon paused">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.pausedJobs }}</div>
          <div class="stat-label">已下架</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon applications">
          <el-icon><UserFilled /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.todayApplications }}</div>
          <div class="stat-label">今日投递</div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <el-card class="main-card" shadow="never">
      <el-tabs v-model="activeTab" class="hr-tabs">

        <!-- 岗位管理 Tab -->
        <el-tab-pane name="jobs">
          <template #label>
            <div class="tab-label">
              <el-icon><Document /></el-icon>
              <span>岗位管理</span>
              <el-badge :value="jobList.length" type="primary" class="tab-badge"/>
            </div>
          </template>

          <div class="tab-content">
            <!-- 操作栏 -->
            <div class="action-bar">
              <div class="action-left">
                <el-button type="primary" @click="dialogVisible = true">
                  <el-icon><Plus /></el-icon>
                  发布新岗位
                </el-button>
                <el-button @click="fetchJobs" :loading="loading">
                  <el-icon><Refresh /></el-icon>
                  刷新列表
                </el-button>
              </div>
              <div class="action-right">
                <el-input
                  v-model="searchKeyword"
                  placeholder="搜索岗位名称..."
                  clearable
                  style="width: 200px;"
                  @input="handleSearch"
                >
                  <template #prefix>
                    <el-icon><Search /></el-icon>
                  </template>
                </el-input>
              </div>
            </div>

            <!-- 岗位表格 -->
            <el-table
              :data="filteredJobs"
              v-loading="loading"
              class="job-table"
              stripe
            >
              <el-table-column prop="id" label="ID" width="80" align="center">
                <template #default="{ row }">
                  <el-tag type="info" effect="plain" size="small">{{ row.id }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="title" label="岗位名称" min-width="180">
                <template #default="{ row }">
                  <div class="job-title-cell">
                    <span class="job-name">{{ row.title }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="description" label="岗位要求" min-width="250" show-overflow-tooltip>
                <template #default="{ row }">
                  <div class="job-desc-cell">{{ row.description }}</div>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="100" align="center">
                <template #default="{ row }">
                  <el-tag
                    :type="row.status === 1 ? 'success' : 'info'"
                    effect="dark"
                    size="small"
                  >
                    <el-icon v-if="row.status === 1"><Check /></el-icon>
                    <el-icon v-else><Clock /></el-icon>
                    {{ row.status === 1 ? '招聘中' : '已下架' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="260" align="center" fixed="right">
                <template #default="{ row }">
                  <div class="action-buttons">
                    <el-button size="small" type="primary" plain @click="viewApplications(row.id)">
                      <el-icon><Files /></el-icon>
                      查看简历
                    </el-button>
                    <el-button size="small" @click="openEditDialog(row)">
                      <el-icon><Edit /></el-icon>
                      编辑
                    </el-button>
                    <el-button
                      size="small"
                      :type="row.status === 1 ? 'warning' : 'success'"
                      @click="toggleJobStatus(row)"
                    >
                      <el-icon v-if="row.status === 1"><Clock /></el-icon>
                      <el-icon v-else><Check /></el-icon>
                      {{ row.status === 1 ? '下架' : '上架' }}
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- AI 问答 Tab -->
        <el-tab-pane name="ai">
          <template #label>
            <div class="tab-label">
              <el-icon><ChatDotRound /></el-icon>
              <span>AI智能助手</span>
            </div>
          </template>

          <div class="chat-container">
            <!-- 聊天消息区域 -->
            <div class="chat-messages" ref="chatBox">
              <div
                v-for="(msg, index) in aiMessages"
                :key="index"
                :class="['message-item', msg.role === 'user' ? 'msg-user' : 'msg-ai']"
              >
                <div class="message-avatar">
                  <el-avatar
                    :size="40"
                    :class="msg.role === 'user' ? 'avatar-user' : 'avatar-ai'"
                  >
                    <el-icon v-if="msg.role === 'user'"><UserFilled /></el-icon>
                    <el-icon v-else><Service /></el-icon>
                  </el-avatar>
                </div>
                <div class="message-content">
                  <div class="message-role">
                    {{ msg.role === 'user' ? '我' : 'Eino AI 助手' }}
                  </div>
                  <div class="message-bubble">
                    <div class="bubble-content">{{ msg.content }}</div>
                  </div>
                </div>
              </div>

              <!-- AI正在输入 -->
              <div v-if="aiLoading" class="message-item msg-ai">
                <div class="message-avatar">
                  <el-avatar :size="40" class="avatar-ai">
                    <el-icon><Service /></el-icon>
                  </el-avatar>
                </div>
                <div class="message-content">
                  <div class="message-role">Eino AI 助手</div>
                  <div class="message-bubble typing">
                    <div class="typing-indicator">
                      <span></span>
                      <span></span>
                      <span></span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 输入区域 -->
            <div class="chat-input-area">
              <div class="input-wrapper">
                <el-input
                  v-model="aiInput"
                  type="textarea"
                  :rows="3"
                  placeholder="询问 AI 助手，例如：帮我写一份 Go 开发工程师的招聘要求..."
                  resize="none"
                  @keyup.enter.ctrl="sendMessage"
                />
                <div class="input-tips">
                  <span>Ctrl + Enter 发送</span>
                </div>
              </div>
              <el-button
                type="primary"
                class="send-btn"
                :loading="aiLoading"
                @click="sendMessage"
              >
                <el-icon><Promotion /></el-icon>
                发送
              </el-button>
            </div>
          </div>
        </el-tab-pane>

      </el-tabs>
    </el-card>

    <!-- 发布岗位对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="发布新岗位"
      width="560px"
      class="job-dialog"
    >
      <div class="dialog-content">
        <el-form :model="jobForm" label-position="top">
          <el-form-item label="岗位名称" required>
            <el-input
              v-model="jobForm.title"
              placeholder="例如：高级 Go 语言开发工程师"
              clearable
            />
          </el-form-item>
          <el-form-item label="岗位要求" required>
            <el-input
              v-model="jobForm.description"
              type="textarea"
              :rows="5"
              placeholder="请详细描述岗位要求、技术栈、工作职责..."
              resize="none"
            />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitLoading" @click="submitJob">
            <el-icon><Check /></el-icon>
            发布岗位
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 编辑岗位对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑岗位"
      width="560px"
      class="job-dialog"
    >
      <div class="dialog-content">
        <el-form :model="editForm" label-position="top">
          <el-form-item label="岗位名称" required>
            <el-input
              v-model="editForm.title"
              placeholder="例如：高级 Go 语言开发工程师"
              clearable
            />
          </el-form-item>
          <el-form-item label="岗位要求" required>
            <el-input
              v-model="editForm.description"
              type="textarea"
              :rows="5"
              placeholder="请详细描述岗位要求、技术栈、工作职责..."
              resize="none"
            />
          </el-form-item>
          <el-form-item label="岗位状态">
            <el-radio-group v-model="editForm.status">
              <el-radio :value="1">
                <el-tag type="success" effect="plain">招聘中</el-tag>
              </el-radio>
              <el-radio :value="0">
                <el-tag type="info" effect="plain">已下架</el-tag>
              </el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="editLoading" @click="updateJob">
            <el-icon><Check /></el-icon>
            保存修改
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 简历详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="岗位投递详情"
      size="50%"
      class="resume-drawer"
    >
      <div class="drawer-content">
        <div class="drawer-header">
          <el-icon><Files /></el-icon>
          <span>已投递简历列表</span>
        </div>
        <el-table :data="applicationList" v-loading="drawerLoading" stripe>
          <el-table-column prop="real_name" label="姓名" width="120" align="center">
            <template #default="{ row }">
              <el-tag effect="plain">{{ row.real_name }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="university" label="毕业院校" min-width="180" />
          <el-table-column prop="phone" label="联系电话" width="140" align="center" />
          <el-table-column prop="highest_education" label="学历" width="100" align="center">
            <template #default="{ row }">
              <el-tag type="success" effect="plain" size="small">{{ row.highest_education }}</el-tag>
            </template>
          </el-table-column>
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

  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document, Check, Clock, UserFilled, Plus, Refresh, Search,
  Files, Edit, ChatDotRound, Service, Promotion, Download
} from '@element-plus/icons-vue'
import request from '../../utils/request'

const activeTab = ref('jobs')

// 搜索
const searchKeyword = ref('')

// 统计数据
const stats = ref({
  totalJobs: 0,
  activeJobs: 0,
  pausedJobs: 0,
  todayApplications: 0
})

// 岗位管理相关变量
const dialogVisible = ref(false)
const loading = ref(false)
const submitLoading = ref(false)
const jobList = ref([])
const jobForm = reactive({ title: '', description: '' })

// 过滤后的岗位
const filteredJobs = computed(() => {
  if (!searchKeyword.value) return jobList.value
  const keyword = searchKeyword.value.toLowerCase()
  return jobList.value.filter(job =>
    job.title.toLowerCase().includes(keyword) ||
    job.description.toLowerCase().includes(keyword)
  )
})

const handleSearch = () => {}

// 编辑岗位相关变量
const editDialogVisible = ref(false)
const editLoading = ref(false)
const editForm = reactive({ id: 0, title: '', description: '', status: 1 })

// 简历抽屉相关变量
const drawerVisible = ref(false)
const drawerLoading = ref(false)
const applicationList = ref([])

// AI 聊天相关变量
const chatBox = ref(null)
const aiMessages = ref([])
const aiInput = ref('')
const aiLoading = ref(false)

// 方法区

const fetchJobs = async () => {
  loading.value = true
  try {
    const res = await request.get('/hr/jobs')
    jobList.value = res.data || []

    // 更新统计
    stats.value.totalJobs = jobList.value.length
    stats.value.activeJobs = jobList.value.filter(j => j.status === 1).length
    stats.value.pausedJobs = jobList.value.filter(j => j.status === 0).length
    stats.value.todayApplications = Math.floor(jobList.value.length * 0.5) // 模拟数据
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const submitJob = async () => {
  if (!jobForm.title || !jobForm.description) {
    ElMessage.warning('请填写完整的岗位信息')
    return
  }
  submitLoading.value = true
  try {
    await request.post('/hr/jobs', jobForm)
    ElMessage.success('岗位发布成功！')
    dialogVisible.value = false
    jobForm.title = ''
    jobForm.description = ''
    fetchJobs()
  } catch (error) {
    console.error(error)
  } finally {
    submitLoading.value = false
  }
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
  } finally {
    drawerLoading.value = false
  }
}

const openEditDialog = (job) => {
  editForm.id = job.id
  editForm.title = job.title
  editForm.description = job.description
  editForm.status = job.status
  editDialogVisible.value = true
}

const updateJob = async () => {
  if (!editForm.title || !editForm.description) {
    ElMessage.warning('请填写完整的岗位信息')
    return
  }
  editLoading.value = true
  try {
    await request.put(`/hr/jobs/${editForm.id}`, {
      title: editForm.title,
      description: editForm.description,
      status: editForm.status
    })
    ElMessage.success('岗位更新成功！')
    editDialogVisible.value = false
    fetchJobs()
  } catch (error) {
    ElMessage.error('更新失败')
    console.error(error)
  } finally {
    editLoading.value = false
  }
}

const toggleJobStatus = async (job) => {
  const newStatus = job.status === 1 ? 0 : 1
  const actionText = newStatus === 0 ? '下架' : '上架'

  try {
    await request.put(`/hr/jobs/${job.id}`, {
      title: job.title,
      description: job.description,
      status: newStatus
    })
    ElMessage.success(`${actionText}成功！`)
    fetchJobs()
  } catch (error) {
    ElMessage.error(`${actionText}失败`)
    console.error(error)
  }
}

// 加载历史对话记录
const loadChatHistory = async () => {
  try {
    const res = await request.get('/hr/ai/history')
    if (res.data && res.data.length > 0) {
      res.data.forEach(msg => {
        aiMessages.value.push({
          role: msg.role,
          content: msg.content
        })
      })
    } else {
      aiMessages.value.push({
        role: 'assistant',
        content: '你好，HR！我是 Eino AI 智能招聘助手。我可以帮你：\n\n• 起草岗位招聘要求\n• 分析简历与岗位的匹配度\n• 解答招聘相关的常见问题\n\n请问有什么可以帮助你的吗？'
      })
    }
  } catch (error) {
    aiMessages.value.push({
      role: 'assistant',
      content: '你好，HR！我是 Eino AI 智能招聘助手。请问有什么可以帮助你的吗？'
    })
    console.error('加载历史对话失败', error)
  }
}

const sendMessage = async () => {
  if (!aiInput.value.trim()) return

  const question = aiInput.value
  aiMessages.value.push({ role: 'user', content: question })
  aiInput.value = ''
  aiLoading.value = true

  // 滚动到底部
  nextTick(() => {
    if (chatBox.value) {
      chatBox.value.scrollTop = chatBox.value.scrollHeight
    }
  })

  try {
    const res = await request.post('/hr/ai/chat', { question })
    aiMessages.value.push({ role: 'assistant', content: res.data })
  } catch (error) {
    console.error(error)
    aiMessages.value.push({ role: 'assistant', content: '抱歉，我的连接好像断开了，请稍后再试。' })
  } finally {
    aiLoading.value = false
    nextTick(() => {
      if (chatBox.value) {
        chatBox.value.scrollTop = chatBox.value.scrollHeight
      }
    })
  }
}

onMounted(() => {
  fetchJobs()
  loadChatHistory()
})
</script>

<style scoped>
.dashboard-container {
  height: 100%;
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  padding: 24px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
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

.stat-icon.active {
  background: linear-gradient(135deg, #67c23a 0%, #95d475 100%);
  color: white;
}

.stat-icon.paused {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
  color: white;
}

.stat-icon.applications {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

/* 主卡片 */
.main-card {
  border-radius: 12px;
  border: none;
  min-height: calc(100vh - 200px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

/* Tab 标签 */
.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tab-badge {
  margin-left: 4px;
}

:deep(.el-tabs__item) {
  font-size: 15px;
}

.tab-content {
  padding: 20px 0;
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.action-left {
  display: flex;
  gap: 12px;
}

/* 岗位表格 */
.job-table {
  border-radius: 8px;
}

:deep(.job-table .el-table__header-wrapper) {
  background: #f5f7fa;
}

.job-title-cell {
  font-weight: 500;
  color: #303133;
}

.job-desc-cell {
  color: #606266;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

/* 聊天容器 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 500px;
  background: #f5f7fa;
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 12px;
  max-width: 80%;
}

.msg-user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.msg-ai {
  align-self: flex-start;
}

.message-avatar {
  flex-shrink: 0;
}

.avatar-user {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.avatar-ai {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.message-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.message-role {
  font-size: 12px;
  color: #909399;
}

.msg-user .message-role {
  text-align: right;
}

.message-bubble {
  padding: 12px 16px;
  border-radius: 12px;
  line-height: 1.6;
}

.msg-user .message-bubble {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-bottom-right-radius: 4px;
}

.msg-ai .message-bubble {
  background: white;
  color: #303133;
  border-bottom-left-radius: 4px;
  border: 1px solid #ebeef5;
}

/* AI正在输入动画 */
.typing-indicator {
  display: flex;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #f5576c;
  border-radius: 50%;
  animation: typing 1.4s infinite both;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.6;
  }
  30% {
    transform: translateY(-4px);
    opacity: 1;
  }
}

/* 输入区域 */
.chat-input-area {
  padding: 16px 20px;
  background: white;
  border-top: 1px solid #ebeef5;
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.input-wrapper {
  flex: 1;
}

.input-tips {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

.send-btn {
  height: 76px;
  font-weight: 600;
}

/* 对话框 */
.dialog-content {
  padding: 16px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 抽屉 */
.drawer-content {
  padding: 0 20px 20px;
}

.drawer-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 20px;
}

.no-resume {
  color: #c0c4cc;
  font-size: 12px;
}
</style>