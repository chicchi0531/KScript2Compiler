package ast

import(
	cm "ks2/compiler/common"
)

// value node
type NodeValue struct {
	Node
	Name string
}

func (n *NodeValue) Push() int {
	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
	}

	n.Driver.OpPushValue(index)
	tag := n.Driver.VariableTable.GetTag(index)
	
	// エラーチェック
	if tag == nil{
		n._err(cm.ERR_0018,"")
		return -1
	}else{
		return tag.VarType
	}
}

func (n *NodeValue) Pop() int {
	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
	}

	n.Driver.OpPop(index)
	tag := n.Driver.VariableTable.GetTag(index)

	// エラーチェック
	if tag == nil{
		n._err(cm.ERR_0019,"")
		return -1
	}else{
		return tag.VarType
	}
}