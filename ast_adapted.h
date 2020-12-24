#pragma once

#include "ast.h"

#include <boost/fusion/include/adapt_struct.hpp>
#include <boost/fusion/include/std_pair.hpp>

// We need to tell fusion about our rexpr and rexpr_key_value
// to make them a first-class fusion citizens

BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::unit,
    entries
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::signed_,
    sign, operand_
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::declaration,
    decl, expression
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::declarator,
    type, identifier_
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::operation,
    sign, operand_
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::expr,
    first, operations
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::assign,
    left, sign, right
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::identifier,
    name
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::function_pre_def,
    return_type, name, args
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::function_def,
    return_type, name, args, states
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::function_call,
    name, args
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::section_statement,
    expression, if_state, else_state
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::for_statement,
    decl, condition, iter, state
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::while_statement,
    condition, state
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::jump_statement,
    ope, expression
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::novel_msg_statement,
    msg, new_page
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::novel_name_statement,
    name
)
BOOST_FUSION_ADAPT_STRUCT(kscript2::ast::import_script,
    file_name
)
