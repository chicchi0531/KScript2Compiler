package vm

type INode interface {
	Push() int
	Pop() int
}

type IStateBlock interface{
	AddStates(IStatement) IStateBlock
	Analyze()
}

type IStatement interface{
	Analyze()
}
