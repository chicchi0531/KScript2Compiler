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
  argList []*vm.Argument
  argument *vm.Argument
  stateBlock vm.IStateBlock
}

// 非終端記号の定義
%type<ival> var_type function_type
%type<node> statement expr assign const define_var uni_expr function_call return_statement
%type<stateBlock> statements
%type<argument> arg_decl
%type<argList> arg_list
%type<nodes> args

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
%token<ival> IF ELSE SWITCH CASE DEFAULT FALLTHROUGH FOR BREAK CONTINUE FUNC RETURN IMPORT TYPE STRUCT 

%token<ival> EOL

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
  : EOL
  | function_decl
  | global_decl

global_decl
  : VAR IDENTIFIER var_type { ast.MakeNodeDecl }
  | VAR IDENTIFIER var_type ASSIGN uni_expr { n := &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line, Right:$5, Driver:driver} }; n.Push()}
  | VAR IDENTIFIER ASSIGN uni_expr { n :=  &ast.NodeDecl{Name:$2, VarType:cm.TYPE_UNKNOWN, Node:ast.Node{Lineno:lexer.line,Right:$4, Driver:driver}}; n.Push()}

uni_expr
  : const
  | MINUS INUM { $$ = ast.MakeIvalNode(lexer.line, -$2, driver) }
  | MINUS FNUM { $$ = ast.MakeFvalNode(lexer.line, -$2, driver)}

function_decl
  : FUNC IDENTIFIER '(' arg_list ')' function_type '{' statements '}' { driver.AddFunction(lexer.line,$6,$2,$4,$8) }

arg_list
  : arg_decl { $$ = []*vm.Argument{$1} }
  | arg_list ',' arg_decl { $$ = append($1,$3) }

arg_decl
  : { $$ = nil }
  | IDENTIFIER var_type { $$ = &vm.Argument{Name:$1, VarType:$2} }

statements
  : statement { s := new(ast.StateBlock); $$ = s.AddStates($1) }
  | statements statement { $$ = $1.AddStates($2) }

statement
  : EOL { $$ = nil }
  | assign EOL { $$ = $1 }
  | expr EOL { $$ = $1 }
  | define_var EOL { $$ = $1 }
  | return_statement EOL { $$ = $1 }

define_var
  : VAR IDENTIFIER var_type { $$ = &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line, Driver:driver} }}
  | VAR IDENTIFIER var_type ASSIGN expr { $$ = &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line,Right:$5, Driver:driver}}}
  | VAR IDENTIFIER ASSIGN expr { $$ =  &ast.NodeDecl{Name:$2, VarType:cm.TYPE_UNKNOWN, Node:ast.Node{Lineno:lexer.line,Right:$4, Driver:driver}}}

return_statement
  : RETURN      { $$ = ast.MakeNodeReturn(nil, lexer.line, driver) }
  | RETURN expr { $$ = ast.MakeNodeReturn($2,  lexer.line, driver) }

assign
  : IDENTIFIER ASSIGN expr
  { 
    varNode := ast.MakeNodeValue(lexer.line, $1, driver)
    $$ = ast.MakeNodeAssign(lexer.line, varNode, $3, driver)
  }

expr
  : const
  | MINUS expr %prec NEG{ $$ = ast.MakeExprNode(lexer.line, $2, nil, ast.OP_NOT, driver)}
  | IDENTIFIER          { $$ = ast.MakeNodeValue(lexer.line, $1, driver) }
  | function_call       { $$ = $1 }
  | expr INCR           { $$ = ast.MakeExprNode(lexer.line, $1, nil, ast.OP_INCR, driver) }
  | expr DECR           { $$ = ast.MakeExprNode(lexer.line, $1, nil, ast.OP_DECR, driver)}
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

function_call
  : IDENTIFIER '(' args ')' { $$ = ast.MakeNodeFunction(lexer.line, $1, $3, driver) }

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

%%

func Parse (filename string, source string) int {
  
  err := &cm.ErrorHandler{ErrorCount:0,WarningCount:0}
  driver = new(vm.Driver)
  driver.Init(filename, err)

  // パース処理
  lexer = &Lexer{src: source, position:0, readPosition:0, line:1, filename:filename, err:err}
	yyParse(lexer)

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
