package vm

import(
	cm "ks2/compiler/common"
)

type Argument struct{
	Name string
	VarType int
}

type FunctionTag struct{
	Name string
	Args []*Argument
	Address int
	RetrunType int
	IsSystem bool
}

// function table
type FunctionTable struct{
	Functions []*FunctionTag
	driver *Driver
}

func MakeFunctionTable(d *Driver) *FunctionTable{
	f := new(FunctionTable)
	f.Functions = make([]*FunctionTag, 0)
	f.driver = d
	return f
}

func (t *FunctionTable) Add(tag *FunctionTag, lineno int){
	//定義済みかチェック
	f := t.Find(tag.Name)
	if f != nil{
		t.driver.Err.LogError(t.driver.Filename, lineno, cm.ERR_0023, "関数："+tag.Name)
		return
	}

	//call用のアドレスを設定
	tag.Address = t.driver.Pc
	t.Functions = append(t.Functions, tag)
}

// 同名の関数を探索
func (t *FunctionTable) Find(name string) *FunctionTag{
	for _, f := range t.Functions{
		if f.Name == name{
			return f
		}
	}
	return nil
}
