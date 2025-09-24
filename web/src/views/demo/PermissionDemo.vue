<template>
  <div class="permission-demo">
    <el-card class="demo-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>权限控制演示</span>
          <el-tag type="info" size="small">Demo</el-tag>
        </div>
      </template>

      <!-- 当前用户权限信息 -->
      <div class="user-info-section">
        <h4>当前用户权限信息</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户角色">
            <el-tag 
              v-for="role in userInfo?.roles" 
              :key="role" 
              type="primary" 
              size="small"
              style="margin-right: 8px;"
            >
              {{ role }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="权限数量">
            {{ permissions?.length || 0 }} 个
          </el-descriptions-item>
          <el-descriptions-item label="是否管理员">
            <el-tag :type="isAdmin() ? 'success' : 'info'" size="small">
              {{ isAdmin() ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="是否超级管理员">
            <el-tag :type="isSuperAdmin() ? 'danger' : 'info'" size="small">
              {{ isSuperAdmin() ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 指令方式演示 -->
      <div class="demo-section">
        <h4>1. 指令方式演示</h4>
        <div class="demo-buttons">
          <!-- 权限指令 -->
          <el-button v-permission="'system:user:create'" type="primary">
            创建用户 (权限: system:user:create)
          </el-button>
          <el-button v-permission="'system:user:update'" type="success">
            编辑用户 (权限: system:user:update)
          </el-button>
          <el-button v-permission="'system:user:delete'" type="danger">
            删除用户 (权限: system:user:delete)
          </el-button>
          
          <!-- 角色指令 -->
          <el-button v-role="'tenant_admin'" type="warning">
            租户管理员功能 (角色: tenant_admin)
          </el-button>
          <el-button v-role="'super_admin'" type="info">
            超级管理员功能 (角色: super_admin)
          </el-button>
          
          <!-- 管理员指令 -->
          <el-button v-admin type="danger">
            租户管理员专用 (v-admin)
          </el-button>
          <el-button v-super-admin type="danger">
            超管专用 (v-super-admin)
          </el-button>
        </div>
      </div>

      <!-- 组件方式演示 -->
      <div class="demo-section">
        <h4>2. 组件方式演示</h4>
        <div class="demo-buttons">
          <Permission permission="system:user:create">
            <el-button type="primary">
              组件包装 - 创建用户
            </el-button>
          </Permission>
          
          <Permission :permission="['system:user:update', 'system:user:delete']">
            <el-button type="warning">
              组件包装 - 编辑或删除用户 (任意一个权限)
            </el-button>
          </Permission>
          
          <Permission role="tenant_admin">
            <el-button type="success">
              组件包装 - 租户管理员功能
            </el-button>
          </Permission>
          
          <Permission :auth="{ permissions: ['system:user:create'], roles: ['tenant_admin'] }">
            <el-button type="info">
              组件包装 - 综合权限检查
            </el-button>
          </Permission>
        </div>
      </div>

      <!-- Composition API 演示 -->
      <div class="demo-section">
        <h4>3. Composition API 演示</h4>
        <div class="demo-buttons">
          <el-button 
            v-if="hasPermission('system:user:create')" 
            type="primary"
          >
            API检查 - 创建用户
          </el-button>
          
          <el-button 
            v-if="hasAnyPermission(['system:user:update', 'system:user:delete'])" 
            type="success"
          >
            API检查 - 编辑或删除用户
          </el-button>
          
          <el-button 
            v-if="hasRole('tenant_admin')" 
            type="warning"
          >
            API检查 - 租户管理员功能
          </el-button>
          
          <el-button 
            v-if="checkAuth({ permissions: ['system:user:view'], roles: ['user'] })" 
            type="info"
          >
            API检查 - 综合权限
          </el-button>
        </div>
      </div>

      <!-- 动态权限检查演示 -->
      <div class="demo-section">
        <h4>4. 动态权限检查演示</h4>
        <div class="permission-checks">
          <el-form :model="checkForm" inline>
            <el-form-item label="权限码">
              <el-input 
                v-model="checkForm.permission" 
                placeholder="输入权限码，如: system:user:create"
                style="width: 250px;"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="checkPermissionDynamic">
                检查权限
              </el-button>
            </el-form-item>
          </el-form>
          <div v-if="checkResult !== null" class="check-result">
            <el-alert 
              :type="checkResult ? 'success' : 'error'"
              :title="checkResult ? '拥有权限' : '无此权限'"
              show-icon
              :closable="false"
            />
          </div>
        </div>
      </div>

      <!-- 权限列表展示 -->
      <div class="demo-section">
        <h4>5. 当前用户权限列表</h4>
        <div class="permissions-list">
          <el-tag 
            v-for="permission in permissions" 
            :key="permission"
            type="success"
            size="small"
            style="margin: 2px;"
          >
            {{ permission }}
          </el-tag>
          <div v-if="!permissions?.length" class="no-permissions">
            <el-empty description="暂无权限" :image-size="80" />
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store'
import { usePermission } from '@/utils/permission'
import Permission from '@/components/Permission/index.vue'

// 用户store
const userStore = useUserStore()
const { userInfo, permissions } = storeToRefs(userStore)

// 权限检查Hook
const {
  hasPermission,
  hasAnyPermission,
  hasAllPermissions,
  hasRole,
  isAdmin,
  isSuperAdmin,
  checkAuth
} = usePermission()

// 动态权限检查
const checkForm = reactive({
  permission: 'system:user:create'
})
const checkResult = ref<boolean | null>(null)

const checkPermissionDynamic = () => {
  if (checkForm.permission.trim()) {
    checkResult.value = hasPermission(checkForm.permission.trim())
  }
}
</script>

<style lang="scss" scoped>
.permission-demo {
  padding: 24px;
  
  .demo-card {
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: 600;
      font-size: 16px;
    }
  }
  
  .user-info-section {
    margin-bottom: 24px;
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;
  }
  
  .demo-section {
    margin-bottom: 24px;
    padding: 16px 0;
    border-bottom: 1px solid #e4e7ed;
    
    &:last-child {
      border-bottom: none;
    }
    
    h4 {
      margin-bottom: 16px;
      color: #303133;
      font-size: 14px;
      font-weight: 600;
    }
  }
  
  .demo-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    
    .el-button {
      margin: 0;
    }
  }
  
  .permission-checks {
    .check-result {
      margin-top: 16px;
    }
  }
  
  .permissions-list {
    max-height: 200px;
    overflow-y: auto;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    padding: 12px;
    background: #fafafa;
    
    .no-permissions {
      text-align: center;
      padding: 20px;
    }
  }
}

// 响应式优化
@media (max-width: 768px) {
  .permission-demo {
    padding: 12px;
    
    .demo-buttons {
      flex-direction: column;
      
      .el-button {
        width: 100%;
      }
    }
    
    .el-form--inline {
      .el-form-item {
        width: 100%;
        margin-bottom: 12px;
        
        .el-input {
          width: 100% !important;
        }
      }
    }
  }
}
</style>