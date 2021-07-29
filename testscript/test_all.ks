// 網羅的なスクリプトテスト

import "scenario/ev_main_001.ks"

// label @entrypoint
func main() {

    ev_main_001()

    __dump("dump.log")
}