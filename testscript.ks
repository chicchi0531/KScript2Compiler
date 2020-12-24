import "stdlib.ks";

// コメント
void main(){

	string chara_name1 = "アリス";
	string chara_name2 = "フィロ";

	//ノベル文
	@{
	- ナレーション
	ここはノベル文です。
	文末になにもない場合は、自動的に改行します。
	明示的な改行命令もあります。<n>

	- %chara_name1%
	%hoge%fuga」
	変数テスト%test%
	変数テスト%test%テスト
	%val1%%val2%%val3%

	- %chara_name2%
	名前が変数
	これはエラー%

	@}

	
	return;
}
