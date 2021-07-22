package vm

type INode interface {
	Push() *VariableTag
	Pop() *VariableTag
}

type IStateBlock interface{
	AddStates(IStatement) IStateBlock
	Analyze()
}

type IStatement interface{
	Analyze()
}
