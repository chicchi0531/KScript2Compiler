// KScript2Compiler.cpp : このファイルには 'main' 関数が含まれています。プログラム実行の開始と終了がそこで行われます。
//

#include "compiler.h"

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

        if (!compiler.compile(argv[i]))
        {
            std::cout << argv[i] << ": コンパイルに失敗しました。エラーメッセージを確認してください。" << std::endl;
        }

        std::cout << "コンパイル終了　(" << i << "/" << argc-1 << ") ======================" << std::endl;
    }

    return 0;
}
