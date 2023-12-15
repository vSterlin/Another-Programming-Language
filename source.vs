func counter() () => int {
    count := 0
    return () => { return count++ }
}

c := counter()

print(c(), c())