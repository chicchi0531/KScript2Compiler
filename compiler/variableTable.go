package compiler

type VariableTag struct{
	name string
	varType int
}

type VariableTable struct{
	variables [][]*VariableTag
	currentTable int
	driver *Driver
}

// グローバル変数の定義
func (t *VariableTable) DefineInGlobal(variableTag *VariableTag) (int,error){
	if len(t.variables) == 0{
		t.variables = make([][]*VariableTag, 1)
	}
	t.variables[0] = append(t.variables[0], variableTag)
	return len(t.variables[0])-1,nil
}

// ローカル変数の定義
func (t *VariableTable) DefineInLocal(name string, varType int) int{
	// 変数テーブルがゼロだった場合に最初のレイヤーだけ確保する
	if len(t.variables) == 0{
		t.variables = make([][]*VariableTag, 1)
	}

	// 定義済みかどうかのチェック
	if t.FindVariable(name) != -1{
		t.driver.err.LogError(t.driver.filename, t.driver.lineno, ERR_0015, "識別子：" + name)
		return -1
	}

	tag := &VariableTag{name:name, varType:varType}
	t.variables[t.currentTable] = append(t.variables[t.currentTable], tag)

	// indexの計算
	index := 0
	for i:=0; i<t.currentTable; i++{
		index += len(t.variables[i])
	}

	println("ローカル変数定義:"+name)

	return index-1
}
// 初期化ありローカル変数定義
func (t *VariableTable) DefineInLocalWithAssign(name string, varType int, expr INode) *AssignNode {
	index := t.DefineInLocal(name, varType)
	varNode := &ValueNode{index:index}
	assignNode := &AssignNode{Node:Node{left:varNode, right:expr, driver:t.driver}}

	println("初期化ありローカル変数定義:"+name)

	return assignNode
}
// 型推論、初期化あり、ローカル変数定義
func (t *VariableTable) DefineInLocalWithAssignAutoType(name string, expr INode) *AssignNode {
	println("型推論、初期化ありローカル変数定義:"+name)

	return t.DefineInLocalWithAssign(name, TYPE_UNKNOWN, expr)
}

// スコープへ入る
func (t *VariableTable) ScopeIn(){
	t.currentTable++
	t.variables = append(t.variables, make([]*VariableTag, 0))
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
	searchIndex := 0
	for i:=0; i<t.currentTable; i++ {
		for j:=0; j<len(t.variables[i]); j++{
			if searchIndex == index{
				return t.variables[i][j]
			}
			searchIndex++
		}
	}
	return nil
}