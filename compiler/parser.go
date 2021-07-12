// Code generated by goyacc -o compiler/parser.go compiler/parser.go.y. DO NOT EDIT.

//line compiler/parser.go.y:1

// プログラムのヘッダを指定
package compiler

import __yyfmt__ "fmt"

//line compiler/parser.go.y:4

import (
	"fmt"
	"ks2/compiler/ast"
	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

var driver *vm.Driver
var lexer *Lexer

//line compiler/parser.go.y:17
type yySymType struct {
	yys  int
	ival int
	fval float32
	sval string

	node       vm.INode
	nodes      []vm.INode
	assign     *ast.Assign
	argList    []*vm.Argument
	argument   *vm.Argument
	stateBlock vm.IStateBlock
	statement  vm.IStatement
}

const INUM = 57346
const FNUM = 57347
const IDENTIFIER = 57348
const STRING_LITERAL = 57349
const PLUS = 57350
const MINUS = 57351
const ASTARISK = 57352
const SLASH = 57353
const PERCENT = 57354
const P_EQ = 57355
const M_EQ = 57356
const A_EQ = 57357
const S_EQ = 57358
const MOD_EQ = 57359
const EQ = 57360
const NEQ = 57361
const GT = 57362
const GE = 57363
const LT = 57364
const LE = 57365
const AND = 57366
const OR = 57367
const INCR = 57368
const DECR = 57369
const ASSIGN = 57370
const VAR = 57371
const INT = 57372
const FLOAT = 57373
const STRING = 57374
const VOID = 57375
const IF = 57376
const ELSE = 57377
const SWITCH = 57378
const CASE = 57379
const DEFAULT = 57380
const FALLTHROUGH = 57381
const FOR = 57382
const BREAK = 57383
const CONTINUE = 57384
const FUNC = 57385
const RETURN = 57386
const IMPORT = 57387
const TYPE = 57388
const STRUCT = 57389
const SYSCALL = 57390
const EOL = 57391
const NEG = 57392

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"INUM",
	"FNUM",
	"IDENTIFIER",
	"STRING_LITERAL",
	"PLUS",
	"MINUS",
	"ASTARISK",
	"SLASH",
	"PERCENT",
	"P_EQ",
	"M_EQ",
	"A_EQ",
	"S_EQ",
	"MOD_EQ",
	"EQ",
	"NEQ",
	"GT",
	"GE",
	"LT",
	"LE",
	"AND",
	"OR",
	"INCR",
	"DECR",
	"ASSIGN",
	"VAR",
	"INT",
	"FLOAT",
	"STRING",
	"VOID",
	"IF",
	"ELSE",
	"SWITCH",
	"CASE",
	"DEFAULT",
	"FALLTHROUGH",
	"FOR",
	"BREAK",
	"CONTINUE",
	"FUNC",
	"RETURN",
	"IMPORT",
	"TYPE",
	"STRUCT",
	"SYSCALL",
	"EOL",
	"'('",
	"')'",
	"'{'",
	"'}'",
	"'['",
	"']'",
	"','",
	"':'",
	"';'",
	"'.'",
	"NEG",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line compiler/parser.go.y:244

func Parse(filename string, source string) int {

	err := &cm.ErrorHandler{ErrorCount: 0, WarningCount: 0}
	driver = new(vm.Driver)
	driver.Init(filename, err)

	// パース処理
	lexer = &Lexer{src: source, position: 0, readPosition: 0, line: 1, filename: filename, err: err}
	yyParse(lexer)

	// ラベル設定
	driver.LabelSettings()

	fmt.Println("Parse End.")

	// パース結果出力
	driver.Dump()

	return 0
}

// 外部用
func GetErrorCount() int {
	return driver.Err.ErrorCount
}
func GetWarningCount() int {
	return driver.Err.WarningCount
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 71,
	18, 0,
	19, 0,
	20, 0,
	21, 0,
	22, 0,
	23, 0,
	-2, 61,
	-1, 72,
	18, 0,
	19, 0,
	20, 0,
	21, 0,
	22, 0,
	23, 0,
	-2, 62,
	-1, 73,
	18, 0,
	19, 0,
	20, 0,
	21, 0,
	22, 0,
	23, 0,
	-2, 63,
	-1, 74,
	18, 0,
	19, 0,
	20, 0,
	21, 0,
	22, 0,
	23, 0,
	-2, 64,
	-1, 75,
	18, 0,
	19, 0,
	20, 0,
	21, 0,
	22, 0,
	23, 0,
	-2, 65,
	-1, 76,
	18, 0,
	19, 0,
	20, 0,
	21, 0,
	22, 0,
	23, 0,
	-2, 66,
	-1, 93,
	49, 26,
	52, 26,
	-2, 54,
	-1, 97,
	49, 31,
	-2, 55,
}

const yyPrivate = 57344

const yyLast = 348

var yyAct = [...]int{
	100, 64, 32, 91, 87, 86, 17, 79, 94, 33,
	84, 132, 142, 103, 25, 98, 95, 104, 104, 26,
	63, 65, 3, 130, 61, 28, 16, 11, 12, 13,
	27, 60, 7, 65, 138, 62, 36, 37, 101, 35,
	136, 30, 18, 19, 20, 9, 66, 67, 68, 69,
	70, 71, 72, 73, 74, 75, 76, 77, 78, 8,
	22, 102, 80, 140, 82, 7, 98, 128, 93, 18,
	19, 20, 99, 58, 59, 97, 96, 24, 23, 15,
	38, 7, 34, 14, 6, 106, 93, 44, 85, 40,
	41, 42, 43, 97, 107, 5, 2, 115, 10, 116,
	118, 4, 1, 21, 92, 129, 85, 90, 119, 108,
	109, 110, 111, 112, 113, 114, 120, 89, 131, 88,
	133, 83, 134, 147, 36, 37, 101, 35, 117, 30,
	29, 80, 121, 139, 39, 135, 0, 141, 137, 0,
	143, 146, 144, 45, 46, 47, 48, 49, 93, 150,
	148, 149, 0, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 45, 46, 47, 48, 49, 38, 0,
	34, 0, 0, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 0, 45, 46, 47, 48, 49, 0,
	0, 0, 0, 145, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 45, 46, 47, 48, 49, 0,
	105, 0, 0, 0, 50, 51, 52, 53, 54, 55,
	56, 57, 58, 59, 36, 37, 101, 35, 65, 30,
	36, 37, 101, 35, 0, 30, 36, 37, 31, 35,
	0, 30, 123, 124, 125, 126, 127, 81, 0, 102,
	47, 48, 49, 0, 98, 102, 0, 122, 0, 0,
	99, 0, 0, 0, 96, 0, 58, 59, 38, 7,
	34, 0, 0, 0, 38, 0, 34, 0, 0, 61,
	38, 0, 34, 45, 46, 47, 48, 49, 0, 0,
	0, 0, 0, 50, 51, 52, 53, 54, 55, 56,
	57, 58, 59, 45, 46, 47, 48, 49, 0, 0,
	0, 0, 0, 50, 51, 52, 53, 54, 55, 56,
	0, 58, 59, 45, 46, 47, 48, 49, 45, 46,
	47, 48, 49, 50, 51, 52, 53, 54, 55, 0,
	0, 58, 59, 0, 0, 0, 58, 59,
}

var yyPact = [...]int{
	16, 16, -1000, -1000, -17, -17, -17, -1000, 77, 73,
	-1000, -1000, -1000, -1000, -24, 39, 72, 49, -1000, -1000,
	-1000, -37, -1000, 39, 232, 59, 72, -1000, 275, -1000,
	232, -26, -1000, -1000, 232, -1000, -1000, -1000, -34, -31,
	-1000, -1000, -1000, -1000, -1000, 232, 232, 232, 232, 232,
	232, 232, 232, 232, 232, 232, 232, 232, -1000, -1000,
	-1000, 232, 196, 232, -1000, 220, 240, 240, 47, 47,
	47, 320, 320, 320, 320, 320, 320, 315, 295, -38,
	275, -1000, 155, 32, -1000, -1000, -17, -17, -17, -17,
	-17, -17, -17, -1000, -1000, -1000, 232, -1000, 232, 226,
	275, 229, 61, -1000, 232, -27, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 275, 176, -47, 176, -1000,
	-1000, 232, -1000, -1000, -1000, -1000, -1000, -1000, 12, 275,
	232, -1, 232, -1000, 275, 35, 232, -39, -19, 135,
	232, 275, -1000, -1000, -1000, 120, 275, -31, -1000, -1000,
	-1000,
}

var yyPgo = [...]int{
	0, 6, 134, 132, 0, 130, 16, 9, 2, 8,
	128, 123, 121, 10, 5, 4, 119, 117, 107, 1,
	3, 104, 60, 103, 7, 102, 96, 22, 101, 95,
	84,
}

var yyR1 = [...]int{
	0, 25, 25, 26, 26, 26, 26, 29, 29, 30,
	28, 23, 23, 23, 22, 19, 12, 12, 13, 13,
	13, 13, 13, 13, 13, 13, 14, 15, 16, 17,
	17, 18, 20, 20, 20, 10, 10, 11, 11, 21,
	21, 3, 3, 3, 3, 3, 3, 9, 6, 6,
	6, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	8, 8, 7, 7, 24, 24, 24, 5, 5, 5,
	1, 1, 1, 2, 2, 2, 2, 2, 27,
}

var yyR2 = [...]int{
	0, 1, 2, 1, 2, 2, 2, 3, 5, 6,
	7, 0, 1, 3, 2, 3, 1, 2, 1, 2,
	2, 2, 2, 2, 2, 2, 1, 1, 1, 1,
	2, 1, 3, 5, 5, 1, 1, 1, 1, 7,
	3, 1, 1, 1, 1, 1, 1, 3, 3, 5,
	4, 1, 2, 1, 1, 1, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 4, 7, 0, 1, 3, 1, 1, 1,
	1, 1, 1, 0, 1, 1, 1, 1, 1,
}

var yyChk = [...]int{
	-1000, -25, -26, -27, -28, -29, -30, 49, 43, 29,
	-26, -27, -27, -27, 6, 6, 50, -1, 30, 31,
	32, -23, -22, 6, 28, 51, 56, -1, -4, -5,
	9, 6, -8, -7, 50, 7, 4, 5, 48, -2,
	30, 31, 32, 33, -22, 8, 9, 10, 11, 12,
	18, 19, 20, 21, 22, 23, 24, 25, 26, 27,
	-4, 50, -4, 54, -19, 52, -4, -4, -4, -4,
	-4, -4, -4, -4, -4, -4, -4, -4, -4, -24,
	-4, 51, -4, -12, -13, -27, -14, -15, -16, -17,
	-18, -20, -21, -8, -9, -6, 44, -7, 34, 40,
	-4, 6, 29, 51, 56, 55, 53, -13, -27, -27,
	-27, -27, -27, -27, -27, -4, -4, -10, -4, -9,
	-6, -3, 28, 13, 14, 15, 16, 17, 6, -4,
	50, -19, 58, -19, -4, -1, 28, -24, 35, -4,
	28, -4, 51, -19, -20, 58, -4, -11, -15, -14,
	-19,
}

var yyDef = [...]int{
	0, -2, 1, 3, 0, 0, 0, 88, 0, 0,
	2, 4, 5, 6, 0, 0, 11, 7, 80, 81,
	82, 0, 12, 0, 0, 83, 0, 14, 8, 51,
	0, 53, 54, 55, 0, 77, 78, 79, 0, 9,
	84, 85, 86, 87, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 70, 71,
	52, 74, 0, 0, 10, 0, 56, 57, 58, 59,
	60, -2, -2, -2, -2, -2, -2, 67, 68, 0,
	75, 69, 0, 0, 16, 18, 0, 0, 0, 0,
	0, 0, 0, -2, 27, 28, 29, -2, 0, 0,
	0, 53, 0, 72, 0, 0, 15, 17, 19, 20,
	21, 22, 23, 24, 25, 30, 0, 0, 0, 35,
	36, 0, 41, 42, 43, 44, 45, 46, 0, 76,
	74, 32, 0, 40, 47, 48, 0, 0, 0, 0,
	0, 50, 73, 33, 34, 0, 49, 0, 37, 38,
	39,
}

var yyTok1 = [...]int{
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	50, 51, 3, 3, 56, 3, 59, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 57, 58,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 54, 3, 55, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 52, 3, 53,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 60,
}

var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:96
		{
			driver.VariableTable.DefineInLocal(lexer.line, yyDollar[2].sval, yyDollar[3].ival)
		}
	case 8:
		yyDollar = yyS[yypt-5 : yypt+1]
//line compiler/parser.go.y:97
		{
			driver.Err.LogError(driver.Filename, lexer.line, cm.ERR_0026, "")
		}
	case 9:
		yyDollar = yyS[yypt-6 : yypt+1]
//line compiler/parser.go.y:100
		{
			driver.DecralateFunction(lexer.line, yyDollar[6].ival, yyDollar[2].sval, yyDollar[4].argList)
		}
	case 10:
		yyDollar = yyS[yypt-7 : yypt+1]
//line compiler/parser.go.y:103
		{
			driver.AddFunction(lexer.line, yyDollar[6].ival, yyDollar[2].sval, yyDollar[4].argList, yyDollar[7].statement)
		}
	case 11:
		yyDollar = yyS[yypt-0 : yypt+1]
//line compiler/parser.go.y:106
		{
			yyVAL.argList = make([]*vm.Argument, 0)
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:107
		{
			yyVAL.argList = []*vm.Argument{yyDollar[1].argument}
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:108
		{
			yyVAL.argList = append(yyDollar[1].argList, yyDollar[3].argument)
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:111
		{
			yyVAL.argument = &vm.Argument{Name: yyDollar[1].sval, VarType: yyDollar[2].ival}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:117
		{
			yyVAL.statement = ast.MakeCompoundStatement(yyDollar[2].stateBlock)
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:120
		{
			s := new(ast.StateBlock)
			yyVAL.stateBlock = s.AddStates(yyDollar[1].statement)
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:121
		{
			yyVAL.stateBlock = yyDollar[1].stateBlock.AddStates(yyDollar[2].statement)
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:124
		{
			yyVAL.statement = nil
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:125
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:126
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:127
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:128
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:129
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:130
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:131
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:134
		{
			yyVAL.statement = ast.MakeExprStatement(yyDollar[1].node, driver)
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:137
		{
			yyVAL.statement = ast.MakeAssignStatement(yyDollar[1].node)
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:140
		{
			yyVAL.statement = ast.MakeVarDefineStatement(yyDollar[1].node)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:143
		{
			yyVAL.statement = ast.MakeReturnStatement(nil, lexer.line, driver)
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:144
		{
			yyVAL.statement = ast.MakeReturnStatement(yyDollar[2].node, lexer.line, driver)
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:147
		{
			yyVAL.statement = ast.MakeFunctionCallStatement(yyDollar[1].node, driver)
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:150
		{
			yyVAL.statement = ast.MakeIfStatement(yyDollar[2].node, yyDollar[3].statement, nil, lexer.line, driver)
		}
	case 33:
		yyDollar = yyS[yypt-5 : yypt+1]
//line compiler/parser.go.y:151
		{
			yyVAL.statement = ast.MakeIfStatement(yyDollar[2].node, yyDollar[3].statement, yyDollar[5].statement, lexer.line, driver)
		}
	case 34:
		yyDollar = yyS[yypt-5 : yypt+1]
//line compiler/parser.go.y:152
		{
			yyVAL.statement = ast.MakeIfStatement(yyDollar[2].node, yyDollar[3].statement, yyDollar[5].statement, lexer.line, driver)
		}
	case 39:
		yyDollar = yyS[yypt-7 : yypt+1]
//line compiler/parser.go.y:163
		{
			yyVAL.statement = ast.MakeForStatement(yyDollar[2].node, yyDollar[4].node, yyDollar[6].statement, yyDollar[7].statement, lexer.line, driver)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:164
		{
			yyVAL.statement = ast.MakeWhileStatement(yyDollar[2].node, yyDollar[3].statement, lexer.line, driver)
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:171
		{
			yyVAL.ival = ast.OP_ASSIGN
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:172
		{
			yyVAL.ival = ast.OP_ADD_ASSIGN
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:173
		{
			yyVAL.ival = ast.OP_SUB_ASSIGN
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:174
		{
			yyVAL.ival = ast.OP_MUL_ASSIGN
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:175
		{
			yyVAL.ival = ast.OP_DIV_ASSIGN
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:176
		{
			yyVAL.ival = ast.OP_MOD_ASSIGN
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:180
		{
			varNode := ast.MakeValueNode(lexer.line, yyDollar[1].sval, driver)
			yyVAL.node = ast.MakeAssign(lexer.line, varNode, yyDollar[3].node, yyDollar[2].ival, driver)
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:186
		{
			yyVAL.node = ast.MakeVarDefineNode(lexer.line, yyDollar[2].sval, yyDollar[3].ival, driver)
		}
	case 49:
		yyDollar = yyS[yypt-5 : yypt+1]
//line compiler/parser.go.y:187
		{
			yyVAL.node = ast.MakeVarDefineNodeWithAssign(lexer.line, yyDollar[2].sval, yyDollar[3].ival, yyDollar[5].node, driver)
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
//line compiler/parser.go.y:188
		{
			yyVAL.node = ast.MakeVarDefineNodeWithAssign(lexer.line, yyDollar[2].sval, cm.TYPE_UNKNOWN, yyDollar[4].node, driver)
		}
	case 52:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:192
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[2].node, nil, ast.OP_NOT, driver)
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:193
		{
			yyVAL.node = ast.MakeValueNode(lexer.line, yyDollar[1].sval, driver)
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:194
		{
			yyVAL.node = yyDollar[1].node
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:195
		{
			yyVAL.node = yyDollar[1].node
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:196
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_ADD, driver)
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:197
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_SUB, driver)
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:198
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_MUL, driver)
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:199
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_DIV, driver)
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:200
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_MOD, driver)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:201
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_EQUAL, driver)
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:202
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_NEQ, driver)
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:203
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_GT, driver)
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:204
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_GE, driver)
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:205
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_LT, driver)
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:206
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_LE, driver)
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:207
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_AND, driver)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:208
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, yyDollar[3].node, ast.OP_OR, driver)
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:209
		{
			yyVAL.node = yyDollar[2].node
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:212
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, nil, ast.OP_INCR, driver)
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
//line compiler/parser.go.y:213
		{
			yyVAL.node = ast.MakeExprNode(lexer.line, yyDollar[1].node, nil, ast.OP_DECR, driver)
		}
	case 72:
		yyDollar = yyS[yypt-4 : yypt+1]
//line compiler/parser.go.y:216
		{
			yyVAL.node = ast.MakeFunctionNode(lexer.line, yyDollar[1].sval, yyDollar[3].nodes, driver)
		}
	case 73:
		yyDollar = yyS[yypt-7 : yypt+1]
//line compiler/parser.go.y:217
		{
			yyVAL.node = ast.MakeSysCallNode(lexer.line, yyDollar[3].node, yyDollar[6].nodes, driver)
		}
	case 74:
		yyDollar = yyS[yypt-0 : yypt+1]
//line compiler/parser.go.y:220
		{
			yyVAL.nodes = make([]vm.INode, 0)
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:221
		{
			yyVAL.nodes = []vm.INode{yyDollar[1].node}
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
//line compiler/parser.go.y:222
		{
			yyVAL.nodes = append(yyDollar[1].nodes, yyDollar[3].node)
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:225
		{
			yyVAL.node = ast.MakeSvalNode(lexer.line, yyDollar[1].sval, driver)
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:226
		{
			yyVAL.node = ast.MakeIvalNode(lexer.line, yyDollar[1].ival, driver)
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:227
		{
			yyVAL.node = ast.MakeFvalNode(lexer.line, yyDollar[1].fval, driver)
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:230
		{
			yyVAL.ival = cm.TYPE_INTEGER
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:231
		{
			yyVAL.ival = cm.TYPE_FLOAT
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:232
		{
			yyVAL.ival = cm.TYPE_STRING
		}
	case 83:
		yyDollar = yyS[yypt-0 : yypt+1]
//line compiler/parser.go.y:235
		{
			yyVAL.ival = cm.TYPE_VOID
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:236
		{
			yyVAL.ival = cm.TYPE_INTEGER
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:237
		{
			yyVAL.ival = cm.TYPE_FLOAT
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:238
		{
			yyVAL.ival = cm.TYPE_STRING
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
//line compiler/parser.go.y:239
		{
			yyVAL.ival = cm.TYPE_VOID
		}
	}
	goto yystack /* stack new state and value */
}
