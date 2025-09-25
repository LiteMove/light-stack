<template>
  <div class="dict-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon class="title-icon"><Operation /></el-icon>
          字典管理
        </h2>
        <p class="page-desc">
          管理系统字典类型和字典数据，维护下拉框选项等基础数据
        </p>
      </div>
      <div class="header-actions">
        <el-button
          type="primary"
          :icon="Plus"
          @click="handleAddType"
          size="default"
        >
          添加字典类型
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

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 左侧字典类型列表 -->
      <el-card class="type-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span class="header-title">
              <el-icon><FolderOpened /></el-icon>
              字典类型
            </span>
            <el-input
              v-model="typeSearchKeyword"
              placeholder="请输入类型名称或编码"
              :prefix-icon="Search"
              clearable
              @input="handleTypeSearch"
              size="small"
              style="width: 200px;"
            />
          </div>
        </template>

        <div class="type-list">
          <div
            v-for="type in filteredTypes"
            :key="type.id"
            :class="['type-item', { 'active': selectedType?.id === type.id }]"
            @click="handleSelectType(type)"
          >
            <div class="type-info">
              <div class="type-name">{{ type.name }}</div>
              <div class="type-code">{{ type.type }}</div>
              <div class="type-desc" v-if="type.description">{{ type.description }}</div>
            </div>
            <div class="type-actions">
              <el-tag
                :type="type.status === 1 ? 'success' : 'danger'"
                size="small"
                effect="light"
              >
                {{ type.status === 1 ? '启用' : '禁用' }}
              </el-tag>
              <el-dropdown @command="handleTypeCommand" placement="bottom-end">
                <el-button type="text" :icon="MoreFilled" size="small" />
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="{ action: 'edit', data: type }">
                      <el-icon><Edit /></el-icon>
                      修改
                    </el-dropdown-item>
                    <el-dropdown-item :command="{ action: 'delete', data: type }" divided>
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>

          <!-- 空状态 -->
          <el-empty v-if="filteredTypes.length === 0" description="暂无字典类型" />
        </div>
      </el-card>

      <!-- 右侧字典数据管理 -->
      <el-card class="data-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span class="header-title" v-if="selectedType">
              <el-icon><Grid /></el-icon>
              {{ selectedType.name }}
              <el-tag size="small" class="type-tag">{{ selectedType.type }}</el-tag>
            </span>
            <span class="header-title" v-else>
              <el-icon><Grid /></el-icon>
              字典数据管理
            </span>
            <div class="header-actions" v-if="selectedType">
              <el-button
                type="primary"
                size="small"
                :icon="Plus"
                @click="handleAddData"
              >
                新增字典值
              </el-button>
            </div>
          </div>
        </template>

        <!-- 搜索区域 -->
        <div class="search-section" v-if="selectedType">
          <el-form :model="dataSearchForm" inline class="search-form">
            <el-form-item label="字典展示值">
              <el-input
                v-model="dataSearchForm.label"
                placeholder="请输入显示值"
                clearable
                @keyup.enter="handleDataSearch"
                @clear="handleDataSearch"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="状态">
              <el-select
                v-model="dataSearchForm.status"
                placeholder="状态筛选"
                clearable
                @change="handleDataSearch"
                style="width: 120px"
              >
                <el-option label="全部" :value="0" />
                <el-option label="启用" :value="1" />
                <el-option label="禁用" :value="2" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :icon="Search"
                @click="handleDataSearch"
                :loading="dataLoading"
              >
                查询
              </el-button>
              <el-button @click="handleResetDataSearch" :disabled="dataLoading">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 字典数据表格 -->
        <div v-if="selectedType">
          <el-table
            :data="dictDataList"
            v-loading="dataLoading"
            stripe
            border
            style="width: 100%"
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column prop="label" label="字典展示值" width="150" show-overflow-tooltip />
            <el-table-column prop="value" label="字典值" width="150" show-overflow-tooltip />
            <el-table-column prop="sortOrder" label="排序" width="80" align="center" />
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="isDefault" label="是否默认" width="100" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.isDefault" type="warning" size="small">默认</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" show-overflow-tooltip />
            <el-table-column prop="createdAt" label="创建时间" width="160" show-overflow-tooltip />
            <el-table-column label="操作" width="150" align="center" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleEditData(row)">
                  <el-icon><Edit /></el-icon>
                  修改
                </el-button>
                <el-button type="danger" link size="small" @click="handleDeleteData(row)">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="dataQuery.page"
              v-model:page-size="dataQuery.pageSize"
              :total="dataTotal"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleDataSizeChange"
              @current-change="handleDataCurrentChange"
            />
          </div>
        </div>

        <!-- 未选择字典类型时的提示 -->
        <div v-else class="no-selection">
          <el-empty description="请先选择左侧的字典类型">
            <template #image>
              <el-icon size="100"><FolderOpened /></el-icon>
            </template>
          </el-empty>
        </div>
      </el-card>
    </div>

    <!-- 字典类型编辑对话框 -->
    <el-dialog
      v-model="typeDialogVisible"
      :title="isEditType ? '修改字典类型' : '新增字典类型'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="typeFormRef"
        :model="typeForm"
        :rules="typeRules"
        label-width="100px"
      >
        <el-form-item label="类型名称" prop="name">
          <el-input v-model="typeForm.name" placeholder="请输入类型名称" />
        </el-form-item>
        <el-form-item label="类型编码" prop="type">
          <el-input
            v-model="typeForm.type"
            placeholder="请输入类型编码"
            :disabled="isEditType"
          />
          <div class="form-tip">类型编码用于程序中调用，创建后不可修改</div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="typeForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="typeForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="typeDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveType" :loading="typeSubmitting">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 字典数据编辑对话框 -->
    <el-dialog
      v-model="dataDialogVisible"
      :title="isEditData ? '修改字典数据' : '新增字典数据'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="dataFormRef"
        :model="dataForm"
        :rules="dataRules"
        label-width="100px"
      >
        <el-form-item label="字典类型" prop="dictType">
          <el-input v-model="dataForm.dictType" disabled />
        </el-form-item>
        <el-form-item label="字典展示值" prop="label">
          <el-input v-model="dataForm.label" placeholder="请输入展示值" />
        </el-form-item>
        <el-form-item label="字典值" prop="value">
          <el-input v-model="dataForm.value" placeholder="请输入字典值" />
        </el-form-item>
        <el-form-item label="排序" prop="sortOrder">
          <el-input-number v-model="dataForm.sortOrder" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="CSS类名" prop="cssClass">
          <el-input v-model="dataForm.cssClass" placeholder="请输入CSS类名" />
        </el-form-item>
        <el-form-item label="列表类名" prop="listClass">
          <el-input v-model="dataForm.listClass" placeholder="请输入列表类名" />
        </el-form-item>
        <el-form-item label="是否默认" prop="isDefault">
          <el-switch v-model="dataForm.isDefault" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="dataForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="dataForm.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dataDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveData" :loading="dataSubmitting">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus'
import {
  Plus,
  RefreshRight,
  Search,
  Edit,
  Delete,
  Operation,
  FolderOpened,
  Grid,
  MoreFilled
} from '@element-plus/icons-vue'
import {
  getDictTypeList,
  createDictType,
  updateDictType,
  deleteDictType,
  getDictDataList,
  createDictData,
  updateDictData,
  deleteDictData,
  type DictType,
  type DictData
} from '@/api/dict'

// 响应式数据
const loading = ref(false)
const dataLoading = ref(false)
const dictTypes = ref<DictType[]>([])
const dictDataList = ref<DictData[]>([])
const selectedType = ref<DictType | null>(null)
const selectedRows = ref<DictData[]>([])

// 字典类型搜索
const typeSearchKeyword = ref('')

// 字典数据搜索
const dataSearchForm = reactive({
  label: '',
  status: 0
})

// 字典数据分页
const dataQuery = reactive({
  page: 1,
  pageSize: 10
})
const dataTotal = ref(0)

// 字典类型对话框
const typeDialogVisible = ref(false)
const typeSubmitting = ref(false)
const isEditType = ref(false)
const typeFormRef = ref<FormInstance>()
const typeForm = reactive({
  id: 0,
  name: '',
  type: '',
  description: '',
  status: 1
})

// 字典数据对话框
const dataDialogVisible = ref(false)
const dataSubmitting = ref(false)
const isEditData = ref(false)
const dataFormRef = ref<FormInstance>()
const dataForm = reactive({
  id: 0,
  dictType: '',
  label: '',
  value: '',
  sortOrder: 0,
  cssClass: '',
  listClass: '',
  isDefault: false,
  status: 1,
  remark: ''
})

// 表单验证规则
const typeRules = {
  name: [
    { required: true, message: '请输入类型名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请输入类型编码', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '编码必须以字母开头，只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  description: [
    { max: 255, message: '长度不能超过 255 个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

const dataRules = {
  label: [
    { required: true, message: '请输入展示值', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  value: [
    { required: true, message: '请输入字典值', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  sortOrder: [
    { type: 'number', message: '排序必须是数字', trigger: 'blur' }
  ],
  cssClass: [
    { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
  ],
  listClass: [
    { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  remark: [
    { max: 255, message: '长度不能超过 255 个字符', trigger: 'blur' }
  ]
}

// 计算属性
const filteredTypes = computed(() => {
  if (!typeSearchKeyword.value) {
    return dictTypes.value
  }
  const keyword = typeSearchKeyword.value.toLowerCase()
  return dictTypes.value.filter(type =>
    type.name.toLowerCase().includes(keyword) ||
    type.type.toLowerCase().includes(keyword)
  )
})

// 监听选中的字典类型变化
watch(selectedType, (newType) => {
  if (newType) {
    loadDictDataList()
  } else {
    dictDataList.value = []
  }
})

// 方法
const loadDictTypes = async () => {
  try {
    loading.value = true
    const response = await getDictTypeList({ page: 1, pageSize: 1000 })
    dictTypes.value = response.data.list || []
  } catch (error) {
    console.error('加载字典类型失败:', error)
    ElMessage.error('加载字典类型失败')
  } finally {
    loading.value = false
  }
}

const loadDictDataList = async () => {
  if (!selectedType.value) return

  try {
    dataLoading.value = true
    const response = await getDictDataList(selectedType.value.type, {
      page: dataQuery.page,
      pageSize: dataQuery.pageSize,
      status: dataSearchForm.status || undefined,
      label: dataSearchForm.label || undefined
    })
    dictDataList.value = response.data.list || []
    dataTotal.value = response.data.total || 0
  } catch (error) {
    console.error('加载字典数据失败:', error)
    ElMessage.error('加载字典数据失败')
  } finally {
    dataLoading.value = false
  }
}

const handleSelectType = (type: DictType) => {
  selectedType.value = type
  // 重置数据搜索条件和分页
  dataSearchForm.label = ''
  dataSearchForm.status = 0
  dataQuery.page = 1
}

const handleAddType = () => {
  isEditType.value = false
  Object.assign(typeForm, {
    id: 0,
    name: '',
    type: '',
    description: '',
    status: 1
  })
  typeDialogVisible.value = true
}

const handleTypeCommand = ({ action, data }: { action: string; data: DictType }) => {
  if (action === 'edit') {
    handleEditType(data)
  } else if (action === 'delete') {
    handleDeleteType(data)
  }
}

const handleEditType = (type: DictType) => {
  isEditType.value = true
  Object.assign(typeForm, {
    id: type.id,
    name: type.name,
    type: type.type,
    description: type.description,
    status: type.status
  })
  typeDialogVisible.value = true
}

const handleDeleteType = async (type: DictType) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除字典类型 "${type.name}" 吗？删除后该类型下的所有字典数据也将被删除！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await deleteDictType(type.id)
    ElMessage.success('删除成功')

    // 如果删除的是当前选中的类型，清空选中状态
    if (selectedType.value?.id === type.id) {
      selectedType.value = null
    }

    loadDictTypes()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleSaveType = async () => {
  if (!typeFormRef.value) return

  try {
    await typeFormRef.value.validate()
    typeSubmitting.value = true

    if (isEditType.value) {
      await updateDictType(typeForm.id, {
        name: typeForm.name,
        type: typeForm.type,
        description: typeForm.description,
        status: typeForm.status
      })
      ElMessage.success('修改成功')
    } else {
      await createDictType({
        name: typeForm.name,
        type: typeForm.type,
        description: typeForm.description,
        status: typeForm.status
      })
      ElMessage.success('创建成功')
    }

    typeDialogVisible.value = false
    loadDictTypes()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    typeSubmitting.value = false
  }
}

const handleAddData = () => {
  if (!selectedType.value) return

  isEditData.value = false
  Object.assign(dataForm, {
    id: 0,
    dictType: selectedType.value.type,
    label: '',
    value: '',
    sortOrder: 0,
    cssClass: '',
    listClass: '',
    isDefault: false,
    status: 1,
    remark: ''
  })
  dataDialogVisible.value = true
}

const handleEditData = (data: DictData) => {
  isEditData.value = true
  Object.assign(dataForm, {
    id: data.id,
    dictType: data.dictType,
    label: data.label,
    value: data.value,
    sortOrder: data.sortOrder,
    cssClass: data.cssClass,
    listClass: data.listClass,
    isDefault: data.isDefault,
    status: data.status,
    remark: data.remark
  })
  dataDialogVisible.value = true
}

const handleDeleteData = async (data: DictData) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除字典数据 "${data.label}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await deleteDictData(data.id)
    ElMessage.success('删除成功')
    loadDictDataList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleSaveData = async () => {
  if (!dataFormRef.value) return

  try {
    await dataFormRef.value.validate()
    dataSubmitting.value = true

    const formData = {
      dictType: dataForm.dictType,
      label: dataForm.label,
      value: dataForm.value,
      sortOrder: dataForm.sortOrder,
      cssClass: dataForm.cssClass,
      listClass: dataForm.listClass,
      isDefault: dataForm.isDefault,
      status: dataForm.status,
      remark: dataForm.remark
    }

    if (isEditData.value) {
      await updateDictData(dataForm.id, formData)
      ElMessage.success('修改成功')
    } else {
      await createDictData(formData)
      ElMessage.success('创建成功')
    }

    dataDialogVisible.value = false
    loadDictDataList()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    dataSubmitting.value = false
  }
}

const handleTypeSearch = () => {
  // 类型搜索通过计算属性实现，这里只需要触发重新计算
}

const handleDataSearch = () => {
  dataQuery.page = 1
  loadDictDataList()
}

const handleResetDataSearch = () => {
  dataSearchForm.label = ''
  dataSearchForm.status = 0
  dataQuery.page = 1
  loadDictDataList()
}

const handleSelectionChange = (selection: DictData[]) => {
  selectedRows.value = selection
}

const handleDataSizeChange = (size: number) => {
  dataQuery.pageSize = size
  dataQuery.page = 1
  loadDictDataList()
}

const handleDataCurrentChange = (page: number) => {
  dataQuery.page = page
  loadDictDataList()
}

const refreshData = () => {
  loadDictTypes()
  if (selectedType.value) {
    loadDictDataList()
  }
}

// 生命周期
onMounted(() => {
  loadDictTypes()
})
</script>

<style lang="scss" scoped>
.dict-management {
  padding: 20px;
  background: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  .header-left {
    flex: 1;

    .page-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 24px;
      font-weight: 600;
      color: #1f2937;
      margin: 0 0 8px 0;

      .title-icon {
        color: #6366f1;
      }
    }

    .page-desc {
      color: #6b7280;
      font-size: 14px;
      margin: 0;
    }
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }
}

.main-content {
  display: flex;
  gap: 20px;
  height: calc(100vh - 180px);
}

.type-card {
  width: 350px;
  flex-shrink: 0;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      color: #1f2937;
    }
  }

  .type-list {
    max-height: calc(100vh - 320px);
    overflow-y: auto;

    .type-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px;
      margin-bottom: 8px;
      border: 1px solid #e5e7eb;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        border-color: #6366f1;
        background-color: #f8fafc;
      }

      &.active {
        border-color: #6366f1;
        background-color: #eef2ff;
      }

      .type-info {
        flex: 1;

        .type-name {
          font-weight: 600;
          color: #1f2937;
          margin-bottom: 4px;
        }

        .type-code {
          font-size: 12px;
          color: #6b7280;
          background: #f3f4f6;
          padding: 2px 6px;
          border-radius: 3px;
          display: inline-block;
          margin-bottom: 4px;
        }

        .type-desc {
          font-size: 12px;
          color: #9ca3af;
        }
      }

      .type-actions {
        display: flex;
        align-items: center;
        gap: 8px;
      }
    }
  }
}

.data-card {
  flex: 1;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      color: #1f2937;

      .type-tag {
        margin-left: 8px;
      }
    }
  }

  .search-section {
    margin-bottom: 20px;
    padding: 16px;
    background: #f9fafb;
    border-radius: 6px;

    .search-form {
      .el-form-item {
        margin-bottom: 0;
      }
    }
  }

  .no-selection {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 400px;
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.form-tip {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}

.dialog-footer {
  text-align: right;
}

// 响应式设计
@media (max-width: 1200px) {
  .main-content {
    flex-direction: column;
    height: auto;
  }

  .type-card {
    width: 100%;
    margin-bottom: 20px;

    .type-list {
      max-height: 300px;
    }
  }
}
</style>