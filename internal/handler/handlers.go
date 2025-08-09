package handler

import (
	"io"
	"strconv"
	"time"

	"zmirror/internal/model"
	"zmirror/internal/service"

	"github.com/gin-gonic/gin"
)

type RegistryHandler struct {
	proxyService    *service.ProxyService
	registryService *service.RegistryService
	logService      *service.LogService
}

func NewRegistryHandler(proxyService *service.ProxyService, registryService *service.RegistryService, logService *service.LogService) *RegistryHandler {
	return &RegistryHandler{
		proxyService:    proxyService,
		registryService: registryService,
		logService:      logService,
	}
}

// ProxyToRegistry 代理请求到上游Registry
func (h *RegistryHandler) ProxyToRegistry(c *gin.Context) {
	// 获取原始请求信息
	method := c.Request.Method
	path := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		path += "?" + c.Request.URL.RawQuery
	}

	// 代理请求
	resp, _, err := h.proxyService.ProxyRequest(method, path, c.Request.Header)
	if err != nil {
		c.JSON(500, gin.H{"errors": []gin.H{{"code": "UNKNOWN", "message": "failed to proxy request"}}})
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for name, values := range resp.Header {
		for _, value := range values {
			c.Header(name, value)
		}
	}

	// 设置状态码
	c.Status(resp.StatusCode)

	// 复制响应体
	io.Copy(c.Writer, resp.Body)

	// 记录代理日志
	if user, exists := c.Get("user"); exists {
		if u, ok := user.(*model.User); ok {
			log := &model.AccessLog{
				ClientIP:   c.ClientIP(),
				Method:     method,
				Path:       path,
				UserAgent:  c.Request.Header.Get("User-Agent"),
				StatusCode: resp.StatusCode,
				Username:   u.Username,
				CreatedAt:  time.Now(),
			}
			h.logService.LogAccess(log)
		}
	}
}

// GetVersion Docker Registry版本检查
func (h *RegistryHandler) GetVersion(c *gin.Context) {
	c.Header("Docker-Distribution-API-Version", "registry/2.0")
	c.JSON(200, gin.H{
		"name":    "Docker Registry Proxy",
		"version": "1.0.0",
	})
}

type AdminHandler struct {
	userService      *service.UserService
	registryService  *service.RegistryService
	whitelistService *service.WhitelistService
	logService       *service.LogService
}

func NewAdminHandler(userService *service.UserService, registryService *service.RegistryService, whitelistService *service.WhitelistService, logService *service.LogService) *AdminHandler {
	return &AdminHandler{
		userService:      userService,
		registryService:  registryService,
		whitelistService: whitelistService,
		logService:       logService,
	}
}

// 用户管理

func (h *AdminHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "user created successfully"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "user deleted successfully"})
}

// 镜像源管理

func (h *AdminHandler) GetRegistries(c *gin.Context) {
	registries, err := h.registryService.GetAllRegistries()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, registries)
}

func (h *AdminHandler) CreateRegistry(c *gin.Context) {
	var registry model.Registry
	if err := c.ShouldBindJSON(&registry); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.registryService.CreateRegistry(&registry); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "registry created successfully"})
}

func (h *AdminHandler) UpdateRegistry(c *gin.Context) {
	var registry model.Registry
	if err := c.ShouldBindJSON(&registry); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.registryService.UpdateRegistry(&registry); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "registry updated successfully"})
}

func (h *AdminHandler) DeleteRegistry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid registry id"})
		return
	}

	if err := h.registryService.DeleteRegistry(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "registry deleted successfully"})
}

// 白名单管理

func (h *AdminHandler) GetWhitelists(c *gin.Context) {
	whitelists, err := h.whitelistService.GetAllWhitelists()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, whitelists)
}

func (h *AdminHandler) CreateWhitelist(c *gin.Context) {
	var whitelist model.Whitelist
	if err := c.ShouldBindJSON(&whitelist); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.whitelistService.CreateWhitelist(&whitelist); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "whitelist created successfully"})
}

func (h *AdminHandler) DeleteWhitelist(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid whitelist id"})
		return
	}

	if err := h.whitelistService.DeleteWhitelist(uint(id)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "whitelist deleted successfully"})
}

// 访问日志

func (h *AdminHandler) GetAccessLogs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	logs, err := h.logService.GetAccessLogs(limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, logs)
}
