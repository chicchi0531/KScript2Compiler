import "stdlib.ks";

void main()
{
    int hAlice = load_sprite("alice.png");
    int hBg = load_sprite("bg.png");

    float posx = 0.0;
    float posy = 0.0;
    int layer = 0;
    float degree = 0.0;

    //背景描画
    draw_sprite(hBg, 0.0, 0.0, 10);

    show_window();
    @{
        - ナレーション
        画像テストです。」
    @}
    hide_window();

    while(1)
    {
        //移動プログラム
        if(get_button_down("up")){ posy += 0.1; }
        if(get_button_down("down")){ posy -= 0.1; }
        if(get_button_down("right")){ posx += 0.1; }
        if(get_button_down("left")){ posx -= 0.1; }

        degree += 50.0;

        draw_sprite(hAlice, posx, posy, layer);
        rotate_sprite(hAlice, degree);
        await();
    }

}