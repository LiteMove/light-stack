import axios from 'axios'
import type { AxiosResponse, AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store'
import { useTenantStore } from '@/store/tenant'
import router from '@/router'

// 获取环境变量配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'
const APP_TITLE = import.meta.env.VITE_APP_TITLE || 'LightStack Admin'

// 创建axios实例
const request = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    const tenantStore = useTenantStore()
    const token = userStore.getToken()

    // 更严格的token验证和处理
    if (token) {
      // 确保token字符串有效
      if (token && token.trim()) {
        config.headers.Authorization = `Bearer ${token.trim()}`
      } else {
        console.log('Token is empty or invalid after processing:', token)
      }
    } else {
      console.log('No token available')
    }

    // 添加租户ID到请求头
    const currentTenant = tenantStore.getCurrentTenant()
    const isSuperAdmin = tenantStore.checkIsSuperAdmin()

    if (currentTenant) {
      config.headers['X-Tenant-Id'] = currentTenant.id.toString()
      console.log('Added tenant ID to request:', currentTenant.id)
    } else if (isSuperAdmin) {
      // 如果是超级管理员但没有选择租户，在某些请求中可能需要特殊处理
      console.log('Super admin without tenant selection')
      // 对于超级管理员，如果没有选择租户，可能需要默认处理
      // 这里可以根据具体业务需求处理
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 如果是文件下载请求（responseType 为 blob），直接返回响应对象
    if (response.config.responseType === 'blob') {
      return response
    }

    const { data } = response

    // 如果返回的状态码为200，说明接口请求成功，可以正常拿到数据
    if (data.code === 200) {
      return data
    }

    // 其他状态码都是错误，直接显示错误信息并抛出异常
    ElMessage.error(data.message || '请求失败')
    return Promise.reject(new Error(data.message || '请求失败'))
  },
  (error) => {
    const { response } = error

    if (response) {
      const { status, data } = response

      // 否则根据HTTP状态码显示通用错误信息
      let errorMessage = ''
      switch (status) {
        case 401:
          errorMessage = '登录过期，请重新登录'
          console.log('Received 401 response, handling logout and redirect')
          // 使用setTimeout延迟执行，避免在响应处理过程中同步执行logout和路由跳转
          setTimeout(() => {
            try {
              const userStore = useUserStore()
              console.log('Calling userStore.logout()')
              userStore.logout()
              // 确保跳转到登录页
              if (router.currentRoute.value.path !== '/login') {
                console.log('Redirecting to login page from:', router.currentRoute.value.path)
                router.replace('/login')
              } else {
                console.log('Already on login page, no redirect needed')
              }
            } catch (err) {
              console.error('Error during logout/redirect:', err)
            }
          }, 100)
          break
        case 403:
          errorMessage = '权限不足'
          break
        case 404:
          errorMessage = '请求的资源不存在'
          break
        case 500:
          errorMessage = '服务器内部错误'
          break
        default:
          errorMessage = '网络错误'
          break
      }
      if (data && data.message) {
        errorMessage = data.message
      }
      ElMessage.error(errorMessage)

      // 抛出包含正确错误信息的新错误对象
      return Promise.reject(new Error(errorMessage))
    } else {
      const errorMessage = '网络连接失败'
      ElMessage.error(errorMessage)
      return Promise.reject(new Error(errorMessage))
    }

    return Promise.reject(error)
  }
)

// 通用请求方法
interface RequestOptions extends AxiosRequestConfig {
  showLoading?: boolean
  showMessage?: boolean
}

export const http = {
  get<T = any>(url: string, config?: RequestOptions): Promise<T> {
    return request.get(url, config)
  },

  post<T = any>(url: string, data?: any, config?: RequestOptions): Promise<T> {
    return request.post(url, data, config)
  },

  put<T = any>(url: string, data?: any, config?: RequestOptions): Promise<T> {
    return request.put(url, data, config)
  },

  delete<T = any>(url: string, config?: RequestOptions): Promise<T> {
    return request.delete(url, config)
  }
}

export default request

// 导出配置信息供其他模块使用
export const config = {
  API_BASE_URL,
  APP_TITLE
}