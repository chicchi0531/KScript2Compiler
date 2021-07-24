package ast

import (
	"ks2/compiler/vm"
)

type NConst struct{
	Node

	Ival int
	Fval float32
	Sval string
}

// make const node
func MakeIvalNode(lineno int, value int, driver *vm.Driver) *NConst {
	n := new(NConst)
	n.Lineno = lineno
	n.Ival = value
	n.Op = OP_INTEGER
	n.Driver = driver
	return n
}
func MakeFvalNode(lineno int, value float32, driver *vm.Driver) *NConst {
	n := new(NConst)
	n.Lineno = lineno
	n.Fval = value
	n.Op = OP_FLOAT
	n.Driver = driver
	return n
}
func MakeSvalNode(lineno int, value string, driver *vm.Driver) *NConst {
	n := new(NConst)
	n.Lineno = lineno
	n.Sval = value
	n.Op = OP_STRING
	n.Driver = driver
	return n
}

func (n *NConst) Push() *vm.VariableTag{
	// const node
	switch n.Op{
	case OP_INTEGER:
		n.Driver.OpPushInteger(n.Ival)
		return vm.MakeVariableTag("", n.Driver.GetType("int", n.Lineno), false, 1, n.Driver)
		
	case OP_FLOAT:
		n.Driver.OpPushFloat(n.Fval)
		return vm.MakeVariableTag("", n.Driver.GetType("float", n.Lineno), false, 1, n.Driver)
		
	case OP_STRING:
		n.Driver.OpPushString(n.Sval)
		return vm.MakeVariableTag("", n.Driver.GetType("string", n.Lineno), false, 1, n.Driver)
		
	}

	panic("予期せぬエラーです。値型以外がconst nodeとしてpushされました。")
}