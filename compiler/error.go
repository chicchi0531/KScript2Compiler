package compiler

import "fmt"

const (
	ERR_0001 = "KS0001: 予期せぬエラーが発生しました。このエラーが頻繁に発生する場合は、Koromosoftに問い合わせてください。"
	ERR_0002 = "KS0002: 未定義の識別子が参照されました。定義を確認してください。"
	ERR_0003 = "KS0003: この関数は値を返さないコードパスがあります。"
	ERR_0004 = "KS0004: シンタックスエラー。"
	ERR_0005 = "KS0005: インクリメント演算子はstring型に適用できません。"
	ERR_0006 = "KS0006: デクリメント演算子はstring型に適用できません。"
	ERR_0007 = "KS0007: not演算子はstring型に適用できません。"
	ERR_0008 = "KS0008: 変数以外のノードがポップされました。このエラーが出た場合Koromosoftに連絡してください。"
	ERR_0009 = "KS0009: 不明なトークンが使用されました。"
	ERR_0010 = "KS0010: サポートされていない文字が検出されました。"
	ERR_0011 = "KS0011: 文字リテラルが閉じられていません。"
	ERR_0012 = "KS0012: 文字列と数値型の演算は出来ません。"
	ERR_0013 = "KS0013: 文字列ではできない演算です。"
	ERR_0014 = "KS0014: 不明な演算子が使用されました。"
)

type ErrorHandler struct{
	errorCount int
	warningCount int
}

func (e *ErrorHandler) LogError(filename string, lineno int, errorcode string, subMsg string) {
	fmt.Printf("%s [%d]: error %s %s",filename, lineno, errorcode, subMsg)
	e.errorCount++
}

func (e *ErrorHandler) LogWarning(filename string, lineno int, warningcode string, subMsg string){
	fmt.Printf("%s [%d]: warning %s %s",filename, lineno, warningcode, subMsg)
	e.warningCount++
}
