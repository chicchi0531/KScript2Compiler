package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chicchi0531/KScript2Compiler_golang/compiler"
)

// メイン関数
func main() {

	if len(os.Args) <= 1 {
		fmt.Println("引数にスクリプトファイルを指定してください。")
		return
	}

	scriptFilePath := os.Args[1]
	_, filename := filepath.Split(scriptFilePath)

	println(" スクリプトロード")
	scriptText := compiler.OpenScriptFile(scriptFilePath)

	println(" スクリプトのトランスパイル")
	scriptText, err := compiler.Transpile(scriptText)
	if err != nil{
		panic(err)
	}
	if _, err = os.Stat("obj"); os.IsNotExist(err){
		if err:=os.Mkdir("obj", os.ModePerm); err!=nil{
			panic(err)
		}
	}
	ioutil.WriteFile("obj/_transpiled_" + filename, []byte(scriptText), os.ModePerm)
	println(" 完了")

	println(" スクリプトのコンパイル")
	compiler.Parse(filename, scriptText)
	
	// 結果出力
	if compiler.GetWarningCount() > 0{
		fmt.Printf("%d件の警告が発生しました。", compiler.GetWarningCount())
	}
	if compiler.GetErrorCount() > 0{
		fmt.Printf("コンパイルに失敗しました。%d件のエラーが発生しました。", compiler.GetErrorCount())
	}else{
		fmt.Printf("コンパイルに成功しました。")
	}
}
  