// KScript2Compiler.cpp : このファイルには 'main' 関数が含まれています。プログラム実行の開始と終了がそこで行われます。
//

#include "compiler.h"
#include <boost/program_options.hpp>

namespace fs = std::filesystem;
namespace po = boost::program_options;

int main(int argc, char* argv[])
{
	kscript2::compiler compiler;

	// コマンドライン引数の処理
	po::options_description opt("オプション");
	opt.add_options()
		("help,h", "ヘルプを表示")
		("output,o", po::value<std::string>(), "出力ファイル名")
		("input,i", po::value<std::string>(), "入力ファイル名");

	// コマンド引数の登録
	po::variables_map m;
	try
	{
		po::store(po::parse_command_line(argc, argv, opt), m);
	}
	catch (const po::error_with_option_name& e)
	{
		std::cout << e.what() << std::endl;
		return -1;
	}
	po::notify(m);

	// 各引数の処理
	if (m.count("help"))
	{
		std::cout << opt << std::endl;
		return 0;
	}

	if (!m.count("input"))
	{
		std::cerr << "コンパイルするファイルを-iオプションで指定してください。" << std::endl;
		return -1;
	}
	std::string input = m["input"].as<std::string>();

	std::string output;
	if (!m.count("output"))
	{
		std::cout << "出力ファイル名が指定されていません。out.ksobjとして出力します。" << std::endl;
		output = "out.ksobj";
	}
	else
	{
		output = m["output"].as<std::string>();
	}

	// コンパイル開始
	std::cout << "コンパイル開始======================" << std::endl;

	auto input_path = fs::path(input);
	auto out_path = fs::path(output);

	if (!compiler.compile(input_path, out_path))
	{
		std::cerr << input_path << " : コンパイルに失敗しました。エラーメッセージを確認してください。" << std::endl;
		return -1;
	}

	std::cout << "コンパイル完了======================" << std::endl;

	return 0;
}
