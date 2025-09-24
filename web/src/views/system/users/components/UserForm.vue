<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑用户' : '新建用户'"
    width="800px"
    :before-close="handleClose"
    destroy-on-close
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      label-position="right"
      v-loading="loading"
      element-loading-text="处理中..."
    >
      <!-- 基本信息 -->
      <el-divider content-position="left">
        <el-icon><User /></el-icon>
        基本信息
      </el-divider>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              maxlength="50"
              show-word-limit
              :disabled="isEdit"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="昵称" prop="nickname">
            <el-input
              v-model="form.nickname"
              placeholder="请输入用户昵称"
              maxlength="50"
              show-word-limit
            >
              <template #prefix>
                <el-icon><UserFilled /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20" v-if="!isEdit">
        <el-col :span="12">
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入登录密码"
              maxlength="50"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="请确认登录密码"
              maxlength="50"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <!-- 联系信息 -->
      <el-divider content-position="left">
        <el-icon><Message /></el-icon>
        联系信息
      </el-divider>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="邮箱" prop="email">
            <el-input
              v-model="form.email"
              placeholder="请输入邮箱地址"
              maxlength="100"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="手机号" prop="phone">
            <el-input
              v-model="form.phone"
              placeholder="请输入手机号码"
              maxlength="11"
            >
              <template #prefix>
                <el-icon><Phone /></el-icon>
              </template>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <!-- 头像上传 -->
      <el-form-item label="用户头像" prop="avatar">
        <AvatarUpload
          v-model="form.avatar"
          @success="handleAvatarSuccess"
          @error="handleAvatarError"
        />
      </el-form-item>

      <!-- 状态设置 -->
      <el-divider content-position="left">
        <el-icon><Setting /></el-icon>
        状态设置
      </el-divider>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="用户状态" prop="status">
            <el-radio-group v-model="form.status" :disabled="isEdit && form.isSystem">
              <el-radio :label="1">
                <el-icon style="color: #67c23a"><Check /></el-icon>
                正常
              </el-radio>
              <el-radio :label="2">
                <el-icon style="color: #f56c6c"><Close /></el-icon>
                禁用
              </el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="系统用户">
            <el-switch
              v-model="form.isSystem"
              active-text="是"
              inactive-text="否"
              :disabled="isEdit && form.isSystem"
            />
            <div class="form-item-tip">
              系统用户不能被删除或禁用
            </div>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="loading">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  User,
  UserFilled,
  Lock,
  Message,
  Phone,
  Plus,
  Setting,
  Check,
  Close,
  Key
} from '@element-plus/icons-vue'
import AvatarUpload from '@/components/AvatarUpload.vue'
import { userApi } from '@/api'
import type { User as UserType, Role } from '@/api/types'

interface UserFormData {
  username: string
  nickname: string
  email: string
  phone?: string
  avatar?: string
  password?: string
  confirmPassword?: string
  status: number
  isSystem: boolean
}

interface Props {
  visible: boolean
  formData: Partial<UserType>
  roles: Role[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const loading = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.formData.id)

// 表单数据
const form = ref<UserFormData>({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  avatar: '',
  password: '',
  confirmPassword: '',
  status: 1,
  isSystem: false
})

// 头像上传成功处理
const handleAvatarSuccess = (file: any) => {
  form.value.avatar = file.accessUrl || file.filePath
  ElMessage.success('头像上传成功')
}

// 头像上传错误处理
const handleAvatarError = (error: string) => {
  ElMessage.error(error || '头像上传失败')
}
const rules = computed<FormRules>(() => {
  const baseRules = {
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' },
      { pattern: /^[a-zA-Z0-9_-]+$/, message: '只能包含字母、数字、下划线和横线', trigger: 'blur' }
    ],
    nickname: [
      { required: true, message: '请输入用户昵称', trigger: 'blur' },
      { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
    ],
    email: [
      { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
      { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
    ],
    phone: [
      { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
    ],
    status: [
      { required: true, message: '请选择用户状态', trigger: 'change' }
    ]
  }

  // 新建用户时需要密码验证
  if (!isEdit.value) {
    return {
      ...baseRules,
      password: [
        { required: true, message: '请输入登录密码', trigger: 'blur' },
        { min: 6, max: 50, message: '长度在 6 到 50 个字符', trigger: 'blur' },
        { 
          pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{6,}$/, 
          message: '密码必须包含大小写字母和数字', 
          trigger: 'blur' 
        }
      ],
      confirmPassword: [
        { required: true, message: '请确认登录密码', trigger: 'blur' },
        { validator: validateConfirmPassword, trigger: 'blur' }
      ]
    }
  }

  return baseRules
})

// 确认密码验证器
function validateConfirmPassword(rule: any, value: string, callback: any) {
  if (!value) {
    callback(new Error('请确认登录密码'))
  } else if (value !== form.value.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 监听表单数据变化
watch(
  () => props.formData,
  (newData) => {
    if (newData) {
      Object.assign(form.value, {
        username: newData.username || '',
        nickname: newData.nickname || '',
        email: newData.email || '',
        phone: newData.phone || '',
        avatar: newData.avatar || '',
        status: newData.status || 1,
        isSystem: newData.isSystem || false,
        password: '',
        confirmPassword: ''
      })
    }
  },
  { immediate: true, deep: true }
)

// 监听表单数据变化
const handleSubmit = async () => {
  try {
    const valid = await formRef.value?.validate()
    if (!valid) return

    loading.value = true

    const submitData = {
      username: form.value.username,
      nickname: form.value.nickname,
      email: form.value.email,
      phone: form.value.phone,
      avatar: form.value.avatar,
      status: form.value.status,
      isSystem: form.value.isSystem
    }

    if (!isEdit.value) {
      // 新建用户包含密码
      Object.assign(submitData, {
        password: form.value.password
      })
    }

    if (isEdit.value) {
      await userApi.updateUser(props.formData.id!, submitData)
      ElMessage.success('用户信息更新成功')
    } else {
      await userApi.createUser(submitData)
      ElMessage.success('用户创建成功')
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    // 错误信息已在响应拦截器中处理
    console.error(isEdit.value ? '更新用户失败:' : '创建用户失败:', error)
  } finally {
    loading.value = false
  }
}

// 关闭弹窗
const handleClose = () => {
  formRef.value?.resetFields()
  emit('update:visible', false)
}

// 生成用户名建议
const generateUsername = () => {
  if (!form.value.nickname) return
  
  // 简单的拼音转换逻辑
  const username = form.value.nickname
    .toLowerCase()
    .replace(/[^a-z0-9]/g, '')
    .substring(0, 20)
  
  if (!form.value.username && username) {
    form.value.username = username + Math.floor(Math.random() * 1000)
  }
}

// 监听昵称变化生成用户名
watch(() => form.value.nickname, generateUsername)
</script>

<style lang="scss" scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-item-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

// 角色显示
.current-roles {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;

  .role-tag {
    margin: 0;
  }
}

// 表单分组标题优化
:deep(.el-divider) {
  margin: 24px 0 16px 0;

  .el-divider__text {
    color: #409eff;
    font-weight: 500;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 6px;
  }
}

// 输入框前缀图标
:deep(.el-input__prefix) {
  color: #909399;
}

// 单选框样式优化
:deep(.el-radio) {
  margin-right: 24px;

  .el-radio__label {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}
</style>