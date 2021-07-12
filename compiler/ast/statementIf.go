package ast

import(
	"ks2/compiler/vm"
)

type IfStatement struct{
	condition vm.INode
	ifState vm.IStatement
	elseState vm.IStatement
	lineno int
	driver *vm.Driver
}

func MakeIfStatement(condition vm.INode, ifstate vm.IStatement, elsestate vm.IStatement, lineno int, driver *vm.Driver) *IfStatement{
	s := new(IfStatement)
	s.condition = condition
	s.ifState = ifstate
	s.elseState = elsestate
	s.driver = driver
	return s
}

// --- ifのみのとき
// <condition>
// jze l1
// <if_state>
// l1:
// --- elseありのとき
// <condition>
// jze l1
// <if_state>
// jmp l2
// l1:
// <else_state>
// l2:
func (s *IfStatement) Analyze(){
	l1 := s.driver.MakeLabel()

	// <condition>
	s.condition.Push()
	// jze l1
	s.driver.OpJze(l1)
	// <if_state>
	s.ifState.Analyze()

	if s.elseState == nil{
		// l1:
		s.driver.SetLabel(l1)
	}else{
		l2 := s.driver.MakeLabel()

		// jmp l2
		s.driver.OpJmp(l2)
		// l1:
		s.driver.SetLabel(l1)
		// <else_state>
		s.elseState.Analyze()
		// l2:
		s.driver.SetLabel(l2)
	}
}