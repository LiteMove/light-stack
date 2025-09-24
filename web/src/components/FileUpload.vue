<template>
  <div class="file-upload">
    <!-- 拖拽上传区域 -->
    <div
      class="upload-dragger"
      :class="{ 'is-dragover': isDragOver, 'is-uploading': uploading }"
      @drop="handleDrop"
      @dragover="handleDragOver"
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @click="handleClick"
    >
      <div class="upload-content">
        <div v-if="uploading" class="upload-status">
          <el-icon class="upload-icon spinning"><Loading /></el-icon>
          <div class="upload-text">上传中...</div>
          <div class="upload-progress">
            <el-progress :percentage="uploadProgress" :show-text="false" />
          </div>
        </div>
        <div v-else class="upload-prompt">
          <el-icon class="upload-icon"><UploadFilled /></el-icon>
          <div class="upload-text">将文件拖拽到此处，或<span class="upload-link">点击上传</span></div>
          <div class="upload-hint">
            支持 {{ allowedExts.join(', ') }} 格式，文件大小不超过 {{ maxSizeMB }}MB
          </div>
        </div>
      </div>
    </div>

    <!-- 隐藏的文件输入框 -->
    <input
      ref="fileInputRef"
      type="file"
      :accept="acceptAttribute"
      style="display: none"
      @change="handleFileChange"
    />

    <!-- 上传选项 -->
    <div v-if="showOptions && !uploading && !uploadedFile" class="upload-options">
      <el-checkbox v-model="isPublic" label="公开文件" />
      <div class="option-hint">
        公开文件可被所有用户访问，私有文件仅上传者可访问
      </div>
    </div>

    <!-- 已上传文件信息 -->
    <div v-if="uploadedFile" class="uploaded-file">
      <div class="file-info">
        <el-icon class="file-icon"><Document /></el-icon>
        <div class="file-details">
          <div class="file-name">{{ uploadedFile.originalName }}</div>
          <div class="file-meta">
            {{ formatFileSize(uploadedFile.fileSize) }} • {{ uploadedFile.fileType }}
          </div>
        </div>
        <div class="file-actions">
          <el-button
            type="primary"
            size="small"
            @click="downloadFile"
            :loading="downloading"
          >
            下载
          </el-button>
          <el-button
            type="danger"
            size="small"
            @click="removeFile"
            :loading="removing"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled, Loading, Document } from '@element-plus/icons-vue'
import { uploadFile, deleteFile } from '@/api/file'

interface FileProfile {
  id: number
  originalName: string
  fileName: string
  fileSize: number
  fileType: string
  mimeType: string
  md5: string
  uploadUserId: number
  usageType: string
  accessUrl: string
  createdAt: string
  updatedAt: string
}

interface Props {
  usageType?: string
  maxSize?: number // MB
  allowedTypes?: string[]
  modelValue?: FileProfile | null
  showOptions?: boolean // 是否显示上传选项
}

interface Emits {
  (e: 'update:modelValue', value: FileProfile | null): void
  (e: 'success', file: FileProfile): void
  (e: 'error', error: string): void
}

const props = withDefaults(defineProps<Props>(), {
  usageType: '',
  maxSize: 50,
  allowedTypes: () => ['jpg', 'jpeg', 'png', 'gif', 'pdf', 'doc', 'docx', 'xls', 'xlsx', 'txt'],
  showOptions: false
})

const emit = defineEmits<Emits>()

// 响应式数据
const isDragOver = ref(false)
const uploading = ref(false)
const downloading = ref(false)
const removing = ref(false)
const uploadProgress = ref(0)
const fileInputRef = ref<HTMLInputElement>()
const isPublic = ref(false)

// 计算属性
const allowedExts = computed(() =>
  props.allowedTypes.map(type => `.${type}`)
)

const acceptAttribute = computed(() =>
  allowedExts.value.join(',')
)

const maxSizeMB = computed(() => props.maxSize)

const uploadedFile = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 文件大小格式化
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 验证文件
const validateFile = (file: File): boolean => {
  // 检查文件大小
  const maxSizeBytes = props.maxSize * 1024 * 1024
  if (file.size > maxSizeBytes) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
    return false
  }

  // 检查文件类型
  const fileExt = file.name.split('.').pop()?.toLowerCase()
  if (!fileExt || !props.allowedTypes.includes(fileExt)) {
    ElMessage.error(`不支持的文件类型，仅支持: ${props.allowedTypes.join(', ')}`)
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
    if (props.usageType) {
      formData.append('usageType', props.usageType)
    }
    formData.append('isPublic', isPublic.value.toString())

    const response = await uploadFile(formData)

    clearInterval(progressInterval)
    uploadProgress.value = 100

    // 延迟一下让用户看到100%
    await new Promise(resolve => setTimeout(resolve, 500))

    uploadedFile.value = response.data
    emit('success', response.data)
    ElMessage.success('文件上传成功')
  } catch (error: any) {
    ElMessage.error(error.message || '文件上传失败')
    emit('error', error.message || '文件上传失败')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

// 拖拽事件处理
const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
  isDragOver.value = true
}

const handleDragEnter = (e: DragEvent) => {
  e.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (e: DragEvent) => {
  e.preventDefault()
  // 只有当离开整个拖拽区域时才设置为false
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  if (
    e.clientX < rect.left ||
    e.clientX > rect.right ||
    e.clientY < rect.top ||
    e.clientY > rect.bottom
  ) {
    isDragOver.value = false
  }
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragOver.value = false

  if (uploading.value) return

  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    uploadFileHandler(files[0])
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

// 下载文件
const downloadFile = async () => {
  if (!uploadedFile.value) return

  try {
    downloading.value = true
    // 使用文件的访问URL进行下载
    if (uploadedFile.value.accessUrl) {
      const link = document.createElement('a')
      link.href = uploadedFile.value.accessUrl
      link.download = uploadedFile.value.originalName
      link.style.display = 'none'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    } else {
      ElMessage.error('无法获取文件下载链接')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '下载失败')
  } finally {
    downloading.value = false
  }
}

// 删除文件
const removeFile = async () => {
  if (!uploadedFile.value) return

  try {
    removing.value = true
    await deleteFile(uploadedFile.value.id)
    uploadedFile.value = null
    ElMessage.success('文件删除成功')
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
  } finally {
    removing.value = false
  }
}
</script>

<style lang="scss" scoped>
.file-upload {
  width: 100%;

  .upload-dragger {
    width: 100%;
    height: 180px;
    border: 2px dashed #dcdfe6;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: all 0.3s;
    background-color: #fafafa;

    &:hover {
      border-color: #409eff;
      background-color: #ecf5ff;
    }

    &.is-dragover {
      border-color: #409eff;
      background-color: #ecf5ff;
      transform: scale(1.02);
    }

    &.is-uploading {
      border-color: #409eff;
      background-color: #ecf5ff;
      cursor: not-allowed;
    }

    .upload-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      padding: 20px;

      .upload-icon {
        font-size: 48px;
        color: #c0c4cc;
        margin-bottom: 16px;
        transition: color 0.3s;

        &.spinning {
          animation: rotate 2s linear infinite;
          color: #409eff;
        }
      }

      .upload-text {
        font-size: 16px;
        color: #606266;
        margin-bottom: 8px;

        .upload-link {
          color: #409eff;
          font-weight: 500;
        }
      }

      .upload-hint {
        font-size: 12px;
        color: #909399;
        text-align: center;
        line-height: 1.4;
      }

      .upload-progress {
        width: 200px;
        margin-top: 16px;
      }
    }
  }

  .upload-options {
    margin-top: 16px;
    padding: 12px;
    background-color: #f8f9fa;
    border-radius: 6px;
    border: 1px solid #e9ecef;

    .option-hint {
      font-size: 12px;
      color: #6c757d;
      margin-top: 4px;
      line-height: 1.4;
    }
  }

  .uploaded-file {
    margin-top: 16px;
    padding: 16px;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    background-color: #f9f9f9;

    .file-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .file-icon {
        font-size: 24px;
        color: #409eff;
        flex-shrink: 0;
      }

      .file-details {
        flex: 1;
        min-width: 0;

        .file-name {
          font-size: 14px;
          font-weight: 500;
          color: #303133;
          margin-bottom: 4px;
          word-break: break-all;
        }

        .file-meta {
          font-size: 12px;
          color: #909399;
        }
      }

      .file-actions {
        display: flex;
        gap: 8px;
        flex-shrink: 0;
      }
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
  .file-upload {
    .upload-dragger {
      height: 120px;

      .upload-content {
        padding: 16px;

        .upload-icon {
          font-size: 36px;
          margin-bottom: 12px;
        }

        .upload-text {
          font-size: 14px;
          text-align: center;
        }

        .upload-hint {
          font-size: 11px;
        }

        .upload-progress {
          width: 150px;
          margin-top: 12px;
        }
      }
    }

    .uploaded-file {
      .file-info {
        .file-actions {
          flex-direction: column;

          .el-button {
            width: 60px;
            font-size: 12px;
          }
        }
      }
    }
  }
}
</style>