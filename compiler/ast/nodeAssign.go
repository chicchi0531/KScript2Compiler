package ast

import(
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

// assign node
type NodeAssign struct {
	Node
}

func MakeNodeAssign(lineno int, valNode *NodeValue, expr vm.INode, driver *vm.Driver) *NodeAssign{
	n := new(NodeAssign)
	n.Lineno = lineno
	n.Left = valNode
	n.Right = expr
	n.Driver = driver
	return n
}

func (n *NodeAssign) Push() int {
	rightType := n.Right.Push()
	leftType := n.Left.Pop()

	// 型チェック
	if leftType == cm.TYPE_UNKNOWN{
		// 左辺が型未推定の場合は、推定して仕込む
		// 直前のポップ命令から変数番号を取得
		varIndex := n.Driver.Program[len(n.Driver.Program)-1].Value
		varTag := n.Driver.VariableTable.GetTag(varIndex)
		varTag.VarType = rightType //右辺の型をそのまま左辺の型とする
		return rightType

	}else if rightType == cm.TYPE_STRING || leftType == cm.TYPE_STRING{
		if rightType != leftType{
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
