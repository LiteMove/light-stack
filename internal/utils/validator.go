package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// PermissionCodeValidator 权限码验证器
func PermissionCodeValidator(fl validator.FieldLevel) bool {
	code := fl.Field().String()

	// 权限码格式: 支持字母、数字、冒号、下划线、短横线
	// 例如: system:user:list, user:create, admin-panel:view
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9:_-]+$`, code)
	return matched
}

// RequiredIfPermissionType 当类型为permission时权限标识必填
func RequiredIfPermissionType(fl validator.FieldLevel) bool {
	// 获取同一结构体中的 Type 字段
	parent := fl.Parent()
	typeField := parent.FieldByName("Type")

	if !typeField.IsValid() {
		return true // 如果找不到Type字段，则通过验证
	}

	menuType := typeField.String()
	code := fl.Field().String()

	// 如果类型是 permission，则 code 必须不为空
	if menuType == "permission" {
		return code != ""
	}

	// 其他类型时，code 可以为空
	return true
}

// RegisterCustomValidators 注册自定义验证器
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("permission_code", PermissionCodeValidator)
	v.RegisterValidation("required_if_permission", RequiredIfPermissionType)
}
