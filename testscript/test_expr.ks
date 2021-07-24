
func main(){
        // 変数テスト
    //----------------
    // 0: a1
    // 1: a2
    // 2: a3(deleted)
    // 3: a3
    // 4: a4(deleted)
    // 5: a4
    //----------------
    // pushfloat(123.4)
    // pushint 1
    // popvalue
    // pushstring(hogehoge)
    // pushint 3
    // popvalue
    // pushint(12345)
    // pushint 5
    // popvalue
    var a1 int
    var a2 float = 123.4
    var a3 = "hogehoge"
    a4 := 12345
    
    // 計算テスト
    // 6: b1
    // 7: b2(deleted)
    // 8: b2
    // 9: b3(deleted)
    //10: b3
    //-----------
    // pushint 7  [b2 := 7]
    // pushint 8
    // popvalue
    // pushint 19 [b3 := -19]
    // not
    // pushint 10
    // popvalue
    // pushint 8  [b1 = b2 + b3 * 5 / 2]
    // pushvalue
    // pushint 10
    // pushvalue
    // pushint 5
    // mul
    // pushint 2
    // div
    // add
    // pushint 6
    // popvalue
    // pushint 6  [b2 = (b1 - b2) * -2 % 2]
    // pushvalue
    // pushint 8
    // pushvalue
    // sub
    // pushint 2
    // not
    // mul
    // pushint 2
    // mod
    // pushint 8
    // popvalue
    // pushint 6  [b3 = (b1 == b2) || (b1 != b3)]
    // pushvalue
    // pushint 8
    // pushvalue
    // equ
    // pushint 6
    // pushvalue
    // pushint 10
    // pushvalue
    // neq
    // or
    // pushint 10
    // popvalue
    // pushint 6  [(b1 > b2) && (b1 < b3)]
    // pushvalue
    // pushint 8
    // pushvalue
    // gt
    // pushint 6
    // pushvalue
    // pushint 10
    // pushvalue
    // lt
    // and
    // pushint 10
    // popvalue
    // pushint 6 [b1 <= b2]
    // pushvalue
    // pushint 8
    // pushvalue
    // le
    // pushint 10
    // popvalue
    // pushint 6 [b1 >= b2]
    // pushvalue
    // pushint 8
    // pushvalue
    // ge
    // pushint 10
    // popvalue
    // pushint 6 [b1 += b2]
    // pushvalue
    // pushint 8
    // pushvalue
    // add
    // pushint 6
    // popvalue
    // pushint 6 [b1 -= b2]
    // pushvalue
    // pushint 8
    // pushvalue
    // sub
    // pushint 6
    // popvalue
    // pushint 6 [b1 *= b2]
    // pushvalue
    // pushint 8
    // pushvalue
    // mul
    // pushint 6
    // popvalue
    // pushint 6 [b1 /= b2]
    // pushvalue
    // pushint 8
    // pushvalue 
    // div
    // pushint 6
    // popvalue
    // pushint 6 [b1++]
    // pushvalue
    // incr
    // pushint 6
    // popvalue
    // pushint 8 [b2--]
    // pushvalue
    // decr
    // pushint 8
    // popvalue
    var b1 int
    b2 := 7
    b3 := -19
    b1 = b2 + b3 * 5 / 2
    b2 = (b1 - b2) * -2 % 2
    b3 = (b1 == b2) || (b1 != b3)
    b3 = (b1 > b2) && (b1 < b3)
    b3 = b1 <= b2
    b3 = b1 >= b2
    b1 += b2
    b1 -= b2
    b1 *= b2
    b1 /= b2
    b1 %= b2
    b1++
    b2--

    __dump("dump.log")
}