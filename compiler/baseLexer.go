package compiler

import(
	cm "ks2/compiler/common"
)

// 最低限必要な構造体を定義
type BaseLexer struct {
	filename string
	src    string
	position  int //現在の位置
	readPosition int //次の読み出し位置
	ch byte //現在の文字
	line int //現在の行数
	err *cm.ErrorHandler //エラーハンドラ
}

func (p *BaseLexer) Error(err string){
	p.err.LogError(p.filename, p.line, cm.ERR_0004, err)
}

// こちらは内部用
func (p *BaseLexer) _err(errcode string, submsg string){
	p.err.LogError(p.filename, p.line, errcode, submsg)
}

// 空白スキップ
func (p *BaseLexer) skipWhiteSpace(){
	for p.ch==' ' || p.ch =='\t' || p.ch=='\r'{
		p.readChar()
	}
}

// コメントスキップ
func (p *BaseLexer) skipComments(){
	// EOFが出てきた場合はスキップ
	if p.ch == 0 || p.nextChar() == 0{
		return
	}

	if p.ch == '/'{
		// 1行コメントをスキップ
		if p.nextChar() == '/'{
			p.readChar()
			for p.ch != '\n'{
				p.readChar()
			}
		// 複数行コメントをスキップ
		}else if p.nextChar() == '*'{
			p.readChar()
			p.readChar()
			for string(p.ch) + string(p.nextChar()) != "*/"{
				p.readChar()
			}
			p.readChar()
			p.readChar()
		}
	}
}

// １文字読み出す
func (p *BaseLexer) readChar(){
	if p.readPosition>=len(p.src){
		p.ch = 0
	}else{
		p.ch = p.src[p.readPosition]
	}
	p.position = p.readPosition
	p.readPosition++

	// 行数カウント
	if p.ch == '\n'{
		p.line++
	}
}

//　次の１文字を確認する（位置は動かさない)
func (p *BaseLexer) nextChar() byte{
	if p.readPosition>=len(p.src){
		return 0
	}else{
		return p.src[p.readPosition]
	}
}

// 1文字読み出し位置を戻す
func (p *BaseLexer) backChar(){
	if p.position > 0{
		// 現在が改行文字なら、戻すときに行数カウントも戻す
		if p.ch == '\n'{
			p.line--
		}
		p.readPosition--
		p.position--
		p.ch = p.src[p.position]
	}
}

// 数字を読み出す
func (p *BaseLexer) readNumber() int{
	num := 0
	for isNumber(p.ch){
		num *= 10
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