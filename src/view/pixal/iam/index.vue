<template>
  <div class="iam-management form-white">
    <div class="page-header">
      <h3>{{ $t('lang.iam.title') }}</h3>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input 
        v-model="search.email" 
        :placeholder="$t('lang.iam.search_email')"
        style="max-width: 250px;"
        clearable>
      </el-input>
      
      <el-select 
        v-model="search.status" 
        :placeholder="$t('lang.iam.select_status')"
        clearable
        style="max-width: 200px;"
        @change="handleSearch">
        <el-option :label="$t('lang.iam.enable_account')" value="1"></el-option>
        <el-option :label="$t('lang.iam.forbidden_account')" value="2"></el-option>
      </el-select>

      <div class="form-actions">
        <el-button type="primary" icon="Search" @click="handleSearch">{{ $t('lang.search') }}</el-button>
        <el-button type="primary" icon="Plus" @click="handleAdd">{{ $t('lang.iam.add_account') }}</el-button>
      </div>
    </div>

    <!-- 账号列表表格 -->
    <el-table :data="tableData" style="width: 100%" border v-loading="loading">
      <el-table-column prop="email" :label="$t('lang.email')" min-width="200"></el-table-column>
      <el-table-column prop="nickname" :label="$t('lang.iam.nickname')" min-width="150"></el-table-column>
      <el-table-column prop="status" :label="$t('lang.status')" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? $t('lang.iam.enable_account') : $t('lang.iam.forbidden_account') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="CreatedAt" :label="$t('lang.created_time')" width="180">
        <template #default="{ row }">
          {{ formatDateYYYYMMDD(row.CreatedAt) }}
        </template>
      </el-table-column>
      <!-- <el-table-column prop="lastLoginAt" :label="$t('lang.iam.last_login')" width="180"></el-table-column> -->
      <el-table-column :label="$t('lang.operation')" width="320" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleEdit(row)">{{ $t('lang.edit') }}</el-button>
          <el-button 
            :type="row.status === 1 ? 'warning' : 'success'" 
            size="small" 
            @click="handleToggleStatus(row)">
            {{ row.status === 1 ? $t('lang.iam.forbidden_account') : $t('lang.iam.enable_account') }}
          </el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">{{ $t('lang.delete') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSearch"
        @current-change="handleSearch">
      </el-pagination>
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? $t('lang.iam.edit_account') : $t('lang.iam.add_account')"
      width="500px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item :label="$t('lang.email')" prop="email">
          <el-input v-model="form.email" :disabled="isEdit"></el-input>
        </el-form-item>
        <el-form-item :label="$t('lang.iam.nickname')" prop="nickname">
          <el-input v-model="form.nickname"></el-input>
        </el-form-item>
        <el-form-item v-if="!isEdit" :label="$t('lang.password')" prop="password">
          <el-input v-model="form.password" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item :label="$t('lang.iam.role')" prop="roles">
          <el-select v-model="form.roles" multiple style="width: 100%;" placement="bottom-start">
            <el-option 
              v-for="role in rolesList" 
              :key="role.ID" 
              :label="$t(`lang.iam.roles.${role.roleName}`)" 
              :value="role.ID"
              :disabled="role.isDefault">
              <span>{{ $t(`lang.iam.roles.${role.roleName}`) }}</span>
              <span v-if="role.description" style="color: #999; margin-left: 8px; font-size: 12px;">{{ $t(`lang.iam.roles.${role.roleName}_desc`) }}</span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('lang.cancel') }}</el-button>
        <el-button v-if="isEdit" type="warning" @click="openPasswordDialog">{{ $t('lang.iam.set_password') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('lang.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 设置密码弹窗 -->
    <el-dialog 
      v-model="passwordDialogVisible" 
      :title="$t('lang.iam.set_password')"
      width="500px">
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordFormRules" label-width="160px">
        <el-form-item :label="$t('lang.new_password')" prop="password">
          <el-input v-model="passwordForm.password" type="password" show-password></el-input>
        </el-form-item>
        <el-form-item :label="$t('lang.confirm_new_password')" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">{{ $t('lang.cancel') }}</el-button>
        <el-button type="primary" @click="handleSetPassword">{{ $t('lang.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getIamUserList, createIamUser, updateIamUser, deleteIamUser, toggleIamUserStatus, getIamRolesList, resetIamUserPassword } from '@/api/iam'
import { formatDateYYYYMMDD } from '@/utils/format'
import { validateEmail, validatePassword } from '@/utils/validates'

const { t } = useI18n()

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref([])
const rolesList = ref([])
const dialogVisible = ref(false)
const passwordDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const passwordFormRef = ref()

const search = reactive({
  email: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  email: '',
  nickname: '',
  password: '',
  roles: []
})

const passwordForm = reactive({
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.password) {
    callback(new Error(t('lang.validation.password_mismatch')))
  } else {
    callback()
  }
}

const passwordFormRules = reactive({
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('lang.validation.confirm_password_required'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
})

const formRules = reactive({
  email: [
    { required: true, message: t('lang.validation.email_required'), trigger: 'blur' },
    { validator: validateEmail, trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: t('lang.iam.nickname_required'), trigger: 'blur' }
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ],
  roles: [
    { required: true, message: t('lang.iam.role_required'), trigger: 'change' }
  ]
})

onMounted(() => {
  handleSearch()
  fetchRolesList()
})

// 获取角色列表
const fetchRolesList = async () => {
  try {
    const res = await getIamRolesList()
    if (res.code === 0) {
      rolesList.value = res.data || []
    }
  } catch (e) {
    console.error('Failed to fetch roles list', e)
  }
}

// 获取默认必选的角色ID列表
const getDefaultRoleIds = () => {
  return rolesList.value.filter(role => role.isDefault).map(role => role.ID)
}

// 检查角色是否为默认必选
const isDefaultRole = (roleId) => {
  const role = rolesList.value.find(r => r.ID === roleId)
  return role?.isDefault === true
}

// 根据ID获取角色信息
const getRoleById = (roleId) => {
  return rolesList.value.find(r => r.ID === roleId)
}

const handleSearch = async () => {
  loading.value = true
  try {
    const res = await getIamUserList({
      email: search.email,
      status: search.status,
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.id = null
  form.email = ''
  form.nickname = ''
  form.password = ''
  form.roles = getDefaultRoleIds()
}

const handleAdd = () => {
  resetForm()
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = (row) => {
  resetForm()
  isEdit.value = true
  form.id = row.ID
  form.email = row.email
  form.nickname = row.nickname
  form.roles = row.roles || getDefaultRoleIds()
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return

  submitLoading.value = true
  try {
    if (isEdit.value) {
      const res = await updateIamUser({
        id: form.id,
        nickname: form.nickname,
        roles: form.roles
      })
      if (res.code === 0) {
        ElMessage.success(t('lang.success'))
        dialogVisible.value = false
        handleSearch()
      }
    } else {
      const res = await createIamUser({
        email: form.email,
        nickname: form.nickname,
        password: form.password,
        roles: form.roles
      })
      if (res.code === 0) {
        ElMessage.success(t('lang.success'))
        dialogVisible.value = false
        handleSearch()
      }
    }
  } finally {
    submitLoading.value = false
  }
}

const handleToggleStatus = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('lang.iam.toggle_status_confirm'),
      t('lang.warning'),
      { type: 'warning' }
    )
    
    const res = await toggleIamUserStatus({
      id: row.ID,
      status: row.status === 1 ? 2 : 1
    })
    if (res.code === 0) {
      ElMessage.success(t('lang.success'))
      handleSearch()
    }
  } catch {
    // 用户取消
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      t('lang.iam.delete_confirm'),
      t('lang.warning'),
      { type: 'warning' }
    )
    
    const res = await deleteIamUser({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success(t('lang.success'))
      handleSearch()
    }
  } catch {
    // 用户取消
  }
}

const openPasswordDialog = () => {
  passwordForm.password = ''
  passwordForm.confirmPassword = ''
  passwordDialogVisible.value = true
}

const handleSetPassword = async () => {
  const valid = await passwordFormRef.value.validate()
  if (!valid) return

  const res = await resetIamUserPassword({
    id: form.id,
    password: passwordForm.password
  })
  if (res.code === 0) {
    ElMessage.success(t('lang.success'))
    passwordDialogVisible.value = false
  }
}
</script>

<style scoped lang="scss">
.iam-management {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
  
  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }
}

.search-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 20px;
  align-items: center;
}

.form-actions {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  margin-top: 20px;
//   display: flex;
  justify-content: flex-start;
}
</style>
