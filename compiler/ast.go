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

	OP_INTEGER
	OP_FLOAT
	OP_STRING
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

	ival    int
	fval    float32
	sval    string
}

func (n *Node) Push() int {

	// 単項演算の場合
	switch n.op {
	case OP_INCR:
		t := n.left.Push()
		if t == TYPE_STRING{
			n._err(ERR_0005,"")
			return TYPE_INTEGER
		}
		n.driver.OpPushInteger(1)
		n.driver.OpAdd()
		return t

	case OP_DECR:
		t := n.left.Push()
		if t == TYPE_STRING{
			n._err(ERR_0006,"")
			return TYPE_INTEGER
		}
		n.driver.OpPushInteger(1)
		n.driver.OpSub()
		return t
		
	case OP_NOT:
		t := n.left.Push()
		if t == TYPE_STRING{
			n._err(ERR_0007, "")
			return TYPE_INTEGER
		}
		n.driver.OpPushInteger(-1)
		n.driver.OpMul()
		return t

	// const node
	case OP_INTEGER:
		n.driver.OpPushInteger(n.ival)
		return TYPE_INTEGER
	case OP_FLOAT:
		n.driver.OpPushFloat(n.fval)
		return TYPE_FLOAT
	case OP_STRING:
		n.driver.OpPushString(n.sval)
		return TYPE_STRING
	}

	// 二項演算の場合
	leftType := n.left.Push()
	rightType := n.right.Push()

	// 型チェック
	if leftType == TYPE_STRING && rightType != TYPE_STRING ||
		leftType != TYPE_STRING && rightType == TYPE_STRING{
		n._err(ERR_0012, "")
	}

	// 文字列演算
	if leftType == TYPE_STRING || rightType == TYPE_STRING{
		switch n.op{
		case OP_ADD:
			n.driver.OpAddString()
		default:
			n._err(ERR_0013, "")
		}
		return TYPE_STRING
	}

	// 数値演算
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
		n._err(ERR_0014, "")
	}

	// どちらかの項がfloatの場合は、float項にキャストする
	if leftType == TYPE_FLOAT || rightType == TYPE_FLOAT{
		return TYPE_FLOAT
	}
	return TYPE_INTEGER
}

func (n *Node) Pop() {
	n._err(ERR_0008, "")
}

// 内部エラー出力用
func (n *Node) _err(errorcode string, submsg string){
	n.driver.err.LogError(n.driver.filename, n.driver.lineno, errorcode, submsg)
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
