// KScript2Compiler.cpp : このファイルには 'main' 関数が含まれています。プログラム実行の開始と終了がそこで行われます。
//

#include "compiler.h"

namespace fs = std::filesystem;

int main(int argc, char* argv[])
{
    kscript2::compiler compiler;

    if (argc < 2)
    {
        std::cout << "ファイル名を引数に指定してください。" << std::endl;
        return -1;
    }


    for (int i = 1; i < argc; i++)
    {
        std::cout << "コンパイル開始　(" << i << "/" << argc-1 <<") ======================" << std::endl;

        auto input_path = fs::path(argv[i]);
        auto out_path = fs::path(argv[i]).replace_extension("ksobj");

        if(!compiler.compile(input_path, out_path))
        {
            std::cerr << input_path << " : コンパイルに失敗しました。エラーメッセージを確認してください。" << std::endl;
            return -1;
        }

        std::cout << "コンパイル完了　(" << i << "/" << argc-1 << ") ======================" << std::endl;
    }

    return 0;
}
