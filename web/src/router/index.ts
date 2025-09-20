import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import Layout from '@/layout/index.vue'
import { useUserStore } from '@/store'

// 静态路由（无需权限的基础路由）
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

// 记录是否已添加动态路由
let dynamicRoutesAdded = false

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

  try {
    // 确保用户数据已加载
    if (!userStore.userInfo) {
      await userStore.initUserData()
    }

    // 确保菜单数据已加载
    if (!userStore.userMenus.length) {
      await userStore.getUserMenus()
    }

    console.log('userStore.userMenus:', userStore.userMenus)

    // 检查是否需要添加动态路由
    const needAddRoutes = userStore.userMenus.length > 0 && !dynamicRoutesAdded

    if (needAddRoutes) {
      console.log('Adding dynamic routes...')
      const dynamicRoutes = userStore.getDynamicRoutes()
      console.log('dynamicRoutes:', dynamicRoutes)

      // 添加动态路由到路由器
      dynamicRoutes.forEach(route => {
        router.addRoute(route)
      })

      dynamicRoutesAdded = true
      console.log('Dynamic routes added successfully')

      // 如果当前要访问的路由是动态路由，重新导航
      const isMatchingDynamicRoute = (routes: any[], targetPath: string): boolean => {
        return routes.some(route => {
          // 检查当前路由
          if (route.path && targetPath.startsWith(route.path)) {
            return true
          }
          // 递归检查子路由
          if (route.children && route.children.length > 0) {
            return isMatchingDynamicRoute(route.children, targetPath)
          }
          return false
        })
      }

      if (isMatchingDynamicRoute(dynamicRoutes, to.path)) {
        console.log('Redirecting to dynamic route:', to.path)
        next({ ...to, replace: true })
        return
      }
    }

    // 如果用户数据已加载但没有菜单（可能是权限问题），仍然允许访问基础页面
    if (!userStore.userMenus.length && to.path !== '/' && to.path !== '/dashboard') {
      console.warn('No user menus available, redirecting to dashboard')
      next('/dashboard')
      return
    }

    next()

  } catch (error) {
    console.error('Failed to init user data:', error)
    // 初始化失败，清除token并跳转到登录页
    userStore.logout()
    next('/login')
    return
  }
})

// 重置动态路由状态的方法
export const resetDynamicRoutes = () => {
  dynamicRoutesAdded = false
}

export default router