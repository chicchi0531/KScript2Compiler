package vm

type INode interface {
	Push() int
	Pop() int
}

type IStateBlock interface{
	AddStates(INode) IStateBlock
	Analyze()
}
