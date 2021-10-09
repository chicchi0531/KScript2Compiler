package compiler

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"strconv"

)

//ノベル構文から、スクリプトにトランスパイルする
func Transpile(script string, setting *CompilerSettings) (string, error){
	scriptLines := strings.Split(script, "\n")
	begintag := regexp.MustCompile(`#[ ]*novel[ ]*{`)
	endtag := regexp.MustCompile(`}`)
	comment := regexp.MustCompile(`//.*`)
	var isNovel bool //ノベル構文中かどうかの状態フラグ


	var result string
	for lineno, str := range scriptLines {
		// コメントの削除
		str = comment.ReplaceAllString(str, "")

		str = strings.TrimSpace(str)

		// # novel行に囲まれた箇所を探す
		if begintag.MatchString(str){
			isNovel = true
			result += "\n"
			continue
		}else if endtag.MatchString(str) && isNovel{
			isNovel = false
			result += "\n"
			continue
		}

		//空行、ノベル構文は、そのまま出力
		if str == "" || !isNovel {
			result += fmt.Sprintln(str)
			continue
		}

		//変換
		trunspiledCode, err := TrunspileLine(str,lineno,setting)
		if err != nil{
			return "", err
		}
		result += fmt.Sprintln(trunspiledCode)
	}
	return result, nil
}

//１行をトランスパイルする
func TrunspileLine (script string, lineno int, setting *CompilerSettings) (string, error){
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
			if result[len(result)-3:] != "<p>"{
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
	return errors.New(fmt.Sprint(lineno, ":" ,errormsg))
}