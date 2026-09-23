// 出典: python_to_go_tutorial.md:2245-2277
// 6. 総合演習：Todo HTTP API（65分） / 6.2 保存先を実装する

package task

import "sync"

// Memory is a concurrency-safe, in-memory task store.
// Its zero value is ready to use. Do not copy it after first use.
type Memory struct {
	mu     sync.Mutex
	nextID int
	items  []Task
}

func (m *Memory) Add(title string) (Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, err := New(m.nextID+1, title)
	if err != nil {
		return Task{}, err
	}
	m.nextID++
	m.items = append(m.items, item)
	return item, nil
}

func (m *Memory) List() []Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]Task, len(m.items))
	copy(result, m.items)
	return result
}
