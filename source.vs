func doStuff(fn () => int) int {
   return fn()
}


func main() int {

    doStuff(() int => { return 1 + 1 })
    return 0
}