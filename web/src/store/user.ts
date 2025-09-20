import { defineStore } from 'pinia'
import { ref } from 'vue'
import {menuApi, authApi, userApi} from '@/api'
import type {Menu, User} from '@/api/types'
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
      if (storedToken && typeof storedToken === 'string') {
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
    console.log('User info set:', info)
    // 检查是否为超级管理员
    const tenantStore = useTenantStore()
    const isSuperAdmin = info.roles?.includes('super_admin') || info.roles?.includes('admin')
    tenantStore.setIsSuperAdmin(isSuperAdmin)
    
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
                  roles: data.role_codes || [],
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
    menuPermissions.value = []
    localStorage.removeItem('userInfo')
    localStorage.removeItem('userMenus')
    localStorage.removeItem('menuPermissions')
  }

  // 获取用户菜单
  const getUserMenus = async (): Promise<User> => {
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
        const { data } = await authApi.getUserInfo()
        userMenus.value = data.menus
        // 保存到本地存储
        localStorage.setItem('userMenus', JSON.stringify(data))
        return data
      } catch (error) {
        console.error('Failed to fetch user menus:', error)
        return []
      }
    }

    return userMenus.value
  }

  // 获取菜单权限
  const getMenuPermissions = async (): Promise<string[]> => {
    // 首先检查本地存储
    if (!menuPermissions.value.length) {
      const storedPermissions = localStorage.getItem('menuPermissions')
      if (storedPermissions) {
        try {
          menuPermissions.value = JSON.parse(storedPermissions)
          return menuPermissions.value
        } catch (error) {
          console.error('Failed to parse stored permissions:', error)
        }
      }
    }

    // 如果本地没有或解析失败，从API获取
    if (!menuPermissions.value.length) {
      try {
        const { data } = await authApi.getUserInfo()
        menuPermissions.value = data.permissions
        // 保存到本地存储
        localStorage.setItem('menuPermissions', JSON.stringify(data.permissions))
        return data.permissions
      } catch (error) {
        console.error('Failed to fetch menu permissions:', error)
        return []
      }
    }

    return menuPermissions.value
  }

  // 检查权限
  const hasPermission = (permission: string): boolean => {
    return permissions.value.includes(permission) || menuPermissions.value.includes(permission)
  }

  // 检查菜单权限
  const hasMenuPermission = (menuCode: string): boolean => {
    return menuPermissions.value.includes(menuCode)
  }

  // 检查角色
  const hasRole = (role: string): boolean => {
    return userInfo.value?.roles?.includes(role) || false
  }

  // 初始化用户数据
  const initUserData = async () => {
    if (getToken()) {
      try {
        // 获取用户基本信息
        await getUserInfo()
        // 确保菜单已加载
        if (!userMenus.value.length) {
          await getUserMenus()
        }
        // 确保权限已加载
        if (!menuPermissions.value.length) {
          await getMenuPermissions()
        }
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

  // 构建菜单树
  const buildMenuTree = (menus: Menu[]): Menu[] => {
    if (!menus || menus.length === 0) {
      return []
    }

    const menuMap = new Map<number, Menu>()
    const roots: Menu[] = []

    // 过滤掉隐藏的菜单和权限类型的菜单，只保留目录和菜单类型
    const visibleMenus = menus.filter(menu =>
      !menu.is_hidden &&
      menu.status === 1 &&
      (menu.type === 'directory' || menu.type === 'menu')
    )

    // 先排序菜单
    const sortedMenus = [...visibleMenus].sort((a, b) => a.sort_order - b.sort_order)

    // 将所有菜单按id存储到map中，并初始化children数组
    sortedMenus.forEach(menu => {
      menuMap.set(menu.id, { ...menu, children: [] })
    })

    // 构建树形结构
    sortedMenus.forEach(menu => {
      const menuItem = menuMap.get(menu.id)!
      if (menu.parent_id === 0) {
        roots.push(menuItem)
      } else {
        const parent = menuMap.get(menu.parent_id)
        if (parent) {
          parent.children = parent.children || []
          parent.children.push(menuItem)
        }
      }
    })

    // 递归排序子菜单
    const sortChildren = (items: Menu[]) => {
      items.forEach(item => {
        if (item.children && item.children.length > 0) {
          item.children.sort((a, b) => a.sort_order - b.sort_order)
          sortChildren(item.children)
        }
      })
    }

    sortChildren(roots)

    // 过滤掉没有子菜单的目录类型菜单（除非它们有路径）
    const filterEmptyDirectories = (items: Menu[]): Menu[] => {
      return items.filter(item => {
        // 如果有子菜单，递归过滤子菜单
        if (item.children && item.children.length > 0) {
          item.children = filterEmptyDirectories(item.children)
        }

        // 保留菜单类型的项目
        if (item.type === 'menu') {
          return true
        }

        // 对于目录类型，只有当它有子菜单或者有路径时才保留
        if (item.type === 'directory') {
          return (item.children && item.children.length > 0) || item.path
        }

        return false
      })
    }

    return filterEmptyDirectories(roots)
  }

  // 将菜单转换为路由
  const menuToRoute = (menu: Menu): RouteRecordRaw => {
    const route: RouteRecordRaw = {
      path: menu.path || `/${menu.code}`,
      name: menu.code,
      meta: {
        title: menu.name,
        icon: menu.icon,
        hidden: menu.is_hidden,
        type: menu.type,
        permission: menu.code // 添加权限标识
      }
    }

    // 根据菜单类型设置组件
    if (menu.type === 'directory') {
      // 目录类型使用Layout组件
      route.component = Layout
      // 如果有子菜单，重定向到第一个可见的子菜单
      if (menu.children && menu.children.length > 0) {
        const firstVisibleChild = menu.children.find(child =>
          !child.is_hidden &&
          child.type === 'menu' &&
          child.status === 1
        )
        if (firstVisibleChild) {
          route.redirect = firstVisibleChild.path || `/${firstVisibleChild.code}`
        }
      }
    } else if (menu.type === 'menu' && menu.component) {
      // 菜单类型动态导入组件
      route.component = () => {
        return import(`@/views/${menu.component}.vue`).catch((error) => {
          console.warn(`Failed to load component: @/views/${menu.component}.vue`, error)
          // 如果组件不存在，使用404组件
          return import('@/views/error/404.vue')
        })
      }
    } else if (menu.type === 'menu' && !menu.component) {
      // 如果菜单没有指定组件，使用默认的空组件或404页面
      route.component = () => import('@/views/error/404.vue')
    }

    // 处理子菜单
    if (menu.children && menu.children.length > 0) {
      const childRoutes = menu.children
        .filter(child =>
          !child.is_hidden &&
          child.status === 1 &&
          child.type !== 'permission'
        )
        .map(child => menuToRoute(child))

      if (childRoutes.length > 0) {
        route.children = childRoutes
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
    const menuTree = buildMenuTree(userMenus.value)

    // 如果没有可用的菜单，返回空数组
    if (!menuTree.length) {
      return []
    }

    // 转换为路由配置
    return menuTree.map(menu => menuToRoute(menu))
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
    getMenuPermissions,
    hasPermission,
    hasMenuPermission,
    hasRole,
    initUserData,
    logout,
    buildMenuTree,
    getDynamicRoutes
  }
})