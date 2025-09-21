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
        <div class="input-wrapper">
          <div class="icon-wrapper">
            <el-icon class="input-icon"><User /></el-icon>
          </div>
          <el-input
            ref="username"
            v-model="loginForm.username"
            placeholder="用户名"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
            class="styled-input"
          />
        </div>
      </el-form-item>

      <el-form-item prop="password">
        <div class="input-wrapper">
          <div class="icon-wrapper">
            <el-icon class="input-icon"><Lock /></el-icon>
          </div>
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
            class="styled-input"
          />
          <div class="password-toggle" @click="showPwd">
            <el-icon class="toggle-icon">
              <View v-if="passwordType === 'password'" />
              <Hide v-else />
            </el-icon>
          </div>
        </div>
      </el-form-item>
      
      <el-form-item prop="captcha" v-if="showCaptcha">
        <div class="input-wrapper captcha-wrapper">
          <div class="icon-wrapper">
            <el-icon class="input-icon"><Key /></el-icon>
          </div>
          <el-input
            ref="captchaInput"
            v-model="loginForm.captcha"
            placeholder="验证码"
            name="captcha"
            type="text"
            tabindex="3"
            class="styled-input captcha-input"
          />
          <div class="captcha-image" @click="refreshCaptcha">
            <img :src="captchaUrl" alt="验证码" />
          </div>
        </div>
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
    console.log('登录响应:', response)
    if (!response.data) {
      ElMessage.error('登录失败')
      return
    }
    // 只保存token，用户信息将在路由守卫中获取
    userStore.setToken(response.data.accessToken.trim())

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
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(8px);
}

.login-form {
  position: relative;
  width: 440px;
  max-width: 90%;
  padding: 60px 40px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 20px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(12px);
  z-index: 1;
  border: 1px solid rgba(255, 255, 255, 0.2);
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

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  transition: all 0.3s;
  height: 52px;
  width: 100%;
  padding: 0 15px;
  margin-bottom: 5px;
  background-color: #f5f7fa;
}

.input-wrapper:hover {
  border-color: #c0c4cc;
  background-color: #f9fafc;
}

.input-wrapper:focus-within {
  border-color: #667eea;
  background-color: #fff;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  color: #889aa4;
  margin-right: 5px;
}

.input-icon {
  font-size: 18px;
}

.password-toggle {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  cursor: pointer;
  color: #889aa4;
  transition: all 0.3s;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.password-toggle:hover {
  color: #667eea;
}

.toggle-icon {
  font-size: 16px;
}

:deep(.el-form-item) {
  border: none;
  background: transparent;
  border-radius: 0;
  margin-bottom: 28px;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.el-form-item:hover) {
  border-color: transparent;
  box-shadow: none;
}

:deep(.styled-input .el-input) {
  display: block;
  height: 52px;
  width: 100%;
  flex: 1;
}

:deep(.styled-input .el-input input) {
  background: transparent;
  border: 0;
  -webkit-appearance: none;
  border-radius: 0;
  padding: 0;
  color: #2d3748;
  height: 52px;
  caret-color: #667eea;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 0.5px;
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