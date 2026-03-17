// 生成密码哈希，用法: go run scripts/gen_password_hash.go [明文密码]
package main

import (
	"fmt"
	"os"

	"github.com/QuantumNous/new-api/common"
)

func main() {
	password := "Abq0zncpyOxdIF1W6J5r1HGULK0oUh2x"
	if len(os.Args) > 1 {
		password = os.Args[1]
	}
	hash, err := common.Password2Hash(password)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(hash)
}
