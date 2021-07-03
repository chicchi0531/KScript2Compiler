package compiler

const(
	VMCODE_PUSHINT = iota
	VMCODE_PUSHFLOAT
	VMCODE_PUSHSTRING
	VMCODE_PUSHVALUE
	VMCODE_POPVALUE

	VMCODE_ADD
	VMCODE_SUB
	VMCODE_MUL
	VMCODE_DIV
	VMCODE_MOD
	VMCODE_EQU
	VMCODE_NEQ
	VMCODE_GT
	VMCODE_GE
	VMCODE_LT
	VMCODE_LE
	VMCODE_NOT
	VMCODE_AND
	VMCODE_OR
	VMCODE_ADDSTRING	
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
