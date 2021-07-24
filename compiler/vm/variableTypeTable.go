package vm

import(
	cm "ks2/compiler/common"
)

// variable type tag
type VariableTypeTag struct{
	TypeName string
	Member []*VariableTag
	Method []*FunctionTag
	Size int

	//各種ステータス類
	IsIncrementable bool
	IsNotable bool
}

func MakeVariableTypeTag(typename string) *VariableTypeTag{
	t := new(VariableTypeTag)
	t.TypeName = typename
	t.Member = make([]*VariableTag, 0)
	t.Method = make([]*FunctionTag, 0)
	t.Size = 1

	t.IsIncrementable = false
	t.IsNotable = false
	return t
}

func (t *VariableTypeTag) FindMember(name string) int {
	for i, m := range t.Member{
		if m.Name == name{
			return i
		}
	}
	return -1
}

func (t *VariableTypeTag) GetMember(index int) *VariableTag{
	return t.Member[index]
}

func (t *VariableTypeTag) FindMethod(name string) int {
	for _, m := range t.Method{
		if m.Name == name{
			return m.Address
		}
	}
	return -1
}

func (t *VariableTypeTag) AddMember(name string, vartype *VariableTypeTag, ispointer bool, arraysize int, driver *Driver){
	tag := MakeVariableTag(name, vartype, ispointer, arraysize, driver)
	tag.Offset = t.Size

	t.Member = append(t.Member, tag)

	// 構造体サイズの計算
	t.Size += vartype.Size * arraysize
}

func (t *VariableTypeTag) AddMethod(name string, returntype int){
	//TODO:
}

// メンバーが循環参照になっていないかチェック
var checkTypeName = ""
func (t *VariableTypeTag) CheckMember(lineno int, driver *Driver) bool {

	for _,m := range t.Member{
		// 調べる型がメンバに含まれていると循環参照とみなす
		if checkTypeName == m.VarType.TypeName{
			return false
		}
		// メンバに潜って探索
		if !m.VarType.CheckMember(lineno, driver){
			return false
		}
	}
	return true
}

func (t *VariableTypeTag) IsDynamic() bool {
	return t.TypeName == "dynamic"
}

func (t *VariableTypeTag) IsString() bool {
	return t.TypeName == "string"
}

func (t *VariableTypeTag) IsUnknown() bool {
	return t.TypeName == "unknown"
}

// variable type table
type VariableTypeTable struct{
	tags []*VariableTypeTag
	driver *Driver
}

func MakeVariableTypeTable(driver *Driver) *VariableTypeTable{
	t := new(VariableTypeTable)
	t.tags = make([]*VariableTypeTag, cm.TYPE_STRUCT)
	t.driver = driver

	// デフォルトの型を定義
	t.addDefaultType("int", cm.TYPE_INTEGER, true, true)
	t.addDefaultType("float", cm.TYPE_FLOAT, false, true)
	t.addDefaultType("string", cm.TYPE_STRING, false, false)
	t.addDefaultType("unknown", cm.TYPE_UNKNOWN, false, false)
	t.addDefaultType("dynamic", cm.TYPE_DYNAMIC, false, false)

	return t
}

func (t *VariableTypeTable) Add(tag *VariableTypeTag){
	t.tags = append(t.tags, tag)
}

func (t *VariableTypeTable) Find(name string) (int, *VariableTypeTag){
	for i, tag := range t.tags{
		if tag == nil { continue }
		if tag.TypeName == name{
			return i, tag
		}
	}
	return -1,nil
}

func (t *VariableTypeTable) GetTag(typeid int) *VariableTypeTag{
	return t.tags[typeid]
}

func (t *VariableTypeTable) IsStruct(tt *VariableTypeTag) bool {
	i,_ := t.Find(tt.TypeName)
	return i >= cm.TYPE_STRUCT
}

func (t *VariableTypeTable) addDefaultType(name string, id int, incrementable bool, isnotable bool){
	tt := MakeVariableTypeTag(name)
	tt.IsIncrementable = incrementable
	tt.IsNotable = isnotable
	t.tags[id] = tt
}