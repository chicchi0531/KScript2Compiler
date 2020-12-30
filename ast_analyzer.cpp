#include "ast_analyzer.h"

using namespace kscript2::ast;
using namespace kscript2;


void ast_analyzer::operator()(unit const& ast) const
{
    for (auto const& entry : ast.entries)
    {
        boost::apply_visitor(ast_analyzer(compiler_), entry);
    }
}

void ast_analyzer::operator()(function_pre_def const& ast) const
{
    // 関数宣言
    compiler_.DefineFunction(ast.return_type, ast.name.name, ast.args);
}

// 関数定義　次のような命令を生成する
// ---------------------------------------
// label                関数のエントリポイント
// (*statement)         関数の中身
// push expr            戻り値（voidの場合は省略）
// jump_return          リターン
// ---------------------------------------
void ast_analyzer::operator()(function_def const& ast) const
{
    // 関数定義
    compiler_.AddFunction(ast.return_type, ast.name.name, ast.args, ast.states);
}

// 関数呼び出し 次のような命令を生成する
// ---------------------------------------
// push expr            引数（任意の数）
// push_const arg_size  引数の数
// call                 呼び出し命令
// ---------------------------------------
int ast_analyzer::operator()(function_call const& ast) const
{
    // 関数呼び出し
    const FunctionTag* tag = compiler_.GetFunctionTag(ast.name.name);
    if (tag == nullptr)
    {
        compiler_.error("関数　" + ast.name.name + "は定義されていません");
    }

    int arg_size = ast.args.size();
    if (tag->ArgSize() != arg_size)
    {
        compiler_.error("引数の数が合いません");
    }

    // 引数をpush
    if (tag->ArgSize() == arg_size)
    {
        int index = 0;
        for (auto const& arg : ast.args)
        {
            int type = tag->GetArg(index++);
            
            if (ast_analyzer(compiler_)(arg) != type)
            {
                compiler_.error("引数の型が合いません。");
            }
        }
    }

    // 引数の数をpush
    compiler_.PushConst(arg_size);

    // システム関数かどうかで命令を分ける
    if (tag->IsSystem())
    {
        compiler_.OpSysCall(tag->GetIndex());
    }
    else
    {
        compiler_.OpCall(tag->GetIndex());
    }

    return tag->GetType();
}

void ast_analyzer::operator()(statements const& ast) const
{
    compiler_.BlockIn();
    for (auto const& entry : ast)
    {
        boost::apply_visitor(ast_analyzer(compiler_), entry);
    }
    compiler_.BlockOut();
}

// IF文処理
// if(expr) A else B
// -------------------------------------
// push expr        条件式
// jump_zero L1     elseへのジャンプ命令
//  A               if文の中身
// jump L2          endifへのジャンプ命令
// label L1         else文の先頭
//  B               else文の中身
// label L2         endifラベル
// -------------------------------------
void ast_analyzer::operator()(section_statement const& ast) const
{
    // 条件式
    int type = ast_analyzer(compiler_)(ast.expression);

    // if文の最後へのラベルジャンプ
    int L1 = compiler_.MakeLabel();
    compiler_.OpJZero(L1);

    // if文の中身
    boost::apply_visitor(ast_analyzer(compiler_), ast.if_state);

    // else文がある場合
    if (ast.else_state.get().which() != 0)
    {
        // if文の終わりにendifまでのジャンプ命令を生成
        int L2 = compiler_.MakeLabel();
        compiler_.OpJmp(L2);

        // L1ラベルを設定(elseの開始地点)
        compiler_.SetLabel(L1);

        // else 文の中身
        boost::apply_visitor(ast_analyzer(compiler_), ast.else_state);

        // L2ラベルをセット
        compiler_.SetLabel(L2);
    }
    else
    {
        // elseがない場合はendifがL1ラベルと対応する
        compiler_.SetLabel(L1);
    }
}

// ループ文
void ast_analyzer::operator()(iteration_statement const& ast) const
{
    boost::apply_visitor(ast_analyzer(compiler_), ast);
}

// for文
// for(init; expr; next) A
// --------------------
// init
// label L1
// push expr
// jmp_zero L2
//  A
// next
// jmp L1
// label L2
// --------------------
void ast_analyzer::operator()(for_statement const& ast) const
{
    int L1 = compiler_.MakeLabel();
    int L2 = compiler_.MakeLabel();
    int break_label = compiler_.SetBreakLabel(L2);
    int continue_label = compiler_.SetContinueLabel(L1);

    // init
    boost::apply_visitor(ast_analyzer(compiler_), ast.decl);
    // label L1
    compiler_.SetLabel(L1);
    // push expr
    int type = ast_analyzer(compiler_)(ast.condition);
    // jmp_zero L2
    compiler_.OpJZero(L2);
    // A
    boost::apply_visitor(ast_analyzer(compiler_),ast.state);
    // next
    boost::apply_visitor(ast_analyzer(compiler_), ast.iter);
    // jmp L1
    compiler_.OpJmp(L1);
    // label L2
    compiler_.SetLabel(L2);
    // ifから抜けたので、breakのジャンプ先を1つ前に戻す
    compiler_.SetBreakLabel(break_label);
}

// while文
// while(expr) A
// ---------------------
// label L1
// push expr
// jmp_zero L2
//  A
// jmp L1
// label L2
// ---------------------
void ast_analyzer::operator()(while_statement const& ast) const
{
    int L1 = compiler_.MakeLabel();
    int L2 = compiler_.MakeLabel();
    int break_label = compiler_.SetBreakLabel(L2);
    int continue_label = compiler_.SetContinueLabel(L1);

    // label L1
    compiler_.SetLabel(L1);
    // push expr
    int type = ast_analyzer(compiler_)(ast.condition);
    // jmp_zero L2
    compiler_.OpJZero(L2);
    // A
    boost::apply_visitor(ast_analyzer(compiler_), ast.state);
    // jmp L1
    compiler_.OpJmp(L1);
    // label L2
    compiler_.SetLabel(L2);
    // break,continueの参照先を1つ前に戻す
    compiler_.SetBreakLabel(break_label);
    compiler_.SetContinueLabel(continue_label);
}

// ジャンプ文
void ast_analyzer::operator()(jump_statement const& ast) const
{
    switch (ast.ope)
    {
    case JUMP_BREAK:
        if (!compiler_.JmpBreakLabel()) compiler_.error("breakがfor/while文の外にあります。");
        break;
    case JUMP_CONTINUE:
        if (!compiler_.JmpContinueLabel()) compiler_.error("continueがfor/while文の外にあります。");
        break;
    case JUMP_RETURN:
        int functype = ast.expression.first.get().which();
        if (compiler_.GetFunctionType() == TYPE_VOID)
        {
            if (functype != 0) compiler_.error("void型の関数に戻り値が設定されています。");
            compiler_.OpJReturn();
        }
        else
        {
            if (functype == 0) compiler_.error("関数の戻り値がありません。");
            else
            {
                int expr_type = ast_analyzer(compiler_)(ast.expression);
                if (expr_type != compiler_.GetFunctionType()) compiler_.error("戻り値の型が合いません。");
            }
            compiler_.OpJReturnV();
        }
        break;
    }
}

// novel statement
void ast_analyzer::operator()(novel_block const& ast) const
{
    for (auto const& s : ast)
    {
        boost::apply_visitor(ast_analyzer(compiler_), s);
    }
}

// novel name
// -----------------
// push_string      名前の文字列もしくは変数
// novel_name
// -----------------
void ast_analyzer::operator()(novel_name_statement const& ast) const
{
    // push string
    boost::apply_visitor(ast_analyzer(compiler_), ast.name);

    // novel_name
    compiler_.OpNovelName();
}

// novel msg
// -----------------------
// push_string
// push_variable ...                    push_stringと繰り返し
// push_const                           上記のメッセージ数をpush
// novel_msg                            メッセージ表示命令
// novel_new_line/novel_new_page        新ページ、もしくは新ライン
// -----------------------
void ast_analyzer::operator()(novel_msg_statement const& ast) const
{
    if (ast.msg.size() == 0 && ast.new_page == 0) return;

    // %が閉じられているかのチェック
    // 偶数の場合は閉じられていない
    if (ast.msg.size() % 2 == 0)
    {
        std::cerr << ast.id_last << ": %が閉じられていません。" << std::endl;
    }

    // push_string
    bool is_identifier = false;
    int count = 0;
    for (auto const& m : ast.msg)
    {
        if (is_identifier)
        {
            identifier i;
            std::string s(m.begin(), m.end());
            i.name = s;
            int type = ast_analyzer(compiler_)(i);
        }
        else
        {
            ast_analyzer(compiler_)(m);
        }

        is_identifier = !is_identifier;
        count++;
    }

    // push_const
    compiler_.PushConst(count);

    // novel_msg
    compiler_.OpNovelMsg();

    //改行、改ページ処理
    switch (ast.new_page)
    {
    case parser::NOVEL_EMPTY:
    case parser::NOVEL_NEWLINE:
        compiler_.OpNovelNewLine();
        break;
    case parser::NOVEL_NEWPAGE:
        compiler_.OpNovelNewPage();
        break;
    }
}

// 変数宣言
void ast_analyzer::operator()(declarator const& ast) const
{
    compiler_.AddValue(ast.type, ast.identifier_.name);
}
void ast_analyzer::operator()(declaration const& ast) const
{
    // 変数宣言
    ast_analyzer(compiler_)(ast.decl);

    // 初期化値を代入
    if (ast.expression.first.get().which() != 0)
    {
        // 内部でassignを構成して呼び出す
        assign a;
        a.left = ast.decl.identifier_;
        a.sign = parser::OP_ASSIGN;
        a.right = ast.expression;
        ast_analyzer(compiler_)(a);
    }
}

// expression
void ast_analyzer::operator()(signed_ const& ast) const
{
    boost::apply_visitor(ast_analyzer(compiler_), ast.operand_);

    switch (ast.sign)
    {
    case OP_NEG: compiler_.OpNeg(); break;
    default: compiler_.error("不明な符号が使用されました。 OPCODE=" + ast.sign);
    }
}
void ast_analyzer::operator()(operation const& ast) const
{
    boost::apply_visitor(ast_analyzer(compiler_), ast.operand_);

    switch (ast.sign)
    {
    case OP_ADD: compiler_.OpAdd(); break;
    case OP_SUB: compiler_.OpSub(); break;
    case OP_MUL: compiler_.OpMul(); break;
    case OP_DIV: compiler_.OpDiv(); break;
    case OP_EQ: compiler_.OpEqu(); break;
    case OP_NE: compiler_.OpNeq(); break;
    case OP_GE: compiler_.OpGe(); break;
    case OP_GT: compiler_.OpGt(); break;
    case OP_LE: compiler_.OpLe(); break;
    case OP_LT: compiler_.OpLt(); break;
    case OP_LOGAND: compiler_.OpLogAnd(); break;
    case OP_LOGOR: compiler_.OpLogOr(); break;
    default: compiler_.error("不明な符号が使用されました。 OPCODE=" + ast.sign);
    }
}

// expr
int ast_analyzer::operator()(expr const& ast) const
{
    // nilの場合は棄却
    if (ast.first.get().which() == 0)return 1;

    boost::apply_visitor(ast_analyzer(compiler_), ast.first);
    for (auto const& ope : ast.operations)
    {
        ast_analyzer(compiler_)(ope);
    }
    return 0;
}

// 代入処理
int pop_variable(compiler& c, const identifier& left)
{
    const ValueTag* tag = c.GetValueTag(left.name);
    if (tag == 0)
    {
        c.error("変数　" + left.name + "　は定義されていません。");
    }
    else
    {
        if (tag->global_)
        {
            c.PopValue(tag->addr_);
        }
        else
        {
            c.PopLocal(tag->addr_);
        }
        return tag->type_;
    }
    return TYPE_INTEGER;
}

// assign命令
// a = b
// ----------------
// push b
// pop a
// ----------------
// a += b
// ----------------
// push a
// push b
// add
// pop a
// ----------------
void ast_analyzer::operator()(assign const& ast) const
{
    // assign 
    int left_type = 0;
    if (ast.sign != OP_ASSIGN)
        left_type = ast_analyzer(compiler_)(ast.left);

    // int型の式を代入
    if (ast_analyzer(compiler_)(ast.right) == TYPE_INTEGER)
    {
        switch (ast.sign)
        {
        case OP_ADD_ASSIGN: compiler_.OpAdd(); break;
        case OP_SUB_ASSIGN: compiler_.OpSub(); break;
        case OP_MUL_ASSIGN: compiler_.OpMul(); break;
        case OP_DIV_ASSIGN: compiler_.OpDiv(); break;
        case OP_MOD_ASSIGN: compiler_.OpMod(); break;
        }

        if (pop_variable(compiler_, ast.left) != TYPE_INTEGER)
        {
            compiler_.error("文字列型に整数を代入しています。");
        }
        return;
    }

    // string型の式を代入
    switch (ast.sign)
    {
    case OP_ADD_ASSIGN:
        compiler_.OpStrAdd();
        break;

    case OP_ASSIGN:
        break;

    default:
        compiler_.error("文字列では許されない計算です。");
        break;
    }
    if (pop_variable(compiler_, ast.left) != TYPE_STRING)
    {
        compiler_.error("整数型に文字列を代入しています。");
    }
}
void ast_analyzer::operator()(assign_list const& ast) const
{
    for (auto const& a : ast)
    {
        ast_analyzer(compiler_)(a);
    }
}

// constant
void ast_analyzer::operator()(constant const& ast) const
{
    boost::apply_visitor(ast_analyzer(compiler_), ast);
}

// identifier
int ast_analyzer::operator()(identifier const& ast) const
{
    const ValueTag* tag = compiler_.GetValueTag(ast.name);
    if (tag == nullptr)
        compiler_.error("変数　" + ast.name + "　は定義されていません。");
    else
    {
        if (tag->global_)
        {
            // グローバル変数
            compiler_.PushValue(tag->addr_);
        }
        else
        {
            // ローカル変数
            compiler_.PushLocal(tag->addr_);
        }
    }
    return tag->type_;
}
void ast_analyzer::operator()(std::wstring const& ast) const
{
    compiler_.PushString(ast);
}
void ast_analyzer::operator()(int value) const
{
    compiler_.PushConst(value);
}
void ast_analyzer::operator()(double value) const
{
    compiler_.PushDouble(value);
}
