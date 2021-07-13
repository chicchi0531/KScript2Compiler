package ast

import (
	"ks2/compiler/vm"
)

type CaseStatement struct{
	expr vm.INode
	block vm.IStatement
	lineno int
	driver *vm.Driver
}

func MakeCaseStatement(expr vm.INode, block vm.IStatement, lineno int, driver *vm.Driver) *CaseStatement{
	s := new(CaseStatement)
	s.expr = expr
	s.block = block
	s.lineno = lineno
	s.driver = driver
	return s
}

func (s *CaseStatement) Analyze(){
	s.block.Analyze()
	s.driver.OpJmp(s.driver.BreakLabel)
}

type SwitchStatement struct{
	expr vm.INode
	caseStates []*CaseStatement
	defaultState vm.IStatement
	lineno int
	driver *vm.Driver
}

func MakeSwitchStatement(expr vm.INode, caseStates []*CaseStatement, defaultState vm.IStatement, lineno int, driver *vm.Driver) *SwitchStatement{
	s := new(SwitchStatement)
	s.expr = expr
	s.caseStates = caseStates
	s.defaultState = defaultState
	s.lineno = lineno
	s.driver = driver
	return s
}

func (s *SwitchStatement) Analyze(){
	// ラベルの作成
	lEnd := s.driver.MakeLabel()//終端ラベル
	labels := make([]int,0)
	for i:=0; i<len(s.caseStates); i++{
		labels = append(labels, s.driver.MakeLabel())
	}
	lDefault := lEnd
	if s.defaultState != nil{
		lDefault = s.driver.MakeLabel()
	}

	// 分岐処理
	for i:=0; i<len(s.caseStates); i++{
		s.expr.Push()
		s.caseStates[i].expr.Push()
		s.driver.OpEqual()
		s.driver.OpJnz(labels[i])
	}
	s.driver.OpJmp(lDefault)

	// case本体の記述
	old_fallthrough := s.driver.FallThroughLabel
	old_break := s.driver.BreakLabel
	s.driver.BreakLabel = lEnd
	for i:=0; i<len(s.caseStates); i++{
		// fallthrough先の設定
		if i+1 >= len(labels){
			s.driver.FallThroughLabel = lDefault
		}else{
			s.driver.FallThroughLabel = labels[i+1]
		}

		s.driver.SetLabel(labels[i])
		s.caseStates[i].Analyze()
	}
	s.driver.FallThroughLabel = old_fallthrough
	
	// defaultの記述
	if s.defaultState != nil{
		s.driver.SetLabel(lDefault)
		s.defaultState.Analyze()
	}
	s.driver.BreakLabel = old_break

	s.driver.SetLabel(lEnd)
}

