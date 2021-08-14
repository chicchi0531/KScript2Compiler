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

func SetState(handle int, poseid int, clothid int, faceid string){
    __syscall[36](handle, poseid, clothid, faceid)
}

func MoveImg(handle int, dx int, dy int, duration float, syncflag int){
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