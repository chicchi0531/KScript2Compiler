package ast

import(
	"ks2/compiler/vm"
)

type ExprStatement struct{
	expr vm.INode
	driver *vm.Driver
}

func MakeExprStatement(expr vm.INode, driver *vm.Driver)*ExprStatement{
	s := new(ExprStatement)
	s.expr = expr
	s.driver = driver
	return s
}

func (s *ExprStatement) Analyze(){
	s.expr.Push()
	s.driver.RemoveLastProg()//余計に詰んだPushを消しておく
	s.driver.RemoveLastProg()
}
