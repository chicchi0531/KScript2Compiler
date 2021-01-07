
#include "grammar_def.h"
#include "config.h"

namespace kscript2 {
    namespace parser
    {
        BOOST_SPIRIT_INSTANTIATE(
            unit_type, iterator_type, context_type);
    }
}
