package database

import (
	"os"
	"path/filepath"

	"zmirror/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// InitDatabase 初始化数据库连接
func InitDatabase(dbPath string) (*gorm.DB, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 获取底层的SQL数据库连接并启用WAL模式
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 显式设置WAL模式
	if _, err := sqlDB.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}

	// 设置其他SQLite优化参数
	if _, err := sqlDB.Exec("PRAGMA synchronous=NORMAL;"); err != nil {
		return nil, err
	}
	if _, err := sqlDB.Exec("PRAGMA cache_size=1000;"); err != nil {
		return nil, err
	}
	if _, err := sqlDB.Exec("PRAGMA temp_store=memory;"); err != nil {
		return nil, err
	}

	// 初始化数据库表结构
	if err := model.InitDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}
