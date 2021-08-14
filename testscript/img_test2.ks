/* KScript2 System ⒸKoromosoft 2021.
https://github.com/chicchi0531/KScript2Compiler.git
*/

import "systemcall.ks"

// 立ち絵関係
func ST(handle int, layer int){
    ShowImg(handle, 0.5, layer, 0)
}

func STAsync(handle int, layer int){
    ShowImg(handle, 0.5, layer, 1)
}

func STZoomAsync(handle int, layer int){
    ScaleImg(handle, 1.8, 1.8, 0.5, 1)
    ShowImg(handle, 0.5, layer, 1)
}

func ClrAsync(handle int){
    HideImg(handle, 0.5, 1)
}

func Face(handle int, pose int, cloth int, face string){
    SetState(handle, pose, cloth, face)
    Wait(0.2)
}

// 背景関係
func BG(handle int, layer int){
    ShowImg(handle, 0.5, layer, 0)
}

func BGAsync(handle int, layer int){
    ShowImg(handle, 0.5, layer, 1)
}

func BGZoomAsync(handle int, layer int){
    // 初期位置セット
    ScaleImg(handle, 1.3, 1.3, 0.0, 0)
    MoveImg(handle, -200, -200, 0.0, 0)
    // 移動処理
    ScaleImg(handle, 1.5, 1.5, 0.5, 1)
    MoveImg(handle, 200, 200, 0.5, 1)
}

func main (){
    少女 := NewImg("Stand/st_000_test")
    少女拡大 := NewImg("Stand/st_000_test")
    背景 := NewImg("BG/bg_000")
    背景拡大 := NewImg("BG/bg_000")

    #novel
    @ShowWindow

    @BGAsync 背景 0    
    @Face 少女 0 1 ""
    @STAsync 少女 1
    @Await
    - nil
    部屋に入ると、少女がいた。」
    @Face 少女 0 1 "笑顔"
    少女はこちらに気づくと、こちらに近づいてきた。」

    @BGZoomAsync 背景拡大 2
    @STZoomAsync 少女拡大 3
    @Await
    - 少女
    こんにちは！」

    @ClrAsync 背景
    @ClrAsync 背景拡大
    @ClrAsync 少女
    @ClrAsync 少女拡大
    @Await

    @HideWindow

    #

    DeleteImg(少女)
    DeleteImg(少女拡大)
    DeleteImg(背景)
    DeleteImg(背景拡大)
}