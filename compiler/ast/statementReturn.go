package ast

import(
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

type ReturnStatement struct{
	lineno int
	driver *vm.Driver
	expr vm.INode
}

func MakeReturnStatement(value vm.INode, lineno int, driver *vm.Driver)*ReturnStatement{
	n := new(ReturnStatement)
	n.expr = value
	n.lineno = lineno
	n.driver = driver
	return n
}

func (n *ReturnStatement) Analyze(){
	// 戻り値のpush
	var retType *vm.VariableTag
	if n.expr != nil{
		retType = n.expr.Push()
		n.driver.OpReturnValue()
	} else {
		tt := n.driver.VariableTypeTable.GetTag(cm.TYPE_VOID)
		retType = vm.MakeVariableTag("", tt, false, 1, n.driver)
		n.driver.OpReturn()
	}

	// 戻り値の型チェック
	if retType != n.driver.CurrentRetType{
		n.driver.Err.LogError(n.driver.Filename, n.lineno, cm.ERR_0024, "")
	}
}