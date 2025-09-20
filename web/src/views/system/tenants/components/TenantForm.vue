<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑租户' : '新建租户'"
    width="600px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
      label-position="right"
    >
      <el-form-item label="租户名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入租户名称"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="域名" prop="domain">
        <el-input
          v-model="formData.domain"
          placeholder="请输入租户域名，如：company.example.com"
          maxlength="100"
          show-word-limit
        >
          <template #prepend>https://</template>
        </el-input>
        <div class="form-tip">
          域名用于租户访问，必须是有效的域名格式
        </div>
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="formData.status">
          <el-radio :label="1">启用</el-radio>
          <el-radio :label="2">禁用</el-radio>
          <el-radio :label="3">试用</el-radio>
          <el-radio :label="4">过期</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="过期时间" prop="expired_at">
        <el-date-picker
          v-model="formData.expired_at"
          type="datetime"
          placeholder="选择过期时间"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DD HH:mm:ss"
          style="width: 100%"
        />
        <div class="form-tip">
          不设置过期时间表示永久有效
        </div>
      </el-form-item>

      <el-form-item label="配置信息" prop="config">
        <el-input
          v-model="formData.config"
          type="textarea"
          :rows="4"
          placeholder="请输入JSON格式的配置信息，如：{&quot;title&quot;:&quot;公司名称&quot;,&quot;logo&quot;:&quot;&quot;,&quot;theme&quot;:&quot;#1890ff&quot;}"
          maxlength="1000"
          show-word-limit
        />
        <div class="form-tip">
          配置信息用于定制租户的显示样式，格式为JSON
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { tenantApi } from '@/api/tenant'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  tenantData: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['update:visible', 'success'])

// 响应式数据
const formRef = ref()
const submitting = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => {
  return props.tenantData && props.tenantData.id
})

// 表单数据
const formData = reactive({
  name: '',
  domain: '',
  status: 1,
  expired_at: '',
  config: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入租户名称', trigger: 'blur' },
    { min: 1, max: 100, message: '租户名称长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  domain: [
    { max: 100, message: '域名长度不能超过 100 个字符', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/,
      message: '请输入有效的域名格式',
      trigger: 'blur'
    }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  config: [
    {
      validator: validateConfig,
      trigger: 'blur'
    }
  ]
}

// 验证配置信息
function validateConfig(rule, value, callback) {
  if (!value) {
    callback()
    return
  }

  try {
    JSON.parse(value)
    callback()
  } catch (error) {
    callback(new Error('配置信息必须是有效的JSON格式'))
  }
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(formData, {
    name: '',
    domain: '',
    status: 1,
    expired_at: '',
    config: ''
  })
}

// 初始化表单数据
const initFormData = () => {
  if (isEdit.value && props.tenantData) {
    Object.assign(formData, {
      name: props.tenantData.name || '',
      domain: props.tenantData.domain || '',
      status: props.tenantData.status || 1,
      expired_at: props.tenantData.expired_at || '',
      config: props.tenantData.config || ''
    })
  } else {
    resetForm()
  }
}

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible) {
    initFormData()
  }
})

// 关闭弹窗
const handleClose = () => {
  resetForm()
  emit('update:visible', false)
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    const valid = await formRef.value.validate()
    if (!valid) return

    submitting.value = true

    const submitData = { ...formData }

    // 如果没有设置过期时间，则移除该字段
    if (!submitData.expired_at) {
      delete submitData.expired_at
    }

    // 如果没有设置配置信息，则移除该字段
    if (!submitData.config) {
      delete submitData.config
    }

    if (isEdit.value) {
      // 更新租户
      await tenantApi.updateTenant(props.tenantData.id, submitData)
      ElMessage.success('更新租户成功')
    } else {
      // 创建租户
      await tenantApi.createTenant(submitData)
      ElMessage.success('创建租户成功')
    }

    emit('success')
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error((isEdit.value ? '更新' : '创建') + '租户失败:', error)
  } finally {
    submitting.value = false
  }
}

// 检查域名可用性
const checkDomainAvailability = async () => {
  if (!formData.domain) return

  try {
    const response = await tenantApi.checkDomain(formData.domain)
    if (!response.data.available) {
      ElMessage.warning('该域名已被使用')
      return false
    }
    return true
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('检查域名失败:', error)
    return false
  }
}

// 检查租户名称可用性
const checkNameAvailability = async () => {
  if (!formData.name) return

  try {
    const response = await tenantApi.checkName(formData.name)
    if (!response.data.available) {
      ElMessage.warning('该租户名称已被使用')
      return false
    }
    return true
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('检查租户名称失败:', error)
    return false
  }
}
</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-input-group__prepend) {
  background-color: #f5f7fa;
  color: #909399;
  border-color: #dcdfe6;
}
</style>