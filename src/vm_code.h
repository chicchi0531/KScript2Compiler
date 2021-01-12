
// enum定義用
#ifdef VM_ENUMDEF
#define VMCODE0(code_, name_) code_,
#define VMCODE1(code_, name_) code_,
#endif

// opコード生成関数用
#ifdef VM_CREATE
#define VMCODE0(code_, name_) void name_() { program.push_back(VMCode(code_)); }
#define VMCODE1(code_, name_) void name_(int arg1) { program.push_back(VMCode(code_, arg1)); }
#endif

// 文字列変換
#ifdef VM_NAMETABLE
#define VMCODE0(code_, name_) #name_,
#define VMCODE1(code_, name_) #name_,
#endif

VMCODE1(VM_PUSHCONST, PushConst)
VMCODE1(VM_PUSHDOUBLE, PushDouble)
VMCODE1(VM_PUSHSTRING, PushString)
VMCODE1(VM_PUSHVALUE, PushValue)
VMCODE1(VM_PUSHLOCAL, PushLocal)
VMCODE1(VM_POPLOCAL, PopLocal)
VMCODE1(VM_POPVALUE, PopValue)
VMCODE0(VM_POP, OpPop)

VMCODE0(VM_CAST_ITOF, CastItoF)// cast int to float
VMCODE0(VM_CAST_FTOI, CastFtoI)// cast float to int

// int計算
VMCODE0(VM_INEG, OpINeg)
VMCODE0(VM_IEQU, OpIEqu)
VMCODE0(VM_INEQ, OpINeq)
VMCODE0(VM_IGT, OpIGt)
VMCODE0(VM_IGE, OpIGe)
VMCODE0(VM_ILE, OpILe)
VMCODE0(VM_ILT, OpILt)
VMCODE0(VM_IADD, OpIAdd)
VMCODE0(VM_ISUB, OpISub)
VMCODE0(VM_IMUL, OpIMul)
VMCODE0(VM_IDIV, OpIDiv)
VMCODE0(VM_IMOD, OpIMod)
VMCODE0(VM_LOGAND, OpLogAnd)
VMCODE0(VM_LOGOR, OpLogOr)

// float計算
VMCODE0(VM_FNEG, OpFNeg)
VMCODE0(VM_FEQU, OpFEqu)
VMCODE0(VM_FNEQ, OpFNeq)
VMCODE0(VM_FGT, OpFGt)
VMCODE0(VM_FGE, OpFGe)
VMCODE0(VM_FLE, OpFLe)
VMCODE0(VM_FLT, OpFLt)
VMCODE0(VM_FADD, OpFAdd)
VMCODE0(VM_FSUB, OpFSub)
VMCODE0(VM_FMUL, OpFMul)
VMCODE0(VM_FDIV, OpFDiv)

// string計算
VMCODE0(VM_STREQ, OpStrEq)
VMCODE0(VM_STRNE, OpStrNe)
VMCODE0(VM_STRADD, OpStrAdd)

VMCODE1(VM_JMP, OpJmp)
VMCODE1(VM_JZE, OpJZero)
VMCODE1(VM_JNZ, OpJNZero)
VMCODE0(VM_RETURN, OpJReturn)
VMCODE0(VM_RETURNV, OpJReturnV)
VMCODE1(VM_CALL, OpCall)
VMCODE1(VM_SYSCALL, OpSysCall)
VMCODE0(VM_NOVEL_MSG, OpNovelMsg)
VMCODE0(VM_NOVEL_NAME, OpNovelName)
VMCODE0(VM_NOVEL_NEWPAGE, OpNovelNewPage)
VMCODE0(VM_NOVEL_NEWLINE, OpNovelNewLine)
VMCODE0(VM_HALT, OpHalt)

#undef VMCODE0
#undef VMCODE1