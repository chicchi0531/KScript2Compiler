package ast

type AssignStatement struct{
	assignNode *Assign
}

func MakeAssignStatement(assign *Assign)*AssignStatement{
	s := new(AssignStatement)
	s.assignNode = assign
	return s
}

func (s *AssignStatement) Analyze(){
	s.assignNode.Push()
}
