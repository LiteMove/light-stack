<template>
  <el-dialog
    v-model="dialogVisible"
    :title="'{{.FunctionName}}详情'"
    width="800px"
    :before-close="handleClose"
    destroy-on-close
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <div v-loading="loading" element-loading-text="加载中...">
      <div class="detail-content">
        {{- $hasGroups := false }}
        {{- range .Fields }}
        {{- if or (contains .ColumnComment "基本") (contains .ColumnComment "信息") }}
        {{- if not $hasGroups }}
        {{- $hasGroups = true }}
        <!-- 基本信息 -->
        <el-divider content-position="left">
          <el-icon><Document /></el-icon>
          基本信息
        </el-divider>
        {{- end }}
        {{- break }}
        {{- end }}
        {{- end }}

        <div class="detail-grid">
          {{- range .Fields }}
          {{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant") (contains .ColumnName "created_at") (contains .ColumnName "updated_at") (contains .ColumnName "deleted_at")) }}
          <div class="detail-item">
            <label class="detail-label">{{.ColumnComment}}:</label>
            <div class="detail-value">
              {{- if eq .HtmlType "switch" }}
              <el-tag
                :type="detailData?.{{generateJSField .ColumnName}} ? 'success' : 'danger'"
                size="small"
              >
                {{`{{ detailData?.`}}{{generateJSField .ColumnName}}{{` ? '是' : '否' }}`}}
              </el-tag>
              {{- else if eq .HtmlType "image" }}
              <el-image
                v-if="detailData?.{{generateJSField .ColumnName}}"
                :src="detailData.{{generateJSField .ColumnName}}"
                :preview-src-list="[detailData.{{generateJSField .ColumnName}}]"
                fit="cover"
                style="width: 100px; height: 100px; border-radius: 4px"
                preview-teleported
              />
              <span v-else class="text-gray-400">无图片</span>
              {{- else if eq .HtmlType "file" }}
              <el-link
                v-if="detailData?.{{generateJSField .ColumnName}}"
                :href="detailData.{{generateJSField .ColumnName}}"
                target="_blank"
                type="primary"
              >
                <el-icon><Download /></el-icon>
                下载文件
              </el-link>
              <span v-else class="text-gray-400">无文件</span>
              {{- else if contains .ColumnName "status" }}
              <el-tag
                :type="getStatusTagType(detailData?.{{generateJSField .ColumnName}})"
                size="small"
              >
                {{`{{ getStatusText(detailData?.`}}{{generateJSField .ColumnName}}{{`) }}`}}
              </el-tag>
              {{- else if eq .HtmlType "datetime" }}
              <span>{{`{{ formatDateTime(detailData?.`}}{{generateJSField .ColumnName}}{{`) }}`}}</span>
              {{- else if eq .HtmlType "date" }}
              <span>{{`{{ formatDate(detailData?.`}}{{generateJSField .ColumnName}}{{`) }}`}}</span>
              {{- else if eq .HtmlType "time" }}
              <span>{{`{{ formatTime(detailData?.`}}{{generateJSField .ColumnName}}{{`) }}`}}</span>
              {{- else if or (eq .HtmlType "textarea") (contains .ColumnName "content") (contains .ColumnName "description") }}
              <div class="detail-text-content">
                {{`{{ detailData?.`}}{{generateJSField .ColumnName}}{{` || '-' }}`}}
              </div>
              {{- else if eq .HtmlType "url" }}
              <el-link
                v-if="detailData?.{{generateJSField .ColumnName}}"
                :href="detailData.{{generateJSField .ColumnName}}"
                target="_blank"
                type="primary"
              >
                {{`{{ detailData.`}}{{generateJSField .ColumnName}}{{` }}`}}
              </el-link>
              <span v-else class="text-gray-400">-</span>
              {{- else if eq .HtmlType "email" }}
              <el-link
                v-if="detailData?.{{generateJSField .ColumnName}}"
                :href="`mailto:${detailData.{{generateJSField .ColumnName}}}`"
                type="primary"
              >
                {{`{{ detailData.`}}{{generateJSField .ColumnName}}{{` }}`}}
              </el-link>
              <span v-else class="text-gray-400">-</span>
              {{- else if contains .ColumnName "phone" }}
              <el-link
                v-if="detailData?.{{generateJSField .ColumnName}}"
                :href="`tel:${detailData.{{generateJSField .ColumnName}}}`"
                type="primary"
              >
                {{`{{ detailData.`}}{{generateJSField .ColumnName}}{{` }}`}}
              </el-link>
              <span v-else class="text-gray-400">-</span>
              {{- else }}
              <span>{{`{{ detailData?.`}}{{generateJSField .ColumnName}}{{` || '-' }}`}}</span>
              {{- end }}
            </div>
          </div>
          {{- end }}
          {{- end }}
        </div>

        <!-- 系统信息 -->
        <el-divider content-position="left">
          <el-icon><Clock /></el-icon>
          系统信息
        </el-divider>

        <div class="detail-grid">
          {{- range .Fields }}
          {{- if or (contains .ColumnName "created_at") (contains .ColumnName "updated_at") }}
          <div class="detail-item">
            <label class="detail-label">
              {{- if contains .ColumnName "created_at" }}创建时间{{- else if contains .ColumnName "updated_at" }}更新时间{{- else }}{{.ColumnComment}}{{- end }}:
            </label>
            <div class="detail-value">
              <span>{{`{{ formatDateTime(detailData?.`}}{{generateJSField .ColumnName}}{{`) }}`}}</span>
            </div>
          </div>
          {{- end }}
          {{- end }}
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button
          type="warning"
          @click="handleEdit"
          v-if="checkPermission('{{generatePermission .ModuleName .BusinessName "edit"}}')"
        >
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Document,
  Clock,
  Download,
  Edit
} from '@element-plus/icons-vue'
import type { {{.ClassName}} } from '@/api/types'
import { usePermission } from '@/composables/usePermission'

interface Props {
  visible: boolean
  detailData: {{.ClassName}} | null
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'edit', data: {{.ClassName}}): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 权限检查
const { checkPermission } = usePermission()

const loading = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 方法
const handleClose = () => {
  emit('update:visible', false)
}

const handleEdit = () => {
  if (props.detailData) {
    emit('edit', props.detailData)
    handleClose()
  }
}

{{- $hasStatusField := false }}
{{- range .Fields }}
{{- if contains .ColumnName "status" }}
{{- $hasStatusField = true }}
// 状态相关方法
const getStatusTagType = (status: any) => {
  switch (status) {
    case 1:
    case true:
      return 'success'
    case 0:
    case false:
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: any) => {
  switch (status) {
    case 1:
    case true:
      return '启用'
    case 0:
    case false:
      return '禁用'
    default:
      return '未知'
  }
}
{{- break }}
{{- end }}
{{- end }}

// 格式化方法
const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString('zh-CN')
}

const formatDate = (date: string | undefined) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const formatTime = (time: string | undefined) => {
  if (!time) return '-'
  return new Date(`1970-01-01 ${time}`).toLocaleTimeString('zh-CN')
}
</script>

<style lang="scss" scoped>
.detail-content {
  max-height: 500px;
  overflow-y: auto;
  padding: 0 8px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px 24px;
  margin-bottom: 24px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  .detail-label {
    font-size: 14px;
    font-weight: 500;
    color: #606266;
    margin-bottom: 4px;
  }

  .detail-value {
    font-size: 14px;
    color: #303133;
    word-break: break-all;
    min-height: 22px;
    display: flex;
    align-items: flex-start;
  }
}

.detail-text-content {
  line-height: 1.6;
  white-space: pre-wrap;
  max-height: 100px;
  overflow-y: auto;
  padding: 8px 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 分组标题样式
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

// 标签样式
:deep(.el-tag) {
  border-radius: 4px;
}

// 图片样式
:deep(.el-image) {
  border: 1px solid #dcdfe6;
}

// 链接样式
:deep(.el-link) {
  font-size: 14px;
}

// 灰色文本
.text-gray-400 {
  color: #909399;
}

// 滚动条样式
.detail-content::-webkit-scrollbar {
  width: 6px;
}

.detail-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.detail-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.detail-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>