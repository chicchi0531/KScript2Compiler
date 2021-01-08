
#pragma once

#include "common.h"

#include <map>
#include <vector>


namespace kscript2 {
    namespace ast
    {
        ///////////////////////////////////////////////////////////////////////////
        //  The AST
        ///////////////////////////////////////////////////////////////////////////
        namespace x3 = boost::spirit::x3;

        struct nil {};
        struct unit;
        struct function_call;
        struct statement;
        struct section_statement;
        struct iteration_statement;
        struct jump_statement;
        struct declaration;
        struct operation;
        struct expr;
        struct signed_;
        struct identifier;

        struct location_info : x3::position_tagged
        {
        };

        // other
        struct import_script : location_info
        {
            std::string file_name;
        };

        //constant
        struct constant :
            x3::variant<
            int,
            double,
            std::wstring
            >
        {
            using base_type::base_type;
            using base_type::operator=;
        };

        struct identifier : location_info
        {
            std::string name;
        };


        //expression
        struct operand : x3::variant<
            nil,
            x3::forward_ast<identifier>,
            x3::forward_ast<constant>,
            x3::forward_ast<signed_>,
            x3::forward_ast<expr>,
            x3::forward_ast<function_call>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };

        struct signed_ : location_info
        {
            int sign;
            operand operand_;
        };

        struct operation : location_info
        {
            int sign;
            operand operand_;
        };

        typedef std::vector<operation> operation_list;
        struct expr : location_info
        {
            operand first;
            operation_list operations;
        };

        struct assign : location_info
        {
            identifier left;
            int sign;
            expr right;
        };
        typedef std::vector<assign> assign_list;

        // decl
        struct declarator : location_info
        {
            int type;
            identifier identifier_;
        };
        struct declaration : location_info
        {
            declarator decl;
            expr expression;
        };

        // novel_statement
        struct novel_msg : x3::variant<
            nil,
            std::wstring,
            x3::forward_ast<identifier>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };
        struct novel_msg_statement : location_info
        {
            std::vector<std::wstring> msg;
            int new_page;
        };
        struct novel_name_statement : location_info
        {
            novel_msg name;
        };
        struct novel_statement : x3::variant<
            nil,
            x3::forward_ast<novel_name_statement>,
            x3::forward_ast<novel_msg_statement>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };
        typedef std::vector<novel_statement> novel_block;

        // statement
        typedef std::vector<statement> statements;
        struct statement : x3::variant<
            nil,
            x3::forward_ast<declaration>,
            x3::forward_ast<assign>,
            x3::forward_ast<function_call>,
            x3::forward_ast<statements>,
            x3::forward_ast<section_statement>,
            x3::forward_ast<iteration_statement>,
            x3::forward_ast<jump_statement>,
            x3::forward_ast<novel_block>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };
        struct section_statement : location_info
        {
            expr expression;
            statement if_state;
            statement else_state;
        };
        struct for_decl : x3::variant<
            nil,
            x3::forward_ast<declaration>,
            x3::forward_ast<assign>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };
        struct for_statement : location_info
        {
            for_decl decl;
            expr condition;
            assign iter;
            statement state;
        };
        struct while_statement : location_info
        {
            expr condition;
            statement state;
        };
        struct iteration_statement : x3::variant<
            x3::forward_ast<for_statement>,
            x3::forward_ast<while_statement>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };
        struct jump_statement : location_info
        {
            int ope;
            expr expression;
        };

        // function
        typedef std::vector<declarator> arg_def_list;
        typedef std::vector<expr> arg_list;
        struct function_pre_def : location_info
        {
            int attr_type;
            int return_type;
            identifier name;
            arg_def_list args;
        };
        struct function_def : location_info
        {
            int return_type;
            identifier name;
            arg_def_list args;
            statements states;
        };
        struct function_call : location_info
        {
            identifier name;
            arg_list args;
        };

        // unit
        struct external_decl : x3::variant<
            x3::forward_ast<declaration>,
            x3::forward_ast<function_def>,
            x3::forward_ast<function_pre_def>,
            x3::forward_ast<import_script>
        >
        {
            using base_type::base_type;
            using base_type::operator=;
        };

        typedef std::vector<external_decl> external_decls;
        struct unit : location_info
        {
            external_decls entries;
        };



    }
}

