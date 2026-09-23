// 出典: python_to_go_tutorial.md:625-655
// 3. struct・ポインタ・interface・error（40分） / 3.1 値を表すpackageを作る

package task

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyTitle = errors.New("title must not be empty")

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func New(id int, title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("create task: %w", ErrEmptyTitle)
	}
	return Task{ID: id, Title: title}, nil
}

func (t *Task) Complete() {
	t.Done = true
}

func (t Task) Summary() string {
	return fmt.Sprintf("%d: %s (done=%t)", t.ID, t.Title, t.Done)
}

// 掲載済みの演習コード（python_to_go_tutorial.md:1110-1117）
func (t *Task) Rename(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("rename task: %w", ErrEmptyTitle)
	}
	t.Title = title
	return nil
}
