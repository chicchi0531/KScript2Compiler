package ast

import(
	"ks2/compiler/vm"
	cm "ks2/compiler/common"
)

type ContinueStatement struct{
	lineno int
	driver *vm.Driver
}

func MakeContinueStatement(lineno int, driver *vm.Driver) *ContinueStatement{
	return &ContinueStatement{lineno:lineno, driver:driver}
}

func (s *ContinueStatement) Analyze(){
	// error check
	if s.driver.ContinueLabel == -1{
		s.driver.Err.LogError(s.driver.Filename, s.lineno, cm.ERR_0030, "")
		return
	}

	s.driver.OpJmp(s.driver.ContinueLabel)
}
