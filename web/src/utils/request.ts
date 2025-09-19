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
      let tokenString = ''
      
      // 如果token是对象，尝试提取字符串值
      if (typeof token === 'object' && token !== null) {
        // 如果token是一个包含token字段的对象
        if ('access_token' in token) {
          tokenString = String(token.access_token)
        } else {
          console.error('Token is an object but does not contain token field:', token)
          return config
        }
      } else if (typeof token === 'string') {
        tokenString = token
      } else {
        console.error('Token is not a string or valid object:', token)
        return config
      }

      // 确保token字符串有效
      if (tokenString && tokenString.trim() && tokenString !== 'undefined' && tokenString !== 'null') {
        config.headers.Authorization = `Bearer ${tokenString.trim()}`
      } else {
        console.log('Token is empty or invalid after processing:', tokenString)
      }
    } else {
      console.log('No token available')
    }

    // 添加租户ID到请求头
    const currentTenant = tenantStore.getCurrentTenant()
    if (currentTenant) {
      config.headers['X-Tenant-ID'] = currentTenant.id.toString()
      console.log('Added tenant ID to request:', currentTenant.id)
    } else {
      // 如果是超级管理员但没有选择租户，在某些请求中可能需要特殊处理
      const isSuperAdmin = tenantStore.checkIsSuperAdmin()
      if (isSuperAdmin) {
        console.log('Super admin without tenant selection')
        // 可以选择性地添加一个特殊标头来标识这是超级管理员请求
        config.headers['X-Super-Admin'] = 'true'
      }
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
    const { data } = response

    // 如果返回的状态码为200，说明接口请求成功，可以正常拿到数据
    if (data.code === 200) {
      return data
    }

    // 其他状态码都是错误
    return Promise.reject(new Error(data.message || '请求失败'))
  },
  (error) => {
    const { response } = error

    if (response) {
      const { status, data } = response

      switch (status) {
        case 401:
          ElMessage.error('登录过期，请重新登录')
          const userStore = useUserStore()
          userStore.logout()
          router.push('/login')
          break
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(data?.message || '网络错误')
      }
    } else {
      ElMessage.error('网络连接失败')
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