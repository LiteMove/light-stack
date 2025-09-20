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
  console.log('[ROUTER] Navigation from:', from.path, 'to:', to.path)
  const userStore = useUserStore()

  // 如果访问登录页面，直接通过
  if (to.path === '/login') {
    console.log('[ROUTER] Accessing login page, allowing...')
    next()
    return
  }

  // 检查是否有token
  const token = userStore.getToken()
  console.log('[ROUTER] Token check:', !!token)
  if (!token) {
    console.log('[ROUTER] No token, redirecting to login')
    next('/login')
    return
  }

  try {
    // 确保用户数据已加载
    console.log('[ROUTER] Current userInfo:', !!userStore.userInfo)
    if (!userStore.userInfo) {
      console.log('[ROUTER] Loading user info...')
      await userStore.initUserData()
      console.log('[ROUTER] User info loaded:', !!userStore.userInfo)
    }

    // 确保菜单数据已加载
    console.log('[ROUTER] Current userMenus length:', userStore.userMenus.length)
    if (!userStore.userMenus.length) {
      console.log('[ROUTER] Loading user menus...')
      await userStore.getUserMenus()
      console.log('[ROUTER] User menus loaded:', userStore.userMenus.length)
    }

    // 打印菜单数据用于调试
    console.log('[ROUTER] userStore.userMenus:', JSON.stringify(userStore.userMenus, null, 2))

    // 检查是否需要添加动态路由
    const needAddRoutes = userStore.userMenus.length > 0 && !dynamicRoutesAdded
    console.log('[ROUTER] Need add routes:', needAddRoutes, 'dynamicRoutesAdded:', dynamicRoutesAdded)

    if (needAddRoutes) {
      console.log('[ROUTER] Adding dynamic routes...')
      const dynamicRoutes = userStore.getDynamicRoutes()
      console.log('[ROUTER] Generated dynamic routes:', JSON.stringify(dynamicRoutes, null, 2))

      // 添加动态路由到路由器
      dynamicRoutes.forEach((route, index) => {
        console.log(`[ROUTER] Adding route ${index}:`, route.path, route.name)
        router.addRoute(route)
      })

      dynamicRoutesAdded = true
      console.log('[ROUTER] Dynamic routes added successfully')

      // 打印当前所有路由
      console.log('[ROUTER] All routes after adding:', router.getRoutes().map(r => ({ name: r.name, path: r.path })))

      // 检查当前URL是否应该匹配动态路由
      // 如果是404页面，检查当前浏览器URL
      const currentUrl = window.location.pathname
      const pathToCheck = to.path === '/404' ? currentUrl : to.path
      console.log('[ROUTER] Checking path:', pathToCheck, 'current URL:', currentUrl)

      const isMatchingDynamicRoute = (routes: any[], targetPath: string): boolean => {
        return routes.some(route => {
          console.log(`[ROUTER] Checking route path: ${route.path} against target: ${targetPath}`)
          // 检查当前路由
          if (route.path && targetPath.startsWith(route.path)) {
            console.log(`[ROUTER] Found matching route: ${route.path}`)
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
      console.log('[ROUTER] Is matching dynamic route:', isMatching, 'for path:', pathToCheck)

      if (isMatching && pathToCheck !== to.path) {
        console.log('[ROUTER] Redirecting to dynamic route:', pathToCheck)
        next({ path: pathToCheck, replace: true })
        return
      }
    }

    // 如果用户数据已加载但没有菜单（可能是权限问题），仍然允许访问基础页面
    if (!userStore.userMenus.length && to.path !== '/' && to.path !== '/dashboard' && to.path !== '/404') {
      console.warn('[ROUTER] No user menus available, redirecting to dashboard')
      next('/dashboard')
      return
    }

    console.log('[ROUTER] Navigation allowed to:', to.path)
    next()

  } catch (error) {
    console.error('[ROUTER] Failed to init user data:', error)
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