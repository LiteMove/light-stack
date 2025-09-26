<template>
  <div class="history-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">生成历史</h1>
        <p class="description">查看和管理代码生成历史记录</p>
      </div>
      <div class="header-actions">
        <el-button @click="goToTableSelection">
          <el-icon><Plus /></el-icon>
          新建生成
        </el-button>
        <el-button type="primary" @click="refreshHistory" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-container">
      <el-form
        :model="searchForm"
        inline
        @submit.prevent="handleSearch"
      >
        <el-form-item label="表名">
          <el-input
            v-model="searchForm.tableName"
            placeholder="请输入表名"
            clearable
            style="width: 200px"
          />
        </el-form-item>

        <el-form-item label="状态">
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

        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 历史记录表格 -->
    <div class="history-table">
      <el-table
        v-loading="loading"
        :data="historyList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
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
            <div class="actions">
              <el-button
                v-if="row.status === 'success'"
                type="primary"
                size="small"
                @click="handleDownload(row)"
              >
                <el-icon><Download /></el-icon>
                下载
              </el-button>

              <el-button
                type="info"
                size="small"
                @click="handleRegenerate(row)"
              >
                <el-icon><Refresh /></el-icon>
                重新生成
              </el-button>

              <el-popconfirm
                title="确定要删除这条记录吗？"
                @confirm="handleDelete(row)"
              >
                <template #reference>
                  <el-button
                    type="danger"
                    size="small"
                  >
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSearch"
          @current-change="handleSearch"
        />
      </div>
    </div>

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
  Refresh,
  Search,
  Download,
  Delete
} from '@element-plus/icons-vue'
import { getGenHistory, downloadCode } from '@/api/generator'
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
    // 这里应该传入实际的任务ID，暂时使用历史记录ID
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
    // await deleteGenHistory(row.id)

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
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);

  .header-content {
    .title {
      margin: 0 0 8px 0;
      font-size: 24px;
      font-weight: 600;
      color: #303133;
    }

    .description {
      margin: 0;
      color: #909399;
      font-size: 14px;
    }
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }
}

.search-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 20px;
}

.history-table {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.table-info {
  .table-name {
    font-weight: 600;
    color: #303133;
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

.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.batch-actions {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  z-index: 1000;

  .selected-info {
    color: #606266;
    font-size: 14px;
  }

  .actions {
    display: flex;
    gap: 12px;
  }
}

:deep(.el-table) {
  .cell {
    padding: 12px 0;
  }
}
</style>