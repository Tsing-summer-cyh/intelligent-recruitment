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
          <span>智能招聘 · 个人中心</span>
        </div>
        <div class="nav-actions">
          <el-button type="primary" plain @click="$router.push('/jobs')">
            <el-icon><Document /></el-icon>
            返回职位大厅
          </el-button>
          <el-button type="danger" plain @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
            退出
          </el-button>
        </div>
      </div>
    </el-header>

    <!-- 主内容区 -->
    <el-main class="main-content">
      <div class="profile-container">
        <!-- 步骤进度指示器 -->
        <div class="progress-section">
          <div class="progress-steps">
            <div class="step" :class="{ active: currentStep >= 1, completed: isCompleted }">
              <div class="step-icon">
                <el-icon v-if="isCompleted"><Check /></el-icon>
                <span v-else>1</span>
              </div>
              <div class="step-info">
                <div class="step-title">基本信息</div>
                <div class="step-desc">姓名、联系方式</div>
              </div>
            </div>
            <div class="step-line" :class="{ active: currentStep >= 2, completed: isCompleted }"></div>
            <div class="step" :class="{ active: currentStep >= 2, completed: isCompleted }">
              <div class="step-icon">
                <el-icon v-if="isCompleted"><Check /></el-icon>
                <span v-else>2</span>
              </div>
              <div class="step-info">
                <div class="step-title">教育背景</div>
                <div class="step-desc">学历、院校</div>
              </div>
            </div>
            <div class="step-line" :class="{ active: currentStep >= 3, completed: isCompleted }"></div>
            <div class="step" :class="{ active: currentStep >= 3, completed: isCompleted }">
              <div class="step-icon">
                <el-icon v-if="isCompleted"><Check /></el-icon>
                <span v-else>3</span>
              </div>
              <div class="step-info">
                <div class="step-title">专业技能</div>
                <div class="step-desc">技能标签、经历</div>
              </div>
            </div>
            <div class="step-line" :class="{ active: currentStep >= 4, completed: isCompleted }"></div>
            <div class="step" :class="{ active: currentStep >= 4, completed: isCompleted }">
              <div class="step-icon">
                <el-icon v-if="isCompleted"><Check /></el-icon>
                <span v-else>4</span>
              </div>
              <div class="step-info">
                <div class="step-title">简历附件</div>
                <div class="step-desc">上传简历文件</div>
              </div>
            </div>
          </div>

          <!-- 完成状态卡片 -->
          <div v-if="isCompleted" class="completed-banner">
            <div class="completed-icon">
              <el-icon><SuccessFilled /></el-icon>
            </div>
            <div class="completed-content">
              <h3>档案已完善</h3>
              <p>您的档案已达到投递标准，现在可以去投递心仪的岗位了！</p>
            </div>
            <el-button type="primary" @click="$router.push('/jobs')">
              去投递岗位
            </el-button>
          </div>
        </div>

        <!-- 表单卡片 -->
        <el-card class="profile-card" shadow="never">
          <template #header>
            <div class="card-header">
              <h2>
                <el-icon><Edit /></el-icon>
                完善个人档案
              </h2>
              <span class="sub-title">一份完善的档案能大大提高你的面试邀约率</span>
            </div>
          </template>

          <el-form ref="profileFormRef" :model="profileForm" :rules="rules" label-position="top" class="form-body">
            <!-- 基本信息 -->
            <div class="form-section">
              <div class="section-title">
                <div class="title-icon basic">
                  <el-icon><User /></el-icon>
                </div>
                <span>基本信息</span>
              </div>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="真实姓名" prop="real_name">
                    <el-input v-model="profileForm.real_name" placeholder="请输入你的真实姓名" clearable>
                      <template #prefix>
                        <el-icon><UserFilled /></el-icon>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="联系电话" prop="phone">
                    <el-input v-model="profileForm.phone" placeholder="请输入11位手机号码" maxlength="11" clearable>
                      <template #prefix>
                        <el-icon><Phone /></el-icon>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- 教育背景 -->
            <div class="form-section">
              <div class="section-title">
                <div class="title-icon education">
                  <el-icon><Reading /></el-icon>
                </div>
                <span>教育背景</span>
              </div>
              <el-row :gutter="24">
                <el-col :span="12">
                  <el-form-item label="最高学历" prop="highest_education">
                    <el-select v-model="profileForm.highest_education" placeholder="请选择学历" style="width: 100%;">
                      <el-option label="大专" value="大专" />
                      <el-option label="本科" value="本科" />
                      <el-option label="硕士" value="硕士" />
                      <el-option label="博士" value="博士" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="毕业院校" prop="university">
                    <el-input v-model="profileForm.university" placeholder="例如：武汉理工大学" clearable>
                      <template #prefix>
                        <el-icon><School /></el-icon>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>

            <!-- 专业技能 -->
            <div class="form-section">
              <div class="section-title">
                <div class="title-icon skills">
                  <el-icon><Medal /></el-icon>
                </div>
                <span>专业技能</span>
              </div>
              <el-form-item label="核心技能标签" prop="skills">
                <el-input v-model="profileForm.skills" placeholder="例如：Go / Vue3 / MySQL / 容器化部署" clearable>
                  <template #prefix>
                    <el-icon><CollectionTag /></el-icon>
                  </template>
                </el-input>
                <div class="skill-preview" v-if="profileForm.skills">
                  <el-tag
                    v-for="(skill, index) in profileForm.skills.split('/').filter(s => s.trim())"
                    :key="index"
                    type="primary"
                    effect="plain"
                    size="small"
                    class="skill-tag"
                  >
                    {{ skill.trim() }}
                  </el-tag>
                </div>
              </el-form-item>

              <el-form-item label="工作/项目经历" prop="experience">
                <el-input
                  v-model="profileForm.experience"
                  type="textarea"
                  :rows="4"
                  placeholder="请简述您的核心经历，包括项目名称、使用技术、主要贡献..."
                  resize="none"
                />
              </el-form-item>
            </div>

            <!-- 简历附件 -->
            <div class="form-section">
              <div class="section-title">
                <div class="title-icon resume">
                  <el-icon><Document /></el-icon>
                </div>
                <span>简历附件</span>
              </div>
              <el-form-item label="上传简历文件" prop="resume_url">
                <div class="upload-area">
                  <el-upload
                    class="resume-upload"
                    action="http://localhost:8080/api/v1/candidate/upload/resume"
                    :headers="{ Authorization: 'Bearer ' + token }"
                    :on-success="handleUploadSuccess"
                    :before-upload="beforeUpload"
                    :show-file-list="false"
                    accept=".pdf,.doc,.docx"
                    drag
                  >
                    <div class="upload-content">
                      <div class="upload-icon">
                        <el-icon><UploadFilled /></el-icon>
                      </div>
                      <div class="upload-text">
                        <p>点击或拖拽文件到此区域上传</p>
                        <p class="upload-tip">仅支持 PDF / DOC / DOCX 格式</p>
                      </div>
                    </div>
                  </el-upload>

                  <div v-if="profileForm.resume_url" class="resume-preview">
                    <div class="preview-icon">
                      <el-icon><Document /></el-icon>
                    </div>
                    <div class="preview-info">
                      <span class="preview-name">简历文件已上传</span>
                      <el-link type="primary" :href="profileForm.resume_url" target="_blank">
                        <el-icon><View /></el-icon>
                        点击预览/下载
                      </el-link>
                    </div>
                    <el-button type="danger" circle size="small" @click="profileForm.resume_url = ''">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                </div>
              </el-form-item>
            </div>

            <!-- 操作按钮 -->
            <div class="form-actions">
              <el-button type="primary" size="large" @click="saveProfile" :loading="saving">
                <el-icon><FolderChecked /></el-icon>
                保存档案
              </el-button>
              <el-button size="large" @click="resetForm">
                <el-icon><RefreshLeft /></el-icon>
                重置表单
              </el-button>
            </div>
          </el-form>
        </el-card>
      </div>
    </el-main>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  UserFilled, SwitchButton, Document, Check, SuccessFilled, Edit, User,
  Phone, Reading, School, Medal, CollectionTag, UploadFilled, View,
  Delete, FolderChecked, RefreshLeft
} from '@element-plus/icons-vue'
import request from '../utils/request'

const router = useRouter()
const saving = ref(false)
const token = ref(localStorage.getItem('token'))
const profileFormRef = ref(null)
const isCompleted = ref(localStorage.getItem('profile_completed') === 'true')

// 计算当前步骤进度
const currentStep = computed(() => {
  let step = 0
  if (profileForm.real_name && profileForm.phone) step = 1
  if (profileForm.highest_education && profileForm.university) step = 2
  if (profileForm.skills && profileForm.experience) step = 3
  if (profileForm.resume_url) step = 4
  return step
})

const profileForm = reactive({
  real_name: '',
  phone: '',
  highest_education: '',
  university: '',
  skills: '',
  experience: '',
  resume_url: ''
})

const rules = reactive({
  real_name: [{ required: true, message: '真实姓名不能为空', trigger: 'blur' }],
  phone: [{ required: true, message: '联系电话不能为空', trigger: 'blur' }],
  highest_education: [{ required: true, message: '请选择最高学历', trigger: 'change' }],
  university: [{ required: true, message: '毕业院校不能为空', trigger: 'blur' }],
  skills: [{ required: true, message: '核心技能标签不能为空', trigger: 'blur' }],
  experience: [{ required: true, message: '工作/项目经历不能为空', trigger: 'blur' }],
  resume_url: [{ required: true, message: '请务必上传合规的简历附件', trigger: 'change' }]
})

const fetchProfile = async () => {
  try {
    const res = await request.get('/candidate/profile')
    if (res.data) {
      Object.assign(profileForm, res.data)
    }
  } catch (error) {
    console.error('获取档案失败', error)
  }
}

const beforeUpload = (file) => {
  const isDoc = file.name.endsWith('.pdf') || file.name.endsWith('.doc') || file.name.endsWith('.docx')
  if (!isDoc) {
    ElMessage.error('只能上传 PDF/DOC/DOCX 格式的简历文件！')
  }
  return isDoc
}

const handleUploadSuccess = (res) => {
  if (res.code === 0) {
    profileForm.resume_url = res.data
    ElMessage.success('简历文件上传成功！')
    if (profileFormRef.value) {
      profileFormRef.value.validateField('resume_url').catch(() => {})
    }
  } else {
    ElMessage.error(res.msg || '上传失败')
  }
}

const saveProfile = async () => {
  if (!profileFormRef.value) return

  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        await request.put('/candidate/profile', profileForm)
        localStorage.setItem('profile_completed', 'true')
        isCompleted.value = true
        ElMessage.success('档案更新成功！HR 现在能看到你最新的信息了。')
      } catch (error) {
        console.error(error)
      } finally {
        saving.value = false
      }
    } else {
      ElMessage.warning('请先完善所有必填项！')
    }
  })
}

const resetForm = () => {
  if (profileFormRef.value) {
    profileFormRef.value.resetFields()
  }
}

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/login')
}

onMounted(() => {
  fetchProfile()
})
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

/* 主内容区 */
.main-content {
  flex: 1;
  padding: 40px;
}

.profile-container {
  max-width: 900px;
  margin: 0 auto;
}

/* 进度步骤 */
.progress-section {
  margin-bottom: 24px;
}

.progress-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0;
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.step {
  display: flex;
  align-items: center;
  gap: 12px;
  opacity: 0.4;
  transition: all 0.3s ease;
}

.step.active {
  opacity: 1;
}

.step.completed {
  opacity: 1;
}

.step-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #e4e7ed;
  color: #909399;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  transition: all 0.3s ease;
}

.step.active .step-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.step.completed .step-icon {
  background: #67c23a;
  color: white;
}

.step-info {
  text-align: left;
}

.step-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.step-desc {
  font-size: 12px;
  color: #909399;
}

.step-line {
  width: 60px;
  height: 2px;
  background: #e4e7ed;
  margin: 0 16px;
  transition: all 0.3s ease;
}

.step-line.active {
  background: #667eea;
}

.step-line.completed {
  background: #67c23a;
}

/* 完成提示 */
.completed-banner {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-top: 16px;
  padding: 20px 24px;
  background: linear-gradient(135deg, #f6ffed 0%, #e8f5e9 100%);
  border-radius: 12px;
  border: 1px solid #67c23a;
}

.completed-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #67c23a;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.completed-icon .el-icon {
  font-size: 24px;
}

.completed-content h3 {
  font-size: 16px;
  color: #67c23a;
  margin: 0 0 4px 0;
}

.completed-content p {
  font-size: 13px;
  color: #606266;
  margin: 0;
}

/* 表单卡片 */
.profile-card {
  border-radius: 16px;
  border: 1px solid #ebeef5;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.card-header h2 {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.sub-title {
  font-size: 13px;
  color: #909399;
}

/* 表单分区 */
.form-body {
  padding: 24px 0;
}

.form-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #ebeef5;
}

.form-section:last-of-type {
  border-bottom: none;
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.title-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.title-icon.basic {
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
}

.title-icon.education {
  background: linear-gradient(135deg, #67c23a 0%, #95d475 100%);
}

.title-icon.skills {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
}

.title-icon.resume {
  background: linear-gradient(135deg, #f56c6c 0%, #fab6b6 100%);
}

.section-title span {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

/* 技能预览 */
.skill-preview {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.skill-tag {
  border-radius: 6px;
}

/* 上传区域 */
.upload-area {
  width: 100%;
}

.resume-upload {
  width: 100%;
}

:deep(.resume-upload .el-upload-dragger) {
  border-radius: 12px;
  border: 2px dashed #d9d9d9;
  background: #fafafa;
  transition: all 0.3s ease;
}

:deep(.resume-upload .el-upload-dragger:hover) {
  border-color: #667eea;
  background: #f5f7fa;
}

.upload-content {
  padding: 32px;
  text-align: center;
}

.upload-icon {
  width: 48px;
  height: 48px;
  color: #c0c4cc;
  margin-bottom: 12px;
}

.upload-icon .el-icon {
  font-size: 48px;
}

.upload-text p {
  margin: 0;
  font-size: 14px;
  color: #606266;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

/* 简历预览 */
.resume-preview {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 16px;
  padding: 16px;
  background: #f0f9ff;
  border-radius: 12px;
  border: 1px solid #409eff;
}

.preview-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-info {
  flex: 1;
}

.preview-name {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

/* 操作按钮 */
.form-actions {
  display: flex;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
}

.form-actions .el-button {
  min-width: 140px;
}
</style>