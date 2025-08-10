package router

import (
	"zmirror/controller"
	"zmirror/internal/handler"
	"zmirror/internal/middleware"
	"zmirror/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置所有路由
func SetupRouter(
	userService *service.UserService,
	registryService *service.RegistryService,
	whitelistService *service.WhitelistService,
	logService *service.LogService,
) *gin.Engine {
	// 初始化处理器
	registryHandler := handler.NewRegistryHandler(
		service.NewProxyService(registryService),
		registryService,
		logService,
	)
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

	// 设置Docker Registry API路由
	setupRegistryRoutes(router, registryHandler, userService, whitelistService, logService)

	// 设置管理API路由
	setupAdminRoutes(router, adminHandler, userService)

	return router
}

// setupRegistryRoutes 设置Docker Registry相关路由
func setupRegistryRoutes(
	router *gin.Engine,
	registryHandler *handler.RegistryHandler,
	userService *service.UserService,
	whitelistService *service.WhitelistService,
	logService *service.LogService,
) {
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
}

// setupAdminRoutes 设置管理API路由
func setupAdminRoutes(router *gin.Engine, adminHandler *handler.AdminHandler, userService *service.UserService) {
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
		api.DELETE("/logs", adminHandler.ClearAccessLogs)
	}
}
