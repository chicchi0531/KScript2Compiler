
func main(){
    var i = 0
    var a = 0
	switch i{

        case 0:
            a=0

        case 1: a=1
            fallthrough
        case 2: a=2
        
        case 3: a=3
            fallthrough
        
        default: fallthrough//error
    }
}