// include用ライブラリ

//システムコール宣言
system int load_sprite(string path);
system void draw_sprite(int handle, float x, float y, int layer);
system void rotate_sprite(int handle, float degree);
system void hide_sprite(int handle);
system void unload_sprite(int handle);
system void await();

system int get_button_down(string name);
system float get_input_value(string name);

system void show_window();
system void hide_window();