<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑权限' : '新增权限'"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="permission-form"
    >
      <el-form-item label="权限名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入权限名称"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="权限编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入权限编码，如：user:list, user:create"
          maxlength="100"
          show-word-limit
          :disabled="isEdit && form.isSystem"
        />
        <div class="form-tip">
          权限编码用于系统权限判断，建议使用英文字母、数字和冒号、下划线
        </div>
      </el-form-item>

      <el-form-item label="权限类型" prop="type">
        <el-radio-group v-model="form.type" @change="handleTypeChange">
          <el-radio label="api">API接口</el-radio>
          <el-radio label="page">页面权限</el-radio>
          <el-radio label="button">按钮权限</el-radio>
          <el-radio label="data">数据权限</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="资源路径" prop="resource" v-if="form.type === 'api' || form.type === 'page'">
        <el-input
          v-model="form.resource"
          :placeholder="form.type === 'api' ? '请输入API路径，如：/api/v1/users' : '请输入页面路径，如：/system/users'"
          maxlength="200"
          show-word-limit
        />
        <div class="form-tip">
          {{ form.type === 'api' ? 'API接口的完整路径' : '页面的路由路径' }}
        </div>
      </el-form-item>

      <el-form-item label="操作动作" prop="action" v-if="form.type === 'api' || form.type === 'button'">
        <el-select
          v-model="form.action"
          :placeholder="form.type === 'api' ? '选择HTTP方法' : '选择按钮操作'"
          style="width: 100%"
          clearable
        >
          <template v-if="form.type === 'api'">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PATCH" value="PATCH" />
          </template>
          <template v-else-if="form.type === 'button'">
            <el-option label="查看" value="view" />
            <el-option label="新增" value="create" />
            <el-option label="编辑" value="edit" />
            <el-option label="删除" value="delete" />
            <el-option label="导入" value="import" />
            <el-option label="导出" value="export" />
            <el-option label="审核" value="audit" />
            <el-option label="发布" value="publish" />
          </template>
        </el-select>
        <div class="form-tip">
          {{ form.type === 'api' ? '指定API接口的HTTP请求方法' : '指定按钮的具体操作类型' }}
        </div>
      </el-form-item>

      <el-form-item label="权限描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          placeholder="请输入权限描述"
          :rows="3"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="排序" prop="sortOrder">
        <el-input-number
          v-model="form.sortOrder"
          :min="0"
          :max="9999"
          controls-position="right"
          placeholder="排序值"
          style="width: 200px"
        />
        <div class="form-tip">
          数值越小排序越靠前
        </div>
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio :label="1">启用</el-radio>
          <el-radio :label="2">禁用</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item v-if="isEdit && form.isSystem" label="系统权限">
        <el-tag type="danger" effect="plain">
          系统内置权限，部分字段不可修改
        </el-tag>
      </el-form-item>

      <!-- 权限示例 -->
      <el-form-item label="编码示例">
        <div class="code-examples">
          <div class="example-section">
            <div class="example-title">权限编码示例：</div>
            <div class="example-list">
              <el-tag size="small" class="example-tag" v-for="example in codeExamples" :key="example">
                {{ example }}
              </el-tag>
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="submitting">
          取消
        </el-button>
        <el-button 
          type="primary" 
          @click="handleSubmit"
          :loading="submitting"
        >
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { permissionApi } from '@/api'
import type { Permission, PermissionFormData } from '@/api/types'

// Props
interface Props {
  visible: boolean
  formData?: Partial<Permission>
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  formData: () => ({})
})

// Emits
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

// 响应式数据
const formRef = ref<FormInstance>()
const submitting = ref(false)

// 表单数据
const form = reactive<PermissionFormData>({
  name: '',
  code: '',
  type: 'api',
  resource: '',
  action: '',
  description: '',
  status: 1,
  sortOrder: 100
})

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!(props.formData && props.formData.id))

// 权限编码示例
const codeExamples = computed(() => {
  const examples: Record<string, string[]> = {
    api: ['user:list', 'user:create', 'user:update', 'user:delete', 'role:assign'],
    page: ['system:user', 'system:role', 'dashboard:view', 'report:analysis'],
    button: ['user:btn:edit', 'user:btn:delete', 'role:btn:assign', 'data:btn:export'],
    data: ['user:data:self', 'user:data:dept', 'user:data:all', 'report:data:sensitive']
  }
  return examples[form.type] || []
})

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' },
    { min: 2, max: 50, message: '权限名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入权限编码', trigger: 'blur' },
    { min: 2, max: 100, message: '权限编码长度在 2 到 100 个字符', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z][a-zA-Z0-9:_]*$/, 
      message: '权限编码必须以字母开头，只能包含字母、数字、冒号和下划线', 
      trigger: 'blur' 
    }
  ],
  type: [
    { required: true, message: '请选择权限类型', trigger: 'change' }
  ],
  resource: [
    { max: 200, message: '资源路径不能超过 200 个字符', trigger: 'blur' }
  ],
  action: [
    { max: 50, message: '操作动作不能超过 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述不能超过 200 个字符', trigger: 'blur' }
  ],
  sortOrder: [
    { required: true, message: '请输入排序值', trigger: 'blur' },
    { type: 'number', min: 0, max: 9999, message: '排序值范围 0-9999', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 权限类型变化处理
const handleTypeChange = (type: string) => {
  // 清空相关字段
  form.resource = ''
  form.action = ''
  
  // 清除验证错误
  nextTick(() => {
    formRef.value?.clearValidate(['resource', 'action'])
  })
}

// 监听表单数据变化
watch(
  () => props.formData,
  (newData) => {
    if (newData && Object.keys(newData).length > 0) {
      Object.assign(form, {
        name: newData.name || '',
        code: newData.code || '',
        type: newData.type || 'api',
        resource: newData.resource || '',
        action: newData.action || '',
        description: newData.description || '',
        status: newData.status || 1,
        sortOrder: newData.sortOrder || 100
      })
    } else {
      // 重置表单
      Object.assign(form, {
        name: '',
        code: '',
        type: 'api',
        resource: '',
        action: '',
        description: '',
        status: 1,
        sortOrder: 100
      })
    }
  },
  { immediate: true, deep: true }
)

// 监听对话框显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      nextTick(() => {
        formRef.value?.clearValidate()
      })
    }
  }
)

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    const valid = await formRef.value.validate()
    if (!valid) return

    submitting.value = true

    // 清理不需要的字段
    const submitData = { ...form }
    if (form.type !== 'api' && form.type !== 'page') {
      submitData.resource = ''
    }
    if (form.type !== 'api' && form.type !== 'button') {
      submitData.action = ''
    }

    if (isEdit.value && props.formData?.id) {
      // 编辑权限
      await permissionApi.updatePermission(props.formData.id, submitData)
      ElMessage.success('权限更新成功')
    } else {
      // 创建权限
      await permissionApi.createPermission(submitData)
      ElMessage.success('权限创建成功')
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
  } finally {
    submitting.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  emit('update:visible', false)
  // 清空表单
  nextTick(() => {
    formRef.value?.resetFields()
  })
}
</script>

<style lang="scss" scoped>
.permission-form {
  .form-tip {
    margin-top: 4px;
    font-size: 12px;
    color: #909399;
    line-height: 1.4;
  }

  .code-examples {
    .example-section {
      .example-title {
        font-size: 13px;
        color: #606266;
        margin-bottom: 8px;
        font-weight: 500;
      }

      .example-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;

        .example-tag {
          margin: 0;
          font-family: 'Monaco', 'Consolas', monospace;
          font-size: 12px;
          cursor: pointer;
          transition: all 0.2s ease;

          &:hover {
            transform: translateY(-1px);
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          }
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// Element Plus Dialog 样式覆盖
:deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;

  .el-dialog__header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 20px 24px;
    border-bottom: none;

    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
    }

    .el-dialog__headerbtn {
      .el-dialog__close {
        color: white;
        font-size: 20px;

        &:hover {
          color: rgba(255, 255, 255, 0.8);
        }
      }
    }
  }

  .el-dialog__body {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;
  }

  .el-dialog__footer {
    padding: 20px 24px;
    background: #fafafa;
    border-top: 1px solid #e4e7ed;
  }
}

// 表单样式优化
:deep(.el-form) {
  .el-form-item {
    margin-bottom: 24px;

    .el-form-item__label {
      color: #606266;
      font-weight: 500;
      line-height: 1.5;
    }

    .el-form-item__content {
      .el-input,
      .el-textarea,
      .el-input-number,
      .el-select {
        .el-input__inner,
        .el-textarea__inner {
          border-radius: 6px;
          border: 1px solid #dcdfe6;
          transition: all 0.2s ease;

          &:hover {
            border-color: #c0c4cc;
          }

          &:focus {
            border-color: #409eff;
            box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
          }
        }
      }

      .el-radio-group {
        .el-radio {
          margin-right: 24px;
          margin-bottom: 8px;

          .el-radio__label {
            color: #606266;
            font-weight: 500;
          }

          &.is-checked {
            .el-radio__label {
              color: #409eff;
            }
          }
        }
      }

      .el-tag {
        padding: 8px 12px;
        border-radius: 6px;
        font-size: 13px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 95% !important;
    margin: 20px auto;

    .el-dialog__header {
      padding: 16px 20px;

      .el-dialog__title {
        font-size: 16px;
      }
    }

    .el-dialog__body {
      padding: 20px;
    }

    .el-dialog__footer {
      padding: 16px 20px;
    }
  }

  .permission-form {
    :deep(.el-form) {
      .el-form-item {
        margin-bottom: 20px;

        .el-form-item__label {
          margin-bottom: 8px;
        }

        .el-form-item__content {
          .el-radio-group {
            .el-radio {
              margin-right: 16px;
              margin-bottom: 8px;
              display: block;
            }
          }
        }
      }
    }

    .code-examples {
      .example-section {
        .example-list {
          .example-tag {
            margin-bottom: 4px;
          }
        }
      }
    }
  }

  .dialog-footer {
    flex-direction: column;

    .el-button {
      width: 100%;
      margin: 0;

      &:first-child {
        margin-bottom: 8px;
      }
    }
  }
}
</style>