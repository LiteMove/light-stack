import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

// 全局租户监听器
export function useGlobalTenantListener() {
  const router = useRouter()

  const handleTenantChange = (event: CustomEvent) => {
    const { newTenant, oldTenant } = event.detail

    // 避免应用启动时的初始化触发（oldTenant 和 newTenant 都为 undefined/null）
    if (oldTenant === undefined && newTenant === null) return

    console.log('🌐 Global tenant change detected:', {
      from: oldTenant?.name || '未选择',
      to: newTenant?.name || '未选择'
    })

    // 所有页面都需要刷新
    console.log('🔄 Reloading current page due to tenant change')

    // 强制刷新当前页面
    router.go(0)
  }

  onMounted(() => {
    window.addEventListener('tenantChanged', handleTenantChange)
    console.log('🔗 Global tenant listener activated')
  })

  onUnmounted(() => {
    window.removeEventListener('tenantChanged', handleTenantChange)
    console.log('🛑 Global tenant listener deactivated')
  })
}