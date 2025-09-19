<template>
  <el-dialog
    v-model="dialogVisible"
    title="租户详情"
    width="700px"
    :close-on-click-modal="false"
  >
    <div class="tenant-detail" v-if="tenantData">
      <!-- 基本信息 -->
      <div class="detail-section">
        <h3 class="section-title">
          <el-icon><InfoFilled /></el-icon>
          基本信息
        </h3>
        <div class="detail-grid">
          <div class="detail-item">
            <span class="label">租户ID：</span>
            <span class="value">{{ tenantData.id }}</span>
          </div>
          <div class="detail-item">
            <span class="label">租户名称：</span>
            <span class="value">{{ tenantData.name }}</span>
          </div>
          <div class="detail-item">
            <span class="label">域名：</span>
            <span class="value">
              <el-link
                v-if="tenantData.domain"
                :href="`http://${tenantData.domain}`"
                target="_blank"
                type="primary"
              >
                {{ tenantData.domain }}
              </el-link>
              <span v-else class="text-gray">未设置</span>
            </span>
          </div>
          <div class="detail-item">
            <span class="label">状态：</span>
            <span class="value">
              <el-tag
                :type="getStatusTagType(tenantData.status)"
                size="small"
                effect="light"
              >
                {{ getStatusText(tenantData.status) }}
              </el-tag>
            </span>
          </div>
        </div>
      </div>

      <!-- 时间信息 -->
      <div class="detail-section">
        <h3 class="section-title">
          <el-icon><Clock /></el-icon>
          时间信息
        </h3>
        <div class="detail-grid">
          <div class="detail-item">
            <span class="label">创建时间：</span>
            <span class="value">{{ formatDateTime(tenantData.created_at) }}</span>
          </div>
          <div class="detail-item">
            <span class="label">更新时间：</span>
            <span class="value">{{ formatDateTime(tenantData.updated_at) }}</span>
          </div>
          <div class="detail-item">
            <span class="label">过期时间：</span>
            <span class="value">
              <div v-if="tenantData.expired_at">
                <span :class="{ 'text-danger': isExpired(tenantData.expired_at) }">
                  {{ formatDateTime(tenantData.expired_at) }}
                </span>
                <el-tag v-if="isExpired(tenantData.expired_at)" type="danger" size="small" class="ml-2">
                  已过期
                </el-tag>
              </div>
              <span v-else class="text-success">永久有效</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 配置信息 -->
      <div class="detail-section" v-if="tenantData.config">
        <h3 class="section-title">
          <el-icon><Setting /></el-icon>
          配置信息
        </h3>
        <div class="config-content">
          <el-input
            :model-value="formatConfig(tenantData.config)"
            type="textarea"
            :rows="6"
            readonly
            class="config-textarea"
          />
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="detail-section">
        <h3 class="section-title">
          <el-icon><DataAnalysis /></el-icon>
          统计信息
        </h3>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ tenantStats.userCount }}</div>
              <div class="stat-label">用户数量</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><Files /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ tenantStats.fileCount }}</div>
              <div class="stat-label">文件数量</div>
            </div>
          </div>
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><Calendar /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ getDaysUntilExpiry() }}</div>
              <div class="stat-label">剩余天数</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="handleEdit">编辑租户</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import {
  InfoFilled,
  Clock,
  Setting,
  DataAnalysis,
  User,
  Files,
  Calendar
} from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/date'

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
const emit = defineEmits(['update:visible', 'edit'])

// 响应式数据
const tenantStats = reactive({
  userCount: 0,
  fileCount: 0,
  storageUsed: 0
})

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 状态映射
const statusMap = {
  1: { text: '启用', type: 'success' },
  2: { text: '禁用', type: 'danger' },
  3: { text: '试用', type: 'warning' },
  4: { text: '过期', type: 'info' }
}

// 获取状态文本
const getStatusText = (status) => {
  return statusMap[status]?.text || '未知'
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  return statusMap[status]?.type || 'info'
}

// 检查是否过期
const isExpired = (expiredAt) => {
  if (!expiredAt) return false
  return new Date(expiredAt) < new Date()
}

// 格式化配置信息
const formatConfig = (config) => {
  if (!config) return ''

  try {
    const parsed = typeof config === 'string' ? JSON.parse(config) : config
    return JSON.stringify(parsed, null, 2)
  } catch (error) {
    return config
  }
}

// 计算剩余天数
const getDaysUntilExpiry = () => {
  if (!props.tenantData?.expired_at) return '永久'

  const expiredDate = new Date(props.tenantData.expired_at)
  const now = new Date()
  const diffTime = expiredDate - now
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return '已过期'
  if (diffDays === 0) return '今天过期'
  return `${diffDays} 天`
}

// 获取租户统计信息
const getTenantStats = async () => {
  if (!props.tenantData?.id) return

  try {
    // TODO: 实现获取租户统计信息的API
    // const response = await tenantApi.getTenantStats(props.tenantData.id)
    // Object.assign(tenantStats, response.data)

    // 模拟数据
    Object.assign(tenantStats, {
      userCount: Math.floor(Math.random() * 100),
      fileCount: Math.floor(Math.random() * 500),
      storageUsed: Math.floor(Math.random() * 1024)
    })
  } catch (error) {
    console.error('获取租户统计信息失败:', error)
  }
}

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible && props.tenantData) {
    getTenantStats()
  }
})

// 关闭弹窗
const handleClose = () => {
  emit('update:visible', false)
}

// 编辑租户
const handleEdit = () => {
  emit('edit', props.tenantData)
  handleClose()
}
</script>

<style scoped>
.tenant-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-section {
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.detail-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.section-title {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title .el-icon {
  color: #3b82f6;
  font-size: 18px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.detail-item {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.label {
  flex-shrink: 0;
  width: 100px;
  color: #6b7280;
  font-size: 14px;
}

.value {
  flex: 1;
  color: #1f2937;
  font-size: 14px;
  word-break: break-all;
}

.config-content {
  margin-top: 12px;
}

.config-textarea {
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 12px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 12px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  gap: 12px;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 18px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.text-gray {
  color: #9ca3af;
}

.text-success {
  color: #10b981;
}

.text-danger {
  color: #ef4444;
}

.ml-2 {
  margin-left: 8px;
}

:deep(.config-textarea .el-textarea__inner) {
  background-color: #f8fafc;
  color: #374151;
}
</style>