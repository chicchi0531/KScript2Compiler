// ----------------------
// systemcall.ks
// systemcallのパース
// Koromosoft (c) 2021
// ----------------------

// [0]Print
// 標準出力
// @param msg 出力文字
func Print(msg string) {
    __syscall[0](msg)
}

// [1]Scan
// 標準入力
// @return 入力文字
func Scan() string {
    var result string = __syscall[1]()
    return result
}

// [2]Assert
// アサーション
// @param msg エラーメッセージ
func Assert(msg string) {
    __syscall[2](msg)
}

// [3]Exit
// プログラム終了
// @param code 終了コード
func Exit(code int) {
    __syscall[3](code)
}

// [4]WaitTime
// 指定秒数のウェイト
// @paran time 待つ時間[ms]
func WaitTime(time int) {
    __syscall[4](time)
}

// [5]WaitForEndOfFrame
// フレームの終わりまでウェイト
func WaitForEndOfFrame() {
    __syscall[5]()
}

// [6]Await
// sync系システムコールの終了を一括で待つ
func Await() {
    __syscall[6]()
}

// [7]GetMemAsInt
// ゲームメモリからの値の取得
func GetMemAsInt() int {
    
}