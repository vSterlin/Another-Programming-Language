func doStuff(fn () => int) int {
   return fn()
}


func main() int {
    x := () string => { return "woa" }
    doStuff(() int => { return 1 + 1 })
    return 0
}