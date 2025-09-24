<template>
  <div class="menu-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><MenuIcon /></el-icon>
          菜单管理
        </h2>
        <p class="page-desc">管理系统菜单结构、权限配置和路由设置</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd" size="default">
          新建菜单
        </el-button>
        <el-button 
          :icon="RefreshRight" 
          @click="refreshMenuTree"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
        <el-button 
          :icon="isTreeView ? List : Operation" 
          @click="toggleView"
          size="default"
        >
          {{ isTreeView ? '列表视图' : '树形视图' }}
        </el-button>
      </div>
    </div>

    <!-- 菜单树形表格 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="table-header">
          <div class="table-title">
            <el-icon class="title-icon"><List /></el-icon>
            <span>菜单列表</span>
            <el-tag v-if="!isTreeView" type="info" size="small" class="total-count">
              共 {{ pagination.total }} 条
            </el-tag>
            <el-tag v-else type="info" size="small" class="total-count">
              树形视图
            </el-tag>
          </div>
          <div class="table-actions">
            <el-tooltip content="展开/收起所有" placement="top">
              <el-button :disabled="!isTreeView"
                size="small"
                :icon="expandAll ? 'Minus' : 'Plus'"
                @click="toggleExpandAll"
                circle
              />
            </el-tooltip>
            <el-tooltip content="刷新数据" placement="top">
              <el-button
                size="small"
                :icon="RefreshRight"
                @click="refreshMenuTree"
                :loading="loading"
                circle
              />
            </el-tooltip>
          </div>
        </div>
      </template>
      <el-table
        ref="tableRef"
        v-loading="loading"
        :data="menuList"
        row-key="id"
        :default-expand-all="false"
        :tree-props="{ children: 'children' }"
        stripe
        border
        style="width: 100%"
        :header-row-style="{ backgroundColor: '#f8f9fa' }"
      >
        <!-- 菜单名称列 -->
        <el-table-column prop="name" label="菜单名称" min-width="280" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="menu-name-cell">
              <div class="menu-icon-wrapper">
                <el-icon v-if="row.icon" class="menu-icon" :class="`icon-${row.type}`">
                  <component :is="getIconComponent(row.icon)" />
                </el-icon>
                <el-icon v-else class="menu-icon default-icon" :class="`icon-${row.type}`">
                  <Folder v-if="row.type === 'directory'" />
                  <Document v-else-if="row.type === 'menu'" />
                  <Key v-else />
                </el-icon>
              </div>
              <div class="menu-info">
                <span class="menu-name">{{ row.name }}</span>
                <div class="menu-meta">
                  <span v-if="row.code" class="menu-code">{{ row.code }}</span>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <!-- 类型列 -->
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="getMenuTypeTagType(row.type)" 
              size="small"
              effect="light"
            >
              {{ getMenuTypeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <!-- 路由路径列 -->
        <el-table-column prop="path" label="路由路径" width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <code v-if="row.path" class="path-code">{{ row.path }}</code>
            <span v-else class="empty-value">-</span>
          </template>
        </el-table-column>
        
        <!-- 组件路径列 -->
        <el-table-column prop="component" label="组件" width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <code v-if="row.component" class="component-code">{{ row.component }}</code>
            <span v-else class="empty-value">-</span>
          </template>
        </el-table-column>
        
        <!-- 权限标识列 -->
        <el-table-column prop="resource" label="权限标识" width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span v-if="row.resource" class="permission-code">{{ row.resource }}:{{ row.action || '*' }}</span>
            <span v-else class="empty-value">-</span>
          </template>
        </el-table-column>
        
        <!-- 排序列 -->
        <el-table-column prop="sortOrder" label="排序" width="80" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small" effect="plain">{{ row.sortOrder }}</el-tag>
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
              :disabled="row.isSystem && row.status === 2"
              size="small"
            />
          </template>
        </el-table-column>
        
        <!-- 显示/隐藏列 -->
        <el-table-column prop="isHidden" label="可见性" width="80" align="center">
          <template #default="{ row }">
            <el-tag 
              v-if="row.isHidden"
              type="warning" 
              size="small"
              effect="light"
            >
              隐藏
            </el-tag>
            <el-tag 
              v-else 
              type="success" 
              size="small"
              effect="light"
            >
              显示
            </el-tag>
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
              <el-tooltip content="添加子菜单" placement="top">
                <el-button
                  type="primary"
                  link
                  size="small"
                  :icon="Plus"
                  @click="handleAdd(row)"
                  v-if="row.type !== 'permission'"
                />
              </el-tooltip>
              <el-tooltip content="编辑菜单" placement="top">
                <el-button
                  type="warning"
                  link
                  size="small"
                  :icon="Edit"
                  @click="handleEdit(row)"
                />
              </el-tooltip>
              <el-tooltip content="复制菜单" placement="top">
                <el-button
                  type="info"
                  link
                  size="small"
                  :icon="CopyDocument"
                  @click="handleCopy(row)"
                />
              </el-tooltip>
              <el-tooltip content="删除菜单" placement="top">
                <el-button
                  type="danger"
                  link
                  size="small"
                  :icon="Delete"
                  @click="handleDelete(row)"
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页器 -->
      <div class="pagination-wrapper" v-if="!isTreeView && pagination.total > 0">
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

    <!-- 菜单表单弹窗 -->
    <MenuForm
      v-model:visible="formVisible"
      :form-data="formData"
      :parent-options="parentOptions"
      @success="handleFormSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  RefreshRight, 
  Search, 
  Edit, 
  Delete, 
  CopyDocument,
  Menu as MenuIcon,
  List,
  InfoFilled,
  Check,
  Close,
  Folder,
  Document,
  Key,
  Operation
} from '@element-plus/icons-vue'
import { menuApi } from '@/api'
import type { Menu, MenuQueryParams } from '@/api/types'
import MenuForm from './components/MenuForm.vue'
import { formatDateTime } from '@/utils/date'

const loading = ref(false)
const menuList = ref<Menu[]>([])
const formVisible = ref(false)
const formData = ref<Partial<Menu>>({})
const parentOptions = ref<Menu[]>([])
const isTreeView = ref(true) // 默认树形视图
const expandAll = ref(false) // 展开状态
const tableRef = ref() // 表格引用

// 搜索表单
const searchForm = reactive<Partial<MenuQueryParams & { type: string }>>({
  name: '',
  status: 0,
  type: '',
  page: 1,
  page_size: 20
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 计算属性
const flatMenuList = computed(() => {
  const flatten = (menus: Menu[]): Menu[] => {
    let result: Menu[] = []
    menus.forEach(menu => {
      result.push(menu)
      if (menu.children?.length) {
        result = result.concat(flatten(menu.children))
      }
    })
    return result
  }
  return flatten(menuList.value)
})

// 切换视图模式
const toggleView = () => {
  isTreeView.value = !isTreeView.value
  fetchMenus()
}

// 切换展开/收起所有
const toggleExpandAll = () => {
  expandAll.value = !expandAll.value

  // 使用 nextTick 确保 DOM 更新后再操作表格
  nextTick(() => {
    if (tableRef.value) {
      if (expandAll.value) {
        // 展开所有行
        flatMenuList.value.forEach(row => {
          if (row.children && row.children.length > 0) {
            tableRef.value.toggleRowExpansion(row, true)
          }
        })
      } else {
        // 收起所有行
        flatMenuList.value.forEach(row => {
          if (row.children && row.children.length > 0) {
            tableRef.value.toggleRowExpansion(row, false)
          }
        })
      }
    }
  })
}

// 获取菜单列表
const fetchMenus = async () => {
  try {
    loading.value = true
    if (isTreeView.value) {
      // 树形视图
      const { data } = await menuApi.getMenuTree()
      menuList.value = data
    } else {
      // 分页列表视图
      const params = {
        page: pagination.page,
        page_size: pagination.pageSize,
        name: searchForm.name,
        status: searchForm.status === 0 ? undefined : searchForm.status
      }
      const { data } = await menuApi.getMenus(params)
      menuList.value = data.list
      pagination.total = data.total
    }
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('获取菜单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取父菜单选项
const fetchParentOptions = async () => {
  try {
    const { data } = await menuApi.getMenuTree()
    parentOptions.value = data.filter(menu => menu.type !== 'permission')
  } catch (error) {
    // 错误信息已在响应拦截器中处理
    console.error('获取父菜单选项失败:', error)
  }
}

// 刷新菜单树
const refreshMenuTree = () => {
  fetchMenus()
  fetchParentOptions()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchMenus()
}

// 重置搜索
const handleResetSearch = () => {
  Object.assign(searchForm, {
    name: '',
    status: 0,
    type: '',
    page: 1,
    page_size: 20
  })
  handleSearch()
}

// 添加菜单
const handleAdd = (parent?: Menu) => {
  formData.value = {
    parentId: parent?.id || 0,
    status: 1,
    sortOrder: 0,
    isHidden: false,
    type: 'menu'
  }
  formVisible.value = true
}

// 编辑菜单
const handleEdit = (row: Menu) => {
  formData.value = { ...row }
  formVisible.value = true
}

// 复制菜单
const handleCopy = (row: Menu) => {
  const copyData = { ...row }
  delete copyData.id
  copyData.name = `${copyData.name}_copy`
  copyData.code = `${copyData.code}_copy`
  formData.value = copyData
  formVisible.value = true
}

// 删除菜单
const handleDelete = async (row: Menu) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜单"${row.name}"吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )
    
    await menuApi.deleteMenu(row.id)
    ElMessage.success('删除成功')
    refreshMenuTree()
  } catch (error: any) {
    if (error !== 'cancel') {
      // 错误信息已在响应拦截器中处理
      console.error('删除菜单失败:', error)
    }
  }
}

// 状态改变
const handleStatusChange = async (row: Menu) => {
  const oldStatus = row.status
  try {
    await menuApi.updateMenuStatus(row.id, { status: row.status })
    ElMessage.success(`菜单已${row.status === 1 ? '启用' : '禁用'}`)
  } catch (error) {
    // 恢复状态
    row.status = oldStatus
    // 错误信息已在响应拦截器中处理
    console.error('状态更新失败:', error)
  }
}

// 分页相关
const handlePageSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchMenus()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchMenus()
}

// 表单成功回调
const handleFormSuccess = () => {
  refreshMenuTree()
}

// 获取菜单类型标签类型
const getMenuTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    directory: 'info',
    menu: 'success',
    permission: 'warning'
  }
  return typeMap[type] || 'info'
}

// 获取菜单类型标签
const getMenuTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    directory: '目录',
    menu: '菜单',
    permission: '权限'
  }
  return typeMap[type] || type
}

// 获取图标组件
const getIconComponent = (icon: string) => {
  // 动态导入 Element Plus 图标
  const iconMap: Record<string, any> = {
    MenuIcon,
    Folder,
    Document,
    Key,
    Plus,
    Edit,
    Delete,
    Operation,
    Setting: 'Setting',
    User: 'User',
    Lock: 'Lock',
    Monitor: 'Monitor',
    DataLine: 'DataLine'
  }
  return iconMap[icon] || MenuIcon
}

// 初始化
onMounted(() => {
  refreshMenuTree()
})
</script>

<style lang="scss" scoped>
.menu-management {
  padding: 24px;
  background-color: #f5f6fa;
  height: 100%;
  overflow-y: auto;

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

        .el-form--inline .el-form-item {
          margin-bottom: 0;
          margin-right: 0;
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

      // 菜单名称单元格
      .menu-name-cell {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 8px 0;

        .menu-icon-wrapper {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 36px;
          height: 36px;
          border-radius: 10px;
          background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
          border: 1px solid #dee2e6;
          transition: all 0.2s ease;

          &:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
          }

          .menu-icon {
            font-size: 18px;
            
            &.icon-directory {
              color: #ff9800;
              background: radial-gradient(circle, #fff3e0 0%, transparent 70%);
            }
            
            &.icon-menu {
              color: #4caf50;
              background: radial-gradient(circle, #e8f5e8 0%, transparent 70%);
            }
            
            &.icon-permission {
              color: #f44336;
              background: radial-gradient(circle, #ffebee 0%, transparent 70%);
            }

            &.default-icon {
              opacity: 0.7;
              color: #9e9e9e;
            }
          }
        }

        .menu-info {
          flex: 1;
          min-width: 0;

          .menu-name {
            font-weight: 500;
            color: #303133;
            font-size: 14px;
            display: block;
          }

          .menu-meta {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-top: 4px;

            .menu-code {
              font-size: 12px;
              color: #909399;
              font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            }
          }
        }
      }

      // 代码样式
      .path-code,
      .component-code {
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 12px;
        background: #f1f3f4;
        padding: 4px 8px;
        border-radius: 4px;
        color: #5f6368;
        display: inline-block;
        max-width: 180px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .permission-code {
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 12px;
        color: #e91e63;
        background: #fce4ec;
        padding: 4px 8px;
        border-radius: 4px;
        display: inline-block;
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
  .menu-management {
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
  .menu-management {
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
          .el-form--inline {
            .el-form-item {
              width: 100%;
              margin-bottom: 12px;
              margin-right: 0;

              .el-input,
              .el-select {
                width: 100% !important;
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
        min-width: 800px;
      }
    }
  }
}

// 深色模式支持（可选）
@media (prefers-color-scheme: dark) {
  .menu-management {
    background-color: #1a1a1a;

    .page-header {
      background: linear-gradient(135deg, #434190 0%, #5a4d7a 100%);
    }

    .toolbar-card,
    .table-card {
      background: #2d2d2d;
      border-color: #404040;
    }
  }
}
</style>