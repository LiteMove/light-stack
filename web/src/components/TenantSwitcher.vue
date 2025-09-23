<template>
  <div class="tenant-switcher" v-if="isSuperAdmin">
    <el-select
      v-model="selectedTenantId"
      placeholder="选择管理租户"
      style="width: 200px"
      @change="handleTenantChange"
      :loading="loading"
    >
      <template #prefix>
        <el-icon><OfficeBuilding /></el-icon>
      </template>
      
      <el-option
        v-for="tenant in tenantList"
        :key="tenant.id"
        :label="tenant.name"
        :value="tenant.id"
        :disabled="tenant.status !== 1"
      >
        <div class="option-content">
          <el-icon class="option-icon"><OfficeBuilding /></el-icon>
          <span>{{ tenant.name }}</span>
          <el-tag 
            v-if="tenant.status !== 1" 
            type="danger" 
            size="small" 
            effect="plain"
            class="option-tag"
          >
            已禁用
          </el-tag>
          <el-tag 
            v-else-if="isExpired(tenant.expiredAt)"
            type="warning" 
            size="small" 
            effect="plain"
            class="option-tag"
          >
            已过期
          </el-tag>
        </div>
      </el-option>
    </el-select>
    
    <!-- 当前租户信息显示 -->
    <div class="current-tenant-info" v-if="currentTenant">
      <el-tooltip :content="`当前管理租户: ${currentTenant.name}`" placement="bottom">
        <el-tag type="primary" effect="light" size="default">
          <el-icon><OfficeBuilding /></el-icon>
          {{ currentTenant.name }}
        </el-tag>
      </el-tooltip>
    </div>
    
    <!-- 未选择租户提示 -->
    <div class="no-tenant-warning" v-else>
      <el-tooltip content="请选择要管理的租户" placement="bottom">
        <el-tag type="warning" effect="light" size="default">
          <el-icon><Warning /></el-icon>
          请选择租户
        </el-tag>
      </el-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { OfficeBuilding, Warning } from '@element-plus/icons-vue'
import { useTenantStore } from '@/store/tenant'
import { useUserStore } from '@/store/user'
import type { Tenant } from '@/store/tenant'

const tenantStore = useTenantStore()
const userStore = useUserStore()
const loading = ref(false)

// 计算属性
const isSuperAdmin = computed(() => tenantStore.checkIsSuperAdmin())
const currentTenant = computed(() => tenantStore.getCurrentTenant())
const tenantList = computed(() => tenantStore.tenantList)

// 当前选中的租户ID
const selectedTenantId = ref<number | null>(
  currentTenant.value ? currentTenant.value.id : null
)

// 检查是否过期
const isExpired = (expiredAt?: string): boolean => {
  if (!expiredAt) return false
  return new Date(expiredAt) < new Date()
}

// 处理租户切换
const handleTenantChange = async (tenantId: number | null) => {
  if (!tenantId) {
    ElMessage.warning('请选择要管理的租户')
    selectedTenantId.value = currentTenant.value?.id || null
    return
  }
  
  try {
    loading.value = true
    
    const tenant = tenantList.value.find(t => t.id === tenantId)
    if (!tenant) {
      ElMessage.error('租户不存在')
      return
    }
    
    // 检查租户状态
    if (tenant.status !== 1) {
      ElMessage.error('该租户已被禁用')
      selectedTenantId.value = currentTenant.value?.id || null
      return
    }
    
    // 检查是否过期
    if (isExpired(tenant.expiredAt)) {
      const confirm = await ElMessageBox.confirm(
        '该租户已过期，确定要切换吗？',
        '租户已过期',
        {
          confirmButtonText: '确定切换',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).catch(() => false)
      
      if (!confirm) {
        selectedTenantId.value = currentTenant.value?.id || null
        return
      }
    }
    
    // 显示切换确认
    const confirm = await ElMessageBox.confirm(
      `确定要切换到租户"${tenant.name}"吗？`,
      '切换租户',
      {
        confirmButtonText: '确定切换',
        cancelButtonText: '取消',
        type: 'info'
      }
    ).catch(() => false)
    
    if (confirm) {
      // 执行切换
      tenantStore.switchTenant(tenant)
      //ElMessage.success(`已切换到租户"${tenant.name}"`)
    } else {
      // 恢复选择
      selectedTenantId.value = currentTenant.value?.id || null
    }
  } catch (error: any) {
    ElMessage.error(error.message || '切换失败')
    selectedTenantId.value = currentTenant.value?.id || null
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(async () => {
  console.log('isSuperAdmin:', isSuperAdmin.value)
  if (isSuperAdmin.value) {
    // 如果是超级管理员，加载租户列表
    loading.value = true
    try {
      await tenantStore.fetchTenantList()
    } catch (error) {
      console.error('Failed to load tenant list:', error)
    } finally {
      loading.value = false
    }
  }
})
</script>

<style lang="scss" scoped>
.tenant-switcher {
  display: flex;
  align-items: center;
  gap: 12px;

  .current-tenant-info {
    .el-tag {
      display: flex;
      align-items: center;
      gap: 4px;
      font-weight: 500;
      
      .el-icon {
        font-size: 14px;
      }
    }
  }
}

// 下拉选项样式
:deep(.el-select) {
  .el-select__wrapper {
    .el-select__prefix {
      color: #409eff;
    }
  }
}

:deep(.el-select-dropdown) {
  .global-option {
    .el-select-dropdown__item {
      font-weight: 600;
      color: #e6a23c;
    }
  }
  
  .option-content {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    
    .option-icon {
      color: #409eff;
      font-size: 14px;
      flex-shrink: 0;
    }
    
    span {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    
    .option-tag {
      margin-left: auto;
      flex-shrink: 0;
    }
  }
  
  .el-select-dropdown__item.is-disabled {
    .option-content {
      .option-icon {
        color: #c0c4cc;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .tenant-switcher {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
    
    .el-select {
      width: 100% !important;
    }
    
    .current-tenant-info {
      display: flex;
      justify-content: center;
    }
  }
}
</style>