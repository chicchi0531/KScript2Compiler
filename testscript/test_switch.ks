// 網羅的なスクリプトテスト

// jmp @entrypoint

// label @entrypoint
func main() {

    //------------
    // 0:a
    // 1:b
    //------------
    // pushint 100
    // pushint 0
    // popvalue
    var a int = 100
    var b string

    // pushint 0
    // pushvalue
    // pushint 1
    // equ
    // jnz @L1   [jump to case 1]
    // pushint 0
    // pushvalue
    // pushint 2
    // equ
    // jnz @L2   [jump to case 2]
    // pushint 0
    // pushvalue
    // pushint 3
    // equ
    // jnz @L3   [jump to case 3]
    // jmp @L4   [jump to default]
    // label @L1
    // pushstring "apple"
    // pushint 1
    // popvalue
    // jmp @L0
    // label @L2
    // jmp @L3   [fallthrough]
    // jmp @L0
    // label @L3
    // pushstring "orange"
    // pushint 1
    // popvalue
    // jmp @L0
    // label @L4
    // pushstring "grape"
    // pushint 1
    // popvalue
    // label @L0
    switch a{
        case 1: b = "apple"
        case 2: fallthrough
        case 3: b = "orange"
        default: b = "grape"
    }

    __dump("dump.log")
}