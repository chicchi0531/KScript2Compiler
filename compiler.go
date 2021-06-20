package main

import(
	"os"
)

// メイン関数
func main() {
	if len(os.Args) <= 1 { return }
	lexer := &Lexer{src: os.Args[1], index:0}
	yyParse(lexer)
	println("計算結果:", lexer.result)
}
  