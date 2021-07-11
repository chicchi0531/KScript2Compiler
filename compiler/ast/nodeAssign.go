package ast

import(
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

// assign node
type Assign struct {
	Node
}

func MakeAssign(lineno int, valNode *NValue, expr vm.INode, driver *vm.Driver) *Assign{
	n := new(Assign)
	n.Lineno = lineno
	n.Left = valNode
	n.Right = expr
	n.Driver = driver
	return n
}

func (n *Assign) Push() int {
	rightType := n.Right.Push()
	leftType := n.Left.Pop()

	// 型チェック
	if leftType == cm.TYPE_UNKNOWN{
		// 型推定でTYPE_DYNAMICを使うことは禁止する
		if rightType == cm.TYPE_DYNAMIC{
			n._err(cm.ERR_0027,"")
			return rightType
		}

		// 左辺が型未推定の場合は、推定して仕込む
		// 直前のポップ命令から変数番号を取得
		varIndex := n.Driver.Program[len(n.Driver.Program)-1].Value
		varTag := n.Driver.VariableTable.GetTag(varIndex)
		varTag.VarType = rightType //右辺の型をそのまま左辺の型とする
		return rightType

	}else if rightType == cm.TYPE_STRING || leftType == cm.TYPE_STRING{

		if rightType == cm.TYPE_DYNAMIC || leftType == cm.TYPE_DYNAMIC{
			return cm.TYPE_STRING
		}else if rightType != leftType{
			n._err(cm.ERR_0017, "")
			return -1
		}else{
			return cm.TYPE_STRING
		}
	}

	// 数値型の場合は、代入先の型をそのまま返す
	// ダウンキャストになる場合は警告を出しておく
	if leftType == cm.TYPE_INTEGER && rightType == cm.TYPE_FLOAT{
		n._warning(cm.WARNING_0001,"float->int")
	}
	return leftType
}
