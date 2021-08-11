package compiler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CompilerSettings struct{
	SystemcallID_showmsg int	`json:"systemcallid_showmsg"`
	SystemcallID_showname int	`json:"systemcallid_showname"`
	SystemcallID_playvoice int	`json:"systemcallid_playvoice"`
}

func MakeDefaultSettings() *CompilerSettings{
	result := new(CompilerSettings)
	result.SystemcallID_playvoice = 2
	result.SystemcallID_showmsg = 0
	result.SystemcallID_showname = 1
	return result
}

func LoadSettingFile(filename string) *CompilerSettings{
	bytes, err := ioutil.ReadFile(filename)
	if err != nil{
		//実行ファイルパスで再試行
		exe, _ := os.Executable()
		path := filepath.Dir(exe)
		bytes, err = ioutil.ReadFile(path + "/" + filename)
		if err != nil{
			fmt.Println("設定ファイルが見つかりません。既存の設定を利用します。")
			return MakeDefaultSettings()
		}
	}

	// json decode
	var settings CompilerSettings
	if err := json.Unmarshal(bytes, &settings); err != nil{
		fmt.Println("設定ファイルの読み出しに失敗しました。既存の設定を利用します。")
		return MakeDefaultSettings()
	}

	return &settings
}