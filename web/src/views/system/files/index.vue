<template>
  <div class="file-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><FolderOpened /></el-icon>
          文件管理
          <!-- 当前租户显示 -->
          <el-tag v-if="isSuperAdmin && currentTenant" type="primary" effect="light" size="default" class="tenant-indicator">
            <el-icon><OfficeBuilding /></el-icon>
            {{ currentTenant.name }}
          </el-tag>
          <el-tag v-else-if="isSuperAdmin && !currentTenant" type="warning" effect="light" size="default" class="tenant-indicator">
            <el-icon><Warning /></el-icon>
            请选择租户
          </el-tag>
        </h2>
        <p class="page-desc">
          {{ isSuperAdmin && currentTenant ? `管理租户 "${currentTenant.name}" 的文件资源和存储空间` : isSuperAdmin ? '请先选择要管理的租户' : '管理系统文件资源和存储空间' }}
        </p>
      </div>
      <div class="header-actions">
        <el-button type="primary" :icon="UploadFilled" @click="showUploadDialog = true" size="default">
          上传文件
        </el-button>
        <el-button
          :icon="RefreshRight"
          @click="refreshFiles"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
        <el-button
          :icon="Download"
          @click="exportFiles"
          size="default"
        >
          导出列表
        </el-button>
      </div>
    </div>

    <!-- 操作工具栏 -->
    <el-card class="toolbar-card" shadow="never">
      <div class="toolbar-content">
        <!-- 搜索区域 -->
        <div class="search-section">
          <el-form :model="searchForm" inline class="search-form">
            <el-form-item label="文件名" class="search-item">
              <el-input
                v-model="searchForm.filename"
                placeholder="搜索文件名"
                clearable
                @keyup.enter="handleSearch"
                @clear="handleSearch"
                :prefix-icon="Search"
                style="width: 260px"
              />
            </el-form-item>
            <el-form-item label="文件类型" class="search-item">
              <el-select
                v-model="searchForm.fileType"
                placeholder="文件类型"
                clearable
                @change="handleSearch"
                style="width: 140px"
              >
                <el-option label="全部" value="" />
                <el-option label="图片" value="jpg,jpeg,png,gif,webp" />
                <el-option label="文档" value="pdf,doc,docx" />
                <el-option label="表格" value="xls,xlsx" />
                <el-option label="文本" value="txt" />
              </el-select>
            </el-form-item>
            <el-form-item label="使用类型" class="search-item">
              <el-select
                v-model="searchForm.usageType"
                placeholder="使用类型"
                clearable
                @change="handleSearch"
                style="width: 140px"
              >
                <el-option label="全部" value="" />
                <el-option label="头像" value="avatar" />
                <el-option label="附件" value="attachment" />
                <el-option label="图片" value="image" />
                <el-option label="文档" value="document" />
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
              <el-button @click="handleResetSearch" :disabled="loading">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 批量操作区域 -->
        <div class="batch-section" v-show="selectedFiles.length > 0">
          <div class="batch-info">
            <el-icon class="info-icon"><InfoFilled /></el-icon>
            已选择 <strong class="selected-count">{{ selectedFiles.length }}</strong> 个文件
          </div>
          <div class="batch-actions">
            <el-button type="danger" size="small" :icon="Delete" @click="confirmBatchDelete">
              批量删除
            </el-button>
            <el-button size="small" @click="clearSelection">
              取消选择
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 统计信息 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon stats-icon-primary">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.totalFiles }}</div>
              <div class="stats-label">总文件数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon stats-icon-success">
              <el-icon><Folder /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ formatFileSize(stats.totalSize) }}</div>
              <div class="stats-label">总存储大小</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon stats-icon-warning">
              <el-icon><Picture /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.imageCount }}</div>
              <div class="stats-label">图片文件</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <div class="stats-icon stats-icon-info">
              <el-icon><DocumentCopy /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.documentCount }}</div>
              <div class="stats-label">文档文件</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 文件列表 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="fileList"
        stripe
        @selection-change="handleSelectionChange"
        style="width: 100%"
        :default-sort="{ prop: 'createdAt', order: 'descending' }"
      >
        <el-table-column type="selection" width="50" />

        <el-table-column label="文件信息" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="file-info">
              <div class="file-preview">
                <el-image
                  v-if="isImageFile(row.mimeType)"
                  :src="getFilePreviewUrl(row.id)"
                  :preview-src-list="[getFilePreviewUrl(row.id)]"
                  fit="cover"
                  class="file-thumbnail"
                  @error="handleImageError"
                >
                  <template #error>
                    <div class="file-icon">
                      <el-icon color="#409eff"><Picture /></el-icon>
                    </div>
                  </template>
                </el-image>
                <div v-else class="file-icon">
                  <el-icon v-if="isDocumentFile(row.mimeType)" color="#67c23a">
                    <Document />
                  </el-icon>
                  <el-icon v-else color="#909399">
                    <Folder />
                  </el-icon>
                </div>
              </div>
              <div class="file-details">
                <div class="file-name" :title="row.originalName">
                  {{ row.originalName }}
                </div>
                <div class="file-meta">
                  <el-tag size="small" type="info" class="file-type-tag">
                    {{ row.fileType.toUpperCase() }}
                  </el-tag>
                  <span class="file-size">{{ formatFileSize(row.fileSize) }}</span>
                  <span class="file-time">{{ formatTime(row.createdAt) }}</span>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="使用类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.usageType" size="small" :type="getUsageTypeColor(row.usageType)">
              {{ getUsageTypeLabel(row.usageType) }}
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>

        <el-table-column label="上传用户" width="140" show-overflow-tooltip>
          <template #default="{ row }">
            <div v-if="row.uploadUser" class="user-info">
              <el-avatar :size="24" :src="row.uploadUser.avatar" class="user-avatar">
                <el-icon><User /></el-icon>
              </el-avatar>
              <div class="user-details">
                <div class="username">{{ row.uploadUser.username }}</div>
                <div class="nickname">{{ row.uploadUser.nickname }}</div>
              </div>
            </div>
            <span v-else class="text-muted">系统</span>
          </template>
        </el-table-column>

        <el-table-column label="MD5" width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tooltip :content="row.md5" placement="top">
              <code class="md5-text">{{ row.md5.substring(0, 8) }}...</code>
            </el-tooltip>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" width="160" align="center" sortable prop="createdAt">
          <template #default="{ row }">
            <div class="time-info">
              <div class="date">{{ formatDate(row.createdAt) }}</div>
              <div class="time">{{ formatTimeOnly(row.createdAt) }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip content="下载文件" placement="top">
                <el-button
                  type="primary"
                  size="small"
                  :icon="Download"
                  @click="downloadFile(row)"
                  :loading="downloadingIds.includes(row.id)"
                  circle
                />
              </el-tooltip>
              <el-tooltip content="复制链接" placement="top">
                <el-button
                  type="info"
                  size="small"
                  :icon="Link"
                  @click="copyFileLink(row)"
                  circle
                />
              </el-tooltip>
              <el-tooltip content="删除文件" placement="top">
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  @click="confirmDeleteFile(row)"
                  :loading="deletingIds.includes(row.id)"
                  circle
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

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
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 上传对话框 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传文件"
      width="600px"
      :close-on-click-modal="false"
    >
      <FileUpload
        v-model="uploadedFile"
        usage-type="document"
        @success="handleUploadSuccess"
      />
      <template #footer>
        <el-button @click="showUploadDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  UploadFilled,
  Search,
  RefreshRight,
  Download,
  FolderOpened,
  OfficeBuilding,
  Warning,
  InfoFilled,
  Delete,
  Picture,
  Document,
  Folder,
  DocumentCopy,
  User,
  Link
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { useTenantStore } from '@/store/tenant'
import {
  getAllFiles,
  downloadFile as downloadFileApi,
  deleteFile,
  type FileProfile
} from '@/api/file'
import FileUpload from '@/components/FileUpload.vue'

const userStore = useUserStore()
const tenantStore = useTenantStore()

// 响应式数据
const loading = ref(false)
const batchDeleting = ref(false)
const showUploadDialog = ref(false)
const uploadedFile = ref<FileProfile | null>(null)
const fileList = ref<FileProfile[]>([])
const selectedFiles = ref<FileProfile[]>([])
const downloadingIds = ref<number[]>([])
const deletingIds = ref<number[]>([])

// 搜索表单
const searchForm = reactive({
  filename: '',
  fileType: '',
  usageType: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 统计数据
const stats = reactive({
  totalFiles: 0,
  totalSize: 0,
  imageCount: 0,
  documentCount: 0
})

// 计算属性
const isSuperAdmin = computed(() => tenantStore.checkIsSuperAdmin())
const currentTenant = computed(() => tenantStore.getCurrentTenant())

// 监听租户变化
watch(currentTenant, () => {
  if (isSuperAdmin.value) {
    pagination.page = 1
    loadFileList()
  }
}, { immediate: false })

// 文件大小格式化
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

// 时间格式化
const formatTime = (timeStr: string): string => {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    return '今天 ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else if (days === 1) {
    return '昨天 ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
}

const formatDate = (timeStr: string): string => {
  return new Date(timeStr).toLocaleDateString('zh-CN')
}

const formatTimeOnly = (timeStr: string): string => {
  return new Date(timeStr).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 判断是否为图片文件
const isImageFile = (mimeType: string): boolean => {
  return mimeType.startsWith('image/')
}

// 判断是否为文档文件
const isDocumentFile = (mimeType: string): boolean => {
  const docTypes = [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'text/plain'
  ]
  return docTypes.includes(mimeType)
}

// 获取使用类型颜色
const getUsageTypeColor = (usageType: string): string => {
  const colorMap: Record<string, string> = {
    avatar: 'success',
    attachment: 'info',
    image: 'warning',
    document: 'primary'
  }
  return colorMap[usageType] || 'info'
}

// 获取使用类型标签
const getUsageTypeLabel = (usageType: string): string => {
  const labelMap: Record<string, string> = {
    avatar: '头像',
    attachment: '附件',
    image: '图片',
    document: '文档'
  }
  return labelMap[usageType] || usageType
}

// 获取文件预览URL
const getFilePreviewUrl = (fileId: number): string => {
  return `/api/v1/files/${fileId}/download`
}

// 处理图片错误
const handleImageError = () => {
  // 图片加载失败时的处理
}

// 加载文件列表
const loadFileList = async () => {
  if (isSuperAdmin.value && !currentTenant.value) {
    fileList.value = []
    pagination.total = 0
    updateStats([])
    return
  }

  try {
    loading.value = true

    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    }

    // 清空空值
    Object.keys(params).forEach(key => {
      if (!params[key as keyof typeof params]) {
        delete params[key as keyof typeof params]
      }
    })

    const response = await getAllFiles(params)
    fileList.value = response.data.files
    pagination.total = response.data.pagination.total

    // 更新统计信息
    updateStats(fileList.value)
  } catch (error: any) {
    ElMessage.error(error.message || '加载文件列表失败')
  } finally {
    loading.value = false
  }
}

// 更新统计信息
const updateStats = (files: FileProfile[]) => {
  stats.totalFiles = files.length
  stats.totalSize = files.reduce((total, file) => total + file.fileSize, 0)
  stats.imageCount = files.filter(file => isImageFile(file.mimeType)).length
  stats.documentCount = files.filter(file => isDocumentFile(file.mimeType)).length
}

// 刷新文件列表
const refreshFiles = () => {
  pagination.page = 1
  loadFileList()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadFileList()
}

// 重置搜索
const handleResetSearch = () => {
  Object.assign(searchForm, {
    filename: '',
    fileType: '',
    usageType: ''
  })
  pagination.page = 1
  loadFileList()
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  loadFileList()
}

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  pagination.page = 1
  loadFileList()
}

// 选择变化
const handleSelectionChange = (selection: FileProfile[]) => {
  selectedFiles.value = selection
}

// 清空选择
const clearSelection = () => {
  selectedFiles.value = []
}

// 下载文件
const downloadFile = async (file: FileProfile) => {
  try {
    downloadingIds.value.push(file.id)
    await downloadFileApi(file.id)
  } catch (error: any) {
    ElMessage.error(error.message || '下载失败')
  } finally {
    downloadingIds.value = downloadingIds.value.filter(id => id !== file.id)
  }
}

// 复制文件链接
const copyFileLink = async (file: FileProfile) => {
  try {
    const url = `${window.location.origin}/api/v1/files/${file.id}/download`
    await navigator.clipboard.writeText(url)
    ElMessage.success('文件链接已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制链接失败')
  }
}

// 确认删除文件
const confirmDeleteFile = async (file: FileProfile) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文件"${file.originalName}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: true
      }
    )

    deletingIds.value.push(file.id)
    await deleteFile(file.id)
    ElMessage.success('文件删除成功')
    await loadFileList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  } finally {
    deletingIds.value = deletingIds.value.filter(id => id !== file.id)
  }
}

// 确认批量删除
const confirmBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedFiles.value.length} 个文件吗？`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    batchDeleting.value = true

    // 逐个删除文件
    for (const file of selectedFiles.value) {
      await deleteFile(file.id)
    }

    ElMessage.success('文件批量删除成功')
    clearSelection()
    await loadFileList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '批量删除失败')
    }
  } finally {
    batchDeleting.value = false
  }
}

// 导出文件列表
const exportFiles = () => {
  ElMessage.info('导出功能开发中...')
}

// 上传成功
const handleUploadSuccess = () => {
  showUploadDialog.value = false
  uploadedFile.value = null
  loadFileList()
}

// 组件挂载
onMounted(() => {
  loadFileList()
})
</script>

<style lang="scss" scoped>
.file-management {
  padding: 24px;
  background: #f6f8fa;
  min-height: calc(100vh - 60px);

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;

    .header-left {
      .page-title {
        display: flex;
        align-items: center;
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        color: #1f2328;

        .title-icon {
          margin-right: 12px;
          font-size: 32px;
          color: #0969da;
        }

        .tenant-indicator {
          margin-left: 16px;
          font-size: 14px;
          font-weight: 500;
        }
      }

      .page-desc {
        margin: 0;
        color: #656d76;
        font-size: 16px;
        line-height: 1.5;
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .toolbar-card {
    margin-bottom: 24px;
    border: 1px solid #d1d9e0;
    border-radius: 12px;

    :deep(.el-card__body) {
      padding: 20px 24px;
    }

    .toolbar-content {
      .search-section {
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

      .batch-section {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: 16px;
        padding: 12px 16px;
        background: #fff3cd;
        border: 1px solid #ffeaa7;
        border-radius: 8px;

        .batch-info {
          display: flex;
          align-items: center;
          color: #856404;

          .info-icon {
            margin-right: 8px;
            font-size: 16px;
          }

          .selected-count {
            color: #d63384;
          }
        }

        .batch-actions {
          display: flex;
          gap: 8px;
        }
      }
    }
  }

  .stats-row {
    margin-bottom: 24px;

    .stats-card {
      border: 1px solid #d1d9e0;
      border-radius: 12px;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
        transform: translateY(-2px);
      }

      :deep(.el-card__body) {
        padding: 20px;
      }

      .stats-content {
        display: flex;
        align-items: center;

        .stats-icon {
          width: 48px;
          height: 48px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 16px;
          font-size: 24px;

          &.stats-icon-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
          }

          &.stats-icon-success {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
            color: white;
          }

          &.stats-icon-warning {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
          }

          &.stats-icon-info {
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
            color: white;
          }
        }

        .stats-info {
          .stats-number {
            font-size: 24px;
            font-weight: 700;
            color: #1f2328;
            line-height: 1.2;
          }

          .stats-label {
            font-size: 14px;
            color: #656d76;
            margin-top: 4px;
          }
        }
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
      .file-info {
        display: flex;
        align-items: center;
        gap: 12px;

        .file-preview {
          width: 48px;
          height: 48px;
          border-radius: 8px;
          overflow: hidden;
          border: 2px solid #f1f3f4;

          .file-thumbnail {
            width: 100%;
            height: 100%;
          }

          .file-icon {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #f8f9fa;
            font-size: 20px;
          }
        }

        .file-details {
          flex: 1;
          min-width: 0;

          .file-name {
            font-weight: 600;
            color: #1f2328;
            margin-bottom: 6px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            cursor: pointer;

            &:hover {
              color: #0969da;
            }
          }

          .file-meta {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;
            color: #656d76;

            .file-type-tag {
              font-weight: 600;
            }
          }
        }
      }

      .user-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .user-avatar {
          flex-shrink: 0;
        }

        .user-details {
          min-width: 0;

          .username {
            font-weight: 500;
            color: #1f2328;
            font-size: 13px;
          }

          .nickname {
            font-size: 12px;
            color: #656d76;
            margin-top: 2px;
          }
        }
      }

      .md5-text {
        font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
        background: #f6f8fa;
        padding: 2px 6px;
        border-radius: 4px;
        font-size: 12px;
        color: #656d76;
        cursor: pointer;

        &:hover {
          background: #eef1f3;
        }
      }

      .time-info {
        text-align: center;

        .date {
          font-weight: 500;
          color: #1f2328;
          font-size: 13px;
        }

        .time {
          font-size: 12px;
          color: #656d76;
          margin-top: 2px;
        }
      }

      .action-buttons {
        display: flex;
        gap: 8px;
        justify-content: center;
      }

      .text-muted {
        color: #8c959f;
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
}

// 响应式设计
@media (max-width: 1200px) {
  .file-management {
    .stats-row {
      .el-col {
        margin-bottom: 16px;
      }
    }
  }
}

@media (max-width: 768px) {
  .file-management {
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

    .toolbar-card {
      .toolbar-content {
        .search-section {
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

        .batch-section {
          flex-direction: column;
          gap: 12px;
          text-align: center;

          .batch-actions {
            justify-content: center;
          }
        }
      }
    }

    .stats-row {
      .el-col {
        margin-bottom: 16px;
      }
    }
  }
}
</style>