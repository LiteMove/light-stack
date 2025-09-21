<template>
  <component v-if="icon" :is="iconComponent" class="menu-icon" />
  <el-icon v-else class="menu-icon">
    <Document />
  </el-icon>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Document } from '@element-plus/icons-vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

interface Props {
  icon?: string
}

const props = defineProps<Props>()

const iconComponent = computed(() => {
  if (!props.icon) return null

  // 常用图标名称映射
  const iconMap: Record<string, string> = {
    'user': 'User',
    'users': 'User', 
    'menu': 'Menu',
    'setting': 'Setting',
    'settings': 'Setting',
    'system': 'Setting',
    'role': 'UserFilled',
    'roles': 'UserFilled',
    'permission': 'Key',
    'permissions': 'Key',
    'dashboard': 'House',
    'home': 'House'
  }

  // 首先尝试映射
  const mappedIconName = iconMap[props.icon.toLowerCase()] || props.icon

  // 如果是Element Plus图标
  const iconComp = (ElementPlusIconsVue as any)[mappedIconName]
  if (iconComp) {
    return iconComp
  }

  // 尝试原始名称（首字母大写）
  const capitalizedIcon = props.icon.charAt(0).toUpperCase() + props.icon.slice(1)
  const iconCompCapitalized = (ElementPlusIconsVue as any)[capitalizedIcon]
  if (iconCompCapitalized) {
    return iconCompCapitalized
  }

  // 记录未找到的图标，但不影响渲染
  console.warn(`Icon not found: ${props.icon}, tried: ${mappedIconName}, ${capitalizedIcon}`)
  return null
})
</script>

<style scoped>
.menu-icon {
  width: 16px;
  height: 16px;
  margin-right: 8px;
}
</style>