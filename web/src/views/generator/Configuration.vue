<template>
  <div class="configuration-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">代码配置 - {{ tableInfo?.tableName }}</h1>
        <p class="description">{{ tableInfo?.tableComment || '配置代码生成参数' }}</p>
      </div>
      <div class="header-actions">
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <el-button type="success" @click="handlePreview" :loading="previewLoading">
          <el-icon><View /></el-icon>
          预览代码
        </el-button>
        <el-button type="primary" @click="handleGenerate" :loading="generateLoading">
          <el-icon><Tools /></el-icon>
          生成代码
        </el-button>
      </div>
    </div>

    <!-- 配置表单 - 添加滚动容器 -->
    <div class="config-scroll-container">
      <div class="config-container">
      <el-form
        ref="configFormRef"
        :model="configForm"
        :rules="configRules"
        label-width="120px"
        label-position="left"
      >
        <!-- 基础信息配置 -->
        <el-card class="config-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>基础信息</span>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="表名" prop="tableName">
                <el-input v-model="configForm.tableName" disabled />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="表注释" prop="tableComment">
                <el-input v-model="configForm.tableComment" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="业务名称" prop="businessName">
                <el-input
                  v-model="configForm.businessName"
                  placeholder="如：user、product"
                  @input="updateDerivedFields"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="功能名称" prop="functionName">
                <el-input
                  v-model="configForm.functionName"
                  placeholder="如：用户管理、产品管理"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="模块名称" prop="moduleName">
                <el-input
                  v-model="configForm.moduleName"
                  placeholder="如：system、business"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="类名" prop="className">
                <el-input v-model="configForm.className" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="包名" prop="packageName">
                <el-input v-model="configForm.packageName" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="作者" prop="author">
                <el-input v-model="configForm.author" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <!-- 菜单配置 -->
        <el-card class="config-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><Menu /></el-icon>
              <span>菜单配置</span>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="父级菜单" prop="parentMenuId">
                <MenuSelector
                  v-model="configForm.parentMenuId"
                  :menus="systemMenus"
                  placeholder="请选择父级菜单"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="菜单名称" prop="menuName">
                <el-input v-model="configForm.menuName" />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="菜单URL" prop="menuUrl">
                <el-input v-model="configForm.menuUrl" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="菜单图标" prop="menuIcon">
                <el-input v-model="configForm.menuIcon" placeholder="如：User、Setting">
                  <template #append>
                    <el-button @click="showIconSelector = true">选择</el-button>
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <!-- 权限配置 - 改为只读显示，不允许编辑 -->
        <el-card class="config-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><Key /></el-icon>
              <span>权限配置</span>
              <span style="margin-left: auto; font-size: 12px; color: #909399; font-weight: normal;">
                权限将根据模块名和业务名自动生成
              </span>
            </div>
          </template>

          <div class="permission-display">
            <el-row v-if="generatedPermissions.length > 0">
              <el-col :span="24">
                <div class="permission-list">
                  <el-tag
                    v-for="permission in generatedPermissions"
                    :key="permission"
                    class="permission-tag"
                    effect="plain"
                  >
                    {{ permission }}
                  </el-tag>
                </div>
              </el-col>
            </el-row>
            <div v-else class="no-permission-tip">
              <el-text type="info">请先配置模块名和业务名，权限将自动生成</el-text>
            </div>
          </div>
        </el-card>

        <!-- 字段配置 -->
        <el-card class="config-card" shadow="never">
          <template #header>
            <div class="card-header">
              <el-icon><Grid /></el-icon>
              <span>字段配置</span>
              <el-button
                size="small"
                type="text"
                @click="resetFieldConfig"
                style="margin-left: auto;"
              >
                重置配置
              </el-button>
            </div>
          </template>

          <FieldConfigTable
            v-model="configForm.columns"
            :loading="fieldLoading"
          />
        </el-card>
      </el-form>
    </div>
    </div>

    <!-- 图标选择器对话框 -->
    <el-dialog v-model="showIconSelector" title="选择图标" width="60%">
      <!-- TODO: 实现图标选择器组件 -->
      <div style="text-align: center; padding: 40px;">
        图标选择器组件待实现
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, FormInstance } from 'element-plus'
import { ArrowLeft, View, Tools, InfoFilled, Menu, Key, Grid } from '@element-plus/icons-vue'
import { getTableColumns, getSystemMenus, saveGenConfig, generateCode, previewCode } from '@/api/generator'
import type { TableColumn, GenTableConfig, SystemMenu, GenTableColumn } from '@/types/generator'
import MenuSelector from '@/components/Generator/MenuSelector.vue'
import FieldConfigTable from '@/components/Generator/FieldConfigTable.vue'

const route = useRoute()
const router = useRouter()

// 响应式数据
const configFormRef = ref<FormInstance>()
const tableInfo = ref<{ tableName: string; tableComment?: string }>()
const systemMenus = ref<SystemMenu[]>([])
const fieldLoading = ref(false)
const previewLoading = ref(false)
const generateLoading = ref(false)
const showIconSelector = ref(false)

const configForm = reactive<Partial<GenTableConfig>>({
  tableName: '',
  tableComment: '',
  businessName: '',
  moduleName: '',
  functionName: '',
  className: '',
  packageName: '',
  author: 'system',
  parentMenuId: undefined,
  menuName: '',
  menuUrl: '',
  menuIcon: '',
  permissions: [],
  columns: []
})

// 表单验证规则
const configRules = {
  businessName: [
    { required: true, message: '请输入业务名称', trigger: 'blur' }
  ],
  functionName: [
    { required: true, message: '请输入功能名称', trigger: 'blur' }
  ],
  moduleName: [
    { required: true, message: '请输入模块名称', trigger: 'blur' }
  ],
  className: [
    { required: true, message: '请输入类名', trigger: 'blur' }
  ],
  packageName: [
    { required: true, message: '请输入包名', trigger: 'blur' }
  ],
  menuName: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' }
  ]
}

// 计算属性
const tableId = computed(() => route.params.tableId as string)

// 根据模块名和业务名自动生成权限
const generatedPermissions = computed(() => {
  const moduleName = configForm.moduleName || ''
  const businessName = configForm.businessName || ''

  if (!moduleName || !businessName) return []

  return [
    `${moduleName}:${businessName}:list`,
    `${moduleName}:${businessName}:add`,
    `${moduleName}:${businessName}:edit`,
    `${moduleName}:${businessName}:delete`,
    `${moduleName}:${businessName}:view`
  ]
})

// 方法
const loadTableColumns = async () => {
  if (!tableId.value) return

  fieldLoading.value = true
  try {
    const response = await getTableColumns(tableId.value)
    const columns = response.data

    // 设置表信息
    tableInfo.value = {
      tableName: tableId.value,
      tableComment: ''
    }

    // 转换字段配置
    configForm.columns = columns.map(col => ({
      id: 0,
      tableConfigId: 0,
      columnName: col.columnName,
      columnComment: col.columnComment,
      columnType: col.columnType,
      goType: col.goType,
      goField: col.goField,
      isPk: col.isPk,
      isIncrement: col.isIncrement,
      isRequired: col.isRequired,
      isInsert: col.isInsert,
      isEdit: col.isEdit,
      isList: col.isList,
      isQuery: col.isQuery,
      queryType: col.queryType,
      htmlType: col.htmlType,
      dictType: col.dictType,
      sort: col.ordinalPosition,
      createdAt: '',
      updatedAt: ''
    }))

    // 初始化基础配置
    initDefaultConfig()
  } catch (error) {
    ElMessage.error('获取表字段信息失败')
    console.error(error)
  } finally {
    fieldLoading.value = false
  }
}

const loadSystemMenus = async () => {
  try {
    const response = await getSystemMenus()
    systemMenus.value = response.data
  } catch (error) {
    ElMessage.error('获取系统菜单失败')
    console.error(error)
  }
}

const initDefaultConfig = () => {
  if (!tableId.value) return

  const tableName = tableId.value
  const businessName = tableName.replace(/^[a-z_]+_/, '').replace(/_/g, '')

  configForm.tableName = tableName
  configForm.businessName = businessName
  configForm.functionName = businessName + '管理'
  configForm.moduleName = 'system'
  configForm.className = toPascalCase(businessName)
  configForm.packageName = 'com.example.' + businessName
  configForm.menuName = configForm.functionName
  configForm.menuUrl = `/${businessName}`
  configForm.menuIcon = 'Grid'
}

const updateDerivedFields = () => {
  const businessName = configForm.businessName || ''
  if (businessName) {
    configForm.className = toPascalCase(businessName)
    configForm.packageName = `com.example.${businessName}`
    configForm.functionName = businessName + '管理'
    configForm.menuName = configForm.functionName
    configForm.menuUrl = `/${businessName}`
    // 自动更新权限配置
    configForm.permissions = generatedPermissions.value
  }
}

const resetFieldConfig = async () => {
  const result = await ElMessageBox.confirm(
    '确定要重置字段配置吗？这将覆盖当前的自定义配置。',
    '重置确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).catch(() => false)

  if (result) {
    await loadTableColumns()
    ElMessage.success('字段配置已重置')
  }
}

const handlePreview = async () => {
  if (!configFormRef.value) return

  const valid = await configFormRef.value.validate().catch(() => false)
  if (!valid) {
    ElMessage.error('请完善配置信息')
    return
  }

  previewLoading.value = true
  try {
    // 确保权限配置是最新的
    configForm.permissions = generatedPermissions.value

    // 先保存配置
    const saveResponse = await saveGenConfig(configForm)
    const configId = saveResponse.data.id

    // 跳转到预览页面
    router.push({
      name: 'GeneratorPreview',
      params: { configId: configId.toString() }
    })
  } catch (error) {
    ElMessage.error('预览代码失败')
    console.error(error)
  } finally {
    previewLoading.value = false
  }
}

const handleGenerate = async () => {
  if (!configFormRef.value) return

  const valid = await configFormRef.value.validate().catch(() => false)
  if (!valid) {
    ElMessage.error('请完善配置信息')
    return
  }

  const result = await ElMessageBox.confirm(
    '确定要生成代码吗？',
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
    // 确保权限配置是最新的
    configForm.permissions = generatedPermissions.value

    // 保存配置
    const saveResponse = await saveGenConfig(configForm)
    const configId = saveResponse.data.id

    // 生成代码
    const generateResponse = await generateCode({
      configId,
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
  router.push({ name: 'GeneratorTables' })
}

// 工具函数
const toPascalCase = (str: string): string => {
  return str.replace(/(?:^|_)([a-z])/g, (_, letter) => letter.toUpperCase())
}

// 生命周期
onMounted(() => {
  loadTableColumns()
  loadSystemMenus()
})
</script>

<style scoped lang="scss">
.configuration-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 20px;
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

.config-scroll-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.config-container {
  max-width: 1200px;
  margin: 0 auto;
  padding-bottom: 40px;
}

.config-card {
  margin-bottom: 20px;

  :deep(.el-card__header) {
    padding: 18px 20px;
    border-bottom: 1px solid #f0f0f0;
  }

  :deep(.el-card__body) {
    padding: 20px;
  }
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;

  .el-icon {
    color: #409eff;
  }
}

.permission-display {
  padding: 20px 0;
}

.permission-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.permission-tag {
  font-family: 'Monaco', 'Consolas', monospace;
  padding: 4px 12px;
  background-color: #f5f7fa;
  border-color: #dcdfe6;
}

.no-permission-tip {
  text-align: center;
  padding: 40px 0;
  color: #909399;
}
</style>