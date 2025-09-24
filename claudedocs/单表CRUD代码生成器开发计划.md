# 单表CRUD代码生成器开发计划

## 项目概述

基于现有light-stack项目架构，开发一个自动化代码生成工具，能够根据数据库表结构生成完整的单表CRUD功能，包括后端API、前端页面和菜单配置。

### 技术栈分析
- **后端**: Go + Gin + GORM + 租户系统
- **前端**: Vue3 + Element Plus + TypeScript
- **数据库**: MySQL/PostgreSQL
- **架构**: RESTful API + 权限控制系统

## 1. 系统架构设计

### 1.1 代码生成器结构
```
tools/codegen/
├── cmd/
│   └── main.go              # CLI入口程序
├── internal/
│   ├── config/
│   │   └── config.go        # 配置解析
│   ├── database/
│   │   ├── analyzer.go      # 数据库表结构分析
│   │   └── schema.go        # 表结构定义
│   ├── generator/
│   │   ├── backend.go       # 后端代码生成
│   │   ├── frontend.go      # 前端代码生成
│   │   └── menu.go          # 菜单SQL生成
│   └── template/
│       ├── model.go.tmpl    # Go模型模板
│       ├── service.go.tmpl  # 服务层模板
│       ├── controller.go.tmpl # 控制器模板
│       ├── list.vue.tmpl    # 列表页面模板
│       ├── form.vue.tmpl    # 表单组件模板
│       ├── api.ts.tmpl      # API接口模板
│       └── menu.sql.tmpl    # 菜单配置SQL模板
├── config/
│   └── config.yaml          # 配置文件示例
└── README.md
```

### 1.2 核心组件设计

#### 1.2.1 配置系统 (config.go)
- 数据库连接配置
- 项目路径配置  
- 表字段映射配置
- 模板文件路径配置
- 菜单权限配置

#### 1.2.2 数据库分析器 (analyzer.go)
- 连接数据库获取表结构
- 解析字段类型和约束
- 生成Go结构体映射
- 识别主键和索引

#### 1.2.3 代码生成引擎
- **后端生成器**: 生成Model、Service、Controller
- **前端生成器**: 生成Vue页面和API接口
- **菜单生成器**: 生成菜单配置SQL

## 2. 配置文件设计

### 2.1 主配置文件 (config.yaml)
```yaml
# 数据库配置
database:
  driver: mysql
  dsn: "user:pass@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

# 项目配置
project:
  name: "light-stack" 
  module: "github.com/LiteMove/light-stack"
  base_path: "../.."

# 文件路径配置
paths:
  model: "internal/model"
  service: "internal/service"
  controller: "internal/controller"
  frontend: "web/src/views"
  api: "web/src/api"

# 模板配置
templates:
  model: "templates/model.go.tmpl"
  service: "templates/service.go.tmpl"
  controller: "templates/controller.go.tmpl"
  list_vue: "templates/list.vue.tmpl"
  form_vue: "templates/form.vue.tmpl"
  api: "templates/api.ts.tmpl"
  menu_sql: "templates/menu.sql.tmpl"

# 表配置
tables:
  - name: "products"           # 数据库表名
    model: "Product"           # Go模型名
    comment: "产品管理"         # 中文描述
    module: "product"          # 模块名(路由路径)
    tenant: true               # 支持租户
    menu:
      name: "产品管理"
      icon: "Goods"
      parent_id: 2
      sort: 10
      permissions:
        - code: "product:list"
          name: "产品列表"
        - code: "product:create"
          name: "创建产品"
        - code: "product:update" 
          name: "更新产品"
        - code: "product:delete"
          name: "删除产品"
```

### 2.2 字段映射配置
```yaml
# 字段类型映射
field_mappings:
  # 数据库类型 -> Go类型
  type_mapping:
    varchar: string
    text: string
    int: int
    bigint: int64
    tinyint: int
    decimal: decimal.Decimal
    datetime: time.Time
    date: time.Time
    
  # 表单控件映射
  form_control:
    string: input
    int: number
    decimal: number
    time.Time: date
    text: textarea
    
# 全局设置
global:
  pagination:
    default_page: 1
    default_size: 20
    max_size: 100
```

## 3. 生成的代码结构

### 3.1 后端代码结构

#### 3.1.1 模型文件 (internal/model/product.go)
```go
type Product struct {
    TenantBaseModel           // 继承租户基础模型
    Name        string        `json:"name" gorm:"type:varchar(100);not null;comment:产品名称"`
    Description string        `json:"description" gorm:"type:text;comment:产品描述"`
    Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
    Status      int           `json:"status" gorm:"type:tinyint;not null;default:1"`
}
```

#### 3.1.2 服务层 (internal/service/product_service.go)
- CreateProduct(product *model.Product) error
- GetProductList(tenantID uint64, page, pageSize int, keyword string, status int) ([]model.Product, int64, error)
- GetProduct(id uint64) (*model.Product, error)
- UpdateProduct(product *model.Product) error
- DeleteProduct(id uint64) error
- UpdateProductStatus(id uint64, status int) error

#### 3.1.3 控制器 (internal/controller/product_controller.go)
- 请求结构体定义
- 参数验证
- 租户ID处理
- CRUD操作接口
- 统一响应格式

### 3.2 前端代码结构

#### 3.2.1 列表页面 (web/src/views/system/products/index.vue)
- 页面头部(标题、操作按钮)
- 搜索表单
- 数据表格(分页、排序、操作列)
- 批量操作
- 权限控制指令

#### 3.2.2 表单组件 (web/src/views/system/products/components/ProductForm.vue)
- 表单验证
- 字段映射
- 提交处理
- 弹框模式

#### 3.2.3 API接口 (web/src/api/product.ts)
```typescript
export interface ProductListParams {
  page: number
  page_size: number
  keyword?: string
  status?: number
}

export const productApi = {
  list: (params: ProductListParams) => request.get('/api/products', { params }),
  create: (data: CreateProductRequest) => request.post('/api/products', data),
  update: (id: number, data: UpdateProductRequest) => request.put(`/api/products/${id}`, data),
  delete: (id: number) => request.delete(`/api/products/${id}`),
  detail: (id: number) => request.get(`/api/products/${id}`)
}
```

## 4. 模板系统设计

### 4.1 后端模板

#### 4.1.1 模型模板 (templates/model.go.tmpl)
```go
package model

import (
    {{- range .Imports }}
    "{{ . }}"
    {{- end }}
)

// {{ .Model }} {{ .Comment }}
type {{ .Model }} struct {
    {{- if .Tenant }}
    TenantBaseModel
    {{- else }}
    BaseModel
    {{- end }}
    {{- range .Fields }}
    {{ .GoName }} {{ .GoType }} `json:"{{ .JsonTag }}" gorm:"{{ .GormTag }}"{{ if .Validate }} validate:"{{ .Validate }}"{{ end }}`
    {{- end }}
}

// Table{{ .Model }}Name 表名
func ({{ .Model }}) TableName() string {
    return "{{ .TableName }}"
}
```

#### 4.1.2 服务模板 (templates/service.go.tmpl)
- 基础CRUD方法
- 分页查询
- 状态更新
- 批量操作
- 错误处理

#### 4.1.3 控制器模板 (templates/controller.go.tmpl)
- 请求结构体
- 路由处理方法
- 参数验证
- 权限检查
- 响应格式化

### 4.2 前端模板

#### 4.2.1 列表页模板 (templates/list.vue.tmpl)
- 响应式布局
- 搜索筛选
- 表格展示
- 分页组件
- 操作按钮
- 权限控制

#### 4.2.2 表单模板 (templates/form.vue.tmpl)
- 表单验证规则
- 字段类型适配
- 提交处理
- 错误提示

## 5. 菜单配置SQL生成

### 5.1 菜单SQL模板 (templates/menu.sql.tmpl)
```sql
-- {{ .Comment }}菜单配置
INSERT INTO `system_menus` (`name`, `path`, `component`, `icon`, `parent_id`, `sort`, `type`, `hidden`, `status`, `created_at`, `updated_at`) VALUES
('{{ .Menu.Name }}', '/system/{{ .Module }}', 'system/{{ .Module }}/index', '{{ .Menu.Icon }}', {{ .Menu.ParentId }}, {{ .Menu.Sort }}, 1, 0, 1, NOW(), NOW());

-- 获取插入的菜单ID
SET @menu_id = LAST_INSERT_ID();

-- {{ .Comment }}权限配置
INSERT INTO `system_permissions` (`code`, `name`, `menu_id`, `type`, `status`, `created_at`, `updated_at`) VALUES
{{- range $index, $perm := .Menu.Permissions }}
('{{ $perm.Code }}', '{{ $perm.Name }}', @menu_id, 1, 1, NOW(), NOW()){{ if ne $index (sub (len $.Menu.Permissions) 1) }},{{ end }}
{{- end }};
```

## 6. CLI工具设计

### 6.1 命令行接口
```bash
# 生成单个表
./codegen generate --table=products --config=config.yaml

# 生成多个表
./codegen generate --tables=products,orders,users --config=config.yaml

# 预览生成内容(不写文件)
./codegen preview --table=products --config=config.yaml

# 分析数据库结构
./codegen analyze --config=config.yaml

# 生成配置文件模板
./codegen init --config=config.yaml
```

### 6.2 参数说明
- `--table`: 指定单个表名
- `--tables`: 指定多个表名(逗号分隔)
- `--config`: 配置文件路径
- `--output`: 输出目录(可选)
- `--preview`: 预览模式(不实际生成文件)
- `--force`: 强制覆盖已存在文件

## 7. 开发阶段规划

### 阶段一：核心架构搭建 (1-2天)
1. 创建项目目录结构
2. 设计配置文件格式
3. 实现配置解析功能
4. 搭建CLI框架

### 阶段二：数据库分析器 (1天)
1. 实现数据库连接
2. 表结构解析
3. 字段类型映射
4. Go结构体生成

### 阶段三：模板系统 (2-3天)
1. 设计模板语法
2. 创建后端代码模板
3. 创建前端页面模板
4. 实现模板渲染引擎

### 阶段四：代码生成器 (2天)
1. 后端代码生成逻辑
2. 前端代码生成逻辑
3. 菜单SQL生成
4. 文件写入和目录管理

### 阶段五：测试和优化 (1天)
1. 单元测试编写
2. 功能测试验证
3. 错误处理完善
4. 性能优化

### 阶段六：文档和部署 (1天)
1. 使用说明文档
2. 配置示例
3. 常见问题解答
4. 部署脚本

## 8. 技术要点

### 8.1 关键技术
- **Go模板引擎**: text/template 实现代码模板渲染
- **数据库反射**: database/sql + GORM 获取表结构
- **YAML解析**: gopkg.in/yaml.v2 处理配置文件
- **CLI框架**: cobra 构建命令行工具
- **文件操作**: os/path 处理文件和目录

### 8.2 设计模式
- **模板方法模式**: 统一代码生成流程
- **策略模式**: 不同数据库类型适配
- **工厂模式**: 生成器对象创建
- **建造者模式**: 复杂配置对象构建

### 8.3 扩展性考虑
- 支持自定义模板
- 支持多种数据库类型
- 支持插件化扩展
- 支持国际化配置

## 9. 预期效果

### 9.1 提升效率
- 单表CRUD开发时间从2-3小时缩短到5-10分钟
- 减少重复代码编写
- 统一代码风格和结构
- 降低人工错误概率

### 9.2 代码质量
- 自动生成规范的代码结构
- 统一的错误处理机制
- 完整的权限控制集成
- 标准化的API接口

### 9.3 维护便利
- 模板化管理代码结构
- 批量更新和修改能力
- 版本控制友好
- 团队协作标准化

## 10. 风险评估

### 10.1 技术风险
- 数据库兼容性问题
- 模板复杂度控制
- 生成代码质量保证
- 性能优化需求

### 10.2 缓解措施
- 充分测试不同数据库类型
- 模块化设计降低复杂度
- 自动化测试验证生成代码
- 增量生成和缓存机制

这个开发计划提供了完整的单表CRUD代码生成器实现方案，涵盖了从架构设计到具体实现的各个方面，可以作为开发的详细指导文档。