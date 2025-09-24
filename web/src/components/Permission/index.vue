<template>
  <slot v-if="hasAuth"></slot>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePermission } from '@/utils/permission'

interface Props {
  // 权限码，支持单个或多个
  permission?: string | string[]
  // 角色，支持单个或多个
  role?: string | string[]
  // 是否需要所有权限/角色（默认false，只需满足其一）
  requireAll?: boolean
  // 综合权限配置
  auth?: {
    permissions?: string[]
    roles?: string[]
    requireAll?: boolean
  }
}

const props = withDefaults(defineProps<Props>(), {
  permission: undefined,
  role: undefined,
  requireAll: false,
  auth: undefined
})

const { 
  hasPermission, 
  hasAnyPermission, 
  hasAllPermissions,
  hasRole,
  hasAnyRole,
  hasAllRoles,
  checkAuth 
} = usePermission()

const hasAuth = computed(() => {
  // 如果提供了综合配置，使用综合检查
  if (props.auth) {
    return checkAuth(props.auth)
  }

  let permissionCheck = true
  let roleCheck = true

  // 检查权限
  if (props.permission) {
    if (Array.isArray(props.permission)) {
      permissionCheck = props.requireAll 
        ? hasAllPermissions(props.permission)
        : hasAnyPermission(props.permission)
    } else {
      permissionCheck = hasPermission(props.permission)
    }
  }

  // 检查角色
  if (props.role) {
    if (Array.isArray(props.role)) {
      roleCheck = props.requireAll
        ? hasAllRoles(props.role)
        : hasAnyRole(props.role)
    } else {
      roleCheck = hasRole(props.role)
    }
  }

  // 权限和角色只要有一个满足即可
  return permissionCheck || roleCheck
})
</script>