<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑{{.FunctionName}}' : '新建{{.FunctionName}}'"
    width="800px"
    :before-close="handleClose"
    destroy-on-close
    :close-on-click-modal="false"
    :close-on-press-escape="false"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      label-position="right"
      v-loading="loading"
      element-loading-text="处理中..."
    >
{{- $hasGroups := false }}
{{- range .FormFields }}
{{- if or (contains .ColumnComment "基本") (contains .ColumnComment "信息") }}
{{- if not $hasGroups }}
{{- $hasGroups = true }}
      <!-- 基本信息 -->
      <el-divider content-position="left">
        <el-icon><Document /></el-icon>
        基本信息
      </el-divider>
{{- end }}
{{- break }}
{{- end }}
{{- end }}

{{- $rowCount := 0 }}
{{- range .FormFields }}
{{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant")) }}
{{- if mod $rowCount 2 | eq 0 }}
      <el-row :gutter="20">
{{- end }}
        <el-col :span="{{if or (eq .HtmlType "textarea") (eq .HtmlType "editor")}}24{{else}}12{{end}}">
          <el-form-item label="{{.ColumnComment}}" prop="{{generateJSField .ColumnName}}">
{{- if eq .HtmlType "input" }}
            <el-input
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="请输入{{.ColumnComment}}"
              maxlength="{{if .ColumnType}}{{.ColumnType}}{{else}}50{{end}}"
              {{if eq .GoType "string"}}show-word-limit{{end}}
              {{if .IsPk}}:disabled="isEdit"{{end}}
            >
              <template #prefix>
                <el-icon><{{if contains .ColumnName "name"}}User{{else if contains .ColumnName "code"}}Key{{else if contains .ColumnName "phone"}}Phone{{else if contains .ColumnName "email"}}Message{{else}}Edit{{end}} /></el-icon>
              </template>
            </el-input>
{{- else if eq .HtmlType "textarea" }}
            <el-input
              v-model="form.{{generateJSField .ColumnName}}"
              type="textarea"
              placeholder="请输入{{.ColumnComment}}"
              :rows="4"
              maxlength="500"
              show-word-limit
            />
{{- else if eq .HtmlType "select" }}
            <el-select
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="请选择{{.ColumnComment}}"
              clearable
              style="width: 100%"
            >
              <!-- TODO: 根据字典类型{{.DictType}}添加选项 -->
              <el-option label="选项1" value="1" />
              <el-option label="选项2" value="2" />
            </el-select>
{{- else if eq .HtmlType "radio" }}
            <el-radio-group v-model="form.{{generateJSField .ColumnName}}">
              <!-- TODO: 根据字典类型{{.DictType}}添加选项 -->
              <el-radio :label="1">选项1</el-radio>
              <el-radio :label="2">选项2</el-radio>
            </el-radio-group>
{{- else if eq .HtmlType "checkbox" }}
            <el-checkbox-group v-model="form.{{generateJSField .ColumnName}}">
              <!-- TODO: 根据字典类型{{.DictType}}添加选项 -->
              <el-checkbox label="1">选项1</el-checkbox>
              <el-checkbox label="2">选项2</el-checkbox>
            </el-checkbox-group>
{{- else if eq .HtmlType "datetime" }}
            <el-date-picker
              v-model="form.{{generateJSField .ColumnName}}"
              type="datetime"
              placeholder="选择{{.ColumnComment}}"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%"
            />
{{- else if eq .HtmlType "date" }}
            <el-date-picker
              v-model="form.{{generateJSField .ColumnName}}"
              type="date"
              placeholder="选择{{.ColumnComment}}"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="width: 100%"
            />
{{- else if eq .HtmlType "time" }}
            <el-time-picker
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="选择{{.ColumnComment}}"
              format="HH:mm:ss"
              value-format="HH:mm:ss"
              style="width: 100%"
            />
{{- else if eq .HtmlType "number" }}
            <el-input-number
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="请输入{{.ColumnComment}}"
              :min="0"
              :precision="{{if eq .GoType "int"}}0{{else}}2{{end}}"
              style="width: 100%"
            />
{{- else if eq .HtmlType "password" }}
            <el-input
              v-model="form.{{generateJSField .ColumnName}}"
              type="password"
              placeholder="请输入{{.ColumnComment}}"
              maxlength="50"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
{{- else if eq .HtmlType "email" }}
            <el-input
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="请输入{{.ColumnComment}}"
              maxlength="100"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
{{- else if eq .HtmlType "url" }}
            <el-input
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="请输入{{.ColumnComment}}"
              maxlength="200"
            >
              <template #prefix>
                <el-icon><Link /></el-icon>
              </template>
            </el-input>
{{- else if eq .HtmlType "image" }}
            <ImageUpload
              v-model="form.{{generateJSField .ColumnName}}"
              @success="handleImageSuccess"
              @error="handleImageError"
            />
{{- else if eq .HtmlType "file" }}
            <FileUpload
              v-model="form.{{generateJSField .ColumnName}}"
              @success="handleFileSuccess"
              @error="handleFileError"
            />
{{- else if eq .HtmlType "switch" }}
            <el-switch
              v-model="form.{{generateJSField .ColumnName}}"
              :active-value="{{if eq .GoType "bool"}}true{{else}}1{{end}}"
              :inactive-value="{{if eq .GoType "bool"}}false{{else}}0{{end}}"
              active-text="是"
              inactive-text="否"
            />
{{- else }}
            <!-- 默认输入框 -->
            <el-input
              v-model="form.{{generateJSField .ColumnName}}"
              placeholder="请输入{{.ColumnComment}}"
              maxlength="50"
            />
{{- end }}
          </el-form-item>
        </el-col>
{{- $rowCount = add $rowCount 1 }}
{{- if or (mod $rowCount 2 | eq 0) (eq .HtmlType "textarea") (eq .HtmlType "editor") }}
      </el-row>
{{- if or (eq .HtmlType "textarea") (eq .HtmlType "editor") }}
{{- $rowCount = 0 }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- if mod $rowCount 2 | eq 1 }}
      </el-row>
{{- end }}
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" :disabled="loading">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{`{{ isEdit ? '保存' : '创建' }}`}}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  Document,
  User,
  Key,
  Phone,
  Message,
  Edit,
  Lock,
  Link
} from '@element-plus/icons-vue'
{{- $hasImageOrFile := false }}
{{- range .FormFields }}
{{- if or (eq .HtmlType "image") (eq .HtmlType "file") }}
{{- $hasImageOrFile = true }}
{{- end }}
{{- end }}
{{- if $hasImageOrFile }}
import ImageUpload from '@/components/ImageUpload.vue'
import FileUpload from '@/components/FileUpload.vue'
{{- end }}
import { {{toLower .BusinessName}}Api } from '@/api'
import type { {{.ClassName}} } from '@/api/types'

interface {{.ClassName}}FormData {
{{- range .FormFields }}
{{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant")) }}
  {{generateJSField .ColumnName}}: {{if eq .GoType "string"}}string{{else if eq .GoType "bool"}}boolean{{else if contains .GoType "int"}}number{{else if eq .GoType "time.Time"}}string{{else}}any{{end}}{{if not .IsRequired}} | null{{end}}
{{- end }}
{{- end }}
}

interface Props {
  visible: boolean
  formData: Partial<{{.ClassName}}>
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

const isEdit = computed(() => !!props.formData.{{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}})

// 表单数据
const form = ref<{{.ClassName}}FormData>({
{{- range .FormFields }}
{{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant")) }}
  {{generateJSField .ColumnName}}: {{getDefaultValue .}},
{{- end }}
{{- end }}
})

{{- if $hasImageOrFile }}
// 文件上传成功处理
const handleImageSuccess = (file: any) => {
  // 根据具体业务逻辑处理图片上传成功
  ElMessage.success('图片上传成功')
}

const handleImageError = (error: string) => {
  ElMessage.error(error || '图片上传失败')
}

const handleFileSuccess = (file: any) => {
  // 根据具体业务逻辑处理文件上传成功
  ElMessage.success('文件上传成功')
}

const handleFileError = (error: string) => {
  ElMessage.error(error || '文件上传失败')
}
{{- end }}

// 表单验证规则
const rules = computed<FormRules>(() => ({
{{- range .FormFields }}
{{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant")) }}
{{- if .IsRequired }}
  {{generateJSField .ColumnName}}: [
    { required: true, message: '请{{if eq .HtmlType "select"}}选择{{else}}输入{{end}}{{.ColumnComment}}', trigger: '{{if eq .HtmlType "select"}}change{{else}}blur{{end}}' }
{{- if eq .GoType "string" }}
    {{if .ColumnType}}, { max: {{.ColumnType}}, message: '长度不能超过 {{.ColumnType}} 个字符', trigger: 'blur' }{{end}}
{{- end }}
{{- if eq .HtmlType "email" }}
    , { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
{{- end }}
{{- if eq .HtmlType "url" }}
    , { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
{{- end }}
{{- if and (eq .HtmlType "input") (contains .ColumnName "phone") }}
    , { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
{{- end }}
  ],
{{- else }}
{{- if or (eq .HtmlType "email") (eq .HtmlType "url") (contains .ColumnName "phone") }}
  {{generateJSField .ColumnName}}: [
{{- if eq .HtmlType "email" }}
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
{{- else if eq .HtmlType "url" }}
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
{{- else if contains .ColumnName "phone" }}
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
{{- end }}
  ],
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}))

// 监听表单数据变化
watch(
  () => props.formData,
  (newData) => {
    if (newData) {
      Object.assign(form.value, {
{{- range .FormFields }}
{{- if not (or (eq .ColumnName "tenant_id") (contains .ColumnName "tenant")) }}
        {{generateJSField .ColumnName}}: newData.{{generateJSField .ColumnName}} {{if eq .GoType "string"}}|| ''{{else if eq .GoType "bool"}}|| false{{else if contains .GoType "int"}}|| 0{{else if eq .GoType "time.Time"}}|| null{{else}}|| null{{end}},
{{- end }}
{{- end }}
      })
    }
  },
  { immediate: true, deep: true }
)

// 提交表单
const handleSubmit = async () => {
  try {
    const valid = await formRef.value?.validate()
    if (!valid) return

    loading.value = true

    const submitData = { ...form.value }

    if (isEdit.value) {
      await {{toLower .BusinessName}}Api.update{{.ClassName}}(props.formData.{{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}}!, submitData)
      ElMessage.success('{{.FunctionName}}信息更新成功')
    } else {
      await {{toLower .BusinessName}}Api.create{{.ClassName}}(submitData)
      ElMessage.success('{{.FunctionName}}创建成功')
    }

    emit('success')
    handleClose()
  } catch (error: any) {
    console.error(isEdit.value ? '更新{{.FunctionName}}失败:' : '创建{{.FunctionName}}失败:', error)
  } finally {
    loading.value = false
  }
}

// 关闭弹窗
const handleClose = () => {
  formRef.value?.resetFields()
  emit('update:visible', false)
}
</script>

<style lang="scss" scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 表单分组标题优化
:deep(.el-divider) {
  margin: 24px 0 16px 0;

  .el-divider__text {
    color: #409eff;
    font-weight: 500;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 6px;
  }
}

// 输入框前缀图标
:deep(.el-input__prefix) {
  color: #909399;
}

// 单选框和复选框样式优化
:deep(.el-radio) {
  margin-right: 24px;

  .el-radio__label {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

:deep(.el-checkbox) {
  margin-right: 24px;
}
</style>