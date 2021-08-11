
func ShowDialog(title string, thumneil int, maintext string, choise [2]string) int {
    var result int
    result = __syscall[22](title,thumneil,maintext,choise)
    return result
}

func Print(msg string){
    __syscall[0](msg)
}

func main(){
    var choise [2]string
    choise[0] = "はい"
    choise[1] = "いいえ"

    c := ShowDialog("テスト質問", 0, "ここは本文です。", choise)

    if c==0{
        Print("はいが選択されました。")
    }else if c==1{
        Print("いいえが選択されました。")
    }else{
        Print("想定しない答えです。")
    }
}