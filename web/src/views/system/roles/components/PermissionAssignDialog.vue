<template>
  <el-dialog
    v-model="dialogVisible"
    :title="'为角色「' + roleData?.name + '」分配权限'"
    width="800px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="permission-assign">
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-form inline class="search-form">
          <el-form-item label="权限类型">
            <el-select
              v-model="searchForm.type"
              placeholder="选择权限类型"
              clearable
              style="width: 150px"
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="API接口" value="api" />
              <el-option label="页面权限" value="page" />
              <el-option label="按钮权限" value="button" />
              <el-option label="数据权限" value="data" />
            </el-select>
          </el-form-item>
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="搜索权限名称或编码"
              clearable
              style="width: 200px"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch" :loading="loading">
              搜索
            </el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 权限列表 -->
      <div class="permission-list">
        <div class="list-header">
          <div class="header-left">
            <el-checkbox
              v-model="selectAll"
              :indeterminate="indeterminate"
              @change="handleSelectAll"
            >
              全选
            </el-checkbox>
            <span class="selected-info">
              已选择 {{ selectedPermissions.length }} / {{ filteredPermissions.length }} 项
            </span>
          </div>
          <div class="header-right">
            <el-button-group>
              <el-button size="small" @click="expandAll">展开全部</el-button>
              <el-button size="small" @click="collapseAll">收起全部</el-button>
            </el-button-group>
          </div>
        </div>

        <div class="list-content" v-loading="loading">
          <div v-if="groupedPermissions.length === 0" class="empty-state">
            <el-empty description="暂无权限数据" />
          </div>
          <div v-else class="permission-groups">
            <div
              v-for="group in groupedPermissions"
              :key="group.type"
              class="permission-group"
            >
              <div class="group-header" @click="toggleGroup(group.type)">
                <el-icon class="expand-icon" :class="{ expanded: expandedGroups.has(group.type) }">
                  <ArrowRight />
                </el-icon>
                <span class="group-title">{{ getTypeLabel(group.type) }}</span>
                <el-tag type="info" size="small" class="group-count">
                  {{ group.permissions.length }} 项
                </el-tag>
              </div>
              <div
                v-show="expandedGroups.has(group.type)"
                class="group-content"
              >
                <div
                  v-for="permission in group.permissions"
                  :key="permission.id"
                  class="permission-item"
                  :class="{
                    'is-selected': selectedPermissions.includes(permission.id),
                    'is-system': permission.isSystem
                  }"
                >
                  <el-checkbox
                    :model-value="selectedPermissions.includes(permission.id)"
                    @change="(checked: boolean) => handlePermissionChange(permission.id, checked)"
                    :disabled="permission.isSystem && !selectedPermissions.includes(permission.id)"
                  >
                    <div class="permission-info">
                      <div class="permission-main">
                        <span class="permission-name">{{ permission.name }}</span>
                        <el-tag
                          v-if="permission.isSystem"
                          type="danger"
                          size="small"
                          effect="plain"
                        >
                          系统权限
                        </el-tag>
                        <el-tag
                          :type="getTypeTagType(permission.type)"
                          size="small"
                          effect="light"
                          class="type-tag"
                        >
                          {{ getTypeLabel(permission.type) }}
                        </el-tag>
                      </div>
                      <div class="permission-meta">
                        <span class="permission-code">{{ permission.code }}</span>
                        <span v-if="permission.resource" class="permission-resource">
                          {{ permission.resource }}
                        </span>
                        <span v-if="permission.action" class="permission-action">
                          {{ permission.action }}
                        </span>
                      </div>
                      <div v-if="permission.description" class="permission-desc">
                        {{ permission.description }}
                      </div>
                    </div>
                  </el-checkbox>
                </div>
              </div>
            </div>
          </div>
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
          保存权限配置
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowRight } from '@element-plus/icons-vue'
import { roleApi, permissionApi } from '@/api'
import type { Role, Permission, PermissionQueryParams } from '@/api/types'

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
const permissions = ref<Permission[]>([])
const selectedPermissions = ref<number[]>([])
const originalPermissions = ref<number[]>([])
const expandedGroups = ref(new Set<string>(['api', 'page', 'button', 'data']))

// 搜索表单
const searchForm = reactive({
  type: '',
  keyword: ''
})

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 过滤后的权限列表
const filteredPermissions = computed(() => {
  return permissions.value.filter(permission => {
    const matchType = !searchForm.type || permission.type === searchForm.type
    const matchKeyword = !searchForm.keyword || 
      permission.name.toLowerCase().includes(searchForm.keyword.toLowerCase()) ||
      permission.code.toLowerCase().includes(searchForm.keyword.toLowerCase())
    
    return matchType && matchKeyword
  })
})

// 按类型分组的权限
const groupedPermissions = computed(() => {
  const groups = new Map<string, Permission[]>()
  
  filteredPermissions.value.forEach(permission => {
    if (!groups.has(permission.type)) {
      groups.set(permission.type, [])
    }
    groups.get(permission.type)!.push(permission)
  })

  return Array.from(groups.entries()).map(([type, permissions]) => ({
    type,
    permissions: permissions.sort((a, b) => a.sortOrder - b.sortOrder)
  })).sort((a, b) => {
    const order = { api: 1, page: 2, button: 3, data: 4 }
    return (order[a.type as keyof typeof order] || 999) - (order[b.type as keyof typeof order] || 999)
  })
})

// 全选状态
const selectAll = computed({
  get: () => {
    const availablePermissions = filteredPermissions.value.filter(p => 
      !p.isSystem || selectedPermissions.value.includes(p.id)
    )
    return availablePermissions.length > 0 && 
           availablePermissions.every(p => selectedPermissions.value.includes(p.id))
  },
  set: (value: boolean) => {
    handleSelectAll(value)
  }
})

// 半选状态
const indeterminate = computed(() => {
  const availablePermissions = filteredPermissions.value.filter(p => 
    !p.isSystem || selectedPermissions.value.includes(p.id)
  )
  const selectedCount = availablePermissions.filter(p => 
    selectedPermissions.value.includes(p.id)
  ).length
  
  return selectedCount > 0 && selectedCount < availablePermissions.length
})

// 获取类型标签
const getTypeLabel = (type: string): string => {
  const typeMap: Record<string, string> = {
    api: 'API接口',
    page: '页面权限',
    button: '按钮权限',
    data: '数据权限'
  }
  return typeMap[type] || type
}

// 获取类型标签样式
const getTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    api: 'primary',
    page: 'success',
    button: 'warning',
    data: 'info'
  }
  return typeMap[type] || 'default'
}

// 获取权限列表
const fetchPermissions = async () => {
  try {
    loading.value = true
    const params: PermissionQueryParams = {
      page: 1,
      page_size: 1000, // 获取所有权限
      keyword: searchForm.keyword || undefined,
      type: searchForm.type || undefined,
      status: 1 // 只获取启用的权限
    }
    
    const { data } = await permissionApi.getPermissions(params)
    permissions.value = data.list
  } catch (error) {
    ElMessage.error('获取权限列表失败')
    permissions.value = []
  } finally {
    loading.value = false
  }
}

// 获取角色权限
const fetchRolePermissions = async () => {
  if (!props.roleData?.id) return

  try {
    const { data } = await roleApi.getRolePermissions(props.roleData.id)
    const permissionIds = data.map(p => p.id)
    selectedPermissions.value = [...permissionIds]
    originalPermissions.value = [...permissionIds]
  } catch (error) {
    console.error('获取角色权限失败:', error)
    selectedPermissions.value = []
    originalPermissions.value = []
  }
}

// 搜索
const handleSearch = () => {
  fetchPermissions()
}

// 重置搜索
const resetSearch = () => {
  searchForm.type = ''
  searchForm.keyword = ''
  fetchPermissions()
}

// 全选/取消全选
const handleSelectAll = (checked: boolean) => {
  const availablePermissions = filteredPermissions.value.filter(p => 
    !p.isSystem || selectedPermissions.value.includes(p.id)
  )
  
  if (checked) {
    const newIds = availablePermissions.map(p => p.id)
    selectedPermissions.value = [...new Set([...selectedPermissions.value, ...newIds])]
  } else {
    const idsToRemove = new Set(availablePermissions.map(p => p.id))
    selectedPermissions.value = selectedPermissions.value.filter(id => !idsToRemove.has(id))
  }
}

// 权限选择变化
const handlePermissionChange = (permissionId: number, checked: boolean) => {
  if (checked) {
    if (!selectedPermissions.value.includes(permissionId)) {
      selectedPermissions.value.push(permissionId)
    }
  } else {
    const index = selectedPermissions.value.indexOf(permissionId)
    if (index > -1) {
      selectedPermissions.value.splice(index, 1)
    }
  }
}

// 切换分组展开状态
const toggleGroup = (type: string) => {
  if (expandedGroups.value.has(type)) {
    expandedGroups.value.delete(type)
  } else {
    expandedGroups.value.add(type)
  }
}

// 展开全部
const expandAll = () => {
  expandedGroups.value = new Set(['api', 'page', 'button', 'data'])
}

// 收起全部
const collapseAll = () => {
  expandedGroups.value.clear()
}

// 提交权限配置
const handleSubmit = async () => {
  if (!props.roleData?.id) return

  try {
    submitting.value = true
    
    await roleApi.assignPermissionsToRole(props.roleData.id, {
      permission_ids: selectedPermissions.value
    })
    
    ElMessage.success('权限配置保存成功')
    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '权限配置保存失败')
  } finally {
    submitting.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  emit('update:visible', false)
  // 重置搜索条件
  searchForm.type = ''
  searchForm.keyword = ''
}

// 监听对话框显示状态
watch(
  () => props.visible,
  (visible) => {
    if (visible && props.roleData) {
      fetchPermissions()
      fetchRolePermissions()
    }
  }
)

// 初始化
onMounted(() => {
  if (props.visible && props.roleData) {
    fetchPermissions()
    fetchRolePermissions()
  }
})
</script>

<style lang="scss" scoped>
.permission-assign {
  .search-section {
    margin-bottom: 20px;
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;

    .search-form {
      display: flex;
      flex-wrap: wrap;
      align-items: flex-end;
      gap: 16px;

      .el-form-item {
        margin-bottom: 0;

        :deep(.el-form-item__label) {
          color: #606266;
          font-weight: 500;
          font-size: 13px;
        }
      }
    }
  }

  .permission-list {
    .list-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      background: #fff;
      border: 1px solid #e4e7ed;
      border-radius: 8px 8px 0 0;

      .header-left {
        display: flex;
        align-items: center;
        gap: 12px;

        .selected-info {
          font-size: 13px;
          color: #606266;
        }
      }

      .header-right {
        .el-button-group {
          .el-button {
            font-size: 12px;
          }
        }
      }
    }

    .list-content {
      min-height: 400px;
      max-height: 500px;
      overflow-y: auto;
      border: 1px solid #e4e7ed;
      border-top: none;
      border-radius: 0 0 8px 8px;

      .empty-state {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 400px;
      }

      .permission-groups {
        .permission-group {
          .group-header {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 12px 16px;
            background: #fafafa;
            border-bottom: 1px solid #e4e7ed;
            cursor: pointer;
            user-select: none;
            transition: all 0.2s ease;

            &:hover {
              background: #f0f0f0;
            }

            .expand-icon {
              transition: transform 0.2s ease;

              &.expanded {
                transform: rotate(90deg);
              }
            }

            .group-title {
              font-weight: 500;
              color: #303133;
            }

            .group-count {
              margin-left: auto;
            }
          }

          .group-content {
            .permission-item {
              padding: 12px 16px;
              border-bottom: 1px solid #f0f0f0;
              transition: all 0.2s ease;

              &:hover {
                background: #f8f9ff;
              }

              &.is-selected {
                background: #e3f2fd;
              }

              &.is-system {
                background: #fff3e0;
              }

              &:last-child {
                border-bottom: none;
              }

              :deep(.el-checkbox) {
                width: 100%;

                .el-checkbox__label {
                  width: 100%;
                  padding-left: 8px;
                }
              }

              .permission-info {
                width: 100%;

                .permission-main {
                  display: flex;
                  align-items: center;
                  gap: 8px;
                  margin-bottom: 4px;

                  .permission-name {
                    font-weight: 500;
                    color: #303133;
                    font-size: 14px;
                  }

                  .type-tag {
                    margin-left: auto;
                  }
                }

                .permission-meta {
                  display: flex;
                  align-items: center;
                  gap: 8px;
                  margin-bottom: 4px;
                  font-size: 12px;
                  color: #909399;

                  .permission-code {
                    font-family: 'Monaco', 'Consolas', monospace;
                    padding: 2px 6px;
                    background: #f0f0f0;
                    border-radius: 4px;
                  }

                  .permission-resource,
                  .permission-action {
                    padding: 2px 6px;
                    background: #e8f4fd;
                    border-radius: 4px;
                    color: #409eff;
                  }
                }

                .permission-desc {
                  font-size: 12px;
                  color: #666;
                  line-height: 1.4;
                }
              }
            }
          }
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

  .permission-assign {
    .search-section {
      .search-form {
        flex-direction: column;
        align-items: stretch;

        .el-form-item {
          width: 100%;
          margin-bottom: 12px;

          :deep(.el-input),
          :deep(.el-select) {
            width: 100% !important;
          }
        }
      }
    }

    .permission-list {
      .list-header {
        flex-direction: column;
        align-items: stretch;
        gap: 12px;

        .header-left {
          justify-content: center;
        }

        .header-right {
          display: flex;
          justify-content: center;
        }
      }
    }
  }
}
</style>