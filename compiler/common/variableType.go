package common

import "strconv"

const (
	TYPE_UNKNOWN = iota //未決定（型推論用)
	TYPE_INTEGER
	TYPE_FLOAT
	TYPE_STRING
	TYPE_VOID
	TYPE_DYNAMIC //ランタイムで決定（syscall用。通常ルーチンでは使えない）
	TYPE_STRUCT = 0x0010 //この値以降は構造体の識別番号
)

func TYPE_TOSTR(t int) string {
	switch t {
	case TYPE_INTEGER:
		return "int"
	case TYPE_FLOAT:
		return "float"
	case TYPE_STRING:
		return "string"
	case TYPE_VOID:
		return "void"
	case TYPE_UNKNOWN:
		return "error:unknown。型推論に失敗しています。このエラーが出た場合はKoromosoftに連絡してください。"
	case TYPE_DYNAMIC:
		return "error:dynamic。ダイナミック型はローカル変数の型としては使用できません。このエラーが出た場合はKoromosoftに連絡してください。"
	}
	return "error:想定されない型です。t:" + strconv.Itoa(t)
}