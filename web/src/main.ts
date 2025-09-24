import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { setupPermissionDirectives, hasPer, hasRole, hasAnyPer, hasAllPer, hasAnyRole, hasAllRole, isAdmin, isSuperAdmin, hasAuth } from './utils/permission'

import './style.css'

const app = createApp(App)

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 注册权限指令
setupPermissionDirectives(app)

// 注册全局权限检查函数
app.config.globalProperties.$hasPer = hasPer
app.config.globalProperties.$hasRole = hasRole
app.config.globalProperties.$hasAnyPer = hasAnyPer
app.config.globalProperties.$hasAllPer = hasAllPer
app.config.globalProperties.$hasAnyRole = hasAnyRole
app.config.globalProperties.$hasAllRole = hasAllRole
app.config.globalProperties.$isAdmin = isAdmin
app.config.globalProperties.$isSuperAdmin = isSuperAdmin
app.config.globalProperties.$hasAuth = hasAuth

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

app.mount('#app')