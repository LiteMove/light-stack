import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tenantApi } from '@/api/tenant'

export interface Tenant {
  id: number
  name: string
  domain?: string
  status: number
  expiredAt?: string
  createdAt: string
}

export const useTenantStore = defineStore('tenant', () => {
  const currentTenant = ref<Tenant | null>(null)
  const tenantList = ref<Tenant[]>([])
  const isSuperAdmin = ref(false)

  // 设置当前租户
  const setCurrentTenant = (tenant: Tenant | null) => {
    currentTenant.value = tenant
    if (tenant) {
      localStorage.setItem('currentTenant', JSON.stringify(tenant))
    } else {
      localStorage.removeItem('currentTenant')
    }
  }

  // 获取当前租户
  const getCurrentTenant = (): Tenant | null => {
    if (!currentTenant.value) {
      const stored = localStorage.getItem('currentTenant')
      if (stored) {
        try {
          currentTenant.value = JSON.parse(stored)
        } catch (error) {
          console.error('Failed to parse stored tenant:', error)
        }
      }
    }
    return currentTenant.value
  }

  // 设置超级管理员状态
  const setIsSuperAdmin = (value: boolean) => {
    isSuperAdmin.value = value
    localStorage.setItem('isSuperAdmin', String(value))
  }

  // 检查是否为超级管理员
  const checkIsSuperAdmin = (): boolean => {
    if (!isSuperAdmin.value) {
      const stored = localStorage.getItem('isSuperAdmin')
      isSuperAdmin.value = stored === 'true'
    }
    return isSuperAdmin.value
  }

  // 获取租户列表
  const fetchTenantList = async () => {
    if (!checkIsSuperAdmin()) return []

    try {
      const response = await tenantApi.getTenantSelectList()
      
      if (response.data) {
        tenantList.value = response.data || []
      }
      return tenantList.value
    } catch (error) {
      console.error('Failed to fetch tenant list:', error)
      return []
    }
  }

  // 切换租户
  const switchTenant = (tenant: Tenant | null) => {
    setCurrentTenant(tenant)
  }

  // 清除租户数据
  const clearTenantData = () => {
    currentTenant.value = null
    tenantList.value = []
    isSuperAdmin.value = false
    localStorage.removeItem('currentTenant')
    localStorage.removeItem('isSuperAdmin')
  }

  // 获取当前租户ID（用于API请求）
  const getCurrentTenantId = (): number | null => {
    const tenant = getCurrentTenant()
    return tenant ? tenant.id : null
  }

  return {
    currentTenant,
    tenantList,
    isSuperAdmin,
    setCurrentTenant,
    getCurrentTenant,
    setIsSuperAdmin,
    checkIsSuperAdmin,
    fetchTenantList,
    switchTenant,
    clearTenantData,
    getCurrentTenantId
  }
})