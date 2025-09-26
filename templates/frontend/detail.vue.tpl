<template>
  <el-dialog
    v-model="dialogVisible"
    title="{{.FunctionName}}详情"
    width="800px"
    :before-close="handleClose"
    destroy-on-close
  >
    <div v-loading="loading" element-loading-text="加载中...">
      <el-descriptions
        class="detail-descriptions"
        :column="2"
        size="large"
        border
      >
{{- range .Fields }}
{{- if not .IsPk }}
        <el-descriptions-item label="{{.ColumnComment}}">
{{- if eq .HtmlType "image" }}
          <el-image
            v-if="data?.{{generateJSField .ColumnName}}"
            :src="data.{{generateJSField .ColumnName}}"
            :preview-src-list="[data.{{generateJSField .ColumnName}}]"
            fit="cover"
            style="width: 60px; height: 60px; border-radius: 4px;"
          />
          <span v-else class="empty-value">-</span>
{{- else if eq .HtmlType "select" }}
          <!-- TODO: 根据字典类型{{.DictType}}显示对应文本 -->
          <el-tag :type="getStatusType(data?.{{generateJSField .ColumnName}})" size="small">
            {{`{{ getStatusText(data?.`}}{{generateJSField .ColumnName}}{{`) }}`}}
          </el-tag>
{{- else if eq .GoType "bool" }}
          <el-tag :type="data?.{{generateJSField .ColumnName}} ? 'success' : 'info'" size="small">
            {{`{{ data?.`}}{{generateJSField .ColumnName}}{{` ? '是' : '否' }}`}}
          </el-tag>
{{- else if eq .HtmlType "switch" }}
          <el-tag :type="data?.{{generateJSField .ColumnName}} ? 'success' : 'info'" size="small">
            {{`{{ data?.`}}{{generateJSField .ColumnName}}{{` ? '启用' : '禁用' }}`}}
          </el-tag>
{{- else if contains .ColumnName "time" }}
          <span class="time-text">{{`{{ formatDateTime(data?.`}}{{generateJSField .ColumnName}}{{`) }}`}}</span>
{{- else if eq .HtmlType "url" }}
          <el-link
            v-if="data?.{{generateJSField .ColumnName}}"
            :href="data.{{generateJSField .ColumnName}}"
            target="_blank"
            type="primary"
          >
            {{`{{ data.`}}{{generateJSField .ColumnName}}{{` }}`}}
          </el-link>
          <span v-else class="empty-value">-</span>
{{- else if eq .HtmlType "file" }}
          <el-link
            v-if="data?.{{generateJSField .ColumnName}}"
            :href="data.{{generateJSField .ColumnName}}"
            target="_blank"
            type="primary"
            :icon="Document"
          >
            查看文件
          </el-link>
          <span v-else class="empty-value">-</span>
{{- else if eq .HtmlType "textarea" }}
          <div class="textarea-content">
            {{`{{ data?.`}}{{generateJSField .ColumnName}}{{` || '-' }}`}}
          </div>
{{- else }}
          {{`{{ data?.`}}{{generateJSField .ColumnName}}{{` || '-' }}`}}
{{- end }}
        </el-descriptions-item>
{{- end }}
{{- end }}

{{- $hasPkField := false }}
{{- range .Fields }}
{{- if .IsPk }}
{{- $hasPkField = true }}
        <!-- 系统信息 -->
        <el-descriptions-item label="{{if eq .ColumnName "id"}}ID{{else}}{{.ColumnComment}}{{end}}" span="2">
          <el-tag type="info" size="small">{{`{{ data?.`}}{{generateJSField .ColumnName}}{{` }}`}}</el-tag>
        </el-descriptions-item>
{{- break }}
{{- end }}
{{- end }}

        <!-- 创建时间 -->
        <el-descriptions-item label="创建时间">
          <span class="time-text">{{`{{ formatDateTime(data?.createdAt) }}`}}</span>
        </el-descriptions-item>

        <!-- 更新时间 -->
        <el-descriptions-item label="更新时间">
          <span class="time-text">{{`{{ formatDateTime(data?.updatedAt) }}`}}</span>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 操作按钮 -->
      <div class="detail-actions">
        <el-button
          v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "update"}}')"
          type="primary"
          :icon="Edit"
          @click="handleEdit"
        >
          编辑{{.FunctionName}}
        </el-button>
        <el-button
          v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "delete"}}')"
          type="danger"
          :icon="Delete"
          @click="handleDelete"
        >
          删除{{.FunctionName}}
        </el-button>
        <el-button @click="handleClose">
          关闭
        </el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Edit,
  Delete,
  Document
} from '@element-plus/icons-vue'
import { {{toLower .BusinessName}}Api } from '@/api'
import type { {{.ClassName}} } from '@/api/types'
import { formatDateTime } from '@/utils/date'

interface Props {
  visible: boolean
  data: {{.ClassName}} | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'edit', data: {{.ClassName}}): void
  (e: 'delete', data: {{.ClassName}}): void
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const loading = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 状态/类型转换函数
{{- $hasSelectField := false }}
{{- range .Fields }}
{{- if eq .HtmlType "select" }}
{{- if not $hasSelectField }}
{{- $hasSelectField = true }}
const getStatusType = (value: any): string => {
  // TODO: 根据具体业务逻辑返回对应的标签类型
  switch (value) {
    case 1:
      return 'success'
    case 2:
      return 'warning'
    case 0:
      return 'info'
    default:
      return 'primary'
  }
}

const getStatusText = (value: any): string => {
  // TODO: 根据字典类型{{.DictType}}返回对应的文本
  switch (value) {
    case 1:
      return '启用'
    case 2:
      return '禁用'
    case 0:
      return '未知'
    default:
      return String(value)
  }
}
{{- end }}
{{- break }}
{{- end }}
{{- end }}

// 编辑{{.FunctionName}}
const handleEdit = () => {
  if (props.data) {
    emit('edit', props.data)
    handleClose()
  }
}

// 删除{{.FunctionName}}
const handleDelete = async () => {
  if (!props.data) return

  try {
    await ElMessageBox.confirm(
      `确定要删除{{.FunctionName}}吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    loading.value = true
    await {{toLower .BusinessName}}Api.delete{{.ClassName}}(props.data.{{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}})
    ElMessage.success('删除成功')
    emit('refresh')
    handleClose()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除{{.FunctionName}}失败:', error)
    }
  } finally {
    loading.value = false
  }
}

// 关闭弹窗
const handleClose = () => {
  emit('update:visible', false)
}
</script>

<style lang="scss" scoped>
.detail-descriptions {
  margin-bottom: 24px;

  :deep(.el-descriptions__body) {
    .el-descriptions__table {
      .el-descriptions__cell {
        &.is-bordered-label {
          background-color: #fafafa;
          font-weight: 500;
          color: #606266;
          width: 120px;
        }

        &.is-bordered-content {
          padding: 12px 16px;
        }
      }
    }
  }

  .empty-value {
    color: #c0c4cc;
    font-style: italic;
    font-size: 12px;
  }

  .time-text {
    font-size: 13px;
    color: #909399;
  }

  .textarea-content {
    max-width: 400px;
    word-wrap: break-word;
    white-space: pre-wrap;
    line-height: 1.5;
    font-size: 13px;
  }
}

.detail-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;

  .el-button {
    min-width: 100px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .detail-descriptions {
    :deep(.el-descriptions__body) {
      .el-descriptions__table {
        .el-descriptions__row {
          display: block;

          .el-descriptions__cell {
            display: block !important;
            width: 100% !important;

            &.is-bordered-label {
              border-right: none !important;
              border-bottom: none !important;
              padding-bottom: 8px;
            }

            &.is-bordered-content {
              border-top: none !important;
              padding-top: 8px;
              padding-bottom: 16px;
            }
          }
        }
      }
    }
  }

  .detail-actions {
    flex-direction: column;
    align-items: stretch;

    .el-button {
      margin: 0;
    }
  }
}
</style>