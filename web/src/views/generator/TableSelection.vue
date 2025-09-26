<template>
  <div class="table-selection">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">代码生成器 - 表选择</h1>
        <p class="description">选择需要生成CRUD代码的数据库表</p>
      </div>
      <div class="header-actions">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索表名或表注释"
          clearable
          style="width: 300px"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="refreshTables" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-container">
      <el-table
        v-loading="loading"
        :data="filteredTables"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        empty-text="暂无数据"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column prop="tableName" label="表名" min-width="200">
          <template #default="{ row }">
            <div class="table-name">
              <el-icon class="table-icon"><Collection /></el-icon>
              <span>{{ row.tableName }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="tableComment" label="表注释" min-width="250">
          <template #default="{ row }">
            <span>{{ row.tableComment || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="createTime" label="创建时间" width="180">
          <template #default="{ row }">
            <span>{{ formatDateTime(row.createTime) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="updateTime" label="更新时间" width="180">
          <template #default="{ row }">
            <span>{{ formatDateTime(row.updateTime) }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="handleGenerate(row)"
            >
              <el-icon><Tools /></el-icon>
              生成代码
            </el-button>

            <el-button
              type="info"
              size="small"
              @click="handlePreviewTable(row)"
            >
              <el-icon><View /></el-icon>
              预览
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 批量操作 -->
    <div class="batch-actions" v-show="selectedTables.length > 0">
      <div class="selected-info">
        已选择 <strong>{{ selectedTables.length }}</strong> 个表
      </div>
      <div class="actions">
        <el-button type="primary" @click="handleBatchGenerate">
          <el-icon><Tools /></el-icon>
          批量生成
        </el-button>
        <el-button @click="clearSelection">
          清除选择
        </el-button>
      </div>
    </div>

    <!-- 表结构预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      :title="`表结构预览 - ${previewTable?.tableName}`"
      width="80%"
      append-to-body
    >
      <div v-if="previewTable">
        <div class="table-info">
          <h3>表信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="表名">
              {{ previewTable.tableName }}
            </el-descriptions-item>
            <el-descriptions-item label="表注释">
              {{ previewTable.tableComment || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDateTime(previewTable.createTime) }}
            </el-descriptions-item>
            <el-descriptions-item label="更新时间">
              {{ formatDateTime(previewTable.updateTime) }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="columns-info" style="margin-top: 20px;">
          <h3>字段信息</h3>
          <el-table
            :data="previewTable.columns"
            style="width: 100%"
            max-height="400"
          >
            <el-table-column prop="columnName" label="字段名" width="150" />
            <el-table-column prop="columnType" label="字段类型" width="120" />
            <el-table-column prop="columnComment" label="注释" min-width="200" />
            <el-table-column prop="isNullable" label="允许空值" width="100">
              <template #default="{ row }">
                <el-tag :type="row.isNullable === 'YES' ? 'warning' : 'success'" size="small">
                  {{ row.isNullable === 'YES' ? '是' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="columnKey" label="键类型" width="100">
              <template #default="{ row }">
                <el-tag v-if="row.columnKey === 'PRI'" type="danger" size="small">主键</el-tag>
                <el-tag v-else-if="row.columnKey === 'UNI'" type="warning" size="small">唯一</el-tag>
                <el-tag v-else-if="row.columnKey === 'MUL'" type="info" size="small">索引</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="columnDefault" label="默认值" width="120">
              <template #default="{ row }">
                <span>{{ row.columnDefault || '-' }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <template #footer>
        <el-button @click="previewVisible = false">关闭</el-button>
        <el-button type="primary" @click="handleGenerateFromPreview">
          <el-icon><Tools /></el-icon>
          生成代码
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Tools, View, Collection } from '@element-plus/icons-vue'
import { getTableList, getTableColumns } from '@/api/generator'
import type { TableInfo } from '@/types/generator'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const searchKeyword = ref('')
const tables = ref<TableInfo[]>([])
const selectedTables = ref<TableInfo[]>([])
const previewVisible = ref(false)
const previewTable = ref<TableInfo | null>(null)

// 计算属性
const filteredTables = computed(() => {
  if (!searchKeyword.value) return tables.value

  const keyword = searchKeyword.value.toLowerCase()
  return tables.value.filter(table =>
    table.tableName.toLowerCase().includes(keyword) ||
    table.tableComment?.toLowerCase().includes(keyword)
  )
})

// 方法
const loadTables = async () => {
  loading.value = true
  try {
    const response = await getTableList()
    tables.value = response.data
  } catch (error) {
    ElMessage.error('获取数据库表列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const refreshTables = () => {
  loadTables()
}

const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
}

const handleSelectionChange = (selection: TableInfo[]) => {
  selectedTables.value = selection
}

const clearSelection = () => {
  selectedTables.value = []
}

const handleGenerate = (table: TableInfo) => {
  router.push({
    name: 'GeneratorConfig',
    params: { tableId: table.tableName }
  })
}

const handleBatchGenerate = async () => {
  if (selectedTables.value.length === 0) {
    ElMessage.warning('请选择要生成的表')
    return
  }

  const result = await ElMessageBox.confirm(
    `确定要为选中的 ${selectedTables.value.length} 个表生成代码吗？`,
    '批量生成确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(() => false)

  if (result) {
    // TODO: 实现批量生成逻辑
    ElMessage.success('批量生成任务已启动')
  }
}

const handlePreviewTable = async (table: TableInfo) => {
  loading.value = true
  try {
    const response = await getTableColumns(table.tableName)
    previewTable.value = {
      ...table,
      columns: response.data
    }
    previewVisible.value = true
  } catch (error) {
    ElMessage.error('获取表结构信息失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleGenerateFromPreview = () => {
  if (previewTable.value) {
    previewVisible.value = false
    handleGenerate(previewTable.value)
  }
}

const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadTables()
})
</script>

<style scoped lang="scss">
.table-selection {
  padding: 20px;
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

.table-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.table-name {
  display: flex;
  align-items: center;
  gap: 8px;

  .table-icon {
    color: #409eff;
  }
}

.batch-actions {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  z-index: 1000;

  .selected-info {
    color: #606266;
    font-size: 14px;
  }

  .actions {
    display: flex;
    gap: 12px;
  }
}

.table-info h3,
.columns-info h3 {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: #303133;
}

:deep(.el-table) {
  .cell {
    padding: 12px 0;
  }
}
</style>