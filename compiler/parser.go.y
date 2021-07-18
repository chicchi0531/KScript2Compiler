%{
// プログラムのヘッダを指定
package compiler

import(
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
  value *vm.VariableTag
  valueList []*vm.VariableTag

  valueNode *ast.NValue

  stateBlock vm.IStateBlock
  statement vm.IStatement
  caseStatement *ast.CaseStatement
  caseStatements []*ast.CaseStatement
}

// 非終端記号の定義
%type<ival> var_type
%type<ival> function_type
%type<ival> assign_op

%type<node> expr const define_var function_call uni_expr assign
%type<valueNode> value
%type<node> for_init
%type<statement> for_iterator
%type<caseStatement> case_statement
%type<caseStatements> cases
%type<statement> default_statement

%type<stateBlock> statements
%type<statement> statement
%type<statement> expr_statement
%type<statement> assign_statement
%type<statement> define_var_statement
%type<statement> return_statement
%type<statement> function_call_statement
%type<statement> block
%type<statement> if_statement
%type<statement> for_statement
%type<statement> break_statement
%type<statement> continue_statement
%type<statement> switch_statement
%type<statement> fallthrough_statement
%type<statement> dump_statement

%type<argument> arg_decl
%type<argList> arg_list
%type<nodes> args

%type<valueList> member_list
%type<value> member

// 終端記号の定義
%token<ival> INUM
%token<fval> FNUM
%token<sval> IDENTIFIER STRING_LITERAL

%token<ival> PLUS MINUS ASTARISK SLASH PERCENT // + - * / %
%token<ival> P_EQ M_EQ A_EQ S_EQ MOD_EQ // += -= *= /= %=
%token<ival> EQ NEQ GT GE LT LE // == != > >= < <=
%token<ival> AND OR // && ||
%token<ival> INCR DECR ASSIGN DECL_ASSIGN// ++ -- = :=

%token<ival> VAR INT FLOAT STRING VOID
%token<ival> IF ELSE SWITCH CASE DEFAULT FALLTHROUGH FOR BREAK CONTINUE FUNC RETURN IMPORT TYPE STRUCT SYSCALL

%token<ival> DUMP // debug用

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
  | import
  | function_define eol
  | global_decl eol
  | function_decl eol
  | type_decl eol

import
  : IMPORT STRING_LITERAL { ImportFile($2) }

global_decl
  : VAR IDENTIFIER var_type              { driver.VariableTable.DefineValue(lexer.line, $2, $3, false, 1) }
  | VAR IDENTIFIER var_type ASSIGN expr  { driver.Err.LogError(driver.Filename, lexer.line, cm.ERR_0026, "") }
  | VAR IDENTIFIER '[' INUM ']' var_type { driver.VariableTable.DefineValue(lexer.line, $2, $6, false, $4) }

function_decl
  : FUNC IDENTIFIER '(' arg_list ')' function_type       { driver.DecralateFunction(lexer.line,$6,$2,$4) }

function_define
  : FUNC IDENTIFIER '(' arg_list ')' function_type block { driver.AddFunction(lexer.line,$6,$2,$4,$7) }

arg_list
  :                       { $$ = make([]*vm.Argument, 0) }
  | arg_decl              { $$ = []*vm.Argument{$1} }
  | arg_list ',' arg_decl { $$ = append($1,$3) }

arg_decl
  : IDENTIFIER var_type               { $$ = &vm.Argument{Name:$1, VarType:$2, IsPointer:false, Size:1} }
  | IDENTIFIER '[' INUM ']' var_type  { $$ = &vm.Argument{Name:$1, VarType:$5, IsPointer:false, Size:$3} }

type_decl
  : TYPE IDENTIFIER STRUCT '{' member_list '}' { driver.AddType($2,$5,lexer.line) }

member_list
  : member             {
    t := make([]*vm.VariableTag,0)
    if $1 != nil {
      t = append(t,$1)
    }
    $$ = t
  }
  | member_list member { if $2 != nil { $$ = append($1, $2) }else{ $$ = $1 } }

member
  : eol          { $$ = nil }
  | arg_decl eol { $$ = vm.MakeVariableTag($1.Name, $1.VarType, $1.IsPointer, $1.Size) }

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
  | for_statement eol { $$ = $1 }
  | break_statement eol { $$ = $1 }
  | continue_statement eol { $$ = $1 }
  | switch_statement eol { $$ = $1 }
  | fallthrough_statement eol { $$ = $1 }
  | dump_statement eol { $$ = $1 }

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

for_init
  : assign
  | define_var

for_iterator
  : assign_statement
  | expr_statement

for_statement
  : FOR for_init ';' expr ';' for_iterator block { $$ = ast.MakeForStatement($2, $4, $6, $7, lexer.line, driver) }
  | FOR expr block { $$ = ast.MakeWhileStatement($2, $3, lexer.line, driver) }

break_statement
  : BREAK { $$ = ast.MakeBreakStatement(lexer.line, driver) }

continue_statement
  : CONTINUE { $$ = ast.MakeContinueStatement(lexer.line, driver) }

switch_statement
  : SWITCH expr '{' eols cases '}'                    { $$ = ast.MakeSwitchStatement($2, $5, nil, lexer.line, driver) }
  | SWITCH expr '{' eols cases default_statement '}'  { $$ = ast.MakeSwitchStatement($2, $5, $6, lexer.line, driver) }

eols
  : eol
  | eols eol

cases
  : case_statement               { $$ = []*ast.CaseStatement{$1} }
  | cases case_statement         { $$ = append($1, $2) }

case_statement
  : CASE expr ':' statements     { $$ = ast.MakeCaseStatement($2, ast.MakeCompoundStatement($4), lexer.line, driver) }

default_statement
  : DEFAULT ':' statements  { $$ = ast.MakeCompoundStatement($3) }

fallthrough_statement
  : FALLTHROUGH { $$ = ast.MakeFallThroughStatement(lexer.line, driver) }

dump_statement
  : DUMP '(' STRING_LITERAL ')' { $$ = ast.MakeDumpStatement($3, lexer.line, driver) }
  | DUMP                        { $$ = ast.MakeDumpStatement("", lexer.line, driver) }

//------------------------------
// expr
//------------------------------

assign_op
  : ASSIGN { $$ = ast.OP_ASSIGN }
  | P_EQ   { $$ = ast.OP_ADD_ASSIGN }
  | M_EQ   { $$ = ast.OP_SUB_ASSIGN }
  | A_EQ   { $$ = ast.OP_MUL_ASSIGN }
  | S_EQ   { $$ = ast.OP_DIV_ASSIGN }
  | MOD_EQ { $$ = ast.OP_MOD_ASSIGN }

assign
  : value assign_op expr { $$ = ast.MakeAssign(lexer.line, $1, $3, $2, driver) }

define_var
  : VAR IDENTIFIER var_type               { $$ = ast.MakeVarDefineNode(lexer.line, $2, $3, false, 1, driver) }
  | VAR IDENTIFIER var_type ASSIGN expr   { $$ = ast.MakeVarDefineNodeWithAssign(lexer.line, $2, $3, $5, driver) }
  | VAR IDENTIFIER ASSIGN expr            { $$ = ast.MakeVarDefineNodeWithAssign(lexer.line, $2, cm.TYPE_UNKNOWN, $4, driver) }
  | IDENTIFIER DECL_ASSIGN expr           { $$ = ast.MakeVarDefineNodeWithAssign(lexer.line, $1, cm.TYPE_UNKNOWN, $3, driver) }
  | VAR IDENTIFIER '[' INUM ']' var_type  { $$ = ast.MakeVarDefineNode(lexer.line, $2, $6, false, $4, driver) }

expr
  : const
  | MINUS expr %prec NEG{ $$ = ast.MakeExprNode(lexer.line, $2, nil, ast.OP_NOT, driver)}
  | value               { $$ = $1 }
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
  : IDENTIFIER '(' args ')'           { $$ = ast.MakeFunctionNode(lexer.line, $1, $3, driver) }
  | SYSCALL '[' expr ']' '(' args ')' { $$ = ast.MakeSysCallNode(lexer.line, $3, $6, driver) }

value
  : IDENTIFIER              { $$ = ast.MakeValueNode(lexer.line, $1, driver) }
  | IDENTIFIER '[' expr ']' { $$ = ast.MakeArrayValueNode(lexer.line, $1, $3, driver) }
  | IDENTIFIER '.' value    { $$ = ast.MakeMemberValueNode(lexer.line, $1, $3, driver) }
  | IDENTIFIER '[' expr ']' '.' value { $$ = ast.MakeArrayMemberValueNode(lexer.line, $1, $3, $6, driver) }

args
  : { $$ = make([]vm.INode,0) }
  | expr { $$ = []vm.INode{$1} }
  | args ',' expr { $$ = append($1,$3) }

const
  : STRING_LITERAL  { $$ = ast.MakeSvalNode(lexer.line, $1, driver) }
  | INUM            { $$ = ast.MakeIvalNode(lexer.line, $1, driver) }
  | FNUM            { $$ = ast.MakeFvalNode(lexer.line, $1, driver) }

var_type
  : INT         { $$ = cm.TYPE_INTEGER }
  | FLOAT       { $$ = cm.TYPE_FLOAT }
  | STRING      { $$ = cm.TYPE_STRING }
  | IDENTIFIER  { $$ = driver.GetType($1, lexer.line) }

function_type
  :         { $$ = cm.TYPE_VOID }
  | INT     { $$ = cm.TYPE_INTEGER }
  | FLOAT   { $$ = cm.TYPE_FLOAT }
  | STRING  { $$ = cm.TYPE_STRING }
  | VOID    { $$ = cm.TYPE_VOID }
  | IDENTIFIER { $$ = driver.GetType($1, lexer.line) }

eol
  : EOL

%%
