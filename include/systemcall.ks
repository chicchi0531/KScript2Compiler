// ----------------------
// systemcall.ks
// systemcallのパース
// Koromosoft (c) 2021
// ----------------------

// [0]メッセージウィンドウにメッセージ表示
func ShowMsg(msg string){
    __syscall[0](msg)
}

// [1]名前ウィンドウにテキスト表示
func ShowName(name string){
    __syscall[1](name)
}

// [2]ボイスの再生
func PlayVoice(id string){
    __syscall[2](id)
}

// [3]標準出力にメッセージ表示
func Print(msg string){
    __syscall[3](msg)
}

// [4]
