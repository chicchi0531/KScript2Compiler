package compiler

import "strconv"

const (
	TYPE_INTEGER = iota
	TYPE_FLOAT
	TYPE_STRING
	TYPE_VOID
	TYPE_UNKNOWN //未決定（型推論用)
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
	}
	return "error:想定されない型です。t:" + strconv.Itoa(t)
}