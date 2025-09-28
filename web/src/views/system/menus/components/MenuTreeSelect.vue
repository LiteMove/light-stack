<template>
  <div class="menu-tree-select">
    <el-tree
      ref="treeRef"
      :data="treeData"
      :props="treeProps"
      show-checkbox
      node-key="id"
      check-strictly
      :default-checked-keys="defaultCheckedKeys"
      :default-expanded-keys="defaultExpandedKeys"
      @check="handleCheck"
      class="menu-tree"
    >
      <template #default="{ node, data }">
        <div class="tree-node">
          <div class="node-content">
            <el-icon v-if="data.icon" class="node-icon">
              <component :is="getIconComponent(data.icon)" />
            </el-icon>
            <span class="node-label">{{ data.name }}</span>
            <el-tag 
              :type="getMenuTypeTagType(data.type)" 
              size="small" 
              class="node-tag"
            >
              {{ getMenuTypeLabel(data.type) }}
            </el-tag>
            <el-tag 
              v-if="data.isSystem"
              type="info" 
              size="small" 
              class="node-tag"
            >
              系统
            </el-tag>
          </div>
          <div class="node-info">
            <span class="node-code">{{ data.code }}</span>
            <span v-if="data.path" class="node-path">{{ data.path }}</span>
          </div>
        </div>
      </template>
    </el-tree>

    <div class="tree-actions" v-if="showActions">
      <el-button @click="expandAll">展开全部</el-button>
      <el-button @click="collapseAll">折叠全部</el-button>
      <el-button @click="checkAll">全选</el-button>
      <el-button @click="uncheckAll">取消全选</el-button>
      <el-button @click="checkParentOnly">仅选父级</el-button>
    </div>

    <div class="selected-info" v-if="showSelectedInfo">
      <el-divider content-position="left">
        <span class="selected-title">已选择的菜单 ({{ selectedMenus.length }})</span>
      </el-divider>
      <div class="selected-items">
        <el-tag
          v-for="menu in selectedMenus"
          :key="menu.id"
          :type="getMenuTypeTagType(menu.type)"
          size="small"
          closable
          @close="removeSelection(menu.id)"
          class="selected-tag"
        >
          {{ menu.name }}
        </el-tag>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import type { ElTree } from 'element-plus'
import type { Menu } from '@/api/types'

interface Props {
  data: Menu[]
  checkedKeys?: number[]
  showActions?: boolean
  showSelectedInfo?: boolean
  checkStrictly?: boolean
  defaultExpandAll?: boolean
}

interface Emits {
  (e: 'check', checkedKeys: number[], checkedNodes: Menu[]): void
  (e: 'change', checkedKeys: number[]): void
}

const props = withDefaults(defineProps<Props>(), {
  checkedKeys: () => [],
  showActions: true,
  showSelectedInfo: true,
  checkStrictly: false,
  defaultExpandAll: true
})

const emit = defineEmits<Emits>()

const treeRef = ref<InstanceType<typeof ElTree>>()

// 树配置
const treeProps = {
  label: 'name',
  children: 'children'
}

// 计算属性
const treeData = computed(() => props.data)

const defaultCheckedKeys = computed(() => props.checkedKeys)

const defaultExpandedKeys = computed(() => {
  if (!props.defaultExpandAll) return []
  return getAllNodeIds(props.data)
})

const selectedMenus = computed(() => {
  const checkedKeys = treeRef.value?.getCheckedKeys() as number[] || []
  return getMenusByIds(props.data, checkedKeys)
})

// 获取所有节点ID
const getAllNodeIds = (nodes: Menu[]): number[] => {
  let ids: number[] = []
  nodes.forEach(node => {
    ids.push(node.id)
    if (node.children?.length) {
      ids = ids.concat(getAllNodeIds(node.children))
    }
  })
  return ids
}

// 根据ID获取菜单对象
const getMenusByIds = (nodes: Menu[], ids: number[]): Menu[] => {
  let menus: Menu[] = []
  nodes.forEach(node => {
    if (ids.includes(node.id)) {
      menus.push(node)
    }
    if (node.children?.length) {
      menus = menus.concat(getMenusByIds(node.children, ids))
    }
  })
  return menus
}

// 展开全部
const expandAll = () => {
  const allNodeIds = getAllNodeIds(props.data)
  allNodeIds.forEach(id => {
    treeRef.value?.setExpanded(id, true)
  })
}

// 折叠全部
const collapseAll = () => {
  const allNodeIds = getAllNodeIds(props.data)
  allNodeIds.forEach(id => {
    treeRef.value?.setExpanded(id, false)
  })
}

// 全选
const checkAll = () => {
  const allNodeIds = getAllNodeIds(props.data)
  treeRef.value?.setCheckedKeys(allNodeIds)
  handleCheck()
}

// 取消全选
const uncheckAll = () => {
  treeRef.value?.setCheckedKeys([])
  handleCheck()
}

// 仅选择父级菜单
const checkParentOnly = () => {
  const parentIds = props.data.map(node => node.id)
  treeRef.value?.setCheckedKeys(parentIds)
  handleCheck()
}

// 移除选择
const removeSelection = (id: number) => {
  treeRef.value?.setChecked(id, false)
  handleCheck()
}

// 处理选择变化
const handleCheck = () => {
  const checkedKeys = treeRef.value?.getCheckedKeys() as number[] || []
  const checkedNodes = treeRef.value?.getCheckedNodes() as Menu[] || []
  
  emit('check', checkedKeys, checkedNodes)
  emit('change', checkedKeys)
}

// 获取菜单类型标签类型
const getMenuTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    directory: 'info',
    menu: 'success',
  }
  return typeMap[type] || 'info'
}

// 获取菜单类型标签
const getMenuTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    directory: '目录',
    menu: '菜单',
  }
  return typeMap[type] || type
}

// 获取图标组件
const getIconComponent = (icon: string) => {
  return 'Menu'
}

// 监听外部选中keys变化
watch(() => props.checkedKeys, (newKeys) => {
  if (treeRef.value && newKeys) {
    nextTick(() => {
      treeRef.value!.setCheckedKeys(newKeys)
    })
  }
})

// 暴露方法给父组件
defineExpose({
  getCheckedKeys: () => treeRef.value?.getCheckedKeys(),
  getCheckedNodes: () => treeRef.value?.getCheckedNodes(),
  setCheckedKeys: (keys: number[]) => treeRef.value?.setCheckedKeys(keys),
  expandAll,
  collapseAll,
  checkAll,
  uncheckAll
})
</script>

<style lang="scss" scoped>
.menu-tree-select {
  .menu-tree {
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    max-height: 400px;
    overflow-y: auto;

    :deep(.el-tree-node__content) {
      height: auto;
      min-height: 32px;
      padding: 8px 20px 8px 8px;
      
      &:hover {
        background-color: #f5f7fa;
      }
    }

    .tree-node {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 4px;

      .node-content {
        display: flex;
        align-items: center;
        gap: 8px;

        .node-icon {
          color: #409eff;
          font-size: 16px;
        }

        .node-label {
          font-weight: 500;
          color: #303133;
        }

        .node-tag {
          margin-left: auto;
        }
      }

      .node-info {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 12px;
        color: #909399;

        .node-code {
          font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
          background-color: #f5f7fa;
          padding: 2px 6px;
          border-radius: 4px;
          border: 1px solid #e4e7ed;
        }

        .node-path {
          font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
          color: #67c23a;
        }
      }
    }
  }

  .tree-actions {
    margin-top: 12px;
    padding: 12px;
    background-color: #f5f7fa;
    border-radius: 4px;
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .selected-info {
    margin-top: 16px;

    .selected-title {
      font-weight: 500;
      color: #303133;
    }

    .selected-items {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      max-height: 120px;
      overflow-y: auto;

      .selected-tag {
        margin: 0;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .menu-tree-select {
    .tree-actions {
      .el-button {
        flex: 1;
        min-width: calc(50% - 4px);
      }
    }

    .selected-items {
      .selected-tag {
        flex-shrink: 0;
      }
    }
  }
}
</style>