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
	SystemcallID_set_choice_title int 	`json:"systemcallid_set_choice_title`
	SystemcallID_set_choice_fravor int 	`json:"systemcallid_set_choice_fravor`
	SystemcallID_set_choice_image int 	`json:"systemcallid_set_choice_image`
	SystemcallID_add_choice_option int 	`json:"systemcallid_add_choice_option`
	SystemcallID_show_choice int 		`json:"systemcallid_show_choice`
}

func MakeDefaultSettings() *CompilerSettings{
	result := new(CompilerSettings)
	result.SystemcallID_playvoice = 2
	result.SystemcallID_showmsg = 0
	result.SystemcallID_showname = 1
	result.SystemcallID_set_choice_title = 3
	result.SystemcallID_set_choice_fravor = 4
	result.SystemcallID_set_choice_image = 5
	result.SystemcallID_add_choice_option = 6
	result.SystemcallID_show_choice = 7
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