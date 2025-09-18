import { http } from '@/utils/request'
import type {
  LoginParams,
  LoginResponse,
  User,
  Role,
  Menu,
  MenuFormData,
  MenuQueryParams,
  MenuStatusData,
  BatchMenuStatusData,
  AssignMenusData,
  ApiResponse,
  PageResponse,
  PageParams
} from './types'

// 认证相关API
export const authApi = {
  // 登录
  login(data: LoginParams): Promise<ApiResponse<LoginResponse>> {
    return http.post('/v1/auth/login', data)
  },

  // 获取用户信息
  getUserInfo(): Promise<ApiResponse<User>> {
    return http.get('/v1/auth/user')
  },

  // 登出
  logout(): Promise<ApiResponse> {
    return http.post('/v1/auth/logout')
  }
}

// 用户管理API
export const userApi = {
  // 获取用户列表
  getUsers(params: PageParams): Promise<ApiResponse<PageResponse<User>>> {
    return http.get('/v1/users', { params })
  },

  // 获取用户详情
  getUser(id: number): Promise<ApiResponse<User>> {
    return http.get(`/v1/users/${id}`)
  },

  // 创建用户
  createUser(data: Partial<User>): Promise<ApiResponse<User>> {
    return http.post('/v1/users', data)
  },

  // 更新用户
  updateUser(id: number, data: Partial<User>): Promise<ApiResponse<User>> {
    return http.put(`/v1/users/${id}`, data)
  },

  // 删除用户
  deleteUser(id: number): Promise<ApiResponse> {
    return http.delete(`/v1/users/${id}`)
  }
}

// 角色管理API
export const roleApi = {
  // 获取角色列表
  getRoles(params: PageParams): Promise<ApiResponse<PageResponse<Role>>> {
    return http.get('/v1/roles', { params })
  },

  // 获取角色详情
  getRole(id: number): Promise<ApiResponse<Role>> {
    return http.get(`/v1/roles/${id}`)
  },

  // 创建角色
  createRole(data: Partial<Role>): Promise<ApiResponse<Role>> {
    return http.post('/v1/roles', data)
  },

  // 更新角色
  updateRole(id: number, data: Partial<Role>): Promise<ApiResponse<Role>> {
    return http.put(`/v1/roles/${id}`, data)
  },

  // 删除角色
  deleteRole(id: number): Promise<ApiResponse> {
    return http.delete(`/v1/roles/${id}`)
  }
}

// 菜单管理API
export const menuApi = {
  // 获取菜单列表(分页)
  getMenus(params: MenuQueryParams): Promise<ApiResponse<PageResponse<Menu>>> {
    return http.get('/v1/admin/menus', { params })
  },

  // 获取菜单树
  getMenuTree(): Promise<ApiResponse<Menu[]>> {
    return http.get('/v1/admin/menus/tree')
  },

  // 获取用户菜单树
  getUserMenuTree(): Promise<ApiResponse<Menu[]>> {
    return http.get('/v1/user/menus')
  },

  // 获取用户菜单权限
  getUserPermissions(): Promise<ApiResponse<{ permissions: string[] }>> {
    return http.get('/v1/user/permissions')
  },

  // 获取菜单详情
  getMenu(id: number): Promise<ApiResponse<Menu>> {
    return http.get(`/v1/admin/menus/${id}`)
  },

  // 创建菜单
  createMenu(data: MenuFormData): Promise<ApiResponse<Menu>> {
    return http.post('/v1/admin/menus', data)
  },

  // 更新菜单
  updateMenu(id: number, data: MenuFormData): Promise<ApiResponse<Menu>> {
    return http.put(`/v1/admin/menus/${id}`, data)
  },

  // 删除菜单
  deleteMenu(id: number): Promise<ApiResponse> {
    return http.delete(`/v1/admin/menus/${id}`)
  },

  // 更新菜单状态
  updateMenuStatus(id: number, data: MenuStatusData): Promise<ApiResponse> {
    return http.put(`/v1/admin/menus/${id}/status`, data)
  },

  // 批量更新菜单状态
  batchUpdateMenuStatus(data: BatchMenuStatusData): Promise<ApiResponse> {
    return http.put('/v1/admin/menus/batch/status', data)
  },

  // 获取角色菜单
  getRoleMenus(roleId: number): Promise<ApiResponse<Menu[]>> {
    return http.get(`/v1/admin/roles/${roleId}/menus`)
  },

  // 为角色分配菜单
  assignMenusToRole(roleId: number, data: AssignMenusData): Promise<ApiResponse> {
    return http.put(`/v1/admin/roles/${roleId}/menus`, data)
  }
}

// 系统API
export const systemApi = {
  // 健康检查
  health(): Promise<ApiResponse> {
    return http.get('/health')
  },

  // ping测试
  ping(): Promise<ApiResponse> {
    return http.get('/v1/ping')
  }
}