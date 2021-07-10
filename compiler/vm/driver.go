package vm

import(
	"fmt"
	cm "ks2/compiler/common"
)

type Op struct{
	Code int
	Value int
}

type Driver struct{
	Filename string
	Pc int
	Program []Op
	StackBase int

	// テーブル類
	VariableTable *VariableTable
	StringTable *StringTable
	FloatTable *FloatTable
	FunctionTable *FunctionTable

	// 解析中の一時情報
	CurrentRetType int

	Err *cm.ErrorHandler
}

func (d *Driver) Init(filename string, err *cm.ErrorHandler){
	d.Pc = 0
	d.Filename = filename
	d.Program = make([]Op,0)
	d.Err = err
	d.VariableTable = MakeVariableTable(d)
	d.FloatTable = MakeFloatTable()
	d.StringTable = MakeStringTable()
	d.FunctionTable = MakeFunctionTable(d)
}

//現在の状態を出力
func (d *Driver) Dump(){
	println("parse result=========")
	for i, op := range d.Program{
	  fmt.Printf("%d:%s %d\n",i,VMCODE_TOSTR(op.Code),op.Value)
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
func (d *Driver) AddFunction(lineno int, returnType int, name string, args []*Argument, block IStateBlock) {
	//関数定義
	d.FunctionTable.Add(&FunctionTag{Name:name, Args:args, RetrunType:returnType},lineno)
	d.CurrentRetType = returnType

	d.VariableTable.ScopeIn()

	//引数を変数定義
	for _,arg := range args{
		d.VariableTable.DefineInLocal(lineno, arg.Name, arg.VarType)
	}
	//ブロック内の命令をpush
	block.Analyze()

	//returnコードパスチェック
	if d.Program[d.Pc-1].Code != VMCODE_RETURN{
		if returnType == cm.TYPE_VOID{
			d.OpReturn()
		}else{
			d.Err.LogError(d.Filename, lineno, cm.ERR_0025, "")
		}
	}

	d.VariableTable.ScopeOut()
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
// pop_value <value_id>
func (d *Driver) OpPop(key int){
	d.addProg(VMCODE_POPVALUE,key)
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
func (d *Driver) OpSysCall(address int){
	d.addProg(VMCODE_SYSCALL, address)
}
// jmp
func (d *Driver) OpJmp(address int){
	d.addProg(VMCODE_JMP, address)
}
// return
func (d *Driver) OpReturn(){
	d.addProg(VMCODE_RETURN, 0)
}

func (d *Driver) addProg(code int, value int){
	prog := &Op{Code:code, Value:value}
	d.Program = append(d.Program, *prog)
	d.Pc++
}