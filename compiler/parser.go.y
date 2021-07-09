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
  argList *vm.ArgList
  argument *vm.Argument
  stateBlock vm.IStateBlock
}

// 非終端記号の定義
%type<ival> var_type function_type
%type<node> statement expr assign const define_var uni_expr
%type<stateBlock> statements
%type<argument> arg
%type<argList> arg_list

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
  | function
  | global_decl

global_decl
  : VAR IDENTIFIER var_type { n := &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line, Driver:driver} }; n.Push() }
  | VAR IDENTIFIER var_type ASSIGN uni_expr { n := &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line, Right:$5, Driver:driver} }; n.Push()}
  | VAR IDENTIFIER ASSIGN uni_expr { n :=  &ast.NodeDecl{Name:$2, VarType:cm.TYPE_UNKNOWN, Node:ast.Node{Lineno:lexer.line,Right:$4, Driver:driver}}; n.Push()}

function
  : FUNC IDENTIFIER '(' arg_list ')' function_type '{' statements '}' { driver.AddFunction($6,$2,$4,$8) }

arg_list
  : arg { l := new(vm.ArgList); $$ = l.Add($1) }
  | arg_list ',' arg { $$ = $1.Add($3) }

arg
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

define_var
  : VAR IDENTIFIER var_type { $$ = &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line, Driver:driver} }}
  | VAR IDENTIFIER var_type ASSIGN expr { $$ = &ast.NodeDecl{Name:$2, VarType:$3, Node:ast.Node{Lineno:lexer.line,Right:$5, Driver:driver}}}
  | VAR IDENTIFIER ASSIGN expr { $$ =  &ast.NodeDecl{Name:$2, VarType:cm.TYPE_UNKNOWN, Node:ast.Node{Lineno:lexer.line,Right:$4, Driver:driver}}}

assign
  : IDENTIFIER ASSIGN expr
  { 
    varNode := &ast.NodeValue{Node:ast.Node{Lineno:lexer.line, Driver:driver }, Name:$1 }
    $$ = &ast.NodeAssign{ Node:ast.Node{Lineno:lexer.line, Left:varNode, Right:$3, Driver:driver} }
  }

expr
  : const
  | MINUS expr %prec NEG { $$ = &ast.Node{Lineno:lexer.line,Left:$2, Right:nil, Op:ast.OP_NOT, Driver:driver}}
  | IDENTIFIER { $$ = &ast.NodeValue{Node:ast.Node{Lineno:lexer.line, Driver:driver }, Name:$1 } }
  | expr INCR { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:nil, Op:ast.OP_INCR, Driver:driver}}
  | expr DECR { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:nil, Op:ast.OP_DECR, Driver:driver}}
  | expr PLUS expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_ADD, Driver:driver} }
  | expr MINUS expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_SUB, Driver:driver} }
  | expr ASTARISK expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_MUL, Driver:driver} }
  | expr SLASH expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_DIV, Driver:driver} }
  | expr PERCENT expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_MOD, Driver:driver} }
  | expr EQ expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_EQUAL, Driver:driver }}
  | expr NEQ expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_NEQ, Driver:driver }}
  | expr GT expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_GT, Driver:driver }}
  | expr GE expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_GE, Driver:driver }}
  | expr LT expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_LT, Driver:driver }}
  | expr LE expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_LE, Driver:driver }}
  | expr AND expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_AND, Driver:driver }}
  | expr OR expr { $$ = &ast.Node{Lineno:lexer.line,Left:$1, Right:$3, Op:ast.OP_OR, Driver:driver }}
  | '(' expr ')' { $$ = $2 }

uni_expr
  : const
  | MINUS INUM { $$ = &ast.Node{Lineno:lexer.line, Op:ast.OP_INTEGER, Driver:driver, Ival:-$2 } }
  | MINUS FNUM { $$ = &ast.Node{Lineno:lexer.line, Op:ast.OP_INTEGER, Driver:driver, Fval:-$2 }}

const
  : STRING_LITERAL { $$ = &ast.Node{Lineno:lexer.line, Op:ast.OP_STRING, Driver:driver, Sval:$1 } }
  | INUM { $$ = &ast.Node{Lineno:lexer.line, Op:ast.OP_INTEGER, Driver:driver, Ival:$1 } }
  | FNUM { $$ = &ast.Node{Lineno:lexer.line, Op:ast.OP_FLOAT, Driver:driver, Fval:$1 } }

var_type
  : INT { $$ = cm.TYPE_INTEGER }
  | FLOAT { $$ = cm.TYPE_FLOAT }
  | STRING { $$ = cm.TYPE_STRING }

function_type
  : { $$ = cm.TYPE_VOID }
  | INT { $$ = cm.TYPE_INTEGER }
  | FLOAT { $$ = cm.TYPE_FLOAT }
  | STRING { $$ = cm.TYPE_STRING }
  | VOID { $$ = cm.TYPE_VOID }

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
