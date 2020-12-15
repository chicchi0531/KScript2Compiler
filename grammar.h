
#pragma once

#include "common.h"
#include "ast.h"

namespace kscript2
{
    namespace x3 = boost::spirit::x3;


    ///////////////////////////////////////////////////////////////////////////
    // rexpr public interface
    ///////////////////////////////////////////////////////////////////////////
    namespace parser
    {
        // comment
        namespace comment
        {
            using x3::standard_wide::char_;
            using x3::standard_wide::space;
            using x3::eol;
            using x3::lexeme;
            static auto const comment = lexeme[
                 "/*" >> *(char_ - "*/") >> "*/"
               | "//" >> *~char_("\r\n") >> eol
            ];
            static auto const skipper = comment | space;
        }

        // grammar
        struct unit_class;
        typedef x3::rule<unit_class, ast::unit> unit_type;
        BOOST_SPIRIT_DECLARE(unit_type);
    }

    parser::unit_type const& unit();
}

