import "stdlib.ks";

// コメント
void main(){

	string chara_name1 = "アリス";
	string chara_name2 = "フィロ";
	int hoge;
	int test;
	test = 0;
	float fvalue = 1.5;

	//ノベル文
	@{
	- ナレーション
	ここはノベル文です。」
	文末になにもない場合は、自動的に改行します。」
	明示的な改行命令もあります。<n>
	これは明示的な改ページ命令です。<p>

	- %chara_name1%
	%hoge%fuga」
	変数テスト%test%」
	変数テスト%test%テスト」

	- %chara_name2%
	名前が変数

	<n>

	@}

	
	return;
}
