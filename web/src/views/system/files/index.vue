<template>
  <div class="file-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-title">
        <h2>文件管理</h2>
        <p class="subtitle">管理和查看上传的文件</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showUploadDialog = true">
          <el-icon><UploadFilled /></el-icon>
          上传文件
        </el-button>
      </div>
    </div>

    <!-- 搜索和过滤 -->
    <div class="search-bar">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-select
            v-model="searchForm.fileType"
            placeholder="文件类型"
            clearable
            @change="handleSearch"
          >
            <el-option label="图片" value="jpg,jpeg,png,gif" />
            <el-option label="文档" value="pdf,doc,docx" />
            <el-option label="表格" value="xls,xlsx" />
            <el-option label="文本" value="txt" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="searchForm.usageType"
            placeholder="使用类型"
            clearable
            @change="handleSearch"
          >
            <el-option label="头像" value="avatar" />
            <el-option label="附件" value="attachment" />
            <el-option label="图片" value="image" />
            <el-option label="文档" value="document" />
          </el-select>
        </el-col>
        <el-col :span="6" v-if="isAdmin">
          <el-input
            v-model="searchForm.uploadUserId"
            placeholder="上传用户ID"
            clearable
            @keyup.enter="handleSearch"
            @clear="handleSearch"
          />
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 文件列表 -->
    <div class="file-list">
      <el-table
        v-loading="loading"
        :data="fileList"
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column label="文件信息" min-width="300">
          <template #default="{ row }">
            <div class="file-info">
              <div class="file-icon">
                <el-icon v-if="isImageFile(row.mimeType)" color="#409eff">
                  <Picture />
                </el-icon>
                <el-icon v-else-if="isDocumentFile(row.mimeType)" color="#67c23a">
                  <Document />
                </el-icon>
                <el-icon v-else color="#909399">
                  <Folder />
                </el-icon>
              </div>
              <div class="file-details">
                <div class="file-name" :title="row.originalName">
                  {{ row.originalName }}
                </div>
                <div class="file-meta">
                  {{ formatFileSize(row.fileSize) }} • {{ row.fileType }} • {{ formatTime(row.createdAt) }}
                </div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="使用类型" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.usageType" size="small" type="info">
              {{ row.usageType }}
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>

        <el-table-column v-if="isAdmin" label="上传用户" width="120">
          <template #default="{ row }">
            <div v-if="row.uploadUser" class="user-info">
              <div class="username">{{ row.uploadUser.username }}</div>
              <div class="nickname text-muted">{{ row.uploadUser.nickname }}</div>
            </div>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>

        <el-table-column label="MD5" width="100">
          <template #default="{ row }">
            <el-tooltip :content="row.md5" placement="top">
              <span class="md5-text">{{ row.md5.substring(0, 8) }}...</span>
            </el-tooltip>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="downloadFile(row)"
              :loading="downloadingIds.includes(row.id)"
            >
              下载
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="confirmDeleteFile(row)"
              :loading="deletingIds.includes(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 批量操作 -->
    <div v-if="selectedFiles.length > 0" class="batch-actions">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedFiles.length }} 个文件</span>
          <div class="batch-buttons">
            <el-button
              type="danger"
              @click="confirmBatchDelete"
              :loading="batchDeleting"
            >
              批量删除
            </el-button>
            <el-button @click="clearSelection">取消选择</el-button>
          </div>
        </div>
      </el-card>
    </div>

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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  UploadFilled,
  Search,
  Refresh,
  Picture,
  Document,
  Folder
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import {
  getAllFiles,
  getUserFiles,
  downloadFile as downloadFileApi,
  deleteFile,
  type FileProfile
} from '@/api/file'
import FileUpload from '@/components/FileUpload.vue'

const userStore = useUserStore()

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
  fileType: '',
  usageType: '',
  uploadUserId: ''
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 计算属性
const isAdmin = computed(() => {
  // 这里应该根据实际的权限逻辑来判断
  return userStore.user?.roles?.some(role => role.name === 'admin' || role.name === 'super_admin')
})

// 文件大小格式化
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 时间格式化
const formatTime = (timeStr: string): string => {
  const date = new Date(timeStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
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

// 加载文件列表
const loadFileList = async () => {
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

    const response = isAdmin.value
      ? await getAllFiles(params)
      : await getUserFiles(params)

    fileList.value = response.data.files
    pagination.total = response.data.pagination.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载文件列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadFileList()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    fileType: '',
    usageType: '',
    uploadUserId: ''
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

// 确认删除文件
const confirmDeleteFile = async (file: FileProfile) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文件"${file.originalName}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
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
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;

    .header-title {
      h2 {
        margin: 0 0 8px 0;
        font-size: 24px;
        color: #303133;
      }

      .subtitle {
        margin: 0;
        color: #606266;
        font-size: 14px;
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .search-bar {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 24px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .file-list {
    background: #fff;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

    .file-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .file-icon {
        font-size: 24px;
        flex-shrink: 0;
      }

      .file-details {
        flex: 1;
        min-width: 0;

        .file-name {
          font-weight: 500;
          color: #303133;
          margin-bottom: 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .file-meta {
          font-size: 12px;
          color: #909399;
        }
      }
    }

    .user-info {
      .username {
        font-weight: 500;
        color: #303133;
        margin-bottom: 2px;
      }

      .nickname {
        font-size: 12px;
      }
    }

    .md5-text {
      font-family: monospace;
      color: #909399;
      cursor: pointer;
    }

    .text-muted {
      color: #c0c4cc;
    }

    .pagination-wrapper {
      padding: 20px;
      text-align: right;
      border-top: 1px solid #ebeef5;
    }
  }

  .batch-actions {
    position: fixed;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 1000;

    .el-card {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);

      :deep(.el-card__body) {
        padding: 16px 24px;
      }
    }

    .batch-info {
      display: flex;
      align-items: center;
      gap: 16px;

      .batch-buttons {
        display: flex;
        gap: 8px;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .file-management {
    padding: 16px;

    .page-header {
      flex-direction: column;
      gap: 16px;

      .header-actions {
        width: 100%;
        justify-content: stretch;

        .el-button {
          flex: 1;
        }
      }
    }

    .search-bar {
      padding: 16px;

      .el-row .el-col {
        margin-bottom: 12px;
      }
    }

    .batch-actions {
      left: 16px;
      right: 16px;
      transform: none;

      .batch-info {
        flex-direction: column;
        gap: 12px;

        .batch-buttons {
          width: 100%;

          .el-button {
            flex: 1;
          }
        }
      }
    }
  }
}
</style>