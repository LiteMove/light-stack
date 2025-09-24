import type { App, DirectiveBinding } from 'vue'
import { useUserStore } from '@/store'

/**
 * 权限指令和工具函数
 * 
 * 使用方式：
 * 1. 指令方式：
 *    v-permission="'user:create'"
 *    v-permission="['user:create', 'user:update']"
 *    v-permission.all="['user:create', 'user:update']"
 *    v-role="'admin'"
 *    v-auth="{ permissions: ['user:create'], roles: ['admin'] }"
 * 
 * 2. v-if 方式：
 *    v-if="hasPer('user:create')"
 *    v-if="hasRole('admin')"
 *    v-if="hasAnyPer(['user:create', 'user:update'])"
 *    v-if="hasAllPer(['user:create', 'user:update'])"
 */

// 全局权限检查函数 - 用于 v-if
export function hasPer(permission: string): boolean {
  const userStore = useUserStore()
  return userStore.hasPermission(permission)
}

export function hasRole(role: string): boolean {
  const userStore = useUserStore()
  return userStore.hasRole(role)
}

export function hasAnyPer(permissions: string[]): boolean {
  const userStore = useUserStore()
  return userStore.hasAnyPermission(permissions)
}

export function hasAllPer(permissions: string[]): boolean {
  const userStore = useUserStore()
  return userStore.hasAllPermissions(permissions)
}

export function hasAnyRole(roles: string[]): boolean {
  const userStore = useUserStore()
  return userStore.hasAnyRole(roles)
}

export function hasAllRole(roles: string[]): boolean {
  const userStore = useUserStore()
  return userStore.hasAllRoles(roles)
}

export function isAdmin(): boolean {
  const userStore = useUserStore()
  return userStore.hasRole('tenant_admin')
}

export function isSuperAdmin(): boolean {
  const userStore = useUserStore()
  return userStore.hasRole('super_admin')
}

// 综合权限检查函数
export function hasAuth(config: { 
  permissions?: string[], 
  roles?: string[], 
  requireAll?: boolean 
}): boolean {
  return checkAuth(config)
}

// 权限检查工具函数
export function checkPermission(permission: string | string[]): boolean {
  const userStore = useUserStore()
  
  if (Array.isArray(permission)) {
    return permission.some(p => userStore.hasPermission(p))
  }
  return userStore.hasPermission(permission)
}

export function checkAllPermissions(permissions: string[]): boolean {
  const userStore = useUserStore()
  return userStore.hasAllPermissions(permissions)
}

export function checkRole(role: string | string[]): boolean {
  const userStore = useUserStore()
  
  if (Array.isArray(role)) {
    return role.some(r => userStore.hasRole(r))
  }
  return userStore.hasRole(role)
}

export function checkAllRoles(roles: string[]): boolean {
  const userStore = useUserStore()
  return userStore.hasAllRoles(roles)
}

export function checkAuth(config: { 
  permissions?: string[], 
  roles?: string[], 
  requireAll?: boolean 
}): boolean {
  const userStore = useUserStore()
  const { permissions, roles, requireAll = false } = config

  let hasPermission = true
  let hasRole = true

  if (permissions && permissions.length > 0) {
    hasPermission = requireAll 
      ? userStore.hasAllPermissions(permissions)
      : userStore.hasAnyPermission(permissions)
  }

  if (roles && roles.length > 0) {
    hasRole = requireAll
      ? userStore.hasAllRoles(roles) 
      : userStore.hasAnyRole(roles)
  }

  // 权限和角色只要有一个满足即可
  return hasPermission || hasRole
}

// 权限检查指令
export const permission = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value, modifiers } = binding
    const userStore = useUserStore()

    if (value) {
      let hasPermission = false

      if (Array.isArray(value)) {
        hasPermission = modifiers.all 
          ? userStore.hasAllPermissions(value)
          : userStore.hasAnyPermission(value)
      } else {
        hasPermission = userStore.hasPermission(value)
      }

      if (!hasPermission) {
        el.style.display = 'none'
        el.setAttribute('data-permission-hidden', 'true')
      }
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const { value, modifiers } = binding
    const userStore = useUserStore()

    if (value) {
      let hasPermission = false

      if (Array.isArray(value)) {
        hasPermission = modifiers.all 
          ? userStore.hasAllPermissions(value)
          : userStore.hasAnyPermission(value)
      } else {
        hasPermission = userStore.hasPermission(value)
      }

      if (hasPermission) {
        el.style.display = ''
        el.removeAttribute('data-permission-hidden')
      } else {
        el.style.display = 'none'
        el.setAttribute('data-permission-hidden', 'true')
      }
    }
  }
}

// 角色检查指令
export const role = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value, modifiers } = binding
    const userStore = useUserStore()

    if (value) {
      let hasRole = false

      if (Array.isArray(value)) {
        hasRole = modifiers.all
          ? userStore.hasAllRoles(value)
          : userStore.hasAnyRole(value)
      } else {
        hasRole = userStore.hasRole(value)
      }

      if (!hasRole) {
        el.style.display = 'none'
        el.setAttribute('data-role-hidden', 'true')
      }
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const { value, modifiers } = binding
    const userStore = useUserStore()

    if (value) {
      let hasRole = false

      if (Array.isArray(value)) {
        hasRole = modifiers.all
          ? userStore.hasAllRoles(value)
          : userStore.hasAnyRole(value)
      } else {
        hasRole = userStore.hasRole(value)
      }

      if (hasRole) {
        el.style.display = ''
        el.removeAttribute('data-role-hidden')
      } else {
        el.style.display = 'none'
        el.setAttribute('data-role-hidden', 'true')
      }
    }
  }
}

// 综合权限检查指令
export const auth = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding

    if (value) {
      const hasAuth = checkAuth(value)
      
      if (!hasAuth) {
        el.style.display = 'none'
        el.setAttribute('data-auth-hidden', 'true')
      }
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding

    if (value) {
      const hasAuth = checkAuth(value)
      
      if (hasAuth) {
        el.style.display = ''
        el.removeAttribute('data-auth-hidden')
      } else {
        el.style.display = 'none'
        el.setAttribute('data-auth-hidden', 'true')
      }
    }
  }
}

// 租户管理员权限指令（简化为直接使用角色检查）
export const admin = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    
    if (!userStore.hasRole('tenant_admin')) {
      el.style.display = 'none'
      el.setAttribute('data-admin-hidden', 'true')
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    
    if (userStore.hasRole('tenant_admin')) {
      el.style.display = ''
      el.removeAttribute('data-admin-hidden')
    } else {
      el.style.display = 'none'
      el.setAttribute('data-admin-hidden', 'true')
    }
  }
}

// 超级管理员权限指令（简化为直接使用角色检查）
export const superAdmin = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    
    if (!userStore.hasRole('super_admin')) {
      el.style.display = 'none'
      el.setAttribute('data-super-admin-hidden', 'true')
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    
    if (userStore.hasRole('super_admin')) {
      el.style.display = ''
      el.removeAttribute('data-super-admin-hidden')
    } else {
      el.style.display = 'none'
      el.setAttribute('data-super-admin-hidden', 'true')
    }
  }
}

// 注册所有指令的函数
export function setupPermissionDirectives(app: App) {
  app.directive('permission', permission)
  app.directive('role', role)
  app.directive('auth', auth)
  app.directive('admin', admin)
  app.directive('super-admin', superAdmin)
}

// Composition API 权限检查 Hook
export function usePermission() {
  const userStore = useUserStore()

  return {
    // 权限检查
    hasPermission: (permission: string) => userStore.hasPermission(permission),
    hasAnyPermission: (permissions: string[]) => userStore.hasAnyPermission(permissions),
    hasAllPermissions: (permissions: string[]) => userStore.hasAllPermissions(permissions),
    
    // 角色检查
    hasRole: (role: string) => userStore.hasRole(role),
    hasAnyRole: (roles: string[]) => userStore.hasAnyRole(roles),
    hasAllRoles: (roles: string[]) => userStore.hasAllRoles(roles),
    
    // 管理员检查（简化为直接使用角色检查）
    isAdmin: () => userStore.hasRole('tenant_admin'),
    isSuperAdmin: () => userStore.hasRole('super_admin'),
    
    // 综合检查
    checkAuth: (config: { permissions?: string[], roles?: string[], requireAll?: boolean }) => 
      checkAuth(config),
    
    // 工具函数
    checkPermission,
    checkAllPermissions,
    checkRole,
    checkAllRoles
  }
}

// 默认导出
export default {
  permission,
  role,
  auth,
  admin,
  superAdmin,
  setupPermissionDirectives,
  usePermission,
  checkPermission,
  checkAllPermissions,
  checkRole,
  checkAllRoles,
  checkAuth
}