// ----------------------
// novel.ks
// 基本的なノベル関数セット
// Koromosoft (c) 2021
// ----------------------

import "vector.ks"
import "systemcall.ks"

// 立ち絵表示
// @param handle 立ち絵のインスタンスハンドル
// @param motion 立ち絵のアニメーションID
// @param layer 表示するレイヤー番号
func ST(handle int, motion string, layer int) {
    ChangeImgState(handle, motion)
    ShowImg(handle, layer)
}

// SE再生
func SE() {

}

// BGM再生
func BGM() {

}

// 背景表示
func BG() {

}