// IF文テスト

// jmp @entrypoint

//---------------
// 0: msg
//---------------
// label @print
// pushint 0
// pushvalue
// pushint 1
// pushint 0
// syscall
// pop
// return
func print(msg string){
    __syscall[0](msg)
}

// label @entrypoint
func main() {

    // 制御文テスト
    // 0: a
    // 1: b
    //-----------
    // pushint 100
    // pushint 0
    // popvalue
    // pushint 200
    // pushint 1
    // popvalue
    // pushint 0    [a == b]
    // pushvalue
    // pushint 1
    // pushvalue
    // equ
    // jze @L1      [if a == b]
    // pushstring "a equals b"
    // pushint 1
    // call @print
    // jmp @L3      [end if]
    // label @L1    [else]
    // pushint 0    [a <= b]
    // pushvalue
    // pushint 1
    // pushvalue
    // le
    // jze @L2      [if a <= b]
    // pushint 0    [a < b]
    // pushvalue
    // pushint 1
    // pushvalue
    // lt
    // jze @L2      [if a < b] nest
    // pushstring "a less than b"
    // pushint 1
    // call @print
    // label @L2    [end if]
    // jmp @L4      [end if]
    // label @L3    [else]
    // pushstring "a greater than b"
    // pushint 1
    // call @print
    // label @L4
    // label @L0    [end if]
    var a int = 100
    var b int = 200
    if a==b {
        print("a equals b")
    } else if a<=b {
        if a < b{
            print("a less than b")
        }
    } else {
        print("a greater than b")
    }

    __dump("dump.log")
}