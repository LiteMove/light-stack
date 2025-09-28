<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/store'
import { useSystemStore } from '@/store/system'
import { useGlobalTenantListener } from '@/composables/useGlobalTenantListener'

// 启用全局租户监听
useGlobalTenantListener()

onMounted(async () => {
  // 应用启动时，从本地存储恢复用户信息
  const userStore = useUserStore()
  userStore.getUserInfo()

  // 初始化系统配置
  const systemStore = useSystemStore()
  // 先尝试从本地存储恢复
  systemStore.restoreFromStorage()
  // 然后异步获取最新配置
  try {
    await systemStore.initSystemConfig()
  } catch (error) {
    console.error('初始化系统配置失败:', error)
    // 即使获取失败也不影响应用启动
  }
})
</script>

<style scoped>
#app {
  height: 100vh;
  width: 100%;
}
</style>