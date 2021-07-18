package vm

import(
	cm "ks2/compiler/common"
)

// variable type tag
type VariableTypeTag struct{
	typename string
	member []*VariableTag
	method []*FunctionTag
	Size int
}

func MakeVariableTypeTag(typename string) *VariableTypeTag{
	t := new(VariableTypeTag)
	t.typename =typename
	t.member = make([]*VariableTag, 0)
	t.method = make([]*FunctionTag, 0)
	t.Size = 1
	return t
}

func (t *VariableTypeTag) FindMember(name string) int {
	for i, m := range t.member{
		if m.Name == name{
			return i
		}
	}
	return -1
}

func (t *VariableTypeTag) GetMember(index int) *VariableTag{
	return t.member[index]
}

func (t *VariableTypeTag) FindMethod(name string) int {
	for _, m := range t.method{
		if m.Name == name{
			return m.Address
		}
	}
	return -1
}

func (t *VariableTypeTag) AddMember(name string, vartype int, ispointer bool, arraysize int, driver *Driver){
	tag := new(VariableTag)
	tag.Name = name
	tag.VarType = vartype
	tag.IsPointer = ispointer
	tag.ArraySize = arraysize
	tag.Offset = t.Size

	t.member = append(t.member, tag)

	// 構造体サイズの計算
	if vartype >= cm.TYPE_STRUCT{
		tt := driver.VariableTypeTable.GetTag(vartype)
		t.Size += tt.Size * arraysize
	}else{
		t.Size += arraysize
	}
}

func (t *VariableTypeTag) AddMethod(name string, returntype int){
	//TODO:
}

// メンバーが循環参照になっていないかチェック
var checkTypeName = ""
func (t *VariableTypeTag) CheckMember(lineno int, driver *Driver)int{
	result := 0
	for _,m := range t.member{
		if m.VarType >= cm.TYPE_STRUCT{
			tt := driver.VariableTypeTable.GetTag(m.VarType)
			// 調べる型がメンバに含まれていると循環参照とみなす
			if checkTypeName == tt.typename{
				return -1
			}
			// メンバに潜って探索
			if tt.CheckMember(lineno, driver) == -1{
				return -1
			}
		}
	}
	return result
}

// variable type table
type VariableTypeTable struct{
	tags []*VariableTypeTag
	driver *Driver
}

func MakeVariableTypeTable(driver *Driver) *VariableTypeTable{
	t := new(VariableTypeTable)
	t.tags = make([]*VariableTypeTag, 0)
	t.driver = driver
	return t
}

func (t *VariableTypeTable) Add(tag *VariableTypeTag){
	t.tags = append(t.tags, tag)
}

func (t *VariableTypeTable) Find(name string) (int, *VariableTypeTag){
	for i, tag := range t.tags{
		if tag.typename == name{
			return i + cm.TYPE_STRUCT, tag
		}
	}
	return -1,nil
}

func (t *VariableTypeTable) GetTag(typeid int) *VariableTypeTag{
	return t.tags[typeid - cm.TYPE_STRUCT]
}
