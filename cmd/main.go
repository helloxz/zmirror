package main

import (
	"fmt"
	"log"
	"os"

	"zmirror/internal/config"
	"zmirror/internal/database"
	"zmirror/internal/router"
	"zmirror/internal/service"
)

func main() {
	// 确保data目录存在
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	db, err := database.InitDatabase(cfg.Database.Path)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化服务
	userService := service.NewUserService(db, cfg.Admin.Username, cfg.Admin.Password)
	registryService := service.NewRegistryService(db)
	whitelistService := service.NewWhitelistService(db)
	logService := service.NewLogService(db)

	// 设置路由
	r := router.SetupRouter(userService, registryService, whitelistService, logService)

	// 启动服务器
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", address)
	log.Fatal(r.Run(address))
}
