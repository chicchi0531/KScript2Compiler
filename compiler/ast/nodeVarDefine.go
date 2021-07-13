package ast

import(
	"ks2/compiler/vm"
)

type NVariableDefine struct{
	Node
	Name string
	VarType int
	IsPointer bool
	Size int
}

func MakeVarDefineNode(lineno int, name string, vartype int, isPointer bool, size int, driver *vm.Driver)*NVariableDefine{
	n := new(NVariableDefine)
	n.Lineno = lineno
	n.Name = name
	n.VarType = vartype
	n.IsPointer = isPointer
	n.Size = size
	n.Driver = driver
	return n
}

func MakeVarDefineNodeWithAssign(lineno int, name string, vartype int, expr vm.INode, driver *vm.Driver)*NVariableDefine{
	n := new(NVariableDefine)
	n.Lineno = lineno
	n.Name = name
	n.VarType = vartype
	n.Right = expr
	n.IsPointer = false
	n.Size = 1
	n.Driver = driver
	return n
}

func (n *NVariableDefine) Push() int{
	// 初期値代入があるかどうか
	if n.Right != nil{
		index := n.Driver.VariableTable.DefineValue(n.Lineno, n.Name, n.VarType, n.IsPointer, n.Size)
		varNode := &NValue{Name:n.Name, Node:Node{Driver:n.Driver}}
		assignNode := &Assign{Node:Node{Left:varNode, Right:n.Right, Driver:n.Driver}}
		assignNode.Push()
		return index		
	}else{
		return n.Driver.VariableTable.DefineValue(n.Lineno, n.Name, n.VarType, n.IsPointer, n.Size)
	}
}
