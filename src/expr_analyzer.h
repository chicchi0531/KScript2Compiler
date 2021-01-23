#pragma once

#include "symbols.h"
#include "ast.h"
#include "compiler.h"

namespace kscript2 {
	namespace ast
	{
		struct expr_return_value
		{
			int type;
			int ival;
			double dval;
			std::wstring sval;

			void operator= (expr_return_value& value)
			{
				type = value.type;
				ival = value.ival;
				dval = value.dval;
				sval = value.sval;
			}
		};

		// ®‚ğ‰ğÍ‚µ‚ÄŒvZŒ‹‰Ê‚ğ•Ô‚·ƒNƒ‰ƒX
		class expr_analyzer
		{
			compiler& compiler_;

		public:
			expr_return_value result;

			expr_analyzer(compiler& c);
			void operator()(signed_ const& ast) const;
			void operator()(operation const& ast) const;
			void operator()(expr const& ast) const;
			void operator()(constant const& ast) const;
			void operator()(identifier const& ast) const;
			void operator()(function_call const& ast) const;

			void operator()(std::wstring const& text) const;
			void operator()(int value) const;
			void operator()(double value) const;
		};
	}
}