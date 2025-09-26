<template>
  <div class="field-config-table">
    <div class="table-header">
      <div class="header-info">
        <span>共 {{ modelValue.length }} 个字段</span>
        <el-divider direction="vertical" />
        <el-button
          size="small"
          type="text"
          @click="selectAllFields(true)"
        >
          全选
        </el-button>
        <el-button
          size="small"
          type="text"
          @click="selectAllFields(false)"
        >
          全不选
        </el-button>
      </div>
      <div class="header-actions">
        <el-button
          size="small"
          @click="expandAll = !expandAll"
        >
          {{ expandAll ? '收起' : '展开' }}配置
        </el-button>
      </div>
    </div>

    <el-table
      :data="modelValue"
      v-loading="loading"
      border
      size="small"
      style="width: 100%"
      max-height="600"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />

      <el-table-column type="expand" width="50" v-if="expandAll">
        <template #default="{ row }">
          <div class="field-expand-config">
            <el-form :model="row" label-width="100px" size="small">
              <el-row :gutter="16">
                <el-col :span="8">
                  <el-form-item label="查询类型">
                    <el-select v-model="row.queryType" style="width: 100%">
                      <el-option label="等于" value="EQ" />
                      <el-option label="不等于" value="NE" />
                      <el-option label="大于" value="GT" />
                      <el-option label="小于" value="LT" />
                      <el-option label="模糊" value="LIKE" />
                      <el-option label="范围" value="BETWEEN" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="显示类型">
                    <el-select v-model="row.htmlType" style="width: 100%">
                      <el-option label="文本框" value="input" />
                      <el-option label="文本域" value="textarea" />
                      <el-option label="下拉框" value="select" />
                      <el-option label="单选框" value="radio" />
                      <el-option label="复选框" value="checkbox" />
                      <el-option label="日期时间" value="datetime" />
                      <el-option label="文件上传" value="upload" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="字典类型">
                    <el-input
                      v-model="row.dictType"
                      placeholder="如：sys_status"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-form>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="columnName" label="字段名" width="150" />

      <el-table-column prop="columnComment" label="字段注释" min-width="150">
        <template #default="{ row }">
          <el-input
            v-model="row.columnComment"
            size="small"
            placeholder="请输入注释"
          />
        </template>
      </el-table-column>

      <el-table-column prop="columnType" label="物理类型" width="120" />

      <el-table-column prop="goType" label="Go类型" width="100">
        <template #default="{ row }">
          <el-tag size="small" :type="getGoTypeTagType(row.goType)">
            {{ row.goType }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="主键" width="60">
        <template #default="{ row }">
          <el-checkbox
            v-model="row.isPk"
            :disabled="true"
          />
        </template>
      </el-table-column>

      <el-table-column label="必填" width="60">
        <template #default="{ row }">
          <el-checkbox
            v-model="row.isRequired"
            :disabled="row.isPk"
          />
        </template>
      </el-table-column>

      <el-table-column label="插入" width="60">
        <template #default="{ row }">
          <el-checkbox
            v-model="row.isInsert"
            :disabled="row.isPk && row.isIncrement"
          />
        </template>
      </el-table-column>

      <el-table-column label="编辑" width="60">
        <template #default="{ row }">
          <el-checkbox
            v-model="row.isEdit"
            :disabled="row.isPk"
          />
        </template>
      </el-table-column>

      <el-table-column label="列表" width="60">
        <template #default="{ row }">
          <el-checkbox v-model="row.isList" />
        </template>
      </el-table-column>

      <el-table-column label="查询" width="60">
        <template #default="{ row }">
          <el-checkbox
            v-model="row.isQuery"
            :disabled="row.isPk"
          />
        </template>
      </el-table-column>

      <el-table-column label="查询方式" width="100">
        <template #default="{ row }">
          <el-select
            v-model="row.queryType"
            size="small"
            :disabled="!row.isQuery"
            style="width: 100%"
          >
            <el-option label="=" value="EQ" />
            <el-option label="!=" value="NE" />
            <el-option label=">" value="GT" />
            <el-option label="<" value="LT" />
            <el-option label="LIKE" value="LIKE" />
            <el-option label="范围" value="BETWEEN" />
          </el-select>
        </template>
      </el-table-column>

      <el-table-column label="显示类型" width="120">
        <template #default="{ row }">
          <el-select
            v-model="row.htmlType"
            size="small"
            style="width: 100%"
          >
            <el-option label="文本框" value="input" />
            <el-option label="文本域" value="textarea" />
            <el-option label="下拉框" value="select" />
            <el-option label="单选框" value="radio" />
            <el-option label="复选框" value="checkbox" />
            <el-option label="日期时间" value="datetime" />
            <el-option label="上传" value="upload" />
          </el-select>
        </template>
      </el-table-column>
    </el-table>

    <div class="table-footer">
      <div class="footer-info">
        <el-text type="info" size="small">
          提示：主键字段自动排除插入和编辑，自增字段自动排除插入
        </el-text>
      </div>
      <div class="footer-actions">
        <el-button
          size="small"
          @click="resetToDefault"
        >
          重置默认配置
        </el-button>
        <el-button
          size="small"
          type="primary"
          @click="applyBatchConfig"
        >
          批量配置
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { GenTableColumn } from '@/types/generator'

interface Props {
  modelValue: GenTableColumn[]
  loading?: boolean
}

interface Emits {
  (event: 'update:modelValue', value: GenTableColumn[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 响应式数据
const expandAll = ref(false)
const selectedFields = ref<GenTableColumn[]>([])

// 计算属性
const fieldList = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 方法
const handleSelectionChange = (selection: GenTableColumn[]) => {
  selectedFields.value = selection
}

const selectAllFields = (selected: boolean) => {
  const updates = fieldList.value.map(field => ({
    ...field,
    isInsert: selected && !field.isPk,
    isEdit: selected && !field.isPk,
    isList: selected,
    isQuery: selected && !field.isPk
  }))

  emit('update:modelValue', updates)
  ElMessage.success(selected ? '已全选所有字段' : '已取消全选')
}

const getGoTypeTagType = (goType: string) => {
  switch (goType) {
    case 'string':
      return ''
    case 'int':
    case 'int64':
    case 'float64':
      return 'warning'
    case 'bool':
      return 'success'
    case 'time.Time':
      return 'info'
    default:
      return 'info'
  }
}

const resetToDefault = async () => {
  const result = await ElMessageBox.confirm(
    '确定要重置为默认配置吗？这将覆盖当前的自定义设置。',
    '重置确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(() => false)

  if (result) {
    // 重置到默认配置逻辑
    const resetFields = fieldList.value.map(field => ({
      ...field,
      isRequired: field.isNullable === 'NO' && !field.isPk,
      isInsert: !field.isPk && !field.isIncrement,
      isEdit: !field.isPk,
      isList: true,
      isQuery: isDefaultQueryField(field.columnName),
      queryType: getDefaultQueryType(field.goType),
      htmlType: getDefaultHtmlType(field.columnName, field.goType),
      dictType: getDefaultDictType(field.columnName)
    }))

    emit('update:modelValue', resetFields)
    ElMessage.success('已重置为默认配置')
  }
}

const applyBatchConfig = () => {
  // TODO: 实现批量配置对话框
  ElMessage.info('批量配置功能开发中')
}

// 工具函数
const isDefaultQueryField = (columnName: string): boolean => {
  const queryFields = ['name', 'title', 'status', 'type', 'code', 'keyword']
  const lowerName = columnName.toLowerCase()
  return queryFields.some(field => lowerName.includes(field))
}

const getDefaultQueryType = (goType: string): string => {
  switch (goType) {
    case 'string':
      return 'LIKE'
    case 'time.Time':
      return 'BETWEEN'
    default:
      return 'EQ'
  }
}

const getDefaultHtmlType = (columnName: string, goType: string): string => {
  const lowerName = columnName.toLowerCase()

  if (lowerName.includes('password')) return 'password'
  if (lowerName.includes('email')) return 'email'
  if (lowerName.includes('url')) return 'url'
  if (lowerName.includes('image') || lowerName.includes('avatar')) return 'upload'
  if (lowerName.includes('content') || lowerName.includes('description') ||
      lowerName.includes('remark')) return 'textarea'
  if (lowerName.includes('status') || lowerName.includes('type')) return 'select'
  if (lowerName.includes('date') || lowerName.includes('time')) return 'datetime'

  switch (goType) {
    case 'bool':
      return 'radio'
    case 'time.Time':
      return 'datetime'
    default:
      return 'input'
  }
}

const getDefaultDictType = (columnName: string): string => {
  const lowerName = columnName.toLowerCase()
  const dictMappings: { [key: string]: string } = {
    'status': 'sys_status',
    'type': 'sys_type',
    'sex': 'sys_user_sex',
    'gender': 'sys_user_sex'
  }

  for (const [key, value] of Object.entries(dictMappings)) {
    if (lowerName.includes(key)) {
      return value
    }
  }

  return ''
}
</script>

<style scoped lang="scss">
.field-config-table {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    padding: 12px 16px;
    background: #f8f9fa;
    border-radius: 6px;

    .header-info {
      display: flex;
      align-items: center;
      font-size: 14px;
      color: #606266;

      span {
        font-weight: 600;
      }
    }

    .header-actions {
      display: flex;
      gap: 8px;
    }
  }

  .table-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
    padding: 12px 16px;
    background: #f8f9fa;
    border-radius: 6px;

    .footer-actions {
      display: flex;
      gap: 8px;
    }
  }
}

.field-expand-config {
  padding: 16px;
  background: #fafafa;
  border-radius: 6px;
  margin: 8px;
}

:deep(.el-table) {
  font-size: 13px;

  .el-table__cell {
    padding: 8px 0;
  }

  .el-input__inner {
    font-size: 12px;
  }

  .el-select {
    font-size: 12px;
  }

  .el-checkbox {
    margin-right: 0;
  }
}

:deep(.el-table__expand-icon) {
  font-size: 14px;
}
</style>