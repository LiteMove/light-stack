// 数据库表信息
export interface TableInfo {
  tableName: string
  tableComment: string
  createTime: string
  updateTime: string
  columns?: TableColumn[]
}

// 数据库字段信息
export interface TableColumn {
  columnName: string
  columnType: string
  columnComment: string
  isNullable: string
  columnDefault: string
  columnKey: string
  extra: string
  ordinalPosition: number

  // 转换后的信息
  goType: string
  goField: string
  isPk: boolean
  isIncrement: boolean
  isRequired: boolean
  isInsert: boolean
  isEdit: boolean
  isList: boolean
  isQuery: boolean
  queryType: string
  htmlType: string
  dictType: string
}

// 代码生成表配置
export interface GenTableConfig {
  id: number
  tableName: string
  tableComment: string
  businessName: string
  moduleName: string
  functionName: string
  className: string
  packageName: string
  author: string
  parentMenuId?: number
  menuName: string
  menuUrl: string
  menuIcon: string
  permissions: string[]
  options: OptionConfig
  remark: string
  createdAt: string
  updatedAt: string
  createdBy?: number
  updatedBy?: number
  columns: GenTableColumn[]
}

// 代码生成字段配置
export interface GenTableColumn {
  id: number
  tableConfigId: number
  columnName: string
  columnComment: string
  columnType: string
  goType: string
  goField: string
  isPk: boolean
  isIncrement: boolean
  isRequired: boolean
  isInsert: boolean
  isEdit: boolean
  isList: boolean
  isQuery: boolean
  queryType: string
  htmlType: string
  dictType: string
  sort: number
  createdAt: string
  updatedAt: string
}

// 配置选项
export interface OptionConfig {
  genPath?: string
  genType?: string
  tplType?: string
  treeCode?: string
  treeParent?: string
  treeName?: string
}

// 生成历史记录
export interface GenHistory {
  id: number
  tableConfigId: number
  tableName: string
  businessName: string
  generateType: string
  fileCount: number
  fileSize: number
  downloadCount: number
  status: string
  errorMessage: string
  filePath: string
  remark: string
  createdAt: string
  createdBy?: number
}

// 系统菜单
export interface SystemMenu {
  id: number
  parentId: number
  name: string
  code: string
  type: string
  path: string
  component: string
  icon: string
  sortOrder: number
  isHidden: boolean
  status: number
  children?: SystemMenu[]
}

// 生成请求
export interface GenerateRequest {
  configId: number
  generateType: 'all' | 'backend' | 'frontend'
}

// 生成响应
export interface GenerateResponse {
  taskId: string
  message: string
  fileCount: number
  fileSize: number
}

// 查询类型枚举
export enum QueryType {
  EQ = 'EQ',      // 等于
  NE = 'NE',      // 不等于
  GT = 'GT',      // 大于
  GTE = 'GTE',    // 大于等于
  LT = 'LT',      // 小于
  LTE = 'LTE',    // 小于等于
  LIKE = 'LIKE',  // 模糊查询
  BETWEEN = 'BETWEEN' // 范围查询
}

// HTML类型枚举
export enum HtmlType {
  INPUT = 'input',         // 输入框
  TEXTAREA = 'textarea',   // 文本域
  SELECT = 'select',       // 下拉框
  RADIO = 'radio',         // 单选框
  CHECKBOX = 'checkbox',   // 复选框
  DATETIME = 'datetime',   // 日期时间
  UPLOAD = 'upload'        // 文件上传
}

// 生成类型枚举
export enum GenerateType {
  ALL = 'all',           // 全部
  BACKEND = 'backend',   // 后端
  FRONTEND = 'frontend'  // 前端
}

// 生成状态枚举
export enum GenerateStatus {
  SUCCESS = 'success',       // 成功
  FAILED = 'failed',         // 失败
  PROCESSING = 'processing'  // 处理中
}