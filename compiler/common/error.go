package common

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
	ERR_0015 = "KS0015: 定義済みの識別子が再定義されました。"
	ERR_0016 = "KS0016: 識別子が定義されていません。"
	ERR_0017 = "KS0017: 互換性のない型の代入が行われました。"
	ERR_0018 = "KS0018: 不正な変数をPushしようとしました。"
	ERR_0019 = "KS0019: 不正な変数をPopしようとしました。"
	ERR_0020 = "KS0020: 関数呼び出しエラー。引数の数が一致しません。"
	ERR_0021 = "KS0021: 関数呼び出しエラー。引数の型が一致しません。"
	ERR_0022 = "KS0022: 関数呼び出しエラー。未定義の関数が呼び出されました。"
	ERR_0023 = "KS0023: 関数定義エラー。同名の関数が定義済みです。"
	ERR_0024 = "KS0024: 関数定義エラー。戻り値の型が定義と一致しません。"
	ERR_0025 = "KS0025: 関数定義エラー。値を返さないコードパスがあります。"
)

const(
	WARNING_0001 = "KS2000: 暗黙的型変換が行われました。"
)

type ErrorHandler struct{
	ErrorCount int
	WarningCount int
}

func (e *ErrorHandler) LogError(filename string, lineno int, errorcode string, subMsg string) {
	fmt.Printf("%s [%d]: error %s %s\n",filename, lineno, errorcode, subMsg)
	e.ErrorCount++
}

func (e *ErrorHandler) LogWarning(filename string, lineno int, warningcode string, subMsg string){
	fmt.Printf("%s [%d]: warning %s %s\n",filename, lineno, warningcode, subMsg)
	e.WarningCount++
}
