<template>
  <div class="menu-management">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>菜单管理</h2>
      <div class="header-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd">
          添加菜单
        </el-button>
        <el-button 
          type="success" 
          :icon="RefreshRight" 
          @click="refreshMenuTree"
          :loading="loading"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 搜索区域 -->
    <el-card class="search-card" shadow="never">
      <el-form :model="searchForm" inline>
        <el-form-item label="菜单名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入菜单名称"
            clearable
            @keyup.enter="handleSearch"
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
            <el-option label="全部" :value="0" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">
            搜索
          </el-button>
          <el-button @click="handleResetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 菜单表格 -->
    <el-card class="table-card" shadow="never">
      <!-- 批量操作 -->
      <div class="batch-actions" v-if="selectedRows.length > 0">
        <span class="selected-info">
          已选择 <strong>{{ selectedRows.length }}</strong> 项
        </span>
        <el-button type="success" size="small" @click="batchEnable">
          批量启用
        </el-button>
        <el-button type="warning" size="small" @click="batchDisable">
          批量禁用
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table
        v-loading="loading"
        :data="menuList"
        row-key="id"
        default-expand-all
        :tree-props="{ children: 'children' }"
        @selection-change="handleSelectionChange"
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="菜单名称" min-width="200">
          <template #default="{ row }">
            <div class="menu-name-cell">
              <el-icon v-if="row.icon" class="menu-icon">
                <component :is="getIconComponent(row.icon)" />
              </el-icon>
              <span>{{ row.name }}</span>
              <el-tag v-if="row.is_system" type="info" size="small">系统</el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="code" label="菜单编码" width="180" />
        
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getMenuTypeTagType(row.type)" size="small">
              {{ getMenuTypeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="path" label="路由路径" width="200" />
        
        <el-table-column prop="component" label="组件路径" width="200">
          <template #default="{ row }">
            <span class="component-path">{{ row.component || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="sort_order" label="排序" width="80" />
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="2"
              @change="handleStatusChange(row)"
              :disabled="row.is_system && row.status === 2"
            />
          </template>
        </el-table-column>
        
        <el-table-column prop="is_hidden" label="隐藏" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_hidden" type="warning" size="small">隐藏</el-tag>
            <el-tag v-else type="success" size="small">显示</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleAdd(row)"
              v-if="row.type !== 'permission'"
            >
              添加子菜单
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
              :disabled="row.is_system"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="!isTreeView">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handleCurrentChange"
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
import { Plus, RefreshRight, Search } from '@element-plus/icons-vue'
import { menuApi } from '@/api'
import type { Menu, MenuQueryParams } from '@/api/types'
import MenuForm from './components/MenuForm.vue'
import { formatDateTime } from '@/utils/date'

// 响应式数据
const loading = ref(false)
const menuList = ref<Menu[]>([])
const selectedRows = ref<Menu[]>([])
const formVisible = ref(false)
const formData = ref<Partial<Menu>>({})
const parentOptions = ref<Menu[]>([])
const isTreeView = ref(true) // 默认树形视图

// 搜索表单
const searchForm = reactive<Partial<MenuQueryParams>>({
  name: '',
  status: 0,
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
    ElMessage.error('获取菜单列表失败')
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
    ElMessage.error('获取父菜单选项失败')
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
    page: 1,
    page_size: 20
  })
  handleSearch()
}

// 添加菜单
const handleAdd = (parent?: Menu) => {
  formData.value = {
    parent_id: parent?.id || 0,
    status: 1,
    sort_order: 0,
    is_hidden: false,
    type: 'menu'
  }
  formVisible.value = true
}

// 编辑菜单
const handleEdit = (row: Menu) => {
  formData.value = { ...row }
  formVisible.value = true
}

// 删除菜单
const handleDelete = async (row: Menu) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜单"${row.name}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await menuApi.deleteMenu(row.id)
    ElMessage.success('删除成功')
    refreshMenuTree()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 状态改变
const handleStatusChange = async (row: Menu) => {
  try {
    await menuApi.updateMenuStatus(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    // 恢复状态
    row.status = row.status === 1 ? 2 : 1
    ElMessage.error('状态更新失败')
  }
}

// 选择变化
const handleSelectionChange = (selection: Menu[]) => {
  selectedRows.value = selection
}

// 批量启用
const batchEnable = async () => {
  const ids = selectedRows.value.map(row => row.id)
  try {
    await menuApi.batchUpdateMenuStatus({ ids, status: 1 })
    ElMessage.success('批量启用成功')
    refreshMenuTree()
  } catch (error) {
    ElMessage.error('批量启用失败')
  }
}

// 批量禁用
const batchDisable = async () => {
  const ids = selectedRows.value.map(row => row.id)
  try {
    await menuApi.batchUpdateMenuStatus({ ids, status: 2 })
    ElMessage.success('批量禁用成功')
    refreshMenuTree()
  } catch (error) {
    ElMessage.error('批量禁用失败')
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
  // 这里可以根据图标名称返回对应的图标组件
  // 暂时返回一个默认图标
  return 'Menu'
}

// 初始化
onMounted(() => {
  refreshMenuTree()
})
</script>

<style lang="scss" scoped>
.menu-management {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      color: #303133;
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .search-card {
    margin-bottom: 20px;

    :deep(.el-card__body) {
      padding: 16px 20px;
    }
  }

  .table-card {
    :deep(.el-card__body) {
      padding: 0;
    }

    .batch-actions {
      padding: 16px 20px;
      background-color: #f5f7fa;
      border-bottom: 1px solid #ebeef5;
      display: flex;
      align-items: center;
      gap: 12px;

      .selected-info {
        color: #606266;
        margin-right: 12px;
      }
    }

    .el-table {
      border: none;

      .menu-name-cell {
        display: flex;
        align-items: center;
        gap: 8px;

        .menu-icon {
          color: #409eff;
        }
      }

      .component-path {
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 12px;
        color: #909399;
      }
    }

    .pagination-wrapper {
      padding: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .menu-management {
    padding: 12px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;

      .header-actions {
        width: 100%;
        justify-content: flex-end;
      }
    }

    .search-card :deep(.el-form--inline) {
      .el-form-item {
        width: 100%;
        margin-bottom: 12px;

        .el-input,
        .el-select {
          width: 100% !important;
        }
      }
    }
  }
}
</style>