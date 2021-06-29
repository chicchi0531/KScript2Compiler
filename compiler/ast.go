package compiler

const (
	TYPE_INTEGER = iota
	TYPE_FLOAT
	TYPE_STRING
	TYPE_VOID
)

const (
	OP_EQUAL = iota
	OP_GT
	OP_GE
	OP_LT
	OP_LE
	OP_NEQU
	OP_NOT
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_AND
	OP_OR
	OP_INCR
	OP_DECR
)

type INode interface {
	Push() int
	Pop()
}

type Node struct {
	left    INode
	right   INode
	op      int
	driver  *Driver
	varType int
	ival    int
	fval    float32
	sval    string
}

func (n *Node) Push() int {
	// 単項演算の場合
	if n.varType != TYPE_STRING {

		switch n.op {
		case OP_INCR:
			if n.left.Push() == 
			n.driver.OpPushInteger(1)
			n.driver.OpAdd()
		case OP_DECR:
			n.left.Push()
			n.driver.OpPushInteger(1)
			n.driver.OpSub()
		case OP_NOT:
			n.left.Push()
			n.driver.OpPushInteger(-1)
			n.driver.OpMul()
		}
		return n.varType
	}

	// 二項演算の場合
	leftType := n.left.Push()
	rightType := n.right.Push()

	if leftType != rightType {
		panic("左辺と右辺の型が違います")
	}

	switch n.op {
	case OP_EQUAL:
		n.driver.OpEqual()
	case OP_ADD:
		n.driver.OpAdd()
	case OP_SUB:
		n.driver.OpSub()
	case OP_MUL:
		n.driver.OpMul()
	case OP_DIV:
		n.driver.OpDiv()
	case OP_GT:
		n.driver.OpGt()
	case OP_GE:
		n.driver.OpGe()
	case OP_LT:
		n.driver.OpLt()
	case OP_LE:
		n.driver.OpLe()
	case OP_AND:
		n.driver.OpAnd()
	case OP_OR:
		n.driver.OpOr()
	default:
		panic("不明な演算子が使用されました")
	}

	return TYPE_INTEGER
}

func (n *Node) Pop() {
	panic("このノードはポップできません")
}

// const node
type ConstNode struct {
	Node
}

func (n *ConstNode) Push() int {
	n.driver.OpPushInteger(n.ival)
	return TYPE_INTEGER
}

// assign node
type AssignNode struct {
	Node
}

func (n *AssignNode) Push() int {
	n.right.Push()
	n.left.Pop()

	return TYPE_INTEGER
}

// value node
type ValueNode struct {
	Node
}

func (n *ValueNode) Push() int {
	n.driver.OpPushValue(-1)
	return TYPE_INTEGER
}

func (n *ValueNode) Pop() {
	n.driver.OpPop(-1)
}
