package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Registry 镜像源模型
type Registry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `gorm:"not null" json:"url"`
	Priority  int       `gorm:"default:0" json:"priority"` // 越小优先级越高
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Whitelist 白名单模型
type Whitelist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Prefix    string    `gorm:"uniqueIndex;not null" json:"prefix"` // 镜像名前缀
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccessLog 访问日志模型
type AccessLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ClientIP   string    `json:"client_ip"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	UserAgent  string    `json:"user_agent"`
	StatusCode int       `json:"status_code"`
	Username   string    `json:"username,omitempty"`
	ImageName  string    `json:"image_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// InitDatabase 初始化数据库
func InitDatabase(db *gorm.DB) error {
	// 自动迁移表结构
	err := db.AutoMigrate(&User{}, &Registry{}, &Whitelist{}, &AccessLog{})
	if err != nil {
		return err
	}

	// 创建默认镜像源
	var count int64
	db.Model(&Registry{}).Count(&count)
	if count == 0 {
		defaultRegistries := []Registry{
			{URL: "https://registry-1.docker.io", Priority: 1, Enabled: true},
		}
		for _, registry := range defaultRegistries {
			db.Create(&registry)
		}
	}

	return nil
}
