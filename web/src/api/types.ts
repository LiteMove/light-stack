// 通用响应接口
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  timestamp: number
}

// 分页响应接口
export interface PageResponse<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 分页请求参数
export interface PageParams {
  page: number
  pageSize: number
  [key: string]: any
}

// 登录请求参数
export interface LoginParams {
  username: string
  password: string
}

// 登录响应数据
export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    email: string
    avatar?: string
  }
}

// 用户信息
export interface User {
  id: number
  tenantId: number
  username: string
  nickname: string
  email: string
  phone?: string
  avatar?: string
  status: number
  isSystem: boolean
  lastLoginAt?: string
  lastLoginIp?: string
  createdAt: string
  updatedAt: string
}

// 角色信息
export interface Role {
  id: number
  tenantId: number
  name: string
  code: string
  description?: string
  status: number
  isSystem: boolean
  sortOrder: number
  createdAt: string
  updatedAt: string
}

// 菜单信息
export interface Menu {
  id: number
  parentId: number
  name: string
  path?: string
  component?: string
  icon?: string
  sortOrder: number
  isHidden: boolean
  isSystem: boolean
  permissionCode?: string
  meta?: any
  children?: Menu[]
  createdAt: string
  updatedAt: string
}