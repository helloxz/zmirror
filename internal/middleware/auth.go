package middleware

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"zmirror/internal/model"
	"zmirror/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware Docker认证中间件
func AuthMiddleware(userService *service.UserService, whitelistService *service.WhitelistService, logService *service.LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("DEBUG: Processing request %s %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Printf("DEBUG: Authorization header: %s\n", c.Request.Header.Get("Authorization"))

		// 记录访问日志
		accessLog := &model.AccessLog{
			ClientIP:  c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			UserAgent: c.Request.Header.Get("User-Agent"),
			CreatedAt: time.Now(),
		}

		// 从路径中提取镜像名
		imageName := extractImageNameFromPath(c.Request.URL.Path)
		accessLog.ImageName = imageName
		fmt.Printf("DEBUG: Extracted image name: %s\n", imageName)

		// 检查是否在白名单中
		if imageName != "" {
			isWhitelisted, err := whitelistService.IsImageWhitelisted(imageName)
			fmt.Printf("DEBUG: Whitelist check - image: %s, whitelisted: %v, error: %v\n", imageName, isWhitelisted, err)
			if err == nil && isWhitelisted {
				fmt.Printf("DEBUG: Image is whitelisted, allowing access\n")
				accessLog.StatusCode = 200
				logService.LogAccess(accessLog)
				c.Next()
				return
			}
		}

		// 需要认证
		authHeader := c.Request.Header.Get("Authorization")
		fmt.Printf("DEBUG: Auth header: %s\n", authHeader)
		if authHeader == "" {
			fmt.Printf("DEBUG: No auth header, returning 401 with Basic realm\n")
			// 使用标准的Basic认证方式，更兼容Docker客户端
			c.Header("WWW-Authenticate", `Basic realm="Docker Registry"`)
			c.Header("Docker-Distribution-API-Version", "registry/2.0")
			accessLog.StatusCode = 401
			logService.LogAccess(accessLog)
			c.JSON(401, gin.H{"errors": []gin.H{{"code": "UNAUTHORIZED", "message": "authentication required"}}})
			c.Abort()
			return
		}

		// 解析认证信息 (只支持Basic认证)
		username, password, ok := parseBasicAuth(authHeader)
		fmt.Printf("DEBUG: Parsed auth - username: %s, password: %s, ok: %v\n", username, password, ok)
		if !ok {
			fmt.Printf("DEBUG: Failed to parse Basic auth\n")
			accessLog.StatusCode = 401
			logService.LogAccess(accessLog)
			c.Header("WWW-Authenticate", `Basic realm="Docker Registry"`)
			c.Header("Docker-Distribution-API-Version", "registry/2.0")
			c.JSON(401, gin.H{"errors": []gin.H{{"code": "UNAUTHORIZED", "message": "invalid authorization header"}}})
			c.Abort()
			return
		}

		// 验证用户
		user, err := userService.AuthenticateUser(username, password)
		fmt.Printf("DEBUG: User authentication - user: %v, error: %v\n", user, err)
		if err != nil {
			fmt.Printf("DEBUG: User authentication failed: %v\n", err)
			accessLog.StatusCode = 401
			logService.LogAccess(accessLog)
			// 确保Docker客户端重新发送认证信息
			c.Header("WWW-Authenticate", `Basic realm="Docker Registry"`)
			c.Header("Docker-Distribution-API-Version", "registry/2.0")
			c.JSON(401, gin.H{"errors": []gin.H{{"code": "UNAUTHORIZED", "message": "invalid credentials"}}})
			c.Abort()
			return
		}

		accessLog.Username = user.Username
		accessLog.StatusCode = 200
		logService.LogAccess(accessLog)

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Next()
	}
}

// AdminAuthMiddleware WEB管理界面认证中间件
func AdminAuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		username, password, ok := parseBasicAuth(authHeader)
		if !ok {
			c.JSON(401, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		user, err := userService.AuthenticateUser(username, password)
		if err != nil || !user.IsAdmin {
			c.JSON(401, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}

		c.Set("admin", user)
		c.Next()
	}
}

// CacheHeadersMiddleware 设置CDN缓存头
func CacheHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 设置缓存策略
		if shouldCache(path) {
			if isManifest(path) {
				// Manifest 文件短时间缓存
				c.Header("Cache-Control", "public, max-age=300") // 5分钟
			} else if isBlob(path) {
				// Blob 文件长时间缓存
				c.Header("Cache-Control", "public, max-age=31536000") // 1年
			}
		} else {
			// 不缓存的请求
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		}

		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// parseBasicAuth 解析Basic认证头
func parseBasicAuth(authHeader string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
		return "", "", false
	}

	encoded := authHeader[len(prefix):]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}

	credentials := string(decoded)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	return parts[0], parts[1], true
}

// extractImageNameFromPath 从请求路径中提取镜像名
func extractImageNameFromPath(path string) string {
	// Docker Registry API路径格式：
	// /v2/{name}/manifests/{reference}
	// /v2/{name}/blobs/{digest}
	if !strings.HasPrefix(path, "/v2/") {
		return ""
	}

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return ""
	}

	// 提取镜像名部分
	nameParts := []string{}
	for i := 2; i < len(parts); i++ {
		if parts[i] == "manifests" || parts[i] == "blobs" || parts[i] == "tags" {
			break
		}
		nameParts = append(nameParts, parts[i])
	}

	return strings.Join(nameParts, "/")
}

// shouldCache 判断路径是否应该被缓存
func shouldCache(path string) bool {
	return isManifest(path) || isBlob(path)
}

// isManifest 判断是否是manifest请求
func isManifest(path string) bool {
	return strings.Contains(path, "/manifests/")
}

// isBlob 判断是否是blob请求
func isBlob(path string) bool {
	return strings.Contains(path, "/blobs/")
}
