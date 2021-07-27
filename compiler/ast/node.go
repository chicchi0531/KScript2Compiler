package ast

import (
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
	"strconv"
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

type Node struct {
	Lineno int
	Left   vm.INode
	Right  vm.INode
	Op     int
	Driver *vm.Driver
}

func MakeExprNode(lineno int, left vm.INode, right vm.INode, op int, driver *vm.Driver) *Node {
	n := new(Node)
	n.Lineno = lineno
	n.Left = left
	n.Right = right
	n.Op = op
	n.Driver = driver
	return n
}

func (n *Node) Push() *vm.VariableTag {

	// 単項演算の場合
	switch n.Op {
	case OP_INCR:
		t := n.Left.Push()
		if t == nil { return vm.MakeErrTag(n.Driver) }

		if !t.VarType.IsIncrementable {
			n._err(cm.ERR_0005, "")
			return vm.MakeErrTag(n.Driver)
		}
		n.Driver.OpIncr()
		n.Left.Pop()
		n.Left.Push()
		return t

	case OP_DECR:
		t := n.Left.Push()
		if t == nil { return vm.MakeErrTag(n.Driver) }

		if !t.VarType.IsIncrementable {
			n._err(cm.ERR_0006, "")
			return vm.MakeErrTag(n.Driver)
		}
		n.Driver.OpDecr()
		n.Left.Pop()
		n.Left.Push()
		return t

	case OP_NOT:
		t := n.Left.Push()
		if t == nil { return vm.MakeErrTag(n.Driver) }

		if !t.VarType.IsNotable {
			n._err(cm.ERR_0007, "")
			return vm.MakeErrTag(n.Driver)
		}
		n.Driver.OpNot()
		return t
	}

	// 二項演算の場合
	leftType := n.Left.Push()
	rightType := n.Right.Push()
	if leftType == nil { return vm.MakeErrTag(n.Driver) }
	if rightType == nil { return vm.MakeErrTag(n.Driver) }

	// 型チェック
	// ダイナミック型は演算に使えない
	if (leftType.VarType.IsDynamic() || rightType.VarType.IsDynamic()){
		n._err(cm.ERR_0039, "")
		return vm.MakeErrTag(n.Driver)
	}
	// 型が一致しない場合は演算できない
	if !leftType.TypeCompare(rightType) {
		n._err(cm.ERR_0012,
			"[" + strconv.Itoa(leftType.ArraySize) + "]" +
			leftType.VarType.TypeName + ":" +
			"[" + strconv.Itoa(rightType.ArraySize) + "]" +
			rightType.VarType.TypeName)
			return vm.MakeErrTag(n.Driver)
	}
	// 構造体型同士は演算できない
	if n.Driver.VariableTypeTable.IsStruct(leftType.VarType) {
		n._err(cm.ERR_0042, "")
		return vm.MakeErrTag(n.Driver)
	}
	// 配列同士は演算できない
	if leftType.ArraySize >= 2{
		n._err(cm.ERR_0043, "")
		return vm.MakeErrTag(n.Driver)
	}

	// 文字列演算
	if leftType.VarType.IsString(){
		switch n.Op {
		case OP_ADD: n.Driver.OpAddString()
		default: n._err(cm.ERR_0013, "")
			return vm.MakeErrTag(n.Driver)
		}
		return leftType
	}

	// 小数点演算
	if leftType.VarType.IsFloat(){
		retType := leftType
		intType := n.Driver.VariableTypeTable.GetTag(cm.TYPE_INTEGER)
		switch n.Op {
		case OP_ADD: n.Driver.OpFAdd()
		case OP_SUB: n.Driver.OpFSub()
		case OP_MUL: n.Driver.OpFMul()
		case OP_DIV: n.Driver.OpFDiv()
		//比較演算の場合、結果はint型にキャストされる
		case OP_GT:	n.Driver.OpFGt(); retType.VarType = intType
		case OP_GE:	n.Driver.OpFGe(); retType.VarType = intType
		case OP_LT:	n.Driver.OpFLt(); retType.VarType = intType
		case OP_LE:	n.Driver.OpFLe(); retType.VarType = intType
		default : n._err(cm.ERR_0044, "")
			return vm.MakeErrTag(n.Driver)
		}
		return retType
	}

	// 数値演算
	switch n.Op {
	case OP_EQUAL: n.Driver.OpEqual()
	case OP_NEQ: n.Driver.OpNequ()
	case OP_ADD: n.Driver.OpAdd()
	case OP_SUB: n.Driver.OpSub()
	case OP_MUL: n.Driver.OpMul()
	case OP_DIV: n.Driver.OpDiv()
	case OP_MOD: n.Driver.OpMod()
	case OP_GT: n.Driver.OpGt()
	case OP_GE: n.Driver.OpGe()
	case OP_LT:	n.Driver.OpLt()
	case OP_LE:	n.Driver.OpLe()
	case OP_AND: n.Driver.OpAnd()
	case OP_OR:	n.Driver.OpOr()
	default: n._err(cm.ERR_0014, "")
	}
	return leftType
}

func (n *Node) Pop() *vm.VariableTag {
	n._err(cm.ERR_0008, "")
	return vm.MakeErrTag(n.Driver)
}

// 内部エラー出力用
func (n *Node) _err(errorcode string, submsg string) {
	n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, errorcode, submsg)
}
func (n *Node) _warning(warningcode string, submsg string) {
	n.Driver.Err.LogWarning(n.Driver.Filename, n.Lineno, warningcode, submsg)
}
