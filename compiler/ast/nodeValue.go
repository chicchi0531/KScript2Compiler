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
	Child *NValue  //メンバー変数用
}

// constracter
func MakeValueNode(lineno int, name string, driver *vm.Driver) *NValue {
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Child = nil
	n.Driver = driver
	return n
}
func MakeArrayValueNode(lineno int, name string, index vm.INode, driver *vm.Driver) *NValue {
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Index = index
	n.Child = nil
	n.Driver = driver
	return n
}
func MakeMemberValueNode(lineno int, name string, child *NValue, driver *vm.Driver) *NValue{
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Child = child
	n.Driver = driver
	return n
}

func (n *NValue) Push() int {
	var vartype int
	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
	}

	// インデックス付きの場合は、インデックスをプッシュして、配列用命令でPUSH
	if n.Index != nil {
		n.Index.Push()
		n.Driver.OpPushArrayValue(index)
	// メンバ変数の場合は、メンバを探索して、オフセットを計算し配列用命令でPUSH
	} else if n.Child != nil{
		var memberIndex int
		vt := n.Driver.VariableTable.GetTag(index)
		memberIndex, vartype = n.calcMemberIndex(vt.VarType, n.Child)
		n.Driver.OpPushInteger(memberIndex+1)
		n.Driver.OpPushArrayValue(index)
	} else {
		n.Driver.OpPushValue(index)
	}

	tag := n.Driver.VariableTable.GetTag(index)

	// エラーチェック
	if tag == nil {
		n._err(cm.ERR_0018, "")
		return -1
	} else if n.Child != nil{
		return vartype
	} else {
		return tag.VarType
	}
}

func (n *NValue) Pop() int {
	var vartype int

	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
	}

	if n.Index != nil{
		n.Index.Push()
		n.Driver.OpPopArrayValue(index)
	} else if n.Child != nil{
		var memberIndex int
		vt := n.Driver.VariableTable.GetTag(index)
		memberIndex, vartype = n.calcMemberIndex(vt.VarType, n.Child)
		n.Driver.OpPushInteger(memberIndex+1)
		n.Driver.OpPopArrayValue(index)
	} else {
		n.Driver.OpPopValue(index)
	}

	tag := n.Driver.VariableTable.GetTag(index)

	// エラーチェック
	if tag == nil {
		n._err(cm.ERR_0019, "")
		return -1
	} else if n.Child != nil {
		return vartype
	} else {
		return tag.VarType
	}
}

// メンバー変数のオフセットを計算
func (n *NValue) calcMemberIndex(vartype int, node *NValue) (int,int) {
	tt := n.Driver.VariableTypeTable.GetTag(vartype)
	index := tt.FindMember(node.Name)
	if index == -1{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0002, node.Name)
		return 0, 0
	}

	// 再帰的に探索し、インデックスを計算
	vartype = tt.GetMember(index).VarType
	if node.Child != nil{
		var i int
		tag := tt.GetMember(index)
		i, vartype = n.calcMemberIndex(tag.VarType, node.Child)
		index += i
	}
	return index, vartype
}
