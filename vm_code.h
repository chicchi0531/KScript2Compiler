
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
VMCODE1(VM_ALLOCSTACK, OpAllocStack)
VMCODE0(VM_POP, OpPop)
VMCODE0(VM_NEG, OpNeg)
VMCODE0(VM_EQU, OpEqu)
VMCODE0(VM_NEQ, OpNeq)
VMCODE0(VM_GT, OpGt)
VMCODE0(VM_GE, OpGe)
VMCODE0(VM_LE, OpLe)
VMCODE0(VM_LT, OpLt)
VMCODE0(VM_LOGAND, OpLogAnd)
VMCODE0(VM_LOGOR, OpLogOr)
VMCODE0(VM_ADD, OpAdd)
VMCODE0(VM_SUB, OpSub)
VMCODE0(VM_MUL, OpMul)
VMCODE0(VM_DIV, OpDiv)
VMCODE0(VM_MOD, OpMod)
VMCODE0(VM_STREQ, OpStrEq)
VMCODE0(VM_STRNE, OpStrNe)
VMCODE0(VM_STRGT, OpStrGt)
VMCODE0(VM_STRGE, OpStrGe)
VMCODE0(VM_STRLT, OpStrLt)
VMCODE0(VM_STRLE, OpStrLe)
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