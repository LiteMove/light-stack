<template>
  <div class="code-preview">
    <div class="code-editor-container">
      <div class="editor-header" v-if="showHeader">
        <div class="language-info">
          <el-tag size="small" :type="getLanguageTagType(language)">
            {{ getLanguageLabel(language) }}
          </el-tag>
          <span class="line-count">{{ lineCount }} 行</span>
        </div>
        <div class="editor-actions">
          <el-button
            size="small"
            @click="copyContent"
            :loading="copyLoading"
          >
            <el-icon><DocumentCopy /></el-icon>
            复制
          </el-button>
          <el-button
            size="small"
            @click="toggleWrap"
          >
            <el-icon><Sort /></el-icon>
            {{ wordWrap ? '取消换行' : '自动换行' }}
          </el-button>
        </div>
      </div>

      <div class="editor-content" :class="{ 'with-header': showHeader }">
        <pre class="code-block" :class="getLanguageClass(language)"><code v-html="highlightedContent"></code></pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { DocumentCopy, Sort } from '@element-plus/icons-vue'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

interface Props {
  content: string
  language?: string
  readonly?: boolean
  showHeader?: boolean
  height?: string | number
}

const props = withDefaults(defineProps<Props>(), {
  language: 'text',
  readonly: true,
  showHeader: true,
  height: 'auto'
})

// 响应式数据
const copyLoading = ref(false)
const wordWrap = ref(false)

// 计算属性
const lineCount = computed(() => {
  return props.content.split('\n').length
})

const highlightedContent = computed(() => {
  if (!props.content) return ''

  try {
    const language = getValidLanguage(props.language)

    if (language === 'text' || language === 'plain') {
      return escapeHtml(props.content)
    }

    const result = hljs.highlight(props.content, {
      language: language,
      ignoreIllegals: true
    })

    return result.value
  } catch (error) {
    // 如果高亮失败，返回转义的HTML
    return escapeHtml(props.content)
  }
})

// 方法
const getValidLanguage = (lang: string): string => {
  const supportedLanguages = [
    'javascript', 'typescript', 'go', 'java', 'python', 'sql',
    'html', 'css', 'scss', 'vue', 'json', 'xml', 'yaml',
    'markdown', 'bash', 'shell', 'dockerfile'
  ]

  const langMap: { [key: string]: string } = {
    'js': 'javascript',
    'ts': 'typescript',
    'vue': 'html', // Vue 文件使用 HTML 高亮
    'yml': 'yaml',
    'md': 'markdown',
    'sh': 'bash'
  }

  const mappedLang = langMap[lang.toLowerCase()] || lang.toLowerCase()

  return supportedLanguages.includes(mappedLang) ? mappedLang : 'text'
}

const getLanguageLabel = (lang: string): string => {
  const labels: { [key: string]: string } = {
    'javascript': 'JavaScript',
    'typescript': 'TypeScript',
    'go': 'Go',
    'java': 'Java',
    'python': 'Python',
    'sql': 'SQL',
    'html': 'HTML',
    'css': 'CSS',
    'scss': 'SCSS',
    'vue': 'Vue',
    'json': 'JSON',
    'xml': 'XML',
    'yaml': 'YAML',
    'markdown': 'Markdown',
    'bash': 'Bash',
    'shell': 'Shell',
    'dockerfile': 'Dockerfile',
    'text': 'Text',
    'plain': 'Plain Text'
  }

  return labels[lang.toLowerCase()] || lang.toUpperCase()
}

const getLanguageTagType = (lang: string): string => {
  const typeMap: { [key: string]: string } = {
    'go': 'primary',
    'javascript': 'warning',
    'typescript': 'info',
    'vue': 'success',
    'sql': 'danger',
    'json': 'info',
    'yaml': 'warning'
  }

  return typeMap[lang.toLowerCase()] || ''
}

const getLanguageClass = (lang: string): string => {
  return `language-${getValidLanguage(lang)}`
}

const escapeHtml = (text: string): string => {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

const copyContent = async () => {
  if (!props.content) {
    ElMessage.warning('没有内容可复制')
    return
  }

  copyLoading.value = true
  try {
    await navigator.clipboard.writeText(props.content)
    ElMessage.success('内容已复制到剪贴板')
  } catch (error) {
    // 降级方案：使用旧的API
    try {
      const textArea = document.createElement('textarea')
      textArea.value = props.content
      textArea.style.position = 'fixed'
      textArea.style.opacity = '0'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success('内容已复制到剪贴板')
    } catch (fallbackError) {
      ElMessage.error('复制失败，请手动复制')
    }
  } finally {
    copyLoading.value = false
  }
}

const toggleWrap = () => {
  wordWrap.value = !wordWrap.value
}

// 生命周期
onMounted(() => {
  nextTick(() => {
    // 确保highlight.js正确初始化
    hljs.configure({
      ignoreUnescapedHTML: true
    })
  })
})
</script>

<style scoped lang="scss">
.code-preview {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.code-editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  overflow: hidden;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #e4e7ed;

  .language-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .line-count {
      font-size: 12px;
      color: #909399;
    }
  }

  .editor-actions {
    display: flex;
    gap: 8px;
  }
}

.editor-content {
  flex: 1;
  overflow: auto;

  &.with-header {
    height: calc(100% - 48px);
  }

  .code-block {
    margin: 0;
    padding: 16px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 13px;
    line-height: 1.6;
    color: #333;
    background: #fff;
    white-space: pre;
    overflow-x: auto;
    min-height: 100%;
    box-sizing: border-box;

    code {
      font-family: inherit;
      font-size: inherit;
      line-height: inherit;
      background: transparent;
      padding: 0;
    }
  }

  .word-wrap {
    .code-block {
      white-space: pre-wrap;
      word-break: break-word;
    }
  }
}

// 自定义滚动条样式
.editor-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.editor-content::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.editor-content::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.editor-content::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

// 覆盖highlight.js默认样式
:deep(.hljs) {
  background: transparent !important;
  padding: 0 !important;
}

:deep(.hljs-comment),
:deep(.hljs-quote) {
  color: #6a737d;
  font-style: italic;
}

:deep(.hljs-keyword),
:deep(.hljs-selector-tag),
:deep(.hljs-subst) {
  color: #d73a49;
  font-weight: bold;
}

:deep(.hljs-number),
:deep(.hljs-literal),
:deep(.hljs-variable),
:deep(.hljs-template-variable),
:deep(.hljs-tag .hljs-attr) {
  color: #005cc5;
}

:deep(.hljs-string),
:deep(.hljs-doctag) {
  color: #032f62;
}

:deep(.hljs-title),
:deep(.hljs-section),
:deep(.hljs-selector-id) {
  color: #6f42c1;
  font-weight: bold;
}

:deep(.hljs-type),
:deep(.hljs-class .hljs-title) {
  color: #445588;
  font-weight: bold;
}

:deep(.hljs-tag),
:deep(.hljs-name),
:deep(.hljs-attribute) {
  color: #22863a;
}

:deep(.hljs-regexp),
:deep(.hljs-link) {
  color: #032f62;
}

:deep(.hljs-symbol),
:deep(.hljs-bullet) {
  color: #990073;
}

:deep(.hljs-built_in),
:deep(.hljs-builtin-name) {
  color: #e36209;
}
</style>