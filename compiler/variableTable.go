package compiler

import(
	"errors"
)

type VariableTag struct{
	name string
	varType int
}

type VariableTable struct{
	variables [][]VariableTag
	currentTable int
}

// グローバル変数の定義
func (t *VariableTable) DefineInGlobal(variableTag VariableTag) (int,error){
	if len(t.variables) == 0{
		t.variables = make([][]VariableTag, 1)
	}
	t.variables[0] = append(t.variables[0], variableTag)
	return len(t.variables[0])-1,nil
}

// ローカル変数の定義
func (t *VariableTable) DefineInLocal(variableTag VariableTag) (int, error){
	// 定義済みかどうかのチェック
	if t.FindVariable(variableTag.name) != -1{
		return 0,errors.New(ERR_0015)
	}
	t.variables[t.currentTable] = append(t.variables[t.currentTable], variableTag)

	// indexの計算
	index := 0
	for i:=0; i<t.currentTable; i++{
		index += len(t.variables[i])
	}

	return index-1,nil
}

// スコープへ入る
func (t *VariableTable) ScopeIn(){
	t.currentTable++
	t.variables = append(t.variables, make([]VariableTag, 0))
}

// スコープから抜ける
func (t *VariableTable) ScopeOut(){
	t.currentTable--
	t.variables = t.variables[:len(t.variables)-1]
}

// 定義済み変数の検索
func (t *VariableTable) FindVariable(name string) int{
	index := 0
	for i:=0; i<t.currentTable; i++{
		for j:=0; j<len(t.variables[i]); j++{
			if t.variables[i][j].name == name{
				return index
			}
			index++
		}
	}
	return -1
}

// 指定インデックスのvariable tagを取得
func (t *VariableTable) GetTag(index int) *VariableTag{
	//TODO
	return nil
}