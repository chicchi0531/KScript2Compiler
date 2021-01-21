// FizzBuzz Program

import "stdlib.ks";

int create_sprite(string path)
{
	int h = load_sprite(path);
	return instantiate_sprite(h);
}
int create_image(string path)
{
	int h = load_sprite(path);
	return instantiate_image(h);
}

void main()
{
	int hAlice = create_sprite("tex_novel_alice.png");
	int hFiro = create_image("tex_novel_firo.png");

	set_window_id(0, "");
	set_window_id(1, "アリス");
	set_window_id(2, "フィロ");

	set_sprite_pos(hAlice, 10.0, 10.0, 0);
	set_image_pos(hFiro, 100.0, 100.0, 0);


	@{
		- null
		ここはナレーション分です。」
		- フィロ
		こんにちは。
		これは掛け合い用のレイアウトのテストよ。
		会話は３行x23文字まで出せるわ。」
		- アリス
		クフフッ！
		話者を変えるとウィンドウが変わるぞ。
		ウィンドウの位置もスクリプトで調整できる。」
		- フィロ
		話者を戻すとこんな感じね。」
		- アリス
		ちゃんと掛け合いができているな。」
	@}
}

