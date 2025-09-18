import { defineStore } from 'pinia'
import { ref } from 'vue'

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

  // 设置token
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 获取token
  const getToken = (): string => {
    if (!token.value) {
      token.value = localStorage.getItem('token') || ''
    }
    return token.value
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
  }

  // 清除用户信息
  const clearUserInfo = () => {
    userInfo.value = null
    permissions.value = []
  }

  // 检查权限
  const hasPermission = (permission: string): boolean => {
    return permissions.value.includes(permission)
  }

  // 检查角色
  const hasRole = (role: string): boolean => {
    return userInfo.value?.roles?.includes(role) || false
  }

  // 登出
  const logout = () => {
    clearToken()
    clearUserInfo()
  }

  return {
    token,
    userInfo,
    permissions,
    setToken,
    getToken,
    clearToken,
    setUserInfo,
    clearUserInfo,
    hasPermission,
    hasRole,
    logout
  }
})