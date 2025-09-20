<template>
  <div class="sidebar-container">
    <div class="sidebar-logo">
      <router-link to="/">
        <img v-show="!collapsed" src="@/assets/logo.png" alt="Logo" class="sidebar-logo-img" />
        <h1 v-show="!collapsed" class="sidebar-title">LightStack</h1>
      </router-link>
    </div>

    <el-scrollbar>
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
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore, useUserStore } from '@/store'
import { constantRoutes } from '@/router'
import SidebarItem from './SidebarItem.vue'

const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()

const collapsed = computed(() => appStore.collapsed)

const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

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
</style>