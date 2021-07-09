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

	// テーブル類
	VariableTable *VariableTable
	StringTable *StringTable
	FloatTable *FloatTable
	FunctionTable *FunctionTable

	Err *cm.ErrorHandler
}

func (d *Driver) Init(filename string, err *cm.ErrorHandler){
	d.Pc = 0
	d.Filename = filename
	d.Program = make([]Op,0)
	d.Err = err
	d.VariableTable = &VariableTable{CurrentTable:0,Variables:make([][]*VariableTag,1),driver:d}
	d.FloatTable = &FloatTable{Values:make([]float32,0)}
	d.StringTable = &StringTable{Values:make([]string,0)}
	d.FunctionTable = &FunctionTable{}
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
func (d *Driver) AddFunction(returnType int, name string, args *ArgList, block IStateBlock) {
	d.VariableTable.ScopeIn()
	block.Analyze()
	d.VariableTable.ScopeOut()
}

// push_integer <value>
func (d *Driver) OpPushInteger(key int){
	prog := &Op{Code:VMCODE_PUSHINT, Value:key}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// push_float <value>
func (d *Driver) OpPushFloat(key float32){
	prog := &Op{Code:VMCODE_PUSHFLOAT, Value:d.FloatTable.Add(key)}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// push_string <value>
func (d *Driver) OpPushString(key string){
	prog := &Op{Code:VMCODE_PUSHSTRING, Value:d.StringTable.Add(key)}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// push_value <value_id>
func (d *Driver) OpPushValue(key int){
	prog := &Op{Code:VMCODE_PUSHVALUE, Value:key}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// pop_value <value_id>
func (d *Driver) OpPop(key int){
	prog := &Op{Code:VMCODE_POPVALUE, Value:key}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// add
func (d *Driver) OpAdd(){
	prog := &Op{Code:VMCODE_ADD}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// sub
func (d *Driver) OpSub(){
	prog := &Op{Code:VMCODE_SUB}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// mul
func (d *Driver) OpMul(){
	prog := &Op{Code:VMCODE_MUL}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// div
func (d *Driver) OpDiv(){
	prog := &Op{Code:VMCODE_DIV}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// mod
func (d *Driver) OpMod(){
	prog := &Op{Code:VMCODE_MOD}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// equal
func (d *Driver) OpEqual(){
	prog := &Op{Code:VMCODE_EQU}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// greater than
func (d *Driver) OpGt(){
	prog := &Op{Code:VMCODE_GT}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// greater equal
func (d *Driver) OpGe(){
	prog := &Op{Code:VMCODE_GE}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// less than
func (d *Driver) OpLt(){
	prog := &Op{Code:VMCODE_LT}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// less equal
func (d *Driver) OpLe(){
	prog := &Op{Code:VMCODE_LE}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// not equal
func (d *Driver) OpNequ(){
	prog := &Op{Code:VMCODE_NEQ}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// and
func (d *Driver) OpAnd(){
	prog := &Op{Code:VMCODE_AND}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// or
func (d *Driver) OpOr(){
	prog := &Op{Code:VMCODE_OR}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// add_string
func (d *Driver) OpAddString(){
	prog := &Op{Code:VMCODE_ADDSTRING}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
// not
func (d *Driver) OpNot(){
	prog := &Op{Code:VMCODE_NOT}
	d.Program = append(d.Program, *prog)
	d.Pc++
}
