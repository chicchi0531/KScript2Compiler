package compiler

import(
	"fmt"
)

type Op struct{
	code int
	value int
}

type Driver struct{
	filename string
	lineno int
	pc int
	program []Op

	// テーブル類
	variableTable *VariableTable
	stringTable *StringTable
	floatTable *FloatTable
	functionTable *FunctionTable

	err *ErrorHandler
}

func (d *Driver) Init(filename string){
	d.pc = 0
	d.lineno = 0
	d.filename = filename
	d.program = make([]Op,0)
	d.err = &ErrorHandler{errorCount:0, warningCount:0}
	d.variableTable = &VariableTable{currentTable:0,driver:d}
	d.floatTable = &FloatTable{values:make([]float32,0)}
	d.stringTable = &StringTable{values:make([]string,0)}
	d.functionTable = &FunctionTable{}
}

//現在の状態を出力
func (d *Driver) Dump(){
	println("parse result=========")
	for i, op := range d.program{
	  fmt.Printf("%d:%s %d\n",i,VMCODE_TOSTR(op.code),op.value)
	}
	println("value table==========")
	for i, v := range d.variableTable.variables{
	  for j, vv := range v{
		fmt.Printf("[%d][%d] name:%s type:%s\n",i,j,vv.name, TYPE_TOSTR(vv.varType))
	  }
	}
	println("float table==========")
	for i, f := range d.floatTable.values{
		fmt.Printf("%d:%g\n",i,f)
	}
	println("string table=========")
	for i, s := range d.stringTable.values{
		fmt.Printf("%d:%s\n",i,s)
	}
	println("=====================")
}

// push_integer <value>
func (d *Driver) OpPushInteger(key int){
	prog := &Op{code:VMCODE_PUSHINT, value:key}
	d.program = append(d.program, *prog)
	d.pc++
}
// push_float <value>
func (d *Driver) OpPushFloat(key float32){
	prog := &Op{code:VMCODE_PUSHFLOAT, value:d.floatTable.Add(key)}
	d.program = append(d.program, *prog)
	d.pc++
}
// push_string <value>
func (d *Driver) OpPushString(key string){
	prog := &Op{code:VMCODE_PUSHSTRING, value:d.stringTable.Add(key)}
	d.program = append(d.program, *prog)
	d.pc++
}
// push_value <value_id>
func (d *Driver) OpPushValue(key int){
	prog := &Op{code:VMCODE_PUSHVALUE, value:key}
	d.program = append(d.program, *prog)
	d.pc++
}
// pop_value <value_id>
func (d *Driver) OpPop(key int){
	prog := &Op{code:VMCODE_POPVALUE, value:key}
	d.program = append(d.program, *prog)
	d.pc++
}
// add
func (d *Driver) OpAdd(){
	prog := &Op{code:VMCODE_ADD}
	d.program = append(d.program, *prog)
	d.pc++
}
// sub
func (d *Driver) OpSub(){
	prog := &Op{code:VMCODE_SUB}
	d.program = append(d.program, *prog)
	d.pc++
}
// mul
func (d *Driver) OpMul(){
	prog := &Op{code:VMCODE_MUL}
	d.program = append(d.program, *prog)
	d.pc++
}
// div
func (d *Driver) OpDiv(){
	prog := &Op{code:VMCODE_DIV}
	d.program = append(d.program, *prog)
	d.pc++
}
// mod
func (d *Driver) OpMod(){
	prog := &Op{code:VMCODE_MOD}
	d.program = append(d.program, *prog)
	d.pc++
}
// equal
func (d *Driver) OpEqual(){
	prog := &Op{code:VMCODE_EQU}
	d.program = append(d.program, *prog)
	d.pc++
}
// greater than
func (d *Driver) OpGt(){
	prog := &Op{code:VMCODE_GT}
	d.program = append(d.program, *prog)
	d.pc++
}
// greater equal
func (d *Driver) OpGe(){
	prog := &Op{code:VMCODE_GE}
	d.program = append(d.program, *prog)
	d.pc++
}
// less than
func (d *Driver) OpLt(){
	prog := &Op{code:VMCODE_LT}
	d.program = append(d.program, *prog)
	d.pc++
}
// less equal
func (d *Driver) OpLe(){
	prog := &Op{code:VMCODE_LE}
	d.program = append(d.program, *prog)
	d.pc++
}
// not equal
func (d *Driver) OpNequ(){
	prog := &Op{code:VMCODE_NEQ}
	d.program = append(d.program, *prog)
	d.pc++
}
// and
func (d *Driver) OpAnd(){
	prog := &Op{code:VMCODE_AND}
	d.program = append(d.program, *prog)
	d.pc++
}
// or
func (d *Driver) OpOr(){
	prog := &Op{code:VMCODE_OR}
	d.program = append(d.program, *prog)
	d.pc++
}
// add_string
func (d *Driver) OpAddString(){
	prog := &Op{code:VMCODE_ADDSTRING}
	d.program = append(d.program, *prog)
	d.pc++
}
// not
func (d *Driver) OpNot(){
	prog := &Op{code:VMCODE_NOT}
	d.program = append(d.program, *prog)
	d.pc++
}
