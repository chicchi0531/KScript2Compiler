// 網羅的なスクリプトテスト

// jmp @entrypoint

// strcut test
type User struct{
    id int
    log [8]string
}
type Room struct{
    id int
    users [3]User
}

// label @entrypoint
func main() {

    __dump("dump.log")
}