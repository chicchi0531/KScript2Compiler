#pragma once

#include "symbols.h"
#include "ast.h"

#include <ostream>

namespace kscript2 {
    namespace ast
    {
        ///////////////////////////////////////////////////////////////////////////
        //  Print out the rexpr tree
        ///////////////////////////////////////////////////////////////////////////
        int const tabsize = 4;

        struct ast_analyzer
        {
            typedef void result_type;

            ast_analyzer(std::ostream& out, int indent = 0)
                : out(out), indent(indent) {}

            void operator()(unit const& ast) const
            {
                for (auto const& entry : ast.entries)
                {
                    boost::apply_visitor(ast_analyzer(out, indent), entry);
                }
            }

            void operator()(function_pre_def const& ast) const
            {
                out << "func_pre_def " << ast.return_type << " " << ast.name.name << " ";

                // 引数出力
                for (auto const& arg : ast.args)
                {
                    out << "arg:" << arg.type << ":" << arg.identifier_.name << " ";
                }
                out << std::endl;

            }

            void operator()(function_def const& ast) const
            {
                out << "func_def" << ast.return_type << " " << ast.name.name << " ";

                // 引数出力
                for (auto const& arg : ast.args)
                {
                    out << "arg:" << arg.type << ":" << arg.identifier_.name << " ";
                }
                out << std::endl;

                ast_analyzer(out, indent)(ast.states);
            }

            void operator()(function_call const& ast) const
            {
                for (auto const& arg : ast.args)
                {
                    ast_analyzer(out, indent)(arg);
                }
                tab(indent);
                out << "func_call " << ast.name.name << std::endl;
            }

            void operator()(nil const& ast) const{}

            void operator()(statements const& ast)const
            {
                tab(indent);
                out << "{" << std::endl;
                for (auto const& entry : ast)
                {
                    boost::apply_visitor(ast_analyzer(out, indent + tabsize), entry);
                }                
                tab(indent);
                out << "}" << std::endl;
            }

            void operator()(section_statement const& ast) const
            {
                ast_analyzer(out, indent)(ast.expression);

                tab(indent);
                out << "jump_zero: endif" << std::endl;
                boost::apply_visitor(ast_analyzer(out, indent),ast.if_state);

                //elseステートメント
                if (ast.else_state.get().which() != 0)
                {
                    tab(indent);
                    out << "jmp: end_else" << std::endl;

                    tab(indent);
                    out << "label: endif" << std::endl;

                    boost::apply_visitor(ast_analyzer(out, indent), ast.else_state);

                    tab(indent);
                    out << "label: end_else" << std::endl;
                }
                else
                {
                    tab(indent);
                    out << "label: endif" << std::endl;
                }
            }

            void operator()(iteration_statement const& ast) const
            {
                boost::apply_visitor(ast_analyzer(out, indent), ast);
            }

            void operator()(jump_statement const& ast) const
            {
                switch (ast.ope)
                {
                case parser::OP_CONTINUE:
                    if (ast.expression.first.get().which() != 0)
                    {
                        std::cerr << "意味解析エラー：continue文に不要な式が渡されました";
                        return;
                    }
                    tab(indent);
                    out << "jmp:直近のループ開始地点" << std::endl;
                    break;

                case parser::OP_BREAK:
                    if (ast.expression.first.get().which() != 0)
                    {
                        std::cerr << "意味解析エラー：break文に不要な式が渡されました";
                        return;
                    }
                    tab(indent);
                    out << "break:直後のループ終了地点" << std::endl;
                    break;

                case parser::OP_RETURN:
                    ast_analyzer(out, indent)(ast.expression);
                    tab(indent);
                    out << "jump_ret" << std::endl;
                    break;
                }
            }

            void operator()(for_statement const& ast) const
            {
                //初期化部
                boost::apply_visitor(ast_analyzer(out, indent), ast.decl);

                tab(indent);
                out << "label: for_begin" << std::endl;

                //条件部
                ast_analyzer(out, indent)(ast.condition);

                tab(indent);
                out << "jump_zero: for_end" << std::endl;

                //処理
                boost::apply_visitor(ast_analyzer(out, indent), ast.state);

                //イテレーション
                ast_analyzer(out, indent)(ast.iter);

                tab(indent);
                out << "jmp: for_bagin" << std::endl;
                tab(indent);
                out << "lebel: for_end" << std::endl;
            }

            void operator()(while_statement const& ast) const
            {
                tab(indent);
                out << "label: while_begin" << std::endl;

                //条件部
                ast_analyzer(out, indent)(ast.condition);

                tab(indent);
                out << "jump_zero: while_end" << std::endl;

                //処理
                boost::apply_visitor(ast_analyzer(out, indent), ast.state);

                tab(indent);
                out << "jmp: while_begin" << std::endl;
                tab(indent);
                out << "label: while_end" << std::endl;
            }

            //novel
            void operator()(novel_block const& ast) const
            {
                for (auto const& s : ast)
                {
                    boost::apply_visitor(ast_analyzer(out, indent), s);
                }
            }
            void operator()(novel_name_statement const& ast) const
            {
                boost::apply_visitor(ast_analyzer(out, indent), ast.name);

                tab(indent);
                out << "novel_name" << std::endl;
            }
            void operator()(novel_msg_statement const& ast) const
            {
                if (ast.msg.size() == 0 && ast.new_page == 0) return;

                // %が閉じられているかのチェック
                // 偶数の場合は閉じられていない
                if (ast.msg.size() % 2 == 0)
                {
                    std::cerr << ast.id_last << ": %が閉じられていません。" << std::endl;
                }

                bool is_identifier = false;
                for (auto const& m : ast.msg)
                {
                    if (is_identifier)
                    {
                        identifier i;
                        std::string s(m.begin(), m.end());
                        i.name = s;
                        ast_analyzer(out, indent)(i);
                    }
                    else
                    {
                        ast_analyzer(out, indent)(m);
                    }

                    tab(indent);
                    out << "novel_msg" << std::endl;

                    is_identifier = !is_identifier;
                }
                //boost::apply_visitor(ast_analyzer(out, indent), ast.msg);
                //ast_analyzer(out, indent)(ast.msg);

                //改行、改ページ処理
                switch (ast.new_page)
                {
                case parser::NOVEL_EMPTY:
                case parser::NOVEL_NEWLINE:
                    tab(indent);
                    out << "■novel_new_line" << std::endl;
                    break;
                case parser::NOVEL_NEWPAGE:
                    tab(indent);
                    out << "■novel_new_page" << std::endl;
                    break;
                }
            }


            void operator()(declarator const& ast) const
            {
                tab(indent);
                out << "decl_var" << " type:" << ast.type << " name:" << ast.identifier_.name << std::endl;
            }

            void operator()(declaration const& ast) const
            {
                ast_analyzer(out, indent)(ast.decl);
                
                // 初期化子つき宣言
                if (ast.expression.first.get().which() != 0)
                {
                    // 内部でassignを構成して呼び出す
                    assign a;
                    a.left = ast.decl.identifier_;
                    a.sign = parser::OP_ASSIGN;
                    a.right = ast.expression;
                    ast_analyzer(out, indent)(a);
                }
            }

            void operator()(signed_ const& ast) const
            {
                boost::apply_visitor(ast_analyzer(out, indent), ast.operand_);

                tab(indent);
                out << "ope_" << ast.sign << std::endl;
            }

            void operator()(operation const& ast) const
            {
                boost::apply_visitor(ast_analyzer(out, indent), ast.operand_);

                tab(indent);
                out << "ope_" << ast.sign << std::endl;
            }

            void operator()(expr const& ast) const
            {
                // nilの場合は棄却
                if (ast.first.get().which() == 0)return;

                boost::apply_visitor(ast_analyzer(out, indent), ast.first);
                for (auto const& ope : ast.operations)
                {
                    ast_analyzer(out, indent)(ope);
                }
            }

            void operator()(assign const& ast) const
            {
                ast_analyzer(out, indent)(ast.right);

                tab(indent);
                out << "assign var:" << ast.left.name << " ope:" << ast.sign << std::endl;
            }
            void operator()(assign_list const& ast) const
            {
                for (auto const& a : ast)
                {
                    ast_analyzer(out, indent)(a);
                }
            }

            void operator()(constant const& ast) const
            {
                boost::apply_visitor(ast_analyzer(out, indent), ast);
            }

            void operator()(identifier const& ast) const
            {
                tab(indent);
                out << "push_var " << ast.name << std::endl;
            }

            void operator()(std::wstring const& text) const
            {
                tab(indent);
                const std::string str(text.begin(), text.end());
                out << "push_wide_str " << '"' << str << '"' << std::endl;
            }

            void operator()(std::string const& text) const
            {
                tab(indent);
                out << "push_str " << '"' << text << '"' << std::endl;
            }

            void operator()(int value) const
            {
                tab(indent);
                out << "push_int " << value << std::endl;
            }

            void operator()(double value) const
            {
                tab(indent);
                out << "push_double " << value << std::endl;
            }

            void operator()(import_script const& ast) const
            {
                tab(indent);
                out << "import: " << ast.file_name << std::endl;
            }

            void tab(int spaces) const
            {
                for (int i = 0; i < spaces; ++i)
                    out << ' ';
            }

            std::ostream& out;
            int indent;
        };
    }
}

