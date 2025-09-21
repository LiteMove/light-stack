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
  page_size: number // Go后端使用下划线命名
  [key: string]: any
}

// 登录请求参数
export interface LoginParams {
  username: string
  password: string
}

// 登录响应数据
export interface TokenResponse {
    accessToken: string
    expiresIn: number
}

// 用户信息
export interface User {
  id: number
  tenantId: number
  tenantName?: string // 租户名称，用于显示
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
  roleCodes: string[]
  permissions: string[]
  menus: Menu[]

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
  code: string
  type: 'directory' | 'menu' | 'permission'
  path?: string
  component?: string
  icon?: string
  resource?: string
  action?: string
  sortOrder: number
  isHidden: boolean
  isSystem: boolean
  status: number
  meta?: string
  children?: Menu[]
  createdAt: string
  updatedAt: string
}

// 菜单创建/更新请求参数
export interface MenuFormData {
  parentId?: number
  name: string
  code: string
  type: 'directory' | 'menu' | 'permission'
  path?: string
  component?: string
  icon?: string
  resource?: string
  action?: string
  sortOrder: number
  isHidden: boolean
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
  menuIds: number[]
}

// 权限信息
export interface Permission {
  id: number
  tenantId: number
  name: string
  code: string
  type: 'api' | 'page' | 'button' | 'data'
  resource?: string
  action?: string
  description?: string
  status: number
  isSystem: boolean
  sortOrder: number
  createdAt: string
  updatedAt: string
}

// 角色权限分配参数
export interface AssignPermissionsData {
  permission_ids: number[]
}

// 角色创建/更新请求参数
export interface RoleFormData {
  name: string
  code: string
  description?: string
  status: number
  sortOrder: number
}

// 角色查询参数
export interface RoleQueryParams {
  page: number
  page_size: number
  keyword?: string
  status?: number
}

// 权限查询参数
export interface PermissionQueryParams {
  page: number
  page_size: number
  keyword?: string
  type?: string
  status?: number
}

// 权限创建/更新请求参数
export interface PermissionFormData {
  name: string
  code: string
  type: 'api' | 'page' | 'button' | 'data'
  resource?: string
  action?: string
  description?: string
  status: number
  sortOrder: number
}

// 用户角色分配参数
export interface AssignUserRolesData {
  roleIds: number[]
}