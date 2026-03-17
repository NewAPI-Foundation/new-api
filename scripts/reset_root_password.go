// 一次性脚本：将 root 用户密码设为指定值。用法: go run scripts/reset_root_password.go [新密码]
package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/QuantumNous/new-api/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load(".env")
	password := "Abq0zncpyOxdIF1W6J5r1HGULK0oUh2x"
	if len(os.Args) > 1 {
		password = os.Args[1]
	}
	dsn := os.Getenv("SQL_DSN")
	if dsn == "" {
		fmt.Println("请设置 .env 中的 SQL_DSN")
		os.Exit(1)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		os.Exit(1)
	}
	hashed, err := common.Password2Hash(password)
	if err != nil {
		fmt.Println("生成哈希失败:", err)
		os.Exit(1)
	}
	res := db.Exec("UPDATE users SET password = ? WHERE username = ?", hashed, "root")
	if res.Error != nil {
		fmt.Println("更新失败:", res.Error)
		os.Exit(1)
	}
	if res.RowsAffected == 0 {
		fmt.Println("未找到用户 root")
		os.Exit(1)
	}
	fmt.Println("已将 root 密码更新为指定值，共更新", res.RowsAffected, "条")
}
