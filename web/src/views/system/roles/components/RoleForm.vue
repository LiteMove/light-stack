<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑角色' : '新增角色'"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="role-form"
    >
      <el-form-item label="角色名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入角色名称"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="角色编码" prop="code">
        <el-input
          v-model="form.code"
          placeholder="请输入角色编码，如：admin, user"
          maxlength="50"
          show-word-limit
          :disabled="isEdit && form.isSystem"
        />
        <div class="form-tip">
          角色编码用于系统权限判断，建议使用英文字母和下划线
        </div>
      </el-form-item>

      <el-form-item label="角色描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          placeholder="请输入角色描述"
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

      <el-form-item v-if="isEdit && form.isSystem" label="系统角色">
        <el-tag type="danger" effect="plain">
          系统内置角色，部分字段不可修改
        </el-tag>
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
import { roleApi } from '@/api'
import type { Role, RoleFormData } from '@/api/types'

// Props
interface Props {
  visible: boolean
  formData?: Partial<Role>
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
const form = reactive<RoleFormData>({
  name: '',
  code: '',
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

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 50, message: '角色名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { min: 2, max: 50, message: '角色编码长度在 2 到 50 个字符', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, 
      message: '角色编码必须以字母开头，只能包含字母、数字和下划线', 
      trigger: 'blur' 
    }
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

// 监听表单数据变化
watch(
  () => props.formData,
  (newData) => {
    if (newData && Object.keys(newData).length > 0) {
      Object.assign(form, {
        name: newData.name || '',
        code: newData.code || '',
        description: newData.description || '',
        status: newData.status || 1,
        sortOrder: newData.sortOrder || 100
      })
    } else {
      // 重置表单
      Object.assign(form, {
        name: '',
        code: '',
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

    if (isEdit.value && props.formData?.id) {
      // 编辑角色
      await roleApi.updateRole(props.formData.id, form)
      ElMessage.success('角色更新成功')
    } else {
      // 创建角色
      await roleApi.createRole(form)
      ElMessage.success('角色创建成功')
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    // 错误信息已在响应拦截器中处理
    console.error(isEdit.value ? '更新角色失败:' : '创建角色失败:', error)
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
.role-form {
  .form-tip {
    margin-top: 4px;
    font-size: 12px;
    color: #909399;
    line-height: 1.4;
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
      .el-input-number {
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

  .role-form {
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
            }
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