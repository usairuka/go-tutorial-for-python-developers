// 出典: python_to_go_tutorial.md:2297-2364
// 6. 総合演習：Todo HTTP API（65分） / 6.3 HTTP処理を実装する

package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"example.com/go-lab/internal/task"
)

type Repository interface {
	Add(title string) (task.Task, error)
	List() []task.Task
}

func New(repo Repository) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, repo.List())
	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		defer r.Body.Close()

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input struct {
			Title string `json:"title"`
		}
		if err := decoder.Decode(&input); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			http.Error(w, "body must contain one JSON value", http.StatusBadRequest)
			return
		}

		item, err := repo.Add(input.Title)
		if errors.Is(err, task.ErrEmptyTitle) {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("add task: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusCreated, item)
	})

	// 掲載済みの演習コード（python_to_go_tutorial.md:2667-2669）
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
