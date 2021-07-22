
type SubType struct{
    sub1 int
    sub2 float
}

type TestType struct{
    member1 [3]int
    member2 float
    member3 [3]SubType
}

func main(){
    var a [3]SubType
    var b TestType

    a = b.member3

    __dump("構造体テスト途中ログ.log")
}
