// 出典: python_to_go_tutorial.md:983-995
// 3. struct・ポインタ・interface・error（40分） / 3.3 errorとdefer / defer

package main

import "fmt"

func main() {
	n := 1

	defer fmt.Println("deferred:", n)
	defer fmt.Println("last registered")

	n = 2
	fmt.Println("now:", n)
}
