package compiler

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//ノベル構文から、スクリプトにトランスパイルする
func Transpile(script string, setting *CompilerSettings) (string, error){
	scriptLines := strings.Split(script, "\n")
	novelBeginTag := regexp.MustCompile(`#[ ]*novel[ ]*{`)
	choiceBeginTag := regexp.MustCompile(`#[ ]*choice[ ]*{`)
	endTag := regexp.MustCompile(`}`)
	comment := regexp.MustCompile(`//.*`)

	// ノベル構文用
	var isNovel bool //ノベル構文中かどうかの状態フラグ

	// 選択肢構文用
	var isChoice bool //選択肢構文中かどうかのフラグ
	var optionStack []string //次回の選択肢をためておくスタック
	var choiceResultVariable string //選択肢の結果が格納される変数名

	var result string
	for lineno, str := range scriptLines {
		// コメントの削除
		str = comment.ReplaceAllString(str, "")
		//空白の削除
		str = strings.TrimSpace(str)

		// ノベル開始構文
		if novelBeginTag.MatchString(str){
			isNovel = true
			result += "\n"
			continue
		// ノベル終了構文
		}else if endTag.MatchString(str) && isNovel{
			isNovel = false
			result += "\n"
			continue
		// 選択肢開始構文
		}else if choiceBeginTag.MatchString(str){
			isChoice = true
			optionStack = nil //スタックを初期化
			choiceResultVariable = ""
			result += "\n"
			continue
		// 選択肢構文終了
		}else if endTag.MatchString(str) && isChoice{
			isChoice = false

			if len(optionStack) == 0{
				return "", logerror("選択肢構文には1つ以上の@option命令を含んでください。", lineno)
			}

			//選択肢追加用の命令を追加
			for _,c := range optionStack{
				result += fmt.Sprintf("__syscall[%d](%s)\n",setting.SystemcallID_add_choice_option, c)
			}

			//選択肢表示命令を追加
			//result指定があれば、代入処理する
			if choiceResultVariable ==""{
				result += fmt.Sprintf("__syscall[%d]()\n",setting.SystemcallID_show_choice)
			}else{
				result += fmt.Sprintf("var %s int\n",choiceResultVariable)
				result += fmt.Sprintf("%s = __syscall[%d]()\n",choiceResultVariable, setting.SystemcallID_show_choice)
			}
			continue
		}

		//空行は、そのまま出力
		if str == "" {
			result += fmt.Sprintln(str)
		
		// ノベル構文はトランスパイル
		}else if isNovel {
			trunspiledCode, err := TruspileNovelLine(str,lineno,setting)
			if err != nil{
				return "", err
			}
			result += fmt.Sprintln(trunspiledCode)
		
		// 選択肢構文はトランスパイル
		}else if isChoice{
			trunspiledCode, err := TrunspileChoiceLine(str,lineno,setting,&optionStack, &choiceResultVariable)
			if err != nil{
				return "", err
			}
			result += trunspiledCode

		// それ以外はそのまま出力
		}else{
			result += fmt.Sprintln(str)
		}
	}
	return result, nil
}

//１行をトランスパイルする
func TruspileNovelLine (script string, lineno int, setting *CompilerSettings) (string, error){
	script = strings.TrimSpace(script)
	if len(script) == 0{
		return "",nil
	}

	var result string
	switch script[0]{
		case '@'://コマンド行　関数に変換
			commands := strings.Split(script[1:], " ")
			commands = RemoveAll(commands, "")
			if len(commands) == 0{
				return "", logerror("コマンド行の書式が正しくありません。",lineno)
			}

			//関数名
			result = commands[0] + "( "

			//引数
			for _, c := range commands[1:]{
				result += c + ","
			}
			result = result[:len(result)-1]
			result = result + ")"

		case '-'://名前行　名前表示関数に変換
			name := strings.TrimSpace(script[1:])
			if name == "nil"{
				name = ""
			}
			result = fmt.Sprintf("__syscall[%d](\"%s\")",setting.SystemcallID_showname, name)

		case '+'://ボイス行　ボイス再生関数に変換
			voice := strings.TrimSpace(script[1:])
			if voice == ""{
				return "", logerror("ボイス行の書式が正しくありません。", lineno)
			}
			result = fmt.Sprintf("__syscall[%d](\"%s\")",setting.SystemcallID_playvoice, voice)

		default://文字表示行　文字表示関数に変換
			//変数埋め込み処理
			variablePattern := regexp.MustCompile((`%([^%]*)%`))
			varList := variablePattern.FindAllStringSubmatch(script, -1)
			result = variablePattern.ReplaceAllString(script, "%s")
			for i,_ := range varList{
				result = strings.Replace(result, "%s", "$"+strconv.Itoa(i),1)
			}

			//ルビ命令の処理
			rubyPattern := regexp.MustCompile(`\[([^|]*)\|([^\]]*)\]`)
			result = rubyPattern.ReplaceAllString(result, "<r>$1|$2</r>")

			//改ページ命令の処理
			result = strings.ReplaceAll(result, "」", "<p>")

			//改行命令の処理
			if len(result) < 3 || result[len(result)-3:] != "<p>"{
				result += "<n>"
			}

			result = fmt.Sprintf("__syscall[%d](\"%s\",",setting.SystemcallID_showmsg, result)
			//変数埋め込みを追加
			for _, v := range varList{
					result += v[1] + ","
			}
			result = result[:len(result)-1] + ")"
	}

	return result, nil
}

// 選択肢構文
func TrunspileChoiceLine(script string, lineno int, setting *CompilerSettings, optionStack *[]string, choiceVal *string) (string, error){
	script = strings.TrimSpace(script)
	if len(script) == 0{
		return "",nil
	}
	if script[0] != '@'{
		return "", logerror("choice構文は@から始まる命令である必要があります。",lineno)
	}

	commands := SplitCommandString(script[1:])
	if len(commands) < 2{
		return "", logerror("コマンドの書式が正しくありません。",lineno)
	}

	var result string
	switch commands[0]{
		case "name":
			result = fmt.Sprintf("__syscall[%d](%s)\n",setting.SystemcallID_set_choice_title, commands[1])
		case "fravor":
			result = fmt.Sprintf("__syscall[%d](%s)\n",setting.SystemcallID_set_choice_fravor, commands[1])
		case "image":
			result = fmt.Sprintf("__syscall[%d](%s)\n",setting.SystemcallID_set_choice_image, commands[1])
		case "result":
			*choiceVal = commands[1]
		case "option":
			*optionStack = append(*optionStack, commands[1])
		default:
			return "", logerror("未定義のコマンド"+commands[0]+"が使用されました。フォーマットを確認してください。",lineno)
	}

	return result,nil
}

// スライスから特定の要素を除外する
func RemoveAll(slice []string, word string) []string{
	var result []string
	for _, str := range slice{
		if str != word{
			result = append(result, str)
		}
	}
	return result
}

//エラー出力
func logerror(errormsg string, lineno int) error {
	return errors.New(fmt.Sprint(lineno+1, ":" ,errormsg))
}

// 命令構文用の文字列スプリット。セパレータはスペース固定
func SplitCommandString(script string) []string{
	var result []string
	var tmpstr []rune
	isInStringLiteral := false
	for _, c := range script{
		
		// 空白文字を検出したら、コマンドを追加
		if !isInStringLiteral && c==' '{
			result = append(result, string(tmpstr))
			tmpstr = nil
		}else if c=='"'{
			isInStringLiteral = !isInStringLiteral
			tmpstr = append(tmpstr,c)
		}else{
			tmpstr = append(tmpstr,c)
		}
	}

	if len(tmpstr) > 0{
		result = append(result, string(tmpstr))
	}

	//空文字列は削除
	result = RemoveAll(result, "")
	return result
}