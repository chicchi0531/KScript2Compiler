package ast

import(
	"ks2/compiler/vm"
)

type ForStatement struct{
	init vm.INode
	expr vm.INode
	iter vm.IStatement
	block vm.IStatement
	lineno int
	driver *vm.Driver
}

func MakeForStatement(init vm.INode, expr vm.INode, iter vm.IStatement, block vm.IStatement, lineno int, driver *vm.Driver) *ForStatement{
	s := new(ForStatement)
	s.init = init
	s.expr = expr
	s.iter = iter
	s.block = block
	s.lineno =lineno
	s.driver = driver
	return s
}

// <init>
// l1:
// <expr>
// jze l2
// <block>
// <iter>
// jmp l1
// l2:
func (s *ForStatement) Analyze(){
	l1 := s.driver.MakeLabel()
	l2 := s.driver.MakeLabel()

	// <init>
	s.init.Push()
	// l1:
	s.driver.SetLabel(l1)
	// <expr>
	s.expr.Push()
	// jze l2
	s.driver.OpJze(l2)
	// <block>
	s.block.Analyze()
	// <iter>
	s.iter.Analyze()
	// jmp l1
	s.driver.OpJmp(l1)
	// l2:
	s.driver.SetLabel(l2)
}