// ----------------------
// fmt.ks
// 基本的な入出力セット
// Koromosoft (c) 2021
// ----------------------

import "systemcall.ks"

// 標準出力に改行付きでメッセージ表示
func Pringln(msg string){
    Print(msg+"\n")
}
