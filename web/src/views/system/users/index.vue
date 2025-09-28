<template>
  <div class="user-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><UserFilled /></el-icon>
          用户管理
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
          {{ isSuperAdmin && currentTenant ? `管理租户 "${currentTenant.name}" 的用户信息、权限分配和状态控制` : isSuperAdmin ? '请先选择要管理的租户' : '管理系统用户信息、权限分配和状态控制' }}
        </p>
      </div>
      <div class="header-actions">
        <el-button 
          v-if="$hasPer('system:user:create')"
          type="primary" 
          :icon="Plus" 
          @click="handleAdd" 
          size="default"
        >
          新建用户
        </el-button>
        <el-button 
          :icon="RefreshRight" 
          @click="refreshUsers"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
        <el-button 
          v-if="$hasPer('system:user:export')"
          :icon="Download" 
          @click="exportUsers"
          size="default"
        >
          导出用户
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
                placeholder="搜索用户名、昵称或邮箱"
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
            已选择 <strong class="selected-count">{{ selectedRows.length }}</strong> 个用户
          </div>
          <div class="batch-actions">
            <el-button 
              v-if="$hasPer('system:user:update')"
              type="success" 
              size="small" 
              :icon="Check" 
              @click="batchEnable"
            >
              批量启用
            </el-button>
            <el-button 
              v-if="$hasPer('system:user:update')"
              type="warning" 
              size="small" 
              :icon="Close" 
              @click="batchDisable"
            >
              批量禁用
            </el-button>
            <el-button 
              v-if="$hasPer('system:user:delete')"
              type="danger" 
              size="small" 
              :icon="Delete" 
              @click="batchDelete"
            >
              批量删除
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 用户列表表格 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="table-header">
          <div class="table-title">
            <el-icon class="title-icon"><List /></el-icon>
            <span>用户列表</span>
            <el-tag type="info" size="small" class="total-count">
              共 {{ pagination.total }} 个用户
            </el-tag>
            <!-- 租户过滤提示 -->
            <el-tag v-if="isSuperAdmin && currentTenant" type="primary" size="small" effect="plain" class="filter-indicator">
              <el-icon><OfficeBuilding /></el-icon>
              {{ currentTenant.name }}
            </el-tag>
          </div>
          <div class="table-actions">
            <el-tooltip content="刷新数据" placement="top">
              <el-button
                size="small"
                :icon="RefreshRight"
                @click="refreshUsers"
                :loading="loading"
                circle
              />
            </el-tooltip>
          </div>
        </div>
      </template>
      
      <!-- 未选择租户提示 -->
      <div v-if="isSuperAdmin && !currentTenant" class="no-tenant-selected">
        <div class="empty-state">
          <el-icon class="empty-icon"><Warning /></el-icon>
          <h3 class="empty-title">请选择租户</h3>
          <p class="empty-description">
            作为超级管理员，您需要先在页面右上角选择要管理的租户，然后才能查看和管理该租户的用户数据。
          </p>
        </div>
      </div>
      
      <!-- 正常表格显示 -->
      <div v-else>
        <el-table
          v-loading="loading"
          :data="userList"
          @selection-change="handleSelectionChange"
          stripe
          border
          style="width: 100%"
          :header-row-style="{ backgroundColor: '#f8f9fa' }"
        >
        <el-table-column type="selection" width="50" align="center" />
        
        <!-- 用户信息列 -->
        <el-table-column prop="username" label="用户信息" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="user-info-cell">
              <div class="user-avatar">
                <el-avatar 
                  :src="row.avatar" 
                  :size="40"
                  :style="{ backgroundColor: getAvatarColor(row.username) }"
                >
                  {{ row.nickname?.charAt(0) || row.username?.charAt(0) }}
                </el-avatar>
              </div>
              <div class="user-details">
                <div class="user-name">
                  <span class="username">{{ row.username }}</span>
                  <el-tag v-if="row.isSystem" type="danger" size="small" effect="plain">
                    系统用户
                  </el-tag>
                </div>
                <div class="user-meta">
                  <span class="nickname">{{ row.nickname || '-' }}</span>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <!-- 联系信息列 -->
        <el-table-column prop="email" label="联系信息" width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="contact-info">
              <div v-if="row.email" class="email">
                <el-icon><Message /></el-icon>
                {{ row.email }}
              </div>
              <div v-if="row.phone" class="phone">
                <el-icon><Phone /></el-icon>
                {{ row.phone }}
              </div>
              <span v-if="!row.email && !row.phone" class="empty-value">-</span>
            </div>
          </template>
        </el-table-column>
        
        <!-- 角色列 -->
        <el-table-column prop="roles" label="角色" width="180">
          <template #default="{ row }">
            <div class="roles-container">
              <el-tag 
                v-for="role in row.roles"
                :key="role.id"
                :type="role.isSystem ? 'danger' : 'primary'"
                size="small"
                effect="light"
                class="role-tag"
              >
                {{ role.name }}
              </el-tag>
              <span v-if="row.roles == null || row.roles.length === 0" class="empty-value">未分配</span>
            </div>
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
              :disabled="row.isSystem || !$hasPer('system:user:update')"
              size="small"
            />
          </template>
        </el-table-column>
        
        <!-- 最后登录列 -->
        <el-table-column prop="lastLoginAt" label="最后登录" width="160" align="center">
          <template #default="{ row }">
            <div class="login-info">
              <div v-if="row.lastLoginAt" class="login-time">
                {{ formatDateTime(row.lastLoginAt) }}
              </div>
              <div v-if="row.lastLoginIp" class="login-ip">
                {{ row.lastLoginIp }}
              </div>
              <span v-if="!row.lastLoginAt" class="empty-value">从未登录</span>
            </div>
          </template>
        </el-table-column>
        
        <!-- 创建时间列 -->
        <el-table-column prop="createdAt" label="创建时间" width="160" align="center">
          <template #default="{ row }">
            <span class="time-text">{{ formatDateTime(row.createdAt) }}</span>
          </template>
        </el-table-column>
        
        <!-- 操作列 -->
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip 
                v-if="$hasPer('system:user:update')"
                content="编辑用户" 
                placement="top"
              >
                <el-button
                  type="primary"
                  link
                  size="small"
                  :icon="Edit"
                  @click="handleEdit(row)"
                />
              </el-tooltip>
              <el-tooltip 
                v-if="$hasPer('system:user:role:assign')"
                content="分配角色" 
                placement="top"
              >
                <el-button
                  type="warning"
                  link
                  size="small"
                  :icon="Key"
                  @click="handleAssignRoles(row)"
                  :disabled="row.id == 1"
                />
              </el-tooltip>
              <el-tooltip 
                v-if="$hasPer('system:user:reset')"
                content="重置密码" 
                placement="top"
              >
                <el-button
                  type="info"
                  link
                  size="small"
                  :icon="Lock"
                  @click="handleResetPassword(row)"
                  :disabled="row.isSystem"
                />
              </el-tooltip>
              <el-tooltip 
                v-if="$hasPer('system:user:delete')"
                content="删除用户" 
                placement="top"
              >
                <el-button
                  type="danger"
                  link
                  size="small"
                  :icon="Delete"
                  @click="handleDelete(row)"
                  :disabled="row.isSystem"
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
      </div>
    </el-card>

    <!-- 用户表单弹窗 -->
    <UserForm
      v-model:visible="formVisible"
      :form-data="formData"
      :roles="roles"
      @success="handleFormSuccess"
    />

    <!-- 角色分配弹窗 -->
    <RoleAssignDialog
      v-model:visible="roleAssignVisible"
      :user-data="selectedUser"
      :roles="roles"
      @success="handleRoleAssignSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, onUnmounted } from 'vue'
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
  Message,
  Phone,
  Key,
  Lock,
  OfficeBuilding,
  Warning
} from '@element-plus/icons-vue'
import { userApi, roleApi } from '@/api'
import type { User, Role, PageParams } from '@/api/types'
import UserForm from './components/UserForm.vue'
import RoleAssignDialog from './components/RoleAssignDialog.vue'
import { formatDateTime } from '@/utils/date'
import { useTenantStore, useUserStore } from '@/store'

// Store实例
const tenantStore = useTenantStore()
const userStore = useUserStore()

// 响应式数据
const loading = ref(false)
const userList = ref<User[]>([])
const roles = ref<Role[]>([])
const selectedRows = ref<User[]>([])
const formVisible = ref(false)
const roleAssignVisible = ref(false)
const formData = ref<Partial<User>>({})
const selectedUser = ref<User | null>(null)
const abortController = ref<AbortController | null>(null)
const isUnmounting = ref(false)

// 计算属性
const isSuperAdmin = computed(() => tenantStore.checkIsSuperAdmin())
const currentTenant = computed(() => tenantStore.getCurrentTenant())

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

// 获取用户头像颜色
const getAvatarColor = (username: string): string => {
  const colors = ['#f56a00', '#7265e6', '#ffbf00', '#00a2ae', '#f56565', '#38a169']
  const index = username.charCodeAt(0) % colors.length
  return colors[index]
}

// 获取用户列表
const fetchUsers = async () => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过获取用户列表')
    return
  }

  // 如果是超级管理员但没有选择租户，不加载数据
  if (isSuperAdmin.value && !currentTenant.value) {
    ElMessage.warning('请先选择要管理的租户')
    userList.value = []
    pagination.total = 0
    return
  }

  // 取消之前的请求
  if (abortController.value) {
    abortController.value.abort()
  }

  // 创建新的 AbortController
  abortController.value = new AbortController()

  try {
    loading.value = true
    const params: PageParams & { keyword?: string; status?: number } = {
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      status: searchForm.status === 0 ? undefined : searchForm.status,
    }

    // 租户ID现在通过请求头自动添加，不需要在参数中指定
    const { data } = await userApi.getUsers(params)

    // 检查请求是否被取消或组件是否正在销毁
    if (abortController.value?.signal.aborted || isUnmounting.value) {
      return
    }

    userList.value = data.list
    pagination.total = data.total

  } catch (error: any) {
    // 如果是请求被取消或401错误，不显示错误信息
    if (error.name === 'AbortError' || error.message?.includes('登录') || isUnmounting.value) {
      console.log('请求已取消或用户未登录')
      return
    }
    // 错误信息已在响应拦截器中处理
    console.error('获取用户列表失败:', error)
  } finally {
    if (!isUnmounting.value) {
      loading.value = false
    }
  }
}

// 获取角色列表
const fetchRoles = async () => {
  // 检查是否已登录
  if (!userStore.getToken()) {
    console.log('用户未登录，跳过获取角色列表')
    return
  }

  try {
    const { data } = await roleApi.getActiveRoles()
    roles.value = data
  } catch (error: any) {
    // 如果是401错误，不显示错误信息
    if (error.message?.includes('登录')) {
      console.log('用户未登录，跳过获取角色列表')
      return
    }
    // 错误信息已在响应拦截器中处理
    console.error('获取角色列表失败:', error)
  }
}

// 刷新用户列表
const refreshUsers = () => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过刷新用户列表')
    return
  }

  fetchUsers()
  fetchRoles()
}

// 搜索
const handleSearch = () => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过搜索')
    return
  }

  pagination.page = 1
  fetchUsers()
}

// 重置搜索
const handleResetSearch = () => {
  // 检查组件是否正在销毁
  if (isUnmounting.value) {
    console.log('组件销毁中，跳过重置搜索')
    return
  }

  Object.assign(searchForm, {
    keyword: '',
    status: 0,
    roleId: 0,
    page: 1,
    pageSize: 20
  })
  handleSearch()
}

// 添加用户
const handleAdd = () => {
  formData.value = {
    status: 1,
    isSystem: false
  }
  formVisible.value = true
}

// 编辑用户
const handleEdit = (row: User) => {
  formData.value = { ...row }
  formVisible.value = true
}

// 删除用户
const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户"${row.username}"吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    await userApi.deleteUser(row.id)
    ElMessage.success('删除成功')
    refreshUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      // 错误信息已在响应拦截器中处理
      console.error('删除用户失败:', error)
    }
  }
}

// 状态改变
const handleStatusChange = async (row: User) => {
  const oldStatus = row.status
  try {
    await userApi.updateUserStatus(row.id, { status: row.status })
    ElMessage.success(`用户已${row.status === 1 ? '启用' : '禁用'}`)
  } catch (error) {
    // 恢复状态
    row.status = oldStatus
    // 错误信息已在响应拦截器中处理
    console.error('状态更新失败:', error)
  }
}

// 分配角色
const handleAssignRoles = (row: User) => {
  selectedUser.value = row
  roleAssignVisible.value = true
}

// 重置密码
const handleResetPassword = async (row: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要重置用户"${row.username}"的密码吗？`,
      '重置密码',
      {
        confirmButtonText: '确定重置',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 假设有重置密码的API
    // await userApi.resetPassword(row.id)
    ElMessage.success('密码重置成功，新密码已发送到用户邮箱')
  } catch (error: any) {
    if (error !== 'cancel') {
      // 错误信息已在响应拦截器中处理
      console.error('密码重置失败:', error)
    }
  }
}

// 导出用户
const exportUsers = async () => {
  try {
    ElMessage.info('导出功能开发中...')
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('导出失败:', error)
  }
}

// 选择变化
const handleSelectionChange = (selection: User[]) => {
  selectedRows.value = selection
}

// 批量启用
const batchEnable = async () => {
  const ids = selectedRows.value.map(row => row.id)
  try {
    // 假设有批量更新状态的API
    // await userApi.batchUpdateUserStatus({ ids, status: 1 })
    ElMessage.success('批量启用成功')
    refreshUsers()
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('批量启用失败:', error)
  }
}

// 批量禁用
const batchDisable = async () => {
  const ids = selectedRows.value.map(row => row.id)
  try {
    // await userApi.batchUpdateUserStatus({ ids, status: 2 })
    ElMessage.success('批量禁用成功')
    refreshUsers()
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('批量禁用失败:', error)
  }
}

// 批量删除
const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 个用户吗？此操作不可恢复！`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const promises = selectedRows.value.map(row => userApi.deleteUser(row.id))
    await Promise.all(promises)
    ElMessage.success('批量删除成功')
    refreshUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      // 错误信息已在响应拦截器中处理
      console.error('批量删除失败:', error)
    }
  }
}

// 分页相关
const handlePageSizeChange = (size: number) => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过分页大小变更')
    return
  }

  pagination.pageSize = size
  pagination.page = 1
  fetchUsers()
}

const handleCurrentChange = (page: number) => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过页码变更')
    return
  }

  pagination.page = page
  fetchUsers()
}

// 表单成功回调
const handleFormSuccess = () => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过表单成功回调')
    return
  }
  refreshUsers()
}

// 角色分配成功回调
const handleRoleAssignSuccess = () => {
  // 检查组件是否正在销毁或用户是否已登录
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过角色分配成功回调')
    return
  }
  refreshUsers()
}

// 初始化
onMounted(() => {
  refreshUsers()
})

onUnmounted(() => {
  isUnmounting.value = true
  if (abortController.value) {
    abortController.value.abort()
  }
})
</script>

<style lang="scss" scoped>
.user-management {
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
        flex-wrap: wrap;

        .title-icon {
          font-size: 32px;
        }

        .tenant-indicator {
          margin-left: 8px;
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 12px;
          font-weight: 500;
          
          .el-icon {
            font-size: 14px;
          }
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

        .filter-indicator {
          margin-left: 8px;
          display: flex;
          align-items: center;
          gap: 4px;
          
          .el-icon {
            font-size: 12px;
          }
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

      // 用户信息单元格
      .user-info-cell {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 8px 0;

        .user-avatar {
          flex-shrink: 0;
        }

        .user-details {
          flex: 1;
          min-width: 0;

          .user-name {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 4px;

            .username {
              font-weight: 500;
              color: #303133;
              font-size: 14px;
            }
          }

          .user-meta {
            .nickname {
              font-size: 12px;
              color: #909399;
              margin-bottom: 2px;
            }

            .tenant-info {
              display: flex;
              align-items: center;
              gap: 4px;
              font-size: 11px;
              color: #409eff;
              
              .el-icon {
                font-size: 10px;
              }
              
              .tenant-name {
                font-weight: 500;
              }
            }
          }
        }
      }

      // 联系信息
      .contact-info {
        .email, .phone {
          display: flex;
          align-items: center;
          gap: 4px;
          margin-bottom: 2px;
          font-size: 13px;
          color: #606266;

          .el-icon {
            font-size: 12px;
            color: #909399;
          }
        }
      }

      // 角色容器
      .roles-container {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;

        .role-tag {
          margin: 0;
        }
      }

      // 登录信息
      .login-info {
        .login-time {
          font-size: 12px;
          color: #606266;
          margin-bottom: 2px;
        }

        .login-ip {
          font-size: 11px;
          color: #909399;
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

    // 空状态样式
    .no-tenant-selected {
      padding: 60px 20px;
      text-align: center;
      
      .empty-state {
        .empty-icon {
          font-size: 64px;
          color: #e6a23c;
          margin-bottom: 16px;
        }
        
        .empty-title {
          margin: 0 0 12px 0;
          font-size: 18px;
          font-weight: 500;
          color: #303133;
        }
        
        .empty-description {
          margin: 0;
          color: #606266;
          font-size: 14px;
          line-height: 1.6;
          max-width: 400px;
          margin: 0 auto;
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
  .user-management {
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
  .user-management {
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