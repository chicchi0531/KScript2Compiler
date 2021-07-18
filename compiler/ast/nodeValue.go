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
func MakeArrayMemberValueNode(lineno int, name string, index vm.INode, child *NValue, driver *vm.Driver) *NValue{
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Index = index
	n.Child = child
	n.Driver = driver
	return n
}

func (n *NValue) Push() int {
	index,vt := n.getVariableTag()
	if vt == nil{
		return -1
	}

	vartype := n.PushAddr(index, vt.VarType)
	n.Driver.OpPushValue()

	return vartype
}

func (n *NValue) Pop() int {
	index,vt := n.getVariableTag()
	if vt == nil{
		return -1
	}

	vartype := n.PushAddr(index, vt.VarType)
	n.Driver.OpPopValue()

	return vartype
}

// 自身のアドレスをプッシュする
func (n *NValue) PushAddr(index int, vartype int) int {
	size := 1
	if vartype >= cm.TYPE_STRUCT{
		size = n.Driver.VariableTypeTable.GetTag(vartype).Size
	}

	// 自身のインデックスをプッシュ
	n.Driver.OpPushInteger(index)

	// 配列の場合
	if n.Index != nil{
		n.Index.Push()
		n.Driver.OpPushInteger(size)
		n.Driver.OpMul()
		n.Driver.OpAdd()
	}

	// チェーンが続いている場合は子供のアドレス計算
	if n.Child != nil{
		return n.pushChildAddr(vartype, n.Child)
	}

	return vartype
}
func (n *NValue) pushChildAddr(parentType int, child *NValue) int {
	if parentType < cm.TYPE_STRUCT{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0037, "")
		return 0
	}
	
	tt := n.Driver.VariableTypeTable.GetTag(parentType)
	memberID := tt.FindMember(child.Name)
	if memberID == -1{
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0002, child.Name)
		return 0
	}

	// 親の構造体から見た自身のメンバオフセットをプッシュ
	memberTt := tt.GetMember(memberID)
	n.Driver.OpPushInteger(memberTt.Offset)
	n.Driver.OpAdd()

	// 配列の場合はサイズ分ずらす
	if child.Index != nil{
		child.Index.Push()
		n.Driver.OpPushInteger(memberTt.ArraySize)
		n.Driver.OpMul()
		n.Driver.OpAdd()
	}

	// チェーンが続いている場合は、子供のアドレスをプッシュしていく
	if child.Child != nil{
		return child.Child.pushChildAddr(memberTt.VarType, child.Child)
	}

	return memberTt.VarType
}

func (n *NValue) getVariableTag() (int,*vm.VariableTag){
	// 変数のインデックスを取得
	index := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
		return 0,nil
	}

	// アドレスをプッシュしてから、プッシュ命令を発行
	vt := n.Driver.VariableTable.GetTag(index)
	if vt == nil {
		n._err(cm.ERR_0018, "")
		return 0,nil
	}
	return index,vt
}