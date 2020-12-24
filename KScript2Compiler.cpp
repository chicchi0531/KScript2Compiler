// KScript2Compiler.cpp : このファイルには 'main' 関数が含まれています。プログラム実行の開始と終了がそこで行われます。
//
#include <iostream>
#include <iterator>
#include <algorithm>
#include <sstream>

#include "ast.h"
#include "grammar.h"
#include "error_handler.h"
#include "config.h"
#include "printer.h"

#include "test.h"

namespace fs = boost::filesystem;
namespace testing = boost::spirit::x3::testing;

auto parse = [](std::string const& source, fs::path input_path)-> std::string
{
    std::stringstream out;

    using kscript2::parser::iterator_type;
    iterator_type iter(source.begin());
    iterator_type const end(source.end());

    // Our AST
    kscript2::ast::unit ast;

    // Our error handler
    using boost::spirit::x3::with;
    using kscript2::parser::error_handler_type;
    using kscript2::parser::error_handler_tag;
    error_handler_type error_handler(iter, end, out, input_path.string()); // Our error handler

    // Our parser
    auto const parser =
        // we pass our error handler to the parser so we can access
        // it later on in our on_error and on_sucess handlers
        with<error_handler_tag>(std::ref(error_handler))
        [
            kscript2::unit()
        ];

    // Go forth and parse!
    namespace x3 = boost::spirit::x3;
    bool success = phrase_parse(iter, end, parser, kscript2::parser::comment::skipper, ast);

    if (success)
    {
        if (iter != end)
            error_handler(iter, "Error! Expecting end of input here: ");
        else
            kscript2::ast::ast_analyzer{ out }(ast);
    }

    return out.str();
};

int num_files_tested = 0;
auto compare = [](fs::path input_path, fs::path expect_path)
{
    std::string output = parse(testing::load(input_path), input_path);
    std::cout << output << std::endl;

    ++num_files_tested;
};

int main(int argc, char* argv[])
{
    /*if (argc < 2)
    {
        std::cout << "usage: " << fs::path(argv[0]).filename() << " path/to/test/files" << std::endl;
        return -1;
    }

    std::cout << "===================================================================================================" << std::endl;
    std::cout << "Testing: " << fs::absolute(fs::path(argv[1])) << std::endl;
    int r = testing::for_each_file(fs::path(argv[1]), compare);
    if (r == 0)
        std::cout << num_files_tested << " files tested." << std::endl;
    std::cout << "===================================================================================================" << std::endl;
    return r;*/

    if (argc < 2)
    {
        std::cout << "ファイル名を引数に指定してください。" << std::endl;
        return -1;
    }

    std::cout << "start test======================" << std::endl;

    for (int i = 1; i < argc; i++)
    {
        std::string output = parse(testing::load(argv[i]), argv[i]);
        std::cout << output << std::endl;
    }
    std::cout << "end test========================" << std::endl;

    return 0;
}
