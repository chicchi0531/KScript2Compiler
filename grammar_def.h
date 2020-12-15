#pragma once

#include "common.h"

#include "ast.h"
#include "ast_adapted.h"
#include "error_handler.h"
#include "grammar.h"
#include "symbols.h"


namespace kscript2 {
    namespace parser
    {
        namespace x3 = boost::spirit::x3;
        namespace unicode = boost::spirit::x3::unicode;
        namespace standard_wide = boost::spirit::x3::standard_wide;

        using x3::lit;
        using x3::lexeme;
        using x3::raw;
        using x3::eol;
        using x3::eoi;

        using standard_wide::char_;
        using x3::string;
        using x3::int_;
        using x3::double_;
        using standard_wide::alpha;
        using standard_wide::alnum;

        ///////////////////////////////////////////////////////////////////////////
        // Rule IDs
        ///////////////////////////////////////////////////////////////////////////

        // external definition
        struct external_decl_class;
        struct function_pre_definition_class;
        struct function_definition_class;
        struct function_call_class;

        // declaration
        struct declaration_class;
        struct declarator_class;
        struct arg_def_list_class;
        struct arg_list_class;

        struct identifier_class;

        // statement
        struct statement_class;
        struct compound_statement_class;
        struct assign_statement_class;
        struct section_statement_class;
        struct function_statement_class;

        // expression
        struct primary_expr_class;
        struct unary_expr_class;
        struct mul_expr_class;
        struct add_expr_class;
        struct expression_class;
        struct relation_expr_class;
        struct equality_expr_class;
        struct and_expr_class;
        struct or_expr_class;

        // constants
        struct const_class;
        struct int_const_class;
        struct float_const_class;
        struct string_literal_class;

        ///////////////////////////////////////////////////////////////////////////
        // Rules
        ///////////////////////////////////////////////////////////////////////////

        // external definition
        unit_type const unit = "unit";
        x3::rule<external_decl_class, ast::external_decl> const external_decl = "external_decl";
        x3::rule<function_pre_definition_class, ast::function_pre_def> const function_pre_definition = "function_pre_definition";
        x3::rule<function_definition_class, ast::function_def> const function_definition = "function_definition";
        x3::rule<function_call_class, ast::function_call> const function_call = "function_call";

        // declaration
        x3::rule<declaration_class, ast::declaration> const declaration = "declaration";
        x3::rule<declarator_class, ast::declarator> const declarator = "declarator";
        x3::rule<arg_def_list_class, ast::arg_def_list> const arg_def_list = "arg_def_list";
        x3::rule<arg_list_class, ast::arg_list> const arg_list = "arg_list";

        // statement
        x3::rule<statement_class, ast::statement> const statement = "statement";
        x3::rule<compound_statement_class, ast::statements> const compound_statement = "compound_statement";
        x3::rule<assign_statement_class, ast::assign_statement> const assign_statement = "assign_statement";
        x3::rule<section_statement_class, ast::section_statement> const section_statement = "section_statement";
        x3::rule<function_statement_class, ast::function_call> const function_statement = "function_statement";

        // expression
        x3::rule<primary_expr_class, ast::operand> const primary_expr = "primary_expr";
        x3::rule<unary_expr_class, ast::operand> const unary_expr = "unary_expr_class";
        x3::rule<mul_expr_class, ast::expr> const mul_expr = "mul_expr";
        x3::rule<add_expr_class, ast::expr> const add_expr = "add_expr";
        x3::rule<relation_expr_class, ast::expr> const relation_expr = "relation_expr";
        x3::rule<equality_expr_class, ast::expr> const equality_expr = "equality_expr";
        x3::rule<and_expr_class, ast::expr> const and_expr = "and_expr";
        x3::rule<or_expr_class, ast::expr> const or_expr = "or_expr";
        x3::rule<expression_class, ast::expr> const expression = "expression";

        // identifier
        x3::rule<identifier_class, ast::identifier> const identifier = "identifier";

        // constants
        x3::rule<const_class, ast::constant> const constant = "constant";
        x3::rule<int_const_class, int> const int_const = "int_value";
        x3::rule<float_const_class, double> const float_const = "float_value";
        x3::rule<string_literal_class, std::wstring> const string_literal = "string_value";


        ///////////////////////////////////////////////////////////////////////////
        // Grammar
        ///////////////////////////////////////////////////////////////////////////

        // external definition
        auto const unit_def = *external_decl;
        auto const external_decl_def =  function_pre_definition | function_definition | declaration;
        auto const function_pre_definition_def = func_type_symbols >> identifier >> '(' >> -arg_def_list >> ')' >> ';';
        auto const function_definition_def = func_type_symbols >> identifier >> '(' >> -arg_def_list >> ')' >> compound_statement;
        auto const function_call_def = identifier >> '(' >> -arg_list >> ')';

        // declaration
        auto const declaration_def = declarator >> -('=' >> expression) >> ';';
        auto const declarator_def = type_symbols > identifier;
        auto const arg_def_list_def = declarator % ',';
        auto const arg_list_def = expression % ',';

        // statement
        auto const statement_def =
            function_statement
            | section_statement
            | declaration
            | assign_statement
            | compound_statement
            ;
        auto const compound_statement_def = '{' > *statement > '}';
        auto const assign_statement_def = identifier >> '=' >> expression >> ';';
        auto const function_statement_def = function_call >> ';';
        auto const section_statement_def =
            lit("if") > '(' > expression > ')' > statement > -("else" > statement)
            ;

        // expression
        auto const primary_expr_def =
            constant									// 定数
            | function_call                             // 関数呼び出し
            | identifier								// 変数
            | '(' > expression > ')'					// カッコの中身
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
        auto const expression_def = or_expr;

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
            primary_expr,
            unary_expr,
            mul_expr,
            add_expr,
            relation_expr,
            equality_expr,
            and_expr,
            or_expr,
            expression,
            identifier,
            constant,
            int_const,
            float_const,
            string_literal);

        ///////////////////////////////////////////////////////////////////////////
        // Annotation and Error handling
        ///////////////////////////////////////////////////////////////////////////

        // We want error-handling only for the start (outermost) rexpr
        // rexpr is the same as rexpr_inner but without error-handling (see error_handler.hpp)
        struct unit_class : x3::annotate_on_success, error_handler_base {};

        // We want these to be annotated with the iterator position.
        struct external_decl_class : x3::annotate_on_success {};
        struct function_pre_definition_class : x3::annotate_on_success {};
        struct function_definition_class : x3::annotate_on_success {};

        struct declaration_class : x3::annotate_on_success {};
        struct declarator_class : x3::annotate_on_success {};
        struct arg_def_class : x3::annotate_on_success {};
        struct arg_def_list_class : x3::annotate_on_success {};


        struct statement_class : x3::annotate_on_success {};
        struct copound_statement_class : x3::annotate_on_success {};
        struct section_statement_class : x3::annotate_on_success {};
        struct assign_statement_class : x3::annotate_on_success {};
        struct rexpr_value_class : x3::annotate_on_success {};

        struct primary_expr_class : x3::annotate_on_success {};
        struct unary_expr_class : x3::annotate_on_success {};
        struct mul_expr_class : x3::annotate_on_success {};
        struct add_expr_class : x3::annotate_on_success {};
        struct relation_expr_class : x3::annotate_on_success {};
        struct equality_expr_class : x3::annotate_on_success {};
        struct and_expr_class : x3::annotate_on_success {};
        struct or_expr_class : x3::annotate_on_success {};
        struct expression_class : x3::annotate_on_success {};

        struct identifier_class : x3::annotate_on_success {};

        struct constant_class : x3::annotate_on_success {};
        struct int_const_class : x3::annotate_on_success {};
        struct float_const_class : x3::annotate_on_success {};
        struct string_literal_class : x3::annotate_on_success {};

    }
}

namespace kscript2
{
    parser::unit_type const& unit()
    {
        return parser::unit;
    }
}
