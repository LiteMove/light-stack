<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`为用户「${userData?.username}」分配角色`"
    width="600px"
    :before-close="handleClose"
    destroy-on-close
  >
    <div class="role-assign-content" v-loading="loading">
      <!-- 用户信息显示 -->
      <div class="user-info-card">
        <el-card shadow="never" class="user-card">
          <div class="user-display">
            <el-avatar 
              :src="userData?.avatar" 
              :size="50"
              :style="{ backgroundColor: getAvatarColor(userData?.username || '') }"
            >
              {{ userData?.nickname?.charAt(0) || userData?.username?.charAt(0) }}
            </el-avatar>
            <div class="user-details">
              <div class="user-name">{{ userData?.username }}</div>
              <div class="user-nickname">{{ userData?.nickname || '-' }}</div>
              <div class="user-status">
                <el-tag 
                  :type="userData?.status === 1 ? 'success' : 'danger'" 
                  size="small"
                >
                  {{ userData?.status === 1 ? '正常' : '禁用' }}
                </el-tag>
                <el-tag v-if="userData?.isSystem" type="danger" size="small">
                  系统用户
                </el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 角色选择 -->
      <div class="role-selection">
        <el-divider content-position="left">
          <el-icon><Key /></el-icon>
          角色分配
        </el-divider>

        <!-- 搜索框 -->
        <div class="role-search">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索角色名称或编码"
            :prefix-icon="Search"
            clearable
            @input="filterRoles"
          />
        </div>

        <!-- 角色列表 -->
        <div class="role-list">
          <el-checkbox-group v-model="selectedRoleIds" @change="handleRoleChange">
            <div class="role-grid">
              <div 
                v-for="role in filteredRoles" 
                :key="role.id"
                class="role-item"
                :class="{ 'role-item--selected': selectedRoleIds.includes(role.id) }"
              >
                <el-checkbox :label="role.id" class="role-checkbox">
                  <div class="role-content">
                    <div class="role-header">
                      <span class="role-name">{{ role.name }}</span>
                      <el-tag 
                        v-if="role.isSystem" 
                        type="danger" 
                        size="small" 
                        effect="plain"
                      >
                        系统角色
                      </el-tag>
                    </div>
                    <div class="role-meta">
                      <span class="role-code">{{ role.code }}</span>
                      <span class="role-sort">排序: {{ role.sortOrder }}</span>
                    </div>
                    <div v-if="role.description" class="role-description">
                      {{ role.description }}
                    </div>
                  </div>
                </el-checkbox>
              </div>
            </div>
          </el-checkbox-group>

          <!-- 空状态 -->
          <el-empty 
            v-if="filteredRoles.length === 0" 
            :image-size="120"
            description="没有找到相关角色"
          />
        </div>

        <!-- 选择统计 -->
        <div class="selection-summary">
          <el-alert
            :title="`已选择 ${selectedRoleIds.length} 个角色`"
            type="info"
            :closable="false"
            show-icon
          >
            <template v-if="selectedRoleIds.length > 0">
              <div class="selected-roles">
                <el-tag
                  v-for="roleId in selectedRoleIds"
                  :key="roleId"
                  :type="getRoleById(roleId)?.isSystem ? 'danger' : 'primary'"
                  size="small"
                  effect="light"
                  closable
                  @close="removeRole(roleId)"
                >
                  {{ getRoleById(roleId)?.name }}
                </el-tag>
              </div>
            </template>
          </el-alert>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="loading">取消</el-button>
        <el-button 
          type="primary" 
          :loading="loading" 
          @click="handleSubmit"
          :disabled="!hasChanged"
        >
          保存分配
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Key, Search } from '@element-plus/icons-vue'
import type { User, Role } from '@/api/types'
import {userApi} from "@/api";

interface Props {
  visible: boolean
  userData: User | null
  roles: Role[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const loading = ref(false)
const searchKeyword = ref('')
const selectedRoleIds = ref<number[]>([])
const originalRoleIds = ref<number[]>([])

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const filteredRoles = computed(() => {
  if (!searchKeyword.value) {
    return props.roles
  }
  const keyword = searchKeyword.value.toLowerCase()
  return props.roles.filter(role => 
    role.name.toLowerCase().includes(keyword) ||
    role.code.toLowerCase().includes(keyword) ||
    (role.description && role.description.toLowerCase().includes(keyword))
  )
})

const hasChanged = computed(() => {
  return JSON.stringify([...selectedRoleIds.value].sort()) !== 
         JSON.stringify([...originalRoleIds.value].sort())
})

// 获取用户头像颜色
const getAvatarColor = (username: string): string => {
  const colors = ['#f56a00', '#7265e6', '#ffbf00', '#00a2ae', '#f56565', '#38a169']
  const index = username.charCodeAt(0) % colors.length
  return colors[index]
}

// 根据ID获取角色
const getRoleById = (id: number): Role | undefined => {
  return props.roles.find(role => role.id === id)
}

// 移除角色
const removeRole = (roleId: number) => {
  const index = selectedRoleIds.value.indexOf(roleId)
  if (index > -1) {
    selectedRoleIds.value.splice(index, 1)
  }
}

// 过滤角色
const filterRoles = () => {
  // 搜索逻辑已在计算属性中处理
}

// 角色选择变化
const handleRoleChange = (value: number[]) => {
  selectedRoleIds.value = value
}

// 初始化用户角色
const initUserRoles = async () => {
  if (!props.userData) return
  
  try {
    loading.value = true
    // 这里应该从API获取用户当前的角色
    const { data } = await userApi.getUserRoles(props.userData.id)
    // 数组对象，提取角色ID
    // 为空判断
    if (!data || data.length === 0) {
      selectedRoleIds.value = []
      originalRoleIds.value = []
      return
    }
    const currentRoleIds = data.map((role: Role) => role.id)
    // 临时使用空数组
    selectedRoleIds.value = [...currentRoleIds]
    originalRoleIds.value = [...currentRoleIds]
  } catch (error) {
    console.error('获取用户角色失败:', error)
    selectedRoleIds.value = []
    originalRoleIds.value = []
  } finally {
    loading.value = false
  }
}

// 监听用户数据变化
watch(() => props.userData, (newUser) => {
  if (newUser && props.visible) {
    initUserRoles()
  }
}, { immediate: true })

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible && props.userData) {
    initUserRoles()
  } else if (!visible) {
    // 清空搜索
    searchKeyword.value = ''
  }
})

// 提交分配
const handleSubmit = async () => {
  if (!props.userData) return

  try {
    loading.value = true
    
    // 这里应该调用API保存角色分配
    await userApi.assignUserRoles(props.userData.id, {
      role_ids: selectedRoleIds.value
    })


    ElMessage.success('角色分配成功')
    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '角色分配失败')
  } finally {
    loading.value = false
  }
}

// 关闭弹窗
const handleClose = () => {
  selectedRoleIds.value = []
  originalRoleIds.value = []
  searchKeyword.value = ''
  emit('update:visible', false)
}
</script>

<style lang="scss" scoped>
.role-assign-content {
  max-height: 70vh;
  overflow-y: auto;
}

.user-info-card {
  margin-bottom: 20px;

  .user-card {
    :deep(.el-card__body) {
      padding: 16px;
    }

    .user-display {
      display: flex;
      align-items: center;
      gap: 16px;

      .user-details {
        flex: 1;

        .user-name {
          font-size: 16px;
          font-weight: 500;
          color: #303133;
          margin-bottom: 4px;
        }

        .user-nickname {
          font-size: 14px;
          color: #606266;
          margin-bottom: 8px;
        }

        .user-status {
          display: flex;
          gap: 8px;
          align-items: center;
        }
      }
    }
  }
}

.role-selection {
  .role-search {
    margin-bottom: 16px;
  }

  .role-list {
    max-height: 400px;
    overflow-y: auto;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    padding: 12px;
    margin-bottom: 16px;

    .role-grid {
      display: grid;
      gap: 12px;
    }

    .role-item {
      border: 1px solid #e4e7ed;
      border-radius: 6px;
      padding: 12px;
      transition: all 0.3s ease;
      background: #fff;

      &:hover {
        border-color: #c6e2ff;
        background-color: #ecf5ff;
      }

      &--selected {
        border-color: #409eff;
        background-color: #ecf5ff;
      }

      .role-checkbox {
        width: 100%;
        margin: 0;

        :deep(.el-checkbox__label) {
          width: 100%;
          padding-left: 8px;
        }

        :deep(.el-checkbox__input) {
          align-self: flex-start;
          margin-top: 2px;
        }
      }

      .role-content {
        .role-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;

          .role-name {
            font-weight: 500;
            color: #303133;
            font-size: 14px;
          }
        }

        .role-meta {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 4px;

          .role-code {
            font-size: 12px;
            color: #909399;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            background: #f5f7fa;
            padding: 2px 6px;
            border-radius: 3px;
          }

          .role-sort {
            font-size: 12px;
            color: #c0c4cc;
          }
        }

        .role-description {
          font-size: 12px;
          color: #606266;
          line-height: 1.4;
          margin-top: 4px;
        }
      }
    }
  }

  .selection-summary {
    .selected-roles {
      display: flex;
      flex-wrap: wrap;
      gap: 6px;
      margin-top: 8px;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 滚动条样式
.role-list::-webkit-scrollbar {
  width: 6px;
}

.role-list::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.role-list::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.role-list::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

// 分组标题优化
:deep(.el-divider) {
  margin: 20px 0 16px 0;

  .el-divider__text {
    color: #409eff;
    font-weight: 500;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 6px;
  }
}

// 空状态
:deep(.el-empty) {
  padding: 40px 0;
}
</style>