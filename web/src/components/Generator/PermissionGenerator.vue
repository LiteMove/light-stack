<template>
  <div class="permission-generator">
    <div class="permission-list">
      <div class="list-header">
        <span>权限列表</span>
        <div class="header-actions">
          <el-button
            size="small"
            type="primary"
            @click="generateDefaultPermissions"
          >
            自动生成
          </el-button>
          <el-button
            size="small"
            @click="addCustomPermission"
          >
            添加权限
          </el-button>
        </div>
      </div>

      <div class="permissions-container">
        <div
          v-for="(permission, index) in permissionList"
          :key="index"
          class="permission-item"
        >
          <el-input
            v-model="permission.code"
            placeholder="权限代码，格式：模块:业务:操作"
            @input="updatePermissions"
          >
            <template #prepend>
              <el-select
                v-model="permission.type"
                style="width: 80px"
                @change="updatePermissionCode(index)"
              >
                <el-option label="查看" value="list" />
                <el-option label="新增" value="add" />
                <el-option label="编辑" value="edit" />
                <el-option label="删除" value="delete" />
                <el-option label="详情" value="view" />
                <el-option label="导出" value="export" />
                <el-option label="导入" value="import" />
                <el-option label="自定义" value="custom" />
              </el-select>
            </template>
            <template #append>
              <el-button
                :icon="Delete"
                @click="removePermission(index)"
                :disabled="permissionList.length <= 1"
              />
            </template>
          </el-input>

          <el-input
            v-model="permission.description"
            placeholder="权限描述"
            style="margin-top: 8px"
          />
        </div>
      </div>
    </div>

    <div class="permission-preview">
      <div class="preview-header">
        <span>生成预览</span>
        <el-text type="info" size="small">
          格式：模块:业务:操作
        </el-text>
      </div>

      <div class="preview-content">
        <el-tag
          v-for="permission in modelValue"
          :key="permission"
          closable
          @close="removePermissionByCode(permission)"
          style="margin: 4px 8px 4px 0"
        >
          {{ permission }}
        </el-tag>

        <div v-if="!modelValue || modelValue.length === 0" class="empty-permissions">
          暂无权限配置
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import { Delete } from '@element-plus/icons-vue'

interface Props {
  modelValue: string[]
  moduleName?: string
  businessName?: string
}

interface Emits {
  (event: 'update:modelValue', value: string[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 权限项接口
interface PermissionItem {
  type: string
  code: string
  description: string
}

// 响应式数据
const permissionList = ref<PermissionItem[]>([
  { type: 'list', code: '', description: '查看列表' },
  { type: 'add', code: '', description: '新增' },
  { type: 'edit', code: '', description: '编辑' },
  { type: 'delete', code: '', description: '删除' },
  { type: 'view', code: '', description: '查看详情' }
])

// 监听属性变化，自动生成权限代码
watch(
  [() => props.moduleName, () => props.businessName],
  () => {
    if (props.moduleName && props.businessName) {
      generateDefaultPermissions()
    }
  },
  { immediate: true }
)

// 监听权限列表变化，同步到父组件
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && newValue.length > 0) {
      syncPermissionList(newValue)
    }
  },
  { immediate: true }
)

// 方法
const generateDefaultPermissions = () => {
  const moduleName = (props.moduleName || 'system').toLowerCase()
  const businessName = (props.businessName || 'example').toLowerCase()

  permissionList.value.forEach(item => {
    if (item.type !== 'custom') {
      item.code = `${moduleName}:${businessName}:${item.type}`
    }
  })

  updatePermissions()
}

const updatePermissionCode = (index: number) => {
  const item = permissionList.value[index]
  const moduleName = (props.moduleName || 'system').toLowerCase()
  const businessName = (props.businessName || 'example').toLowerCase()

  if (item.type !== 'custom') {
    item.code = `${moduleName}:${businessName}:${item.type}`
  }

  updatePermissions()
}

const addCustomPermission = () => {
  permissionList.value.push({
    type: 'custom',
    code: '',
    description: '自定义权限'
  })
}

const removePermission = (index: number) => {
  permissionList.value.splice(index, 1)
  updatePermissions()
}

const removePermissionByCode = (code: string) => {
  const index = permissionList.value.findIndex(item => item.code === code)
  if (index > -1) {
    removePermission(index)
  }
}

const updatePermissions = () => {
  const permissions = permissionList.value
    .map(item => item.code)
    .filter(code => code.trim() !== '')

  emit('update:modelValue', permissions)
}

const syncPermissionList = (permissions: string[]) => {
  // 根据现有权限同步权限列表
  if (!permissions || permissions.length === 0) {
    return
  }

  const newList: PermissionItem[] = []

  permissions.forEach(permission => {
    const parts = permission.split(':')
    if (parts.length === 3) {
      const operation = parts[2]
      const descriptions: { [key: string]: string } = {
        'list': '查看列表',
        'add': '新增',
        'edit': '编辑',
        'delete': '删除',
        'view': '查看详情',
        'export': '导出',
        'import': '导入'
      }

      newList.push({
        type: descriptions[operation] ? operation : 'custom',
        code: permission,
        description: descriptions[operation] || '自定义权限'
      })
    }
  })

  if (newList.length > 0) {
    permissionList.value = newList
  }
}
</script>

<style scoped lang="scss">
.permission-generator {
  display: flex;
  gap: 20px;
}

.permission-list {
  flex: 1;

  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;

    span {
      font-weight: 600;
      color: #303133;
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .permissions-container {
    .permission-item {
      margin-bottom: 12px;
    }
  }
}

.permission-preview {
  flex: 1;
  background: #fafafa;
  border-radius: 6px;
  padding: 16px;

  .preview-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #e4e7ed;

    span {
      font-weight: 600;
      color: #303133;
    }
  }

  .preview-content {
    min-height: 100px;

    .empty-permissions {
      color: #909399;
      text-align: center;
      padding: 40px 0;
      font-size: 14px;
    }
  }
}

:deep(.el-input-group__prepend) {
  padding: 0;
}

:deep(.el-input-group__append) {
  padding: 0 8px;
}
</style>