package vm

import (
	cm "ks2/compiler/common"
)

type VariableTag struct {
	Name      string
	VarType   int
	IsPointer bool
	Size int
}

type VariableTable struct {
	Variables    [][]*VariableTag
	CurrentTable int
	driver       *Driver
}

func MakeVariableTable(driver *Driver) *VariableTable {
	t := new(VariableTable)
	t.Variables = make([][]*VariableTag, 1)
	t.CurrentTable = 0
	t.driver = driver
	return t
}

// ローカル変数の定義
func (t *VariableTable) DefineValue(lineno int, name string, varType int, isPointer bool, size int) int {
	// 定義済みかどうかのチェック
	if t.FindVariable(name) != -1 {
		t.driver.Err.LogError(t.driver.Filename, lineno, cm.ERR_0015, "識別子："+name)
		return -1
	}
	// サイズが1以上かのチェック
	if size <= 0{
		t.driver.Err.LogError(t.driver.Filename,lineno, cm.ERR_0033, "")
		return -1
	}

	// メモリ確保
	tag := &VariableTag{Name: name, VarType: varType, IsPointer: isPointer, Size:size }
	t.Variables[t.CurrentTable] = append(t.Variables[t.CurrentTable], tag)

	// 配列の場合は、サイズ分メモリを埋める
	for i:=0; i<size-1; i++{
		t.Variables[t.CurrentTable] = append(t.Variables[t.CurrentTable],
			&VariableTag{Name: "", VarType: varType, IsPointer: isPointer, Size:1 })
	}

	// indexの計算
	index := 0
	for i := 0; i <= t.CurrentTable; i++ {
		index += len(t.Variables[i])
	}
	index -= size

	return index
}

// スコープへ入る
func (t *VariableTable) ScopeIn() {
	t.CurrentTable++
	t.Variables = append(t.Variables, make([]*VariableTag, 0))
}

// スコープから抜ける
func (t *VariableTable) ScopeOut() {
	t.CurrentTable--
	t.Variables = t.Variables[:len(t.Variables)-1]
}

// 定義済み変数の検索
func (t *VariableTable) FindVariable(name string) int {
	index := 0
	for i := 0; i <= t.CurrentTable; i++ {
		for j := 0; j < len(t.Variables[i]); j++ {
			if t.Variables[i][j].Name == name {
				return index
			}
			index++
		}
	}
	return -1
}

// 指定インデックスのvariable tagを取得
func (t *VariableTable) GetTag(index int) *VariableTag {
	searchIndex := 0
	for i := 0; i <= t.CurrentTable; i++ {
		for j := 0; j < len(t.Variables[i]); j++ {
			if searchIndex == index {
				return t.Variables[i][j]
			}
			searchIndex++
		}
	}
	return nil
}
