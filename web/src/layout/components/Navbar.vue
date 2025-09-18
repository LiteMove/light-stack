<template>
  <div class="navbar">
    <div class="navbar-left">
      <hamburger
        :is-active="collapsed"
        class="hamburger-container"
        @toggle-click="toggleSidebar"
      />
      <breadcrumb />
    </div>

    <div class="navbar-right">
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="avatar-container">
          <el-avatar
            :size="36"
            :src="userInfo?.avatar"
            class="user-avatar"
          >
            {{ userInfo?.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="user-name">{{ userInfo?.nickname || userInfo?.username }}</span>
          <el-icon class="caret-down">
            <CaretBottom />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="settings">
              <el-icon><Setting /></el-icon>
              系统设置
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { CaretBottom, User, Setting, SwitchButton } from '@element-plus/icons-vue'
import { useAppStore, useUserStore } from '@/store'
import Hamburger from './Hamburger.vue'
import Breadcrumb from './Breadcrumb.vue'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const collapsed = computed(() => appStore.collapsed)
const userInfo = computed(() => userStore.userInfo)

const toggleSidebar = () => {
  appStore.toggleSidebar()
}

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      // 跳转到个人中心
      break
    case 'settings':
      // 跳转到系统设置
      break
    case 'logout':
      await handleLogout()
      break
  }
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    userStore.logout()
    router.push('/login')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.navbar {
  height: 60px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.navbar-left {
  display: flex;
  align-items: center;
}

.hamburger-container {
  line-height: 46px;
  height: 100%;
  float: left;
  cursor: pointer;
  transition: background 0.3s;
  -webkit-tap-highlight-color: transparent;
  margin-right: 20px;
}

.hamburger-container:hover {
  background: rgba(0, 0, 0, 0.025);
}

.navbar-right {
  display: flex;
  align-items: center;
}

.avatar-container {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.avatar-container:hover {
  background-color: #f5f5f5;
}

.user-avatar {
  margin-right: 10px;
}

.user-name {
  margin-right: 8px;
  font-size: 14px;
  color: #333;
}

.caret-down {
  font-size: 12px;
  color: #999;
}
</style>