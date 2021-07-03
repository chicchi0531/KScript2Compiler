package compiler

const (
	TYPE_INTEGER = iota
	TYPE_FLOAT
	TYPE_STRING
	TYPE_VOID
	TYPE_UNKNOWN //未決定（型推論用)
)

const (
	OP_EQUAL = iota
	OP_GT
	OP_GE
	OP_LT
	OP_LE
	OP_NEQ
	OP_NOT
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
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
	Pop() int
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
		n.left.Pop()//一旦代入してから再度PUSH
		n.left.Push()
		return t

	case OP_DECR:
		t := n.left.Push()
		if t == TYPE_STRING{
			n._err(ERR_0006,"")
			return TYPE_INTEGER
		}
		n.driver.OpPushInteger(1)
		n.driver.OpSub()
		n.left.Pop()//一旦代入してから再度PUSH
		n.left.Push()
		return t
		
	case OP_NOT:
		t := n.left.Push()
		if t == TYPE_STRING{
			n._err(ERR_0007, "")
			return TYPE_INTEGER
		}
		n.driver.OpNot()
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
	case OP_NEQ:
		n.driver.OpNequ()
	case OP_ADD:
		n.driver.OpAdd()
	case OP_SUB:
		n.driver.OpSub()
	case OP_MUL:
		n.driver.OpMul()
	case OP_DIV:
		n.driver.OpDiv()
	case OP_MOD:
		n.driver.OpMod()
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

func (n *Node) Pop() int {
	n._err(ERR_0008, "")
	return -1
}

// 内部エラー出力用
func (n *Node) _err(errorcode string, submsg string){
	n.driver.err.LogError(n.driver.filename, n.driver.lineno, errorcode, submsg)
}
func (n *Node) _warning(warningcode string, submsg string){
	n.driver.err.LogWarning(n.driver.filename, n.driver.lineno, warningcode, submsg)
}

// assign node
type AssignNode struct {
	Node
}

func (n *AssignNode) Push() int {
	rightType := n.right.Push()
	leftType := n.left.Pop()

	// 型チェック
	if leftType == TYPE_UNKNOWN{
		// 左辺が型未推定の場合は、推定して仕込む
		// 直前のポップ命令から変数番号を取得
		varIndex := n.driver.program[len(n.driver.program)-1].value
		varTag := n.driver.variableTable.GetTag(varIndex)
		varTag.varType = rightType //右辺の型をそのまま左辺の型とする
		return rightType

	}else if rightType == TYPE_STRING || leftType == TYPE_STRING{
		if rightType != leftType{
			n._err(ERR_0017, "")
			return -1
		}else{
			return TYPE_STRING
		}
	}

	// 数値型の場合は、代入先の型をそのまま返す
	// ダウンキャストになる場合は警告を出しておく
	if leftType == TYPE_INTEGER && rightType == TYPE_FLOAT{
		n._warning(WARNING_0001,"float->int")
	}
	return leftType
}

// value node
type ValueNode struct {
	Node
	index int
}

func MakeValueNode(name string, driver *Driver) *ValueNode{
	index := driver.variableTable.FindVariable(name)
	if index == -1{
		driver.err.LogError(driver.filename, driver.lineno, ERR_0016, "不明な識別子："+name)
	}
	return &ValueNode{Node:Node{ driver:driver }, index:index }
}

func (n *ValueNode) Push() int {
	n.driver.OpPushValue(n.index)
	tag := n.driver.variableTable.GetTag(n.index)
	return tag.varType
}

func (n *ValueNode) Pop() int {
	n.driver.OpPop(n.index)
	tag := n.driver.variableTable.GetTag(n.index)
	return tag.varType
}
