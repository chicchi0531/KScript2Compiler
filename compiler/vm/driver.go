package vm

import(
	"fmt"
	cm "ks2/compiler/common"
)

type Op struct{
	Pc int
	Code int
	Value int
}

type Label struct{
	Index int
	Pos int
}

type Driver struct{
	Filename string
	Program []*Op
	StackBase int

	// テーブル類
	VariableTable *VariableTable
	StringTable *StringTable
	FloatTable *FloatTable
	FunctionTable *FunctionTable
	Labels []*Label

	// 解析中の一時情報
	CurrentRetType int
	BreakLabel int
	ContinueLabel int
	FallThroughLabel int

	Err *cm.ErrorHandler
}

func (d *Driver) Init(filename string, err *cm.ErrorHandler){
	d.Filename = filename
	d.Program = make([]*Op,0)
	d.Err = err
	d.VariableTable = MakeVariableTable(d)
	d.FloatTable = MakeFloatTable()
	d.StringTable = MakeStringTable()
	d.FunctionTable = MakeFunctionTable(d)
	d.Labels = make([]*Label, 0)

	d.CurrentRetType = cm.TYPE_INTEGER
	d.BreakLabel = -1
	d.ContinueLabel = -1
	d.FallThroughLabel = -1
}

//現在の状態を出力
func (d *Driver) Dump(){
	println("parse result=========")
	for _, op := range d.Program{
	  fmt.Printf("%d:%s %d\n",op.Pc,VMCODE_TOSTR(op.Code),op.Value)
	}
	println("value table==========")
	for i, v := range d.VariableTable.Variables{
	  for j, vv := range v{
		fmt.Printf("[%d][%d] name:%s type:%s\n",i,j,vv.Name, cm.TYPE_TOSTR(vv.VarType))
	  }
	}
	println("float table==========")
	for i, f := range d.FloatTable.Values{
		fmt.Printf("%d:%g\n",i,f)
	}
	println("string table=========")
	for i, s := range d.StringTable.Values{
		fmt.Printf("%d:%s\n",i,s)
	}
	println("=====================")
}

// function関係
func (d *Driver) DecralateFunction(lineno int, returnType int, name string, args []*Argument){
	tag := &FunctionTag{Name:name, Args:args, ReturnType:returnType, Defined:false }
	if d.FunctionTable.Find(name) != nil{
		d.Err.LogError(d.Filename, lineno, cm.ERR_0028, "関数："+name)
		return
	}

	d.FunctionTable.Add(tag, lineno)
}
func (d *Driver) AddFunction(lineno int, returnType int, name string, args []*Argument, statement IStatement) {
	
	//宣言済みの場合は、テーブルへの追加は行わず
	//Definedフラグを立てるだけ
	f := d.FunctionTable.Find(name)
	if f != nil{
		if !f.Defined{
			f.Defined = true
		}else{
			d.Err.LogError(d.Filename, lineno, cm.ERR_0023, "関数："+name)
			return
		}
	}else{
		//関数定義
		f = d.FunctionTable.Add(&FunctionTag{Name:name, Args:args, ReturnType:returnType},lineno)
	}
	d.CurrentRetType = returnType

	//ジャンプ用ラベルをセット
	d.SetLabel(f.Address)

	d.VariableTable.ScopeIn()

	//引数を変数定義
	for _,arg := range args{
		d.VariableTable.DefineValue(lineno, arg.Name, arg.VarType, arg.IsPointer, arg.Size)
	}
	//命令をpush
	statement.Analyze()

	//returnコードパスチェック
	if d.Program[len(d.Program)-1].Code != VMCODE_RETURN{
		if returnType == cm.TYPE_VOID{
			d.OpPushInteger(0)//ダミーの戻り値を積んでおく
			d.OpReturn()
		}else{
			d.Err.LogError(d.Filename, lineno, cm.ERR_0025, "")
		}
	}

	d.VariableTable.ScopeOut()
}

// label関係
func (d *Driver) MakeLabel() int{
	index := len(d.Labels)
	d.Labels = append(d.Labels, &Label{Index:index, Pos:0})
	return index
}

func (d *Driver) SetLabel(index int){
	d.addProg(VMCODE_DUMMYLABEL, index)
}

func (d *Driver) LabelSettings() int{
	pos := 0
	for _,p := range d.Program{
		// ラベル命令にアドレスを代入
		p.Pc = pos
		if p.Code == VMCODE_DUMMYLABEL{
			d.Labels[p.Value].Pos = pos
		}else{
			pos++
		}
	}

	// 各命令のアドレスを差し替える
	for _,p := range d.Program{
		switch p.Code{
		case VMCODE_JMP: fallthrough
		case VMCODE_JZE: fallthrough
		case VMCODE_JNZ: fallthrough
		case VMCODE_CALL:
			p.Value = d.Labels[p.Value].Pos
		}
	}
	return pos
}

// push_integer <value>
func (d *Driver) OpPushInteger(key int){
	d.addProg(VMCODE_PUSHINT,key)
}
// push_float <value>
func (d *Driver) OpPushFloat(key float32){
	d.addProg(VMCODE_PUSHFLOAT,d.FloatTable.Add(key))
}
// push_string <value>
func (d *Driver) OpPushString(key string){
	d.addProg(VMCODE_PUSHSTRING,d.StringTable.Add(key))
}
// push_value <value_id>
func (d *Driver) OpPushValue(key int){
	d.addProg(VMCODE_PUSHVALUE,key)
}
// push_array_value <value_id>
func (d *Driver) OpPushArrayValue(key int){
	d.addProg(VMCODE_PUSHARRAYVALUE, key)
}
// pop_value <value_id>
func (d *Driver) OpPopValue(key int){
	d.addProg(VMCODE_POPVALUE,key)
}
// pop_array_value <value_id>
func (d *Driver) OpPopArrayValue(key int){
	d.addProg(VMCODE_POPARRAYVALUE, key)
}
// pop
func (d *Driver) OpPop(){
	d.addProg(VMCODE_POP, 0)
}
// add
func (d *Driver) OpAdd(){
	d.addProg(VMCODE_ADD,0)
}
// sub
func (d *Driver) OpSub(){
	d.addProg(VMCODE_SUB,0)
}
// mul
func (d *Driver) OpMul(){
	d.addProg(VMCODE_MUL,0)
}
// div
func (d *Driver) OpDiv(){
	d.addProg(VMCODE_DIV,0)
}
// mod
func (d *Driver) OpMod(){
	d.addProg(VMCODE_MOD,0)
}
// incr
func (d *Driver) OpIncr(){
	d.addProg(VMCODE_INCR,0)
}
// decr
func (d *Driver) OpDecr(){
	d.addProg(VMCODE_DECR,0)
}
// equal
func (d *Driver) OpEqual(){
	d.addProg(VMCODE_EQU,0)
}
// greater than
func (d *Driver) OpGt(){
	d.addProg(VMCODE_GT,0)
}
// greater equal
func (d *Driver) OpGe(){
	d.addProg(VMCODE_GE,0)
}
// less than
func (d *Driver) OpLt(){
	d.addProg(VMCODE_LT,0)
}
// less equal
func (d *Driver) OpLe(){
	d.addProg(VMCODE_LE,0)
}
// not equal
func (d *Driver) OpNequ(){
	d.addProg(VMCODE_NEQ,0)
}
// and
func (d *Driver) OpAnd(){
	d.addProg(VMCODE_AND,0)
}
// or
func (d *Driver) OpOr(){
	d.addProg(VMCODE_OR,0)
}
// add_string
func (d *Driver) OpAddString(){
	d.addProg(VMCODE_ADDSTRING,0)
}
// not
func (d *Driver) OpNot(){
	d.addProg(VMCODE_NOT,0)
}
// call
func (d *Driver) OpCall(address int){
	d.addProg(VMCODE_CALL,address)
}
// sys call
func (d *Driver) OpSysCall(){
	d.addProg(VMCODE_SYSCALL, 0)
}
// jmp
func (d *Driver) OpJmp(address int){
	d.addProg(VMCODE_JMP, address)
}
// jzero
func (d *Driver) OpJze(address int){
	d.addProg(VMCODE_JZE, address)
}
// jnotzero
func (d *Driver) OpJnz(address int){
	d.addProg(VMCODE_JNZ, address)
}
// return
func (d *Driver) OpReturn(){
	d.addProg(VMCODE_RETURN, 0)
}

func (d *Driver) addProg(code int, value int){
	prog := &Op{Code:code, Value:value}
	d.Program = append(d.Program, prog)
}

func (d *Driver) RemoveLastProg(){
	d.Program = d.Program[:len(d.Program)-1]
}