package ast

import(
	"ks2/compiler/vm"
	cm "ks2/compiler/common"
)

type BreakStatement struct{
	lineno int
	driver *vm.Driver
}

func MakeBreakStatement(lineno int, driver *vm.Driver) *BreakStatement{
	return &BreakStatement{lineno:lineno, driver:driver}
}

func (s *BreakStatement) Analyze(){
	// error check
	if s.driver.BreakLabel == -1{
		s.driver.Err.LogError(s.driver.Filename, s.lineno, cm.ERR_0029, "")
		return
	}

	s.driver.OpJmp(s.driver.BreakLabel)
}
