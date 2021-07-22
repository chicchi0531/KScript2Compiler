// ----------------------
// stdlib.ks
// 基本的な機能セット
// Koromosoft (c) 2021
// ----------------------

import "systemcall.ks"

// 標準出力に改行付きでメッセージ表示
func Pringln(msg string){
    Print(msg+"\n")
}
