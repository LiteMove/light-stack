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

  // 如果是Element Plus图标
  const iconComp = (ElementPlusIconsVue as any)[props.icon]
  if (iconComp) {
    return iconComp
  }

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