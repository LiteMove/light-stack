<template>
  <div class="profile-container">
    <div class="profile-header">
      <h2>个人中心</h2>
      <p>管理您的个人信息和账户设置</p>
    </div>

    <el-tabs v-model="activeTab" class="profile-tabs">
      <!-- 个人信息标签页 -->
      <el-tab-pane label="个人信息" name="profile">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>

          <el-form
            ref="profileFormRef"
            :model="profileForm"
            :rules="profileRules"
            label-width="120px"
            class="profile-form"
          >
            <el-form-item label="用户名">
              <el-input v-model="profileForm.username" disabled />
            </el-form-item>

            <el-form-item label="昵称" prop="nickname">
              <el-input
                v-model="profileForm.nickname"
                placeholder="请输入昵称"
                clearable
              />
            </el-form-item>

            <el-form-item label="邮箱" prop="email">
              <el-input
                v-model="profileForm.email"
                placeholder="请输入邮箱地址"
                clearable
              />
            </el-form-item>

            <el-form-item label="手机号" prop="phone">
              <el-input
                v-model="profileForm.phone"
                placeholder="请输入手机号码"
                clearable
              />
            </el-form-item>

            <el-form-item label="角色">
              <el-tag
                v-for="role in profileInfo?.roles || []"
                :key="role.id"
                type="primary"
                class="role-tag"
              >
                {{ role.name }}
              </el-tag>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="updateProfile" :loading="updating">
                保存修改
              </el-button>
              <el-button @click="resetProfileForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 修改密码卡片 -->
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>修改密码</span>
            </div>
          </template>

          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="120px"
            class="profile-form"
          >
            <el-form-item label="原密码" prop="oldPassword">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入原密码"
                show-password
                clearable
              />
            </el-form-item>

            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码（至少6位）"
                show-password
                clearable
              />
            </el-form-item>

            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
                show-password
                clearable
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="changePassword" :loading="changingPassword">
                修改密码
              </el-button>
              <el-button @click="resetPasswordForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 租户配置标签页（仅租户管理员可见） -->
      <el-tab-pane v-if="isTenantAdmin" label="租户配置" name="tenantConfig">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>文件存储配置</span>
              <el-tag type="warning" size="small">仅租户管理员可修改</el-tag>
            </div>
          </template>

          <el-form
            ref="tenantConfigFormRef"
            :model="tenantConfigForm"
            :rules="tenantConfigRules"
            label-width="150px"
            class="profile-form"
          >
            <el-form-item label="存储类型" prop="fileStorage.type">
              <el-radio-group v-model="tenantConfigForm.fileStorage.type">
                <el-radio label="local">本地存储</el-radio>
                <el-radio label="oss" disabled>OSS存储（暂未实现）</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="默认公开访问" prop="fileStorage.defaultPublic">
              <el-switch
                v-model="tenantConfigForm.fileStorage.defaultPublic"
                active-text="是"
                inactive-text="否"
              />
            </el-form-item>

            <el-form-item label="最大文件大小" prop="fileStorage.maxFileSize">
              <el-input-number
                v-model="fileSizeMB"
                :min="1"
                :max="1000"
                controls-position="right"
                @change="updateMaxFileSize"
              />
              <span style="margin-left: 8px;">MB</span>
            </el-form-item>

            <el-form-item label="允许的文件类型" prop="fileStorage.allowedTypes">
              <el-checkbox-group v-model="tenantConfigForm.fileStorage.allowedTypes">
                <el-checkbox label=".jpg">JPG</el-checkbox>
                <el-checkbox label=".jpeg">JPEG</el-checkbox>
                <el-checkbox label=".png">PNG</el-checkbox>
                <el-checkbox label=".gif">GIF</el-checkbox>
                <el-checkbox label=".pdf">PDF</el-checkbox>
                <el-checkbox label=".doc">DOC</el-checkbox>
                <el-checkbox label=".docx">DOCX</el-checkbox>
                <el-checkbox label=".xls">XLS</el-checkbox>
                <el-checkbox label=".xlsx">XLSX</el-checkbox>
                <el-checkbox label=".txt">TXT</el-checkbox>
              </el-checkbox-group>
            </el-form-item>

            <!-- OSS配置（预留，当前禁用） -->
            <template v-if="tenantConfigForm.fileStorage.type === 'oss'">
              <el-divider content-position="left">OSS配置</el-divider>

              <el-form-item label="OSS提供商">
                <el-select v-model="tenantConfigForm.fileStorage.ossProvider" disabled>
                  <el-option label="阿里云OSS" value="aliyun" />
                  <el-option label="腾讯云COS" value="tencent" />
                  <el-option label="AWS S3" value="aws" />
                  <el-option label="七牛云Kodo" value="qiniu" />
                  <el-option label="又拍云USS" value="upyun" />
                </el-select>
              </el-form-item>

              <el-form-item label="Endpoint">
                <el-input v-model="tenantConfigForm.fileStorage.ossEndpoint" disabled />
              </el-form-item>

              <el-form-item label="Bucket">
                <el-input v-model="tenantConfigForm.fileStorage.ossBucket" disabled />
              </el-form-item>

              <el-form-item label="Access Key">
                <el-input v-model="tenantConfigForm.fileStorage.ossAccessKey" disabled />
              </el-form-item>

              <el-form-item label="Secret Key">
                <el-input v-model="tenantConfigForm.fileStorage.ossSecretKey" type="password" disabled />
              </el-form-item>
            </template>

            <el-form-item>
              <el-button type="primary" @click="updateTenantConfig" :loading="updatingTenantConfig">
                保存配置
              </el-button>
              <el-button @click="resetTenantConfigForm">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { profileApi, type ProfileInfo, type TenantConfig } from '@/api/profile'

// 响应式数据
const activeTab = ref('profile')
const profileInfo = ref<ProfileInfo | null>(null)
const updating = ref(false)
const changingPassword = ref(false)
const updatingTenantConfig = ref(false)

// 表单引用
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const tenantConfigFormRef = ref<FormInstance>()

// 个人信息表单
const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: ''
})

// 密码修改表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 租户配置表单
const tenantConfigForm = reactive<TenantConfig>({
  fileStorage: {
    type: 'local',
    baseUrl: '',
    defaultPublic: false,
    maxFileSize: 50 * 1024 * 1024, // 50MB
    allowedTypes: ['.jpg', '.jpeg', '.png', '.gif', '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.txt'],
    localPath: 'uploads',
    ossProvider: 'aliyun',
    ossEndpoint: '',
    ossRegion: '',
    ossBucket: '',
    ossAccessKey: '',
    ossSecretKey: '',
    ossCustomDomain: ''
  }
})

// 文件大小（MB）
const fileSizeMB = ref(50)

// 是否为租户管理员
const isTenantAdmin = computed(() => {
  return profileInfo.value?.roles?.some(role =>
    role.name === '租户管理员' || role.name === 'tenant_admin'
  ) || false
})

// 表单验证规则
const profileRules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 50, message: '昵称长度在1到50个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const tenantConfigRules: FormRules = {
  'fileStorage.type': [
    { required: true, message: '请选择存储类型', trigger: 'change' }
  ],
  'fileStorage.maxFileSize': [
    { required: true, message: '请设置最大文件大小', trigger: 'blur' }
  ],
  'fileStorage.allowedTypes': [
    { required: true, type: 'array', min: 1, message: '请至少选择一种文件类型', trigger: 'change' }
  ]
}

// 获取个人信息
const getProfileInfo = async () => {
  try {
    const { data } = await profileApi.getProfile()
    profileInfo.value = data
    // 填充表单数据
    profileForm.username = data.username
    profileForm.nickname = data.nickname
    profileForm.email = data.email || ''
    profileForm.phone = data.phone || ''
  } catch (error) {
    console.error('获取个人信息失败:', error)
    ElMessage.error('获取个人信息失败')
  }
}

// 获取租户配置
const getTenantConfig = async () => {
  if (!isTenantAdmin.value) return

  try {
    const { data } = await profileApi.getTenantConfig()
    Object.assign(tenantConfigForm, data)
    fileSizeMB.value = Math.round(data.fileStorage.maxFileSize / (1024 * 1024))
  } catch (error) {
    console.error('获取租户配置失败:', error)
    ElMessage.error('获取租户配置失败')
  }
}

// 更新个人信息
const updateProfile = async () => {
  if (!profileFormRef.value) return

  const isValid = await profileFormRef.value.validate().catch(() => false)
  if (!isValid) return

  try {
    updating.value = true
    await profileApi.updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone || undefined
    })
    ElMessage.success('个人信息更新成功')
    // 重新获取信息
    await getProfileInfo()
  } catch (error) {
    console.error('更新个人信息失败:', error)
    ElMessage.error('更新个人信息失败')
  } finally {
    updating.value = false
  }
}

// 修改密码
const changePassword = async () => {
  if (!passwordFormRef.value) return

  const isValid = await passwordFormRef.value.validate().catch(() => false)
  if (!isValid) return

  try {
    await ElMessageBox.confirm('修改密码后需要重新登录，是否继续？', '确认修改', {
      type: 'warning'
    })

    changingPassword.value = true
    await profileApi.changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    ElMessage.success('密码修改成功，请重新登录')
    // 清空表单
    resetPasswordForm()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('修改密码失败:', error)
      ElMessage.error('修改密码失败')
    }
  } finally {
    changingPassword.value = false
  }
}

// 更新租户配置
const updateTenantConfig = async () => {
  if (!tenantConfigFormRef.value) return

  const isValid = await tenantConfigFormRef.value.validate().catch(() => false)
  if (!isValid) return

  try {
    updatingTenantConfig.value = true
    await profileApi.updateTenantConfig(tenantConfigForm)
    ElMessage.success('租户配置更新成功')
  } catch (error) {
    console.error('更新租户配置失败:', error)
    ElMessage.error('更新租户配置失败')
  } finally {
    updatingTenantConfig.value = false
  }
}

// 重置表单
const resetProfileForm = () => {
  if (profileInfo.value) {
    profileForm.nickname = profileInfo.value.nickname
    profileForm.email = profileInfo.value.email || ''
    profileForm.phone = profileInfo.value.phone || ''
  }
  profileFormRef.value?.clearValidate()
}

const resetPasswordForm = () => {
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  passwordFormRef.value?.clearValidate()
}

const resetTenantConfigForm = () => {
  getTenantConfig()
  tenantConfigFormRef.value?.clearValidate()
}

// 更新最大文件大小
const updateMaxFileSize = () => {
  tenantConfigForm.fileStorage.maxFileSize = fileSizeMB.value * 1024 * 1024
}

// 组件挂载时获取数据
onMounted(() => {
  getProfileInfo()
  // 如果是租户管理员，获取租户配置
  setTimeout(() => {
    if (isTenantAdmin.value) {
      getTenantConfig()
    }
  }, 500) // 等待角色信息加载完成
})
</script>

<style scoped lang="scss">
.profile-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.profile-header {
  margin-bottom: 24px;

  h2 {
    margin: 0 0 8px 0;
    color: #303133;
    font-size: 24px;
    font-weight: 600;
  }

  p {
    margin: 0;
    color: #606266;
    font-size: 14px;
  }
}

.profile-tabs {
  :deep(.el-tabs__header) {
    margin: 0 0 20px 0;
  }
}

.profile-card {
  margin-bottom: 20px;

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-weight: 600;
  }
}

.profile-form {
  max-width: 600px;

  .role-tag {
    margin-right: 8px;
    margin-bottom: 4px;
  }

  :deep(.el-form-item__label) {
    font-weight: 500;
  }
}

:deep(.el-tabs__content) {
  padding-top: 0;
}

:deep(.el-card__body) {
  padding: 24px;
}

:deep(.el-divider__text) {
  color: #909399;
  font-weight: 500;
}
</style>