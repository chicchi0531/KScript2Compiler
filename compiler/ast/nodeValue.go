package ast

import (
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

// value node
type NValue struct {
	Node
	Name  string
	Index vm.INode //配列用
}

func MakeValueNode(lineno int, name string, driver *vm.Driver) *NValue {
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Driver = driver
	return n
}

func MakeArrayValueNode(lineno int, name string, index vm.INode, driver *vm.Driver) *NValue {
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Index = index
	n.Driver = driver
	return n
}

func (n *NValue) Push() int {
	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
	}

	// インデックス付きの場合は、インデックスをプッシュして、配列用命令でPUSH
	if n.Index != nil {
		n.Index.Push()
		n.Driver.OpPushArrayValue(index)
	} else {
		n.Driver.OpPushValue(index)
	}

	tag := n.Driver.VariableTable.GetTag(index)

	// エラーチェック
	if tag == nil {
		n._err(cm.ERR_0018, "")
		return -1
	} else {
		return tag.VarType
	}
}

func (n *NValue) Pop() int {
	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
	}

	if n.Index != nil{
		n.Index.Push()
		n.Driver.OpPopArrayValue(index)
	} else {
		n.Driver.OpPopValue(index)
	}

	tag := n.Driver.VariableTable.GetTag(index)

	// エラーチェック
	if tag == nil {
		n._err(cm.ERR_0019, "")
		return -1
	} else {
		return tag.VarType
	}
}
