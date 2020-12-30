#pragma once

#include "common.h"

namespace kscript2{namespace parser
{
	namespace x3 = boost::spirit::x3;

	using x3::symbols;

	// 変数、関数の戻り値の型
	enum
	{
		TYPE_INTEGER,
		TYPE_FLOAT,
		TYPE_STRING,
		TYPE_VOID,
	};

	// pragma命令の種類
	enum
	{
		PRAGMA_INCLUDE,
	};

	// 仮想マシン命令定義
	enum
	{
#define VM_ENUMDEF
#include "vm_code.h"
#undef VM_ENUMDEF
		VM_MAXCOMMAND	// dummy label
	};

	// OPEコード定義
	enum OPCODE
	{
		OP_NEG,
		OP_ADD,
		OP_SUB,
		OP_MUL,
		OP_DIV,
		OP_MOD,
		OP_AND,
		OP_OR,
		OP_LSHIFT,
		OP_RSHIFT,
		OP_LOGAND,
		OP_LOGOR,
		OP_EQ,
		OP_NE,
		OP_GT,
		OP_GE,
		OP_LT,
		OP_LE,
		OP_ASSIGN,
		OP_ADD_ASSIGN,
		OP_SUB_ASSIGN,
		OP_MUL_ASSIGN,
		OP_DIV_ASSIGN,
		OP_MOD_ASSIGN,
		OP_NUMBER,
		OP_IDENTIFIER,
		OP_STRING,
		OP_FUNCTION,
		OP_ARRAY
	};

#define STR(name_) #name_
	static const char* op_string
	{
		STR(OP_ADD)
		STR(OP_SUB)

	};

	//-------------------------------
	// Symbols
	//-------------------------------
	struct type_symbols_ : symbols<int>
	{
		type_symbols_()
		{
			add
				("int", TYPE_INTEGER)
				("string", TYPE_STRING)
				("float", TYPE_FLOAT)
				;
		}
	}static type_symbols;

	struct func_type_symbols_ : symbols<int>
	{
		func_type_symbols_()
		{
			add
				("int", TYPE_INTEGER)
				("string", TYPE_STRING)
				("float", TYPE_FLOAT)
				("void", TYPE_VOID)
				;
		}
	}static func_type_symbols;

	struct op_sign_symbols_ : symbols<int>
	{
		op_sign_symbols_()
		{
			add
				("-", OP_NEG)
				;
		}
	}static op_sign_symbols;

	struct op_mul_symbols_ : symbols<int>
	{
		op_mul_symbols_()
		{
			add
				("*", OP_MUL)
				("/", OP_DIV)
				("%", OP_MOD)
				;
		}
	}static op_mul_symbols;

	struct op_add_symbols_ : symbols<int>
	{
		op_add_symbols_()
		{
			add
				("+", OP_ADD)
				("-", OP_SUB)
				;

		}
	}static op_add_symbols;

	struct op_equ_symbols_ : symbols<int>
	{
		op_equ_symbols_()
		{
			add
				("==", OP_EQ)
				("!=", OP_NE)
				;

		}
	}static op_equ_symbols;

	struct op_relation_symbols_ : symbols<int>
	{
		op_relation_symbols_()
		{
			add
				(">=", OP_GE)
				(">", OP_GT)
				("<=", OP_LE)
				("<", OP_LT)
				;
		}
	}static op_relation_symbols;

	struct op_logand_symbols_ : symbols<int>
	{
		op_logand_symbols_()
		{
			add
				("&&", OP_LOGAND)
				;
		}
	}static op_logand_symbols;

	struct op_logor_symbols_ : symbols<int>
	{
		op_logor_symbols_()
		{
			add
				("||", OP_LOGOR)
				;
		}
	}static op_logor_symbols;

	struct op_assign_symbols_ : symbols<int>
	{
		op_assign_symbols_()
		{
			add
				("=", OP_ASSIGN)
				("+=", OP_ADD_ASSIGN)
				("-=", OP_SUB_ASSIGN)
				("*=", OP_MUL_ASSIGN)
				("/=", OP_DIV_ASSIGN)
				("%=", OP_MOD_ASSIGN)
				;

		}
	}static op_assign_symbols;

	// ジャンプ命令
	enum
	{
		JUMP_CONTINUE,
		JUMP_BREAK,
		JUMP_RETURN
	};
	struct op_jump_symbols_ : symbols<int>
	{
		op_jump_symbols_()
		{
			add
				("continue", JUMP_CONTINUE)
				("break", JUMP_BREAK)
				("return", JUMP_RETURN)
				;
		}
	}static op_jump_symbols;
	
	// ノベル命令
	enum
	{
		NOVEL_EMPTY,
		NOVEL_NEWLINE,
		NOVEL_NEWPAGE,
	};

	struct novel_new_page_ : symbols<int>
	{
		novel_new_page_()
		{
			add
				("」", NOVEL_NEWPAGE)
				("<p>", NOVEL_NEWPAGE)
				;
		}
	}static novel_new_page;

	struct novel_new_line_ : symbols<int>
	{
		novel_new_line_()
		{
			add
				("<n>", NOVEL_NEWLINE)
				;
		}
	}static novel_new_line;
	
	static symbols<> keywords =
	{
		"if",
		"else",
		"while",
		"for",
		"break",
		"continue",
		"return",
		"yield",
		"import"
	};
}}