%{
// プログラムのヘッダを指定
package compiler

import (
  "fmt"
)
%}
%union {
  ival int
  fval float32
  sval string
}
// 非終端記号の定義
%type<ival> program expr

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
%token<ival> IF ELSE SWITCH CASE DEFAULT FALLTHROWGH FOR BREAK CONTINUE FUNC RETURN IMPORT TYPE STRUCT 

%token<ival> EOL

// 演算の優先度の指定
%left PLUS, MINUS
%left ASTARISK, SLASH

%%
// 文法規則を指定
program
  : expr
  {
    $$ = $1
  }

expr
  : INUM
  | expr PLUS expr { $$ = $1 + $3 }
  | expr MINUS expr { $$ = $1 - $3 }
  | expr ASTARISK expr { $$ = $1 * $3 }
  | expr SLASH expr { $$ = $1 / $3 }

%%

func Parse (source string) int {
  lexer := &Lexer{src: source, position:0, readPosition:0, line:0}
	//yyParse(lexer)

  token := 0
  lval := new(yySymType)
  for token != -1{
    token = lexer.Lex(lval)
    fmt.Printf("token:%d, i:%d f:%g s:%s\n", token, lval.ival, lval.fval, lval.sval)
  }

  return 0
}
