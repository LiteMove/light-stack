<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`权限详情 - ${permissionData?.name}`"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="permission-detail" v-if="permissionData">
      <div class="detail-section">
        <div class="section-header">
          <el-icon class="section-icon"><Key /></el-icon>
          <span class="section-title">基本信息</span>
        </div>
        <div class="section-content">
          <div class="detail-grid">
            <div class="detail-item">
              <label class="item-label">权限名称:</label>
              <div class="item-value">
                {{ permissionData.name }}
                <el-tag v-if="permissionData.isSystem" type="danger" size="small" effect="plain" class="system-tag">
                  系统权限
                </el-tag>
              </div>
            </div>
            <div class="detail-item">
              <label class="item-label">权限编码:</label>
              <div class="item-value code-value">{{ permissionData.code }}</div>
            </div>
            <div class="detail-item">
              <label class="item-label">权限类型:</label>
              <div class="item-value">
                <el-tag
                  :type="getTypeTagType(permissionData.type)"
                  size="default"
                  effect="light"
                >
                  {{ getTypeLabel(permissionData.type) }}
                </el-tag>
              </div>
            </div>
            <div class="detail-item">
              <label class="item-label">状态:</label>
              <div class="item-value">
                <el-tag
                  :type="permissionData.status === 1 ? 'success' : 'danger'"
                  size="default"
                  effect="light"
                >
                  {{ permissionData.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </div>
            </div>
            <div class="detail-item">
              <label class="item-label">排序:</label>
              <div class="item-value">{{ permissionData.sortOrder }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="detail-section" v-if="permissionData.resource || permissionData.action">
        <div class="section-header">
          <el-icon class="section-icon"><Setting /></el-icon>
          <span class="section-title">资源配置</span>
        </div>
        <div class="section-content">
          <div class="detail-grid">
            <div class="detail-item" v-if="permissionData.resource">
              <label class="item-label">资源路径:</label>
              <div class="item-value resource-value">{{ permissionData.resource }}</div>
            </div>
            <div class="detail-item" v-if="permissionData.action">
              <label class="item-label">操作动作:</label>
              <div class="item-value">
                <el-tag type="info" size="default" effect="light">
                  {{ permissionData.action }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="detail-section" v-if="permissionData.description">
        <div class="section-header">
          <el-icon class="section-icon"><Document /></el-icon>
          <span class="section-title">权限描述</span>
        </div>
        <div class="section-content">
          <div class="description-content">
            {{ permissionData.description }}
          </div>
        </div>
      </div>

      <div class="detail-section">
        <div class="section-header">
          <el-icon class="section-icon"><Clock /></el-icon>
          <span class="section-title">时间信息</span>
        </div>
        <div class="section-content">
          <div class="detail-grid">
            <div class="detail-item">
              <label class="item-label">创建时间:</label>
              <div class="item-value">{{ formatDateTime(permissionData.createdAt) }}</div>
            </div>
            <div class="detail-item">
              <label class="item-label">更新时间:</label>
              <div class="item-value">{{ formatDateTime(permissionData.updatedAt) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 使用示例 -->
      <div class="detail-section">
        <div class="section-header">
          <el-icon class="section-icon"><Memo /></el-icon>
          <span class="section-title">使用示例</span>
        </div>
        <div class="section-content">
          <div class="usage-examples">
            <div class="example-item">
              <div class="example-title">前端权限检查</div>
              <div class="example-code">
                <code>hasPermission('{{ permissionData.code }}')</code>
              </div>
            </div>
            <div class="example-item" v-if="permissionData.type === 'api'">
              <div class="example-title">API接口调用</div>
              <div class="example-code">
                <code>{{ permissionData.action || 'GET' }} {{ permissionData.resource || '/api/example' }}</code>
              </div>
            </div>
            <div class="example-item" v-if="permissionData.type === 'button'">
              <div class="example-title">按钮权限控制</div>
              <div class="example-code">
                <code>v-permission="'{{ permissionData.code }}'"</code>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- JSON数据 -->
      <div class="detail-section">
        <div class="section-header">
          <el-icon class="section-icon"><DataLine /></el-icon>
          <span class="section-title">JSON数据</span>
          <el-button 
            type="text" 
            size="small" 
            @click="copyToClipboard"
            :icon="DocumentCopy"
            class="copy-btn"
          >
            复制
          </el-button>
        </div>
        <div class="section-content">
          <div class="json-viewer">
            <pre><code>{{ formatJson(permissionData) }}</code></pre>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">
          关闭
        </el-button>
        <el-button type="primary" @click="handleEdit" v-if="!permissionData?.isSystem">
          编辑权限
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Key, 
  Setting, 
  Document, 
  Clock, 
  Memo,
  DataLine,
  DocumentCopy
} from '@element-plus/icons-vue'
import type { Permission } from '@/api/types'
import { formatDateTime } from '@/utils/date'

// Props
interface Props {
  visible: boolean
  permissionData?: Permission | null
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  permissionData: null
})

// Emits
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'edit', permission: Permission): void
}>()

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 获取类型标签
const getTypeLabel = (type: string): string => {
  const typeMap: Record<string, string> = {
    api: 'API接口',
    page: '页面权限',
    button: '按钮权限',
    data: '数据权限'
  }
  return typeMap[type] || type
}

// 获取类型标签样式
const getTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    api: 'primary',
    page: 'success',
    button: 'warning',
    data: 'info'
  }
  return typeMap[type] || 'default'
}

// 格式化JSON
const formatJson = (data: any): string => {
  try {
    return JSON.stringify(data, null, 2)
  } catch (error) {
    return String(data)
  }
}

// 复制到剪贴板
const copyToClipboard = async () => {
  if (!props.permissionData) return

  try {
    const jsonText = formatJson(props.permissionData)
    await navigator.clipboard.writeText(jsonText)
    ElMessage.success('JSON数据已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败，请手动选择复制')
  }
}

// 编辑权限
const handleEdit = () => {
  if (props.permissionData) {
    emit('edit', props.permissionData)
    handleClose()
  }
}

// 关闭对话框
const handleClose = () => {
  emit('update:visible', false)
}
</script>

<style lang="scss" scoped>
.permission-detail {
  .detail-section {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }

    .section-header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 16px;
      padding-bottom: 8px;
      border-bottom: 1px solid #e4e7ed;

      .section-icon {
        color: #409eff;
        font-size: 16px;
      }

      .section-title {
        font-weight: 600;
        color: #303133;
        font-size: 16px;
      }

      .copy-btn {
        margin-left: auto;
        font-size: 12px;
        padding: 4px 8px;
      }
    }

    .section-content {
      .detail-grid {
        display: grid;
        grid-template-columns: 1fr;
        gap: 16px;

        .detail-item {
          display: flex;
          align-items: flex-start;
          gap: 12px;

          .item-label {
            min-width: 80px;
            font-weight: 500;
            color: #606266;
            font-size: 14px;
            flex-shrink: 0;
          }

          .item-value {
            flex: 1;
            color: #303133;
            font-size: 14px;
            word-break: break-all;

            &.code-value {
              font-family: 'Monaco', 'Consolas', monospace;
              background: #f5f7fa;
              padding: 4px 8px;
              border-radius: 4px;
              font-size: 13px;
              color: #409eff;
            }

            &.resource-value {
              font-family: 'Monaco', 'Consolas', monospace;
              background: #f0f9ff;
              padding: 4px 8px;
              border-radius: 4px;
              font-size: 13px;
              color: #0369a1;
            }

            .system-tag {
              margin-left: 8px;
            }
          }
        }
      }

      .description-content {
        padding: 16px;
        background: #f8f9fa;
        border-radius: 8px;
        color: #606266;
        line-height: 1.6;
        font-size: 14px;
      }

      .usage-examples {
        .example-item {
          margin-bottom: 16px;

          &:last-child {
            margin-bottom: 0;
          }

          .example-title {
            font-weight: 500;
            color: #303133;
            margin-bottom: 8px;
            font-size: 14px;
          }

          .example-code {
            background: #f5f7fa;
            border: 1px solid #e4e7ed;
            border-radius: 6px;
            padding: 12px;

            code {
              font-family: 'Monaco', 'Consolas', monospace;
              font-size: 13px;
              color: #e83e8c;
              background: none;
            }
          }
        }
      }

      .json-viewer {
        background: #f8f9fa;
        border: 1px solid #e4e7ed;
        border-radius: 8px;
        max-height: 300px;
        overflow-y: auto;

        pre {
          margin: 0;
          padding: 16px;
          font-family: 'Monaco', 'Consolas', monospace;
          font-size: 12px;
          line-height: 1.5;

          code {
            color: #24292e;
            background: none;
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

  .permission-detail {
    .detail-section {
      margin-bottom: 20px;

      .section-header {
        margin-bottom: 12px;

        .section-title {
          font-size: 14px;
        }
      }

      .section-content {
        .detail-grid {
          .detail-item {
            flex-direction: column;
            align-items: flex-start;
            gap: 4px;

            .item-label {
              min-width: auto;
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