// 出典: python_to_go_tutorial.md:2478-2586
// 6. 総合演習：Todo HTTP API（65分） / 6.6 HTTPと並行アクセスをテストする

package httpapi_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/go-lab/internal/httpapi"
	"example.com/go-lab/internal/task"
)

func TestCreateAndList(t *testing.T) {
	handler := httpapi.New(&task.Memory{})

	empty := httptest.NewRecorder()
	handler.ServeHTTP(empty, httptest.NewRequest(http.MethodGet, "/tasks", nil))
	if empty.Code != http.StatusOK || strings.TrimSpace(empty.Body.String()) != "[]" {
		t.Fatalf("empty list: status=%d body=%s", empty.Code, empty.Body.String())
	}

	created := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"  Learn Go  "}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(created, request)
	if created.Code != http.StatusCreated {
		t.Fatalf("create: status=%d body=%s", created.Code, created.Body.String())
	}

	var item task.Task
	if err := json.NewDecoder(created.Body).Decode(&item); err != nil {
		t.Fatal(err)
	}
	if item.ID != 1 || item.Title != "Learn Go" || item.Done {
		t.Fatalf("unexpected created task: %+v", item)
	}

	listed := httptest.NewRecorder()
	handler.ServeHTTP(listed, httptest.NewRequest(http.MethodGet, "/tasks", nil))
	if listed.Code != http.StatusOK {
		t.Fatalf("list: status=%d", listed.Code)
	}

	var items []task.Task
	if err := json.NewDecoder(listed.Body).Decode(&items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0] != item {
		t.Fatalf("unexpected list: %+v", items)
	}
}

func TestRejectInvalidInput(t *testing.T) {
	for _, body := range []string{
		`{"title":" "}`,
		`{"title":`,
		`{"title":"Go","extra":true}`,
		`{"title":"Go"} {"title":"another"}`,
	} {
		t.Run(body, func(t *testing.T) {
			store := &task.Memory{}
			handler := httpapi.New(store)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
			handler.ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
			}
			if len(store.List()) != 0 {
				t.Fatal("invalid request must not create a task")
			}
		})
	}
}

var errStoreUnavailable = errors.New("store unavailable: internal detail")

type failingRepo struct{}

func (failingRepo) Add(string) (task.Task, error) {
	return task.Task{}, errStoreUnavailable
}

func (failingRepo) List() []task.Task {
	return []task.Task{}
}

func TestCreateUnexpectedError(t *testing.T) {
	handler := httpapi.New(failingRepo{})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Learn Go"}`))
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	body := response.Body.String()
	if !strings.Contains(body, "internal server error") {
		t.Fatalf("body should contain public error message, got %q", body)
	}
	if strings.Contains(body, errStoreUnavailable.Error()) {
		t.Fatalf("body must not expose internal error detail: %q", body)
	}
}

// 掲載済みの演習コード（python_to_go_tutorial.md:2680-2697）
func TestHealth(t *testing.T) {
	handler := httpapi.New(&task.Memory{})

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status=%d", response.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body=%v", body)
	}
}
