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
  parent_id: number
  name: string
  code: string
  type: 'directory' | 'menu' | 'permission'
  path?: string
  component?: string
  icon?: string
  resource?: string
  action?: string
  sort_order: number
  is_hidden: boolean
  is_system: boolean
  status: number
  meta?: string
  children?: Menu[]
  created_at: string
  updated_at: string
}

// 菜单创建/更新请求参数
export interface MenuFormData {
  parent_id?: number
  name: string
  code: string
  type: 'directory' | 'menu' | 'permission'
  path?: string
  component?: string
  icon?: string
  resource?: string
  action?: string
  sort_order: number
  is_hidden: boolean
  status: number
  meta?: string
}

// 菜单列表查询参数
export interface MenuQueryParams {
  page: number
  page_size: number
  name?: string
  status?: number
}

// 菜单状态更新参数
export interface MenuStatusData {
  status: number
}

// 批量状态更新参数
export interface BatchMenuStatusData {
  ids: number[]
  status: number
}

// 角色菜单分配参数
export interface AssignMenusData {
  menu_ids: number[]
}