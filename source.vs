func counter(){
    c := 0
    count := () => {
        c = c + 1
    }

    return count
}

c := counter()

c()