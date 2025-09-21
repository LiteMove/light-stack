<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑菜单' : '添加菜单'"
    width="700px"
    :before-close="handleClose"
    destroy-on-close
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      label-position="right"
      v-loading="loading"
      element-loading-text="处理中..."
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="父级菜单" prop="parentId">
            <el-tree-select
              v-model="form.parentId"
              :data="parentOptions"
              :props="treeProps"
              check-strictly
              placeholder="请选择父级菜单"
              clearable
              style="width: 100%"
              :render-after-expand="false"
              default-expand-all
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="菜单类型" prop="type">
            <el-radio-group v-model="form.type" @change="handleTypeChange">
              <el-radio label="directory">目录</el-radio>
              <el-radio label="menu">菜单</el-radio>
              <el-radio label="permission">权限</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="菜单名称" prop="name">
            <el-input
              v-model="form.name"
              placeholder="请输入菜单名称"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="菜单编码" prop="code">
            <el-input
              v-model="form.code"
              placeholder="请输入菜单编码"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>
        </el-col>
      </el-row>

      <!-- 路由和组件配置 -->
      <template v-if="form.type === 'directory' || form.type === 'menu'">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="路由路径" prop="path">
              <el-input
                v-model="form.path"
                placeholder="请输入路由路径，如：/system/menu"
                maxlength="255"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单图标" prop="icon">
              <el-input
                v-model="form.icon"
                placeholder="请输入图标名称"
                maxlength="100"
              >
                <template #prefix>
                  <el-icon v-if="form.icon">
                    <component :is="getIconComponent(form.icon)" />
                  </el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item v-if="form.type === 'menu'" label="组件路径" prop="component">
          <el-input
            v-model="form.component"
            placeholder="请输入组件路径，如：system/menu/index"
            maxlength="255"
          >
            <template #prepend>src/views/</template>
            <template #append>.vue</template>
          </el-input>
        </el-form-item>
      </template>

      <!-- 权限配置 -->
      <template v-if="form.type === 'permission'">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="资源" prop="resource">
              <el-input
                v-model="form.resource"
                placeholder="请输入资源路径"
                maxlength="255"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作" prop="action">
              <el-input
                v-model="form.action"
                placeholder="请输入操作名称"
                maxlength="50"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </template>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="显示排序" prop="sortOrder">
            <el-input-number
              v-model="form.sortOrder"
              :min="0"
              :max="999"
              controls-position="right"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="状态" prop="status">
            <el-radio-group v-model="form.status">
              <el-radio :label="1">启用</el-radio>
              <el-radio :label="2">禁用</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="是否隐藏">
            <el-switch
              v-model="form.isHidden"
              active-text="隐藏"
              inactive-text="显示"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="元数据" prop="meta">
        <el-input
          v-model="form.meta"
          type="textarea"
          :rows="4"
          placeholder="请输入JSON格式的元数据，如：{&quot;keepAlive&quot;: true, &quot;title&quot;: &quot;页面标题&quot;}"
          maxlength="1000"
          show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { menuApi } from '@/api'
import type { Menu, MenuFormData } from '@/api/types'

interface Props {
  visible: boolean
  formData: Partial<Menu>
  parentOptions: Menu[]
}

interface Emits {
  (e: 'update:visible', visible: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const loading = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.formData.id)

// 表单数据
const form = ref<MenuFormData>({
  parentId: 0,
  name: '',
  code: '',
  type: 'menu',
  path: '',
  component: '',
  icon: '',
  resource: '',
  action: '',
  sortOrder: 0,
  isHidden: false,
  status: 1,
  meta: ''
})

// 树选择器配置
const treeProps = {
  label: 'name',
  value: 'id',
  children: 'children',
  disabled: (data: Menu) => data.type === 'permission'
}

// 表单验证规则
const rules = computed<FormRules>(() => ({
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入菜单编码', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_-]+$/, message: '只能包含字母、数字、下划线和横线', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择菜单类型', trigger: 'change' }
  ],
  path: [
    { 
      required: form.value.type === 'directory' || form.value.type === 'menu', 
      message: '请输入路由路径', 
      trigger: 'blur' 
    },
    { max: 255, message: '长度不能超过 255 个字符', trigger: 'blur' }
  ],
  component: [
    { 
      required: form.value.type === 'menu', 
      message: '请输入组件路径', 
      trigger: 'blur' 
    },
    { max: 255, message: '长度不能超过 255 个字符', trigger: 'blur' }
  ],
  icon: [
    { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
  ],
  resource: [
    { max: 255, message: '长度不能超过 255 个字符', trigger: 'blur' }
  ],
  action: [
    { max: 50, message: '长度不能超过 50 个字符', trigger: 'blur' }
  ],
  sortOrder: [
    { type: 'number', min: 0, max: 999, message: '排序值应在 0-999 之间', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  meta: [
    { validator: validateMeta, trigger: 'blur' }
  ]
}))

// 自定义验证器
function validateMeta(rule: any, value: string, callback: any) {
  if (!value) {
    callback()
    return
  }
  
  try {
    JSON.parse(value)
    callback()
  } catch (error) {
    callback(new Error('请输入有效的JSON格式'))
  }
}

// 监听表单数据变化
watch(
  () => props.formData,
  (newData) => {
    if (newData) {
      Object.assign(form.value, {
        parentId: newData.parentId || 0,
        name: newData.name || '',
        code: newData.code || '',
        type: newData.type || 'menu',
        path: newData.path || '',
        component: newData.component || '',
        icon: newData.icon || '',
        resource: newData.resource || '',
        action: newData.action || '',
        sort_order: newData.sortOrder || 0,
        isHidden: newData.isHidden || false,
        status: newData.status || 1,
        meta: newData.meta || ''
      })
    }
  },
  { immediate: true, deep: true }
)

// 菜单类型变化处理
const handleTypeChange = (type: string) => {
  // 清空相关字段
  if (type === 'permission') {
    form.value.path = ''
    form.value.component = ''
    form.value.icon = ''
  } else if (type === 'directory') {
    form.value.component = ''
    form.value.resource = ''
    form.value.action = ''
  } else if (type === 'menu') {
    form.value.resource = ''
    form.value.action = ''
  }
  
  // 触发表单验证
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

// 获取图标组件
const getIconComponent = (icon: string) => {
  // 动态导入 Element Plus 图标
  const iconMap: Record<string, any> = {
    Menu: 'Menu',
    Folder: 'Folder', 
    Document: 'Document',
    Key: 'Key',
    Setting: 'Setting',
    User: 'User',
    Lock: 'Lock',
    Monitor: 'Monitor',
    DataLine: 'DataLine',
    House: 'House',
    UserFilled: 'UserFilled',
    Tools: 'Tools'
  }
  return iconMap[icon] || 'Menu'
}

// 提交表单
const handleSubmit = async () => {
  try {
    const valid = await formRef.value?.validate()
    if (!valid) return

    loading.value = true

    if (isEdit.value) {
      await menuApi.updateMenu(props.formData.id!, form.value)
      ElMessage.success('更新成功')
    } else {
      await menuApi.createMenu(form.value)
      ElMessage.success('创建成功')
    }

    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    loading.value = false
  }
}

// 关闭弹窗
const handleClose = () => {
  formRef.value?.resetFields()
  emit('update:visible', false)
}

// 生成建议的菜单编码
const generateCode = () => {
  const name = form.value.name
  if (!name) return

  // 简单的拼音转换逻辑，实际项目中可以使用专门的拼音库
  const code = name
    .replace(/[\u4e00-\u9fa5]/g, (char) => {
      // 这里可以添加中文转拼音的逻辑
      return char
    })
    .toLowerCase()
    .replace(/[^a-z0-9]/g, '_')

  if (!form.value.code) {
    form.value.code = code
  }
}

// 路径建议
const suggestPath = () => {
  if (!form.value.code || form.value.path) return

  const parentPath = getParentPath(form.value.parentId)
  form.value.path = parentPath ? `${parentPath}/${form.value.code}` : `/${form.value.code}`
}

// 获取父路径
const getParentPath = (parentId: number): string => {
  if (!parentId) return ''
  
  const findPath = (menus: Menu[], targetId: number): string => {
    for (const menu of menus) {
      if (menu.id === targetId) {
        return menu.path || ''
      }
      if (menu.children?.length) {
        const path = findPath(menu.children, targetId)
        if (path) return path
      }
    }
    return ''
  }
  
  return findPath(props.parentOptions, parentId)
}

// 监听名称变化，自动生成编码和路径
watch(() => form.value.name, () => {
  if (!isEdit.value) {
    generateCode()
    suggestPath()
  }
})

watch(() => form.value.code, () => {
  if (!isEdit.value) {
    suggestPath()
  }
})
</script>

<style lang="scss" scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-input-group__prepend) {
  background-color: #f5f7fa;
  color: #909399;
  border: 1px solid #dcdfe6;
}

:deep(.el-input-group__append) {
  background-color: #f5f7fa;
  color: #909399;
  border: 1px solid #dcdfe6;
}

:deep(.el-textarea__inner) {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
}
</style>