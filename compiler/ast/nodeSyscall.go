package ast

import(
	"ks2/compiler/vm"
	cm "ks2/compiler/common"
)

type NSysCall struct{
	Node
	index vm.INode
	args []vm.INode
}

func MakeSysCallNode(lineno int, index vm.INode, args []vm.INode, driver *vm.Driver) *NSysCall{
	n := new(NSysCall)
	n.Lineno = lineno
	n.index = index
	n.args = args
	n.Driver = driver
	return n
}

func (n *NSysCall) Push() *vm.VariableTag{
	// 引数逆積み
	size := 0
	for i:=len(n.args)-1; i>=0; i--{
		vt := n.args[i].Push()
		size += vt.VarType.Size * vt.ArraySize
	}
	// 引数の数積み
	n.Driver.OpPushInteger(size)

	// システムコールの番号積み
	n.index.Push()

	n.Driver.OpSysCall()

	tt := n.Driver.VariableTypeTable.GetTag(cm.TYPE_DYNAMIC)
	vt := vm.MakeVariableTag("", tt, false, 1, n.Driver)
	return vt
}