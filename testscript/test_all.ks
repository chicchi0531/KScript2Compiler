// 網羅的なスクリプトテスト

// jmp @entrypoint

// strcut test
type User struct{
    id float
    log [8]string
}
type Room struct{
    id int
    users [3]User
}

func MakeRoom() Room{
    var r Room
    r.id = 1
    for u:=0; u<3; u++{
        r.users[u].id = 1.0
        for l:=0; l<8; l++{
            r.users[u].log[l] = "default log."
        }
    }
    return r
}

// label @entrypoint
func main() {

    var room Room
    var room2 Room
    var fval float = 123.5

    // 1 :pushint 123
    // 2 :pushint 0 [addr:room]
    // 3 :pushint 1 [addr:Room.id]
    // 4 :add       [addr:room.id]
    // 5 :popvalue
    room.id = 123

    // 6 :pushint 10
    // 7 :pushint 0 [addr:room]
    // 8 :pushint 2 [addr:Room.users]
    // 9 :add       [addr:room.users]
    // 10:pushint 1 [[1]]
    // 11:pushint 10 [sizeof(User)]
    // 12:mul       [addr:User[1]]
    // 13:add       [addr:room.users[1]]
    // 14:pushint 1 [addr:User.id]
    // 15:add       [addr:room.users[1].id]
    // 16:popvalue
    room.users[1].id = 20.5 + fval

    // 17:pushint 0
    // 18:pushvalue 32
    // 19:pushint 32
    // 20:popvalue 32
    room2 = MakeRoom()

    __dump("dump.log")
}