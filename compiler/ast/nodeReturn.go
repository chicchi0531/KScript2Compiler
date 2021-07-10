package ast

import(
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

type NodeReturn struct{
	Node
}

func MakeNodeReturn(value vm.INode, lineno int, driver *vm.Driver)*NodeReturn{
	n := new(NodeReturn)
	n.Left = value
	n.Lineno = lineno
	n.Driver = driver
	return n
}

func (n *NodeReturn) Push() int{
	// 戻り値のpush
	retType := cm.TYPE_VOID
	if n.Left != nil{
		retType = n.Left.Push()
	}

	// 戻り値の型チェック
	if retType != n.Driver.CurrentRetType{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0024, "")
	}

	n.Driver.OpReturn()

	return retType
}