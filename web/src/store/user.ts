import { defineStore } from 'pinia'
import { ref } from 'vue'
import {menuApi, authApi, userApi} from '@/api'
import type { Menu } from '@/api/types'
import type { RouteRecordRaw } from 'vue-router'
import { useTenantStore } from './tenant'
import { resetDynamicRoutes } from '@/router'
import Layout from '@/layout/index.vue'

export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  avatar?: string
  roles: string[]
  permissions: string[]
    menus: Menu[]
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)
  const permissions = ref<string[]>([])
  const userMenus = ref<Menu[]>([])
  const menuPermissions = ref<string[]>([])

  // 设置token
  const setToken = (newToken: string) => {

    // 确保输入是有效的字符串
    if (!newToken || typeof newToken !== 'string') {
      console.error('setToken: Invalid token type:', typeof newToken, newToken)
      return
    }
    
    const tokenString = newToken.trim()
    if (!tokenString) {
      console.error('setToken: Empty token string')
      return
    }
    
    token.value = tokenString
    localStorage.setItem('token', tokenString)
  }

  // 获取token
  const getToken = (): string => {

    // 如果token为空，从本地存储中获取
    if (!token.value) {
      const storedToken = localStorage.getItem('token') || ''

      // 确保从localStorage获取的也是有效字符串
      if (storedToken) {
        token.value = storedToken.trim()
      } else {
        token.value = ''
      }
    }
    
    // 最终验证返回值
    const returnValue = token.value || ''

    // 确保返回值一定是字符串
    return String(returnValue)
  }

  // 清除token
  const clearToken = () => {
    token.value = ''
    localStorage.removeItem('token')
  }

  // 设置用户信息
  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
    permissions.value = info.permissions || []
    userMenus.value = info.menus || []
    // 检查是否为超级管理员（只检查super_admin角色）
    const tenantStore = useTenantStore()
    const isSuperAdminRole = info.roles?.includes('super_admin')
    tenantStore.setIsSuperAdmin(isSuperAdminRole)
    
    // 将用户信息持久化到本地存储
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  // 获取用户信息
  const getUserInfo = async (): Promise<UserInfo | null> => {
    if (!userInfo.value) {
      const storedUserInfo = localStorage.getItem('userInfo')
      if (storedUserInfo) {
        try {
          const parsedUserInfo = JSON.parse(storedUserInfo)
          setUserInfo(parsedUserInfo) // 使用 setUserInfo 来正确设置超级管理员状态
        } catch (error) {
          console.error('Failed to parse stored user info:', error)
        }
      } else {
          // 从API获取用户信息
          try {
              const { data } = await authApi.getUserInfo()

              // 处理用户基本信息
              const userInfoData = {
                  id: data.id,
                  username: data.username,
                  nickname: data.nickname,
                  email: data.email,
                  avatar: data.avatar,
                  roles: data.roleCodes || [],
                  permissions: data.permissions || [],
                  menus: data.menus || []
              }

              setUserInfo(userInfoData)

          } catch (error) {
              console.error('Failed to fetch user info:', error)
              throw error
          }
      }
    }
    return userInfo.value
  }

  // 清除用户信息
  const clearUserInfo = () => {
    userInfo.value = null
    permissions.value = []
    userMenus.value = []
    localStorage.removeItem('userInfo')
    localStorage.removeItem('userMenus')
    localStorage.removeItem('permissions')
  }

  // 获取用户菜单
  const getUserMenus = async (): Promise<Menu[]> => {
    // 首先检查本地存储
    if (!userMenus.value.length) {
      const storedMenus = localStorage.getItem('userMenus')
      if (storedMenus) {
        try {
          userMenus.value = JSON.parse(storedMenus)
          return userMenus.value
        } catch (error) {
          console.error('Failed to parse stored menus:', error)
        }
      }
    }

    // 如果本地没有或解析失败，从API获取
    if (!userMenus.value.length) {
      try {
        const { data } = await menuApi.getUserMenuTree()
        userMenus.value = data || []
        // 保存到本地存储
        localStorage.setItem('userMenus', JSON.stringify(data || []))
        return data || []
      } catch (error) {
        console.error('Failed to fetch user menus:', error)
        return []
      }
    }
    return userMenus.value
  }

  // 获取菜单权限
  const getPermissions = async (): Promise<string[]> => {
    // 首先检查本地存储
    if (!permissions.value.length) {
      const storedPermissions = localStorage.getItem('permissions')
      if (storedPermissions) {
        try {
            permissions.value = JSON.parse(storedPermissions)
          return permissions.value
        } catch (error) {
          console.error('Failed to parse stored permissions:', error)
        }
      }
    }

    // 如果本地没有或解析失败，从API获取
    if (!permissions.value.length) {
      try {
        const { data } = await menuApi.getUserPermissions()
        permissions.value = data.permissions || []
        // 保存到本地存储
        localStorage.setItem('permissions', JSON.stringify(data.permissions || []))
        return data.permissions || []
      } catch (error) {
        console.error('Failed to fetch menu permissions:', error)
        return []
      }
    }

    return permissions.value
  }

  // 检查权限
  const hasPermission = (permission: string): boolean => {
    return permissions.value.includes(permission) || menuPermissions.value.includes(permission)
  }

  // 检查菜单权限
  const hasMenuPermission = (permission: string): boolean => {
    return menuPermissions.value.includes(permission)
  }

  // 检查是否拥有任意一个权限
  const hasAnyPermission = (permissionList: string[]): boolean => {
    return permissionList.some(permission => hasPermission(permission))
  }

  // 检查是否拥有所有权限
  const hasAllPermissions = (permissionList: string[]): boolean => {
    return permissionList.every(permission => hasPermission(permission))
  }

  // 检查角色
  const hasRole = (role: string): boolean => {
    return userInfo.value?.roles?.includes(role) || false
  }

  // 检查是否拥有任意一个角色
  const hasAnyRole = (roleList: string[]): boolean => {
    return roleList.some(role => hasRole(role))
  }

  // 检查是否拥有所有角色
  const hasAllRoles = (roleList: string[]): boolean => {
    return roleList.every(role => hasRole(role))
  }

  // 检查是否为租户管理员
  const isAdmin = (): boolean => {
    return hasRole('tenant_admin')
  }

  // 检查是否为超级管理员
  const isSuperAdmin = (): boolean => {
    return hasRole('super_admin')
  }

  // 初始化用户数据
  const initUserData = async () => {
    if (getToken()) {
      try {
        // 获取用户基本信息
        await getUserInfo()
      } catch (error) {
        console.error('Failed to init user data:', error)
        throw error
      }
    }
  }

  // 登出
  const logout = () => {
    clearToken()
    clearUserInfo()
    // 清除租户数据
    const tenantStore = useTenantStore()
    tenantStore.clearTenantData()
    // 重置动态路由状态
    resetDynamicRoutes()
  }

  // 将菜单转换为路由
  const menuToRoute = (menu: Menu, isChild: boolean = false): RouteRecordRaw => {
      let modules = import.meta.glob('../views/**/*.vue')

      const route: RouteRecordRaw = {
      path: menu.path || `/${menu.code}`,
      name: menu.code,
      meta: {
        title: menu.name,
        icon: menu.icon,
        hidden: menu.isHidden,
        type: menu.type,
        permission: menu.code // 添加权限标识
      }
    }

    // 根据菜单类型设置组件
    if (menu.type === 'directory') {
      // 顶级目录使用Layout组件，子级目录不使用Layout
      if (!isChild) {
        route.component = Layout
        // 如果有子菜单，重定向到第一个可见的子菜单
        if (menu.children && menu.children.length > 0) {
          const firstVisibleChild = menu.children.find(child =>
            !child.isHidden &&
            child.type === 'menu' &&
            child.status === 1
          )
          if (firstVisibleChild) {
            route.redirect = firstVisibleChild.path || `/${firstVisibleChild.code}`
          }
        }
      } else {
        // 子级目录不使用组件，只是路由容器
        route.component = undefined
      }
    } else if (menu.type === 'menu') {

      if (menu.component) {
        const componentPath = `../views/${menu.component}.vue`

        if (!isChild) {
          // 顶级菜单：需要包装在Layout中，作为子路由
          route.component = Layout
          
          // 确保组件被正确加载
          let componentLoader = modules[componentPath]
          if (!componentLoader) {
            // 尝试不同的路径匹配方式
            const altPaths = [
              `../views/${menu.component}.vue`,
              `../views/${menu.component}/index.vue`,
              `../views/${menu.component.toLowerCase()}.vue`,
              `../views/${menu.component.toLowerCase()}/index.vue`
            ]
            
            for (const altPath of altPaths) {
              if (modules[altPath]) {
                componentLoader = modules[altPath]
                break
              }
            }
          }
          
          route.children = [{
            path: '',
            name: menu.code + 'Child',
            component: componentLoader || (() => import('../views/error/404.vue')),
            meta: {
              title: menu.name,
              icon: menu.icon,
              hidden: menu.isHidden,
              type: menu.type,
              permission: menu.code
            }
          }]
        } else {
          // 子级菜单：直接使用组件，不再包装Layout
          let componentLoader = modules[componentPath]
          if (!componentLoader) {
            // 尝试不同的路径匹配方式
            const altPaths = [
              `../views/${menu.component}.vue`,
              `../views/${menu.component}/index.vue`,
              `../views/${menu.component.toLowerCase()}.vue`,
              `../views/${menu.component.toLowerCase()}/index.vue`
            ]
            
            for (const altPath of altPaths) {
              if (modules[altPath]) {
                componentLoader = modules[altPath]
                break
              }
            }
          }
          
          route.component = componentLoader || (() => import('../views/error/404.vue'))
        }

      } else {

        if (!isChild) {
          // 顶级菜单：仍然使用布局包装
          route.component = Layout
          route.children = [{
            path: '',
            name: menu.code + 'Child',
            component: () => import('../views/error/404.vue'),
            meta: {
              title: menu.name,
              icon: menu.icon,
              hidden: menu.isHidden,
              type: menu.type,
              permission: menu.code
            }
          }]
        } else {
          // 子级菜单：直接使用404组件
          route.component = () => import('../views/error/404.vue')
        }
      }
    }

    // 处理子菜单
    if (menu.children && menu.children.length > 0) {
      const childRoutes = menu.children
        .filter(child => {
          const isValid = !child.isHidden && child.status === 1 && child.type !== 'permission'
          return isValid
        })
        .map(child => menuToRoute(child, true)) // 传递isChild=true

      if (childRoutes.length > 0) {
        // 如果当前路由已有children（如顶级menu类型），则合并
        if (route.children && route.children.length > 0) {
          route.children.push(...childRoutes)
        } else {
          route.children = childRoutes
        }
      }
    }

    return route
  }

  // 获取动态路由
  const getDynamicRoutes = (): RouteRecordRaw[] => {
    if (!userMenus.value.length) {
      return []
    }

    // 构建菜单树
    const menuTree = userMenus.value

    // 如果没有可用的菜单，返回空数组
    if (!menuTree.length) {
      return []
    }

    // 转换为路由配置
    const routes = menuTree.map(menu => menuToRoute(menu))
    return routes
  }

  return {
    token,
    userInfo,
    permissions,
    userMenus,
    menuPermissions,
    setToken,
    getToken,
    clearToken,
    setUserInfo,
    getUserInfo,
    clearUserInfo,
    getUserMenus,
    getPermissions,
    hasPermission,
    hasMenuPermission,
    hasAnyPermission,
    hasAllPermissions,
    hasRole,
    hasAnyRole,
    hasAllRoles,
    isAdmin,
    isSuperAdmin,
    initUserData,
    logout,
    getDynamicRoutes
  }
})