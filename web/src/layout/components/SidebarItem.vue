<template>
  <div v-if="!item.meta?.hidden">
    <!-- 如果有子菜单且子菜单数量大于1 -->
    <el-sub-menu
      v-if="hasChildren && !hasOneShowingChild"
      :index="resolvePath(item.path)"
    >
      <template #title>
        <menu-icon :icon="item.meta?.icon" />
        <span>{{ item.meta?.title }}</span>
      </template>
      <sidebar-item
        v-for="child in showingChildren"
        :key="child.path"
        :item="child"
        :base-path="resolvePath(child.path)"
      />
    </el-sub-menu>

    <!-- 没有子菜单或只有一个子菜单 -->
    <el-menu-item
      v-else-if="!hasChildren || hasOneShowingChild"
      :index="resolvePath(onlyOneChild.path)"
    >
      <menu-icon :icon="onlyOneChild.meta?.icon || item.meta?.icon" />
      <template #title>
        <span>{{ onlyOneChild.meta?.title || item.meta?.title }}</span>
      </template>
    </el-menu-item>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import MenuIcon from './MenuIcon.vue'

interface Props {
  item: RouteRecordRaw
  basePath: string
}

const props = defineProps<Props>()

const hasChildren = computed(() => {
  return props.item.children && props.item.children.length > 0
})

const showingChildren = computed(() => {
  return props.item.children?.filter(item => !item.meta?.hidden) || []
})

const hasOneShowingChild = computed(() => {
  return showingChildren.value.length === 1
})

const onlyOneChild = computed(() => {
  if (hasOneShowingChild.value) {
    return showingChildren.value[0]
  }
  // 如果没有显示的子菜单，返回当前菜单项作为单独项目
  return { ...props.item, path: '' }
})

const resolvePath = (routePath: string) => {
  if (!routePath) {
    return props.basePath
  }
  if (routePath.startsWith('/')) {
    return routePath
  }
  return `${props.basePath}/${routePath}`.replace(/\/+/g, '/')
}
</script>