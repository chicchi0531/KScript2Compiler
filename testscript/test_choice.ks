func SetChoiceTitle(title string){

}
func SetChoiceFravor(fravor string){

}
func SetChoiceImage(imgpath string){

}
func AddChoiceOption(name string){

}
func ShowChoice()int{
    return 1
}

func main(){
    // 結果を入れるタイプ
    #choice{
        @name "時間を進めよう "
        @fravor "このゲームは、リアルタイムで時間が経過します。左下の再生ボタンを押して、時間を進めてみましょう。もう一度ボタンを押すことで、時間を止められます。"
        @image "tutorial/image_001"
        @option "わかった"
        @result choiceResult
    }
    // 結果無視のタイプ
    #choice{
        @name "時間を進めよう"
        @fravor "このゲームは、リアルタイムで時間が経過します。左下の再生ボタンを押して、時間を進めてみましょう。もう一度ボタンを押すことで、時間を止められます。"
        @image "tutorial/image_001"
        @option "なるほど"
    }
}