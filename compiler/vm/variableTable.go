package vm

import (
	cm "ks2/compiler/common"
)

type VariableTag struct {
	Name      string
	VarType   *VariableTypeTag
	IsPointer bool
	ArraySize int
	Offset    int //構造体メンバ用
}

func MakeVariableTag(name string, vartype *VariableTypeTag, ispointer bool, arraysize int, driver *Driver) *VariableTag {
	t := new(VariableTag)
	t.Name = name
	t.VarType = vartype
	t.IsPointer = ispointer
	t.ArraySize = arraysize
	t.Offset = 0

	return t
}

func MakeErrTag(driver *Driver) *VariableTag{
	t := new(VariableTag)
	t.Name = ""
	t.VarType = driver.VariableTypeTable.GetTag(0)
	t.IsPointer = false
	t.ArraySize = 1
	t.Offset = 0
	return t
}

func (t *VariableTag) TypeCompare(dst *VariableTag) bool {
	return (t.VarType == dst.VarType) && (t.ArraySize == dst.ArraySize)
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
func (t *VariableTable) DefineValue(lineno int, name string, varType *VariableTypeTag, isPointer bool, arraysize int) (int, bool) {
	// 定義済みかどうかのチェック
	index, _ := t.FindVariable(name)
	if index != -1 && name != "" {
		t.driver.Err.LogError(t.driver.Filename, lineno, cm.ERR_0015, "識別子："+name)
		return -1, false
	}
	// サイズが1以上かのチェック
	if arraysize <= 0 {
		t.driver.Err.LogError(t.driver.Filename, lineno, cm.ERR_0033, "")
		return -1, false
	}

	// 定義
	// 配列の場合は、サイズ分メモリを埋める
	// 配列の場合は、２番目以降の要素の名前は空にする
	tmpName := name
	for i := 0; i < arraysize; i++ {
		vt := MakeVariableTag(tmpName, varType, isPointer, arraysize, t.driver)
		t.Variables[t.CurrentTable] = append(t.Variables[t.CurrentTable], vt)
		// メンバー分のメモリ確保
		if t.driver.VariableTypeTable.IsStruct(varType) {
			t.defineStructMember(lineno, varType)
		}
		tmpName = ""
	}

	// indexの計算
	if t.CurrentTable == 0 {
		// global変数の場合
		index = len(t.Variables[0])
		index -= varType.Size * arraysize
		return index, true
	} else {
		index = 0
		for i := 1; i <= t.CurrentTable; i++ {
			index += len(t.Variables[i])
		}
		index -= varType.Size * arraysize
		return index, false
	}
}

func (t *VariableTable) defineStructMember(lineno int, tt *VariableTypeTag) {
	// 構造体のメンバは直接検索に引っかからないよう、空名にしておく
	for _, m := range tt.Member {
		t.DefineValue(lineno, "", m.VarType, m.IsPointer, m.ArraySize)
	}
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
func (t *VariableTable) FindVariable(name string) (int, bool) {

	// グローバル変数テーブルから検索
	index := 0
	for j := 0; j < len(t.Variables[0]); j++ {
		if t.Variables[0][j].Name == name {
			return index,true
		}
		index++;
	}

	// グローバル変数テーブル以外から検索
	index = 0
	for i := 1; i <= t.CurrentTable; i++ {
		for j := 0; j < len(t.Variables[i]); j++ {
			if t.Variables[i][j].Name == name {
				return index,false
			}
			index++
		}
	}
	return -1, false
}

// 指定インデックスのvariable tagを取得
func (t *VariableTable) GetTag(index int, isglobal bool) *VariableTag {
	searchIndex := 0

	if isglobal {
		for j := 0; j < len(t.Variables[0]); j++ {
			if searchIndex == index {
				return t.Variables[0][j]
			}
			searchIndex++
		}
	} else {
		for i := 1; i <= t.CurrentTable; i++ {
			for j := 0; j < len(t.Variables[i]); j++ {
				if searchIndex == index {
					return t.Variables[i][j]
				}
				searchIndex++
			}
		}
	}
	return nil
}

// 指定インデックスの変数を削除
// 実際は名前を消して検索できなくするだけなので、
// 他の変数のインデックスには影響しない
// func (t *VariableTable) DeleteTag(index int) {
// 	vt := t.GetTag(index, false)
// 	if vt != nil {
// 		vt.Name = ""
// 	}
// }

// テーブルの最後の変数を削除
func (t *VariableTable) RemoveLast() {
		t.Variables[t.CurrentTable] = 
			t.Variables[t.CurrentTable][:len(t.Variables[t.CurrentTable])-1]
}

// グローバル変数のサイズを取得
func (t *VariableTable) GetGlobalValueSize() int{
	return len(t.Variables[0])
}