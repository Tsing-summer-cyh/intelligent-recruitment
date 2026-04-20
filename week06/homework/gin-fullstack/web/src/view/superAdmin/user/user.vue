<template>
  <div>
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="searchInfo.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addUser">
          新增用户
        </el-button>
      </div>

      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="头像" min-width="75">
          <template #default="scope">
            <CustomPic style="margin-top: 8px" :pic-src="scope.row.headerImg" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="70" prop="ID" />
        <el-table-column align="left" label="用户名" min-width="140" prop="userName" />
        <el-table-column align="left" label="昵称" min-width="140" prop="nickName" />
        <el-table-column align="left" label="手机号" min-width="160" prop="phone" />
        <el-table-column align="left" label="邮箱" min-width="180" prop="email" />
        <el-table-column align="left" label="登录IP" min-width="160" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.lastLoginIp || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="登录时间" min-width="180">
          <template #default="scope">
            {{ formatLastLoginAt(scope.row.lastLoginAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="用户角色" min-width="220">
          <template #default="scope">
            <el-cascader
              v-model="scope.row.authorityIds"
              :options="authOptions"
              :show-all-levels="false"
              collapse-tags
              :props="cascaderProps"
              :clearable="false"
              @visible-change="
                (flag) => {
                  changeAuthority(scope.row, flag)
                }
              "
              @remove-tag="
                () => {
                  changeAuthority(scope.row, false)
                }
              "
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="启用" min-width="120">
          <template #default="scope">
            <el-switch
              v-model="scope.row.enable"
              inline-prompt
              :active-value="1"
              :inactive-value="2"
              @change="
                () => {
                  switchEnable(scope.row)
                }
              "
            />
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          :min-width="appStore.operateMinWith"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteUserFunc(scope.row)"
            >
              删除
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              @click="openEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="primary"
              link
              icon="magic-stick"
              @click="resetPasswordFunc(scope.row)"
            >
              重置密码
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog
      v-model="resetPwdDialog"
      title="重置密码"
      width="500px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form ref="resetPwdForm" :model="resetPwdInfo" label-width="100px">
        <el-form-item label="用户账号">
          <el-input v-model="resetPwdInfo.userName" disabled />
        </el-form-item>
        <el-form-item label="用户昵称">
          <el-input v-model="resetPwdInfo.nickName" disabled />
        </el-form-item>
        <el-form-item label="新密码">
          <div class="flex w-full">
            <el-input
              v-model="resetPwdInfo.password"
              class="flex-1"
              placeholder="请输入新密码"
              show-password
            />
            <el-button
              type="primary"
              style="margin-left: 10px"
              @click="generateRandomPassword"
            >
              生成随机密码
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeResetPwdDialog">
            取消
          </el-button>
          <el-button type="primary" @click="confirmResetPassword">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer
      v-model="addUserDialog"
      :size="appStore.drawerSize"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogFlag === 'add' ? '新增用户' : '编辑用户' }}</span>
          <div>
            <el-button @click="closeAddUserDialog">
              取消
            </el-button>
            <el-button type="primary" @click="enterAddUserDialog">
              确定
            </el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="userForm"
        :rules="rules"
        :model="userInfo"
        label-width="80px"
      >
        <el-form-item v-if="dialogFlag === 'add'" label="用户名" prop="userName">
          <el-input v-model="userInfo.userName" />
        </el-form-item>
        <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password">
          <el-input v-model="userInfo.password" show-password />
        </el-form-item>
        <el-form-item label="昵称" prop="nickName">
          <el-input v-model="userInfo.nickName" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="userInfo.phone" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userInfo.email" />
        </el-form-item>
        <el-form-item label="用户角色" prop="authorityIds">
          <el-cascader
            v-model="userInfo.authorityIds"
            style="width: 100%"
            :options="authOptions"
            :show-all-levels="false"
            :props="cascaderProps"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item label="启用" prop="enable">
          <el-switch
            v-model="userInfo.enable"
            inline-prompt
            :active-value="1"
            :inactive-value="2"
          />
        </el-form-item>
        <el-form-item label="头像">
          <SelectImage v-model="userInfo.headerImg" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import { nextTick, ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'

  import { getAuthorityList } from '@/api/authority'
  import {
    deleteUser,
    getUserList,
    register,
    resetPassword,
    setUserAuthorities,
    setUserInfo
  } from '@/api/user'
  import CustomPic from '@/components/customPic/index.vue'
  import SelectImage from '@/components/selectImage/selectImage.vue'
  import WarningBar from '@/components/warningBar/warningBar.vue'
  import { formatTimeToStr } from '@/utils/date'
  import { useAppStore } from '@/pinia'

  defineOptions({
    name: 'User'
  })

  const appStore = useAppStore()

  const searchForm = ref(null)
  const userForm = ref(null)
  const resetPwdForm = ref(null)

  const createSearchInfo = () => ({
    username: '',
    nickname: '',
    phone: '',
    email: ''
  })

  const createUserInfo = () => ({
    ID: undefined,
    userName: '',
    password: '',
    nickName: '',
    phone: '',
    email: '',
    headerImg: '',
    authorityId: '',
    authorityIds: [],
    enable: 1
  })

  const createResetPwdInfo = () => ({
    ID: '',
    userName: '',
    nickName: '',
    password: ''
  })

  const searchInfo = ref(createSearchInfo())
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const authOptions = ref([])
  const addUserDialog = ref(false)
  const dialogFlag = ref('add')
  const userInfo = ref(createUserInfo())
  const resetPwdDialog = ref(false)
  const resetPwdInfo = ref(createResetPwdInfo())
  const tempAuth = {}

  const cascaderProps = {
    multiple: true,
    checkStrictly: true,
    label: 'authorityName',
    value: 'authorityId',
    disabled: 'disabled',
    emitPath: false
  }

  const rules = {
    userName: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 5, message: '至少 5 个字符', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入用户密码', trigger: 'blur' },
      { min: 6, message: '至少 6 个字符', trigger: 'blur' }
    ],
    nickName: [{ required: true, message: '请输入用户昵称', trigger: 'blur' }],
    phone: [
      {
        pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/,
        message: '请输入合法手机号',
        trigger: 'blur'
      }
    ],
    email: [
      {
        pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/,
        message: '请输入正确的邮箱',
        trigger: 'blur'
      }
    ],
    authorityIds: [
      { required: true, message: '请选择用户角色', trigger: 'change' }
    ]
  }

  const formatLastLoginAt = (value) => {
    if (!value) {
      return '-'
    }
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) {
      return '-'
    }
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm')
  }

  const setAuthorityOptions = (authorityData, optionsData) => {
    authorityData?.forEach((item) => {
      const option = {
        authorityId: item.authorityId,
        authorityName: item.authorityName
      }
      if (item.children?.length) {
        option.children = []
        setAuthorityOptions(item.children, option.children)
      }
      optionsData.push(option)
    })
  }

  const setOptions = (authData) => {
    authOptions.value = []
    setAuthorityOptions(authData, authOptions.value)
  }

  const setAuthorityIds = () => {
    tableData.value.forEach((user) => {
      user.authorityIds =
        user.authorities?.map((item) => item.authorityId) || []
    })
  }

  const getTableData = async () => {
    const table = await getUserList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list || []
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
      setAuthorityIds()
    }
  }

  const initPage = async () => {
    await getTableData()
    const res = await getAuthorityList()
    setOptions(res.data || [])
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const onReset = () => {
    searchInfo.value = createSearchInfo()
    searchForm.value?.resetFields?.()
    page.value = 1
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const addUser = () => {
    dialogFlag.value = 'add'
    userInfo.value = createUserInfo()
    addUserDialog.value = true
  }

  const openEdit = (row) => {
    dialogFlag.value = 'edit'
    userInfo.value = {
      ...createUserInfo(),
      ...JSON.parse(JSON.stringify(row)),
      authorityIds: row.authorityIds || []
    }
    addUserDialog.value = true
  }

  const closeAddUserDialog = () => {
    userForm.value?.resetFields?.()
    userInfo.value = createUserInfo()
    addUserDialog.value = false
  }

  const enterAddUserDialog = async () => {
    if (!userInfo.value.authorityIds.length) {
      ElMessage.warning('请选择至少一个角色')
      return
    }

    await userForm.value.validate(async (valid) => {
      if (!valid) {
        return
      }

      const req = {
        ...userInfo.value,
        authorityId: userInfo.value.authorityIds[0]
      }

      if (dialogFlag.value === 'add') {
        const res = await register(req)
        if (res.code === 0) {
          ElMessage.success('创建成功')
          await getTableData()
          closeAddUserDialog()
        }
        return
      }

      const res = await setUserInfo(req)
      if (res.code === 0) {
        ElMessage.success('编辑成功')
        await getTableData()
        closeAddUserDialog()
      }
    })
  }

  const deleteUserFunc = async (row) => {
    await ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const res = await deleteUser({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  }

  const changeAuthority = async (row, flag) => {
    if (flag) {
      tempAuth[row.ID] = [...(row.authorityIds || [])]
      return
    }

    await nextTick()

    if (!row.authorityIds?.length) {
      row.authorityIds = [...(tempAuth[row.ID] || [])]
      ElMessage.warning('用户至少保留一个角色')
      return
    }

    const res = await setUserAuthorities({
      ID: row.ID,
      authorityIds: row.authorityIds
    })

    if (res.code === 0) {
      ElMessage.success('角色设置成功')
      delete tempAuth[row.ID]
      return
    }

    row.authorityIds = [...(tempAuth[row.ID] || [])]
    delete tempAuth[row.ID]
  }

  const switchEnable = async (row) => {
    const req = {
      ...JSON.parse(JSON.stringify(row)),
      authorityId: row.authorityIds?.[0]
    }
    const res = await setUserInfo(req)
    if (res.code === 0) {
      ElMessage.success(`${req.enable === 2 ? '禁用' : '启用'}成功`)
      await getTableData()
    }
  }

  const resetPasswordFunc = (row) => {
    resetPwdInfo.value = {
      ID: row.ID,
      userName: row.userName,
      nickName: row.nickName,
      password: ''
    }
    resetPwdDialog.value = true
  }

  const closeResetPwdDialog = () => {
    resetPwdForm.value?.resetFields?.()
    resetPwdInfo.value = createResetPwdInfo()
    resetPwdDialog.value = false
  }

  const generateRandomPassword = async () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    let password = ''
    for (let i = 0; i < 12; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    resetPwdInfo.value.password = password
    try {
      await navigator.clipboard.writeText(password)
      ElMessage.success('随机密码已复制到剪贴板')
    } catch {
      ElMessage.warning('复制失败，请手动复制')
    }
  }

  const confirmResetPassword = async () => {
    if (!resetPwdInfo.value.password) {
      ElMessage.warning('请输入或生成新密码')
      return
    }

    const res = await resetPassword({
      ID: resetPwdInfo.value.ID,
      password: resetPwdInfo.value.password
    })

    if (res.code === 0) {
      ElMessage.success(res.msg || '密码重置成功')
      closeResetPwdDialog()
    }
  }

  initPage()
</script>
