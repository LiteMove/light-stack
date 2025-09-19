import { defineStore } from 'pinia'
import { ref } from 'vue'
import { menuApi, authApi } from '@/api'
import type { Menu } from '@/api/types'
import { useTenantStore } from './tenant'

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
          setUserInfo(data)
        } catch (error) {
          console.error('Failed to fetch user info:', error)
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
  }

  // 获取用户菜单
  const getUserMenus = async (): Promise<Menu[]> => {
    try {
      const { data } = await menuApi.getUserMenuTree()
      userMenus.value = data
      return data
    } catch (error) {
      console.error('Failed to fetch user menus:', error)
      return []
    }
  }

  // 获取菜单权限
  const getMenuPermissions = async (): Promise<string[]> => {
    try {
      const { data } = await menuApi.getUserPermissions()
      menuPermissions.value = data.permissions
      return data.permissions
    } catch (error) {
      console.error('Failed to fetch menu permissions:', error)
      return []
    }
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
      await getUserInfo()
      await getUserMenus()
      await getMenuPermissions()
    }
  }

  // 登出
  const logout = () => {
    clearToken()
    clearUserInfo()
    // 清除租户数据
    const tenantStore = useTenantStore()
    tenantStore.clearTenantData()
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
    logout
  }
})