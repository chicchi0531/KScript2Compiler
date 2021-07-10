package vm

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
	
	VMCODE_JMP
	VMCODE_CALL
	VMCODE_SYSCALL
	VMCODE_RETURN
)

func VMCODE_TOSTR (code int) string{
	switch code{
	case VMCODE_PUSHINT:return "PushInt"
	case VMCODE_PUSHFLOAT:return "PushFloat"
	case VMCODE_PUSHSTRING:return "PushString"
	case VMCODE_PUSHVALUE:return "PushValue"
	case VMCODE_POPVALUE:return "PopValue"

	case VMCODE_ADD:return "Add"
	case VMCODE_SUB:return "Sub"
	case VMCODE_MUL:return "Mul"
	case VMCODE_DIV:return "Div"
	case VMCODE_MOD:return "Mod"
	case VMCODE_EQU:return "Equ"
	case VMCODE_NEQ:return "Neq"
	case VMCODE_GT:return "Gt"
	case VMCODE_GE:return "Ge"
	case VMCODE_LT:return "Lt"
	case VMCODE_LE:return "Le"
	case VMCODE_NOT:return "Not"
	case VMCODE_AND:return "And"
	case VMCODE_OR:return "Or"
	case VMCODE_ADDSTRING:return "AddString"

	case VMCODE_JMP:return "Jmp"
	case VMCODE_CALL:return "Call"
	case VMCODE_SYSCALL:return "Syscall"
	case VMCODE_RETURN:return "Return"
}
	return "VMTOSTRに登録されていないコードが呼ばれました。"
}