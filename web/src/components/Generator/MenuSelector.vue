<template>
  <el-tree-select
    v-model="selectedValue"
    :data="treeData"
    :props="treeProps"
    :placeholder="placeholder"
    :clearable="true"
    :check-strictly="true"
    :render-after-expand="false"
    style="width: 100%"
    @node-click="handleNodeClick"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SystemMenu } from '@/types/generator'

interface Props {
  modelValue?: number
  menus: SystemMenu[]
  placeholder?: string
}

interface Emits {
  (event: 'update:modelValue', value: number | undefined): void
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请选择菜单'
})

const emit = defineEmits<Emits>()

// 计算属性
const selectedValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const treeData = computed(() => {
  return convertToTreeData(props.menus)
})

const treeProps = {
  label: 'name',
  value: 'id',
  children: 'children'
}

// 方法
const convertToTreeData = (menus: SystemMenu[]): any[] => {
  return menus.map(menu => ({
    id: menu.id,
    name: menu.name,
    icon: menu.icon,
    path: menu.path,
    children: menu.children ? convertToTreeData(menu.children) : undefined
  }))
}

const handleNodeClick = (data: any) => {
  emit('update:modelValue', data.id)
}
</script>

<style scoped lang="scss">
:deep(.el-tree-select) {
  width: 100%;
}

:deep(.el-tree-node__content) {
  height: 36px;
  line-height: 36px;
}

:deep(.el-tree-node__label) {
  display: flex;
  align-items: center;
  gap: 6px;
}
</style>