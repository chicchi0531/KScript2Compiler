package ast

import(
	"ks2/compiler/vm"
)

type VarDefineStatement struct{
	node vm.INode
}

func MakeVarDefineStatement(node vm.INode) *VarDefineStatement{
	s := new(VarDefineStatement)
	s.node = node
	return s
}

func (s *VarDefineStatement) Analyze(){
	s.node.Push()
}
