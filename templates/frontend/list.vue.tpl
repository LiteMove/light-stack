<template>
  <div class="{{toLower .BusinessName}}-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">{{.FunctionName}}</h1>
        <p class="description">{{.TableComment}}</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd" v-if="checkPermission('{{generatePermission .ModuleName .BusinessName "add"}}')">
          <el-icon><Plus /></el-icon>
          新建{{.FunctionName}}
        </el-button>
      </div>
    </div>

    <!-- 搜索区域 -->
    {{- if .HasQuery }}
    <el-card class="search-card" shadow="never">
      <el-form
        :model="queryParams"
        ref="queryFormRef"
        :inline="true"
        label-width="auto"
        label-position="right"
      >
        {{- range .QueryFields }}
        <el-form-item label="{{.ColumnComment}}" prop="{{generateJSField .ColumnName}}">
          {{- if eq .QueryType "LIKE" }}
          <el-input
            v-model="queryParams.{{generateJSField .ColumnName}}"
            placeholder="请输入{{.ColumnComment}}"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
          {{- else if eq .QueryType "EQ" }}
          {{- if eq .HtmlType "select" }}
          <el-select
            v-model="queryParams.{{generateJSField .ColumnName}}"
            placeholder="请选择{{.ColumnComment}}"
            clearable
            style="width: 200px"
          >
            <!-- TODO: 根据字典类型{{.DictType}}添加选项 -->
            <el-option label="选项1" value="1" />
            <el-option label="选项2" value="2" />
          </el-select>
          {{- else }}
          <el-input
            v-model="queryParams.{{generateJSField .ColumnName}}"
            placeholder="请输入{{.ColumnComment}}"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
          {{- end }}
          {{- else if eq .QueryType "BETWEEN" }}
          {{- if eq .HtmlType "datetime" }}
          <el-date-picker
            v-model="queryParams.{{generateJSField .ColumnName}}Range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 300px"
          />
          {{- else if eq .HtmlType "date" }}
          <el-date-picker
            v-model="queryParams.{{generateJSField .ColumnName}}Range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 250px"
          />
          {{- else }}
          <el-input
            v-model="queryParams.{{generateJSField .ColumnName}}"
            placeholder="请输入{{.ColumnComment}}"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
          {{- end }}
          {{- end }}
        </el-form-item>
        {{- end }}
        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
          <el-button type="success" @click="handleExport" v-if="checkPermission('{{generatePermission .ModuleName .BusinessName "export"}}')">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
    {{- end }}

    <!-- 数据表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        :height="tableHeight"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column label="序号" type="index" width="60" align="center" />

        {{- range .ListFields }}
        {{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant")) }}
        <el-table-column
          prop="{{generateJSField .ColumnName}}"
          label="{{.ColumnComment}}"
          {{- if .IsPk }}
          width="100"
          {{- else if eq .HtmlType "datetime" }}
          width="180"
          {{- else if eq .HtmlType "date" }}
          width="120"
          {{- else if eq .HtmlType "switch" }}
          width="100"
          {{- else if contains .ColumnName "status" }}
          width="100"
          {{- else if contains .ColumnName "type" }}
          width="120"
          {{- end }}
          {{- if not (or .IsPk (eq .HtmlType "textarea") (contains .ColumnName "content") (contains .ColumnName "description")) }}
          sortable="custom"
          {{- end }}
          {{- if or (contains .ColumnName "content") (contains .ColumnName "description") (eq .HtmlType "textarea") }}
          show-overflow-tooltip
          {{- end }}
        >
          {{- if eq .HtmlType "switch" }}
          <template #default="{ row }">
            <el-switch
              v-model="row.{{generateJSField .ColumnName}}"
              :active-value="{{if eq .GoType "bool"}}true{{else}}1{{end}}"
              :inactive-value="{{if eq .GoType "bool"}}false{{else}}0{{end}}"
              @change="handleStatusChange(row)"
              :disabled="!checkPermission('{{generatePermission $.ModuleName $.BusinessName "edit"}}')"
            />
          </template>
          {{- else if eq .HtmlType "image" }}
          <template #default="{ row }">
            <el-image
              v-if="row.{{generateJSField .ColumnName}}"
              :src="row.{{generateJSField .ColumnName}}"
              :preview-src-list="[row.{{generateJSField .ColumnName}}]"
              fit="cover"
              style="width: 60px; height: 60px"
              preview-teleported
            />
            <span v-else class="text-gray-400">无图片</span>
          </template>
          {{- else if contains .ColumnName "status" }}
          <template #default="{ row }">
            <el-tag
              :type="getStatusTagType(row.{{generateJSField .ColumnName}})"
              size="small"
            >
              {{`{{ getStatusText(row.`}}{{generateJSField .ColumnName}}{{`) }}`}}
            </el-tag>
          </template>
          {{- else if eq .HtmlType "datetime" }}
          <template #default="{ row }">
            {{`{{ formatDateTime(row.`}}{{generateJSField .ColumnName}}{{`) }}`}}
          </template>
          {{- else if eq .HtmlType "date" }}
          <template #default="{ row }">
            {{`{{ formatDate(row.`}}{{generateJSField .ColumnName}}{{`) }}`}}
          </template>
          {{- else if or (contains .ColumnName "content") (contains .ColumnName "description") (eq .HtmlType "textarea") }}
          <template #default="{ row }">
            <span class="text-truncate">{{`{{ row.`}}{{generateJSField .ColumnName}}{{` || '-' }}`}}</span>
          </template>
          {{- end }}
        </el-table-column>
        {{- end }}
        {{- end }}

        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              text
              @click="handleView(row)"
              v-if="checkPermission('{{generatePermission .ModuleName .BusinessName "query"}}')"
            >
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button
              type="warning"
              size="small"
              text
              @click="handleEdit(row)"
              v-if="checkPermission('{{generatePermission .ModuleName .BusinessName "edit"}}')"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              text
              @click="handleDelete(row)"
              v-if="checkPermission('{{generatePermission .ModuleName .BusinessName "remove"}}')"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页组件 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          :small="false"
          :disabled="loading"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 表单弹窗 -->
    <{{.ClassName}}Form
      v-model:visible="formVisible"
      :form-data="formData"
      @success="handleFormSuccess"
    />

    <!-- 详情弹窗 -->
    <{{.ClassName}}Detail
      v-model:visible="detailVisible"
      :detail-data="detailData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Download,
  View,
  Edit,
  Delete
} from '@element-plus/icons-vue'
import { {{toLower .BusinessName}}Api } from '@/api'
import type { {{.ClassName}}, {{.ClassName}}QueryParams } from '@/api/types'
import {{.ClassName}}Form from './components/{{.ClassName}}Form.vue'
import {{.ClassName}}Detail from './components/{{.ClassName}}Detail.vue'
import { usePermission } from '@/composables/usePermission'
import { usePagination } from '@/composables/usePagination'
import { useTableHeight } from '@/composables/useTableHeight'

// 权限检查
const { checkPermission } = usePermission()

// 表格高度计算
const { tableHeight } = useTableHeight()

// 响应式数据
const queryFormRef = ref<FormInstance>()
const loading = ref(false)
const tableData = ref<{{.ClassName}}[]>([])
const selectedRows = ref<{{.ClassName}}[]>([])
const formVisible = ref(false)
const detailVisible = ref(false)
const formData = ref<Partial<{{.ClassName}}>>({})
const detailData = ref<{{.ClassName}} | null>(null)

// 查询参数
const queryParams = reactive<{{.ClassName}}QueryParams>({
  pageNum: 1,
  pageSize: 20,
{{- range .QueryFields }}
  {{generateJSField .ColumnName}}: {{if eq .QueryType "BETWEEN"}}{{if or (eq .HtmlType "datetime") (eq .HtmlType "date")}}null{{else}}''{{end}}{{else}}{{if eq .GoType "string"}}''{{else if eq .GoType "bool"}}null{{else if contains .GoType "int"}}null{{else}}null{{end}}{{end}},
{{- if eq .QueryType "BETWEEN" }}
{{- if or (eq .HtmlType "datetime") (eq .HtmlType "date") }}
  {{generateJSField .ColumnName}}Range: null,
{{- end }}
{{- end }}
{{- end }}
})

// 分页数据
const { total, handleSizeChange, handleCurrentChange } = usePagination(() => {
  queryParams.pageNum = 1
  getList()
})

// 计算属性
const hasSelected = computed(() => selectedRows.value.length > 0)

// 方法定义
const getList = async () => {
  loading.value = true
  try {
    // 处理时间范围查询参数
    const params = { ...queryParams }
    {{- range .QueryFields }}
    {{- if eq .QueryType "BETWEEN" }}
    {{- if or (eq .HtmlType "datetime") (eq .HtmlType "date") }}
    if (params.{{generateJSField .ColumnName}}Range) {
      params.{{generateJSField .ColumnName}}Start = params.{{generateJSField .ColumnName}}Range[0]
      params.{{generateJSField .ColumnName}}End = params.{{generateJSField .ColumnName}}Range[1]
    }
    delete params.{{generateJSField .ColumnName}}Range
    {{- end }}
    {{- end }}
    {{- end }}

    const response = await {{toLower .BusinessName}}Api.get{{.ClassName}}List(params)
    tableData.value = response.data.list
    total.value = response.data.total
  } catch (error) {
    console.error('获取{{.FunctionName}}列表失败:', error)
    ElMessage.error('获取{{.FunctionName}}列表失败')
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

const resetQuery = () => {
  queryFormRef.value?.resetFields()
  queryParams.pageNum = 1
  getList()
}

const handleAdd = () => {
  formData.value = {}
  formVisible.value = true
}

const handleEdit = (row: {{.ClassName}}) => {
  formData.value = { ...row }
  formVisible.value = true
}

const handleView = (row: {{.ClassName}}) => {
  detailData.value = row
  detailVisible.value = true
}

const handleDelete = async (row: {{.ClassName}}) => {
  const result = await ElMessageBox.confirm(
    `确定要删除"{{range .Fields}}{{if .IsPk}}${row.{{generateJSField .ColumnName}}}{{break}}{{end}}{{end}}"吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(() => false)

  if (!result) return

  try {
    await {{toLower .BusinessName}}Api.delete{{.ClassName}}(row.{{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}})
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    console.error('删除{{.FunctionName}}失败:', error)
    ElMessage.error('删除失败')
  }
}

const handleSelectionChange = (selection: {{.ClassName}}[]) => {
  selectedRows.value = selection
}

const handleSortChange = ({ prop, order }: { prop: string; order: string | null }) => {
  if (order) {
    queryParams.orderBy = prop
    queryParams.isAsc = order === 'ascending'
  } else {
    queryParams.orderBy = undefined
    queryParams.isAsc = undefined
  }
  getList()
}

{{- $hasStatusField := false }}
{{- range .ListFields }}
{{- if contains .ColumnName "status" }}
{{- $hasStatusField = true }}
const handleStatusChange = async (row: {{$.ClassName}}) => {
  try {
    await {{toLower $.BusinessName}}Api.update{{$.ClassName}}Status({
      {{range $.Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}: row.{{generateJSField .ColumnName}},{{break}}{{end}}{{end}}
      {{generateJSField .ColumnName}}: row.{{generateJSField .ColumnName}}
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    ElMessage.error('状态更新失败')
    // 还原状态
    row.{{generateJSField .ColumnName}} = !row.{{generateJSField .ColumnName}}
  }
}

const getStatusTagType = (status: any) => {
  switch (status) {
    case 1:
    case true:
      return 'success'
    case 0:
    case false:
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: any) => {
  switch (status) {
    case 1:
    case true:
      return '启用'
    case 0:
    case false:
      return '禁用'
    default:
      return '未知'
  }
}
{{- break }}
{{- end }}
{{- end }}

const handleExport = async () => {
  try {
    const response = await {{toLower .BusinessName}}Api.export{{.ClassName}}List(queryParams)
    // 处理文件下载
    const blob = new Blob([response.data], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `{{.FunctionName}}_${new Date().getTime()}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出{{.FunctionName}}失败:', error)
    ElMessage.error('导出失败')
  }
}

const handleFormSuccess = () => {
  getList()
}

// 工具函数
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString('zh-CN')
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  getList()
})
</script>

<style scoped lang="scss">
.{{toLower .BusinessName}}-page {
  padding: 20px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
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
  flex-shrink: 0;

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

.search-card {
  margin-bottom: 20px;
  flex-shrink: 0;

  :deep(.el-card__body) {
    padding: 20px 20px 8px 20px;
  }
}

.table-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  :deep(.el-card__body) {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 20px;
    overflow: hidden;
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  flex-shrink: 0;
}

.text-truncate {
  display: inline-block;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: top;
}

// 状态标签样式
:deep(.el-tag) {
  border-radius: 4px;
}

// 图片预览样式
:deep(.el-image) {
  border-radius: 4px;
  border: 1px solid #dcdfe6;
}

// 表格样式优化
:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;

  .el-table__header {
    th {
      background-color: #f5f7fa;
      color: #606266;
      font-weight: 600;
    }
  }

  .el-table__body {
    tr:hover {
      background-color: #f5f7fa;
    }
  }
}

// 搜索表单样式
:deep(.el-form--inline) {
  .el-form-item {
    margin-bottom: 12px;
  }
}
</style>