import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const collapsed = ref(false)
  const device = ref('desktop')
  const size = ref('default')

  // 切换侧边栏
  const toggleSidebar = () => {
    collapsed.value = !collapsed.value
  }

  // 设置设备类型
  const setDevice = (deviceType: string) => {
    device.value = deviceType
  }

  // 设置组件尺寸
  const setSize = (componentSize: string) => {
    size.value = componentSize
  }

  return {
    collapsed,
    device,
    size,
    toggleSidebar,
    setDevice,
    setSize
  }
})