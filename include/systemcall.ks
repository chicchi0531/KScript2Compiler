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

// [5]標準出力にメッセージ表示
func Print(msg string){
    __syscall[5](msg)
}

// [4]SEサウンドの再生
func PlaySE(name string){
    __syscall[4](name)
}

// [5]BGMの再生
func PlayBGM(name string){
    __syscall[5](name)
}

// [6]BGMの停止
func StopBGM(name string){
    __syscall[6](name)
}

// [7]画像の表示
func ShowImg(handle int, layer int){
    __syscall[7](handle, layer)
}

// [8]画像の消去
func HideImg(handle int){
    __syscall[8](handle)
}

// [9]画像のロード
func LoadImg(name string)int{
    handle := __syscall[9](name)
    return handle
}

// [10]画像の移動
func MoveImg(handle int, posx float, posy float, duration float, sync int){
    __syscall[10](handle, posx, posy, duration, sync)
}

// [11]画像の回転
func RotateImg(handle int, rot float, duration float, sync int){
    __syscall[11](handle, rot, duration, sync)
}

// [12]画像の拡縮
func ScaleImg(handle int, sx float, sy float, duration float, sync int){
    __syscall[12](handle, sx, sy, duration, sync)
}

// [13]画像のアニメーション変化
func ChangeImgState(handle int, state string, sync int){
    __syscall[13](handle, state, sync)
}