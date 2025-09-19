package main

import (
	"log"

	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/database"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	db := database.GetDB()

	// 初始化存储库和服务
	tenantRepo := repository.NewTenantRepository(db)
	userRepo := repository.NewUserRepository(db)
	tenantService := service.NewTenantService(tenantRepo, userRepo)

	// 测试租户服务
	log.Println("开始测试租户管理功能...")

	// 测试获取租户列表
	tenants, total, err := tenantService.GetTenantList(1, 10, "", 0)
	if err != nil {
		log.Printf("获取租户列表失败: %v", err)
	} else {
		log.Printf("获取租户列表成功，共 %d 个租户", total)
		for _, tenant := range tenants {
			log.Printf("租户: ID=%d, 名称=%s, 域名=%s, 状态=%d",
				tenant.ID, tenant.Name, tenant.Domain, tenant.Status)
		}
	}

	// 测试检查域名是否存在
	exists, err := tenantService.CheckDomainExists("system")
	if err != nil {
		log.Printf("检查域名失败: %v", err)
	} else {
		log.Printf("域名 'system' 是否存在: %t", exists)
	}

	// 测试根据域名获取租户
	tenant, err := tenantService.ValidateTenant("system")
	if err != nil {
		log.Printf("验证租户失败: %v", err)
	} else {
		log.Printf("验证租户成功: ID=%d, 名称=%s", tenant.ID, tenant.Name)
	}

	log.Println("租户管理功能测试完成！")
}
