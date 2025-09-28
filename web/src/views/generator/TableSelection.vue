<template>
  <div class="table-selection">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><Collection /></el-icon>
          代码生成器
        </h2>
        <p class="page-desc">选择需要生成CRUD代码的数据库表</p>
      </div>
      <div class="header-actions">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索表名或表注释"
          clearable
          style="width: 300px"
          @input="handleSearch"
          :prefix-icon="Search"
        />
        <el-button
          :icon="RefreshRight"
          @click="refreshTables"
          :loading="loading"
          size="default"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-card class="table-card" shadow="never">
      <!-- 空状态提示 -->
      <div v-if="filteredTables.length === 0 && !loading" class="empty-state">
        <el-empty
          description="暂无数据库表"
          :image-size="120"
        >
          <template #image>
            <el-icon class="empty-icon"><Collection /></el-icon>
          </template>
        </el-empty>
      </div>

      <!-- 数据表格 -->
      <transition name="fade" mode="out-in" v-else>
        <el-table
          v-loading="loading"
          :data="filteredTables"
          stripe
          @selection-change="handleSelectionChange"
          style="width: 100%"
          element-loading-text="正在加载数据库表..."
          element-loading-background="rgba(0, 0, 0, 0.1)"
        >
          <el-table-column type="selection" width="50" />

          <el-table-column label="表信息" min-width="300" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="table-info">
                <div class="table-preview">
                  <div class="table-icon">
                    <el-icon color="#409eff"><Collection /></el-icon>
                  </div>
                </div>
                <div class="table-details">
                  <div class="table-name" :title="row.tableName">
                    {{ row.tableName }}
                  </div>
                  <div class="table-meta">
                    <span class="table-comment">{{ row.tableComment || '暂无注释' }}</span>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="创建时间" width="160" align="center" sortable>
            <template #default="{ row }">
              <div class="time-info">
                <div class="date">{{ formatDate(row.createTime) }}</div>
                <div class="time">{{ formatTimeOnly(row.createTime) }}</div>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="更新时间" width="160" align="center" sortable>
            <template #default="{ row }">
              <div class="time-info">
                <div class="date">{{ formatDate(row.updateTime) }}</div>
                <div class="time">{{ formatTimeOnly(row.updateTime) }}</div>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="220" fixed="right" align="center">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-tooltip content="预览表结构" placement="top">
                  <el-button
                    type="success"
                    size="small"
                    :icon="View"
                    @click="handlePreviewTable(row)"
                    circle
                  />
                </el-tooltip>
                <el-tooltip content="生成代码" placement="top">
                  <el-button
                    type="primary"
                    size="small"
                    :icon="Tools"
                    @click="handleGenerate(row)"
                    circle
                  />
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </transition>
    </el-card>

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
import {
  Search,
  RefreshRight,
  Tools,
  View,
  Collection,
  Document,
  Files
} from '@element-plus/icons-vue'
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

const formatDate = (timeStr: string): string => {
  return new Date(timeStr).toLocaleDateString('zh-CN')
}

const formatTimeOnly = (timeStr: string): string => {
  return new Date(timeStr).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 生命周期
onMounted(() => {
  loadTables()
})
</script>

<style scoped lang="scss">
.table-selection {
  padding: 24px;
  background: #f6f8fa;
  min-height: calc(100vh - 60px);

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
      }

      :deep(.el-input) {
        .el-input__wrapper {
          background: rgba(255, 255, 255, 0.1);
          border-color: rgba(255, 255, 255, 0.2);
          backdrop-filter: blur(10px);

          &:hover {
            border-color: rgba(255, 255, 255, 0.3);
          }

          .el-input__inner {
            color: white;

            &::placeholder {
              color: rgba(255, 255, 255, 0.7);
            }
          }

          .el-input__prefix {
            color: rgba(255, 255, 255, 0.8);
          }
        }
      }
    }
  }

  .stats-row {
    margin-bottom: 24px;

    .stats-card {
      border: 1px solid #d1d9e0;
      border-radius: 12px;
      transition: all 0.3s ease;

      &:hover {
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
        transform: translateY(-2px);
      }

      :deep(.el-card__body) {
        padding: 20px;
      }

      .stats-content {
        display: flex;
        align-items: center;

        .stats-icon {
          width: 48px;
          height: 48px;
          border-radius: 12px;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 16px;
          font-size: 24px;

          &.stats-icon-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
          }

          &.stats-icon-success {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
            color: white;
          }

          &.stats-icon-warning {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
          }

          &.stats-icon-info {
            background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
            color: white;
          }
        }

        .stats-info {
          .stats-number {
            font-size: 24px;
            font-weight: 700;
            color: #1f2328;
            line-height: 1.2;
          }

          .stats-label {
            font-size: 14px;
            color: #656d76;
            margin-top: 4px;
          }
        }
      }
    }
  }

  .table-card {
    border: 1px solid #d1d9e0;
    border-radius: 12px;

    :deep(.el-card__body) {
      padding: 0;
    }

    .el-table {
      .table-info {
        display: flex;
        align-items: center;
        gap: 12px;

        .table-preview {
          width: 48px;
          height: 48px;
          border-radius: 8px;
          overflow: hidden;
          border: 2px solid #f1f3f4;

          .table-icon {
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #f8f9fa;
            font-size: 20px;
          }
        }

        .table-details {
          flex: 1;
          min-width: 0;

          .table-name {
            font-weight: 600;
            color: #1f2328;
            margin-bottom: 6px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            cursor: pointer;

            &:hover {
              color: #0969da;
            }
          }

          .table-meta {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;
            color: #656d76;

            .table-comment {
              opacity: 0.8;
            }
          }
        }
      }

      .time-info {
        text-align: center;

        .date {
          font-weight: 500;
          color: #1f2328;
          font-size: 13px;
        }

        .time {
          font-size: 12px;
          color: #656d76;
          margin-top: 2px;
        }
      }

      .action-buttons {
        display: flex;
        gap: 8px;
        justify-content: center;
      }
    }

    .empty-state {
      padding: 60px 20px;
      text-align: center;

      .empty-icon {
        font-size: 120px;
        color: #dcdfe6;
        margin-bottom: 16px;
      }

      :deep(.el-empty__description) {
        color: #909399;
        font-size: 16px;
      }
    }
  }

  .batch-actions {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.16);
    padding: 16px 20px;
    display: flex;
    align-items: center;
    gap: 20px;
    z-index: 1000;
    backdrop-filter: blur(10px);
    border: 1px solid rgba(255, 255, 255, 0.2);

    .selected-info {
      color: #606266;
      font-size: 14px;
      font-weight: 500;
    }

    .actions {
      display: flex;
      gap: 12px;
    }
  }
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

  .cell {
    padding: 12px 0;
  }
}

// 过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 响应式设计
@media (max-width: 1200px) {
  .table-selection {
    .stats-row {
      .el-col {
        margin-bottom: 16px;
      }
    }
  }
}

@media (max-width: 768px) {
  .table-selection {
    padding: 16px;

    .page-header {
      flex-direction: column;
      gap: 16px;

      .header-actions {
        width: 100%;

        .el-button {
          flex: 1;
        }

        :deep(.el-input) {
          flex: 2;
        }
      }
    }

    .stats-row {
      .el-col {
        margin-bottom: 16px;
      }
    }

    .batch-actions {
      bottom: 16px;
      right: 16px;
      left: 16px;
      flex-direction: column;
      gap: 12px;

      .actions {
        justify-content: center;
        width: 100%;

        .el-button {
          flex: 1;
        }
      }
    }
  }
}
</style>