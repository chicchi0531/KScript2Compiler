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
%type<ival> var_type
%type<node> statement expr assign const define_var

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
  :
  | program statement { if $2!=nil { $2.Push() } }

statement
  : EOL { $$ = nil; driver.lineno++ }
  | assign EOL { $$ = $1; driver.lineno++ }
  | expr EOL { $$ = $1; driver.lineno++ }
  | define_var EOL { $$ = $1; driver.lineno++ }

define_var
  : VAR IDENTIFIER var_type { $$ = nil; driver.variableTable.DefineInLocal($2, $3) }
  | VAR IDENTIFIER var_type ASSIGN expr { $$ = driver.variableTable.DefineInLocalWithAssign($2, $3, $5) }
  | VAR IDENTIFIER ASSIGN expr { $$ = driver.variableTable.DefineInLocalWithAssignAutoType($2, $4) }

assign
  : IDENTIFIER ASSIGN expr
  { 
    varNode := MakeValueNode($1, driver)
    $$ = &AssignNode{ Node:Node{ left:varNode, right:$3, driver:driver} }
  }

expr
  : const
  | IDENTIFIER { $$ = MakeValueNode($1, driver) }
  | MINUS expr %prec NEG { $$ = &Node{left:$2, right:nil, op:OP_NOT, driver:driver}}
  | expr INCR { $$ = &Node{left:$1, right:nil, op:OP_INCR, driver:driver}}
  | expr DECR { $$ = &Node{left:$1, right:nil, op:OP_DECR, driver:driver}}
  | expr PLUS expr { $$ = &Node{left:$1, right:$3, op:OP_ADD, driver:driver} }
  | expr MINUS expr { $$ = &Node{left:$1, right:$3, op:OP_SUB, driver:driver} }
  | expr ASTARISK expr { $$ = &Node{left:$1, right:$3, op:OP_MUL, driver:driver} }
  | expr SLASH expr { $$ = &Node{left:$1, right:$3, op:OP_DIV, driver:driver} }
  | expr PERCENT expr { $$ = &Node{left:$1, right:$3, op:OP_MOD, driver:driver} }
  | expr EQ expr { $$ = &Node{left:$1, right:$3, op:OP_EQUAL, driver:driver }}
  | expr NEQ expr { $$ = &Node{left:$1, right:$3, op:OP_NEQ, driver:driver }}
  | expr GT expr { $$ = &Node{left:$1, right:$3, op:OP_GT, driver:driver }}
  | expr GE expr { $$ = &Node{left:$1, right:$3, op:OP_GE, driver:driver }}
  | expr LT expr { $$ = &Node{left:$1, right:$3, op:OP_LT, driver:driver }}
  | expr LE expr { $$ = &Node{left:$1, right:$3, op:OP_LE, driver:driver }}
  | expr AND expr { $$ = &Node{left:$1, right:$3, op:OP_AND, driver:driver }}
  | expr OR expr { $$ = &Node{left:$1, right:$3, op:OP_OR, driver:driver }}
  | '(' expr ')' { $$ = $2 }

const
  : INUM { $$ = &Node{ op:OP_INTEGER, driver:driver, ival:$1 } }
  | FNUM { $$ = &Node{ op:OP_FLOAT, driver:driver, fval:$1 } }
  | STRING_LITERAL { $$ = &Node{ op:OP_STRING, driver:driver, sval:$1 } }

var_type
  : INT { $$ = $1 }
  | FLOAT { $$ = $1 }
  | STRING { $$ = $1 }

%%

func Parse (filename string, source string) int {

  driver = &Driver{
    pc:0, lineno:1, filename:filename,
    program:make([]Op,0),
    err:&ErrorHandler{errorCount:0,warningCount:0},
    variableTable:&VariableTable{currentTable:0}}
  driver.variableTable.driver = driver

  // パース処理
  lexer := &Lexer{src: source, position:0, readPosition:0, line:1, filename:filename, driver:driver}
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