package main

import (
	"GolangTask4/internal/api/routes"
	"GolangTask4/internal/domain"
	"GolangTask4/pkg/database"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	DB, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("Database connected successfully")

	// 创建表
	DB.AutoMigrate(&domain.User{}, &domain.Post{}, &domain.Comment{})

	// 初始化Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	port := fmt.Sprintf(":%s", 8080)
	fmt.Printf("Server is running on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
