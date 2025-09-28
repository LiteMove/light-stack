<template>
  <div class="app-wrapper">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside :width="sidebarWidth">
        <sidebar />
      </el-aside>

      <!-- 主内容区 -->
      <el-container direction="vertical">
        <!-- 头部导航 -->
        <el-header>
          <navbar />
        </el-header>

        <!-- 主内容 -->
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/store'
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'

const appStore = useAppStore()

const sidebarWidth = computed(() => {
  return appStore.collapsed ? '64px' : '200px'
})
</script>

<style scoped>
.app-wrapper {
  position: relative;
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

.el-container {
  overflow-y: auto;
}

.el-aside {
  background-color: #304156;
  transition: width 0.28s;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  height: 60px !important;
}

.el-main {
  background-color: #f5f5f5;
  padding: 0;
  height: calc(100vh - 60px);
  overflow: hidden;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .el-aside {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 999;
    height: 100vh;
  }

  .el-container[direction="vertical"] {
    margin-left: 0;
  }

  .el-main {
    padding: 10px;
  }
}

@media (max-width: 480px) {
  .el-main {
    padding: 5px;
  }
}
</style>