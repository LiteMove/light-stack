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
      },
      {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: {
          title: '个人中心',
          icon: 'User',
          hidden: true
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

    // 检查是否需要添加动态路由
    const needAddRoutes = userStore.userMenus.length > 0 && !dynamicRoutesAdded

    if (needAddRoutes) {
      const dynamicRoutes = userStore.getDynamicRoutes()

      // 添加动态路由到路由器
      dynamicRoutes.forEach((route, index) => {
        router.addRoute(route)
      })

      dynamicRoutesAdded = true

      // 检查当前URL是否应该匹配动态路由
      // 如果是404页面，检查当前浏览器URL
      const currentUrl = window.location.pathname
      const pathToCheck = to.path === '/404' ? currentUrl : to.path

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

      const isMatching = isMatchingDynamicRoute(dynamicRoutes, pathToCheck)

      if (isMatching && pathToCheck !== to.path) {
        next({ path: pathToCheck, replace: true })
        return
      }
    }

    // 如果用户数据已加载但没有菜单（可能是权限问题），仍然允许访问基础页面
    if (!userStore.userMenus.length && to.path !== '/' && to.path !== '/dashboard' && to.path !== '/404') {
      next('/dashboard')
      return
    }

    next()

  } catch (error) {
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