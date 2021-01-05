#pragma once

#include "common.h"

#include <map>

namespace kscript2 {
    namespace parser
    {
        namespace x3 = boost::spirit::x3;

        ////////////////////////////////////////////////////////////////////////////
        //  Our error handler
        ////////////////////////////////////////////////////////////////////////////
        // X3 Error Handler Utility
        template <typename Iterator>
        using error_handler = x3::error_handler<Iterator>;

        // tag used to get our error handler from the context
        using error_handler_tag = x3::error_handler_tag;

        struct error_handler_base
        {
            error_handler_base();

            template <typename Iterator, typename Exception, typename Context>
            x3::error_handler_result on_error(
                Iterator& first, Iterator const& last
                , Exception const& x, Context const& context);

            template <typename T, typename Iterator, typename Context>
            inline void on_success(Iterator const& first, Iterator const& last
            , T& ast, Context const& context);

        };

        ////////////////////////////////////////////////////////////////////////////
        // Implementation
        ////////////////////////////////////////////////////////////////////////////

        inline error_handler_base::error_handler_base()
        {
        }

        template <typename Iterator, typename Exception, typename Context>
        inline x3::error_handler_result
            error_handler_base::on_error(
                Iterator& first, Iterator const& last
                , Exception const& x, Context const& context)
        {
            std::string which = x.which();
            std::string message = "Error! Expecting: " + which + " here:";
            auto& error_handler = x3::get<error_handler_tag>(context).get();
            error_handler(x.where(), message);
            return x3::error_handler_result::fail;
        }

        ///////////////////////////////////////////////////////////////////////
        //  Our annotation handler
        ///////////////////////////////////////////////////////////////////////

        // tag used to get the position cache from the context
        template <typename T, typename Iterator, typename Context>
        inline void
            error_handler_base::on_success(Iterator const& first, Iterator const& last
        , T& ast, Context const& context)
        {
            auto& position_cache = x3::get<error_handler_tag>(context).get();
            position_cache.annotate(ast, first, last);
        }
    }
}
