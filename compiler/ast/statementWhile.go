package ast

import(
	"ks2/compiler/vm"
)

type WhileStatement struct{
	expr vm.INode
	block vm.IStatement
	lineno int
	driver *vm.Driver
}

func MakeWhileStatement(expr vm.INode, block vm.IStatement, lineno int, driver *vm.Driver) *WhileStatement{
	s := new(WhileStatement)
	s.expr = expr
	s.block = block
	s.lineno = lineno
	s.driver = driver
	return s
}

// l1:
// <expr>
// jze l2
// <block>
// jmp l1
// l2:
func (s *WhileStatement) Analyze(){
	l1 := s.driver.MakeLabel()
	l2 := s.driver.MakeLabel()

	// l1:
	s.driver.SetLabel(l1)
	// <expr>
	s.expr.Push()
	// jze l2
	s.driver.OpJze(l2)
	// <block>
	s.block.Analyze()
	// jmp l1
	s.driver.OpJmp(l1)
	// l2:
	s.driver.SetLabel(l2)
}