func doStuff(fn () => int) int {
   return fn()
}


func main() int {

    x := () int => { return 1 + 1 }
    doStuff(() int => { return 1 + 1 })
    return 0
}