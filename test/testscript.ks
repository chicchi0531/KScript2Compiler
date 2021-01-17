// FizzBuzz Program

import "stdlib.ks";

void main()
{
	for(int i=0; i<20; i+=1)
	{
		if(i%3 == 0) print("Fizz");
		if(i%5 == 0) print("Buzz");
		if(i%3 !=0 && i%5 !=0)
		{
			string istr = itos(i);
			print(istr);
		}
		print(",");
	}

	//演算子テスト
	int a = 1;
	float b = 2.5;
	string c = "hogehoge";

	// downcast
	int x = a + b;
	x = b + a;

	// upcast
	float y = a + b;
	y = b + a;

	// nocast
	float z = b + b;
	int w = a + a;

	// string
	string f = c + "fugafuga";
}

