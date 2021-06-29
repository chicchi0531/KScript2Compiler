package compiler

import "fmt"

type Driver struct{
	pc int
}

func (d *Driver) OpPushInteger(key int){
	fmt.Printf("%d:push_int %d\n", d.pc, key)
	d.pc++
}
func (d *Driver) OpPushFloat(key int){
	fmt.Printf("%d:push_float %d\n", d.pc, key)
	d.pc++
}
func (d *Driver) OpPushString(key int){
	fmt.Printf("%d:push_string %d\n", d.pc, key)
	d.pc++
}
func (d *Driver) OpPushValue(key int){
	fmt.Printf("%d:push_value %d\n", d.pc, key)
	d.pc++
}
func (d *Driver) OpPop(key int){
	fmt.Printf("%d:pop %d\n", d.pc, key)
	d.pc++
}

func (d *Driver) OpAdd(){
	fmt.Printf("%d:add\n", d.pc)
	d.pc++
}
func (d *Driver) OpSub(){
	fmt.Printf("%d:sub\n", d.pc)
	d.pc++
}
func (d *Driver) OpMul(){
	fmt.Printf("%d:mul\n", d.pc)
	d.pc++
}
func (d *Driver) OpDiv(){
	fmt.Printf("%d:div\n", d.pc)
	d.pc++
}
func (d *Driver) OpEqual(){
	fmt.Printf("%d:equ\n", d.pc)
	d.pc++
}
func (d *Driver) OpGt(){
	fmt.Printf("%d:gt\n", d.pc)
	d.pc++
}
func (d *Driver) OpGe(){
	fmt.Printf("%d:ge\n", d.pc)
	d.pc++
}
func (d *Driver) OpLt(){
	fmt.Printf("%d:lt\n", d.pc)
	d.pc++
}
func (d *Driver) OpLe(){
	fmt.Printf("%d:le\n", d.pc)
	d.pc++
}
func (d *Driver) OpNequ(){
	fmt.Printf("%d:nequ\n", d.pc)
	d.pc++
}
func (d *Driver) OpAnd(){
	fmt.Printf("%d:and\n", d.pc)
	d.pc++
}
func (d *Driver) OpOr(){
	fmt.Printf("%d:or\n", d.pc)
	d.pc++
}