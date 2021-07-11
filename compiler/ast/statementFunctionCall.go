package ast

import(
	"ks2/compiler/vm"
)

type FunctionCallStatement struct{
	node vm.INode
	driver *vm.Driver
}

func MakeFunctionCallStatement(node vm.INode, driver *vm.Driver) *FunctionCallStatement{
	s := new(FunctionCallStatement)
	s.node = node
	s.driver = driver
	return s
}

func (s *FunctionCallStatement) Analyze(){
	s.node.Push()
	s.driver.OpPop() //戻り値をポップしておく
}