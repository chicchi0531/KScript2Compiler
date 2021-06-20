%{
// プログラムのヘッダを指定
package compiler

import (
  "os"
  "fmt"
)
%}
%union {
  ival int
  dval double
  sval string
}
// 非終端記号の定義
%type<ival> program expr

// 終端記号の定義
%token<ival> NUMBER
%token<dval> DOUBLE
%token<sval> IDENTIFIER STRINGLITERAL

// 演算の優先度の指定
%left '+','-'
%left '*','/'

%%
// 文法規則を指定
program
  : expr
  {
    $$ = $1
    yylex.(*Lexer).result = $$
  }

expr
  : NUMBER
  | expr '+' expr { $$ = $1 + $3 }
  | expr '-' expr { $$ = $1 - $3 }
  | expr '*' expr { $$ = $1 * $3 }
  | expr '/' expr { $$ = $1 / $3 }

%%

// 最低限必要な構造体を定義
type Lexer struct {
  src    string
  index  int
  result int
}

// ここでトークン（最小限の要素）を一つずつ返す
func (p *Lexer) Lex(lval *yySymType) int {
  for p.index < len(p.src) {
    c := p.src[p.index]
    p.index++
    if c == '+' { return int(c) }
    if c == '-' { return int(c) }
    if c == '*' { return int(c) }
    if c == '/' { return int(c) }
    if '0' <= c && c <= '9' {
      lval.ival = int(c - '0')
      return NUMBER
    }
  }
  return -1
}
// エラー報告用
func (p *Lexer) Error(e string) {
  fmt.Println("[error] " + e)
}
