#include "ast_analyzer.h"

using namespace kscript2::ast;
using namespace kscript2;
using namespace boost::spirit::x3;


template<typename T>
inline bool IsNil(T value)
{
    return value.get().which() == 0;
}


void ast_analyzer::operator()(unit const& ast) const
{
    for (auto const& entry : ast.entries)
    {
        boost::apply_visitor(ast_analyzer(compiler_,positions_), entry);
    }
}

void ast_analyzer::operator()(function_pre_def const& ast) const
{
    // 関数宣言
    compiler_.DeclFunction(ast.attr_type, ast.return_type, ast.name.name, ast.args, ast);
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
    compiler_.DefineFunction(ast.return_type, ast.name.name, ast.args, ast.states, ast);
}

// 関数呼び出し 次のような命令を生成する
// ---------------------------------------
// push expr            引数（任意の数）
// push_const arg_size  引数の数
// call                 呼び出し命令
// ---------------------------------------
void ast_analyzer::operator()(function_call const& ast) const
{
    // 関数呼び出し
    const FunctionTag* tag = compiler_.GetFunctionTag(ast.name.name);
    if (tag == nullptr)
    {
        compiler_.error("関数　" + ast.name.name + "は定義されていません", ast);
        return;
    }

    int arg_size = ast.args.size();
    if (tag->ArgSize() != arg_size)
    {
        compiler_.error("引数の数が合いません", ast);
    }

    // 引数をpush
    if (tag->ArgSize() == arg_size)
    {
        int index = 0;
        for (auto const& arg : ast.args)
        {
            int type = tag->GetArg(index++);
            
            auto a = ast_analyzer(compiler_, positions_);
            a(arg);
            if (compiler_.GetAstReturn() != type)
            {
                compiler_.error("引数の型が合いません。", ast);
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

    compiler_.SetAstReturn(tag->GetType());
}

void ast_analyzer::operator()(statements const& ast) const
{
    compiler_.BlockIn();
    for (auto const& s : ast)
    {
        boost::apply_visitor(ast_analyzer(compiler_, positions_), s);
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
    auto a = ast_analyzer(compiler_, positions_);
    a(ast.expression);

    // if文の最後へのラベルジャンプ
    int L1 = compiler_.MakeLabel();
    compiler_.OpJZero(L1);

    // if文の中身
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.if_state);

    // else文がある場合
    if (ast.else_state.get().which() != 0)
    {
        // if文の終わりにendifまでのジャンプ命令を生成
        int L2 = compiler_.MakeLabel();
        compiler_.OpJmp(L2);

        // L1ラベルを設定(elseの開始地点)
        compiler_.SetLabel(L1);

        // else 文の中身
        boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.else_state);

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
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast);
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
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.decl);
    // label L1
    compiler_.SetLabel(L1);
    // push expr
    auto ana_con = ast_analyzer(compiler_, positions_);
    ana_con(ast.condition);
    // jmp_zero L2
    compiler_.OpJZero(L2);
    // A
    boost::apply_visitor(ast_analyzer(compiler_, positions_),ast.state);
    // next
    auto ana_iter = ast_analyzer(compiler_, positions_);
    ana_iter(ast.iter);
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
    auto a = ast_analyzer(compiler_, positions_);
    a(ast.condition);
    // jmp_zero L2
    compiler_.OpJZero(L2);
    // A
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.state);
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
        if (!compiler_.JmpBreakLabel()) compiler_.error("breakがfor/while文の外にあります。", ast);
        break;
    case JUMP_CONTINUE:
        if (!compiler_.JmpContinueLabel()) compiler_.error("continueがfor/while文の外にあります。", ast);
        break;
    case JUMP_RETURN:
        int functype = ast.expression.first.get().which();
        if (compiler_.GetFunctionType() == TYPE_VOID)
        {
            if (functype != 0) compiler_.error("void型の関数に戻り値が設定されています。", ast);
            compiler_.OpJReturn();
        }
        else
        {
            if (functype == 0) compiler_.error("関数の戻り値がありません。", ast);
            else
            {
                auto a = ast_analyzer(compiler_, positions_);
                a(ast.expression);
                int expr_type = compiler_.GetAstReturn();
                if (expr_type != compiler_.GetFunctionType()) compiler_.error("戻り値の型が合いません。", ast);
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
        boost::apply_visitor(ast_analyzer(compiler_, positions_), s);
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
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.name);

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
    // 空行の場合はスキップ
    if (ast.msg.size() == 1 && ast.new_page == 0 && ast.msg[0] == L"")
        return;

    // %が閉じられているかのチェック
    // 偶数の場合は閉じられていない
    if (ast.msg.size() % 2 == 0)
    {
        compiler_.error("%が閉じられていません", ast);
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
            auto a = ast_analyzer(compiler_, positions_);
            a(i);
        }
        else
        {
            this->operator()(m);
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
class value_selector
{
    compiler& c_;
public:
    value_selector(compiler& c):c_(c){}
    void operator()(ast::nil var) { auto v = VMVariableValue(); c_.SetGlobalValue(v); }
    void operator()(int var)
    {
        auto v = VMVariableValue();
        v.type = TYPE_INTEGER;
        v.ival = var;
        c_.SetGlobalValue(v);
    }
    void operator()(double var)
    {
        auto v = VMVariableValue();
        v.type = TYPE_FLOAT;
        v.fval = var;
        c_.SetGlobalValue(v);
    }
    void operator()(std::wstring var)
    {
        auto v = VMVariableValue();
        v.type = TYPE_STRING;
        v.sval = var;
        c_.SetGlobalValue(v);
    }
};

void ast_analyzer::operator()(global_declaration const& ast) const
{
    // 変数宣言
    this->operator()(ast.decl);

    // 初期値を代入。グローバル変数なので、値は計算してバイナリに出力する
    boost::apply_visitor(value_selector(compiler_), ast.value);
}
void ast_analyzer::operator()(declarator const& ast) const
{
    compiler_.AddValue(ast.type, ast.identifier_.name, ast);
}
void ast_analyzer::operator()(declaration const& ast) const
{
    // 変数宣言
    this->operator()(ast.decl);

    // 初期化値を代入
    if (ast.expression.first.get().which() != 0)
    {
        // 内部でassignを構成して呼び出す
        assign a;
        a.left = ast.decl.identifier_;
        a.sign = parser::OP_ASSIGN;
        a.right = ast.expression;

        auto ana = ast_analyzer(compiler_, positions_);
        ana(a);
    }

}

// expression

//単項演算
void ast_analyzer::operator()(signed_ const& ast) const
{
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.operand_);

    switch (ast.sign)
    {
    case OP_NEG:
        switch (compiler_.GetAstReturn())
        {
        case TYPE_INTEGER: compiler_.OpINeg(); break;
        case TYPE_FLOAT: compiler_.OpFNeg(); break;
        default:
            compiler_.error("文字列定数/変数に単項演算子が使用されました。単項演算子は数値型のみに有効です。", ast);
        }
        break;
    default: compiler_.error("不明な符号が使用されました。 OPCODE=" + std::to_string(ast.sign), ast);
    }
}

//二項演算
void ast_analyzer::operator()(operation const& ast) const
{
    //左辺の型
    int left_type = compiler_.GetAstReturn();
    auto dummyindex = compiler_.DummyOp(); //ダミー命令を置いておき、キャストする場合はあとから命令を書き換える

    //右辺の処理
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.operand_);
    int right_type = compiler_.GetAstReturn();

    //型チェック
    if(left_type != right_type)
    {
        if(left_type == TYPE_STRING ||
        right_type == TYPE_STRING)
        {
            compiler_.error("文字列と数値の演算はサポートされていません。", ast);
            return;
        }
    }

    //演算結果の型ぎめ
    int return_type = TYPE_INTEGER;
    if (right_type == TYPE_FLOAT || left_type == TYPE_FLOAT) return_type = TYPE_FLOAT; //どっちかFLOATの場合はFLOATにキャスト
    else if (right_type == TYPE_STRING && left_type == TYPE_STRING) return_type = TYPE_STRING;

    // キャスト命令を差し込む（左辺）
    if (return_type == TYPE_FLOAT && left_type == TYPE_INTEGER)
        compiler_.ReplaceOp(VMCode(VM_CAST_ITOF), dummyindex);
    else
        compiler_.EraseOp(dummyindex);

    // キャスト命令を差し込む（右辺）
    if (return_type == TYPE_FLOAT && right_type == TYPE_INTEGER)
        compiler_.CastItoF();

    //演算結果の型を設定
    compiler_.SetAstReturn(return_type);

    //演算子の処理
    switch (return_type)
    {
    case TYPE_INTEGER:
        switch (ast.sign)
        {
        case OP_ADD: compiler_.OpIAdd(); break;
        case OP_SUB: compiler_.OpISub(); break;
        case OP_MUL: compiler_.OpIMul(); break;
        case OP_DIV: compiler_.OpIDiv(); break;
        case OP_MOD: compiler_.OpIMod(); break;
        case OP_EQ: compiler_.OpIEqu(); break;
        case OP_NE: compiler_.OpINeq(); break;
        case OP_GE: compiler_.OpIGe(); break;
        case OP_GT: compiler_.OpIGt(); break;
        case OP_LE: compiler_.OpILe(); break;
        case OP_LT: compiler_.OpILt(); break;
        case OP_LOGAND: compiler_.OpLogAnd(); break;
        case OP_LOGOR: compiler_.OpLogOr(); break;
        default: compiler_.error("不明な符号が使用されました。 OPCODE=" + std::to_string(ast.sign), ast);
        }
        break;

    case TYPE_FLOAT:
        switch (ast.sign)
        {
        case OP_ADD: compiler_.OpFAdd(); break;
        case OP_SUB: compiler_.OpFSub(); break;
        case OP_MUL: compiler_.OpFMul(); break;
        case OP_DIV: compiler_.OpFDiv(); break;
        case OP_EQ: compiler_.OpFEqu(); break;
        case OP_NE: compiler_.OpFNeq(); break;
        case OP_GE: compiler_.OpFGe(); break;
        case OP_GT: compiler_.OpFGt(); break;
        case OP_LE: compiler_.OpFLe(); break;
        case OP_LT: compiler_.OpFLt(); break;
        case OP_LOGAND: compiler_.OpLogAnd(); break;
        case OP_LOGOR: compiler_.OpLogOr(); break;
        default: compiler_.error("不明な符号が使用されました。 OPCODE=" + std::to_string(ast.sign), ast);
        }
        break;

    case TYPE_STRING:
        switch (ast.sign)
        {
        case OP_ADD: compiler_.OpStrAdd(); break;
        case OP_EQ: compiler_.OpStrEq(); break;
        case OP_NE: compiler_.OpStrNe(); break;
        case OP_LOGAND: compiler_.OpLogAnd(); break;
        case OP_LOGOR: compiler_.OpLogOr(); break;
        default: compiler_.error("不明な符号が使用されました。 OPCODE=" + std::to_string(ast.sign), ast);
        }
        break;
    }
}

// expr
void ast_analyzer::operator()(expr const& ast) const
{
    // nilの場合は棄却
    if (ast.first.get().which() == 0)return;

    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast.first);
    for (auto const& ope : ast.operations)
    {
        auto a = ast_analyzer(compiler_, positions_);
        a(ope);
    }
}

// 代入処理
int pop_variable(compiler& c, const identifier& left)
{
    const ValueTag* tag = c.GetValueTag(left.name);
    if (tag == 0)
    {
        c.error("変数　" + left.name + "　は定義されていません。", left);
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
    // 左辺の型を拾ってくる
    const ValueTag* tag = compiler_.GetValueTag(ast.left.name);
    if (tag == 0)
    {
        compiler_.error("変数　" + ast.left.name + "　は定義されていません。", ast);
        return;
    }
    int left_type = tag->type_;

    // 計算式を含む演算子の場合、左辺をpush
    int dummyindex;
    if (ast.sign != OP_ASSIGN)
    {
        auto a = ast_analyzer(compiler_, positions_);
        a(ast.left);
        dummyindex = compiler_.DummyOp();
    }

    // 右辺をpush
    auto a = ast_analyzer(compiler_, positions_);
    a(ast.right);
    int right_type = compiler_.GetAstReturn();

    //型チェック
    if (left_type != right_type)
    {
        if (left_type == TYPE_STRING ||
            right_type == TYPE_STRING)
        {
            compiler_.error("文字列と数値の演算はサポートされていません。", ast);
            return;
        }
    }

    // 演算する場合の型チェック
    int return_type = right_type;
    if (ast.sign != OP_ASSIGN)
    {
        //演算結果の型ぎめ
        if (right_type == TYPE_FLOAT || left_type == TYPE_FLOAT) return_type = TYPE_FLOAT; //どっちかFLOATの場合はFLOATにキャスト
        else if (right_type == TYPE_STRING && left_type == TYPE_STRING) return_type = TYPE_STRING;

        // キャスト命令を差し込む（左辺）
        if (return_type == TYPE_FLOAT && left_type == TYPE_INTEGER)
            compiler_.ReplaceOp(VMCode(VM_CAST_ITOF), dummyindex);
        else
            compiler_.EraseOp(dummyindex);

        // キャスト命令を差し込む（右辺）
        if (return_type == TYPE_FLOAT && right_type == TYPE_INTEGER)
            compiler_.CastItoF();
    }
    

    // 型に応じて演算命令を設定
    switch (return_type)
    {
    case TYPE_INTEGER:
        switch (ast.sign)
        {
        case OP_ADD_ASSIGN: compiler_.OpIAdd(); break;
        case OP_SUB_ASSIGN: compiler_.OpISub(); break;
        case OP_MUL_ASSIGN: compiler_.OpIMul(); break;
        case OP_DIV_ASSIGN: compiler_.OpIDiv(); break;
        case OP_MOD_ASSIGN: compiler_.OpIMod(); break;
        }

        // 代入命令
        if (left_type == TYPE_FLOAT) compiler_.CastItoF();
        pop_variable(compiler_, ast.left);
        break;

    case TYPE_FLOAT:
        switch (ast.sign)
        {
        case OP_ADD_ASSIGN: compiler_.OpFAdd(); break;
        case OP_SUB_ASSIGN: compiler_.OpFSub(); break;
        case OP_MUL_ASSIGN: compiler_.OpFMul(); break;
        case OP_DIV_ASSIGN: compiler_.OpFDiv(); break;
        case OP_MOD_ASSIGN: compiler_.error("float型は余算に対応していません", ast); break;
        }

        // 代入命令
        if (left_type == TYPE_INTEGER) compiler_.CastFtoI();
        pop_variable(compiler_, ast.left);
        break;

    case TYPE_STRING:
        // string型の式を代入
        switch (ast.sign)
        {
        case OP_ADD_ASSIGN:
            compiler_.OpStrAdd();
            break;
        case OP_ASSIGN:
            break;

        default:
            compiler_.error("文字列では許されない計算です。", ast);
            break;
        }
        
        // 代入命令
        pop_variable(compiler_, ast.left);
        break;
    }


}
/*
void ast_analyzer::operator()(assign_list const& ast) const
{
    for (auto const& a : ast)
    {
        this->operator()(a);
    }
}
*/

// constant
void ast_analyzer::operator()(constant const& ast) const
{
    boost::apply_visitor(ast_analyzer(compiler_, positions_), ast);
}

// identifier
void ast_analyzer::operator()(identifier const& ast) const
{
    if (ast.name == "") return;

    const ValueTag* tag = compiler_.GetValueTag(ast.name);
    if (tag == nullptr)
        compiler_.error("変数　" + ast.name + "　は定義されていません。", ast);
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
        compiler_.SetAstReturn(tag->type_);
    }
}
void ast_analyzer::operator()(std::wstring const& ast) const
{
    compiler_.PushString(ast);
    compiler_.SetAstReturn(TYPE_STRING);
}
void ast_analyzer::operator()(int value) const
{
    compiler_.PushConst(value);
    compiler_.SetAstReturn(TYPE_INTEGER);
}
void ast_analyzer::operator()(double value) const
{
    compiler_.PushDouble(value);
    compiler_.SetAstReturn(TYPE_FLOAT);
}
