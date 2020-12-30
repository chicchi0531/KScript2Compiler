int hoge;

/* 
複数行コメント
*/

void main();

int main()
{
	int i;
	i = 100;
	int i=0; //同時初期化
	float f=1.1; //同時初期化
	comp = a < b && a >= c;
	uni = "unicode test";
	uni = "UTF8てすと";
	hoge = 200;
	hoge = "hogehoge";

	// if文
	if(a>0)
	{
		b="if";
	}

	//ifelse文
	if(a>0)
	{
		b="if";
	}else
	{
		b="else";
	}

	//ネスト文
	if(a>0)
	{
		b="if";
		if(a==b)
		{
			b="if in if";
		}
	}
	else if(a==b)
	{
		b="else";
		if(a==b)
		{
			b="if in elseif";
		}
	}
	else
	{
		b="else";
	}

	// for文テスト
	for(int i=0; i<10; i+=1)
	{
		print("hogehoge");
		if(z>0)
		{
			continue;
		}
	}

	// func call test
	a = add(1,b);
	b = add(a+b, c+d);
	add(1,2);

	//while
	while(a<b)
	{
		a+=1;
		break;
	}

}


int add(int a, int b)
{
	result = a + b;
	return result;
}