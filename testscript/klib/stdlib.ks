// KScript2 標準ライブラリ
// 
import "unityutil.ks"

// メッセージウィンドウにメッセージ表示
func ShowMsg(msg string){
    __syscall[0](msg)
}

// 名前ウィンドウにテキスト表示
func ShowName(name string){
    __syscall[1](name)
}

// ボイスの再生
func PlayVoice(id string){
    __syscall[2](id)
}

// 標準出力にメッセージ表示
func Print(msg string){
    __syscall[3](msg)
}

// 標準出力に改行付きでメッセージ表示
func Pringln(msg string){
    __syscall[3](msg+"\n")
}
