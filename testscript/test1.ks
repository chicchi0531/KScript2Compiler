
func main(){
    var a = 10 + 20 * 20 / 40
    var b = a + 12345
    var c = a * b
    var mod = a % 3

    // 数値演算
    a = 12.3 + 456 + a

    // 文字列演算
    var str = "hogehoge" + "fugafuga"

    // increment
    a++
    b = a++*(c-b--)

    // compare
    c = a == b
    var a2 = a != b
    var a3 = a > b
    var a4 = a < b
    var a5 = a >= b
    var a6 = a <= b
    var a7 = a == b && a != c
    var a8 = a == b && a != c || b > a && b < c

    // minus
    a = -1 + 2
    a = a---(a+b)
}