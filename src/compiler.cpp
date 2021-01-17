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

bool compiler::compile(const fs::path &in_filepath, const fs::path& out_filepath)
{
	// グローバル変数用、変数テーブルをセット
	variables.push_back(ValueTable());
	variables[0].set_global();

	// 先頭はHalt命令にしておく
	OpHalt();

	// コンパイル
	if (!parse(in_filepath))
		return false;

	// ラベル設定
	int code_size = LabelSetting();

	// バイナリファイル出力
	if (!CreateData(out_filepath))
		return false;

	return true;
}

// 外部変数の定義
struct define_value
{
	compiler *c_;
	int type_;
	define_value(compiler *c, int type) : c_(c), type_(type)
	{
	}

	void operator()(const ast::declarator &node) const
	{
		c_->AddValue(type_, node.identifier_.name, node);
	}
};
void compiler::DefineValue(int type, const std::vector<ast::declarator> &nodes)
{
	for (auto const &node : nodes)
	{
		define_value(this, type)(node);
	}
}

// 関数宣言
void compiler::DeclFunction(int attr, int type, const std::string &name, const ast::arg_def_list &args, const ast::function_pre_def &ast)
{
	const FunctionTag *tag = functions.find(name);
	if (tag)
	{
		if (!tag->CheckArgList(args))
		{
			error("関数 " + name + " に異なる型の引数が指定されました。", ast);
			return;
		}
	}
	else
	{
		FunctionTag func(type);

		//引数設定
		func.SetArgs(args);

		//宣言済みに変更
		func.SetDeclaration();

		//システムコール設定
		if (attr == ATTR_SYSTEM)
		{
			func.SetSystem();
			func.SetIndex(system_function_num++); //システムコールは専用のインデックスを割り当てる
			systemcalls.push_back(name);
		}
		else
		{
			//開始地点のラベル登録
			func.SetIndex(MakeLabel());
		}

		if (functions.add(name, func) == 0)
		{
			error("内部エラー：関数テーブルに登録できません", ast);
		}
	}
}

// 関数定義
struct add_value
{
	compiler *c_;
	ValueTable &values_;
	mutable int addr_;
	add_value(compiler *c, ValueTable &values)
		: c_(c), values_(values), addr_(-4)
	{
	}

	void operator()(const ast::declarator &arg) const
	{
		if (!values_.add_arg(arg.type, arg.identifier_.name, addr_))
		{
			c_->error("引数" + arg.identifier_.name + "は既に登録されています。", arg);
		}
		addr_--;
	}
};

void compiler::DefineFunction(int type, const std::string &name, const ast::arg_def_list &args, const ast::statements &block, const ast::function_def &ast)
{
	FunctionTag *tag = functions.find(name);
	if (tag)
	{
		if (tag->IsDefinition())
		{
			error("関数 " + name + " は既に定義されています。", ast);
			return;
		}
		if (tag->IsDeclaration() && !tag->CheckArgList(args))
		{
			error("関数 " + name + " に異なる型の引数が指定されています。", ast);
			return;
		}
		if (tag->IsDeclaration() && tag->IsSystem())
		{
			error("システムコールを定義することはできません。システムコールは宣言のみ記述すると、VM側でC#コードとリンクします", ast);
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
			return;
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
	const VMCode &code = program.back();
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
			error("関数　" + name + "　の最後にreturn文が必要です。", ast);
		}
	}

	// 変数スタックを減らす
	BlockOut();

	// 処理中の関数名を消去
	current_function_name.clear();
}

void compiler::AddValue(int type, const std::string &name, const ast::declarator &ast)
{
	ValueTable &values = variables.back();
	if (!values.add(type, name, 1))
	{
		error("変数 " + name + " は既に定義済みのシンボルです。", ast);
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
void compiler::PushString(const std::wstring &str)
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

// ラベル解決
// 1.アドレスの生成
// 2.ダミーのラベルコマンドがあったアドレスを、ラベルテーブルに登録
// 3.jmpコマンドの飛び先をラベルテーブルに登録されたアドレスで置き換える

int compiler::LabelSetting()
{
	// アドレス計算
	int pos = 0;
	for (auto const &s : program)
	{
		if (s.op_ == VM_MAXCOMMAND)
		{
			labels[s.arg1_].pos_ = pos;
		}
		else
		{
			pos++;
		}
	}
	// ジャンプアドレス設定
	for (auto &s : program)
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
void compiler::Include(const std::string &filepath, const x3::position_tagged& ast)
{
	// コンパイル中に同インスタンスでコンパイルすることでincludeを実装する
	auto current_dir = filepathes_.back().parent_path();
	fs::path import_path(current_dir.string() + "/" + filepath);
	if(!parse(import_path, ast))
	{
		throw CompilerErrorException("import命令に失敗しました。");
	}
}

// バイナリデータ生成
bool compiler::CreateData(const fs::path &path)
{
	//ディレクトリの構成
	auto dir = path.parent_path();
	if(!fs::exists(dir) && dir != "")
	{
		fs::create_directory(dir);
	}

	auto create_utf8_bom = [](const auto path) {
		std::ofstream ofs(path);
		const unsigned char bom[] = {0xEF, 0xBB, 0xBF};
		ofs.write(reinterpret_cast<const char *>(bom), sizeof(bom));
		return ofs;
	};

	// UTF8-BOMでデータを書き込む
	auto ofs = create_utf8_bom(path);
	if (!ofs)
	{
		std::cerr << "出力ファイルが作成できませんでした。 : " << path << std::endl;
		return false;
	}

	// json形式で書き出し（汚いが手書き）
	std::streampos fp = ofs.tellp();
	ofs << "{" << std::endl;

	// program section
	ofs << "\"Program\" : [";
	fp = ofs.tellp();
	for (auto const &p : program)
	{
		// Label命令はスキップ
		if(p.op_ == VM_MAXCOMMAND) continue;

		ofs << "{" << "\"op\":" << std::to_string(p.op_) << ","
			<< "\"arg\":" << std::to_string(p.arg1_) << "}";
		fp = ofs.tellp();
		ofs << ",";
	}
	ofs.seekp(fp); // 要素の最後の,を消す
	ofs << "]," << std::endl;

	// エントリポイント
	ofs << "\"EntryPoint\":";
	auto tag = functions.find("main");
	if (!tag)
	{
		std::cerr << "エントリポイントが見つかりません。main関数を定義してください。" << std::endl;
		return false;
	}
	auto entrypoint = labels[tag->GetIndex()].pos_;
	ofs << std::to_string(entrypoint) << "," << std::endl;

	// text section
	ofs << "\"Text\" : [";
	fp = ofs.tellp();
	for (auto const &t : text_table)
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
	fp = ofs.tellp();
	for (auto d : double_table)
	{
		ofs << std::to_string(d);
		fp = ofs.tellp();
		ofs << ",";
	}
	ofs.seekp(fp);
	ofs << "]," << std::endl;

	// systemcall section
	ofs << "\"SystemCall\" : [";
	fp = ofs.tellp();
	for (auto &fname : systemcalls)
	{
		auto tag = functions.find(fname);
		if (tag)
		{
			ofs << "{"
				<< "\"name\":" << "\"" << fname << "\"" << ","
				<< "\"argnum\":"<< std::to_string(tag->ArgSize()) << ","
				<< "\"argtypes\":" << "\"";
			for (int i = 0; i < tag->ArgSize(); i++)
			{
				switch (tag->GetArg(i))
				{
				case TYPE_INTEGER:
					ofs << "i";
					break;
				case TYPE_STRING:
					ofs << "s";
					break;
				case TYPE_FLOAT:
					ofs << "d";
					break;
				}
			}
			ofs << "\"";
			ofs << "}";
			fp = ofs.tellp();
			ofs << ",";
		}
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

//###########################
// error handling
//###########################
void err_msg_output(std::stringstream &out, const std::string &msg)
{
	out << msg << std::endl;
}
void err_output(const fs::path &filepath, std::stringstream &out, const position_cache &positions, const std::string &message, const x3::position_tagged &ast)
{
	auto first = positions.first();

	//位置が仕込まれていない場合は、メッセージだけ出力
	if (ast.id_first == -1)
	{
		err_msg_output(out, message);
		return;
	}

	boost::iterator_range<iterator_type> pos;
	pos = positions.position_of(ast);

	//行を推定
	int row = 1, column = 1;
	while (first != pos.end())
	{
		if (*first == '\n')
		{
			row++;
			column = 0;
		}
		column++;
		first++;
	}

	out << filepath << " [" << row << ":" << column << "] " << message << std::endl;
	out << "> " << std::string(pos.begin(), pos.end()) << std::endl;
}
void compiler::error(const std::string &message, const x3::position_tagged &ast)
{
	std::stringstream out;

	//標準エラー出力に表示
	std::cerr << "[Error] : ";
	err_output(filepathes_.back(), out, positions, message, ast);
	std::cerr << out.str() << std::endl;

	error_count++;
}
void compiler::error(const std::string &msg)
{
	std::stringstream out;

	//標準エラー出力に表示
	std::cerr << "[Error] : ";
	err_msg_output(out, msg);
	std::cerr << out.str() << std::endl;

	error_count++;
}
void compiler::warning(const std::string &message, const x3::position_tagged &ast)
{
	std::stringstream out;

	//標準エラー出力に表示
	std::cerr << "[Warning] : ";
	err_output(filepathes_.back(), out, positions, message, ast);
	std::cerr << out.str() << std::endl;

	warning_count++;
}
void compiler::warning(const std::string &msg)
{
	std::stringstream out;

	//標準エラー出力に表示
	std::cerr << "[Warning] : ";
	err_msg_output(out, msg);
	std::cerr << out.str() << std::endl;

	warning_count++;
}
void compiler::info(const std::string &message, const x3::position_tagged &ast)
{
	std::stringstream out;

	// 標準出力に表示
	std::cout << "[Info] : ";
	err_output(filepathes_.back(), out, positions, message, ast);
	std::cout << out.str() << std::endl;
}
void compiler::info(const std::string &msg)
{
	std::stringstream out;

	// 標準出力に表示
	std::cout << "[Info] : ";
	err_msg_output(out, msg);
	std::cout << out.str() << std::endl;
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

	const char *op_name[] =
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
		std::cout << ", " << program[i].arg1_;
		std::cout << std::endl;

		if (program[i].op_ != VM_MAXCOMMAND)
		{
			pos++;
		}
	}
}
#endif

bool compiler::parse(const fs::path &path, const x3::position_tagged& called_pos)
{
	// すでにコンパイル済みの場合は破棄
	auto fp_iter = filepathes_.begin();
	auto fp_end = filepathes_.end();
	if (std::find(fp_iter, fp_end, path) != fp_end)
	{
		auto msg = "import file:" + path.string() + "はすでにimport済みです。コンパイルをスキップします。";
		info(msg, called_pos);
		
		return true;
	}

	try
	{
		//スクリプトファイルの読み込み
		filepathes_.push_back(path);
		auto source = FileLoad(path);

		using parser::iterator_type;
		iterator_type iter(source.begin());
		iterator_type const end(source.end());

		// 構文解析用 構文木
		ast::unit ast;

		// 構文解析用エラーハンドラ
		using parser::error_handler_tag;
		using parser::error_handler_type;
		std::stringstream out;
		error_handler_type error_handler(iter, end, out, path.string()); // Our error handler

		// パーサー定義
		auto const parser =
			// we pass our error handler to the parser so we can access
			// it later on in our on_error and on_sucess handlers
			x3::with<error_handler_tag>(std::ref(error_handler))
				[kscript2::unit()];

		// ###################
		// 構文解析
		// ###################
		info("■ 構文解析開始==================");
		if(!phrase_parse(iter, end, parser, parser::comment::skipper, ast))
		{
			// print error
			error(out.str(), called_pos);
			throw CompilerErrorException("★ 構文解析エラー：エラーメッセージを確認してください。");
		}

		// 予期せぬエラー処理
		if (iter != end)
		{
			error_handler(iter, "ファイルの終端に到達する前に解析が終了しました。");
			error(out.str(), called_pos);
			throw CompilerErrorException("★ 構文解析エラー：エラーメッセージを確認してください。");
		}
		info("■ 完了=======================");

		// ###################
		// 意味解析
		// ###################
		info("■ 意味解析開始===============");
		positions = error_handler.get_position_cache();
		auto analyzer = ast::ast_analyzer(*this, positions);
		analyzer(ast);

		// error check
		if (error_count > 0)
			throw CompilerErrorException("★ 意味解析エラー：" + std::to_string(error_count) + "個のエラーが見つかりました。");

		info("■ 完了=======================");
	}
	catch (CompilerErrorException e)
	{
		error(e.error_message);
		return false;
	}

	return true;
}