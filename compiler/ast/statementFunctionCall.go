package ast

import (
	"ks2/compiler/vm"
)

type FunctionCallStatement struct {
	node   vm.INode
	driver *vm.Driver
}

func MakeFunctionCallStatement(node vm.INode, driver *vm.Driver) *FunctionCallStatement {
	s := new(FunctionCallStatement)
	s.node = node
	s.driver = driver
	return s
}

func (s *FunctionCallStatement) Analyze() {
	retType := s.node.Push()

	if retType != nil && !retType.VarType.IsDynamic() {
		s.driver.OpPushInteger(retType.ArraySize * retType.VarType.Size)
		s.driver.OpPop() //戻り値をポップしておく
	}
}
