package ast

import(
	"ks2/compiler/vm"
)

type CompoundStatement struct{
	statements vm.IStateBlock
}

func MakeCompoundStatement(statement vm.IStateBlock) *CompoundStatement{
	s := new(CompoundStatement)
	s.statements = statement
	return s
}

func (s *CompoundStatement) Analyze(){
	s.statements.Analyze()
}