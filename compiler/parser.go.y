%{
// プログラムのヘッダを指定
package compiler

import(
  "fmt"
  "ks2/compiler/ast"
  "ks2/compiler/vm"
  cm "ks2/compiler/common"
)

var driver *vm.Driver
var lexer *Lexer

%}

%union {
  ival int
  fval float32
  sval string

  node vm.INode
  nodes []vm.INode
  assign *ast.Assign
  argList []*vm.Argument
  argument *vm.Argument
  stateBlock vm.IStateBlock
  statement vm.IStatement
}

// 非終端記号の定義
%type<ival> var_type
%type<ival> function_type

%type<node> expr const define_var function_call uni_expr

%type<stateBlock> statements
%type<statement> statement
%type<statement> expr_statement
%type<statement> assign_statement
%type<statement> define_var_statement
%type<statement> return_statement
%type<statement> function_call_statement
%type<statement> block
%type<statement> if_statement

%type<argument> arg_decl
%type<argList> arg_list
%type<nodes> args
%type<assign> assign

// 終端記号の定義
%token<ival> INUM
%token<fval> FNUM
%token<sval> IDENTIFIER STRING_LITERAL

%token<ival> PLUS MINUS ASTARISK SLASH PERCENT // + - * / %
%token<ival> P_EQ M_EQ A_EQ S_EQ // += -= *= /=
%token<ival> EQ NEQ GT GE LT LE // == != > >= < <=
%token<ival> AND OR // && ||
%token<ival> INCR DECR ASSIGN// ++ -- =

%token<ival> VAR INT FLOAT STRING VOID
%token<ival> IF ELSE SWITCH CASE DEFAULT FALLTHROUGH FOR BREAK CONTINUE FUNC RETURN IMPORT TYPE STRUCT SYSCALL

%token EOL
%token<ival> '(' ')' '{' '}' '[' ']' ',' ':' ';' '.'

// 演算の優先度の指定
%left '(' ')'
%left OR
%left AND
%nonassoc EQ, NEQ, GT, GE, LT, LE
%left PLUS, MINUS
%left ASTARISK, SLASH, PERCENT
%left INCR, DECR
%left NEG

%%

// 文法規則を指定
program
  : define_or_state
  | program define_or_state

define_or_state
  : eol
  | function_decl eol
  | global_decl eol

global_decl
  : VAR IDENTIFIER var_type              { driver.VariableTable.DefineInLocal(lexer.line, $2, $3) }
  | VAR IDENTIFIER var_type ASSIGN expr  { driver.Err.LogError(driver.Filename, lexer.line, cm.ERR_0026, "") }

function_decl
  : FUNC IDENTIFIER '(' arg_list ')' function_type block { driver.AddFunction(lexer.line,$6,$2,$4,$7) }

arg_list
  : { $$ = make([]*vm.Argument, 0) }
  | arg_decl { $$ = []*vm.Argument{$1} }
  | arg_list ',' arg_decl { $$ = append($1,$3) }

arg_decl
  : IDENTIFIER var_type { $$ = &vm.Argument{Name:$1, VarType:$2} }

//---------------------------
// statements
//---------------------------
block
  : '{' statements '}' { $$ = ast.MakeCompoundStatement($2) }

statements
  : statement { s := new(ast.StateBlock); $$ = s.AddStates($1) }
  | statements statement { $$ = $1.AddStates($2) }

statement
  : eol { $$ = nil }
  | expr_statement eol { $$ = $1 }
  | assign_statement eol { $$ = $1 }
  | define_var_statement eol { $$ = $1 }
  | return_statement eol { $$ = $1 }
  | function_call_statement eol { $$ = $1 }
  | if_statement eol { $$ = $1 }

expr_statement
  : uni_expr    { $$ = ast.MakeExprStatement($1, driver) }

assign_statement
  : assign      { $$ = ast.MakeAssignStatement($1) }

define_var_statement
  : define_var  { $$ = ast.MakeVarDefineStatement($1) }

return_statement
  : RETURN      { $$ = ast.MakeReturnStatement(nil, lexer.line, driver) }
  | RETURN expr { $$ = ast.MakeReturnStatement($2,  lexer.line, driver) }

function_call_statement
  : function_call { $$ = ast.MakeFunctionCallStatement($1, driver) }

if_statement
  : IF expr block                   { $$ = ast.MakeIfStatement($2, $3, nil, lexer.line, driver) }
  | IF expr block ELSE block        { $$ = ast.MakeIfStatement($2, $3, $5, lexer.line, driver) }
  | IF expr block ELSE if_statement { $$ = ast.MakeIfStatement($2, $3, $5, lexer.line, driver) }

//------------------------------
// expr
//------------------------------

assign
  : IDENTIFIER ASSIGN expr
  { 
    varNode := ast.MakeValueNode(lexer.line, $1, driver)
    $$ = ast.MakeAssign(lexer.line, varNode, $3, driver)
  }

define_var
  : VAR IDENTIFIER var_type             { $$ = ast.MakeVarDefineNode(lexer.line, $2, $3, driver) }
  | VAR IDENTIFIER var_type ASSIGN expr { $$ = ast.MakeVarDefineNodeWithAssign(lexer.line, $2, $3, $5, driver) }
  | VAR IDENTIFIER ASSIGN expr          { $$ = ast.MakeVarDefineNodeWithAssign(lexer.line, $2, cm.TYPE_UNKNOWN, $4, driver) }

expr
  : const
  | MINUS expr %prec NEG{ $$ = ast.MakeExprNode(lexer.line, $2, nil, ast.OP_NOT, driver)}
  | IDENTIFIER          { $$ = ast.MakeValueNode(lexer.line, $1, driver) }
  | uni_expr            { $$ = $1 }
  | function_call       { $$ = $1 }
  | expr PLUS expr      { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_ADD, driver) }
  | expr MINUS expr     { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_SUB, driver) }
  | expr ASTARISK expr  { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_MUL, driver)}
  | expr SLASH expr     { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_DIV, driver) }
  | expr PERCENT expr   { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_MOD, driver) }
  | expr EQ expr        { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_EQUAL, driver)}
  | expr NEQ expr       { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_NEQ, driver) }
  | expr GT expr        { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_GT, driver) }
  | expr GE expr        { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_GE, driver) }
  | expr LT expr        { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_LT, driver) }
  | expr LE expr        { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_LE, driver) }
  | expr AND expr       { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_AND, driver) }
  | expr OR expr        { $$ = ast.MakeExprNode(lexer.line, $1, $3, ast.OP_OR, driver) }
  | '(' expr ')'        { $$ = $2 }

uni_expr
  : expr INCR           { $$ = ast.MakeExprNode(lexer.line, $1, nil, ast.OP_INCR, driver) }
  | expr DECR           { $$ = ast.MakeExprNode(lexer.line, $1, nil, ast.OP_DECR, driver) }

function_call
  : IDENTIFIER '(' args ')' { $$ = ast.MakeFunctionNode(lexer.line, $1, $3, driver) }
  | SYSCALL '[' expr ']' '(' args ')' { $$ = ast.MakeSysCallNode(lexer.line, $3, $6, driver) }

args
  : { $$ = make([]vm.INode,0) }
  | expr { $$ = []vm.INode{$1} }
  | args ',' expr { $$ = append($1,$3) }

const
  : STRING_LITERAL  { $$ = ast.MakeSvalNode(lexer.line, $1, driver) }
  | INUM            { $$ = ast.MakeIvalNode(lexer.line, $1, driver) }
  | FNUM            { $$ = ast.MakeFvalNode(lexer.line, $1, driver) }

var_type
  : INT     { $$ = cm.TYPE_INTEGER }
  | FLOAT   { $$ = cm.TYPE_FLOAT }
  | STRING  { $$ = cm.TYPE_STRING }

function_type
  :         { $$ = cm.TYPE_VOID }
  | INT     { $$ = cm.TYPE_INTEGER }
  | FLOAT   { $$ = cm.TYPE_FLOAT }
  | STRING  { $$ = cm.TYPE_STRING }
  | VOID    { $$ = cm.TYPE_VOID }

eol
  : EOL

%%

func Parse (filename string, source string) int {
  
  err := &cm.ErrorHandler{ErrorCount:0,WarningCount:0}
  driver = new(vm.Driver)
  driver.Init(filename, err)

  // パース処理
  lexer = &Lexer{src: source, position:0, readPosition:0, line:1, filename:filename, err:err}
	yyParse(lexer)

  // ラベル設定
  driver.LabelSettings()

  fmt.Println("Parse End.")

  // パース結果出力
  driver.Dump()

  return 0
}

// 外部用
func GetErrorCount() int{
  return driver.Err.ErrorCount
}
func GetWarningCount() int{
  return driver.Err.WarningCount
}
