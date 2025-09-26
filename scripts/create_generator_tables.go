package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LiteMove/light-stack/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	cfg := config.Get()

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 读取SQL文件
	sqlContent, err := os.ReadFile("doc/generator_schema.sql")
	if err != nil {
		log.Fatalf("读取SQL文件失败: %v", err)
	}

	// 分割SQL语句
	statements := strings.Split(string(sqlContent), ";")

	// 执行SQL语句
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		if err := db.Exec(statement).Error; err != nil {
			log.Printf("执行SQL语句失败: %v", err)
			log.Printf("SQL语句: %s", statement)
		} else {
			log.Printf("成功执行SQL语句")
		}
	}

	fmt.Println("数据库表创建完成!")
}
