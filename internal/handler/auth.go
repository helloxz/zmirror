package handler

import (
	"encoding/base64"
	"strings"
	"time"

	"zmirror/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userService *service.UserService
	jwtSecret   []byte
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   []byte("your-secret-key"), // 在生产环境中应该从配置读取
	}
}

// TokenResponse JWT token响应
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// GetToken 获取Bearer token
func (h *AuthHandler) GetToken(c *gin.Context) {
	// 获取service和scope参数
	service := c.Query("service")
	scope := c.Query("scope")

	// 检查Basic认证
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.Header("WWW-Authenticate", `Basic realm="Docker Registry"`)
		c.JSON(401, gin.H{"error": "authentication required"})
		return
	}

	username, password, ok := parseBasicAuth(authHeader)
	if !ok {
		c.Header("WWW-Authenticate", `Basic realm="Docker Registry"`)
		c.JSON(401, gin.H{"error": "invalid authorization header"})
		return
	}

	// 验证用户
	user, err := h.userService.AuthenticateUser(username, password)
	if err != nil {
		c.Header("WWW-Authenticate", `Basic realm="Docker Registry"`)
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	// 创建JWT claims
	claims := jwt.MapClaims{
		"iss":    "localhost:8080",                 // issuer
		"sub":    user.Username,                    // subject (username)
		"aud":    service,                          // audience (service)
		"exp":    time.Now().Add(time.Hour).Unix(), // expires in 1 hour
		"nbf":    time.Now().Unix(),                // not before
		"iat":    time.Now().Unix(),                // issued at
		"jti":    "random-jwt-id",                  // JWT ID
		"access": []map[string]interface{}{},       // access permissions
	}

	// 如果有scope，解析并添加到access中
	if scope != "" {
		access := parseScope(scope)
		claims["access"] = access
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	// 返回token
	c.JSON(200, TokenResponse{
		Token:     tokenString,
		ExpiresIn: 3600, // 1 hour
	})
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

// parseScope 解析Docker Registry scope
func parseScope(scope string) []map[string]interface{} {
	// scope格式: "repository:helloz/zurl:pull"
	access := []map[string]interface{}{}

	parts := strings.SplitN(scope, ":", 3)
	if len(parts) == 3 && parts[0] == "repository" {
		actions := strings.Split(parts[2], ",")
		access = append(access, map[string]interface{}{
			"type":    parts[0],
			"name":    parts[1],
			"actions": actions,
		})
	}

	return access
}
