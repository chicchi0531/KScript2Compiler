// for loop test

// label @entrypoint
func main() {

    sum := 0
    // normal type for
    for i:=0; i<10; i++{
        if sum == 5{
            sum++
            continue
        }
        sum += i
    }

    sum = 0
    // while type for
    for sum < 10{
        if sum == 5{
            break
        }
        sum++
    }

    //以下はerrortest
    //continue
    //break
