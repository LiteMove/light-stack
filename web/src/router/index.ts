import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Layout from '@/layout/index.vue'
import { useUserStore } from '@/store'

// 静态路由
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: '登录',
      hidden: true
    }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '首页',
          icon: 'House',
          affix: true
        }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    name: 'System',
    meta: {
      title: '系统管理',
      icon: 'Setting'
    },
    children: [
        {
            path: '/system/tenants',
            name: 'tenants',
            component: () => import('@/views/system/tenants/index.vue'),
            meta: {
                title: '租户管理',
                icon: 'User'
            }
        },
      {
        path: '/system/users',
        name: 'Users',
        component: () => import('@/views/system/users/index.vue'),
        meta: {
          title: '用户管理',
          icon: 'User'
        }
      },
      {
        path: '/system/roles',
        name: 'Roles',
        component: () => import('@/views/system/roles/index.vue'),
        meta: {
          title: '角色管理',
          icon: 'UserFilled'
        }
      },
      {
        path: '/system/permissions',
        name: 'Permissions',
        component: () => import('@/views/system/permissions/index.vue'),
        meta: {
          title: '权限管理',
          icon: 'Key'
        }
      },
      {
        path: '/system/menus',
        name: 'Menus',
        component: () => import('@/views/system/menus/index.vue'),
        meta: {
          title: '菜单管理',
          icon: 'Menu'
        }
      }
    ]
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '404',
      hidden: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
    meta: {
      hidden: true
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes,
  scrollBehavior: () => ({ left: 0, top: 0 })
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  
  // 如果访问登录页面，直接通过
  if (to.path === '/login') {
    next()
    return
  }
  
  // 检查是否有token
  const token = userStore.getToken()
  if (!token) {
    // 没有token，跳转到登录页
    next('/login')
    return
  }
  
  // 确保用户数据已加载
  if (!userStore.userInfo) {
    try {
      await userStore.initUserData()
    } catch (error) {
      console.error('Failed to init user data:', error)
      // 初始化失败，清除token并跳转到登录页
      userStore.logout()
      next('/login')
      return
    }
  }
  
  next()
})

export default router