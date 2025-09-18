import type { App, DirectiveBinding } from 'vue'
import { useUserStore } from '@/store'

/**
 * 权限指令
 * 用法：
 * v-permission="'menu_management'"
 * v-permission="['menu_management', 'user_management']"
 * v-role="'admin'"
 * v-role="['admin', 'manager']"
 */

// 权限检查指令
export const permission = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding
    const userStore = useUserStore()

    if (value) {
      let hasPermission = false

      if (Array.isArray(value)) {
        // 数组形式，只要有一个权限即可显示
        hasPermission = value.some(permission => 
          userStore.hasPermission(permission) || userStore.hasMenuPermission(permission)
        )
      } else {
        // 字符串形式
        hasPermission = userStore.hasPermission(value) || userStore.hasMenuPermission(value)
      }

      if (!hasPermission) {
        el.style.display = 'none'
      }
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding
    const userStore = useUserStore()

    if (value) {
      let hasPermission = false

      if (Array.isArray(value)) {
        hasPermission = value.some(permission => 
          userStore.hasPermission(permission) || userStore.hasMenuPermission(permission)
        )
      } else {
        hasPermission = userStore.hasPermission(value) || userStore.hasMenuPermission(value)
      }

      el.style.display = hasPermission ? '' : 'none'
    }
  }
}

// 角色检查指令
export const role = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding
    const userStore = useUserStore()

    if (value) {
      let hasRole = false

      if (Array.isArray(value)) {
        // 数组形式，只要有一个角色即可显示
        hasRole = value.some(role => userStore.hasRole(role))
      } else {
        // 字符串形式
        hasRole = userStore.hasRole(value)
      }

      if (!hasRole) {
        el.style.display = 'none'
      }
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding
    const userStore = useUserStore()

    if (value) {
      let hasRole = false

      if (Array.isArray(value)) {
        hasRole = value.some(role => userStore.hasRole(role))
      } else {
        hasRole = userStore.hasRole(value)
      }

      el.style.display = hasRole ? '' : 'none'
    }
  }
}

// 管理员权限指令
export const admin = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    
    if (!userStore.hasRole('admin') && !userStore.hasRole('super_admin')) {
      el.style.display = 'none'
    }
  },
  updated(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore()
    
    el.style.display = (userStore.hasRole('admin') || userStore.hasRole('super_admin')) ? '' : 'none'
  }
}

// 注册所有指令的函数
export function setupPermissionDirectives(app: App) {
  app.directive('permission', permission)
  app.directive('role', role)
  app.directive('admin', admin)
}

// 默认导出
export default {
  permission,
  role,
  admin,
  setupPermissionDirectives
}