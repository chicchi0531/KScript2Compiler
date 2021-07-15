package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cm "ks2/compiler/common"
	"ks2/compiler/vm"
)

var errHandler *cm.ErrorHandler
var compilerVersion_Major byte = 1
var compilerVersion_Minor byte = 0

func OpenScriptFile(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// ks -> ksobjへのコンパイル
func Compile(path string) int {
	errHandler := cm.MakeErrorHandler()
	_, filename := filepath.Split(path)

	println(" スクリプトロード")
	source := OpenScriptFile(path)

	println(" スクリプトのトランスパイル")
	scriptText, err := Transpile(source)
	if err != nil{
		panic(err)
	}
	if _, err = os.Stat("obj"); os.IsNotExist(err){
		if err := os.Mkdir("obj", os.ModePerm); err != nil{
			panic(err)
		}
	}
	ioutil.WriteFile("obj/_transpiled_" + filename, []byte(scriptText), os.ModePerm)
	println(" 完了")

	println(" スクリプトのコンパイル")
	lexer = MakeLexer(filename, scriptText, errHandler)
	driver = vm.MakeDriver(filename, errHandler)
	result := Parse()

	// コンパイルご処理
	driver.LabelSettings()

	// デバッグ出力
	driver.Dump()

	// 結果出力
	if errHandler.WarningCount > 0{
		fmt.Printf("%d件の警告が発生しました。\n", errHandler.WarningCount)
	}
	if errHandler.ErrorCount > 0{
		fmt.Printf("コンパイルに失敗しました。%d件のエラーが発生しました。\n", errHandler.ErrorCount)
		return result
	}else{
		fmt.Printf("コンパイルに成功しました。\n")
	}

	// ファイルへの書き出し
	err = OutputFiles(driver)
	if err != nil{
		println(err.Error())
	}

	return result
}

// スクリプト内でのimport命令処理
func ImportFile(filename string) int {
	currentLexer := lexer

	// パース処理
	result := Compile(filename)

	// lexerを元に復帰
	lexer = currentLexer

	return result
}

// ksil -> ksobjへの変換
func Parse() int {

	fmt.Println("■コンパイル開始 " + driver.Filename)
	result := yyParse(lexer)
	fmt.Println("■コンパイル完了 " + driver.Filename)

	return result
}

// ファイルへの書き出し
func OutputFiles(d *vm.Driver) error {

	// バイナリファイルの書き出し
	makeDirectories("bin/")
	outpath := "bin/" + getFilenameWithoutExt(d.Filename) + ".ksobj"
	file, err := os.Create(outpath)
	if err != nil{
		return err
	}
	defer file.Close()

	// ヘッダ情報の書き込み
	// 0x0000-0x0002 prefix "ks2"
	// 0x0003-0x0004  file version

	_, err = file.Write([]byte("ks2"))
	if err != nil {return err}

	_, err = file.Write([]byte{compilerVersion_Major, compilerVersion_Minor})
	if err != nil {return err}

	// データの書き込み
	// 0x0004 data
	for _,prog := range d.Program{
		if prog.Code != vm.VMCODE_DUMMYLABEL{
			// 命令の書き込み
			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.BigEndian, int8(prog.Code))
			if err != nil {return err}

			_,err = file.Write(buf.Bytes())
			if err != nil {return err}

			// 値の書き込み
			buf = new(bytes.Buffer)
			err = binary.Write(buf, binary.BigEndian, int32(prog.Value))
			if err != nil {return err}

			_, err = file.Write(buf.Bytes())
			if err != nil {return err}
		}
	}

	return nil
}

func getFilenameWithoutExt(path string) string {
	return filepath.Base( path[:len(path) - len(filepath.Ext(path))] )
}

func makeDirectories(path string) error{
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err){
		os.Mkdir(path, 0777)
	}
	return err
}