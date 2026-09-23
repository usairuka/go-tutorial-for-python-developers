// 出典: python_to_go_tutorial.md:2596-2642
// 6. 総合演習：Todo HTTP API（65分） / 6.6 HTTPと並行アクセスをテストする

package task_test

import (
	"sync"
	"testing"

	"example.com/go-lab/internal/task"
)

func TestMemoryConcurrentAdd(t *testing.T) {
	var store task.Memory

	const count = 20
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			if _, err := store.Add("Learn Go"); err != nil {
				t.Error(err)
			}
			_ = store.List()
		}()
	}

	wg.Wait()

	items := store.List()
	if len(items) != count {
		t.Fatalf("count=%d, want %d", len(items), count)
	}

	seen := make(map[int]bool)
	for _, item := range items {
		if item.ID < 1 || item.ID > count || seen[item.ID] {
			t.Fatalf("invalid or duplicate ID: %d", item.ID)
		}
		seen[item.ID] = true
	}

	items[0].Title = "changed outside the store"
	if store.List()[0].Title != "Learn Go" {
		t.Fatal("List must return a copy")
	}
}
