// 出典: python_to_go_tutorial.md:181-217
// 1. 型・変数・制御構文・関数（35分） / 1.1 まず実行

package main

import "fmt"

func divide(a, b int) (int, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

func main() {
	var count int
	name := "Go"
	const limit = 3
	fmt.Printf("%s count=%d\n", name, count)

	for i := 0; i < limit; i++ {
		fmt.Println(i)
	}

	if value, ok := divide(7, 2); ok {
		fmt.Println("quotient:", value)
	}
	fmt.Println("decimal:", float64(7)/2)

	score := 80
	switch {
	case score >= 80:
		fmt.Println("pass")
	default:
		fmt.Println("retry")
	}

	double := func(n int) int { return n * 2 }
	fmt.Println("double:", double(4))

	// 掲載済みの演習コード（python_to_go_tutorial.md:290-291）
	fmt.Println(sumEven(6))
	fmt.Println(sumEven(0))
}

// 掲載済みの演習コード（python_to_go_tutorial.md:298-306）
func sumEven(limit int) int {
	total := 0
	for n := 1; n <= limit; n++ {
		if n%2 == 0 {
			total += n
		}
	}
	return total
}
