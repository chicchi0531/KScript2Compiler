#pragma once

#include "ast.h"

#include <ostream>

namespace kscript2 {
    namespace ast
    {
        ///////////////////////////////////////////////////////////////////////////
        //  Print out the rexpr tree
        ///////////////////////////////////////////////////////////////////////////
        int const tabsize = 4;

        struct rexpr_printer
        {
            typedef void result_type;

            rexpr_printer(std::ostream& out, int indent = 0)
                : out(out), indent(indent) {}

            void operator()(unit const& ast) const
            {
                for (auto const& entry : ast.entries)
                {
                    boost::apply_visitor(rexpr_printer(out, indent), entry);
                }
            }

            void operator()(external_decl const& ast) const
            {
                boost::apply_visitor(rexpr_printer(out, indent + tabsize), ast);
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

                rexpr_printer(out, indent)(ast.states);
            }

            void operator()(function_call const& ast) const
            {
                for (auto const& arg : ast.args)
                {
                    rexpr_printer(out, indent)(arg);
                }
                tab(indent);
                out << "func_call " << ast.name.name << std::endl;
            }

            void operator()(statement const& ast) const
            {
                boost::apply_visitor(rexpr_printer(out, indent), ast);
            }

            void operator()(nil const& ast) const{}

            void operator()(statements const& ast)const
            {
                tab(indent);
                out << "{" << std::endl;
                for (auto const& entry : ast)
                {
                    boost::apply_visitor(rexpr_printer(out, indent + tabsize), entry);
                }                
                tab(indent);
                out << "}" << std::endl;
            }

            void operator()(assign_statement const& ast) const
            {
                rexpr_printer(out, indent)(ast.first);
                rexpr_printer(out, indent)(ast.second);

                tab(indent);
                out << "assign" << std::endl;
            }

            void operator()(section_statement const& ast) const
            {
                rexpr_printer(out, indent)(ast.expression);

                tab(indent);
                out << "jump_zero: endif" << std::endl;
                rexpr_printer(out, indent)(ast.if_state);

                //elseステートメント
                if (ast.else_state.get().which() != 0)
                {
                    tab(indent);
                    out << "jmp: end_else" << std::endl;

                    tab(indent);
                    out << "label: endif" << std::endl;

                    boost::apply_visitor(rexpr_printer(out, indent), ast.else_state);

                    tab(indent);
                    out << "label: end_else" << std::endl;
                }
                else
                {
                    tab(indent);
                    out << "label: endif" << std::endl;
                }
            }

            void operator()(declarator const& ast) const
            {
                tab(indent);
                out << "decl_var" << " type:" << ast.type << " name:" << ast.identifier_.name << std::endl;
            }

            void operator()(declaration const& ast) const
            {
                rexpr_printer(out, indent)(ast.decl);
                
                // 初期化子つき宣言
                if (ast.expression.first.get().which() != 0)
                {
                    assign_statement assign;
                    assign.first = ast.decl.identifier_;
                    assign.second = ast.expression;
                    rexpr_printer(out, indent)(assign);
                }
            }


            void operator()(operand const& ast) const
            {
                boost::apply_visitor(rexpr_printer(out, indent), ast);
            }

            void operator()(signed_ const& ast) const
            {
                boost::apply_visitor(rexpr_printer(out, indent), ast.operand_);

                tab(indent);
                out << "ope_" << ast.sign << std::endl;
            }

            void operator()(operation const& ast) const
            {
                boost::apply_visitor(rexpr_printer(out, indent), ast.operand_);

                tab(indent);
                out << "ope_" << ast.sign << std::endl;
            }

            void operator()(expr const& ast) const
            {
                // nilの場合は棄却
                if (ast.first.get().which() == 0)return;

                boost::apply_visitor(rexpr_printer(out, indent), ast.first);
                for (auto const& ope : ast.operations)
                {
                    rexpr_printer(out, indent)(ope);
                }
            }

            void operator()(identifier const& ast) const
            {
                tab(indent);
                out << "push_var " << ast.name << std::endl;
            }

            void operator()(constant const& ast) const
            {
                boost::apply_visitor(rexpr_printer(out, indent), ast);
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

