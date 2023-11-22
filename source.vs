
func call(callback () => void) {
    callback()
}

func fib(n int) int {
    if n <= 1 {
        return 1
    } else {
        return fib(n - 1) + fib(n - 2) 
    }
}

// can comment now wow
func main() int {

    x := 1
    y := true
    z := "hello world"


    call(() void => { print(z) })
 

    n := 0
    f := fib(n)
    print(f)
    n = n + 1



    return 0
}