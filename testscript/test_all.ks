// 網羅的なスクリプトテスト

func print(msg string){
    __syscall[0](msg)
}

// label @entrypoint
func main() {

    print("Hello world")

    __dump("dump.log")
}