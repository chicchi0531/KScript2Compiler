// FizzBuzz Program

import "stdlib.ks";
system void print(string msg);
system string itos(int value);

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

	int a = 0;
	int b = 1;
	int result = add(a,b);
}

