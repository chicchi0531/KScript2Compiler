#pragma once

#include "symbols.h"
#include "ast.h"
#include "compiler.h"

#include <ostream>

namespace kscript2 {
    namespace ast
    {
        ///////////////////////////////////////////////////////////////////////////
        //  Analyze AST
        ///////////////////////////////////////////////////////////////////////////
        int const tabsize = 4;

        struct ast_analyzer
        {
        private:
            compiler& compiler_;

        public:
            typedef void result_type;

            ast_analyzer(compiler& c)
                : compiler_(c) {}

            // external definition
            void operator()(nil const& ast) const {}
            void operator()(unit const& ast) const;

            // functions
            void operator()(function_pre_def const& ast) const;
            void operator()(function_def const& ast) const;
            int operator()(function_call const& ast)const;

            // statements
            void operator()(statements const& ast) const;
            void operator()(section_statement const& ast) const;
            void operator()(iteration_statement const& ast) const;
            void operator()(jump_statement const& ast) const;

            void operator()(for_statement const& ast) const;
            void operator()(while_statement const& ast) const;

            //novel
            void operator()(novel_block const& ast) const;
            void operator()(novel_name_statement const& ast) const;
            void operator()(novel_msg_statement const& ast) const;

            // declaration
            void operator()(declarator const& ast) const;
            void operator()(declaration const& ast) const;

            void operator()(signed_ const& ast) const;
            void operator()(operation const& ast) const;

            int operator()(expr const& ast) const
            {
                // nilÇÃèÍçáÇÕä¸ãp
                if (ast.first.get().which() == 0)return 1;

                boost::apply_visitor(ast_analyzer(compiler_), ast.first);
                for (auto const& ope : ast.operations)
                {
                    ast_analyzer(compiler_)(ope);
                }
                return 0;
            }

            void operator()(assign const& ast) const;
            void operator()(assign_list const& ast) const;

            // constant
            void operator()(constant const& ast) const;
            int operator()(identifier const& ast) const;
            void operator()(std::wstring const& text) const;
            void operator()(int value) const;
            void operator()(double value) const;

            void operator()(import_script const& ast) const
            {
                compiler_.Include(ast.file_name);
            }
        };
    }
}

