package ast

import(
	"ks2/compiler/vm"
)

type AssignStatement struct{
	assignNode vm.INode
}

func MakeAssignStatement(assign vm.INode) *AssignStatement{
	s := new(AssignStatement)
	s.assignNode = assign
	return s
}

func (s *AssignStatement) Analyze(){
	s.assignNode.Push()
}
