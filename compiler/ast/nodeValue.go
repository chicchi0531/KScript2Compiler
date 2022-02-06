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
func MakeMemberValueNode(lineno int, name string, child *NValue, driver *vm.Driver) *NValue {
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Child = child
	n.Driver = driver
	return n
}
func MakeArrayMemberValueNode(lineno int, name string, index vm.INode, child *NValue, driver *vm.Driver) *NValue {
	n := new(NValue)
	n.Lineno = lineno
	n.Name = name
	n.Index = index
	n.Child = child
	n.Driver = driver
	return n
}

func (n *NValue) Push() *vm.VariableTag {
	index, vt, isglobal := n.getVariableTag(n.Name)
	if vt == nil {
		return vm.MakeErrTag(n.Driver)
	}

	var lastNode *NValue
	vt, lastNode = n.PushAddr(index, vt)

	// 配列アクセスしている場合は単体のサイズを、
	// 配列そのものの場合は配列込みのサイズをプッシュ
	if lastNode.Index == nil{
		if isglobal {
			n.Driver.OpPushValueRange(vt.VarType.Size * vt.ArraySize)
		} else {
			n.Driver.OpPushLocalRange(vt.VarType.Size * vt.ArraySize)
		}
	} else {
		if isglobal {
			n.Driver.OpPushValueRange(vt.VarType.Size)
		} else {
			n.Driver.OpPushLocalRange(vt.VarType.Size)
		}
	}

	// 型の作成
	arraySize := vt.ArraySize
	if lastNode.Index != nil{
		arraySize = 1
	}
	resultType := vm.MakeVariableTag("", vt.VarType, vt.IsPointer, arraySize, n.Driver)

	return resultType
}

func (n *NValue) Pop() *vm.VariableTag {
	index, vt, isglobal := n.getVariableTag(n.Name)
	if vt == nil {
		return vm.MakeErrTag(n.Driver)
	}

	var lastNode *NValue
	vt, lastNode = n.PushAddr(index, vt)

	// 構造体の場合はサイズ分同様にpush
	if n.Driver.VariableTypeTable.IsStruct(vt.VarType) {
		// 配列アクセスしている場合は単体のサイズを、
		// 配列そのものの場合は配列込みのサイズをプッシュ
		if lastNode.Index == nil{
			if isglobal {
				n.Driver.OpPopValueRange(vt.VarType.Size * vt.ArraySize)
			} else {
				n.Driver.OpPopLocalRange(vt.VarType.Size * vt.ArraySize)
			}
		} else {
			if isglobal {
				n.Driver.OpPopValueRange(vt.VarType.Size)
			} else {
				n.Driver.OpPopLocalRange(vt.VarType.Size)
			}
		}
	} else {
		if isglobal {
			n.Driver.OpPopValue()
		} else {
			n.Driver.OpPopLocal()
		}
	}

	// 型の作成
	arraySize := vt.ArraySize
	if lastNode.Index != nil{
		arraySize = 1
	}
	resultType := vm.MakeVariableTag("", vt.VarType, vt.IsPointer, arraySize, n.Driver)
	
	return resultType
}

// 自身のアドレスをプッシュする
func (n *NValue) PushAddr(index int, vt *vm.VariableTag) (*vm.VariableTag, *NValue) {

	// 自身のインデックスをプッシュ
	n.Driver.OpPushInteger(index)

	// 配列の場合
	if n.Index != nil {
		n.Index.Push()
		n.Driver.OpPushInteger(vt.VarType.Size)
		n.Driver.OpMul()
		n.Driver.OpAdd()
	}

	// チェーンが続いている場合は子供のアドレス計算
	if n.Child != nil {
		return n.pushChildAddr(vt.VarType, n.Child)
	}

	return vt, n
}

func (n *NValue) pushChildAddr(parentType *vm.VariableTypeTag, child *NValue) (*vm.VariableTag, *NValue) {
	if !n.Driver.VariableTypeTable.IsStruct(parentType) {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0037, "")
		return nil, nil
	}

	memberID := parentType.FindMember(child.Name)
	if memberID == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0002, child.Name)
		return nil, nil
	}

	// 親の構造体から見た自身のメンバオフセットをプッシュ
	memberTt := parentType.GetMember(memberID)
	n.Driver.OpPushInteger(memberTt.Offset)
	n.Driver.OpAdd()

	// 配列の場合はサイズ分ずらす
	if child.Index != nil {
		child.Index.Push()
		n.Driver.OpPushInteger(memberTt.VarType.Size)
		n.Driver.OpMul()
		n.Driver.OpAdd()
	}

	// チェーンが続いている場合は、子供のアドレスをプッシュしていく
	if child.Child != nil {
		return child.Child.pushChildAddr(memberTt.VarType, child.Child)
	}

	return memberTt, child
}

func (n *NValue) getVariableTag(name string) (int, *vm.VariableTag, bool) {
	// 変数のインデックスを取得
	index, isglobal := n.Driver.VariableTable.FindVariable(n.Name)
	if index == -1 {
		n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0016, "不明な識別子："+n.Name)
		return 0, nil, false
	}

	// アドレスをプッシュしてから、プッシュ命令を発行
	vt := n.Driver.VariableTable.GetTag(index, isglobal)
	if vt == nil {
		n._err(cm.ERR_0018, "")
		return 0, nil, false
	}
	return index, vt, isglobal
}
