// img_test.ks
// 画像関係の表示テストプログラム

import "systemcall.ks"

// ノベル用短縮命令
func ST(handle int){
    ShowImg(handle, 0.5)
}
func FA(handle int, pose int, cloth int, face string){
    SetPose(handle,pose)
    SetFace(handle,face)
    SetCloth(handle,cloth)
    Wait(0.2)
}

func main(){
    var version string = "ver.1.0"

    hTest := NewImg("st_000_test")

    #novel
    @ShowWindow
    @FA hTest 0 1 ""
    @ST hTest
    - 女性
    こんにちは。これはノベルテスト%version%です」
    @FA hTest 0 1 "通常2"
    このノベルスクリプトシステムは、
    独自のスクリプト言語で制御します。」
    @FA hTest 1 1 ""
    Unity上で動作していて、
    既存のUnityプロジェクトに簡単に組み込めます。」
    @FA hTest 1 1 "笑顔2"
    立ち絵はこのように、ポーズや表情や衣装を、
    自由に切り替えられます。」
    @FA hTest 0 1 "笑顔"
    文章の途中でも、
    @FA hTest 0 1 "怒り"
    こんなかんじで、
    @FA hTest 1 1 "にやり2"
    簡単に表情を切り替えられます。」
    @FA hTest 0 1 "ほほえみ"
    @HideWindow
    #

    HideImg(hTest, 1.0)
    DeleteImg(hTest)
}