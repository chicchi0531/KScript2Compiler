package ast

import (
	cm "ks2/compiler/common"
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
		defNode := MakeVarDefineNode(n.Lineno, "",
				n.Driver.VariableTypeTable.GetTag(cm.TYPE_INTEGER),
				false, 1 ,n.Driver)
		n.Driver.OpPushInteger(n.Ival)
		return defNode.Push()
	case OP_FLOAT:
		defNode := MakeVarDefineNode(n.Lineno, "",
				n.Driver.VariableTypeTable.GetTag(cm.TYPE_FLOAT),
				false, 1, n.Driver)
		n.Driver.OpPushFloat(n.Fval)
		return defNode.Push()
	case OP_STRING:
		defNode := MakeVarDefineNode(n.Lineno, "",
				n.Driver.VariableTypeTable.GetTag(cm.TYPE_STRING),
				false, 1, n.Driver)
		n.Driver.OpPushString(n.Sval)
		return defNode.Push()
	}

	panic("予期せぬエラーです。値型以外がconst nodeとしてpushされました。")

	return nil
}