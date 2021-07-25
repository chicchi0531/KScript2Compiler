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

func OpenScriptFile(filename string) (string, error) {
	buf, err := ioutil.ReadFile(filename)
	return string(buf), err
}

func GetScriptFullpath(filename string) string {
	path := driver.CurrentDirectory + "/" + filename
	_, err := os.Stat(path)
	if err == nil {return path}

	// 見つからない場合は標準ライブラリフォルダから探す
	kspath := os.Getenv("KS_PATH")
	// KS_PATHがない場合は実行フォルダから探す
	if kspath == ""{
		kspath, _ = os.Executable()
		kspath = filepath.Dir(kspath)	
	}
	path = kspath + "/include/" + filename

	_, err = os.Stat(path)
	if err == nil {return path}

	return path
}

// ks -> ksobjへのコンパイル
func Compile(path string, isImport bool) int {
	errHandler := cm.MakeErrorHandler()
	currentdir, filename := filepath.Split(path)
	currentdir,_ = filepath.Abs(currentdir)

	// load script
	source, err := OpenScriptFile(path)
	if err != nil{
		fmt.Println(filename + " : スクリプトのロードに失敗しました。 " + err.Error())
		return -1
	}

	// ks -> ksil transpile
	scriptText, err := Transpile(source)
	if err != nil{
		fmt.Println(filename + " : スクリプトのトランスパイル [ks->ksil] に失敗しました" + err.Error())
		return -1
	}

	// output ksil
	makeDirectories("obj")
	ksilFilepath := "obj/" + getFilenameWithoutExt(filename) + ".ksil"
	ioutil.WriteFile(ksilFilepath, []byte(scriptText), os.ModePerm)

	// ksil -> ksobj compile
	lexer = MakeLexer(filename, scriptText, errHandler)
	if !isImport {
		driver = vm.MakeDriver(path, errHandler)
	}else{
		driver.Err = errHandler
		driver.Filename = path
	}
	driver.CurrentDirectory = currentdir
	result := Parse()

	// コンパイル後処理
	if !isImport{
		driver.LabelSettings()
	}

	// デバッグ出力
	dumpFile, err := os.Create("log.txt")
	if err != nil{
		fmt.Println("logファイルを開けません。")
	}
	defer dumpFile.Close()
	driver.Dump(dumpFile)

	// 結果出力
	if errHandler.WarningCount > 0{
		fmt.Printf(driver.Filename + "%d件の警告が発生しました。\n", errHandler.WarningCount)
	}
	if errHandler.ErrorCount > 0{
		fmt.Printf(driver.Filename + " : コンパイルに失敗しました。%d件のエラーが発生しました。\n", errHandler.ErrorCount)
		return result
	}

	// ファイルへの書き出し
	if !isImport{
		err = OutputFiles(driver)
		if err != nil{
			fmt.Println("ファイルの出力に失敗しました。 : " + err.Error())
		}
	}

	return result
}

// スクリプト内でのimport命令処理
var imported []string
func ImportFile(path string) int {
	currentLexer := lexer
	currentDir := driver.CurrentDirectory
	currentFile := driver.Filename
	currentErr := driver.Err

	// ファイルの重複ロードチェック
	path = GetScriptFullpath(path)
	for _, i := range imported{
		if path == i{
			return 0
		}
	}
	imported = append(imported, path)

	// パース処理
	result := Compile(path, true)

	// 復帰
	lexer = currentLexer
	driver.CurrentDirectory = currentDir
	driver.Filename = currentFile
	driver.Err = currentErr

	return result
}

// ksil -> ksobjへの変換
func Parse() int {

	result := yyParse(lexer)

	return result
}

// ファイルへの書き出し
func OutputFiles(d *vm.Driver) error {

	// バイナリファイルの書き出し
	outdir := "bin/"
	makeDirectories(outdir)
	outpath := outdir + getFilenameWithoutExt(d.Filename) + ".ksobj"
	file, err := os.Create(outpath)
	if err != nil{
		return err
	}
	defer file.Close()

	// ヘッダ情報の書き込み
	// 0x0000-0x0002 prefix "ks2"
	// 0x0003-0x0004 file version

	_, err = file.Write([]byte("ks2"))
	if err != nil {return err}

	_, err = file.Write([]byte{compilerVersion_Major, compilerVersion_Minor})
	if err != nil {return err}

	// データの書き込み
	// 0x0005 data
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

	// float tableの書き込み
	_,err = file.Write([]byte{vm.VMCODE_FLOATTABLE})
	if err != nil {return err}

	for _,val := range d.FloatTable.Values{
		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.BigEndian, float32(val))
		if err != nil {return err}

		_, err = file.Write(buf.Bytes())
		if err != nil {return err}
	}

	// string tableの書き出し(csv形式)
	locale := "jp"
	strTableOutDir := outdir + locale + "/"
	makeDirectories(strTableOutDir)
	strTableOutPath := strTableOutDir + getFilenameWithoutExt(d.Filename) + ".ksdat"

	stFile,err := os.Create(strTableOutPath)
	if err != nil{return err}
	defer stFile.Close()

	for i,str := range d.StringTable.Values{
		_,err = fmt.Fprintf(stFile, "%d,\"%s\"\n",i,str)
		if err != nil{return err}
	}

	return nil
}

func getFilenameWithoutExt(path string) string {
	return filepath.Base( path[:len(path) - len(filepath.Ext(path))] )
}

func makeDirectories(path string) error{
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err){
		os.MkdirAll(path, os.ModePerm)
	}
	return err
}