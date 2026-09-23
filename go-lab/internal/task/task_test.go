// 出典: python_to_go_tutorial.md:1200-1259
// 4. package設計・テスト・ジェネリクス（35分） / 4.2 テーブル駆動テスト

package task_test

import (
	"errors"
	"testing"

	"example.com/go-lab/internal/task"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantTitle string
		wantErr   error
	}{
		{name: "valid", input: "Learn Go", wantTitle: "Learn Go"},
		{name: "trim", input: "  Learn Go  ", wantTitle: "Learn Go"},
		{name: "empty", input: "", wantErr: task.ErrEmptyTitle},
		{name: "whitespace", input: " \t\n", wantErr: task.ErrEmptyTitle},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := task.New(1, tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}

			if tc.wantErr != nil {
				return
			}

			if got.Title != tc.wantTitle || got.ID != 1 || got.Done {
				t.Errorf("unexpected task: %+v", got)
			}
		})
	}
}

func TestComplete(t *testing.T) {
	original := task.Task{
		ID:    1,
		Title: "Learn Go",
	}

	copied := original
	copied.Complete()

	if original.Done || !copied.Done {
		t.Fatal("completing a copy must not change the original")
	}

	original.Complete()

	if !original.Done {
		t.Fatal("Complete must set Done")
	}
}

// 掲載済みの演習コード（python_to_go_tutorial.md:1528-1542）
func TestRename(t *testing.T) {
	item := task.Task{ID: 1, Title: "Before"}
	if err := item.Rename("  After  "); err != nil {
		t.Fatal(err)
	}
	if item.Title != "After" {
		t.Fatalf("title = %q, want After", item.Title)
	}
	if err := item.Rename("  "); !errors.Is(err, task.ErrEmptyTitle) {
		t.Fatalf("error = %v, want ErrEmptyTitle", err)
	}
	if item.Title != "After" {
		t.Fatalf("invalid rename changed title to %q", item.Title)
	}
}
