<template>
  <div class="login-container">
    <el-form
      ref="loginFormRef"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      autocomplete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">LightStack 管理平台</h3>
      </div>

      <el-form-item prop="username">
        <span class="svg-container">
          <el-icon><User /></el-icon>
        </span>
        <el-input
          ref="username"
          v-model="loginForm.username"
          placeholder="用户名"
          name="username"
          type="text"
          tabindex="1"
          autocomplete="on"
        />
      </el-form-item>

      <el-form-item prop="password">
        <span class="svg-container">
          <el-icon><Lock /></el-icon>
        </span>
        <el-input
          :key="passwordType"
          ref="password"
          v-model="loginForm.password"
          :type="passwordType"
          placeholder="密码"
          name="password"
          tabindex="2"
          autocomplete="on"
          @keyup.enter="handleLogin"
        />
        <span class="show-pwd" @click="showPwd">
          <el-icon>
            <View v-if="passwordType === 'password'" />
            <Hide v-else />
          </el-icon>
        </span>
      </el-form-item>

      <el-button
        :loading="loading"
        type="primary"
        style="width: 100%; margin-bottom: 20px; height: 48px; font-size: 16px; font-weight: 600; border-radius: 12px; background: linear-gradient(135deg, #667eea, #764ba2); border: none;"
        @click.prevent="handleLogin"
      >
        {{ loading ? '登录中...' : '登录' }}
      </el-button>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'
import { useUserStore } from '@/store'
import { authApi } from '@/api'
import type { FormInstance } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const passwordType = ref('password')

const loginForm = reactive({
  username: 'admin',
  password: 'admin123'
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const showPwd = () => {
  if (passwordType.value === 'password') {
    passwordType.value = ''
  } else {
    passwordType.value = 'password'
  }
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    const response = await authApi.login(loginForm)
    const { token, user } = response.data

    // 保存token和用户信息
    userStore.setToken(token)
    userStore.setUserInfo({
      id: user.id,
      username: user.username,
      nickname: user.nickname,
      email: user.email,
      avatar: user.avatar,
      roles: [], // 后续从API获取
      permissions: [] // 后续从API获取
    })

    ElMessage.success('登录成功')
    router.push('/')
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(5px);
}

.login-form {
  position: relative;
  width: 420px;
  max-width: 90%;
  padding: 60px 40px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  z-index: 1;
}

.title-container {
  position: relative;
  margin-bottom: 40px;
}

.title-container .title {
  font-size: 28px;
  color: #2d3748;
  margin: 0;
  text-align: center;
  font-weight: 700;
  letter-spacing: 1px;
}

.svg-container {
  padding: 0 15px 0 0;
  color: #667eea;
  vertical-align: middle;
  width: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.show-pwd {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 18px;
  color: #a0aec0;
  cursor: pointer;
  user-select: none;
  transition: color 0.3s ease;
}

.show-pwd:hover {
  color: #667eea;
}

:deep(.el-form-item) {
  border: 1px solid #e2e8f0;
  background: #fff;
  border-radius: 12px;
  margin-bottom: 24px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

:deep(.el-form-item:hover) {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
}

:deep(.el-input) {
  display: inline-block;
  height: 48px;
  width: calc(100% - 40px);
}

:deep(.el-input input) {
  background: transparent;
  border: 0;
  -webkit-appearance: none;
  border-radius: 0;
  padding: 14px 0;
  color: #2d3748;
  height: 48px;
  caret-color: #667eea;
  font-size: 16px;
}

:deep(.el-input input::placeholder) {
  color: #a0aec0;
}

:deep(.el-input input:focus) {
  outline: none;
  box-shadow: none;
}

:deep(.el-input input:-webkit-autofill) {
  box-shadow: 0 0 0px 1000px #fff inset !important;
  -webkit-text-fill-color: #2d3748 !important;
}
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-form {
  animation: fadeInUp 0.6s ease-out;
}

:deep(.el-button--primary:hover) {
  opacity: 0.9;
  transform: translateY(-1px);
  transition: all 0.2s ease;
}

:deep(.el-button--primary:active) {
  transform: translateY(0);
}

@media (max-width: 480px) {
  .login-form {
    padding: 40px 24px;
    margin: 20px;
  }
  
  .title-container .title {
    font-size: 24px;
  }
  
  :deep(.el-form-item) {
    margin-bottom: 20px;
  }
}</style>