package ast

import(
	"ks2/compiler/vm"
)

type NVariableDefine struct{
	Node
	Name string
	VarType *vm.VariableTypeTag
	IsPointer bool
	ArraySize int
}

func MakeVarDefineNode(lineno int, name string, vartype *vm.VariableTypeTag, isPointer bool, arraysize int, driver *vm.Driver)*NVariableDefine{
	n := new(NVariableDefine)
	n.Lineno = lineno
	n.Name = name
	n.VarType = vartype
	n.IsPointer = isPointer
	n.ArraySize = arraysize
	n.Driver = driver
	return n
}

func MakeVarDefineNodeWithAssign(lineno int, name string, vartype *vm.VariableTypeTag, expr vm.INode, driver *vm.Driver)*NVariableDefine{
	n := new(NVariableDefine)
	n.Lineno = lineno
	n.Name = name
	n.VarType = vartype
	n.Right = expr
	n.IsPointer = false
	n.ArraySize = 1
	n.Driver = driver
	return n
}

func (n *NVariableDefine) Push() *vm.VariableTag{
	// 初期値代入があるかどうか
	var index int
	if n.Right != nil{
		index = n.Driver.VariableTable.DefineValue(n.Lineno, n.Name, n.VarType, n.IsPointer, n.ArraySize)
		varNode := MakeValueNode(n.Lineno, n.Name, n.Driver)
		assignNode := MakeAssignAsInit(n.Lineno, varNode, n.Right, OP_ASSIGN, index, n.Driver)
		assignNode.Push()
	}else{
		index = n.Driver.VariableTable.DefineValue(n.Lineno, n.Name, n.VarType, n.IsPointer, n.ArraySize)	
	}
	return n.Driver.VariableTable.GetTag(index)
}
