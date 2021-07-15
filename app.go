package main

import (
	"fmt"
	"os"

	"ks2/compiler"
)

// メイン関数
func main() {

	if len(os.Args) <= 1 {
		fmt.Println("引数にスクリプトファイルを指定してください。")
		return
	}

	scriptFilePath := os.Args[1]

	compiler.Compile(scriptFilePath)
}
  