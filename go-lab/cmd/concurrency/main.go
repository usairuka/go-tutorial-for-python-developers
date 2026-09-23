// 出典: python_to_go_tutorial.md:1622-1700
// 5. goroutine・channel・context・排他制御（45分） / 5.2 完了を待ち、結果を集める

package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

func doubleAll(ctx context.Context, values []int) ([]int, error) {
	results := make(chan int)
	var wg sync.WaitGroup

	for _, value := range values {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			timer := time.NewTimer(20 * time.Millisecond)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			}

			select {
			case <-ctx.Done():
				return
			case results <- n * 2:
			}
		}(value)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	doubled := make([]int, 0, len(values))

	for value := range results {
		doubled = append(doubled, value)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	sort.Ints(doubled)

	return doubled, nil
}

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	values, err := doubleAll(ctx, []int{1, 2, 3})
	if err != nil {
		fmt.Println("failed:", err)
		return
	}

	fmt.Println("results:", values)

	canceled, stop := context.WithCancel(context.Background())
	stop()

	_, err = doubleAll(canceled, []int{1, 2, 3})
	fmt.Println("canceled:", errors.Is(err, context.Canceled))

	// 掲載済みの演習コード（python_to_go_tutorial.md:2160-2164）
	expired, release := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer release()

	_, err = doubleAll(expired, []int{1, 2, 3})
	fmt.Println("deadline:", errors.Is(err, context.DeadlineExceeded))
}
