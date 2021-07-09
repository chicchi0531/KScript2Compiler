package vm

type Argument struct{
	Name string
	VarType int
}

type ArgList []Argument
func (a *ArgList) Add(arg *Argument) *ArgList{
	*a = append(*a, *arg)
	return a
}

type FunctionTag struct{
	Name string
	
}

type FunctionTable struct{
	
}

