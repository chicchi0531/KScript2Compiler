package compiler

import (
	"fmt"
)

// 最低限必要な構造体を定義
type Lexer struct {
	src    string
	position  int //現在の位置
	readPosition int //次の読み出し位置
	ch byte //現在の文字
	line int //現在の行数
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
		}
		tok = PLUS // +

	case '-':
		ch := p.nextChar()
		if ch == '-'{
			p.readChar()
			tok = DECR // --
		}else if ch == '=' {
			p.readChar()
			tok = M_EQ // -=
		}
		tok = MINUS // -

	case '*':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = A_EQ // *=
		}
		tok = ASTARISK

	case '/':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = S_EQ // /=
		}
		tok = SLASH

	case '=':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = EQ
		}
		tok = ASSIGN

	case '!':
		p.readChar()
		if p.ch == '='{
			tok = NEQ
		}
		p.logError("シンタックスエラー。!のあとに=がありません。")
	
	case '>':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = GE
		}
		tok = GT

	case '<':
		ch := p.nextChar()
		if ch == '='{
			p.readChar()
			tok = LE
		}
		tok = LT

	case '&':
		p.readChar()
		if p.ch == '&'{
			tok = AND
		}
		p.logError("シンタックスエラー。&のあとに&がありません。")

	case '|':
		p.readChar()
		if p.ch == '|'{
			tok = OR
		}
		p.logError("シンタックスエラー。|のあとに|がありません。")

	case '\n':
		p.line++
		tok = EOL

	case '(': fallthrough
	case ')': fallthrough
	case '{': fallthrough
	case '}': fallthrough
	case '[': fallthrough
	case ']': fallthrough
	case ':': fallthrough
	case ';': fallthrough
	case ',': fallthrough
	case '.': tok = int(p.ch)

	default:
		if isNumber(p.ch){
			num := p.readNumber()

			// 少数点数の場合
			if p.ch == '.'{
				p.readChar()
				decimal := p.readNumber()
				lval.fval = getFloat(num, decimal)

				// floatのsuffixを読み飛ばす
				if p.ch == 'f' || p.ch == 'F'{
					p.readChar()
				}
				tok = FNUM
			}

			// 整数の場合
			lval.ival = num
			tok = INUM
		}else if isLetter(p.ch){
			lval.sval = p.readIdentifier()
			tok = IDENTIFIER
		}else{
			p.logError("意図しない文字が検出されました。" + string(p.ch))
		}
	}

	return tok
}

// 空白スキップ
func (p *Lexer) skipWhiteSpace(){
	for p.ch==' ' || p.ch =='\t' || p.ch=='\r'{
		p.readChar()
	}
}

// コメントスキップ
func (p *Lexer) skipComments(){
	// EOFが出てきた場合はスキップ
	if p.ch == 0 || p.nextChar() == 0{
		return
	}

	str := p.src[p.position:p.position+1]

	// 一行コメントをスキップ
	if str=="//"{
		for p.ch!='\n'{
			p.readChar()
		}
	}
	// マルチラインコメントをスキップ
	if str=="/*"{
		for str != "*/"{
			str = p.src[p.position-1:p.position]
			p.readChar()
		}
	}
}

// １文字読み出す
func (p *Lexer) readChar(){
	if p.readPosition>=len(p.src){
		p.ch = 0
	}else{
		p.ch = p.src[p.readPosition]
	}
	p.position = p.readPosition
	p.readPosition++
}

//　次の１文字を確認する（位置は動かさない)
func (p *Lexer) nextChar() byte{
	if p.readPosition>=len(p.src){
		return 0
	}else{
		return p.src[p.readPosition]
	}
}

// 1文字読み出し位置を戻す
func (p *Lexer) backChar(){
	if p.position > 0{
		p.readPosition--
		p.position--
		p.ch = p.src[p.position]
	}
}

// 数字を読み出す
func (p *Lexer) readNumber() int{
	num := 0
	for isNumber(p.ch){
		num += (int)(p.ch - '0')
		p.readChar()
	}
	p.backChar()//余計に読んだ文を１つ戻す
	return num
}

// 識別子を読み出す
func (p *Lexer) readIdentifier() string{
	head := p.position
	for isLetter(p.ch) || isNumber(p.ch){
		p.readChar()
	}
	p.backChar()//余計に読んだ分を１つ戻す
	return p.src[head:p.position+1]
}

// 文字列リテラルを読み出す
func (p *Lexer) readStringLiteral() string{
	head := p.position
	p.readChar()
	for p.ch!='"'{
		if p.ch == 0 || p.ch == '\n'{
			p.logError("文字リテラルが閉じられていません。")
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

// 文字かどうかの判定
func isLetter(ch byte)bool{
	// アルファベットの範囲にあるかどうか
	if ('a'<=ch && ch<='z') || ('A'<=ch && ch<='Z') || ch=='_'{
		return true
	
	// asciiの外にあるかどうか
	}else if ch>=0x80{
		return true
	
	// asciiの範囲かつ、アルファベットじゃなければ記号か数字とみなす
	}else{
		return false
	}
}

// 数字かどうかの判定
func isNumber(ch byte)bool{
	return '0' <= ch && ch <= '9'
}

// エラー報告用
func (p *Lexer) logError(e string) {
	panic(fmt.Sprintf("Error[%d] %s", p.line ,e))
}

// 整数部と小数点部からfloat32を構成
func getFloat(number int, decimal int) float32{
	var result float32
	result = float32(decimal)

	// 0.xxになるまで10で割り続ける
	for result > 1{
		result /= 10.0
	}
	return float32(number) + result
}