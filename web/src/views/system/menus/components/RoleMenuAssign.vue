<template>
  <el-dialog
    v-model="dialogVisible"
    title="分配菜单权限"
    width="800px"
    :before-close="handleClose"
    destroy-on-close
  >
    <div class="role-menu-assign">
      <div class="role-info">
        <el-descriptions :column="2" size="small" border>
          <el-descriptions-item label="角色名称">
            {{ roleInfo.name }}
          </el-descriptions-item>
          <el-descriptions-item label="角色编码">
            {{ roleInfo.code }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ roleInfo.description || '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <el-divider content-position="left">菜单权限配置</el-divider>

      <div class="menu-tree-container" v-loading="loading">
        <MenuTreeSelect
          ref="menuTreeRef"
          :data="menuTree"
          :checked-keys="checkedMenuIds"
          :show-actions="true"
          :show-selected-info="true"
          @change="handleMenuChange"
        />
      </div>

      <div class="permission-summary" v-if="selectedMenus.length > 0">
        <el-alert
          :title="`已选择 ${selectedMenus.length} 个菜单项`"
          type="info"
          :closable="false"
        >
          <template #default>
            <div class="summary-content">
              <div class="summary-item">
                <span class="summary-label">目录:</span>
                <span class="summary-value">{{ getMenuCountByType('directory') }} 个</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">菜单:</span>
                <span class="summary-value">{{ getMenuCountByType('menu') }} 个</span>
              </div>
              <div class="summary-item">
                <span class="summary-label">权限:</span>
                <span class="summary-value">{{ getMenuCountByType('permission') }} 个</span>
              </div>
            </div>
          </template>
        </el-alert>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">
          保存分配
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { menuApi } from '@/api'
import type { Menu, Role } from '@/api/types'
import MenuTreeSelect from './MenuTreeSelect.vue'

interface Props {
  visible: boolean
  roleInfo: Partial<Role>
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const menuTreeRef = ref<InstanceType<typeof MenuTreeSelect>>()
const loading = ref(false)
const saving = ref(false)
const menuTree = ref<Menu[]>([])
const checkedMenuIds = ref<number[]>([])
const selectedMenuIds = ref<number[]>([])

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const selectedMenus = computed(() => {
  return getMenusByIds(menuTree.value, selectedMenuIds.value)
})

// 获取菜单数据
const fetchMenuTree = async () => {
  try {
    loading.value = true
    const { data } = await menuApi.getMenuTree()
    menuTree.value = data
  } catch (error) {
    ElMessage.error('获取菜单数据失败')
  } finally {
    loading.value = false
  }
}

// 获取角色已有菜单
const fetchRoleMenus = async () => {
  if (!props.roleInfo.id) return
  
  try {
    const { data } = await menuApi.getRoleMenus(props.roleInfo.id)
    checkedMenuIds.value = data.map(menu => menu.id)
    selectedMenuIds.value = [...checkedMenuIds.value]
  } catch (error) {
    ElMessage.error('获取角色菜单失败')
  }
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

// 获取指定类型的菜单数量
const getMenuCountByType = (type: string): number => {
  return selectedMenus.value.filter(menu => menu.type === type).length
}

// 菜单选择变化
const handleMenuChange = (checkedKeys: number[]) => {
  selectedMenuIds.value = checkedKeys
}

// 保存分配
const handleSave = async () => {
  if (!props.roleInfo.id) {
    ElMessage.error('角色信息不完整')
    return
  }

  try {
    saving.value = true
    await menuApi.assignMenusToRole(props.roleInfo.id, {
      menuIds: selectedMenuIds.value
    })
    
    ElMessage.success('菜单权限分配成功')
    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error('菜单权限分配失败')
  } finally {
    saving.value = false
  }
}

// 关闭弹窗
const handleClose = () => {
  emit('update:visible', false)
}

// 监听弹窗显示状态
watch(() => props.visible, async (visible) => {
  if (visible) {
    await fetchMenuTree()
    await fetchRoleMenus()
  }
})
</script>

<style lang="scss" scoped>
.role-menu-assign {
  .role-info {
    margin-bottom: 20px;
  }

  .menu-tree-container {
    min-height: 300px;
    max-height: 500px;
  }

  .permission-summary {
    margin-top: 16px;

    .summary-content {
      display: flex;
      gap: 20px;
      margin-top: 8px;

      .summary-item {
        display: flex;
        align-items: center;
        gap: 4px;

        .summary-label {
          color: #606266;
          font-size: 13px;
        }

        .summary-value {
          font-weight: 600;
          color: #409eff;
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 响应式设计
@media (max-width: 768px) {
  .role-menu-assign {
    .permission-summary {
      .summary-content {
        flex-direction: column;
        gap: 8px;
      }
    }
  }
}
</style>