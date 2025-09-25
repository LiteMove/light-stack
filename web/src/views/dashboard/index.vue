<template>
  <div class="dashboard">
    <!-- 欢迎信息 -->
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template v-if="userRole === 'user'">
            <div class="welcome-section">
              <h2>{{ dashboardData.welcomeText }}</h2>
              <div class="user-info">
                <p><strong>用户名：</strong>{{ dashboardData.username }}</p>
                <p><strong>所属租户：</strong>{{ dashboardData.tenantName }}</p>
                <p v-if="dashboardData.lastLogin"><strong>上次登录：</strong>{{ dashboardData.lastLogin }}</p>
              </div>
            </div>
          </template>
          <template v-else-if="userRole === 'super_admin'">
            <h2>超级管理员控制台</h2>
            <p>欢迎使用 Light Stack 系统管理平台</p>
          </template>
          <template v-else>
            <h2>租户管理控制台</h2>
            <p>欢迎使用 Light Stack 租户管理系统</p>
          </template>
        </el-card>
      </el-col>
    </el-row>

    <!-- 统计卡片 - 根据角色显示不同内容 -->
    <el-row :gutter="20" style="margin-top: 20px" v-if="userRole !== 'user'">
      <!-- 超级管理员统计 -->
      <template v-if="userRole === 'super_admin'">
        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon user">
                <el-icon><User /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.userCount || 0 }}</div>
                <div class="stat-label">用户总数</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon tenant">
                <el-icon><OfficeBuilding /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.tenantCount || 0 }}</div>
                <div class="stat-label">租户总数</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon file">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.fileCount || 0 }}</div>
                <div class="stat-label">文件总数</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="6">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon log">
                <el-icon><List /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.logCount || 0 }}</div>
                <div class="stat-label">日志总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </template>

      <!-- 租户管理员统计 -->
      <template v-else-if="userRole === 'admin'">
        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon user">
                <el-icon><User /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.userCount || 0 }}</div>
                <div class="stat-label">租户用户数</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon file">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.fileCount || 0 }}</div>
                <div class="stat-label">文件总数</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="stat-card">
            <div class="stat-item">
              <div class="stat-icon log">
                <el-icon><List /></el-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">{{ dashboardData.logCount || 0 }}</div>
                <div class="stat-label">日志总数</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </template>
    </el-row>

    <!-- 租户管理员的租户信息 -->
    <el-row :gutter="20" style="margin-top: 20px" v-if="userRole === 'admin'">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>租户信息</span>
          </template>
          <div class="tenant-info">
            <div class="info-item">
              <span class="label">租户名称：</span>
              <span class="value">{{ dashboardData.tenantName }}</span>
            </div>
            <div class="info-item" v-if="dashboardData.expiredAt">
              <span class="label">到期时间：</span>
              <span class="value" :class="getExpireStatusClass()">{{ formatDate(dashboardData.expiredAt) }}</span>
            </div>
            <div class="info-item">
              <span class="label">租户状态：</span>
              <el-tag :type="getStatusType(dashboardData.status)">{{ getStatusText(dashboardData.status) }}</el-tag>
            </div>
            <div class="info-item" v-if="dashboardData.storageUsed !== undefined">
              <span class="label">存储使用：</span>
              <span class="value">{{ formatBytes(dashboardData.storageUsed) }} / {{ formatBytes(dashboardData.storageLimit) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 超级管理员的系统信息 -->
    <el-row :gutter="20" style="margin-top: 20px" v-if="userRole === 'super_admin'">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>系统信息</span>
          </template>
          <div class="system-info" v-if="systemInfo">
            <div class="info-item">
              <span class="label">系统版本：</span>
              <span class="value">{{ systemInfo.version }}</span>
            </div>
            <div class="info-item">
              <span class="label">Go版本：</span>
              <span class="value">{{ systemInfo.goVersion }}</span>
            </div>
            <div class="info-item">
              <span class="label">启动时间：</span>
              <span class="value">{{ systemInfo.startTime }}</span>
            </div>
            <div class="info-item">
              <span class="label">运行时间：</span>
              <span class="value">{{ systemInfo.uptime }}</span>
            </div>
            <div class="info-item">
              <span class="label">操作系统：</span>
              <span class="value">{{ systemInfo.osInfo }}</span>
            </div>
            <div class="info-item">
              <span class="label">服务器时间：</span>
              <span class="value">{{ systemInfo.serverTime }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { User, OfficeBuilding, Document, List } from '@element-plus/icons-vue'
import { getDashboardStats, getSystemInfo } from '@/api/dashboard'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const dashboardData = ref<any>({})
const systemInfo = ref<any>(null)
const loading = ref(false)

// 计算用户角色
const userRole = computed(() => {
  const roles = userStore.userInfo?.roles || []
  if (roles.includes('super_admin')) return 'super_admin'
  if (roles.includes('admin') || roles.includes('tenant_admin')) return 'admin'
  return 'user'
})

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    loading.value = true
    const response = await getDashboardStats()
    dashboardData.value = response.data
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
    ElMessage.error('获取仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

// 获取系统信息（超管专用）
const fetchSystemInfo = async () => {
  if (userRole.value !== 'super_admin') return

  try {
    const response = await getSystemInfo()
    systemInfo.value = response.data
  } catch (error) {
    console.error('获取系统信息失败:', error)
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 格式化字节大小
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 获取状态类型
const getStatusType = (status: number) => {
  const statusMap: Record<number, string> = {
    1: 'success',  // 启用
    2: 'danger',   // 禁用
    3: 'warning',  // 试用
    4: 'info'      // 过期
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status: number) => {
  const statusMap: Record<number, string> = {
    1: '启用',
    2: '禁用',
    3: '试用',
    4: '过期'
  }
  return statusMap[status] || '未知'
}

// 获取过期状态样式类
const getExpireStatusClass = () => {
  if (!dashboardData.value.expiredAt) return ''
  const expireTime = new Date(dashboardData.value.expiredAt).getTime()
  const now = Date.now()
  const daysLeft = (expireTime - now) / (1000 * 60 * 60 * 24)

  if (daysLeft < 0) return 'expired'
  if (daysLeft < 7) return 'warning'
  if (daysLeft < 30) return 'notice'
  return ''
}

onMounted(async () => {
  // 确保用户信息已加载
  if (!userStore.userInfo) {
    try {
      await userStore.getUserInfo()
    } catch (error) {
      console.error('Failed to load user info:', error)
    }
  }

  fetchDashboardData()
  if (userRole.value === 'super_admin') {
    fetchSystemInfo()
  }
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: calc(100vh - 60px);
}

.welcome-section {
  text-align: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
  position: relative;
  overflow: hidden;
}

.welcome-section::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  animation: float 6s ease-in-out infinite;
}

.welcome-section h2 {
  font-size: 2.5rem;
  margin: 0 0 20px 0;
  text-shadow: 0 2px 4px rgba(0,0,0,0.3);
}

.user-info {
  margin-top: 20px;
  text-align: left;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
  background: rgba(255,255,255,0.1);
  padding: 20px;
  border-radius: 10px;
  backdrop-filter: blur(10px);
}

.user-info p {
  margin: 12px 0;
  font-size: 1.1rem;
}

.stat-card {
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border: none;
  border-radius: 16px;
  overflow: hidden;
  position: relative;
  background: white;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2);
}

.stat-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.15);
}

.stat-item {
  display: flex;
  align-items: center;
  padding: 25px 20px;
  position: relative;
}

.stat-icon {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  font-size: 28px;
  color: white;
  position: relative;
  box-shadow: 0 4px 15px rgba(0,0,0,0.2);
}

.stat-icon::after {
  content: '';
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  background: inherit;
  z-index: -1;
  opacity: 0.3;
  filter: blur(8px);
}

.stat-icon.user {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.tenant {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-icon.file {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.log {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 36px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 15px;
  color: #8892a6;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.tenant-info,
.system-info {
  padding: 30px;
  background: white;
  border-radius: 16px;
  position: relative;
}

.tenant-info::before,
.system-info::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #43e97b, #38f9d7);
  border-radius: 16px 16px 0 0;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f2f5;
}

.info-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.info-item .label {
  font-weight: 600;
  color: #4a5568;
  min-width: 120px;
  font-size: 15px;
}

.info-item .value {
  color: #2d3748;
  font-weight: 500;
  font-size: 15px;
}

.info-item .value.warning {
  color: #ed8936;
  font-weight: 600;
}

.info-item .value.expired {
  color: #e53e3e;
  font-weight: 600;
}

.info-item .value.notice {
  color: #3182ce;
  font-weight: 600;
}

/* 动画效果 */
@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

/* 卡片进入动画 */
.stat-card {
  animation: slideInUp 0.6s ease-out;
}

.stat-card:nth-child(1) { animation-delay: 0.1s; }
.stat-card:nth-child(2) { animation-delay: 0.2s; }
.stat-card:nth-child(3) { animation-delay: 0.3s; }
.stat-card:nth-child(4) { animation-delay: 0.4s; }

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 欢迎卡片 */
.el-card {
  border: none;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.el-card .el-card__header {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  border-radius: 16px 16px 0 0;
  font-weight: 600;
  color: #495057;
}

/* 超级管理员和租户管理员标题样式 */
.dashboard h2 {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 700;
  margin: 0 0 10px 0;
  font-size: 2rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard {
    padding: 10px;
  }

  .welcome-section h2 {
    font-size: 1.8rem;
  }

  .user-info {
    max-width: 100%;
    padding: 15px;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }

  .stat-number {
    font-size: 28px;
  }

  .tenant-info,
  .system-info {
    padding: 20px;
  }

  .info-item .label {
    min-width: 100px;
    font-size: 14px;
  }

  .info-item .value {
    font-size: 14px;
  }
}
</style>