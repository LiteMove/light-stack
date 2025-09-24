<template>
  <div class="avatar-upload">
    <!-- 头像显示区域 -->
    <div class="avatar-section">
      <div
        class="avatar-container"
        @click="handleClick"
        :class="{ 'is-uploading': uploading }"
      >
        <div v-if="uploading" class="upload-mask">
          <el-icon class="upload-icon spinning"><Loading /></el-icon>
          <div class="upload-text">上传中...</div>
          <div class="upload-progress">
            <el-progress :percentage="uploadProgress" :show-text="false" />
          </div>
        </div>
        <el-avatar
          v-else
          :size="size"
          :src="avatarUrl"
          class="user-avatar"
        >
          <el-icon><UserFilled /></el-icon>
        </el-avatar>
        <div v-if="!uploading" class="avatar-overlay">
          <el-icon><Camera /></el-icon>
          <span>点击更换</span>
        </div>
      </div>
    </div>

    <!-- 隐藏的文件输入框 -->
    <input
      ref="fileInputRef"
      type="file"
      accept="image/*"
      style="display: none"
      @change="handleFileChange"
    />

    <!-- 操作按钮 -->
    <div v-if="avatarUrl && !uploading" class="avatar-actions">
      <el-button size="small" @click.stop="handleClick">
        <el-icon><Upload /></el-icon>
        更换头像
      </el-button>
      <el-button type="danger" size="small" @click.stop="removeAvatar">
        <el-icon><Delete /></el-icon>
        删除头像
      </el-button>
    </div>

    <!-- 提示信息 -->
    <div class="avatar-tips">
      <p>• 支持 JPG、PNG、GIF 格式</p>
      <p>• 文件大小不超过 {{ maxSize }}MB</p>
      <p>• 建议尺寸 {{ size }}x{{ size }} 像素</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserFilled, Camera, Upload, Delete, Loading } from '@element-plus/icons-vue'
import { uploadFile } from '@/api/file'

interface Props {
  modelValue?: string
  size?: number
  maxSize?: number // MB
  usageType?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'success', file: any): void
  (e: 'error', error: string): void
}

const props = withDefaults(defineProps<Props>(), {
  size: 120,
  maxSize: 2,
  usageType: 'avatar'
})

const emit = defineEmits<Emits>()

// 响应式数据
const uploading = ref(false)
const uploadProgress = ref(0)
const fileInputRef = ref<HTMLInputElement>()

// 计算属性
const avatarUrl = computed({
  get: () => props.modelValue || '',
  set: (value) => emit('update:modelValue', value)
})

// 验证文件
const validateFile = (file: File): boolean => {
  // 检查文件大小
  const maxSizeBytes = props.maxSize * 1024 * 1024
  if (file.size > maxSizeBytes) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
    return false
  }

  // 检查文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('头像只能是 JPG、PNG、GIF 格式')
    return false
  }

  return true
}

// 上传文件
const uploadFileHandler = async (file: File) => {
  if (!validateFile(file)) return

  try {
    uploading.value = true
    uploadProgress.value = 0

    // 模拟上传进度
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 90) {
        uploadProgress.value += Math.random() * 30
      }
    }, 200)

    const formData = new FormData()
    formData.append('file', file)
    formData.append('usageType', props.usageType)
    formData.append('isPublic', 'true') // TODO, 头像默认公有，头像通常设为公有

    const response = await uploadFile(formData)

    clearInterval(progressInterval)
    uploadProgress.value = 100

    // 延迟一下让用户看到100%
    await new Promise(resolve => setTimeout(resolve, 500))

    const fileData = response.data
    avatarUrl.value = fileData.accessUrl || fileData.filePath
    emit('success', fileData)
    ElMessage.success('头像上传成功')
  } catch (error: any) {
    ElMessage.error(error.message || '头像上传失败')
    emit('error', error.message || '头像上传失败')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

// 点击上传
const handleClick = () => {
  if (uploading.value) return
  fileInputRef.value?.click()
}

// 文件选择
const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    uploadFileHandler(file)
  }
  // 清空input，允许重复选择同一个文件
  target.value = ''
}

// 删除头像
const removeAvatar = async () => {
  try {
    await ElMessageBox.confirm('确定要删除当前头像吗？', '确认删除', {
      type: 'warning'
    })

    avatarUrl.value = ''
    ElMessage.success('头像删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除头像失败')
    }
  }
}
</script>

<style lang="scss" scoped>
.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;

  .avatar-section {
    position: relative;
    display: flex;
    justify-content: center;

    .avatar-container {
      position: relative;
      cursor: pointer;
      border-radius: 50%;
      overflow: hidden;
      transition: all 0.3s ease;

      &:hover {
        transform: scale(1.05);

        .avatar-overlay {
          opacity: 1;
        }
      }

      &.is-uploading {
        cursor: not-allowed;

        .upload-mask {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.6);
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          border-radius: 50%;
          color: white;
          z-index: 1;

          .upload-icon {
            font-size: 32px;
            margin-bottom: 8px;

            &.spinning {
              animation: rotate 2s linear infinite;
            }
          }

          .upload-text {
            font-size: 14px;
            margin-bottom: 12px;
          }

          .upload-progress {
            width: 80px;
          }
        }
      }

      .user-avatar {
        border: 3px solid #fff;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        background: #f5f7fa;

        :deep(.el-icon) {
          font-size: 48px;
          color: #c0c4cc;
        }
      }

      .avatar-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.6);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        color: white;
        opacity: 0;
        transition: opacity 0.3s ease;
        font-size: 14px;

        .el-icon {
          font-size: 24px;
          margin-bottom: 4px;
        }
      }
    }
  }

  .avatar-actions {
    display: flex;
    gap: 8px;
  }

  .avatar-tips {
    text-align: center;

    p {
      margin: 0 0 4px 0;
      font-size: 12px;
      color: #909399;
      line-height: 1.4;
    }
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

// 响应式设计
@media (max-width: 768px) {
  .avatar-upload {
    .avatar-actions {
      flex-direction: column;
      width: 120px;

      .el-button {
        font-size: 12px;
        padding: 6px 8px;
      }
    }
  }
}
</style>