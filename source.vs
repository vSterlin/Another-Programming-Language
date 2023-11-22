func getStr() string {
    return "hello"
}

func fib(n int) int {
    if n <= 1 {
        return n
    } else {
        return fib(n-1) + fib(n-2)
    }
}


func main() int {
    fib(10)
    return 0
}