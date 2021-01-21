// include用ライブラリ

//システムコール宣言
system void print(string msg);
system string itos(int value);

// 画像関係
system int load_sprite(string path);
system void unload_sprite(int handle);

system int instantiate_sprite(int handle);
system int instantiate_image(int handle);

system void set_sprite_pos(int handle, float x, float y, int layer);
system void set_image_pos(int handle, float x, float y, int layer);

system void rotate_sprite(int handle, float degree);
system void rotate_image(int handle, float degree);

system void delete_sprite(int handle);
system void delete_image(int handle);

// 同期命令
system void await();

// input関係
system int get_button_down(string name);
system float get_input_value(string name);

// ノベル関係
system void show_window();
system void hide_window();
system void set_window_id(int id, string name);