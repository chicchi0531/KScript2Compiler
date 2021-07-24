package compiler

import(
	cm "ks2/compiler/common"
	"strconv"
)

// 最低限必要な構造体を定義
type Lexer struct {
	BaseLexer
}

func MakeLexer(filename string, src string, err *cm.ErrorHandler) *Lexer{
	l := new(Lexer)
	l.filename = filename
	l.src = src
	l.position = 0
	l.readPosition = 0
	l.line = 1
	l.err = err
	return l
}

// ここでトークン（最小限の要素）を一つずつ返す
func (p *Lexer) Lex(lval *yySymType) int {
	tok := 0

	if p.position >= len(p.src){
		return -1
	}

	// コメントスキップ
	p.readChar()
	p.skipComments()
	p.skipWhiteSpace()

	switch(p.ch){
	case 0:
		tok = -1
	case '"':
		lval.sval = p.readStringLiteral()
		tok = STRING_LITERAL
	case '+':
		ch := p.nextChar()
		if ch == '+'{
			p.readChar()
			tok = INCR //++
		}else if ch == '='{
			p.readChar()
			tok = P_EQ //+=
		}else{
			tok = PLUS // +
		}

	case '-':
		ch := p.nextChar()
		if ch == '-'{
			p.readChar()
			tok = DECR // --
		}else if ch == '=' {
			p.readChar()
			tok = M_EQ // -=
		}else{
			tok = MINUS // -
		}

	case '*':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = A_EQ // *=
		}else{
			tok = ASTARISK
		}

	case '/':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = S_EQ // /=
		}else{
			tok = SLASH
		}
		
	case '%':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = MOD_EQ //%=
		}else{
			tok = PERCENT
		}

	case '=':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = EQ
		}else{
			tok = ASSIGN
		}

	case '!':
		p.readChar()
		if p.ch == '='{
			tok = NEQ
		}else{
			p._err(cm.ERR_0009, "ErrorToken: !")
		}
	case '>':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = GE
		}else{
			tok = GT
		}

	case '<':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = LE
		}else{
			tok = LT
		}

	case '&':
		p.readChar()
		if p.ch == '&'{
			tok = AND
		}else{
			p._err(cm.ERR_0009, "ErrorToken: &")
		}

	case '|':
		p.readChar()
		if p.ch == '|'{
			tok = OR
		}else{
			p._err(cm.ERR_0009, "ErrorToken: |")
		}
	
	case ':':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = DECL_ASSIGN
		}else{
			tok = int(p.ch)
		}

	case '\n':
		tok = EOL

	case '(': fallthrough
	case ')': fallthrough
	case '{': fallthrough
	case '}': fallthrough
	case '[': fallthrough
	case ']': fallthrough
	case ';': fallthrough
	case ',': fallthrough
	case '.': tok = int(p.ch)

	default:
		if isNumber(p.ch){
			num := p.readNumber()

			// 少数点数の場合
			if p.nextChar() == '.'{
				p.readChar()//.の読み飛ばし
				
				p.readChar()
				lval.fval = float32(num) + p.readRealNumber()

				// floatのsuffixを読み飛ばす
				if p.ch == 'f' || p.ch == 'F'{
					p.readChar()
				}
				tok = FNUM
			}else{
				// 整数の場合
				lval.ival = num
				tok = INUM
			}
		}else if isLetter(p.ch){
			key := p.readIdentifier()
			tok = getKeywordToken(key)
			lval.sval = key
		}else{
			p._err(cm.ERR_0010, "Unsuppoted character: " + string(p.ch))
		}
	}

	// debug output
	//fmt.Printf("token:%d, i:%d f:%g s:%s\n", tok, lval.ival, lval.fval, lval.sval)

	return tok
}

// 実数を読み出す(小数点以下だけ)
func (p *Lexer) readRealNumber() float32{
	var num float64 = 0.0
	str := "0."

	for isNumber(p.ch){
		str += string(p.ch)
		p.readChar()
	}
	num,err := strconv.ParseFloat(str, 32)
	if err != nil{
		p._err(cm.ERR_0034, "value:"+str)
		return 0.0
	}
	p.backChar() //余計に読んだ文を１つ戻す

	return float32(num)
}

// 文字列リテラルを読み出す
func (p *Lexer) readStringLiteral() string{
	head := p.position
	p.readChar()
	for p.ch!='"'{
		if p.ch == 0 || p.ch == '\n'{
			p._err(cm.ERR_0011, "")
		}else if p.ch == '\\'{
			// \が出た場合はもう一文字無条件で読み出しておく
			// エスケープ文字の処理は、実行マシンに任せる
			p.readChar()
		}
		p.readChar()
	}

	// ""を抜いた形で返す
	return p.src[head+1:p.position]
}

// 特定の文字列からキーワードを検索
func getKeywordToken(key string)int {
	switch key{
	case "var": return VAR
	case "int": return INT
	case "float": return FLOAT
	case "string": return STRING
	case "void": return VOID
	case "if": return IF
	case "else": return ELSE
	case "switch": return SWITCH
	case "case": return CASE
	case "default": return DEFAULT
	case "fallthrough": return FALLTHROUGH
	case "for": return FOR
	case "break": return BREAK
	case "continue": return CONTINUE
	case "func": return FUNC
	case "return": return RETURN
	case "import": return IMPORT
	case "type": return TYPE
	case "struct": return STRUCT
	case "__syscall": return SYSCALL
	case "__dump" : return DUMP
	default: return IDENTIFIER
	}
}