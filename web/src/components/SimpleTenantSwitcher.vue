<template>
  <div class="simple-tenant-switcher" v-if="isSuperAdmin">
    <!-- 当前租户显示按钮 -->
    <el-dropdown trigger="click" @command="handleTenantSwitch">
      <el-button
        v-if="currentTenant"
        type="primary"
        :icon="OfficeBuilding"
        class="tenant-btn"
      >
        {{ currentTenant.name }}
        <el-icon class="dropdown-icon">
          <ArrowDown />
        </el-icon>
      </el-button>

      <!-- 未选择租户按钮 -->
      <el-button
        v-else
        type="warning"
        :icon="Warning"
        class="tenant-btn"
      >
        请选择租户
        <el-icon class="dropdown-icon">
          <ArrowDown />
        </el-icon>
      </el-button>

      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item
            v-for="tenant in tenantList"
            :key="tenant.id"
            :command="tenant.id"
            :disabled="tenant.status !== 1"
          >
            <div class="tenant-option">
              <el-icon class="tenant-icon"><OfficeBuilding /></el-icon>
              <span class="tenant-name">{{ tenant.name }}</span>
              <el-icon v-if="currentTenant?.id === tenant.id" class="current-icon">
                <Check />
              </el-icon>
              <el-tag v-else-if="tenant.status !== 1" type="danger" size="small">
                禁用
              </el-tag>
              <el-tag v-else-if="isExpired(tenant.expiredAt)" type="warning" size="small">
                过期
              </el-tag>
            </div>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  OfficeBuilding,
  Warning,
  ArrowDown,
  Check
} from '@element-plus/icons-vue'
import { useTenantStore } from '@/store/tenant'
import type { Tenant } from '@/store/tenant'

const tenantStore = useTenantStore()

// 计算属性
const isSuperAdmin = computed(() => tenantStore.checkIsSuperAdmin())
const currentTenant = computed(() => tenantStore.getCurrentTenant())
const tenantList = computed(() => tenantStore.tenantList)

// 检查是否过期
const isExpired = (expiredAt?: string): boolean => {
  if (!expiredAt) return false
  return new Date(expiredAt) < new Date()
}

// 处理租户切换
const handleTenantSwitch = (tenantId: number) => {
  const tenant = tenantList.value.find(t => t.id === tenantId)
  if (!tenant) return

  if (tenant.status !== 1) {
    ElMessage.error('该租户已被禁用')
    return
  }

  if (currentTenant.value?.id === tenant.id) {
    return
  }

  try {
    // 直接切换
    tenantStore.switchTenant(tenant)

    // 显示结果消息
    if (isExpired(tenant.expiredAt)) {
      ElMessage.warning(`已切换到"${tenant.name}"（该租户已过期）`)
    } else {
      ElMessage.success(`已切换到"${tenant.name}"`)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '切换失败')
  }
}

// 初始化
onMounted(async () => {
  if (isSuperAdmin.value) {
    try {
      await tenantStore.fetchTenantList()
    } catch (error) {
      console.error('Failed to load tenant list:', error)
    }
  }
})
</script>

<style lang="scss" scoped>
.simple-tenant-switcher {
  .tenant-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 120px;
    transition: all 0.3s ease;

    .dropdown-icon {
      margin-left: auto;
    }

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
    }
  }
}

:deep(.el-dropdown-menu) {
  .tenant-option {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    min-width: 200px;

    .tenant-icon {
      color: #409eff;
      font-size: 14px;
      flex-shrink: 0;
    }

    .tenant-name {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .current-icon {
      color: #67c23a;
      font-size: 14px;
    }
  }

  .el-dropdown-menu__item.is-disabled {
    .tenant-option {
      .tenant-icon {
        color: #c0c4cc;
      }
      .tenant-name {
        color: #909399;
      }
    }
  }
}
</style>