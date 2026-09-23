// 出典: python_to_go_tutorial.md:1396-1413
// 4. package設計・テスト・ジェネリクス（35分） / 4.3 ジェネリクスを1つ使う

package main

import "fmt"

func contains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func main() {
	fmt.Println(contains([]int{1, 2, 3}, 2))
	fmt.Println(contains([]string{"Go", "Python"}, "Rust"))
}
