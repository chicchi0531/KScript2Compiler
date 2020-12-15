#pragma once

#include "common.h"

namespace kscript2{namespace parser
{
	namespace x3 = boost::spirit::x3;

	using x3::symbols;

	// ïœêîÅAä÷êîÇÃñﬂÇËílÇÃå^
	enum
	{
		TYPE_INTEGER,
		TYPE_FLOAT,
		TYPE_STRING,
		TYPE_VOID,
	};

	// pragmañΩóﬂÇÃéÌóﬁ
	enum
	{
		PRAGMA_INCLUDE,
	};

	// ââéZéq
	enum
	{
		OP_ADD,
		OP_SUB,
		OP_MUL,
		OP_DIV,
		OP_MOD,
		OP_EQU,
		OP_NEQ,
		OP_GE,
		OP_GT,
		OP_LE,
		OP_LT,
		OP_ASSIGN,
		OP_ADDEQU,
		OP_SUBEQU,
		OP_MULEQU,
		OP_DIVEQU,
		OP_MODEQU,
		OP_LOGAND,
		OP_LOGOR,
		OP_NEG,
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
	}type_symbols;

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
	}func_type_symbols;

	struct op_sign_symbols_ : symbols<int>
	{
		op_sign_symbols_()
		{
			add
			("-", OP_NEG)
				;
		}
	}op_sign_symbols;

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
	}op_mul_symbols;

	struct op_add_symbols_ : symbols<int>
	{
		op_add_symbols_()
		{
			add
			("+", OP_ADD)
				("-", OP_SUB)
				;

		}
	}op_add_symbols;

	struct op_equ_symbols_ : symbols<int>
	{
		op_equ_symbols_()
		{
			add
			("==", OP_EQU)
				("!=", OP_NEQ)
				;

		}
	}op_equ_symbols;

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
	}op_relation_symbols;

	struct op_logand_symbols_ : symbols<int>
	{
		op_logand_symbols_()
		{
			add
			("&&", OP_LOGAND)
				;
		}
	}op_logand_symbols;

	struct op_logor_symbols_ : symbols<int>
	{
		op_logor_symbols_()
		{
			add
			("||", OP_LOGOR)
				;
		}
	}op_logor_symbols;

	struct op_assign_symbols_ : symbols<int>
	{
		op_assign_symbols_()
		{
			add
			("=", OP_ASSIGN)
				("+=", OP_ADDEQU)
				("-=", OP_SUBEQU)
				("*=", OP_MULEQU)
				("/=", OP_DIVEQU)
				("%=", OP_MODEQU)
				;

		}
	}op_assign_symbols;

	symbols<> keywords = 
	{
		"if",
		"else",
		"while",
		"for",
		"break",
		"continue",
		"return",
		"yield"
	};


}}