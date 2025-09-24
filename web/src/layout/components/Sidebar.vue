<template>
  <div class="sidebar-container">
    <div class="sidebar-logo">
      <router-link to="/">
        <img
          v-if="!collapsed && (systemLogo || hasDefaultLogo)"
          :src="systemLogo || defaultLogoSrc"
          :alt="systemName"
          class="sidebar-logo-img"
          @error="handleLogoError"
        />
        <h1 v-show="!collapsed" class="sidebar-title">{{ systemName }}</h1>
      </router-link>
    </div>

    <div class="sidebar-menu-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="collapsed"
        :unique-opened="false"
        :collapse-transition="false"
        mode="vertical"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <sidebar-item
          v-for="route in routes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
        />
      </el-menu>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore, useUserStore } from '@/store'
import { useSystemStore } from '@/store/system'
import { constantRoutes } from '@/router'
import SidebarItem from './SidebarItem.vue'

const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()
const systemStore = useSystemStore()

const hasDefaultLogo = ref(false)
const defaultLogoSrc = ref('')

// 尝试加载默认logo
try {
  // 使用动态import检查logo文件是否存在
  defaultLogoSrc.value = new URL('@/assets/logo.png', import.meta.url).href
  hasDefaultLogo.value = true
} catch (error) {
  console.log('默认logo文件不存在:', error)
  hasDefaultLogo.value = false
}

const collapsed = computed(() => appStore.collapsed)

const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

// 系统配置
const systemName = computed(() => systemStore.getSystemName())
const systemLogo = computed(() => systemStore.getSystemLogo())

// logo加载错误处理
const handleLogoError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.style.display = 'none'
  hasDefaultLogo.value = false
}

// 静态路由（登录页等）
const staticRoutes = computed(() => {
  return constantRoutes.filter(route => !route.meta?.hidden)
})

// 动态路由（用户菜单）
const dynamicRoutes = computed(() => {
  return userStore.getDynamicRoutes()
})

// 合并所有路由
const routes = computed(() => {
  return [...staticRoutes.value, ...dynamicRoutes.value]
})
</script>

<style scoped>
.sidebar-container {
  height: 100%;
  background-color: #304156;
}

.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b2f3a;
}

.sidebar-logo a {
  display: flex;
  align-items: center;
  text-decoration: none;
}

.sidebar-logo-img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
}

.sidebar-title {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.el-menu {
  border-right: none;
}

.sidebar-menu-wrapper {
  height: calc(100% - 60px);
  overflow-y: auto;
}

.sidebar-menu-wrapper::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu-wrapper::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-menu-wrapper::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}

.sidebar-menu-wrapper::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>