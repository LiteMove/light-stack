<template>
  <div class="image-upload">
    <!-- 图片显示区域 -->
    <div class="image-section">
      <div
        class="image-container"
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
          :src="imageUrl"
          class="user-image"
          :class="imageClass"
        >
          <el-icon><UserFilled /></el-icon>
        </el-avatar>
        <div v-if="!uploading" class="image-overlay">
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
    <div v-if="imageUrl && !uploading" class="image-actions">
      <el-button size="small" @click.stop="handleClick">
        <el-icon><Upload /></el-icon>
        {{ uploadButtonText }}
      </el-button>
      <el-button type="danger" size="small" @click.stop="removeImage">
        <el-icon><Delete /></el-icon>
        {{ deleteButtonText }}
      </el-button>
    </div>

    <!-- 提示信息 -->
    <div class="image-tips">
      <p>• 支持 JPG、PNG、GIF 格式</p>
      <p>• 文件大小不超过 {{ maxSize }}MB</p>
      <p>• 建议尺寸 {{ size }}x{{ size }} 像素</p>
      <p v-if="usageType === 'system-logo'">• 系统Logo将默认设为公开访问</p>
      <p v-if="usageType === 'avatar'">• 头像将默认设为公开访问</p>
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
const imageUrl = computed({
  get: () => props.modelValue || '',
  set: (value) => emit('update:modelValue', value)
})

// 根据使用类型动态设置按钮文本
const uploadButtonText = computed(() => {
  const textMap: Record<string, string> = {
    'avatar': '更换头像',
    'system-logo': '更换Logo',
    'banner': '更换横幅',
    'icon': '更换图标'
  }
  return textMap[props.usageType] || '更换图片'
})

const deleteButtonText = computed(() => {
  const textMap: Record<string, string> = {
    'avatar': '删除头像',
    'system-logo': '删除Logo',
    'banner': '删除横幅',
    'icon': '删除图标'
  }
  return textMap[props.usageType] || '删除图片'
})

// 根据使用类型设置图片样式类
const imageClass = computed(() => {
  const classMap: Record<string, string> = {
    'system-logo': 'logo-style',
    'banner': 'banner-style',
    'icon': 'icon-style'
  }
  return classMap[props.usageType] || ''
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
    const fileTypeName = props.usageType === 'avatar' ? '头像' :
                        props.usageType === 'system-logo' ? '系统Logo' : '图片'
    ElMessage.error(`${fileTypeName}只能是 JPG、PNG、GIF 格式`)
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

    // 根据使用类型设置默认公开访问权限
    const shouldBePublic = props.usageType === 'system-logo' || props.usageType === 'avatar'
    formData.append('isPublic', shouldBePublic.toString())

    const response = await uploadFile(formData)

    clearInterval(progressInterval)
    uploadProgress.value = 100

    // 延迟一下让用户看到100%
    await new Promise(resolve => setTimeout(resolve, 500))

    const fileData = response.data
    imageUrl.value = fileData.accessUrl || fileData.filePath
    emit('success', fileData)

    const fileTypeName = props.usageType === 'avatar' ? '头像' :
                        props.usageType === 'system-logo' ? '系统Logo' : '图片'
    ElMessage.success(`${fileTypeName}上传成功`)
  } catch (error: any) {
    const fileTypeName = props.usageType === 'avatar' ? '头像' :
                        props.usageType === 'system-logo' ? '系统Logo' : '图片'
    ElMessage.error(error.message || `${fileTypeName}上传失败`)
    emit('error', error.message || `${fileTypeName}上传失败`)
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

// 删除图片
const removeImage = async () => {
  try {
    const fileTypeName = props.usageType === 'avatar' ? '头像' :
                        props.usageType === 'system-logo' ? 'Logo' : '图片'

    await ElMessageBox.confirm(`确定要删除当前${fileTypeName}吗？`, '确认删除', {
      type: 'warning'
    })

    imageUrl.value = ''
    ElMessage.success(`${fileTypeName}删除成功`)
  } catch (error) {
    if (error !== 'cancel') {
      const fileTypeName = props.usageType === 'avatar' ? '头像' :
                          props.usageType === 'system-logo' ? 'Logo' : '图片'
      ElMessage.error(`删除${fileTypeName}失败`)
    }
  }
}
</script>

<style lang="scss" scoped>
.image-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;

  .image-section {
    position: relative;
    display: flex;
    justify-content: center;

    .image-container {
      position: relative;
      cursor: pointer;
      border-radius: 50%; // 默认圆形，会被 logo-style 等覆盖
      overflow: hidden;
      transition: all 0.3s ease;

      // 根据使用类型调整容器形状
      &:has(.logo-style) {
        border-radius: 8px;
      }

      &:has(.banner-style) {
        border-radius: 12px;
      }

      &:has(.icon-style) {
        border-radius: 6px;
      }

      &:hover {
        transform: scale(1.05);

        .image-overlay {
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
          border-radius: inherit; // 继承父元素的圆角
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

      .user-image {
        border: 3px solid #fff;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        background: #f5f7fa;

        // 系统Logo样式
        &.logo-style {
          border-radius: 8px; // Logo通常用圆角矩形
          background: #ffffff;
          border: 2px solid #dcdfe6;
        }

        // 横幅样式
        &.banner-style {
          border-radius: 12px;
          background: #f8f9fa;
        }

        // 图标样式
        &.icon-style {
          border-radius: 6px;
          background: #fafafa;
        }

        :deep(.el-icon) {
          font-size: 48px;
          color: #c0c4cc;
        }
      }

      .image-overlay {
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
        border-radius: inherit; // 继承父元素的圆角
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

  .image-actions {
    display: flex;
    gap: 8px;
  }

  .image-tips {
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
  .image-upload {
    .image-actions {
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