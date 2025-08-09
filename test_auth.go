package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"zmirror/internal/model"
	"zmirror/internal/service"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 创建临时测试数据库
	tmpDir, err := os.MkdirTemp("", "zmirror_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 初始化数据库
	if err := model.InitDatabase(db); err != nil {
		log.Fatal(err)
	}

	// 创建用户服务，管理员账号为 admin/test123
	userService := service.NewUserService(db, "admin", "test123")

	// 测试1: 管理员认证
	fmt.Println("=== 测试管理员认证 ===")
	adminUser, err := userService.AuthenticateUser("admin", "test123")
	if err != nil {
		fmt.Printf("管理员认证失败: %v\n", err)
	} else {
		fmt.Printf("管理员认证成功: 用户名=%s, 是否管理员=%v\n", adminUser.Username, adminUser.IsAdmin)
	}

	// 测试2: 创建普通用户
	fmt.Println("\n=== 测试创建普通用户 ===")
	normalUser := &model.User{
		Username: "testuser",
		Password: "testpass",
	}
	if err := userService.CreateUser(normalUser); err != nil {
		fmt.Printf("创建普通用户失败: %v\n", err)
	} else {
		fmt.Printf("创建普通用户成功: ID=%d, 用户名=%s, 是否管理员=%v\n", normalUser.ID, normalUser.Username, normalUser.IsAdmin)
	}

	// 测试3: 普通用户认证
	fmt.Println("\n=== 测试普通用户认证 ===")
	authUser, err := userService.AuthenticateUser("testuser", "testpass")
	if err != nil {
		fmt.Printf("普通用户认证失败: %v\n", err)
	} else {
		fmt.Printf("普通用户认证成功: 用户名=%s, 是否管理员=%v\n", authUser.Username, authUser.IsAdmin)
	}

	// 测试4: 错误密码认证
	fmt.Println("\n=== 测试错误密码认证 ===")
	_, err = userService.AuthenticateUser("testuser", "wrongpass")
	if err != nil {
		fmt.Printf("错误密码认证失败(符合预期): %v\n", err)
	} else {
		fmt.Println("错误密码认证成功(不符合预期)")
	}

	// 测试5: 获取所有用户
	fmt.Println("\n=== 测试获取所有用户 ===")
	users, err := userService.GetAllUsers()
	if err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
	} else {
		fmt.Printf("用户总数: %d\n", len(users))
		for _, user := range users {
			fmt.Printf("  用户: ID=%d, 用户名=%s, 是否管理员=%v\n", user.ID, user.Username, user.IsAdmin)
		}
	}

	fmt.Println("\n测试完成!")
}
