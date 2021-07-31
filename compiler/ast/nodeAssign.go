package ast

import (
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
	"strconv"
)

const(
	OP_ASSIGN = iota
	OP_ADD_ASSIGN
	OP_SUB_ASSIGN
	OP_MUL_ASSIGN
	OP_DIV_ASSIGN
	OP_MOD_ASSIGN
)

// assign node
type Assign struct {
	Node
}

func MakeAssign(lineno int, varNode *NValue, expr vm.INode, op int, driver *vm.Driver) *Assign{
	n := new(Assign)
	n.Lineno = lineno
	n.Left = varNode
	n.Right = expr
	n.Op = op
	n.Driver = driver
	return n
}

func (n *Assign) Push() *vm.VariableTag {

	// ただの代入以外は計算処理をするので、左辺をプッシュしておく
	if n.Op != OP_ASSIGN{
		t := n.Left.Push()
		if n.Driver.VariableTypeTable.IsStruct(t.VarType){
			n._err(cm.ERR_0042,"")
			return vm.MakeErrTag(n.Driver)
		}
		if t.ArraySize > 1{
			n._err(cm.ERR_0043, "")
			return vm.MakeErrTag(n.Driver)
		}
	}
	rightType := n.Right.Push()

	// ただの代入以外は計算命令を挿入
	if n.Op != OP_ASSIGN{
		// 文字列型
		if rightType.VarType.IsString() {
			if n.Op == OP_ADD_ASSIGN{
				n.Driver.OpAddString()
			} else {
				n._err(cm.ERR_0013, "")
				return vm.MakeErrTag(n.Driver)
			}
		// 小数点型
		} else if rightType.VarType.IsFloat() {
			switch n.Op{
			case OP_ADD_ASSIGN: n.Driver.OpFAdd()
			case OP_SUB_ASSIGN: n.Driver.OpFSub()
			case OP_MUL_ASSIGN: n.Driver.OpFMul()
			case OP_DIV_ASSIGN: n.Driver.OpFDiv()
			default:
				n._err(cm.ERR_0044, "")
				return vm.MakeErrTag(n.Driver)
			}
		// 整数型
		} else {
			switch n.Op{
			case OP_ADD_ASSIGN: n.Driver.OpAdd()
			case OP_SUB_ASSIGN: n.Driver.OpSub()
			case OP_MUL_ASSIGN: n.Driver.OpMul()
			case OP_DIV_ASSIGN: n.Driver.OpDiv()
			case OP_MOD_ASSIGN: n.Driver.OpMod()
			}
		}
	}

	leftType := n.Left.Pop()

	// 型チェック
	if rightType.VarType.IsDynamic(){
		return leftType
	} else if !rightType.TypeCompare(leftType) {
		n._err(cm.ERR_0017,
			"[" + strconv.Itoa(rightType.ArraySize) + "]" +
			rightType.VarType.TypeName + "->" +
			"[" + strconv.Itoa(leftType.ArraySize) + "]" +
			leftType.VarType.TypeName )
		return nil
	}
	return leftType
}

// 初期化時に使うAssign
type AssignAsInit struct{
	Assign
	index int
}

func MakeAssignAsInit(lineno int, varNode *NValue, expr vm.INode, op int, index int, driver *vm.Driver) *AssignAsInit{
	n := new(AssignAsInit)
	n.Lineno = lineno
	n.Left = varNode
	n.Right = expr
	n.Op = op
	n.Driver = driver
	
	n.index = index
	return n
}

func (n *AssignAsInit) Push() *vm.VariableTag{
	// 初期化時は演算禁止
	if n.Op != OP_ASSIGN{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0040, "")
		return nil
	}
	rightType := n.Right.Push()
	backPtr := n.Driver.GetPc()
	leftType := n.Left.Pop()

	// 型が不定の場合は推定する
	if leftType.VarType.IsUnknown(){
		if rightType.VarType.IsDynamic(){
			n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0041, "")
			return nil
		}

		// Unknown変数を削除し、新しく型が決まった変数を作り直す
		name := n.Driver.VariableTable.GetTag(n.index).Name
		n.Driver.VariableTable.RemoveLast()
		n.Driver.VariableTable.DefineValue(n.Lineno, name, rightType.VarType, rightType.IsPointer, rightType.ArraySize)

		// popをやり直す
		n.Driver.BackIndex(backPtr)
		leftType = n.Left.Pop()
	}

	// 型チェック
	if rightType.VarType.IsDynamic(){
		return leftType
	} else if !rightType.TypeCompare(leftType){
		n._err(cm.ERR_0017,
			"[" + strconv.Itoa(rightType.ArraySize) + "]" +
			rightType.VarType.TypeName + "->" +
			"[" + strconv.Itoa(leftType.ArraySize) + "]" +
			leftType.VarType.TypeName )
		return nil
	}
	return leftType
}