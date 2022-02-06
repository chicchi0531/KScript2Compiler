// ----------------------
// systemcall.ks
// systemcallのラップ
// Koromosoft (c) 2021
// ----------------------


func Print(msg string){
    __syscall[0](msg)
}

func Scan(){
    __syscall[1]()
}

func Assert(msg string){
    __syscall[2](msg)
}

func Exit(){
    __syscall[3]()
}

func Wait(time float){
    __syscall[4](time)
}

func WaitEndOfFrame(){
    __syscall[5]()
}

func Await(){
    __syscall[6]()
}

func ShowMsg(msg string){
    var caps [1]string
    __syscall[16](msg, caps)
}

func ShowName(name string){
    __syscall[17](name)
}

func ShowWindow(){
    __syscall[18]()
}

func HideWindow(){
    __syscall[19]()
}

func Choice1(title string, thumneil int, msg string, choise [1]string){
    __syscall[22](title, thumneil, msg, choise)
}
func Choice2(title string, thumneil int, msg string, choise [2]string){
    __syscall[22](title, thumneil, msg, choise)
}
func Choice3(title string, thumneil int, msg string, choise [3]string){
    __syscall[22](title, thumneil, msg, choise)
}

func NewImg(name string) int{
    var handle int
    handle = __syscall[32](name)
    return handle
}

func DeleteImg(handle int){
    __syscall[33](handle)
}

func ShowImg(handle int, duration float, layer int, asyncflag int){
    __syscall[34](handle, duration, layer, asyncflag)
}

func HideImg(handle int, duration float, asyncflag int){
    __syscall[35](handle, duration, asyncflag)
}

func SetState(handle int, poseid int, clothid string, faceid string){
    __syscall[36](handle, poseid, clothid, faceid)
}

func SetPosImg(handle int, dx int, dy int, duration float, syncflag int){
    __syscall[37](handle, dx, dy, duration, syncflag)
}

func RotateImg(handle int, rot float, duration float, syncflag int){
    __syscall[38](handle, rot, duration, syncflag)
}

func ScaleImg(handle int, sx float, sy float, duration float, syncflag int){
    __syscall[39](handle, sx, sy, duration, syncflag)
}

func BlurImg(handle int, radius float, duration float, syncflag int){
    __syscall[40](handle, radius, duration, syncflag)
}

func ChangeStillState(handle int, id int, duration float, syncflag int){
    __syscall[41](handle, id, duration, syncflag)
}

func ShakeScreen(duration float, strength int, vibrato int, randomness int, syncflag int){
    __syscall[42](duration, strength, vibrato, randomness, syncflag)
}

func Fade(r int, g int, b int, duration float, syncflag int){
    __syscall[43](r, g, b, duration, syncflag)
}

func AddFilter(r int, g int, b int, duration float, syncflag int){
    __syscall[44](r, g, b, duration, syncflag)
}

func MulFilter(r int, g int, b int, duration float, syncflag int){
    __syscall[45](r, g, b, duration, syncflag)
}

func ScreenFilter(r int, g int, b int, duration float, syncflag int){
    __syscall[46](r, g, b, duration, syncflag)
}

func ClearFilter(duration float, syncflag int){
    __syscall[47](duration, syncflag)
}

func TransitionIn(r int, g int, b int, transType int, duration float, syncflag int){
    __syscall[48](r, g, b, transType, duration, syncflag)
}

func TransitionOut(r int, g int, b int, transType int, duration float, syncflag int){
    __syscall[49](r, g, b, duration, syncflag)
}

func BlackBoard(alpha float, duration float, syncflag int){
    __syscall[50](alpha, duration, syncflag)
}

func PlayBGM(id int){
    __syscall[64](id)
}

func StopBGM(){
    __syscall[65]()
}

func PlaySE(id int){
    __syscall[66](id)
}

func PlayVoice(id int){
    __syscall[67](id)
}
