<template>
  <div class="history-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><DocumentCopy /></el-icon>
          生成历史
        </h2>
        <p class="page-desc">查看和管理代码生成历史记录</p>
      </div>
      <div class="header-actions">
        <el-button
          :icon="Plus"
          @click="goToTableSelection"
          size="default"
        >
          新建生成
        </el-button>
        <el-button
          :icon="RefreshRight"
          @click="refreshHistory"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="search-card" shadow="never">
      <el-form
        :model="searchForm"
        inline
        @submit.prevent="handleSearch"
        class="search-form"
      >
        <el-form-item label="表名" class="search-item">
          <el-input
            v-model="searchForm.tableName"
            placeholder="请输入表名"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
            :prefix-icon="Search"
          />
        </el-form-item>

        <el-form-item label="状态" class="search-item">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
            <el-option label="处理中" value="processing" />
          </el-select>
        </el-form-item>

        <el-form-item class="search-actions">
          <el-button
            type="primary"
            :icon="Search"
            @click="handleSearch"
            :loading="loading"
          >
            搜索
          </el-button>
          <el-button @click="resetSearch" :disabled="loading">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 历史记录表格 -->
    <el-card class="table-card" shadow="never">
      <!-- 空状态提示 -->
      <div v-if="historyList.length === 0 && !loading" class="empty-state">
        <el-empty
          description="暂无生成历史记录"
          :image-size="120"
        >
          <template #image>
            <el-icon class="empty-icon"><DocumentCopy /></el-icon>
          </template>
          <el-button type="primary" @click="goToTableSelection">
            <el-icon><Plus /></el-icon>
            开始第一次生成
          </el-button>
        </el-empty>
      </div>

      <!-- 数据表格 -->
      <transition name="fade" mode="out-in" v-else>
        <el-table
          v-loading="loading"
          :data="historyList"
          stripe
          @selection-change="handleSelectionChange"
          style="width: 100%"
          element-loading-text="正在加载生成历史..."
          element-loading-background="rgba(0, 0, 0, 0.1)"
        >
        <el-table-column type="selection" width="55" />

        <el-table-column prop="tableName" label="表名" width="150">
          <template #default="{ row }">
            <div class="table-info">
              <span class="table-name">{{ row.tableName }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="businessName" label="业务名称" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.businessName }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="generateType" label="生成类型" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getGenerateTypeTagType(row.generateType)"
              size="small"
            >
              {{ getGenerateTypeLabel(row.generateType) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              size="small"
            >
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="文件信息" width="150">
          <template #default="{ row }">
            <div class="file-info">
              <div class="file-count">{{ row.fileCount }} 个文件</div>
              <div class="file-size">{{ formatFileSize(row.fileSize) }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="downloadCount" label="下载次数" width="100">
          <template #default="{ row }">
            <span>{{ row.downloadCount }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="生成时间" width="180">
          <template #default="{ row }">
            <span>{{ formatDateTime(row.createdAt) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="errorMessage" label="错误信息" min-width="200">
          <template #default="{ row }">
            <div v-if="row.status === 'failed' && row.errorMessage" class="error-message">
              <el-tooltip :content="row.errorMessage" placement="top">
                <span>{{ truncateText(row.errorMessage, 30) }}</span>
              </el-tooltip>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip v-if="row.status === 'success'" content="下载代码包" placement="top">
                <el-button
                  type="primary"
                  size="small"
                  :icon="Download"
                  @click="handleDownload(row)"
                  circle
                />
              </el-tooltip>

              <el-tooltip content="重新生成" placement="top">
                <el-button
                  type="info"
                  size="small"
                  :icon="RefreshRight"
                  @click="handleRegenerate(row)"
                  circle
                />
              </el-tooltip>

              <el-tooltip content="删除记录" placement="top">
                <el-popconfirm
                  title="确定要删除这条记录吗？"
                  @confirm="handleDelete(row)"
                >
                  <template #reference>
                    <el-button
                      type="danger"
                      size="small"
                      :icon="Delete"
                      circle
                    />
                  </template>
                </el-popconfirm>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>
      </transition>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <div class="pagination-info">
          <span class="pagination-text">
            共 {{ pagination.total }} 条记录，显示第 {{ (pagination.page - 1) * pagination.pageSize + 1 }} -
            {{ Math.min(pagination.page * pagination.pageSize, pagination.total) }} 条
          </span>
        </div>
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="sizes, prev, pager, next, jumper"
          @size-change="handleSearch"
          @current-change="handleSearch"
        />
      </div>
    </el-card>

    <!-- 批量操作 -->
    <div class="batch-actions" v-show="selectedRows.length > 0">
      <div class="selected-info">
        已选择 <strong>{{ selectedRows.length }}</strong> 条记录
      </div>
      <div class="actions">
        <el-button @click="handleBatchDownload" :disabled="!canBatchDownload">
          <el-icon><Download /></el-icon>
          批量下载
        </el-button>
        <el-button type="danger" @click="handleBatchDelete">
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
        <el-button @click="clearSelection">
          清除选择
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  RefreshRight,
  Search,
  Download,
  Delete,
  DocumentCopy
} from '@element-plus/icons-vue'
import { getGenHistory, downloadCode,deleteGenHistory } from '@/api/generator'
import type { GenHistory } from '@/types/generator'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const historyList = ref<GenHistory[]>([])
const selectedRows = ref<GenHistory[]>([])

const searchForm = reactive({
  tableName: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 计算属性
const canBatchDownload = computed(() => {
  return selectedRows.value.some(row => row.status === 'success')
})

// 方法
const loadHistory = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }

    const response = await getGenHistory(params)
    historyList.value = response.data.list
    pagination.total = response.data.total
  } catch (error) {
    ElMessage.error('获取历史记录失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const refreshHistory = () => {
  loadHistory()
}

const handleSearch = () => {
  pagination.page = 1
  loadHistory()
}

const resetSearch = () => {
  searchForm.tableName = ''
  searchForm.status = ''
  handleSearch()
}

const handleSelectionChange = (selection: GenHistory[]) => {
  selectedRows.value = selection
}

const clearSelection = () => {
  selectedRows.value = []
}

const handleDownload = async (row: GenHistory) => {
  if (row.status !== 'success') {
    ElMessage.warning('只能下载生成成功的文件')
    return
  }

  try {
    // 使用历史记录ID下载对应的代码包
    const response = await downloadCode(row.id.toString())

    // 创建下载链接
    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${row.businessName}_code.zip`
    a.click()
    URL.revokeObjectURL(url)

    ElMessage.success('文件下载成功')

    // 更新下载次数（这里应该调用API更新）
    row.downloadCount += 1
  } catch (error) {
    ElMessage.error('文件下载失败')
    console.error(error)
  }
}

const handleRegenerate = (row: GenHistory) => {
  router.push({
    name: 'GeneratorConfig',
    params: { tableId: row.tableName }
  })
}

const handleDelete = async (row: GenHistory) => {
  try {
    // TODO: 调用删除API
    await deleteGenHistory(row.id)

    ElMessage.success('删除成功')
    loadHistory()
  } catch (error) {
    ElMessage.error('删除失败')
    console.error(error)
  }
}

const handleBatchDownload = async () => {
  const successRows = selectedRows.value.filter(row => row.status === 'success')
  if (successRows.length === 0) {
    ElMessage.warning('请选择生成成功的记录')
    return
  }

  const result = await ElMessageBox.confirm(
    `确定要下载选中的 ${successRows.length} 个文件吗？`,
    '批量下载确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).catch(() => false)

  if (result) {
    // TODO: 实现批量下载逻辑
    ElMessage.success('批量下载任务已启动')
  }
}

const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要删除的记录')
    return
  }

  const result = await ElMessageBox.confirm(
    `确定要删除选中的 ${selectedRows.value.length} 条记录吗？此操作不可恢复。`,
    '批量删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(() => false)

  if (result) {
    try {
      // TODO: 实现批量删除API调用
      ElMessage.success('批量删除成功')
      clearSelection()
      loadHistory()
    } catch (error) {
      ElMessage.error('批量删除失败')
      console.error(error)
    }
  }
}

const goToTableSelection = () => {
  router.push({ name: 'GeneratorTables' })
}

// 工具函数
const getStatusLabel = (status: string): string => {
  const labels: { [key: string]: string } = {
    'success': '成功',
    'failed': '失败',
    'processing': '处理中'
  }
  return labels[status] || status
}

const getStatusTagType = (status: string): string => {
  const types: { [key: string]: string } = {
    'success': 'success',
    'failed': 'danger',
    'processing': 'warning'
  }
  return types[status] || ''
}

const getGenerateTypeLabel = (type: string): string => {
  const labels: { [key: string]: string } = {
    'all': '全部',
    'backend': '后端',
    'frontend': '前端'
  }
  return labels[type] || type
}

const getGenerateTypeTagType = (type: string): string => {
  const types: { [key: string]: string } = {
    'all': 'primary',
    'backend': 'success',
    'frontend': 'info'
  }
  return types[type] || ''
}

const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(size) / Math.log(k))
  return parseFloat((size / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatDateTime = (dateTime: string): string => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString('zh-CN')
}

const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

// 生命周期
onMounted(() => {
  loadHistory()
})
</script>

<style scoped lang="scss">
.history-page {
  padding: 24px;
  background: #f6f8fa;
  min-height: calc(100vh - 60px);

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
    padding: 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);

    .header-left {
      color: white;

      .page-title {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 12px;

        .title-icon {
          font-size: 32px;
        }
      }

      .page-desc {
        margin: 0;
        font-size: 14px;
        opacity: 0.9;
        line-height: 1.4;
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;

      .el-button {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.2);
        color: white;
        backdrop-filter: blur(10px);

        &:hover {
          background: rgba(255, 255, 255, 0.2);
          border-color: rgba(255, 255, 255, 0.3);
          transform: translateY(-2px);
        }
      }
    }
  }

  .search-card {
    margin-bottom: 24px;
    border: 1px solid #d1d9e0;
    border-radius: 12px;

    :deep(.el-card__body) {
      padding: 20px 24px;
    }

    .search-form {
      .search-item {
        margin-right: 16px;
        margin-bottom: 0;
      }

      .search-actions {
        margin-left: auto;
      }
    }
  }

  .table-card {
    border: 1px solid #d1d9e0;
    border-radius: 12px;

    :deep(.el-card__body) {
      padding: 0;
    }

    .el-table {
      .table-info {
        .table-name {
          font-weight: 600;
          color: #1f2328;
        }
      }

      .file-info {
        font-size: 12px;
        color: #909399;

        .file-count {
          margin-bottom: 2px;
          color: #606266;
        }
      }

      .error-message {
        color: #f56c6c;
        font-size: 12px;
        cursor: pointer;
      }

      .action-buttons {
        display: flex;
        gap: 8px;
        justify-content: center;
      }
    }

    .empty-state {
      padding: 60px 20px;
      text-align: center;

      .empty-icon {
        font-size: 120px;
        color: #dcdfe6;
        margin-bottom: 16px;
      }

      :deep(.el-empty__description) {
        color: #909399;
        font-size: 16px;
      }

      .el-button {
        margin-top: 16px;
      }
    }

    .pagination-wrapper {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px 24px;
      border-top: 1px solid #d1d9e0;
      background: #f8f9fa;

      .pagination-info {
        .pagination-text {
          font-size: 14px;
          color: #656d76;
        }
      }
    }
  }

  .batch-actions {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.16);
    padding: 16px 20px;
    display: flex;
    align-items: center;
    gap: 20px;
    z-index: 1000;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);

    .selected-info {
      color: #606266;
      font-size: 14px;
      font-weight: 500;
    }

    .actions {
      display: flex;
      gap: 12px;
    }
  }
}

// 表格样式优化
:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;

  .el-table__header {
    th {
      background-color: #f5f7fa;
      color: #606266;
      font-weight: 600;
    }
  }

  .el-table__body {
    tr:hover {
      background-color: #f5f7fa;
    }
  }

  .cell {
    padding: 12px 0;
  }
}

// 过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 响应式设计
@media (max-width: 768px) {
  .history-page {
    padding: 16px;

    .page-header {
      flex-direction: column;
      gap: 16px;

      .header-actions {
        width: 100%;

        .el-button {
          flex: 1;
        }
      }
    }

    .search-card {
      .search-form {
        .el-form-item {
          width: 100%;
          margin-right: 0;
          margin-bottom: 16px;

          .el-input,
          .el-select {
            width: 100% !important;
          }
        }
      }
    }

    .batch-actions {
      bottom: 16px;
      right: 16px;
      left: 16px;
      flex-direction: column;
      gap: 12px;

      .actions {
        justify-content: center;
        width: 100%;

        .el-button {
          flex: 1;
        }
      }
    }
  }
}
</style>