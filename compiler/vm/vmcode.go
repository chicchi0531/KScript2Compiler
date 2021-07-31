package vm

const(
	VMCODE_PUSHINT = iota
	VMCODE_PUSHFLOAT
	VMCODE_PUSHSTRING
	VMCODE_PUSHVALUE
	VMCODE_POPVALUE
	VMCODE_POP

	VMCODE_ADD
	VMCODE_SUB
	VMCODE_MUL
	VMCODE_DIV

	VMCODE_FADD
	VMCODE_FSUB
	VMCODE_FMUL
	VMCODE_FDIV

	VMCODE_MOD
	VMCODE_INCR
	VMCODE_DECR
	VMCODE_EQU
	VMCODE_NEQ

	VMCODE_FGT
	VMCODE_FGE
	VMCODE_FLT
	VMCODE_FLE

	VMCODE_GT
	VMCODE_GE
	VMCODE_LT
	VMCODE_LE

	VMCODE_NOT
	VMCODE_FNOT
	VMCODE_AND
	VMCODE_OR
	VMCODE_ADDSTRING
	
	VMCODE_JMP
	VMCODE_JZE
	VMCODE_JNZ
	VMCODE_CALL
	VMCODE_SYSCALL
	VMCODE_RETURN
	VMCODE_RETURNV
	VMCODE_NOP

	VMCODE_DUMMYLABEL

	VMCODE_FLOATTABLE = 0xef //float table識別用のダミー命令
)

func VMCODE_TOSTR (code int) string{
	switch code{
	case VMCODE_PUSHINT:return "PushInt"
	case VMCODE_PUSHFLOAT:return "PushFloat"
	case VMCODE_PUSHSTRING:return "PushString"
	case VMCODE_PUSHVALUE:return "PushValue"
	case VMCODE_POPVALUE:return "PopValue"
	case VMCODE_POP:return "Pop"

	case VMCODE_ADD:return "Add"
	case VMCODE_SUB:return "Sub"
	case VMCODE_MUL:return "Mul"
	case VMCODE_DIV:return "Div"
	case VMCODE_MOD:return "Mod"
	case VMCODE_INCR:return "Incr"
	case VMCODE_DECR:return "Decr"
	case VMCODE_EQU:return "Equ"
	case VMCODE_NEQ:return "Neq"
	case VMCODE_GT:return "Gt"
	case VMCODE_GE:return "Ge"
	case VMCODE_LT:return "Lt"
	case VMCODE_LE:return "Le"
	case VMCODE_NOT:return "Not"
	case VMCODE_AND:return "And"
	case VMCODE_OR:return "Or"
	case VMCODE_FADD:return "FAdd"
	case VMCODE_FSUB:return "FSub"
	case VMCODE_FMUL:return "FMul"
	case VMCODE_FDIV:return "FDiv"
	case VMCODE_FGT:return "FGt"
	case VMCODE_FGE:return "FGe"
	case VMCODE_FLT:return "FLt"
	case VMCODE_FLE:return "FLe"
	case VMCODE_ADDSTRING:return "AddString"

	case VMCODE_JMP:return "Jmp"
	case VMCODE_JZE:return "Jze"
	case VMCODE_JNZ:return "Jnz"
	case VMCODE_CALL:return "Call"
	case VMCODE_SYSCALL:return "Syscall"
	case VMCODE_RETURN:return "Return"
	case VMCODE_RETURNV: return "ReturnV"
	case VMCODE_NOP: return "Nop"

	case VMCODE_DUMMYLABEL:return "Label"
}
	return "VMTOSTRに登録されていないコードが呼ばれました。"
}