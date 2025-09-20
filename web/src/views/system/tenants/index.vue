<template>
  <div class="tenant-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><OfficeBuilding /></el-icon>
          租户管理
        </h2>
        <p class="page-desc">管理多租户系统中的租户信息、域名配置和状态控制</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd" size="default">
          新建租户
        </el-button>
        <el-button
          :icon="RefreshRight"
          @click="refreshTenants"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
        <el-button
          :icon="Download"
          @click="exportTenants"
          size="default"
        >
          导出数据
        </el-button>
      </div>
    </div>

    <!-- 操作工具栏 -->
    <el-card class="toolbar-card" shadow="never">
      <div class="toolbar-content">
        <!-- 搜索区域 -->
        <div class="search-section">
          <el-form :model="searchForm" inline class="search-form">
            <el-form-item label="关键词" class="search-item">
              <el-input
                v-model="searchForm.keyword"
                placeholder="搜索租户名称或域名"
                clearable
                @keyup.enter="handleSearch"
                @clear="handleSearch"
                :prefix-icon="Search"
                style="width: 260px"
              />
            </el-form-item>

            <el-form-item label="状态" class="search-item">
              <el-select
                v-model="searchForm.status"
                placeholder="选择状态"
                clearable
                style="width: 140px"
                @change="handleSearch"
              >
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="2" />
                <el-option label="试用" :value="3" />
                <el-option label="过期" :value="4" />
              </el-select>
            </el-form-item>

            <el-form-item class="search-actions">
              <el-button type="primary" :icon="Search" @click="handleSearch">
                搜索
              </el-button>
              <el-button :icon="RefreshRight" @click="resetSearch">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 批量操作 -->
        <div class="batch-section" v-if="selectedRows.length > 0">
          <div class="batch-info">
            <span>已选中 {{ selectedRows.length }} 项</span>
          </div>
          <div class="batch-actions">
            <el-button
              type="success"
              size="small"
              @click="batchUpdateStatus(1)"
              :icon="Check"
            >
              批量启用
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="batchUpdateStatus(2)"
              :icon="Close"
            >
              批量禁用
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="tenantList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
        row-key="id"
        :default-sort="{ prop: 'created_at', order: 'descending' }"
      >
        <el-table-column type="selection" width="50" />

        <el-table-column prop="id" label="租户ID" width="80" sortable="custom" />

        <el-table-column prop="name" label="租户名称" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="tenant-name">
              <span class="name-text">{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="domain" label="域名" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="domain-info" v-if="row.domain">
              <el-link :href="`http://${row.domain}`" target="_blank" type="primary">
                {{ row.domain }}
              </el-link>
            </div>
            <span v-else class="text-gray">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100" sortable="custom">
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.status)"
              size="small"
              effect="light"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="expired_at" label="过期时间" width="180" sortable="custom">
          <template #default="{ row }">
            <div v-if="row.expired_at">
              <span :class="{ 'text-danger': isExpired(row.expired_at) }">
                {{ formatDateTime(row.expired_at) }}
              </span>
              <el-tag v-if="isExpired(row.expired_at)" type="danger" size="small" class="ml-2">
                已过期
              </el-tag>
            </div>
            <span v-else class="text-gray">永久</span>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <el-button
                type="primary"
                size="small"
                link
                @click="handleEdit(row)"
                :icon="Edit"
              >
                编辑
              </el-button>

              <el-button
                :type="row.status === 1 ? 'warning' : 'success'"
                size="small"
                link
                @click="toggleStatus(row)"
                :icon="row.status === 1 ? 'Close' : 'Check'"
              >
                {{ row.status === 1 ? '禁用' : '启用' }}
              </el-button>

              <el-button
                type="info"
                size="small"
                link
                @click="viewDetail(row)"
                :icon="View"
              >
                详情
              </el-button>

              <el-popconfirm
                title="确定要删除这个租户吗？"
                confirm-button-text="确定"
                cancel-button-text="取消"
                @confirm="handleDelete(row)"
                v-if="row.id !== 0"
              >
                <template #reference>
                  <el-button
                    type="danger"
                    size="small"
                    link
                    :icon="Delete"
                  >
                    删除
                  </el-button>
                </template>
              </el-popconfirm>
            </div>
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
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 租户表单弹窗 -->
    <TenantForm
      v-model:visible="formVisible"
      :tenant-data="currentTenant"
      @success="handleFormSuccess"
    />

    <!-- 租户详情弹窗 -->
    <TenantDetail
      v-model:visible="detailVisible"
      :tenant-data="currentTenant"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  RefreshRight,
  Download,
  Search,
  Edit,
  Delete,
  View,
  Check,
  Close,
  OfficeBuilding
} from '@element-plus/icons-vue'
import TenantForm from './components/TenantForm.vue'
import TenantDetail from './components/TenantDetail.vue'
import { tenantApi } from '@/api/tenant'
import { formatDateTime } from '@/utils/date'

// 响应式数据
const loading = ref(false)
const tenantList = ref([])
const selectedRows = ref([])
const formVisible = ref(false)
const detailVisible = ref(false)
const currentTenant = ref(null)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: null
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
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

// 获取租户列表
const getTenantList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword,
      status: searchForm.status
    }

    const response = await tenantApi.getTenantList(params)

    if (response.data) {
      tenantList.value = response.data.list || []
      pagination.total = response.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取租户列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  pagination.page = 1
  getTenantList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = null
  pagination.page = 1
  getTenantList()
}

// 刷新数据
const refreshTenants = () => {
  getTenantList()
}

// 新增租户
const handleAdd = () => {
  currentTenant.value = null
  formVisible.value = true
}

// 编辑租户
const handleEdit = (row) => {
  currentTenant.value = { ...row }
  formVisible.value = true
}

// 查看详情
const viewDetail = (row) => {
  currentTenant.value = { ...row }
  detailVisible.value = true
}

// 删除租户
const handleDelete = async (row) => {
  try {
    await tenantApi.deleteTenant(row.id)
    ElMessage.success('删除成功')
    getTenantList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

// 切换状态
const toggleStatus = async (row) => {
  const newStatus = row.status === 1 ? 2 : 1
  const statusText = newStatus === 1 ? '启用' : '禁用'

  try {
    await ElMessageBox.confirm(
      `确定要${statusText}租户"${row.name}"吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await tenantApi.updateTenantStatus(row.id, { status: newStatus })
    ElMessage.success(`${statusText}成功`)
    getTenantList()
  } catch (error) {
    if (error !== 'cancel') {
      // 错误信息已经在request拦截器中显示了，这里不需要再显示
      console.error('状态更新失败:', error)
    }
  }
}

// 批量更新状态
const batchUpdateStatus = async (status) => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要操作的租户')
    return
  }

  const statusText = status === 1 ? '启用' : '禁用'

  try {
    await ElMessageBox.confirm(
      `确定要${statusText}选中的 ${selectedRows.value.length} 个租户吗？`,
      '批量操作确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 批量更新
    const promises = selectedRows.value.map(row =>
      tenantApi.updateTenantStatus(row.id, { status })
    )

    await Promise.all(promises)
    ElMessage.success(`批量${statusText}成功`)
    getTenantList()
    selectedRows.value = []
  } catch (error) {
    if (error !== 'cancel') {
      // 错误信息已经在request拦截器中显示了，这里不需要再显示
      console.error('批量状态更新失败:', error)
    }
  }
}

// 导出租户数据
const exportTenants = () => {
  ElMessage.info('导出功能开发中...')
}

// 表格选择变化
const handleSelectionChange = (selection) => {
  selectedRows.value = selection
}

// 排序变化
const handleSortChange = ({ prop, order }) => {
  // 处理排序逻辑
  console.log('排序变化:', prop, order)
}

// 页面大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
  getTenantList()
}

// 当前页变化
const handleCurrentChange = (page) => {
  pagination.page = page
  getTenantList()
}

// 表单提交成功
const handleFormSuccess = () => {
  formVisible.value = false
  getTenantList()
}

// 初始化
onMounted(() => {
  getTenantList()
})
</script>

<style lang="scss" scoped>
.tenant-management {
  padding: 24px;
  background-color: #f5f6fa;
  min-height: calc(100vh - 64px);
}

// 页面头部样式
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

      &.el-button--primary {
        background: rgba(255, 255, 255, 0.2);
        
        &:hover {
          background: rgba(255, 255, 255, 0.3);
        }
      }
    }
  }
}

// 工具栏样式
.toolbar-card {
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e4e7ed;

  :deep(.el-card__body) {
    padding: 20px;
  }

  .toolbar-content {
    .search-section {
      margin-bottom: 16px;

      .search-form {
        display: flex;
        flex-wrap: wrap;
        align-items: flex-end;
        gap: 16px;

        .search-item {
          margin-bottom: 0;
          margin-right: 0;

          :deep(.el-form-item__label) {
            color: #606266;
            font-weight: 500;
            font-size: 13px;
          }
        }

        .search-actions {
          margin-bottom: 0;
          margin-right: 0;
          display: flex;
          gap: 12px;

          .el-button {
            border-radius: 6px;
            transition: all 0.2s ease;

            &:hover {
              transform: translateY(-1px);
            }
          }
        }
      }
    }

    .batch-section {
      padding: 16px;
      background: linear-gradient(90deg, #e3f2fd 0%, #f3e5f5 100%);
      border-radius: 8px;
      border: 1px solid #bbdefb;
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 16px;

      .batch-info {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #1976d2;
        font-weight: 500;

        .info-icon {
          font-size: 18px;
          color: #2196f3;
        }

        .selected-count {
          color: #1565c0;
          font-size: 16px;
        }
      }

      .batch-actions {
        display: flex;
        gap: 8px;
      }
    }
  }
}

// 表格卡片样式
.table-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  border: 1px solid #e4e7ed;

  :deep(.el-card__body) {
    padding: 0;
  }

  :deep(.el-card__header) {
    padding: 16px 20px;
    border-bottom: 1px solid #f0f0f0;
    background: #fafafa;
  }

  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .table-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      color: #303133;

      .title-icon {
        color: #409eff;
        font-size: 16px;
      }

      .total-count {
        margin-left: 8px;
        background: #e3f2fd;
        color: #1976d2;
        border: 1px solid #bbdefb;
      }
    }

    .table-actions {
      display: flex;
      gap: 8px;

      .el-button {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        transition: all 0.2s ease;

        &:hover {
          transform: scale(1.1);
        }
      }
    }
  }

  // 表格样式优化
  .el-table {
    border: none;
    
    // 表头样式
    :deep(.el-table__header-wrapper) {
      .el-table__header {
        th {
          background: #f8f9fa;
          color: #495057;
          font-weight: 600;
          font-size: 14px;
          border-bottom: 2px solid #dee2e6;
        }
      }
    }

    // 行样式
    :deep(.el-table__body) {
      tr {
        transition: all 0.2s ease;
        
        &:hover {
          background-color: #f8f9ff;
        }
      }
    }

    // 租户名称单元格
    .tenant-name {
      display: flex;
      align-items: center;
      gap: 8px;

      .name-text {
        font-weight: 500;
        color: #303133;
        font-size: 14px;
      }
    }

    // 域名信息
    .domain-info {
      font-family: 'Monaco', 'Consolas', monospace;
      font-size: 13px;
    }

    // 操作按钮
    .table-actions {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 4px;

      .el-button {
        padding: 4px;
        width: 28px;
        height: 28px;
        border-radius: 6px;
        
        &:hover {
          transform: translateY(-1px);
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
      }
    }

    .empty-value {
      color: #c0c4cc;
      font-style: italic;
      font-size: 12px;
    }

    .time-text {
      font-size: 12px;
      color: #909399;
    }
  }

  // 分页样式
  .pagination-wrapper {
    padding: 20px;
    background: #fafafa;
    border-top: 1px solid #e4e7ed;
    display: flex;
    justify-content: center;

    :deep(.el-pagination) {
      .el-pager li {
        border-radius: 6px;
        margin: 0 2px;
      }
      
      .btn-prev,
      .btn-next {
        border-radius: 6px;
      }
    }
  }
}

.text-gray {
  color: #9ca3af;
}

.text-danger {
  color: #ef4444;
}

.ml-2 {
  margin-left: 8px;
}

// 响应式设计优化
@media (max-width: 1200px) {
  .tenant-management {
    padding: 16px;

    .page-header {
      flex-direction: column;
      align-items: stretch;
      gap: 16px;

      .header-actions {
        justify-content: flex-end;
      }
    }
  }
}

@media (max-width: 768px) {
  .tenant-management {
    padding: 12px;

    .page-header {
      padding: 20px 16px;
      
      .header-left {
        .page-title {
          font-size: 24px;
        }
      }

      .header-actions {
        flex-wrap: wrap;
        gap: 8px;
        
        .el-button {
          flex: 1;
          min-width: 100px;
        }
      }
    }

    .toolbar-card {
      :deep(.el-card__body) {
        padding: 16px;
      }

      .toolbar-content {
        .search-section {
          .search-form {
            flex-direction: column;
            align-items: stretch;

            .search-item {
              width: 100%;
              margin-bottom: 12px;
              margin-right: 0;

              :deep(.el-input),
              :deep(.el-select) {
                width: 100% !important;
              }
            }

            .search-actions {
              justify-content: center;
              
              .el-button {
                flex: 1;
              }
            }
          }
        }

        .batch-section {
          flex-direction: column;
          align-items: stretch;
          gap: 12px;

          .batch-actions {
            justify-content: center;
          }
        }
      }
    }

    // 移动端表格滚动
    .table-card {
      .el-table {
        min-width: 1000px;
      }
    }
  }
}
</style>