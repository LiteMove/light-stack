package main

import (
	"fmt"
	"strings"

	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/utils"
	"github.com/go-playground/validator/v10"
)

// TestMenuRequest 测试用的菜单请求结构
type TestMenuRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Code string `json:"code" validate:"required_if_permission,permission_code,max=100"`
	Type string `json:"type" validate:"required,oneof=directory menu permission"`
}

func main() {
	config.Init()

	// 注册自定义验证器
	validate := validator.New()
	utils.RegisterCustomValidators(validate)

	// 测试案例
	testCases := []TestMenuRequest{
		{Name: "用户管理", Code: "", Type: "directory"},                    // 目录类型，code可以为空 - 应该通过
		{Name: "用户列表", Code: "", Type: "menu"},                         // 菜单类型，code可以为空 - 应该通过
		{Name: "创建用户", Code: "system:user:create", Type: "permission"}, // 权限类型，有权限标识 - 应该通过
		{Name: "删除用户", Code: "", Type: "permission"},                   // 权限类型，无权限标识 - 应该失败
		{Name: "编辑用户", Code: "user_edit", Type: "permission"},          // 权限类型，下划线格式 - 应该通过
		{Name: "查看用户", Code: "user-view", Type: "permission"},          // 权限类型，短横线格式 - 应该通过
		{Name: "导出用户", Code: "user:export@", Type: "permission"},       // 权限类型，无效字符 - 应该失败
	}

	fmt.Println("菜单权限标识验证测试")
	fmt.Println(strings.Repeat("=", 50))

	for i, testCase := range testCases {
		fmt.Printf("测试 %d: %+v\n", i+1, testCase)

		err := validate.Struct(testCase)
		if err != nil {
			fmt.Printf("  ❌ 验证失败: %v\n", err)
		} else {
			fmt.Printf("  ✅ 验证通过\n")
		}
		fmt.Println()
	}

	fmt.Println("测试完成")
}
