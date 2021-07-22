// 網羅的なスクリプトテスト
import "scenario/ev_main_001.ks"

func TestFunc(a int, b int) int {
    val1 := 0
    return val1 + a + b
}

func main() {

    // ノベルスタート
    ev_main_001()
    __dump("dump.log")
}