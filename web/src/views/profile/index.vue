<template>
  <div class="profile-container">
    <div class="profile-header">
      <h2>个人中心</h2>
      <p>管理您的个人信息和账户设置</p>
    </div>

    <div class="profile-content">
      <el-tabs v-model="activeTab" class="profile-tabs">
      <!-- 个人信息标签页 -->
      <el-tab-pane name="profile">
        <template #label>
          <span class="tab-label">
            <el-icon><User /></el-icon>
            个人信息
          </span>
        </template>
        
        <div class="profile-layout">
          <!-- 左侧：个人信息 -->
          <div class="profile-left">
            <el-card class="profile-card">
              <template #header>
                <div class="card-header">
                  <span>
                    <el-icon><UserFilled /></el-icon>
                    基本信息
                  </span>
                </div>
              </template>

              <!-- 头像区域 -->
              <div class="avatar-section">
                <AvatarUpload
                  v-model="profileForm.avatar"
                  @success="handleAvatarSuccess"
                  @error="handleAvatarError"
                />
              </div>

              <el-form
                ref="profileFormRef"
                :model="profileForm"
                :rules="profileRules"
                label-width="120px"
                class="profile-form"
              >
                <el-form-item label="用户名">
                  <el-input v-model="profileForm.username" disabled>
                    <template #prefix>
                      <el-icon><User /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="昵称" prop="nickname">
                  <el-input
                    v-model="profileForm.nickname"
                    placeholder="请输入昵称"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Avatar /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="邮箱" prop="email">
                  <el-input
                    v-model="profileForm.email"
                    placeholder="请输入邮箱地址"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Message /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="手机号" prop="phone">
                  <el-input
                    v-model="profileForm.phone"
                    placeholder="请输入手机号码"
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Phone /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="角色">
                  <div class="roles-container">
                    <el-tag
                      v-for="role in profileInfo?.roles || []"
                      :key="role.id"
                      type="primary"
                      size="large"
                      class="role-tag"
                      effect="plain"
                    >
                      <el-icon><Star /></el-icon>
                      {{ role.name }}
                    </el-tag>
                    <el-text v-if="!profileInfo?.roles?.length" type="info">
                      暂无角色信息
                    </el-text>
                  </div>
                </el-form-item>

                <el-form-item>
                  <el-button type="primary" @click="updateProfile" :loading="updating">
                    <el-icon v-if="!updating"><Check /></el-icon>
                    保存修改
                  </el-button>
                  <el-button @click="resetProfileForm">
                    <el-icon><Refresh /></el-icon>
                    重置
                  </el-button>
                </el-form-item>
              </el-form>
            </el-card>
          </div>

          <!-- 右侧：密码修改 -->
          <div class="profile-right">
            <el-card class="profile-card">
              <template #header>
                <div class="card-header">
                  <span>
                    <el-icon><Key /></el-icon>
                    修改密码
                  </span>
                  <el-tag type="warning" size="small">修改后需重新登录</el-tag>
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
                  >
                    <template #prefix>
                      <el-icon><Lock /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="新密码" prop="newPassword">
                  <el-input
                    v-model="passwordForm.newPassword"
                    type="password"
                    placeholder="请输入新密码（至少6位）"
                    show-password
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Key /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item label="确认密码" prop="confirmPassword">
                  <el-input
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    placeholder="请再次输入新密码"
                    show-password
                    clearable
                  >
                    <template #prefix>
                      <el-icon><Key /></el-icon>
                    </template>
                  </el-input>
                </el-form-item>

                <el-form-item>
                  <el-button type="primary" @click="changePassword" :loading="changingPassword">
                    <el-icon v-if="!changingPassword"><Check /></el-icon>
                    修改密码
                  </el-button>
                  <el-button @click="resetPasswordForm">
                    <el-icon><Refresh /></el-icon>
                    重置
                  </el-button>
                </el-form-item>
              </el-form>
            </el-card>
          </div>
        </div>
      </el-tab-pane>

      <!-- 租户配置标签页（仅租户管理员可见） -->
      <el-tab-pane v-if="isTenantAdmin" name="tenantConfig">
        <template #label>
          <span class="tab-label">
            <el-icon><Setting /></el-icon>
            租户配置
          </span>
        </template>

        <!-- 超级管理员没有选择租户时的提示 -->
        <div v-if="isSuperAdmin && !currentTenantInfo" class="tenant-select-tip">
          <el-alert
            title="请先选择要配置的租户"
            type="warning"
            description="作为超级管理员，您需要先选择一个租户才能进行配置管理。请使用顶部的租户切换器选择目标租户。"
            :closable="false"
            show-icon>
          </el-alert>
        </div>

        <!-- 有租户信息时显示配置界面 -->
        <div v-else>
          <!-- 当前租户信息显示 -->
          <div class="tenant-info-bar" v-if="currentTenantInfo">
            <el-tag type="info" size="large" effect="dark">
              <el-icon><Setting /></el-icon>
              正在配置租户：{{ currentTenantInfo.name }}
            </el-tag>
          </div>
        
        <div class="tenant-config-layout">
          <!-- 左侧：系统基本信息 -->
          <div class="config-left">
            <el-card class="profile-card">
              <template #header>
                <div class="card-header">
                  <span>
                    <el-icon><Setting /></el-icon>
                    系统基本信息
                  </span>
                  <el-tag type="warning" size="small">仅租户管理员可修改</el-tag>
                </div>
              </template>

              <el-form
                ref="systemInfoFormRef"
                :model="tenantConfigForm"
                :rules="systemInfoRules"
                label-width="120px"
                class="profile-form"
              >
                <el-form-item label="系统名称" prop="systemName">
                  <el-input
                    v-model="tenantConfigForm.systemName"
                    placeholder="请输入系统名称"
                    clearable
                  />
                </el-form-item>

                <el-form-item label="系统Logo" prop="logo">
                  <div class="logo-upload-wrapper">
<!--                    <div class="logo-preview" v-if="tenantConfigForm.logo">-->
<!--                      <img-->
<!--                        :src="tenantConfigForm.logo"-->
<!--                        alt="系统Logo预览"-->
<!--                        class="logo-preview-image"-->
<!--                        @error="handleLogoImageError"-->
<!--                      />-->
<!--                    </div>-->
                    <AvatarUpload
                      v-model="tenantConfigForm.logo"
                      :size="80"
                      :max-size="5"
                      usage-type="system-logo"
                      @success="handleLogoUploadSuccess"
                      @error="handleLogoUploadError"
                    />
                    <el-input
                      v-model="tenantConfigForm.logo"
                      placeholder="或直接输入Logo URL地址"
                      clearable
                      class="logo-url-input"
                    />
                  </div>
                  <div class="form-tip">支持上传图片文件（推荐）或输入网络图片URL，建议尺寸80x80像素</div>
                </el-form-item>

                <el-form-item label="系统描述" prop="description">
                  <el-input
                    v-model="tenantConfigForm.description"
                    type="textarea"
                    :rows="4"
                    placeholder="请输入系统描述信息"
                    maxlength="200"
                    show-word-limit
                  />
                </el-form-item>

                <el-form-item label="版权信息" prop="copyright">
                  <el-input
                    v-model="tenantConfigForm.copyright"
                    placeholder="请输入版权信息"
                    clearable
                  />
                  <div class="form-tip">例如：© 2024 公司名称 版权所有</div>
                </el-form-item>
              </el-form>
            </el-card>
          </div>

          <!-- 右侧：文件存储配置 -->
          <div class="config-right">
            <el-card class="profile-card">
              <template #header>
                <div class="card-header">
                  <span>
                    <el-icon><Folder /></el-icon>
                    文件存储配置
                  </span>
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
                    <el-radio label="oss">OSS存储</el-radio>
                  </el-radio-group>
                </el-form-item>

                <!-- 本地存储配置 -->
                <template v-if="tenantConfigForm.fileStorage.type === 'local'">
                  <el-divider content-position="left">本地存储配置</el-divider>
                  
                  <el-form-item label="访问域名" prop="fileStorage.localAccessDomain">
                    <el-input
                      v-model="tenantConfigForm.fileStorage.localAccessDomain"
                      placeholder="请输入文件访问域名（如：https://files.example.com）"
                      clearable
                    />
                    <div class="form-tip">留空则使用默认的 /static 路径</div>
                  </el-form-item>
                </template>

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

                <!-- OSS配置 -->
                <template v-if="tenantConfigForm.fileStorage.type === 'oss'">
                  <el-divider content-position="left">OSS配置</el-divider>

                  <el-form-item label="OSS提供商" prop="fileStorage.ossProvider">
                    <el-select 
                      v-model="tenantConfigForm.fileStorage.ossProvider" 
                      placeholder="请选择OSS提供商"
                    >
                      <el-option label="阿里云OSS" value="aliyun" />
<!--                      <el-option label="腾讯云COS" value="tencent" />-->
<!--                      <el-option label="AWS S3" value="aws" />-->
<!--                      <el-option label="七牛云Kodo" value="qiniu" />-->
<!--                      <el-option label="又拍云USS" value="upyun" />-->
                    </el-select>
                  </el-form-item>

                  <el-form-item label="Endpoint" prop="fileStorage.ossEndpoint">
                    <el-input 
                      v-model="tenantConfigForm.fileStorage.ossEndpoint" 
                      placeholder="请输入OSS服务端点"
                      clearable
                    />
                    <div class="form-tip">例如：oss-cn-beijing.aliyuncs.com</div>
                  </el-form-item>

                  <el-form-item 
                    v-if="tenantConfigForm.fileStorage.ossProvider === 'aws'" 
                    label="Region" 
                    prop="fileStorage.ossRegion"
                  >
                    <el-input 
                      v-model="tenantConfigForm.fileStorage.ossRegion" 
                      placeholder="请输入AWS区域"
                      clearable
                    />
                    <div class="form-tip">例如：us-east-1</div>
                  </el-form-item>

                  <el-form-item label="Bucket" prop="fileStorage.ossBucket">
                    <el-input 
                      v-model="tenantConfigForm.fileStorage.ossBucket" 
                      placeholder="请输入存储桶名称"
                      clearable
                    />
                  </el-form-item>

                  <el-form-item label="Access Key" prop="fileStorage.ossAccessKey">
                    <el-input 
                      v-model="tenantConfigForm.fileStorage.ossAccessKey" 
                      placeholder="请输入访问密钥"
                      clearable
                    />
                  </el-form-item>

                  <el-form-item label="Secret Key" prop="fileStorage.ossSecretKey">
                    <el-input 
                      v-model="tenantConfigForm.fileStorage.ossSecretKey" 
                      type="password" 
                      placeholder="请输入访问密钥"
                      show-password
                      clearable
                    />
                  </el-form-item>

                  <el-form-item label="自定义域名" prop="fileStorage.ossCustomDomain">
                    <el-input 
                      v-model="tenantConfigForm.fileStorage.ossCustomDomain" 
                      placeholder="请输入OSS自定义域名（如：https://cdn.example.com）"
                      clearable
                    />
                    <div class="form-tip">设置后将使用此域名访问文件，否则使用OSS默认域名</div>
                  </el-form-item>
                </template>
              </el-form>
            </el-card>
          </div>
        </div>
        
        <!-- 底部操作按钮 -->
        <div class="config-actions" v-if="currentTenantInfo">
          <el-button type="primary" @click="updateTenantConfig" :loading="updatingTenantConfig">
            <el-icon v-if="!updatingTenantConfig"><Check /></el-icon>
            保存配置
          </el-button>
          <el-button @click="resetTenantConfigForm">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </div>
        </div> <!-- 关闭 v-else div -->
      </el-tab-pane>
    </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  User, UserFilled, Avatar, Message, Phone, Star, Check, Refresh,
  Lock, Key, Setting, Upload, Folder, Delete
} from '@element-plus/icons-vue'
import AvatarUpload from '@/components/AvatarUpload.vue'
import { profileApi, type ProfileInfo, type TenantConfig } from '@/api/profile'
import { useTenantStore } from '@/store/tenant'

// 租户store
const tenantStore = useTenantStore()

// 响应式数据
const activeTab = ref('profile')
const profileInfo = ref<ProfileInfo | null>(null)
const currentTenantInfo = ref<any>(null)
const updating = ref(false)
const changingPassword = ref(false)
const updatingTenantConfig = ref(false)

// 表单引用
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const systemInfoFormRef = ref<FormInstance>()
const tenantConfigFormRef = ref<FormInstance>()

// 个人信息表单
const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: ''
})

// 密码修改表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 租户配置表单
const tenantConfigForm = reactive<TenantConfig>({
  systemName: '',
  logo: '',
  description: '',
  copyright: '',
  fileStorage: {
    type: 'local',
    defaultPublic: false,
    maxFileSize: 50 * 1024 * 1024, // 50MB
    allowedTypes: ['.jpg', '.jpeg', '.png', '.gif', '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.txt'],
    // 本地存储配置
    localAccessDomain: '',
    // OSS配置
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

// 是否为租户管理员或超级管理员
const isTenantAdmin = computed(() => {
  if (!profileInfo.value?.roles) return false

  return profileInfo.value.roles.some(role =>
    role.name === '租户管理员' ||
    role.name === 'tenant_admin' ||
    role.name === 'super_admin' ||
    role.name === '超级管理员' ||
    role.code === 'super_admin' ||
    role.code === 'tenant_admin'
  )
})

// 是否为超级管理员
const isSuperAdmin = computed(() => {
  if (!profileInfo.value?.roles) return false

  return profileInfo.value.roles.some(role =>
    role.name === 'super_admin' ||
    role.name === '超级管理员' ||
    role.code === 'super_admin'
  )
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
  ],
  'fileStorage.localAccessDomain': [
    { max: 500, message: '访问域名长度不能超过500个字符', trigger: 'blur' }
  ],
  'fileStorage.ossProvider': [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (tenantConfigForm.fileStorage.type === 'oss' && !value) {
          callback(new Error('使用OSS存储时，请选择OSS提供商'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  'fileStorage.ossEndpoint': [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (tenantConfigForm.fileStorage.type === 'oss' && !value) {
          callback(new Error('使用OSS存储时，请填写Endpoint'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  'fileStorage.ossBucket': [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (tenantConfigForm.fileStorage.type === 'oss' && !value) {
          callback(new Error('使用OSS存储时，请填写Bucket名称'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  'fileStorage.ossAccessKey': [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (tenantConfigForm.fileStorage.type === 'oss' && !value) {
          callback(new Error('使用OSS存储时，请填写Access Key'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  'fileStorage.ossSecretKey': [
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (tenantConfigForm.fileStorage.type === 'oss' && !value) {
          callback(new Error('使用OSS存储时，请填写Secret Key'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const systemInfoRules: FormRules = {
  systemName: [
    { max: 100, message: '系统名称长度不能超过100个字符', trigger: 'blur' }
  ],
  logo: [
    { max: 500, message: 'Logo URL长度不能超过500个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '系统描述长度不能超过200个字符', trigger: 'blur' }
  ],
  copyright: [
    { max: 100, message: '版权信息长度不能超过100个字符', trigger: 'blur' }
  ]
}

// 头像上传成功处理
const handleAvatarSuccess = async (file: any) => {
  ElMessage.success('头像上传成功')
  // 自动保存个人信息以更新头像
  await updateProfileAvatar()
}

// 头像上传错误处理
const handleAvatarError = (error: string) => {
  ElMessage.error(error || '头像上传失败')
}

// Logo上传成功处理
const handleLogoUploadSuccess = (file: any) => {
  ElMessage.success('系统Logo上传成功')
  // Logo URL已经通过v-model自动更新到tenantConfigForm.logo
}

// Logo上传错误处理
const handleLogoUploadError = (error: string) => {
  ElMessage.error(error || '系统Logo上传失败')
}

// Logo图片加载错误处理
const handleLogoImageError = (event: Event) => {
  console.error('Logo图片加载失败:', tenantConfigForm.logo)
  // 可以在这里显示默认图片或提示
}

// 单独更新头像
const updateProfileAvatar = async () => {
  try {
    await profileApi.updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone || undefined,
      avatar: profileForm.avatar
    })
  } catch (error) {
    console.error('更新头像失败:', error)
    throw error
  }
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
    profileForm.avatar = data.avatar || ''
  } catch (error) {
    console.error('获取个人信息失败:', error)
    ElMessage.error('获取个人信息失败')
  }
}

// 获取租户配置
const getTenantConfig = async () => {
  if (!isTenantAdmin.value) return

  try {
    // 获取当前租户信息用于显示
    const currentTenant = tenantStore.getCurrentTenant()

    // 如果是超级管理员但没有选择租户，提示用户选择
    if (isSuperAdmin.value && !currentTenant) {
      ElMessage.warning('请先选择要配置的租户')
      return
    }

    currentTenantInfo.value = currentTenant

    const { data } = await profileApi.getTenantConfig()

    // 正确更新配置数据，确保响应式更新
    tenantConfigForm.systemName = data.systemName || ''
    tenantConfigForm.logo = data.logo || ''
    tenantConfigForm.description = data.description || ''
    tenantConfigForm.copyright = data.copyright || ''

    // 更新文件存储配置
    const fileStorage = data.fileStorage || {}
    tenantConfigForm.fileStorage.type = fileStorage.type || 'local'
    tenantConfigForm.fileStorage.defaultPublic = fileStorage.defaultPublic || false
    tenantConfigForm.fileStorage.maxFileSize = fileStorage.maxFileSize || (50 * 1024 * 1024)

    // 重要：正确更新数组，确保Vue能够检测到变化
    const allowedTypes = fileStorage.allowedTypes || ['.jpg', '.jpeg', '.png', '.gif', '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.txt']
    tenantConfigForm.fileStorage.allowedTypes.length = 0  // 清空现有数组
    tenantConfigForm.fileStorage.allowedTypes.push(...allowedTypes)  // 添加新数据

    // 本地存储配置
    tenantConfigForm.fileStorage.localAccessDomain = fileStorage.localAccessDomain || ''

    // OSS配置
    tenantConfigForm.fileStorage.ossProvider = fileStorage.ossProvider || 'aliyun'
    tenantConfigForm.fileStorage.ossEndpoint = fileStorage.ossEndpoint || ''
    tenantConfigForm.fileStorage.ossRegion = fileStorage.ossRegion || ''
    tenantConfigForm.fileStorage.ossBucket = fileStorage.ossBucket || ''
    tenantConfigForm.fileStorage.ossAccessKey = fileStorage.ossAccessKey || ''
    tenantConfigForm.fileStorage.ossSecretKey = fileStorage.ossSecretKey || ''
    tenantConfigForm.fileStorage.ossCustomDomain = fileStorage.ossCustomDomain || ''

    fileSizeMB.value = Math.round((tenantConfigForm.fileStorage.maxFileSize || 0) / (1024 * 1024))
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
  // 检查权限
  if (!isTenantAdmin.value) {
    ElMessage.error('没有权限修改租户配置')
    return
  }

  // 如果是超级管理员，检查是否选择了租户
  if (isSuperAdmin.value && !tenantStore.getCurrentTenant()) {
    ElMessage.warning('请先选择要配置的租户')
    return
  }

  // 验证系统信息表单
  let systemInfoValid = true
  if (systemInfoFormRef.value) {
    systemInfoValid = await systemInfoFormRef.value.validate().catch(() => false)
  }

  // 验证文件存储配置表单
  let tenantConfigValid = true
  if (tenantConfigFormRef.value) {
    tenantConfigValid = await tenantConfigFormRef.value.validate().catch(() => false)
  }

  if (!systemInfoValid || !tenantConfigValid) return

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
    profileForm.avatar = profileInfo.value.avatar || ''
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
  systemInfoFormRef.value?.clearValidate()
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

// 监听租户切换（用于超级管理员）
watch(
  () => tenantStore.currentTenant,
  (newTenant, oldTenant) => {
    if (newTenant?.id !== oldTenant?.id && isTenantAdmin.value) {
      // 租户切换时重新获取配置
      setTimeout(() => {
        getTenantConfig()
      }, 100)
    }
  },
  { deep: true }
)
</script>

<style scoped lang="scss">
.profile-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
  height: calc(100vh - 80px); // 减去导航栏高度等
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.profile-header {
  margin-bottom: 24px;
  flex-shrink: 0;

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

.profile-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.profile-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  
  :deep(.el-tabs__header) {
    margin: 0 0 20px 0;
    flex-shrink: 0;
  }
  
  :deep(.el-tabs__content) {
    flex: 1;
    overflow-y: auto;
    padding-top: 0;
    
    // 自定义滚动条样式
    &::-webkit-scrollbar {
      width: 6px;
    }
    
    &::-webkit-scrollbar-track {
      background: #f1f1f1;
      border-radius: 3px;
    }
    
    &::-webkit-scrollbar-thumb {
      background: #c1c1c1;
      border-radius: 3px;
      
      &:hover {
        background: #a1a1a1;
      }
    }
  }
  
  :deep(.el-tab-pane) {
    height: 100%;
  }
}

.profile-layout {
  display: flex;
  gap: 24px;
  height: 100%;
  
  @media (max-width: 1200px) {
    flex-direction: column;
    gap: 16px;
  }
}

.profile-left,
.profile-right {
  flex: 1;
  min-width: 0; // 防止flex子元素溢出
}

.profile-right {
  @media (max-width: 1200px) {
    flex: none;
  }
}

.avatar-section {
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
}

.profile-card {
  margin-bottom: 20px;
  height: fit-content;
  
  :deep(.el-card__header) {
    background-color: #fafbfc;
    border-bottom: 1px solid #ebeef5;
  }
}

.tenant-info-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;

  .el-tag {
    padding: 8px 16px;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;

    .el-icon {
      font-size: 16px;
    }
  }
}

.tenant-select-tip {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;

  .el-alert {
    max-width: 600px;
  }
}

.tenant-config-layout {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
  
  @media (max-width: 1400px) {
    flex-direction: column;
    gap: 20px;
  }
}

.config-left,
.config-right {
  flex: 1;
  min-width: 0;
}

.config-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding: 20px 0;
  border-top: 1px solid #ebeef5;
  margin-top: 20px;
}

.profile-form {
  max-width: 600px;

  .roles-container {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    min-height: 32px;
  }

  .role-tag {
    margin: 0;
    display: inline-flex;
    align-items: center;
    gap: 4px;
    
    .el-icon {
      font-size: 14px;
    }
  }

  :deep(.el-form-item__label) {
    font-weight: 500;
  }
  
  :deep(.el-input__prefix) {
    display: flex;
    align-items: center;
    
    .el-icon {
      color: #a8abb2;
      font-size: 16px;
    }
  }
  
  :deep(.el-button) {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    
    .el-icon {
      font-size: 16px;
    }
  }
}

.tab-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  
  .el-icon {
    font-size: 16px;
  }
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  
  > span {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    
    .el-icon {
      color: #409eff;
      font-size: 18px;
    }
  }
}

:deep(.el-card__body) {
  padding: 24px;
}

:deep(.el-divider__text) {
  color: #909399;
  font-weight: 500;
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
  line-height: 1.4;
}

// 响应式设计
@media (max-height: 800px) {
  .profile-container {
    height: calc(100vh - 60px);
  }
}

.logo-upload-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .logo-preview {
    display: flex;
    justify-content: center;
    margin-bottom: 8px;

    .logo-preview-image {
      width: 80px;
      height: 80px;
      object-fit: contain;
      border: 1px solid #dcdfe6;
      border-radius: 6px;
      background-color: #f8f9fa;
      padding: 8px;
    }
  }

  .logo-url-input {
    margin-top: 8px;
  }
}

@media (max-width: 768px) {
  .profile-container {
    padding: 16px;
    height: calc(100vh - 60px);
  }

  .profile-form {
    max-width: 100%;
  }

  .profile-card {
    margin-bottom: 16px;
  }

  .logo-upload-wrapper {
    gap: 12px;

    .logo-preview {
      margin-bottom: 6px;

      .logo-preview-image {
        width: 60px;
        height: 60px;
        padding: 6px;
      }
    }

    .logo-url-input {
      margin-top: 6px;
    }
  }
}
</style>