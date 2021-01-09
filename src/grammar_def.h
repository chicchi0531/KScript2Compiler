#pragma once

#include "common.h"

#include "ast.h"
#include "ast_adapted.h"
#include "error_handler.h"
#include "grammar.h"
#include "symbols.h"

#include <iostream>


namespace kscript2 {
    namespace parser
    {
        namespace x3 = boost::spirit::x3;
        namespace unicode = boost::spirit::x3::unicode;
        namespace standard_wide = boost::spirit::x3::standard_wide;

        using x3::lexeme;
        using x3::no_skip;
        using x3::skip;
        using x3::raw;
        using x3::eol;
        using x3::eoi;
        using x3::eps;
        using x3::matches;

        using standard_wide::char_;
        using x3::lit;
        using x3::string;
        using x3::int_;
        using x3::double_;
        using standard_wide::alpha;
        using standard_wide::alnum;
        using standard_wide::blank;

        ///////////////////////////////////////////////////////////////////////////
        // Rules
        ///////////////////////////////////////////////////////////////////////////

        // external definition
        unit_type const unit = "unit";
        x3::rule<struct external_decl_class, ast::external_decl> const external_decl = "external_decl";
        x3::rule<struct function_pre_definition_class, ast::function_pre_def> const function_pre_definition = "function_pre_definition";
        x3::rule<struct function_definition_class, ast::function_def> const function_definition = "function_definition";
        x3::rule<struct function_call_class, ast::function_call> const function_call = "function_call";

        // declaration
        x3::rule<struct declaration_class, ast::declaration> const declaration = "declaration";
        x3::rule<struct declarator_class, ast::declarator> const declarator = "declarator";
        x3::rule<struct arg_def_list_class, ast::arg_def_list> const arg_def_list = "arg_def_list";
        x3::rule<struct arg_list_class, ast::arg_list> const arg_list = "arg_list";

        // statement
        x3::rule<struct statement_class, ast::statement> const statement = "statement";
        x3::rule<struct compound_statement_class, ast::statements> const compound_statement = "compound_statement";
        x3::rule<struct assign_statement_class, ast::assign> const assign_statement = "assign_statement";
        x3::rule<struct section_statement_class, ast::section_statement> const section_statement = "section_statement";
        x3::rule<struct iteration_statement_class, ast::iteration_statement> const iteration_statement = "iteration_statement";
        x3::rule<struct jump_statement_class, ast::jump_statement> const jump_statement = "jump_statement";
        x3::rule<struct function_statement_class, ast::function_call> const function_statement = "function_statement";

        // iteration statement
        x3::rule<struct for_init_class, ast::for_decl> const for_init = "for_init";
        x3::rule<struct for_condition_class, ast::expr> const for_condition = "for_condition";
        x3::rule<struct for_iteration_class, ast::assign> const for_iteration = "for_iteration";
        x3::rule<struct for_statement_class, ast::for_statement> const for_statement = "for_statement";
        x3::rule<struct while_statement_class, ast::while_statement> const while_statement = "while_statement";

        // novel statement
        x3::rule<struct novel_block_class, ast::novel_block> const novel_block = "novel_block";
        x3::rule<struct novel_statement_class, ast::novel_statement> const novel_statement = "novel_statement";
        x3::rule<struct novel_name_statement_class, ast::novel_name_statement> const novel_name_statement = "novel_name_statement";
        x3::rule<struct novel_msg_statement_class, ast::novel_msg_statement> const novel_msg_statement = "novel_msg_statement";
        x3::rule<struct novel_identifier_class, ast::identifier> const novel_identifier = "novel_identifier";
        x3::rule<struct novel_string_class, std::wstring> const novel_string = "novel_string";

        // expression
        x3::rule<struct primary_expr_class, ast::operand> const primary_expr = "primary_expr";
        x3::rule<struct unary_expr_class, ast::operand> const unary_expr = "unary_expr_class";
        x3::rule<struct mul_expr_class, ast::expr> const mul_expr = "mul_expr";
        x3::rule<struct add_expr_class, ast::expr> const add_expr = "add_expr";
        x3::rule<struct relation_expr_class, ast::expr> const relation_expr = "relation_expr";
        x3::rule<struct equality_expr_class, ast::expr> const equality_expr = "equality_expr";
        x3::rule<struct and_expr_class, ast::expr> const and_expr = "and_expr";
        x3::rule<struct or_expr_class, ast::expr> const or_expr = "or_expr";
        x3::rule<struct assign_expr_class, ast::assign> const assign_expr = "assign";
        x3::rule<struct expression_class, ast::expr> const expression = "expression";

        // identifier
        x3::rule<struct identifier_class, ast::identifier> const identifier = "identifier";

        // constants
        x3::rule<struct const_class, ast::constant> const constant = "constant";
        x3::rule<struct int_const_class, int> const int_const = "int_value";
        x3::rule<struct float_const_class, double> const float_const = "float_value";
        x3::rule<struct string_literal_class, std::wstring> const string_literal = "string_value";

        // other
        x3::rule<struct import_class, ast::import_script> const import_script = "import";


        ///////////////////////////////////////////////////////////////////////////
        // Grammar
        ///////////////////////////////////////////////////////////////////////////

        // external definition
        auto const unit_def = *external_decl;
        auto const external_decl_def = import_script | function_pre_definition | function_definition | declaration;
        auto const function_pre_definition_def = -(func_attribute_symbols) >> func_type_symbols >> identifier >> '(' >> -arg_def_list >> ')' >> ';';
        auto const function_definition_def = func_type_symbols >> identifier >> '(' >> -arg_def_list >> ')' >> compound_statement;
        auto const function_call_def = identifier >> '(' >> -arg_list >> ')';

        // declaration
        auto const declaration_def = declarator >> -('=' >> expression) >> ';';
        auto const declarator_def = type_symbols > identifier;
        auto const arg_def_list_def = declarator % ',';
        auto const arg_list_def = expression % ',';

        // statement
        auto const statement_def =
            declaration
            | function_statement
            | section_statement
            | iteration_statement
            | jump_statement
            | assign_statement
            | compound_statement
            | novel_block
            ;
        auto const compound_statement_def = '{' > *statement > '}';
        auto const assign_statement_def = assign_expr >> ';';
        auto const function_statement_def = function_call >> ';';
        auto const section_statement_def =
            lit("if") > '(' > for_condition > ')' > statement > -("else" > statement)
            ;
        auto const iteration_statement_def = for_statement | while_statement;
        auto const jump_statement_def = op_jump_symbols > -expression > ';'
            ;

        // novel statement
        auto const novel_block_def = lit("@{") > skip(blank)[novel_statement % eol] > "@}";
        auto const novel_statement_def = novel_name_statement | novel_msg_statement;
        auto const novel_name_statement_def = 
            lit("-") > (novel_identifier | novel_string)
            ;
        auto const novel_msg_statement_def = 
            (novel_string % L'%') >> -(novel_new_page | novel_new_line)
            ;
        auto const novel_string_def =
            lexeme[
            *(char_ - (char_(L"@%") | eol | novel_new_page | novel_new_line))
            ]
            ;
        auto const novel_identifier_def = lexeme[
            '%' > *(char_ - ('%' | eol)) > '%'
        ];

        // iteration statement
        auto const for_init_def = declaration | (assign_expr >> ';');
        auto const for_condition_def = expression;
        auto const for_iteration_def = assign_expr;
        auto const for_statement_def = 
            lit("for") >> '(' >> -for_init >> -for_condition >> ';' >> -for_iteration >> ')' >> statement
            ;
        auto const while_statement_def =
            lit("while") >> '(' >> for_condition >> ')' >> statement
            ;

        // expression
        auto const primary_expr_def =
            constant									// 定数
            | function_call                             // 関数呼び出し
            | identifier								// 変数
            | '(' > or_expr > ')'					    // カッコの中身
            ;
        auto const unary_expr_def = primary_expr
            | op_sign_symbols > primary_expr
            ;
        auto const mul_expr_def = unary_expr
            >> *(op_mul_symbols > unary_expr)
            ;
        auto const add_expr_def = mul_expr
            >> *(op_add_symbols > mul_expr)
            ;
        auto const relation_expr_def = add_expr
            >> *(op_relation_symbols > add_expr)
            ;
        auto const equality_expr_def = relation_expr
            >> *(op_equ_symbols > relation_expr)
            ;
        auto const and_expr_def = equality_expr
            >> *(op_logand_symbols > equality_expr)
            ;
        auto const or_expr_def = and_expr
            >> *(op_logor_symbols > and_expr)
            ;
        auto const assign_expr_def =
            identifier >> op_assign_symbols >> or_expr
            ;
        auto const expression_def =
            or_expr;

        // identifier
        auto const identifier_def = raw[
            lexeme[
                ((alpha | '_') >> *(alnum | '_'))
                - (keywords >> char_ - (alnum | '_')) //予約語を除く
            ]];

        // constant
        x3::real_parser<double, x3::strict_real_policies<double>> strict_double; //double_はすべての数値にマッチするので、厳密なrealパーサーを使う
        auto const constant_def = float_const | int_const | string_literal;
        auto const int_const_def = int_;
        auto const float_const_def = strict_double;
        auto const string_literal_def = lexeme[
            L'"' >> *(char_ - (L'"' | eol)) >> L'"'
        ];

        // other
        auto const import_script_def = lit("import") > string_literal > ';';
        
        // definition macro
        BOOST_SPIRIT_DEFINE(
            unit,
            external_decl,
            function_pre_definition,
            function_definition,
            function_call,
            arg_def_list,
            arg_list,
            declaration,
            declarator,
            statement,
            compound_statement,
            function_statement,
            assign_statement,
            section_statement,
            jump_statement,
            iteration_statement,
            for_init,
            for_condition,
            for_iteration,
            for_statement,
            while_statement,
            novel_block,
            novel_statement,
            novel_name_statement,
            novel_msg_statement,
            novel_string,
            novel_identifier,
            primary_expr,
            unary_expr,
            mul_expr,
            add_expr,
            relation_expr,
            equality_expr,
            and_expr,
            or_expr,
            assign_expr,
            expression,
            identifier,
            constant,
            int_const,
            float_const,
            string_literal,
            import_script
        )

        ///////////////////////////////////////////////////////////////////////////
        // Annotation and Error handling
        ///////////////////////////////////////////////////////////////////////////
        #define GRAMMAR_PRINTER(name) name(){std::cout<< #name << std::endl; }
        #define GRAMMAR_ID_CLASS_DEF(name) struct name : annotation_base { /*GRAMMAR_PRINTER(name)*/ };

        // We want error-handling only for the start (outermost) rexpr
        // rexpr is the same as rexpr_inner but without error-handling (see error_handler.hpp)
        struct unit_class : annotation_base ,error_handler_base { /*GRAMMAR_PRINTER(unit_class)*/ };

        // We want these to be annotated with the iterator position.
        GRAMMAR_ID_CLASS_DEF(external_decl_class)
        GRAMMAR_ID_CLASS_DEF(function_pre_definition_class)
        GRAMMAR_ID_CLASS_DEF(function_definition_class)

        GRAMMAR_ID_CLASS_DEF(declaration_class)
        GRAMMAR_ID_CLASS_DEF(declarator_class)
        GRAMMAR_ID_CLASS_DEF(arg_def_class)
        GRAMMAR_ID_CLASS_DEF(arg_def_list_class)


        GRAMMAR_ID_CLASS_DEF(statement_class)
        GRAMMAR_ID_CLASS_DEF(compound_statement_class)
        GRAMMAR_ID_CLASS_DEF(section_statement_class)
        GRAMMAR_ID_CLASS_DEF(iteration_statement_class)
        GRAMMAR_ID_CLASS_DEF(jump_statement_class)
        GRAMMAR_ID_CLASS_DEF(assign_statement_class)

        GRAMMAR_ID_CLASS_DEF(for_statement_class)
        GRAMMAR_ID_CLASS_DEF(while_statement_class)

        GRAMMAR_ID_CLASS_DEF(novel_block_class)
        GRAMMAR_ID_CLASS_DEF(novel_statement_class)
        GRAMMAR_ID_CLASS_DEF(novel_name_statement_class)
        GRAMMAR_ID_CLASS_DEF(novel_msg_statement_class)
        GRAMMAR_ID_CLASS_DEF(novel_identifier_class)
        GRAMMAR_ID_CLASS_DEF(novel_string_class)

        GRAMMAR_ID_CLASS_DEF(primary_expr_class)
        GRAMMAR_ID_CLASS_DEF(unary_expr_class)
        GRAMMAR_ID_CLASS_DEF(mul_expr_class)
        GRAMMAR_ID_CLASS_DEF(add_expr_class)
        GRAMMAR_ID_CLASS_DEF(relation_expr_class)
        GRAMMAR_ID_CLASS_DEF(equality_expr_class)
        GRAMMAR_ID_CLASS_DEF(and_expr_class)
        GRAMMAR_ID_CLASS_DEF(or_expr_class)
        GRAMMAR_ID_CLASS_DEF(assign_expr_class)
        GRAMMAR_ID_CLASS_DEF(expression_class)

        GRAMMAR_ID_CLASS_DEF(identifier_class)

        GRAMMAR_ID_CLASS_DEF(constant_class)
        GRAMMAR_ID_CLASS_DEF(int_const_class)
        GRAMMAR_ID_CLASS_DEF(float_const_class)
        GRAMMAR_ID_CLASS_DEF(string_literal_class)

        GRAMMAR_ID_CLASS_DEF(import_class)

    }
}

namespace kscript2
{
    parser::unit_type const& unit()
    {
        return parser::unit;
    }
}
