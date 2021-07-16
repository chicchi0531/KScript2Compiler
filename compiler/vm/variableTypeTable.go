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
	t.Size = 0
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

func (t *VariableTypeTag) AddMember(name string, vartype int, ispointer bool, size int){
	tag := new(VariableTag)
	tag.Name = name
	tag.VarType = vartype
	tag.IsPointer = ispointer
	tag.Size = size

	t.member = append(t.member, tag)
}

func (t *VariableTypeTag) AddMethod(name string, returntype int){
	//TODO:
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
