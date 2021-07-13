package ast

import (
	"ks2/compiler/vm"
	cm "ks2/compiler/common"
)

type FallThroughStatement struct{
	driver *vm.Driver
	lineno int
}

func MakeFallThroughStatement(lineno int, driver *vm.Driver) *FallThroughStatement{
	return &FallThroughStatement{lineno:lineno, driver:driver}
}

func (s *FallThroughStatement) Analyze(){
	// error check
	if s.driver.FallThroughLabel == -1{
		s.driver.Err.LogError(s.driver.Filename, s.lineno, cm.ERR_0032, "")
		return
	}

	s.driver.OpJmp(s.driver.FallThroughLabel)
}