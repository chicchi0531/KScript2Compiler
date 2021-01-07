#include <fstream>
#include <iostream>
#include <string>
#include <iomanip>
#include <iterator>
#include <algorithm>
#include <sstream>

#include "compiler.h"
#include "grammar.h"
#include "ast_analyzer.h"
#include "config.h"
#include "error_handler.h"

namespace x3 = boost::spirit::x3;
using namespace kscript2;

bool compiler::compile(const std::string& filepath)
{
	filepath_ = std::filesystem::path(filepath);
	auto source = FileLoad(filepath);

	std::stringstream out;

	using kscript2::parser::iterator_type;
	iterator_type iter(source.begin());
	iterator_type const end(source.end());

	// Our AST
	kscript2::ast::unit ast;

	// Our error handler
	using boost::spirit::x3::with;
	using kscript2::parser::error_handler_type;
	using kscript2::parser::error_handler_tag;
	error_handler_type error_handler(iter, end, out, filepath); // Our error handler

	// Our parser
	auto const parser =
		// we pass our error handler to the parser so we can access
		// it later on in our on_error and on_sucess handlers
		with<error_handler_tag>(std::ref(error_handler))
		[
			kscript2::unit()
		];

	std::cout << "■ 構文解析開始==================" << std::endl;

	// Go forth and parse!
	namespace x3 = boost::spirit::x3;
	bool success = phrase_parse(iter, end, parser, kscript2::parser::comment::skipper, ast);

	try
	{
		if (success)
		{
			if (iter != end)
			{
				error_handler(iter, "★ 構文解析エラー：ファイルの終端に到達する前に解析が終了しました。");
				return false;
			}
			else
			{
				std::cout << "■ 完了=======================" << std::endl;

				std::cout << "■ 意味解析開始===============" << std::endl;
				positions = error_handler.get_position_cache();
				auto analyzer = kscript2::ast::ast_analyzer(*this, positions);
				analyzer(ast);

				// error check
				if (error_count > 0) throw CompilerErrorException("★ 意味解析エラー：" + std::to_string(error_count) + "個のエラーが見つかりました。");

				std::cout << "■ 完了=======================" << std::endl;
			}
		}
		else
		{
			// print error
			std::cout << out.str() << std::endl;
			throw CompilerErrorException("★ 構文解析エラー");
		}
	}
	catch (CompilerErrorException e)
	{
		error_handler(iter, e.error_message);
		return false;
	}

	// 出力
	CreateData(1);

	return true;
}

void compiler::error(const std::string& message, const x3::position_tagged& ast)
{
	auto first = positions.first();
	auto pos = positions.position_of(ast);

	//行を推定
	int row=1,column=1;
	while(first != pos.end())
	{
		if(*first == '\n')
		{
			row++;
			column = 0;
		}
		column++;
		first++;
	}

	std::cerr << "[" << row << ":" << column << "] " << message << std::endl;
	std::cerr << "> " << std::string(pos.begin(), pos.end()) << std::endl;
	error_count++;
}

// 内部関数の定義
bool compiler::add_function(int index, int type, const char* name, const char* args)
{
	FunctionTag func(type);
	if (!func.SetArgs(args))
		return false;

	//定義済みにする
	func.SetDeclaration();

	//システム関数としてセットする
	func.SetSystem();

	//
	func.SetIndex(index);
	if (functions.add(name, func) == 0)
	{
		return false;
	}
	return true;
}

// 外部変数の定義
struct define_value
{
	compiler* c_;
	int type_;
	define_value(compiler* c, int type):c_(c), type_(type)
	{}

	void operator()(const ast::declarator& node) const
	{
		c_->AddValue(type_, node.identifier_.name, node);
	}
};
void compiler::DefineValue(int type, const std::vector<ast::declarator>& nodes)
{
	for (auto const& node : nodes)
	{
		define_value(this, type)(node);
	}
}

// 関数宣言
void compiler::DefineFunction(int type, const std::string& name, const ast::arg_def_list& args, const ast::function_pre_def& ast)
{
	const FunctionTag* tag = functions.find(name);
	if (tag)
	{
		if (!tag->CheckArgList(args))
		{
			error("関数 " + name + " に異なる型の引数が指定されました。", ast);
			return;
		}
		else
		{
			FunctionTag func(type);

			//引数設定
			func.SetArgs(args);

			//宣言済みに変更
			func.SetDeclaration();

			//開始地点のラベル登録
			func.SetIndex(MakeLabel());

			if (functions.add(name, func) == 0)
			{
				error("内部エラー：関数テーブルに登録できません", ast);
			}
		}
	}
}

// 関数定義
struct add_value
{
	compiler* c_;
	ValueTable& values_;
	mutable int addr_;
	add_value(compiler* c, ValueTable& values)
		: c_(c), values_(values), addr_(-4)
	{}

	void operator()(const ast::declarator& arg) const
	{
		if (!values_.add_arg(arg.type, arg.identifier_.name, addr_))
		{
			c_->error("引数" + arg.identifier_.name + "は既に登録されています。", arg);
		}
		addr_--;
	}
};

void compiler::AddFunction(int type, const std::string& name, const ast::arg_def_list& args, const ast::statements& block, const ast::function_def& ast)
{
	FunctionTag* tag = functions.find(name);
	if (tag)
	{
		if (tag->IsDefinition())
		{
			error("関数 " + name + " は既に定義されています。",ast );
			return;
		}
		if (tag->IsDeclaration() && !tag->CheckArgList(args))
		{
			error("関数 " + name + " に異なる方の引数が指定されています。", ast);
			return;
		}
		tag->SetDefinition();
	}
	else
	{
		FunctionTag func(type);

		//引数を設定
		func.SetArgs(args);
		
		//定義済みに設定
		func.SetDefinition();

		//ラベル登録
		func.SetIndex(MakeLabel());
		tag = functions.add(name, func);
		if (tag == nullptr)
		{
			error("内部エラー：関数テーブルに登録できませんでした。", ast);
		}
	}

	// 処理中の関数情報を記録
	current_function_name = name;
	current_function_type = type;

	// エントリポイントにラベルを置く
	SetLabel(tag->GetIndex());

	// 変数スタックを増やす
	BlockIn();

	// 引数リストを構築
	std::for_each(args.rbegin(), args.rend(), add_value(this, variables.back()));

	// 文があれば、分を登録
	auto a = ast::ast_analyzer(*this, positions);
	a(block);

	// 戻り値の処理
	const VMCode& code = program.back();
	if (type == TYPE_VOID)
	{
		// 関数の最期にreturnがなければ追加
		if (code.op_ != VM_RETURN)
			OpJReturn();
	}
	else
	{
		if (code.op_ != VM_RETURNV)
		{
			error("関数　" + name + "　の最後にreturn文が必要です。",ast);
		}
	}

	// 変数スタックを減らす
	BlockOut();

	// 処理中の関数名を消去
	current_function_name.clear();
}

void compiler::AddValue(int type, const std::string& name, const ast::declarator& ast)
{
	ValueTable& values = variables.back();
	if (!values.add(type, name, 1))
	{
		error("変数 " + name + " は既に定義済みのシンボルです。",ast);
	}
}

// ラベル生成
int compiler::MakeLabel()
{
	int index = (int)labels.size();
	labels.push_back(Label(index));
	return index;
}

// ラベルのダミーコマンドを登録
void compiler::SetLabel(int label)
{
	program.push_back(VMCode(VM_MAXCOMMAND, label));
}

// 文字列定数をpush
void compiler::PushString(const std::wstring& str)
{
	PushString((int)text_table.size());
	text_table.push_back(str);
}

// 小数点をpush
void compiler::PushDouble(double value)
{
	PushDouble((int)double_table.size());
	double_table.push_back(value);
}

// break文に対応したJmpコマンド生成
bool compiler::JmpBreakLabel()
{
	if (break_index < 0)
		return false;
	OpJmp(break_index);
	return true;
}

// continue文に対応したJmpコマンド生成
bool compiler::JmpContinueLabel()
{
	if (continue_index < 0)
		return false;
	OpJmp(continue_index);
	return true;
}

// ブロック内で、新しい変数セットに変数を登録する
void compiler::BlockIn()
{
	int start_addr = 0;
	if (variables.size() > 1)
	{
		start_addr = variables.back().size();
	}
	variables.push_back(ValueTable(start_addr));
}

// ブロックの終了で変数スコープを削除
void compiler::BlockOut()
{
	variables.pop_back();
}

// ローカル変数用にスタックを確保
void compiler::AllocStack()
{
	OpAllocStack(variables.back().size());
}

// ラベル解決
// 1.アドレスの生成
// 2.ダミーのラベルコマンドがあったアドレスを、ラベルテーブルに登録
// 3.jmpコマンドの飛び先をラベルテーブルに登録されたアドレスで置き換える

int compiler::LabelSetting()
{
	// アドレス計算
	int pos = 0;
	for (auto const& s : program)
	{
		if (s.op_ == VM_MAXCOMMAND)
		{
			labels[s.arg1_].pos_ = pos;
		}
		else
		{
			pos += s.size_;
		}
	}
	// ジャンプアドレス設定
	for (auto& s : program)
	{
		switch (s.op_)
		{
		case VM_JMP:
		case VM_JZE:
		case VM_JNZ:
		case VM_CALL:
			s.arg1_ = labels[s.arg1_].pos_;
			break;
		}
	}
	return pos;
}

// include文
void compiler::Include(const std::string& filepath)
{
	// コンパイル中に同インスタンスでコンパイルすることでincludeを実装する

}

// バイナリデータ生成
bool compiler::CreateData(int code_size)
{
	auto create_utf8_bom = [](const auto path)
	{
		std::ofstream ofs(path);
		const unsigned char bom[] = {0xEF,0xBB,0xBF};
		ofs.write(reinterpret_cast<const char*>(bom), sizeof(bom));
		return ofs;
	};

	// UTF8-BOMでデータを書き込む
	auto path = filepath_.replace_extension("ksobj");
	auto ofs = create_utf8_bom(path);

	// json形式で書き出し（汚いが手書き）
	std::streampos fp = ofs.tellp();
	ofs << "{" << std::endl;

	// program section
	ofs << "\"Program\" : [";
	for(auto const& p : program)
	{
		ofs << "[" << std::to_string(p.op_) << "," << std::to_string(p.arg1_) << "]";
		fp = ofs.tellp();
		ofs << ",";
	}
	ofs.seekp(fp); // 要素の最後の,を消す
	ofs << "]," << std::endl;

	// text section
	ofs << "\"Text\" : [";
	for(auto const& t : text_table)
	{
		std::string str = std::string(t.begin(), t.end());
		ofs << "\"" << str << "\"";
		fp = ofs.tellp();
		ofs << ",";
	}
	ofs.seekp(fp); // 要素の最後の,を消す
	ofs << "]," << std::endl;

	// double section
	ofs << "\"Double\" : [";
	for(auto d : double_table)
	{
		ofs << std::to_string(d);
		fp = ofs.tellp();
		ofs << ",";
	}
	ofs.seekp(fp);
	ofs << "]" << std::endl;

	// end tag
	ofs << "}" << std::endl;

#ifdef _DEBUG
	debug_dump();
#endif

	return true;
}

// デバッグダンプ
#ifdef _DEBUG
void compiler::debug_dump()
{
	std::cout << "---variables---" << std::endl;
	size_t vsize = variables.size();
	std::cout << "value stack = " << vsize << std::endl;
	for (size_t i = 0; i < vsize; i++)
	{
		variables[i].dump();
	}
	std::cout << "---code---" << std::endl;

	const char* op_name[] =
	{
#define VM_NAMETABLE
#include "vm_code.h"
#undef VM_NAMETABLE
		"LABEL",
	};

	int pos = 0;
	size_t size = program.size();
	for (size_t i = 0; i < size; i++)
	{
		std::cout << std::setw(6) << pos << ": " << op_name[program[i].op_];
		if (program[i].size_ > 1)
		{
			std::cout << ", " << program[i].arg1_;
		}
		std::cout << std::endl;

		if (program[i].op_ != VM_MAXCOMMAND)
		{
			pos += program[i].size_;
		}
	}
}
#endif