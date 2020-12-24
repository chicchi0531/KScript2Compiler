#pragma once

#include <string>
#include <vector>
#include <map>

#include "ast.h"
#include "symbols.h"

namespace kscript2
{
	using namespace parser;

	// 命令コード
	class VMCode
	{
	public:
		int op_;
		int arg1_;
		int size_;

	public:
		VMCode(int op)
			: size_(1), op_(op), arg1_(0)
		{}
		VMCode(int op, int arg1)
			: size_(5), op_(op), arg1_(arg1)
		{}

		int* Get(int* p) const
		{
			// ラベル以外を取得
			// ラベルは命令として載せないのでスキップ
			if (op_ != OP_MAXCOMMAND)
			{
				*p++ = op_;
				if (size_ > 1)
				{
					*(int*)p = arg1_;
					p += 4;
				}
			}
			return p;
		}
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
		struct dump_action
		{
			void operator()(const std::pair<std::string, ValueTag>& it)
			{
				std::cout << it.first << ", addr = " << it.second.addr_
					<< ", type = " << it.second.type_
					<< ", size = " << it.second.size_
					<< ", global = " << it.second.global_ << std::endl;
			}
		};

		void dump() const
		{
			std::cout << "--------------- value ---------------" << std::endl;
			for (auto const& i : variables_)
			{
				dump_action(i);
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
		FunctionTag(){}
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
		bool IsDeclaration() const { return (flags_ & flag_declaration) != 0; }
		bool IsDefinition() const { return(flags_ & flag_definition) != 0; }
		bool IsSystem() const { return (flags_ & flag_system) != 0; }

	};

	class FunctionTable
	{
	private:
		typedef std::map<std::string, FunctionTag>::iterator iter;
		typedef std::map<std::string, FunctionTag>::const_iterator const_iter;

	public:
		FunctionTable(){}

		FunctionTag* add(const std::string& name, const FunctionTag& tag)
		{

		}
	};

	class compiler
	{
	private:
		std::vector<VMCode> _program;

	public:
		bool compile(const std::string& filepath);

		// 命令の発行
		void op_add();
		void op_sub();
		void op_mul();
		void op_div();
		void op_mod();
		void op_equ();
		void op_neg();


	private:
		void op_(int code) { _program.push_back(VMCode(code, arg_type())); }
		void op_(int code, arg_type args) { _program.push_back(VMCode(code, args)); }

	};
}