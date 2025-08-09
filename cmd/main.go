package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"zmirror/internal/handler"
	"zmirror/internal/middleware"
	"zmirror/internal/model"
	"zmirror/internal/service"

	"zmirror/controller"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`

	Admin struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"admin"`

	Database struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"database"`
}

func main() {
	// 确保data目录存在
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	db, err := initDatabase(config.Database.Path)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化服务
	userService := service.NewUserService(db, config.Admin.Username, config.Admin.Password)
	registryService := service.NewRegistryService(db)
	whitelistService := service.NewWhitelistService(db)
	proxyService := service.NewProxyService(registryService)
	logService := service.NewLogService(db)

	// 初始化处理器
	registryHandler := handler.NewRegistryHandler(proxyService, registryService, logService)
	adminHandler := handler.NewAdminHandler(userService, registryService, whitelistService, logService)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 添加中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.CacheHeadersMiddleware())

	// 静态文件服务
	router.LoadHTMLGlob("static/html/*")
	router.Static("/assets", "./static/assets")
	router.GET("/", controller.Home)

	// Docker Registry API路由 - 使用NoRoute处理
	// 使用NoRoute处理所有v2请求，包括/v2/
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 处理v2路径（包括/v2/和其他v2路径）
		if path == "/v2/" || path == "/v2" || (len(path) > 4 && path[:4] == "/v2/") {
			// 应用认证中间件
			middleware.AuthMiddleware(userService, whitelistService, logService)(c)
			if c.IsAborted() {
				return
			}

			// 如果是版本检查路径，调用GetVersion
			if path == "/v2/" || path == "/v2" {
				registryHandler.GetVersion(c)
			} else {
				registryHandler.ProxyToRegistry(c)
			}
			return
		}

		// 其他路径返回404
		c.JSON(404, gin.H{"error": "not found"})
	})

	// 管理API路由
	api := router.Group("/api")
	api.Use(middleware.AdminAuthMiddleware(userService))
	{
		// 用户管理
		api.GET("/users", adminHandler.GetUsers)
		api.POST("/users", adminHandler.CreateUser)
		api.DELETE("/users/:id", adminHandler.DeleteUser)

		// 镜像源管理
		api.GET("/registries", adminHandler.GetRegistries)
		api.POST("/registries", adminHandler.CreateRegistry)
		api.PUT("/registries", adminHandler.UpdateRegistry)
		api.DELETE("/registries/:id", adminHandler.DeleteRegistry)

		// 白名单管理
		api.GET("/whitelists", adminHandler.GetWhitelists)
		api.POST("/whitelists", adminHandler.CreateWhitelist)
		api.DELETE("/whitelists/:id", adminHandler.DeleteWhitelist)

		// 访问日志
		api.GET("/logs", adminHandler.GetAccessLogs)
	}

	// 启动服务器
	address := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Server starting on %s", address)
	log.Fatal(router.Run(address))
}

func loadConfig() (*Config, error) {
	configPath := "./data/config.toml"

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return nil, err
		}
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func createDefaultConfig(configPath string) error {
	defaultConfig := `[server]
host = "0.0.0.0"
port = "8080"

[admin]
username = "admin"
password = "admin123"

[database]
path = "./data/registry.db"
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}

func initDatabase(dbPath string) (*gorm.DB, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	// 连接数据库并启用WAL模式
	db, err := gorm.Open(sqlite.Open(dbPath+"?_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 初始化数据库表结构
	if err := model.InitDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}
