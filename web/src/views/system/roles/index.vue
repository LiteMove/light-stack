<template>
  <div class="role-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><UserFilled /></el-icon>
          角色管理
        </h2>
        <p class="page-desc">管理系统角色、权限分配和访问控制</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd" size="default">
          新建角色
        </el-button>
        <el-button
          :icon="RefreshRight"
          @click="refreshRoles"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
        <el-button
          :icon="Download"
          @click="exportRoles"
          size="default"
        >
          导出角色
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
                placeholder="搜索角色名称或编码"
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
                placeholder="状态筛选"
                clearable
                @change="handleSearch"
                style="width: 140px"
              >
                <el-option label="全部" :value="0" />
                <el-option label="正常" :value="1" />
                <el-option label="禁用" :value="2" />
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
        <div class="batch-section" v-show="selectedRows.length > 0">
          <div class="batch-info">
            <el-icon class="info-icon"><InfoFilled /></el-icon>
            已选择 <strong class="selected-count">{{ selectedRows.length }}</strong> 个角色
          </div>
          <div class="batch-actions">
            <el-button type="success" size="small" :icon="Check" @click="batchEnable">
              批量启用
            </el-button>
            <el-button type="warning" size="small" :icon="Close" @click="batchDisable">
              批量禁用
            </el-button>
            <el-button type="danger" size="small" :icon="Delete" @click="batchDelete">
              批量删除
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 角色列表表格 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="table-header">
          <div class="table-title">
            <el-icon class="title-icon"><List /></el-icon>
            <span>角色列表</span>
            <el-tag type="info" size="small" class="total-count">
              共 {{ pagination.total }} 个角色
            </el-tag>
          </div>
          <div class="table-actions">
            <el-tooltip content="刷新数据" placement="top">
              <el-button
                size="small"
                :icon="RefreshRight"
                @click="refreshRoles"
                :loading="loading"
                circle
              />
            </el-tooltip>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="roleList"
        @selection-change="handleSelectionChange"
        stripe
        border
        style="width: 100%"
        :header-row-style="{ backgroundColor: '#f8f9fa' }"
      >
        <el-table-column type="selection" width="50" align="center" />
        
        <!-- 角色信息列 -->
        <el-table-column prop="name" label="角色信息" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="role-info-cell">
              <div class="role-avatar">
                <el-avatar 
                  :size="40"
                  :style="{ backgroundColor: getRoleColor(row.name) }"
                >
                  {{ row.name?.charAt(0) }}
                </el-avatar>
              </div>
              <div class="role-details">
                <div class="role-name">
                  <span class="name">{{ row.name }}</span>
                  <el-tag v-if="row.isSystem" type="danger" size="small" effect="plain">
                    系统角色
                  </el-tag>
                </div>
                <div class="role-meta">
                  <span class="code">{{ row.code }}</span>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <!-- 描述列 -->
        <el-table-column prop="description" label="描述" width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="description-text">{{ row.description || '-' }}</span>
          </template>
        </el-table-column>
        
        <!-- 状态列 -->
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="2"
              @change="handleStatusChange(row)"
              :disabled="row.isSystem"
              size="small"
            />
          </template>
        </el-table-column>
        
        <!-- 排序列 -->
        <el-table-column prop="sortOrder" label="排序" width="80" align="center">
          <template #default="{ row }">
            <span class="sort-order">{{ row.sortOrder }}</span>
          </template>
        </el-table-column>
        
        <!-- 创建时间列 -->
        <el-table-column prop="createdAt" label="创建时间" width="160" align="center">
          <template #default="{ row }">
            <span class="time-text">{{ formatDateTime(row.createdAt) }}</span>
          </template>
        </el-table-column>
        
        <!-- 操作列 -->
        <el-table-column label="操作" width="280" align="center" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip 
                :content="row.code === 'super_admin' ? '超级管理员角色禁止编辑' : '编辑角色'" 
                placement="top"
              >
                <el-button
                  type="primary"
                  link
                  size="small"
                  :icon="Edit"
                  @click="handleEdit(row)"
                  :disabled="row.code === 'super_admin'"
                />
              </el-tooltip>

              <el-tooltip 
                :content="row.code === 'super_admin' ? '超级管理员角色禁止分配菜单' : '菜单权限'" 
                placement="top"
              >
                <el-button
                  type="info"
                  link
                  size="small"
                  :icon="Menu"
                  @click="handleMenus(row)"
                  :disabled="row.code === 'super_admin'"
                />
              </el-tooltip>
              <el-tooltip content="复制角色" placement="top">
                <el-button
                  type="success"
                  link
                  size="small"
                  :icon="CopyDocument"
                  @click="handleCopy(row)"
                />
              </el-tooltip>
              <el-tooltip 
                :content="row.code === 'super_admin' ? '超级管理员角色禁止删除' : '删除角色'" 
                placement="top"
              >
                <el-button
                  type="danger"
                  link
                  size="small"
                  :icon="Delete"
                  @click="handleDelete(row)"
                  :disabled="row.isSystem || row.code === 'super_admin'"
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页器 -->
      <div class="pagination-wrapper" v-if="pagination.total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>

    <!-- 角色表单弹窗 -->
    <RoleForm
      v-model:visible="formVisible"
      :form-data="formData"
      @success="handleFormSuccess"
    />



    <!-- 菜单权限分配弹窗 -->
    <MenuAssignDialog
      v-model:visible="menuAssignVisible"
      :role-data="selectedRole"
      @success="handleMenuAssignSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  RefreshRight, 
  Search, 
  Edit, 
  Delete, 
  UserFilled,
  List,
  InfoFilled,
  Check,
  Close,
  Download,
  Menu,
  CopyDocument
} from '@element-plus/icons-vue'
import { roleApi } from '@/api'
import type { Role, RoleQueryParams } from '@/api/types'
import RoleForm from './components/RoleForm.vue'

import MenuAssignDialog from './components/MenuAssignDialog.vue'
import { formatDateTime } from '@/utils/date'

// 响应式数据
const loading = ref(false)
const roleList = ref<Role[]>([])
const selectedRows = ref<Role[]>([])
const formVisible = ref(false)

const menuAssignVisible = ref(false)
const formData = ref<Partial<Role>>({})
const selectedRole = ref<Role | null>(null)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: 0,
  page: 1,
  pageSize: 20
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 获取角色头像颜色
const getRoleColor = (name: string): string => {
  const colors = ['#f56a00', '#7265e6', '#ffbf00', '#00a2ae', '#f56565', '#38a169']
  const index = name?.charCodeAt(0) % colors.length || 0
  return colors[index]
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    loading.value = true
    const params: RoleQueryParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status === 0 ? undefined : searchForm.status
    }
    
    const { data } = await roleApi.getRoles(params)
    roleList.value = data.list
    pagination.total = data.total
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('获取角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 刷新角色列表
const refreshRoles = () => {
  fetchRoles()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchRoles()
}

// 重置搜索
const handleResetSearch = () => {
  Object.assign(searchForm, {
    keyword: '',
    status: 0,
    page: 1,
    pageSize: 20
  })
  handleSearch()
}

// 添加角色
const handleAdd = () => {
  formData.value = {
    status: 1,
    isSystem: false,
    sortOrder: 100
  }
  formVisible.value = true
}

// 编辑角色
const handleEdit = (row: Role) => {
  formData.value = { ...row }
  formVisible.value = true
}

// 复制角色
const handleCopy = (row: Role) => {
  formData.value = {
    ...row,
    id: undefined,
    name: `${row.name}_副本`,
    code: `${row.code}_copy`,
    isSystem: false
  }
  formVisible.value = true
}

// 删除角色
const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除角色"${row.name}"吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )
    
    await roleApi.deleteRole(row.id)
    ElMessage.success('删除成功')
    refreshRoles()
  } catch (error: any) {
    if (error !== 'cancel') {
      // 错误信息已在响应拦截器中处理
      console.error('删除角色失败:', error)
    }
  }
}

// 状态改变
const handleStatusChange = async (row: Role) => {
  const oldStatus = row.status
  try {
    await roleApi.updateRoleStatus(row.id, { status: row.status })
    ElMessage.success(`角色已${row.status === 1 ? '启用' : '禁用'}`)
  } catch (error) {
    // 恢复状态
    row.status = oldStatus
    // 错误信息已在响应拦截器中处理
    console.error('状态更新失败:', error)
  }
}



// 菜单权限管理
const handleMenus = (row: Role) => {
  selectedRole.value = row
  menuAssignVisible.value = true
}

// 导出角色
const exportRoles = async () => {
  try {
    ElMessage.info('导出功能开发中...')
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('导出失败:', error)
  }
}

// 选择变化
const handleSelectionChange = (selection: Role[]) => {
  selectedRows.value = selection
}

// 批量启用
const batchEnable = async () => {
  const ids = selectedRows.value.map(row => row.id)
  try {
    await roleApi.batchUpdateRoleStatus({ ids, status: 1 })
    ElMessage.success('批量启用成功')
    refreshRoles()
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('批量启用失败:', error)
  }
}

// 批量禁用
const batchDisable = async () => {
  const ids = selectedRows.value.map(row => row.id)
  try {
    await roleApi.batchUpdateRoleStatus({ ids, status: 2 })
    ElMessage.success('批量禁用成功')
    refreshRoles()
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('批量禁用失败:', error)
  }
}

// 批量删除
const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个角色吗？此操作不可恢复！`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const promises = selectedRows.value.map(row => roleApi.deleteRole(row.id))
    await Promise.all(promises)
    ElMessage.success('批量删除成功')
    refreshRoles()
  } catch (error: any) {
    if (error !== 'cancel') {
      // 错误信息已在响应拦截器中处理
      console.error('批量删除失败:', error)
    }
  }
}

// 分页相关
const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchRoles()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchRoles()
}

// 表单成功回调
const handleFormSuccess = () => {
  refreshRoles()
}



// 菜单分配成功回调
const handleMenuAssignSuccess = () => {
  ElMessage.success('菜单权限分配成功')
}

// 初始化
onMounted(() => {
  refreshRoles()
})
</script>

<style lang="scss" scoped>
.role-management {
  padding: 24px;
  background-color: #f5f6fa;
  min-height: calc(100vh - 64px);

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
    overflow: hidden;

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

      // 角色信息单元格
      .role-info-cell {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 8px 0;

        .role-avatar {
          flex-shrink: 0;
        }

        .role-details {
          flex: 1;
          min-width: 0;

          .role-name {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 4px;

            .name {
              font-weight: 500;
              color: #303133;
              font-size: 14px;
            }
          }

          .role-meta {
            .code {
              font-size: 12px;
              color: #909399;
              font-family: 'Monaco', 'Consolas', monospace;
            }
          }
        }
      }

      .description-text {
        color: #606266;
        font-size: 13px;
      }

      .sort-order {
        font-weight: 500;
        color: #606266;
      }

      .time-text {
        font-size: 12px;
        color: #909399;
      }

      // 操作按钮
      .action-buttons {
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
}

// 响应式设计优化
@media (max-width: 1200px) {
  .role-management {
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
  .role-management {
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
</style>ev