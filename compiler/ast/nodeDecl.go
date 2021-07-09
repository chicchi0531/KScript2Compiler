package ast

type NodeDecl struct{
	Node
	Name string
	VarType int
}

func (n *NodeDecl) Push() int{
	// 初期値代入があるかどうか
	if n.Right != nil{
		index := n.Driver.VariableTable.DefineInLocal(n.Lineno, n.Name, n.VarType)
		varNode := &NodeValue{Name:n.Name, Node:Node{Driver:n.Driver}}
		assignNode := &NodeAssign{Node:Node{Left:varNode, Right:n.Right, Driver:n.Driver}}
		assignNode.Push()
		return index		
	}else{
		return n.Driver.VariableTable.DefineInLocal(n.Lineno, n.Name, n.VarType)
	}
}
