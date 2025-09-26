<template>
  <div class="preview-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">代码预览</h1>
        <p class="description">预览生成的代码文件</p>
      </div>
      <div class="header-actions">
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <el-button type="success" @click="copyAllCode" v-if="hasFiles">
          <el-icon><DocumentCopy /></el-icon>
          复制所有
        </el-button>
        <el-button type="primary" @click="handleGenerate" :loading="generateLoading">
          <el-icon><Tools /></el-icon>
          生成并下载
        </el-button>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="preview-container">
      <!-- 文件树 -->
      <div class="file-tree">
        <div class="tree-header">
          <h3>文件结构</h3>
          <el-button
            size="small"
            type="text"
            @click="expandAllNodes = !expandAllNodes"
          >
            {{ expandAllNodes ? '收起' : '展开' }}全部
          </el-button>
        </div>

        <el-tree
          ref="fileTreeRef"
          :data="fileTree"
          :props="treeProps"
          :default-expand-all="expandAllNodes"
          :highlight-current="true"
          @node-click="handleFileSelect"
          class="file-tree-content"
        >
          <template #default="{ node, data }">
            <div class="tree-node">
              <el-icon v-if="data.isDirectory" class="folder-icon">
                <Folder />
              </el-icon>
              <el-icon v-else class="file-icon">
                <Document />
              </el-icon>
              <span class="node-label">{{ node.label }}</span>
              <el-tag v-if="!data.isDirectory" size="small" type="info" class="file-type">
                {{ getFileExtension(data.name) }}
              </el-tag>
            </div>
          </template>
        </el-tree>
      </div>

      <!-- 代码内容 -->
      <div class="code-content">
        <div v-if="!currentFile" class="empty-content">
          <el-empty description="请从左侧选择文件进行预览" />
        </div>

        <div v-else class="file-content">
          <!-- 文件头部 -->
          <div class="file-header">
            <div class="file-info">
              <h4 class="file-name">{{ currentFile.name }}</h4>
              <span class="file-path">{{ currentFile.path }}</span>
            </div>
            <div class="file-actions">
              <el-button
                size="small"
                @click="copyFileContent"
              >
                <el-icon><DocumentCopy /></el-icon>
                复制
              </el-button>
              <el-button
                size="small"
                @click="downloadFile"
              >
                <el-icon><Download /></el-icon>
                下载
              </el-button>
            </div>
          </div>

          <!-- 代码编辑器 -->
          <div class="code-editor">
            <CodePreview
              :content="currentFile.content"
              :language="getLanguage(currentFile.name)"
              :readonly="true"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <el-loading
      v-loading="loading"
      text="正在生成预览..."
      background="rgba(0, 0, 0, 0.7)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  Tools,
  DocumentCopy,
  Download,
  Folder,
  Document
} from '@element-plus/icons-vue'
import { previewCode, generateCode } from '@/api/generator'
import CodePreview from '@/components/Generator/CodePreview.vue'

const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(false)
const generateLoading = ref(false)
const expandAllNodes = ref(false)
const fileTreeRef = ref()

interface FileNode {
  name: string
  path: string
  content: string
  isDirectory: boolean
  children?: FileNode[]
}

const fileTree = ref<FileNode[]>([])
const allFiles = ref<{ [path: string]: string }>({})
const currentFile = ref<FileNode | null>(null)

// 计算属性
const configId = computed(() => parseInt(route.params.configId as string))

const hasFiles = computed(() => Object.keys(allFiles.value).length > 0)

const treeProps = {
  label: 'name',
  children: 'children'
}

// 方法
const loadPreview = async () => {
  loading.value = true
  try {
    const response = await previewCode(configId.value)
    allFiles.value = response.data.files

    // 构建文件树
    buildFileTree()

    ElMessage.success('预览加载成功')
  } catch (error) {
    ElMessage.error('加载预览失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const buildFileTree = () => {
  const tree: FileNode[] = []
  const pathMap = new Map<string, FileNode>()

  // 处理所有文件路径
  Object.keys(allFiles.value).forEach(filePath => {
    const parts = filePath.split('/')
    let currentPath = ''
    let currentLevel = tree

    parts.forEach((part, index) => {
      const isLast = index === parts.length - 1
      currentPath = currentPath ? `${currentPath}/${part}` : part

      let node = pathMap.get(currentPath)
      if (!node) {
        node = {
          name: part,
          path: currentPath,
          content: isLast ? allFiles.value[filePath] : '',
          isDirectory: !isLast,
          children: isLast ? undefined : []
        }

        pathMap.set(currentPath, node)
        currentLevel.push(node)
      }

      if (!isLast && node.children) {
        currentLevel = node.children
      }
    })
  })

  fileTree.value = tree
}

const handleFileSelect = (data: FileNode) => {
  if (!data.isDirectory) {
    currentFile.value = data
  }
}

const getFileExtension = (fileName: string): string => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  const extMap: { [key: string]: string } = {
    'go': 'Go',
    'js': 'JS',
    'ts': 'TS',
    'vue': 'Vue',
    'sql': 'SQL',
    'json': 'JSON',
    'md': 'MD',
    'yml': 'YAML',
    'yaml': 'YAML'
  }
  return extMap[ext || ''] || ext?.toUpperCase() || ''
}

const getLanguage = (fileName: string): string => {
  const ext = fileName.split('.').pop()?.toLowerCase()
  const langMap: { [key: string]: string } = {
    'go': 'go',
    'js': 'javascript',
    'ts': 'typescript',
    'vue': 'vue',
    'sql': 'sql',
    'json': 'json',
    'md': 'markdown',
    'yml': 'yaml',
    'yaml': 'yaml',
    'html': 'html',
    'css': 'css',
    'scss': 'scss'
  }
  return langMap[ext || ''] || 'text'
}

const copyFileContent = async () => {
  if (!currentFile.value) return

  try {
    await navigator.clipboard.writeText(currentFile.value.content)
    ElMessage.success('文件内容已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const copyAllCode = async () => {
  try {
    const allContent = Object.entries(allFiles.value)
      .map(([path, content]) => `// ${path}\n${content}`)
      .join('\n\n' + '='.repeat(50) + '\n\n')

    await navigator.clipboard.writeText(allContent)
    ElMessage.success('所有代码已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const downloadFile = () => {
  if (!currentFile.value) return

  const blob = new Blob([currentFile.value.content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = currentFile.value.name
  a.click()
  URL.revokeObjectURL(url)

  ElMessage.success('文件下载成功')
}

const handleGenerate = async () => {
  const result = await ElMessageBox.confirm(
    '确定要生成并下载代码包吗？',
    '生成确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).catch(() => false)

  if (!result) return

  generateLoading.value = true
  try {
    const response = await generateCode({
      configId: configId.value,
      generateBackend: true,
      generateFrontend: true,
      generateSQL: true,
      outputFormat: 'zip',
      generateType: 'all'
    })

    ElMessage.success('代码生成成功！')

    // 跳转到历史记录页面
    router.push({ name: 'GeneratorHistory' })
  } catch (error) {
    ElMessage.error('生成代码失败')
    console.error(error)
  } finally {
    generateLoading.value = false
  }
}

const goBack = () => {
  router.back()
}

// 生命周期
onMounted(() => {
  loadPreview()
})
</script>

<style scoped lang="scss">
.preview-page {
  padding: 20px;
  height: calc(100vh - 40px);
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);

  .header-content {
    .title {
      margin: 0 0 8px 0;
      font-size: 24px;
      font-weight: 600;
      color: #303133;
    }

    .description {
      margin: 0;
      color: #909399;
      font-size: 14px;
    }
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }
}

.preview-container {
  flex: 1;
  display: flex;
  gap: 20px;
  min-height: 0;
}

.file-tree {
  width: 300px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;

  .tree-header {
    padding: 16px 20px;
    border-bottom: 1px solid #f0f0f0;
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 16px;
      color: #303133;
    }
  }

  .file-tree-content {
    flex: 1;
    padding: 16px;
    overflow: auto;
  }
}

.code-content {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  min-height: 0;

  .empty-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .file-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;

    .file-header {
      padding: 16px 20px;
      border-bottom: 1px solid #f0f0f0;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .file-info {
        .file-name {
          margin: 0 0 4px 0;
          font-size: 16px;
          color: #303133;
        }

        .file-path {
          font-size: 12px;
          color: #909399;
        }
      }

      .file-actions {
        display: flex;
        gap: 8px;
      }
    }

    .code-editor {
      flex: 1;
      min-height: 0;
    }
  }
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;

  .folder-icon {
    color: #409eff;
  }

  .file-icon {
    color: #67c23a;
  }

  .node-label {
    flex: 1;
    font-size: 14px;
  }

  .file-type {
    font-size: 11px;
  }
}

:deep(.el-tree) {
  .el-tree-node__content {
    height: 32px;
    line-height: 32px;
  }

  .el-tree-node__content:hover {
    background-color: #f5f7fa;
  }

  .is-current > .el-tree-node__content {
    background-color: #e6f7ff;
    color: #409eff;
  }
}
</style>