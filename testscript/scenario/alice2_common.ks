/** @lisence KScript2 System ⒸKoromosoft 2021.
* @ks_version 2.0
* K2Alice2で使用する共通ルーチンライブラリ
* @scripter Chicchi
*/

import "systemcall.ks"

// 立ち絵関係
func ST(handle int, layer int){
    ShowImg(handle, 0.5, layer, 0)
}

func STAsync(handle int, pos int, layer int){
    SetPosImg(handle, 300*pos, 0, 0.0, 1)
    ShowImg(handle, 0.5, layer, 1)
}

func STZoomAsync(handle int, layer int){
    // 初期位置セット
    SetPosImg(handle, 0, -100, 0.0, 0)
    ScaleImg(handle, 1.8, 1.8, 0.0, 0)

    ShowImg(handle, 0.5, layer, 1)
}

func ZoomChangeST(bg int, ch1 int, ch2 int, dir int){
    d := 0.5
    // カメラ移動
    SetPosImg(bg, 100*dir, 0, d, 1)
    SetPosImg(ch1, 500*dir, -100, d, 1)
    HideImg(ch1, d, 1)

    // 新しいキャラ表示
    STZoomAsync(ch2, 3)
    SetPosImg(ch2, -300*dir, -100, 0.0, 0)
    SetPosImg(ch2, 0, -100, d, 1)
    Await()
}

func ClrAsync(handle int){
    HideImg(handle, 0.5, 1)
}

func Face(handle int, pose int, cloth string, face string){
    SetState(handle, pose, cloth, face)
}

// 背景関係
func BG(handle int, layer int){
    ShowImg(handle, 0.5, layer, 0)
}

func BGAsync(handle int, layer int){
    ShowImg(handle, 0.5, layer, 1)
}

func BGZoomAsync(handle int, layer int){
    d := 0.5
    // 初期位置セット
    ScaleImg(handle, 1.3, 1.3, 0.0, 0)
    SetPosImg(handle, -200, -200, 0.0, 0)
    // 移動処理
    ScaleImg(handle, 1.5, 1.5, d, 1)
    SetPosImg(handle, 0, 0, d, 1)
    BlurImg(handle, 20.0, d, 1)
    ShowImg(handle, d, layer, 1)
}

// 下地関係
func BB_In(){
    BlackBoard(1.0, 0.2, 1)
}
func BB_Trans_In(){
    BlackBoard(0.7, 0.2, 1)
}
func BB_Out(){
    BlackBoard(0.0, 0.2, 0)
}

// スチル関係
func Still(handle int, id int, layer int){
    ChangeStillState(handle, id, 0.0, 0)
    ShowImg(handle, 0.5, layer, 0)
}
func ChangeStill(handle int, id int){
    ChangeStillState(handle, id, 0.5, 0)
}

// フィルタ系
func WhiteFade(){
    Fade(255,255,255, 0.5, 0)
}
func BlackFade(){
    Fade(0,0,0, 0.5, 0)
}
func ClrFilter(){
    ClearFilter(0.5, 0)
}
func ClrFilter_3s(){
    ClearFilter(3.0, 0)
}
func WhiteOut(){
    AddFilter(255,255,255, 1.0, 0)
}

//------------
// 背景定義
//------------
var P_BG_000 string
var P_BG_001 string

func InitBG(){
    P_BG_000 = "BG/bg_000_0"
    P_BG_001 = "BG/bg_001_0"
}

//------------
// スチル定義
//------------
var P_STILL_000 string
var P_STILL_001 string
var P_STILL_002 string

func InitStill(){
    P_STILL_000 = "Still/still_000"
    P_STILL_001 = "Still/still_001"
    P_STILL_002 = "Still/still_002"
}

//------------
// 立ち絵定義
//------------
var P_ST_000_Alice string
var P_ST_001_Firo string
var P_ST_007_Wreath string

func InitStand(){
    P_ST_000_Alice = "Stand/st_000_alice"
    P_ST_001_Firo = "Stand/st_001_firo"
    P_ST_007_Wreath = "Stand/st_007_wreath"
}

// 初期化
func Alice2_Init(){
    InitBG()
    InitStill()
    InitStand()
}

// 終了
func Alice2_Shut(){

}