%{
// プログラムのヘッダを指定
package compiler

import(
  "fmt"
)

var driver *Driver

%}

%union {
  ival int
  fval float32
  sval string

  node INode
}
// 非終端記号の定義
%type<ival> define_or_state
%type<node> expr assign const

// 終端記号の定義
%token<ival> INUM
%token<fval> FNUM
%token<sval> IDENTIFIER STRING_LITERAL

%token<ival> PLUS MINUS ASTARISK SLASH // + - * /
%token<ival> P_EQ M_EQ A_EQ S_EQ // += -= *= /=
%token<ival> EQ NEQ GT GE LT LE // == != > >= < <=
%token<ival> AND OR // && ||
%token<ival> INCR DECR ASSIGN// ++ -- =

%token<ival> VAR INT FLOAT STRING VOID
%token<ival> IF ELSE SWITCH CASE DEFAULT FALLTHROUGH FOR BREAK CONTINUE FUNC RETURN IMPORT TYPE STRUCT 

%token<ival> EOL

// 演算の優先度の指定
%left OR
%left AND
%left EQ, NEQ, GT, GE, LT, LE
%left PLUS, MINUS
%left ASTARISK, SLASH
%left INCR, DECR

%%
// 文法規則を指定
program
  :
  | program define_or_state

define_or_state
  : EOL { fmt.Println("eol") }
  | assign EOL { $$ = $1.Push() }

assign
  : IDENTIFIER ASSIGN expr
  { 
    varNode := &ValueNode{ Node:Node{ driver:driver }}
    $$ = &AssignNode{ Node:Node{ left:varNode, right:$3, driver:driver} }
  }

expr
  : const
  | IDENTIFIER { $$ = &ValueNode{ Node:Node{driver:driver }}}
  | IDENTIFIER INCR {}
  | IDENTIFIER DECR {}
  | expr PLUS expr { $$ = &Node{left:$1, right:$3, op:OP_ADD, driver:driver} }
  | expr MINUS expr { $$ = &Node{left:$1, right:$3, op:OP_SUB, driver:driver} }
  | expr ASTARISK expr { $$ = &Node{left:$1, right:$3, op:OP_MUL, driver:driver} }
  | expr SLASH expr { $$ = &Node{left:$1, right:$3, op:OP_DIV, driver:driver} }
  | expr EQ expr {}
  | expr NEQ expr {}
  | expr GT expr {}
  | expr GE expr {}
  | expr LT expr {}
  | expr LE expr {}
  | expr AND expr {}
  | expr OR expr {}

const
  : INUM { $$ = &Node{ op:OP_INTEGER, driver:driver, ival:$1 } }
  | FNUM { $$ = &Node{ op:OP_FLOAT, driver:driver, fval:$1 } }
  | STRING_LITERAL { $$ = &Node{ op:OP_STRING, driver:driver, sval:$1 } }

%%

func Parse (filename string, source string) int {
  driver = &Driver{pc:0, lineno:0, filename:filename, err:&ErrorHandler{errorCount:0,warningCount:0}}

  // パース処理
  lexer := &Lexer{src: source, position:0, readPosition:0, line:0, filename:filename, driver:driver}
	yyParse(lexer)

  fmt.Println("Parse End.")

  return 0
}

// 外部用
func GetErrorCount() int{
  return driver.err.errorCount
}
func GetWarningCount() int{
  return driver.err.warningCount
}