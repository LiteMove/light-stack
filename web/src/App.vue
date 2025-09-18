<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/store'

onMounted(() => {
  // 应用启动时，从本地存储恢复用户信息
  const userStore = useUserStore()
  userStore.getUserInfo()
})
</script>

<style scoped>
#app {
  height: 100vh;
  width: 100%;
}
</style>