#pragma once

#include "error_handler.h"
#include "common.h"
#include "grammar.h"

namespace kscript2 {
    namespace parser
    {
        // Our Iterator Type
        typedef std::string::const_iterator iterator_type;
        namespace x3 = boost::spirit::x3;

        // The Phrase Parse Context
        //typedef
        //    x3::phrase_parse_context<x3::ascii::space_type>::type
        //    phrase_context_type;
        using skipper_type = decltype(kscript2::parser::comment::skipper);
        using phrase_context_type = x3::phrase_parse_context<skipper_type>::type;

        // Our Error Handler
        typedef error_handler<iterator_type> error_handler_type;

        // Combined Error Handler and Phrase Parse Context
        typedef x3::context<
            error_handler_tag
            , std::reference_wrapper<error_handler_type>
            , phrase_context_type>
            context_type;
        
        // position cache
        using position_cache = x3::position_cache<std::vector<iterator_type>>;

    }
}

