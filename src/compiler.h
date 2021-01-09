#pragma once

#include <string>
#include <vector>
#include <map>
#include <iostream>
#include <exception>
#include <filesystem>

#include <boost/filesystem.hpp>

#include "ast.h"
#include "symbols.h"
#include "config.h"

namespace kscript2
{
	using namespace parser;
	namespace fs = std::filesystem;

	// エラーオブジェクト
	class CompilerErrorException : std::exception
	{
	public:
		CompilerErrorException(std::string str) : error_message(str){}
		std::string error_message;
	};

	// 命令コード
	class VMCode
	{
	public:
		int op_;
		int arg1_;

	public:
		VMCode(int op)
			: op_(op), arg1_(0)
		{}
		VMCode(int op, int arg1)
			: op_(op), arg1_(arg1)
		{}

		/*int* Get(int* p) const
		{
			// ラベル以外を取得
			// ラベルは命令として載せないのでスキップ
			if (op_ != VM_MAXCOMMAND)
			{
				*p++ = op_;
				if (size_ > 1)
				{
					*(int*)p = arg1_;
					p += 4;
				}
			}
			return p;
		}*/
	};

	//ラベル
	class Label
	{
	public:
		int index_;
		int pos_;

	public:
		Label(int index)
			: index_(index), pos_(0)
		{
		}
	};

	// 変数テーブル
	class ValueTag
	{
	public:
		int addr_;
		int type_;
		int size_;
		bool global_;

	public:
		ValueTag(): addr_(-1), type_(TYPE_INTEGER), size_(1), global_(false)
		{}
		ValueTag(int addr, int type, int size, bool global)
			: addr_(addr), type_(type), size_(size), global_(global)
		{}
	};
	class ValueTable
	{
	private:
		typedef std::map<std::string, ValueTag>::iterator iter;
		typedef std::map<std::string, ValueTag>::const_iterator const_iter;
		std::map<std::string, ValueTag> variables_;
		int addr_;
		bool global_;

	public:
		ValueTable(int start_addr=0): addr_(start_addr), global_(false)
		{}

		// スコープをグローバルに
		void set_global() { global_ = true; }

		// 変数追加
		bool add(int type, const std::string& name, int size = 1)
		{
			std::pair<iter, bool> result = variables_.insert(
				{ name, ValueTag(addr_, type, size, global_) }
			);
			
			if (result.second)
			{
				addr_ += size;
				return true;
			}
			return false;
		}

		// 変数検索
		const ValueTag* find(const std::string& name) const
		{
			const_iter it = variables_.find(name);
			if (it != variables_.end())
				return &it->second;
			return nullptr;
		}

		// 引数追加
		bool add_arg(int type, const std::string& name, int addr)
		{
			std::pair<iter, bool> result = variables_.insert(
				{ name, ValueTag(addr, type, 1, false) }
			);
			return result.second;
		}

		// 現在のテーブルサイズ
		int size() const { return addr_; }

		// テーブルクリア
		void clear()
		{
			variables_.clear();
			addr_ = 0;
		}

#ifdef _DEBUG
		void dump() const
		{
			std::cout << "--------------- value ---------------" << std::endl;
			for (auto const& it : variables_)
			{
				std::cout << it.first << ", addr = " << it.second.addr_
					<< ", type = " << it.second.type_
					<< ", size = " << it.second.size_
					<< ", global = " << it.second.global_ << std::endl;
			}
		}
#endif
	};

	class FunctionTag
	{
	private:
		enum
		{
			flag_declaration = 1 << 0,
			flag_definition = 1 << 1,
			flag_system = 1 << 2,
		};

		int type_;
		int flags_;
		int index_;
		std::vector<int> args_;

	public:
		FunctionTag(int type) : type_(type), flags_(0), index_(0) {}

		void SetArg(int type)
		{
			args_.push_back((char)type);
		}

		void SetArgs(const ast::arg_def_list& args)
		{
			for (auto const& arg : args)
			{
				args_.push_back(arg.type);
			}
		}
		bool SetArgs(const char* args)
		{
			if (args)
			{
				for (int i = 0; args[i] != 0; i++)
				{
					switch (args[i])
					{
					case 'I': case 'i':
						args_.push_back(TYPE_INTEGER);
						break;
					case 'S': case 's':
						args_.push_back(TYPE_STRING);
						break;
					default:
						return false;
					}
				}
			}
			return true;
		}
		void SetArgs(const std::vector<int>& args)
		{
			for (auto const& arg : args)
			{
				args_.push_back(arg);
			}
		}

		bool CheckArgList(const ast::arg_def_list& args) const
		{
			// 引数がない場合
			if (args.empty())
				return args_.empty();

			// 引数の個数が異なる場合
			if (args.size() != args_.size())
				return false;

			// 全引数の型をチェック
			for (size_t i=0; i<args_.size(); i++)
			{
				if (args[i].type != args_[i])
					return false;
			}
			return true;
		}

		bool CheckArgList(const std::vector<int>& args) const
		{
			// 引数がない場合
			if (args.empty())
				return args_.empty();

			// 引数の個数が異なる
			if (args.size() != args_.size())
				return false;

			// 全引数の型チェック
			for (size_t i = 0; i < args_.size(); i++)
			{
				if (args[i] != args_[i])
					return false;
			}
			return true;
		}

		int GetArg(int index) const
		{
			return args_[index];
		}

		int ArgSize() const { return (int)args_.size(); }
		
		void SetIndex(int index) { index_ = index; }
		void SetDeclaration() { flags_ |= flag_declaration; }
		void SetDefinition() { flags_ |= flag_definition; }
		void SetSystem() { flags_ |= flag_system; }

		int GetIndex() const { return index_; }
		int GetType() const { return type_; }
		bool IsDeclaration() const { return (flags_ & flag_declaration) != 0; }
		bool IsDefinition() const { return(flags_ & flag_definition) != 0; }
		bool IsSystem() const { return (flags_ & flag_system) != 0; }

	};

	class FunctionTable
	{
	private:
		typedef std::map<std::string, FunctionTag>::iterator iter;
		typedef std::map<std::string, FunctionTag>::const_iterator const_iter;

		std::map<std::string, FunctionTag> functions_;

	public:
		FunctionTable(){}

		FunctionTag* add(const std::string& name, const FunctionTag& tag)
		{
			auto result = functions_.insert({ name, tag });
			if (result.second)
				return &result.first->second;
			return nullptr;
		}

		const FunctionTag* find(const std::string& name) const
		{
			const_iter it = functions_.find(name);
			if (it != functions_.end())
				return &it->second;
			return nullptr;
		}

		FunctionTag* find(const std::string& name)
		{
			iter it = functions_.find(name);
			if (it != functions_.end())
			{
				return &it->second;
			}
			return nullptr;
		}

		void clear()
		{
			functions_.clear();
		}

	};

	class compiler
	{
	private:
		FunctionTable functions;
		std::vector<std::string> systemcalls;
		std::vector<ValueTable> variables;
		std::vector<VMCode> program;
		std::vector<Label> labels;
		std::vector<std::wstring> text_table;
		std::vector<double> double_table;

		int break_index;
		int continue_index;
		int ast_return;

		// error handle
		int error_count;
		int warning_count;

		position_cache positions;

		std::string current_function_name;
		int current_function_type;
		int system_function_num;

		std::vector<std::filesystem::path> filepathes_;

	public:

		compiler():
			break_index(-1), continue_index(-1), 
			error_count(0), warning_count(0),
			ast_return(0), current_function_type(TYPE_INTEGER),
			positions(std::string("").begin(), std::string("").end()),
			system_function_num(0){
		}

		void SetAstReturn(int value){ast_return = value;}
		int GetAstReturn()const{return ast_return;}

		bool compile(const fs::path& in_filepath, const fs::path& out_filepath);

		// 命令の発行
		#define VM_CREATE
		#include "vm_code.h"
		#undef VM_CREATE

#ifdef _DEBUG
		void debug_dump();
#endif

		// 変数宣言
		void DefineValue(int type, const std::vector<ast::declarator>& node);
		// 関数宣言
		void DeclFunction(int attr, int type, const std::string& name, const ast::arg_def_list& args, const ast::function_pre_def& ast);
		// 関数定義
		void DefineFunction(int type, const std::string& name, const ast::arg_def_list& args, const ast::statements& block, const ast::function_def& ast);

		// 変数定義
		void AddValue(int type, const std::string& name, const ast::declarator& ast);
		// 変数検索
		const ValueTag* GetValueTag(const std::string& name) const
		{
			int size = (int)variables.size();
			for (int i = size - 1; i >= 0; i--)
			{
				const ValueTag* tag = variables[i].find(name);
				if (tag) return tag;
			}
			return nullptr;
		}

		// 関数の検索
		const FunctionTag* GetFunctionTag(const std::string& name) const
		{
			return functions.find(name);
		}

		// ステートメント処理
		void BlockIn();
		void BlockOut();

		// Break分のジャンプ先設定
		int SetBreakLabel(int label)
		{
			int old_index = break_index;
			break_index = label;
			return old_index;
		}
		// Continue文のジャンプ先設定
		int SetContinueLabel(int label)
		{
			int old_index = continue_index;
			continue_index = label;
			return old_index;
		}

		// jmp命令
		bool JmpBreakLabel();
		bool JmpContinueLabel();

		// label処理
		int LabelSetting();
		int MakeLabel();
		void SetLabel(int label);

		void PushString(const std::wstring& name);
		void PushDouble(double value);
		int GetFunctionType() const { return current_function_type; }

		// include命令
		void Include(const std::string& filepath, const x3::position_tagged& ast);

		// 実行データを生成
		bool CreateData(const fs::path& path);

		// error handling
		void error(const std::string& m, const x3::position_tagged& ast);
		void error(const std::string& m);
		void warning(const std::string& m, const x3::position_tagged& ast);
		void warning(const std::string& m);
		void info(const std::string& m, const x3::position_tagged& ast);
		void info(const std::string& m);
	private:
		
		// utf8-bomのbom部分のスキップ
		inline void skip_utf8_bom(std::ifstream& fs)
		{
			int dst[3];
			for (auto& i : dst) i = fs.get();
			constexpr int utf8[] = { 0xEF, 0xBB, 0xBF };

			// bomがついていない場合は読み取り位置を先頭に戻す
			if (!std::equal(std::begin(dst), std::end(dst), utf8))
			{
				info("ファイルフォーマット：utf8");
				fs.seekg(0);
			}
			else
			{
				info("ファイルフォーマット：utf8 with bom");
			}
		}

		// スクリプトファイルの読み出し
		std::string FileLoad(const fs::path& p)
		{
			//ファイルの存在チェック
			if(!fs::exists(p))
			{
				throw CompilerErrorException(p.string() + " : ファイルが見つかりませんでした。パスが正しいか確認してください。");
			}

			//ファイルを開く
			std::ifstream file(p);
			if (!file)
			{
				throw CompilerErrorException(p.string() + " : ファイルが開けませんでした。権限を確認してください。");
			}

			// utf8 bomの処理
			file.imbue(std::locale());
			skip_utf8_bom(file);

			//stringに読み出す
			std::istreambuf_iterator<char> fbegin(file);
			std::istreambuf_iterator<char> fend;
			std::string contents(fbegin, fend);

			return contents;
		}

		// コンパイルのパース部分
		bool parse(const fs::path& path, const x3::position_tagged& called_pos = x3::position_tagged());

	};
}