// include用ライブラリ

//システムコール宣言
system void print(string msg);
system string itos(int value);

int hoge = 12345;

int sub(int a, int b)
{
	return a-b;
}

int add(int a, int b)
{
	print(itos(sub(a,b)));
	return a + b;
}
