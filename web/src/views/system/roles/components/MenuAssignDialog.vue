<template>
  <el-dialog
    v-model="dialogVisible"
    :title="'为角色「' + roleData?.name + '」分配权限'"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="menu-assign">
      <!-- 操作区域 -->
      <div class="operation-section">
        <div class="operation-left">
          <el-checkbox
            v-model="checkAll"
            :indeterminate="indeterminate"
            @change="handleCheckAllChange"
          >
            全选
          </el-checkbox>
          <span class="selected-info">
            已选择 {{ selectedMenus.length }} 项菜单
          </span>
        </div>
        <div class="operation-right">
          <el-button-group>
            <el-button size="small" @click="expandAll">展开全部</el-button>
            <el-button size="small" @click="collapseAll">收起全部</el-button>
          </el-button-group>
        </div>
      </div>

      <!-- 菜单树 -->
      <div class="menu-tree-container" v-loading="loading">
        <el-tree
          ref="treeRef"
          :data="menuTree"
          show-checkbox
          node-key="id"
          :default-expanded-keys="expandedKeys"
          :default-checked-keys="selectedMenus"
          :props="treeProps"
          :check-strictly="false"
          :expand-on-click-node="false"
          @check="handleNodeCheck"
          class="menu-tree"
        >
          <template #default="{ node, data }">
            <div class="menu-node">
              <div class="node-content">
                <el-icon class="node-icon" v-if="data.icon">
                  <component :is="data.icon" />
                </el-icon>
                <el-icon class="node-icon" v-else>
                  <Folder />
                </el-icon>
                <span class="node-label">{{ data.name }}</span>
                <el-tag
                  v-if="data.isSystem"
                  type="danger"
                  size="small"
                  effect="plain"
                  class="system-tag"
                >
                  系统
                </el-tag>
                <el-tag
                  :type="getTypeTagType(data.type)"
                  size="small"
                  effect="light"
                  class="type-tag"
                >
                  {{ getTypeLabel(data.type) }}
                </el-tag>
              </div>
              <div class="node-meta">
                <span class="node-code">{{ data.code }}</span>
                <span v-if="data.path" class="node-path">{{ data.path }}</span>
              </div>
            </div>
          </template>
        </el-tree>
      </div>

      <!-- 选中的菜单预览 -->
      <div class="selected-preview" v-if="selectedMenus.length > 0">
        <div class="preview-header">
          <span class="preview-title">已选择的菜单 ({{ selectedMenus.length }})</span>
          <el-button
            type="text"
            size="small"
            @click="clearSelection"
          >
            清空选择
          </el-button>
        </div>
        <div class="preview-content">
          <el-tag
            v-for="menuId in selectedMenus"
            :key="menuId"
            closable
            size="small"
            class="menu-tag"
            @close="removeMenu(menuId)"
          >
            {{ getMenuLabel(menuId) }}
          </el-tag>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="submitting">
          取消
        </el-button>
        <el-button 
          type="primary" 
          @click="handleSubmit"
          :loading="submitting"
        >
          保存菜单配置
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage, type ElTree } from 'element-plus'
import { Folder } from '@element-plus/icons-vue'
import { roleApi, menuApi } from '@/api'
import type { Role, Menu } from '@/api/types'

// Props
interface Props {
  visible: boolean
  roleData?: Role | null
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  roleData: null
})

// Emits
const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const treeRef = ref<InstanceType<typeof ElTree>>()
const menuTree = ref<Menu[]>([])
const selectedMenus = ref<number[]>([])
const originalMenus = ref<number[]>([])
const expandedKeys = ref<number[]>([])
const menuMap = ref(new Map<number, Menu>())

// 树形控件配置
const treeProps = {
  children: 'children',
  label: 'name',
  //disabled: (data: Menu) => data.isSystem && !selectedMenus.value.includes(data.id)
}

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 全选状态
const checkAll = computed({
  get: () => {
    const allMenuIds = getAllMenuIds(menuTree.value)
    const availableMenuIds = allMenuIds.filter(id => {
      const menu = menuMap.value.get(id)
      return !menu?.isSystem || selectedMenus.value.includes(id)
    })
    return availableMenuIds.length > 0 && availableMenuIds.every(id => selectedMenus.value.includes(id))
  },
  set: (value: boolean) => {
    handleCheckAllChange(value)
  }
})

// 半选状态
const indeterminate = computed(() => {
  const allMenuIds = getAllMenuIds(menuTree.value)
  const availableMenuIds = allMenuIds.filter(id => {
    const menu = menuMap.value.get(id)
    return !menu?.isSystem || selectedMenus.value.includes(id)
  })
  const selectedCount = availableMenuIds.filter(id => selectedMenus.value.includes(id)).length
  return selectedCount > 0 && selectedCount < availableMenuIds.length
})

// 获取所有菜单ID
const getAllMenuIds = (menus: Menu[]): number[] => {
  const ids: number[] = []
  const traverse = (items: Menu[]) => {
    items.forEach(item => {
      ids.push(item.id)
      if (item.children && item.children.length > 0) {
        traverse(item.children)
      }
    })
  }
  traverse(menus)
  return ids
}

// 构建菜单映射
const buildMenuMap = (menus: Menu[]) => {
  const map = new Map<number, Menu>()
  const traverse = (items: Menu[]) => {
    items.forEach(item => {
      map.set(item.id, item)
      if (item.children && item.children.length > 0) {
        traverse(item.children)
      }
    })
  }
  traverse(menus)
  menuMap.value = map
}

// 获取类型标签
const getTypeLabel = (type: string): string => {
  const typeMap: Record<string, string> = {
    directory: '目录',
    menu: '菜单',
    permission: '权限'
  }
  return typeMap[type] || type
}

// 获取类型标签样式
const getTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    directory: 'info',
    menu: 'primary',
    permission: 'warning'
  }
  return typeMap[type] || 'default'
}

// 获取菜单标签
const getMenuLabel = (menuId: number): string => {
  const menu = menuMap.value.get(menuId)
  return menu ? menu.name : `菜单${menuId}`
}

// 获取菜单树
const fetchMenuTree = async () => {
  try {
    loading.value = true
    const { data } = await menuApi.getMenuTree()
    menuTree.value = data
    buildMenuMap(data)
    
    // 设置默认展开的节点
    const firstLevelIds = data.map(item => item.id)
    expandedKeys.value = firstLevelIds
  } catch (error) {
    ElMessage.error('获取菜单树失败')
    menuTree.value = []
  } finally {
    loading.value = false
  }
}

// 获取菜单的叶子节点（用于回显时只显示用户真正选择的菜单）
const getLeafMenuIds = (allMenuIds: number[]): number[] => {
  const leafIds: number[] = []
  
  allMenuIds.forEach(menuId => {
    const menu = menuMap.value.get(menuId)
    if (!menu) return
    
    // 检查是否有子菜单也在选中列表中
    const hasSelectedChildren = allMenuIds.some(id => {
      const childMenu = menuMap.value.get(id)
      return childMenu && childMenu.parentId === menuId
    })
    
    // 如果没有子菜单被选中，则这是一个叶子节点
    if (!hasSelectedChildren) {
      leafIds.push(menuId)
    }
  })
  
  return leafIds
}

// 获取角色菜单
const fetchRoleMenus = async () => {
  if (!props.roleData?.id) return

  try {
    const { data } = await roleApi.getRoleMenus(props.roleData.id)
    const menuIds = data.map(menu => menu.id)
    
    // 获取叶子节点（用户真正选择的菜单）
    const leafMenuIds = getLeafMenuIds(menuIds)
    
    console.log('从后端获取的所有菜单ID:', menuIds)
    console.log('提取的叶子节点ID（用于显示）:', leafMenuIds)
    
    selectedMenus.value = [...leafMenuIds]
    originalMenus.value = [...leafMenuIds]
    
    // 设置树形控件的选中状态
    if (treeRef.value) {
      treeRef.value.setCheckedKeys(leafMenuIds)
    }
  } catch (error) {
    console.error('获取角色菜单失败:', error)
    selectedMenus.value = []
    originalMenus.value = []
  }
}

// 节点选中状态变化
const handleNodeCheck = (data: Menu, checked: any) => {
  selectedMenus.value = checked.checkedKeys
}

// 全选/取消全选
const handleCheckAllChange = (checked: boolean) => {
  const allMenuIds = getAllMenuIds(menuTree.value)
  const availableMenuIds = allMenuIds.filter(id => {
    const menu = menuMap.value.get(id)
    return !menu?.isSystem || selectedMenus.value.includes(id)
  })

  if (checked) {
    selectedMenus.value = [...new Set([...selectedMenus.value, ...availableMenuIds])]
  } else {
    const systemMenuIds = selectedMenus.value.filter(id => {
      const menu = menuMap.value.get(id)
      return menu?.isSystem
    })
    selectedMenus.value = systemMenuIds
  }

  // 更新树形控件的选中状态
  if (treeRef.value) {
    treeRef.value.setCheckedKeys(selectedMenus.value)
  }
}

// 展开全部
const expandAll = () => {
  const allMenuIds = getAllMenuIds(menuTree.value)
  expandedKeys.value = allMenuIds
  if (treeRef.value) {
    allMenuIds.forEach(id => {
      treeRef.value!.store.nodesMap[id]?.expand()
    })
  }
}

// 收起全部
const collapseAll = () => {
  expandedKeys.value = []
  if (treeRef.value) {
    const allMenuIds = getAllMenuIds(menuTree.value)
    allMenuIds.forEach(id => {
      treeRef.value!.store.nodesMap[id]?.collapse()
    })
  }
}

// 清空选择
const clearSelection = () => {
  selectedMenus.value = []
  if (treeRef.value) {
    treeRef.value.setCheckedKeys([])
  }
}

// 移除单个菜单
const removeMenu = (menuId: number) => {
  const index = selectedMenus.value.indexOf(menuId)
  if (index > -1) {
    selectedMenus.value.splice(index, 1)
    if (treeRef.value) {
      treeRef.value.setCheckedKeys(selectedMenus.value)
    }
  }
}

// 获取包含父菜单的完整菜单ID列表
const getCompleteMenuIds = (selectedIds: number[]): number[] => {
  const completeIds = new Set<number>()
  
  // 递归函数：获取菜单的所有父菜单ID
  const getParentIds = (menuId: number): number[] => {
    const parentIds: number[] = []
    const menu = menuMap.value.get(menuId)
    
    if (menu && menu.parentId && menu.parentId !== 0) {
      parentIds.push(menu.parentId)
      // 递归获取父菜单的父菜单
      parentIds.push(...getParentIds(menu.parentId))
    }
    
    return parentIds
  }
  
  // 遍历所有选中的菜单ID
  selectedIds.forEach(id => {
    // 添加当前菜单ID
    completeIds.add(id)
    
    // 添加所有父菜单ID
    const parentIds = getParentIds(id)
    parentIds.forEach(parentId => completeIds.add(parentId))
  })
  
  return Array.from(completeIds)
}

// 提交菜单配置
const handleSubmit = async () => {
  if (!props.roleData?.id) return

  try {
    submitting.value = true
    
    // 获取包含父菜单的完整菜单ID列表
    const completeMenuIds = getCompleteMenuIds(selectedMenus.value)
    
    console.log('原始选中的菜单ID:', selectedMenus.value)
    console.log('包含父菜单的完整ID列表:', completeMenuIds)
    
    await roleApi.assignMenusToRole(props.roleData.id, {
      menuIds: completeMenuIds
    })
    
    //ElMessage.success('菜单权限配置保存成功')
    console.log('菜单权限配置保存成功')
    emit('success')
    handleClose()
  } catch (error: any) {
    // 错误信息已在响应拦截器中处理
    console.error('菜单权限配置保存失败:', error)
  } finally {
    submitting.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  emit('update:visible', false)
}

// 监听对话框显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible && props.roleData) {
      fetchMenuTree().then(() => {
        fetchRoleMenus()
      })
    }
  }
)

// 初始化
onMounted(() => {
  if (props.visible && props.roleData) {
    fetchMenuTree().then(() => {
      fetchRoleMenus()
    })
  }
})
</script>

<style lang="scss" scoped>
.menu-assign {
  .operation-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: #f8f9fa;
    border-radius: 8px;
    margin-bottom: 16px;

    .operation-left {
      display: flex;
      align-items: center;
      gap: 12px;

      .selected-info {
        font-size: 13px;
        color: #606266;
      }
    }

    .operation-right {
      .el-button-group {
        .el-button {
          font-size: 12px;
        }
      }
    }
  }

  .menu-tree-container {
    min-height: 400px;
    max-height: 500px;
    overflow-y: auto;
    border: 1px solid #e4e7ed;
    border-radius: 8px;
    padding: 16px;

    .menu-tree {
      :deep(.el-tree-node) {
        .el-tree-node__content {
          height: auto;
          padding: 8px 0;

          &:hover {
            background: #f8f9ff;
          }
        }

        .el-tree-node__expand-icon {
          color: #c0c4cc;
          font-size: 16px;

          &.is-leaf {
            color: transparent;
          }
        }

        .el-checkbox {
          margin-right: 12px;
        }
      }

      .menu-node {
        flex: 1;
        padding-right: 12px;

        .node-content {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 4px;

          .node-icon {
            color: #606266;
            font-size: 16px;
          }

          .node-label {
            font-weight: 500;
            color: #303133;
            font-size: 14px;
          }

          .system-tag {
            margin-left: 8px;
          }

          .type-tag {
            margin-left: auto;
          }
        }

        .node-meta {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 12px;
          color: #909399;

          .node-code {
            font-family: 'Monaco', 'Consolas', monospace;
            padding: 2px 6px;
            background: #f0f0f0;
            border-radius: 4px;
          }

          .node-path {
            padding: 2px 6px;
            background: #e8f4fd;
            border-radius: 4px;
            color: #409eff;
          }
        }
      }
    }
  }

  .selected-preview {
    margin-top: 16px;
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;

    .preview-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      .preview-title {
        font-weight: 500;
        color: #303133;
        font-size: 14px;
      }
    }

    .preview-content {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .menu-tag {
        margin: 0;
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// Element Plus Dialog 样式覆盖
:deep(.el-dialog) {
  border-radius: 12px;
  overflow: hidden;

  .el-dialog__header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 20px 24px;
    border-bottom: none;

    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
    }

    .el-dialog__headerbtn {
      .el-dialog__close {
        color: white;
        font-size: 20px;

        &:hover {
          color: rgba(255, 255, 255, 0.8);
        }
      }
    }
  }

  .el-dialog__body {
    padding: 24px;
    max-height: 70vh;
    overflow-y: auto;
  }

  .el-dialog__footer {
    padding: 20px 24px;
    background: #fafafa;
    border-top: 1px solid #e4e7ed;
  }
}

// 响应式设计
@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 95% !important;
    margin: 20px auto;
  }

  .menu-assign {
    .operation-section {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;

      .operation-left {
        justify-content: center;
      }

      .operation-right {
        display: flex;
        justify-content: center;
      }
    }

    .selected-preview {
      .preview-header {
        flex-direction: column;
        align-items: stretch;
        gap: 8px;

        .preview-title {
          text-align: center;
        }
      }
    }
  }
}
</style>