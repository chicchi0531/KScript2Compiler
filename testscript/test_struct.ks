
type SubType struct{
    subMember int
}

type TestType struct{
    member1 int
    member2 float
    member3 string
    member4 SubType
}

var user TestType
var user2 TestType

func main(){
    user.member1 = 12345
    user.member2 = 123.45
    user.member3 = "hogehoge"
    user.member4.subMember = user.member1

    user2 = user
}