<template>
  <div class="{{toLower .BusinessName}}-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><{{.MenuIcon}} /></el-icon>
          {{.FunctionName}}管理
        </h2>
        <p class="page-desc">管理系统{{.FunctionName}}信息</p>
      </div>
      <div class="header-actions">
        <el-button
          v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "create"}}')"
          type="primary"
          :icon="Plus"
          @click="handleAdd"
          size="default"
        >
          新建{{.FunctionName}}
        </el-button>
        <el-button
          :icon="RefreshRight"
          @click="refreshData"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 操作工具栏 -->
    <el-card class="toolbar-card" shadow="never">
      <div class="toolbar-content">
        <!-- 搜索区域 -->
        <div class="search-section">
          <el-form :model="searchForm" inline class="search-form">
{{- if .HasQuery }}
{{- range .QueryFields }}
            <el-form-item label="{{.ColumnComment}}" class="search-item">
{{- if eq .QueryType "LIKE" }}
              <el-input
                v-model="searchForm.{{generateJSField .ColumnName}}"
                placeholder="搜索{{.ColumnComment}}"
                clearable
                @keyup.enter="handleSearch"
                @clear="handleSearch"
                :prefix-icon="Search"
                style="width: 200px"
              />
{{- else if eq .QueryType "EQ" }}
  {{- if eq .HtmlType "select" }}
              <el-select
                v-model="searchForm.{{generateJSField .ColumnName}}"
                placeholder="选择{{.ColumnComment}}"
                clearable
                @change="handleSearch"
                style="width: 150px"
              >
                <el-option label="全部" value="" />
                <!-- TODO: 根据字典类型添加选项 -->
              </el-select>
  {{- else }}
              <el-input
                v-model="searchForm.{{generateJSField .ColumnName}}"
                placeholder="{{.ColumnComment}}"
                clearable
                @keyup.enter="handleSearch"
                @clear="handleSearch"
                style="width: 150px"
              />
  {{- end }}
{{- else if eq .QueryType "BETWEEN" }}
              <el-date-picker
                v-model="searchForm.{{generateJSField .ColumnName}}Range"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                @change="handleSearch"
                style="width: 240px"
              />
{{- end }}
            </el-form-item>
{{- end }}
{{- else }}
            <el-form-item label="关键词" class="search-item">
              <el-input
                v-model="searchForm.keyword"
                placeholder="搜索关键词"
                clearable
                @keyup.enter="handleSearch"
                @clear="handleSearch"
                :prefix-icon="Search"
                style="width: 260px"
              />
            </el-form-item>
{{- end }}
            <el-form-item class="search-actions">
              <el-button
                type="primary"
                :icon="Search"
                @click="handleSearch"
                :loading="loading"
              >
                搜索
              </el-button>
              <el-button @click="handleResetSearch" :disabled="loading">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 批量操作区域 -->
        <div class="batch-section" v-show="selectedRows.length > 0">
          <div class="batch-info">
            <el-icon class="info-icon"><InfoFilled /></el-icon>
            已选择 <strong class="selected-count">{{`{{ selectedRows.length }}`}}</strong> 个{{.FunctionName}}
          </div>
          <div class="batch-actions">
            <el-button
              v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "delete"}}')"
              type="danger"
              size="small"
              :icon="Delete"
              @click="batchDelete"
            >
              批量删除
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- {{.FunctionName}}列表表格 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="table-header">
          <div class="table-title">
            <el-icon class="title-icon"><List /></el-icon>
            <span>{{.FunctionName}}列表</span>
            <el-tag type="info" size="small" class="total-count">
              共 {{`{{ pagination.total }}`}} 个{{.FunctionName}}
            </el-tag>
          </div>
          <div class="table-actions">
            <el-tooltip content="刷新数据" placement="top">
              <el-button
                size="small"
                :icon="RefreshRight"
                @click="refreshData"
                :loading="loading"
                circle
              />
            </el-tooltip>
          </div>
        </div>
      </template>

      <div>
        <el-table
          v-loading="loading"
          :data="dataList"
          @selection-change="handleSelectionChange"
          stripe
          border
          style="width: 100%"
          :header-row-style="{ backgroundColor: '#f8f9fa' }"
        >
        <el-table-column type="selection" width="50" align="center" />

{{- range .ListFields }}
        <!-- {{.ColumnComment}}列 -->
        <el-table-column prop="{{generateJSField .ColumnName}}" label="{{.ColumnComment}}"{{if gt (len .ColumnComment) 8}} min-width="{{mul (len .ColumnComment) 12}}"{{else}} width="120"{{end}} {{if ne .HtmlType "image"}}show-overflow-tooltip{{end}}>
          <template #default="{ row }">
{{- if eq .HtmlType "image" }}
            <el-image
              v-if="row.{{generateJSField .ColumnName}}"
              :src="row.{{generateJSField .ColumnName}}"
              :preview-src-list="[row.{{generateJSField .ColumnName}}]"
              fit="cover"
              style="width: 40px; height: 40px; border-radius: 4px;"
            />
            <span v-else class="empty-value">-</span>
{{- else if eq .HtmlType "select" }}
            <!-- TODO: 根据字典类型显示对应文本 -->
            <el-tag :type="getStatusType(row.{{generateJSField .ColumnName}})" size="small">
              {{`{{ getStatusText(row.`}}{{generateJSField .ColumnName}}{{`) }}`}}
            </el-tag>
{{- else if eq .GoType "bool" }}
            <el-tag :type="row.{{generateJSField .ColumnName}} ? 'success' : 'info'" size="small">
              {{`{{ row.`}}{{generateJSField .ColumnName}}{{` ? '是' : '否' }}`}}
            </el-tag>
{{- else if contains .ColumnName "time" }}
            <span class="time-text">{{`{{ formatDateTime(row.`}}{{generateJSField .ColumnName}}{{`) }}`}}</span>
{{- else }}
            {{`{{ row.`}}{{generateJSField .ColumnName}}{{` || '-' }}`}}
{{- end }}
          </template>
        </el-table-column>

{{- end }}
        <!-- 操作列 -->
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip
                v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "update"}}')"
                content="编辑{{.FunctionName}}"
                placement="top"
              >
                <el-button
                  type="primary"
                  link
                  size="small"
                  :icon="Edit"
                  @click="handleEdit(row)"
                />
              </el-tooltip>
              <el-tooltip
                v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "view"}}')"
                content="查看详情"
                placement="top"
              >
                <el-button
                  type="info"
                  link
                  size="small"
                  :icon="View"
                  @click="handleView(row)"
                />
              </el-tooltip>
              <el-tooltip
                v-if="$hasPer('{{generatePermission .ModuleName .BusinessName "delete"}}')"
                content="删除{{.FunctionName}}"
                placement="top"
              >
                <el-button
                  type="danger"
                  link
                  size="small"
                  :icon="Delete"
                  @click="handleDelete(row)"
                />
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页器 -->
      <div class="pagination-wrapper" v-if="pagination.total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
      </div>
    </el-card>

    <!-- {{.FunctionName}}表单弹窗 -->
    <{{.ClassName}}Form
      v-model:visible="formVisible"
      :form-data="formData"
      @success="handleFormSuccess"
    />

    <!-- {{.FunctionName}}详情弹窗 -->
    <{{.ClassName}}Detail
      v-model:visible="detailVisible"
      :data="detailData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  RefreshRight,
  Search,
  Edit,
  Delete,
  View,
  List,
  InfoFilled,
  Download,
  {{.MenuIcon}}
} from '@element-plus/icons-vue'
import { {{toLower .BusinessName}}Api } from '@/api'
import type { {{.ClassName}}, PageParams } from '@/api/types'
import {{.ClassName}}Form from './components/{{.ClassName}}Form.vue'
import {{.ClassName}}Detail from './components/{{.ClassName}}Detail.vue'
import { formatDateTime } from '@/utils/date'
import { useUserStore } from '@/store'

// Store实例
const userStore = useUserStore()

// 响应式数据
const loading = ref(false)
const dataList = ref<{{.ClassName}}[]>([])
const selectedRows = ref<{{.ClassName}}[]>([])
const formVisible = ref(false)
const detailVisible = ref(false)
const formData = ref<Partial<{{.ClassName}}>>({})
const detailData = ref<{{.ClassName}} | null>(null)
const abortController = ref<AbortController | null>(null)
const isUnmounting = ref(false)

// 搜索表单
const searchForm = reactive({
{{- if .HasQuery }}
{{- range .QueryFields }}
{{- if eq .QueryType "BETWEEN" }}
  {{generateJSField .ColumnName}}Range: [] as string[],
{{- else }}
  {{generateJSField .ColumnName}}: '{{getDefaultValue .}}',
{{- end }}
{{- end }}
{{- else }}
  keyword: '',
{{- end }}
  page: 1,
  pageSize: 20
})

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 获取{{.FunctionName}}列表
const fetchData = async () => {
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过获取{{.FunctionName}}列表')
    return
  }

  if (abortController.value) {
    abortController.value.abort()
  }

  abortController.value = new AbortController()

  try {
    loading.value = true
    const params: PageParams & Record<string, any> = {
      page: pagination.page,
      page_size: pagination.pageSize,
{{- if .HasQuery }}
{{- range .QueryFields }}
{{- if eq .QueryType "BETWEEN" }}
      {{generateJSField .ColumnName}}_start: searchForm.{{generateJSField .ColumnName}}Range[0] || undefined,
      {{generateJSField .ColumnName}}_end: searchForm.{{generateJSField .ColumnName}}Range[1] || undefined,
{{- else if eq .QueryType "LIKE" }}
      {{generateJSField .ColumnName}}: searchForm.{{generateJSField .ColumnName}} || undefined,
{{- else }}
      {{generateJSField .ColumnName}}: searchForm.{{generateJSField .ColumnName}} === '{{getDefaultValue .}}' ? undefined : searchForm.{{generateJSField .ColumnName}},
{{- end }}
{{- end }}
{{- else }}
      keyword: searchForm.keyword || undefined,
{{- end }}
    }

    const { data } = await {{toLower .BusinessName}}Api.get{{pluralize .ClassName}}(params)

    if (abortController.value?.signal.aborted || isUnmounting.value) {
      return
    }

    dataList.value = data.list
    pagination.total = data.total

  } catch (error: any) {
    if (error.name === 'AbortError' || error.message?.includes('登录') || isUnmounting.value) {
      console.log('请求已取消或用户未登录')
      return
    }
    console.error('获取{{.FunctionName}}列表失败:', error)
  } finally {
    if (!isUnmounting.value) {
      loading.value = false
    }
  }
}

// 刷新{{.FunctionName}}列表
const refreshData = () => {
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过刷新{{.FunctionName}}列表')
    return
  }
  fetchData()
}

// 搜索
const handleSearch = () => {
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过搜索')
    return
  }
  pagination.page = 1
  fetchData()
}

// 重置搜索
const handleResetSearch = () => {
  if (isUnmounting.value) {
    console.log('组件销毁中，跳过重置搜索')
    return
  }
  Object.assign(searchForm, {
{{- if .HasQuery }}
{{- range .QueryFields }}
{{- if eq .QueryType "BETWEEN" }}
    {{generateJSField .ColumnName}}Range: [],
{{- else }}
    {{generateJSField .ColumnName}}: '{{getDefaultValue .}}',
{{- end }}
{{- end }}
{{- else }}
    keyword: '',
{{- end }}
    page: 1,
    pageSize: 20
  })
  handleSearch()
}

// 添加{{.FunctionName}}
const handleAdd = () => {
  formData.value = {}
  formVisible.value = true
}

// 编辑{{.FunctionName}}
const handleEdit = (row: {{.ClassName}}) => {
  formData.value = { ...row }
  formVisible.value = true
}

// 查看详情
const handleView = (row: {{.ClassName}}) => {
  detailData.value = row
  detailVisible.value = true
}

// 删除{{.FunctionName}}
const handleDelete = async (row: {{.ClassName}}) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除{{.FunctionName}}吗？此操作不可恢复！`,
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: false
      }
    )

    await {{toLower .BusinessName}}Api.delete{{.ClassName}}(row.{{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}})
    ElMessage.success('删除成功')
    refreshData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除{{.FunctionName}}失败:', error)
    }
  }
}

// 选择变化
const handleSelectionChange = (selection: {{.ClassName}}[]) => {
  selectedRows.value = selection
}

// 批量删除
const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 \${selectedRows.value.length} 个{{.FunctionName}}吗？此操作不可恢复！`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const promises = selectedRows.value.map(row => {{toLower .BusinessName}}Api.delete{{.ClassName}}(row.{{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}}))
    await Promise.all(promises)
    ElMessage.success('批量删除成功')
    refreshData()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
    }
  }
}

// 分页相关
const handlePageSizeChange = (size: number) => {
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过分页大小变更')
    return
  }
  pagination.pageSize = size
  pagination.page = 1
  fetchData()
}

const handleCurrentChange = (page: number) => {
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过页码变更')
    return
  }
  pagination.page = page
  fetchData()
}

// 表单成功回调
const handleFormSuccess = () => {
  if (isUnmounting.value || !userStore.getToken()) {
    console.log('组件销毁中或用户未登录，跳过表单成功回调')
    return
  }
  refreshData()
}

// 状态/类型转换函数
{{- if .HasQuery }}
{{- range .QueryFields }}
{{- if eq .HtmlType "select" }}
const getStatusType = (value: any): string => {
  // TODO: 根据具体业务逻辑返回对应的标签类型
  return 'primary'
}

const getStatusText = (value: any): string => {
  // TODO: 根据字典类型返回对应的文本
  return String(value)
}
{{- break }}
{{- end }}
{{- end }}
{{- end }}

// 初始化
onMounted(() => {
  refreshData()
})

// 组件销毁时停止监听并清理资源
onUnmounted(() => {
  console.log('{{.FunctionName}}管理组件开始销毁，清理资源')
  isUnmounting.value = true

  if (abortController.value) {
    abortController.value.abort()
    abortController.value = null
  }

  dataList.value = []
  selectedRows.value = []
  loading.value = false

  console.log('{{.FunctionName}}管理组件资源清理完成')
})
</script>

<style lang="scss" scoped>
.{{toLower .BusinessName}}-management {
  padding: 24px;
  background-color: #f5f6fa;
  min-height: calc(100vh - 64px);

  // 页面头部样式
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
    padding: 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);

    .header-left {
      color: white;

      .page-title {
        margin: 0 0 8px 0;
        font-size: 28px;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 12px;
        flex-wrap: wrap;

        .title-icon {
          font-size: 32px;
        }
      }

      .page-desc {
        margin: 0;
        font-size: 14px;
        opacity: 0.9;
        line-height: 1.4;
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;

      .el-button {
        background: rgba(255, 255, 255, 0.1);
        border-color: rgba(255, 255, 255, 0.2);
        color: white;
        backdrop-filter: blur(10px);

        &:hover {
          background: rgba(255, 255, 255, 0.2);
          border-color: rgba(255, 255, 255, 0.3);
          transform: translateY(-2px);
        }

        &.el-button--primary {
          background: rgba(255, 255, 255, 0.2);

          &:hover {
            background: rgba(255, 255, 255, 0.3);
          }
        }
      }
    }
  }

  // 工具栏样式
  .toolbar-card {
    margin-bottom: 20px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    border: 1px solid #e4e7ed;

    :deep(.el-card__body) {
      padding: 20px;
    }

    .toolbar-content {
      .search-section {
        margin-bottom: 16px;

        .search-form {
          display: flex;
          flex-wrap: wrap;
          align-items: flex-end;
          gap: 16px;

          .search-item {
            margin-bottom: 0;
            margin-right: 0;

            :deep(.el-form-item__label) {
              color: #606266;
              font-weight: 500;
              font-size: 13px;
            }
          }

          .search-actions {
            margin-bottom: 0;
            margin-right: 0;
            display: flex;
            gap: 12px;

            .el-button {
              border-radius: 6px;
              transition: all 0.2s ease;

              &:hover {
                transform: translateY(-1px);
              }
            }
          }
        }
      }

      .batch-section {
        padding: 16px;
        background: linear-gradient(90deg, #e3f2fd 0%, #f3e5f5 100%);
        border-radius: 8px;
        border: 1px solid #bbdefb;
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-top: 16px;

        .batch-info {
          display: flex;
          align-items: center;
          gap: 8px;
          color: #1976d2;
          font-weight: 500;

          .info-icon {
            font-size: 18px;
            color: #2196f3;
          }

          .selected-count {
            color: #1565c0;
            font-size: 16px;
          }
        }

        .batch-actions {
          display: flex;
          gap: 8px;
        }
      }
    }
  }

  // 表格卡片样式
  .table-card {
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    border: 1px solid #e4e7ed;
    overflow: hidden;

    :deep(.el-card__body) {
      padding: 0;
    }

    :deep(.el-card__header) {
      padding: 16px 20px;
      border-bottom: 1px solid #f0f0f0;
      background: #fafafa;
    }

    .table-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .table-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-weight: 600;
        color: #303133;

        .title-icon {
          color: #409eff;
          font-size: 16px;
        }

        .total-count {
          margin-left: 8px;
          background: #e3f2fd;
          color: #1976d2;
          border: 1px solid #bbdefb;
        }
      }

      .table-actions {
        display: flex;
        gap: 8px;

        .el-button {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          transition: all 0.2s ease;

          &:hover {
            transform: scale(1.1);
          }
        }
      }
    }

    // 表格样式优化
    .el-table {
      border: none;

      // 表头样式
      :deep(.el-table__header-wrapper) {
        .el-table__header {
          th {
            background: #f8f9fa;
            color: #495057;
            font-weight: 600;
            font-size: 14px;
            border-bottom: 2px solid #dee2e6;
          }
        }
      }

      // 行样式
      :deep(.el-table__body) {
        tr {
          transition: all 0.2s ease;

          &:hover {
            background-color: #f8f9ff;
          }
        }
      }

      .empty-value {
        color: #c0c4cc;
        font-style: italic;
        font-size: 12px;
      }

      .time-text {
        font-size: 12px;
        color: #909399;
      }

      // 操作按钮
      .action-buttons {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 4px;

        .el-button {
          padding: 4px;
          width: 28px;
          height: 28px;
          border-radius: 6px;

          &:hover {
            transform: translateY(-1px);
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          }
        }
      }
    }

    // 分页样式
    .pagination-wrapper {
      padding: 20px;
      background: #fafafa;
      border-top: 1px solid #e4e7ed;
      display: flex;
      justify-content: center;

      :deep(.el-pagination) {
        .el-pager li {
          border-radius: 6px;
          margin: 0 2px;
        }

        .btn-prev,
        .btn-next {
          border-radius: 6px;
        }
      }
    }
  }
}

// 响应式设计优化
@media (max-width: 1200px) {
  .{{toLower .BusinessName}}-management {
    padding: 16px;

    .page-header {
      flex-direction: column;
      align-items: stretch;
      gap: 16px;

      .header-actions {
        justify-content: flex-end;
      }
    }
  }
}

@media (max-width: 768px) {
  .{{toLower .BusinessName}}-management {
    padding: 12px;

    .page-header {
      padding: 20px 16px;

      .header-left {
        .page-title {
          font-size: 24px;
        }
      }

      .header-actions {
        flex-wrap: wrap;
        gap: 8px;

        .el-button {
          flex: 1;
          min-width: 100px;
        }
      }
    }

    .toolbar-card {
      :deep(.el-card__body) {
        padding: 16px;
      }

      .toolbar-content {
        .search-section {
          .search-form {
            flex-direction: column;
            align-items: stretch;

            .search-item {
              width: 100%;
              margin-bottom: 12px;
              margin-right: 0;

              :deep(.el-input),
              :deep(.el-select) {
                width: 100% !important;
              }
            }

            .search-actions {
              justify-content: center;

              .el-button {
                flex: 1;
              }
            }
          }
        }

        .batch-section {
          flex-direction: column;
          align-items: stretch;
          gap: 12px;

          .batch-actions {
            justify-content: center;
          }
        }
      }
    }

    // 移动端表格滚动
    .table-card {
      .el-table {
        min-width: 1000px;
      }
    }
  }
}
</style>