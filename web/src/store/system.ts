import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tenantApi } from '@/api/tenant'

export interface SystemConfig {
  id: number
  name: string
  systemName: string
  logo: string
  description: string
  copyright: string
}

export const useSystemStore = defineStore('system', () => {
  const systemConfig = ref<SystemConfig | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 初始化系统配置
  const initSystemConfig = async (domain?: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await tenantApi.getTenantInfo(domain)

      if (response.data) {
        systemConfig.value = response.data
        // 更新页面标题
        if (response.data.systemName) {
          document.title = response.data.systemName
        }

        // 更新favicon（如果有logo）
        if (response.data.logo) {
          updateFavicon(response.data.logo)
        }

        localStorage.setItem('systemConfig', JSON.stringify(response.data))
      }
    } catch (err: any) {
      console.error('获取系统配置失败:', err)
      error.value = err.message || '获取系统配置失败'
      // 如果获取失败，尝试从本地存储恢复
      const stored = localStorage.getItem('systemConfig')
      if (stored) {
        try {
          systemConfig.value = JSON.parse(stored)
        } catch (parseError) {
          console.error('解析本地系统配置失败:', parseError)
        }
      }
    } finally {
      loading.value = false
    }
  }

  // 更新favicon
  const updateFavicon = (logoUrl: string) => {
    try {
      let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
      if (!link) {
        link = document.createElement('link')
        link.rel = 'icon'
        document.getElementsByTagName('head')[0].appendChild(link)
      }
      link.href = logoUrl
    } catch (error) {
      console.error('更新favicon失败:', error)
    }
  }

  // 获取系统名称
  const getSystemName = () => {
    return systemConfig.value?.systemName || systemConfig.value?.name || 'Light Stack'
  }

  // 获取系统logo
  const getSystemLogo = () => {
    return systemConfig.value?.logo || ''
  }

  // 获取版权信息
  const getCopyright = () => {
    return systemConfig.value?.copyright || ''
  }

  // 获取系统描述
  const getDescription = () => {
    return systemConfig.value?.description || ''
  }

  // 清除系统配置
  const clearSystemConfig = () => {
    systemConfig.value = null
    error.value = null
    localStorage.removeItem('systemConfig')
  }

  // 从本地存储恢复配置
  const restoreFromStorage = () => {
    const stored = localStorage.getItem('systemConfig')
    if (stored) {
      try {
        systemConfig.value = JSON.parse(stored)
        // 恢复页面标题
        if (systemConfig.value?.systemName) {
          document.title = systemConfig.value.systemName
        }
        // 恢复favicon
        if (systemConfig.value?.logo) {
          updateFavicon(systemConfig.value.logo)
        }
      } catch (error) {
        console.error('从本地存储恢复系统配置失败:', error)
      }
    }
  }

  return {
    systemConfig,
    loading,
    error,
    initSystemConfig,
    getSystemName,
    getSystemLogo,
    getCopyright,
    getDescription,
    clearSystemConfig,
    restoreFromStorage
  }
})