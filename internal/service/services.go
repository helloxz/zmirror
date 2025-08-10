package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"zmirror/internal/model"

	"gorm.io/gorm"
)

type RegistryService struct {
	db *gorm.DB
}

func NewRegistryService(db *gorm.DB) *RegistryService {
	return &RegistryService{db: db}
}

// GetEnabledRegistries 获取启用的镜像源，按优先级排序
func (s *RegistryService) GetEnabledRegistries() ([]model.Registry, error) {
	var registries []model.Registry
	err := s.db.Where("enabled = ?", true).Order("priority ASC").Find(&registries).Error
	return registries, err
}

// GetAllRegistries 获取所有镜像源
func (s *RegistryService) GetAllRegistries() ([]model.Registry, error) {
	var registries []model.Registry
	err := s.db.Order("priority ASC").Find(&registries).Error
	return registries, err
}

// CreateRegistry 创建镜像源
func (s *RegistryService) CreateRegistry(registry *model.Registry) error {
	return s.db.Create(registry).Error
}

// UpdateRegistry 更新镜像源
func (s *RegistryService) UpdateRegistry(registry *model.Registry) error {
	return s.db.Save(registry).Error
}

// DeleteRegistry 删除镜像源
func (s *RegistryService) DeleteRegistry(id uint) error {
	return s.db.Delete(&model.Registry{}, id).Error
}

type WhitelistService struct {
	db *gorm.DB
}

func NewWhitelistService(db *gorm.DB) *WhitelistService {
	return &WhitelistService{db: db}
}

// IsImageWhitelisted 检查镜像是否在白名单中
func (s *WhitelistService) IsImageWhitelisted(imageName string) (bool, error) {
	var whitelists []model.Whitelist
	err := s.db.Where("enabled = ?", true).Find(&whitelists).Error
	if err != nil {
		return false, err
	}

	for _, wl := range whitelists {
		if strings.HasPrefix(imageName, wl.Prefix) {
			return true, nil
		}
	}
	return false, nil
}

// GetAllWhitelists 获取所有白名单
func (s *WhitelistService) GetAllWhitelists() ([]model.Whitelist, error) {
	var whitelists []model.Whitelist
	err := s.db.Find(&whitelists).Error
	return whitelists, err
}

// CreateWhitelist 创建白名单
func (s *WhitelistService) CreateWhitelist(whitelist *model.Whitelist) error {
	return s.db.Create(whitelist).Error
}

// DeleteWhitelist 删除白名单
func (s *WhitelistService) DeleteWhitelist(id uint) error {
	return s.db.Delete(&model.Whitelist{}, id).Error
}

type UserService struct {
	db        *gorm.DB
	adminUser string
	adminPass string
}

func NewUserService(db *gorm.DB, adminUser, adminPass string) *UserService {
	return &UserService{
		db:        db,
		adminUser: adminUser,
		adminPass: adminPass,
	}
}

// HashPassword 哈希密码
func (s *UserService) HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// AuthenticateUser 验证用户
func (s *UserService) AuthenticateUser(username, password string) (*model.User, error) {
	// 首先检查是否是管理员用户（从配置文件）
	if username == s.adminUser && password == s.adminPass {
		return &model.User{
			Username: username,
			IsAdmin:  true,
		}, nil
	}

	// 从数据库查找普通用户（只查找非管理员用户）
	var user model.User
	hashedPassword := s.HashPassword(password)
	err := s.db.Where("username = ? AND password = ?", username, hashedPassword).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers 获取所有普通用户（不包括管理员）
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var users []model.User
	// 普通用户在数据库中，管理员在配置文件中，所以只查询数据库
	err := s.db.Find(&users).Error
	return users, err
}

// CreateUser 创建普通用户
func (s *UserService) CreateUser(user *model.User) error {
	user.Password = s.HashPassword(user.Password)
	user.IsAdmin = false // 确保创建的是普通用户
	return s.db.Create(user).Error
}

// DeleteUser 删除普通用户
func (s *UserService) DeleteUser(id uint) error {
	// 只删除普通用户，不能删除管理员（管理员在配置文件中）
	return s.db.Delete(&model.User{}, id).Error
}

type ProxyService struct {
	registryService *RegistryService
	client          *http.Client
}

func NewProxyService(registryService *RegistryService) *ProxyService {
	return &ProxyService{
		registryService: registryService,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyRequest 代理请求到上游镜像源
func (s *ProxyService) ProxyRequest(method, path string, headers http.Header) (*http.Response, string, error) {
	fmt.Printf("PROXY DEBUG: Starting proxy request %s %s\n", method, path)
	registries, err := s.registryService.GetEnabledRegistries()
	if err != nil {
		return nil, "", err
	}

	for _, registry := range registries {
		fmt.Printf("PROXY DEBUG: Trying registry %s\n", registry.URL)
		targetURL, err := url.JoinPath(registry.URL, path)
		if err != nil {
			continue
		}

		// 首次请求
		fmt.Printf("PROXY DEBUG: Making first request to %s\n", targetURL)
		resp, err := s.makeRequest(method, targetURL, headers, "")
		if err != nil {
			fmt.Printf("PROXY DEBUG: First request failed: %v\n", err)
			continue
		}

		fmt.Printf("PROXY DEBUG: First request status: %d\n", resp.StatusCode)

		// 如果是401且响应头包含WWW-Authenticate，尝试获取token
		if resp.StatusCode == 401 {
			authHeader := resp.Header.Get("WWW-Authenticate")
			fmt.Printf("PROXY DEBUG: Got 401, WWW-Authenticate: %s\n", authHeader)
			if authHeader != "" && strings.Contains(strings.ToLower(authHeader), "bearer") {
				resp.Body.Close()

				// 尝试获取匿名token
				fmt.Printf("PROXY DEBUG: Trying to get anonymous token\n")
				token, err := s.getAnonymousToken(authHeader, path)
				if err == nil && token != "" {
					fmt.Printf("PROXY DEBUG: Got token, making second request\n")
					// 用token重新请求
					newResp, err := s.makeRequest(method, targetURL, headers, token)
					if err == nil {
						fmt.Printf("PROXY DEBUG: Second request successful: %d\n", newResp.StatusCode)
						// 只有成功才返回，否则继续尝试下一个镜像源
						if newResp.StatusCode >= 200 && newResp.StatusCode < 300 {
							return newResp, registry.URL, nil
						}
						newResp.Body.Close()
					} else {
						fmt.Printf("PROXY DEBUG: Second request failed: %v\n", err)
					}
				} else {
					fmt.Printf("PROXY DEBUG: Failed to get token: %v\n", err)
				}
			} else {
				resp.Body.Close()
			}
			continue // 401认证失败，尝试下一个镜像源
		}

		// 如果成功（2xx状态码），返回结果
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, registry.URL, nil
		}

		// 其他状态码（包括404），关闭响应体并尝试下一个镜像源
		resp.Body.Close()
		fmt.Printf("PROXY DEBUG: Registry %s returned %d, trying next registry\n", registry.URL, resp.StatusCode)
	}

	return nil, "", fmt.Errorf("all registries failed")
}

// makeRequest 发送HTTP请求
func (s *ProxyService) makeRequest(method, targetURL string, headers http.Header, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, targetURL, nil)
	if err != nil {
		return nil, err
	}

	// 复制请求头，但排除一些不需要的头
	for name, values := range headers {
		if shouldSkipHeader(name) {
			continue
		}
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// 如果有token，添加Authorization头
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return s.client.Do(req)
}

// getAnonymousToken 获取匿名访问token
func (s *ProxyService) getAnonymousToken(authHeader, requestPath string) (string, error) {
	// 解析WWW-Authenticate头
	// 格式: Bearer realm="https://auth.docker.io/token",service="registry.docker.io",scope="repository:helloz/onenav:pull"

	var realm, service, scope string

	// 移除 "Bearer " 前缀
	authHeader = strings.TrimPrefix(authHeader, "Bearer ")

	// 简单的解析WWW-Authenticate头
	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "realm=") {
			realm = strings.Trim(strings.TrimPrefix(part, "realm="), "\"")
		} else if strings.HasPrefix(part, "service=") {
			service = strings.Trim(strings.TrimPrefix(part, "service="), "\"")
		} else if strings.HasPrefix(part, "scope=") {
			scope = strings.Trim(strings.TrimPrefix(part, "scope="), "\"")
		}
	}

	if realm == "" {
		return "", fmt.Errorf("no realm found in auth header")
	}

	// 构建token请求URL
	tokenURL := realm
	params := url.Values{}
	if service != "" {
		params.Add("service", service)
	}
	if scope != "" {
		params.Add("scope", scope)
	}

	if len(params) > 0 {
		tokenURL += "?" + params.Encode()
	}

	// 创建一个新的HTTP客户端用于token请求，避免超时问题
	tokenClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 请求token
	resp, err := tokenClient.Get(tokenURL)
	if err != nil {
		return "", fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	// 解析token响应
	var tokenResp struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %v", err)
	}

	// 返回token（优先使用token字段，如果没有则使用access_token）
	if tokenResp.Token != "" {
		return tokenResp.Token, nil
	}
	if tokenResp.AccessToken != "" {
		return tokenResp.AccessToken, nil
	}

	return "", fmt.Errorf("no token found in response")
}

// shouldSkipHeader 判断是否应该跳过某个请求头
func shouldSkipHeader(name string) bool {
	skipHeaders := []string{
		"Connection",
		"Proxy-Connection",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Authorization", // 不要传递用户的认证信息给上游
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
		"Host",
	}

	lowerName := strings.ToLower(name)
	for _, skip := range skipHeaders {
		if strings.ToLower(skip) == lowerName {
			return true
		}
	}
	return false
}

type LogService struct {
	db *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{db: db}
}

// shouldLogRequest 判断是否应该记录该请求的日志
func (s *LogService) shouldLogRequest(method, path string) bool {
	// 1. 所有HEAD请求都不记录
	if method == "HEAD" {
		return false
	}

	// 2. 路径包含"blobs"的不记录
	if strings.Contains(path, "blobs") {
		return false
	}

	// 3. 路径等于"/v2/"的不记录
	if path == "/v2/" {
		return false
	}

	// 4. 只记录GET请求且路径包含"manifests"的，但排除带有sha256的请求
	if method == "GET" && strings.Contains(path, "manifests") {
		// 如果路径包含sha256，不记录
		if strings.Contains(path, "sha256:") {
			return false
		}
		return true
	}

	// 其他情况都不记录
	return false
}

// LogAccess 记录访问日志
func (s *LogService) LogAccess(log *model.AccessLog) {
	// 检查是否应该记录该请求
	if !s.shouldLogRequest(log.Method, log.Path) {
		return
	}

	// 异步记录日志，不影响主要业务流程
	go func() {
		s.db.Create(log)
	}()
}

// GetAccessLogs 获取访问日志
func (s *LogService) GetAccessLogs(limit int) ([]model.AccessLog, error) {
	var logs []model.AccessLog
	err := s.db.Order("created_at DESC").Limit(limit).Find(&logs).Error
	return logs, err
}

// ClearAccessLogs 清空访问日志
func (s *LogService) ClearAccessLogs() error {
	return s.db.Where("1 = 1").Delete(&model.AccessLog{}).Error
}
