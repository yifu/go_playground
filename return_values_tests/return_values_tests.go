package main

import(
"fmt"
)

func main() {
     // 1- Of course returning one value is accepted by the compiler.
     fmt.Println("return one value =", return_one_value())

     // 2- A function returning multiple values is not composable with
     // a function accepting only one value.
     //fmt.Println("return two values =", return_two_values())

     // 3- But! The compiler is ok with map lookups, which return two values: elt,ok.
     m := map[int]int{0:1, 1:2, 2:3}
     fmt.Println("m[1]=", m[1])

     elt, ok := m[1]
     fmt.Println("elt=", elt, ", ok=", ok)

     // 4- A function accepting two values is composable with a function returning two values. Nice!
     accept_two_values(return_two_values())
}

func return_one_value() int {
     return 1
}

func return_two_values() (int,int) {
     return 1,2
}

func accept_two_values(a,b int) {
     fmt.Println("a=", a, ", b=", b)
}