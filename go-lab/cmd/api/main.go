// 出典: python_to_go_tutorial.md:2384-2412
// 6. 総合演習：Todo HTTP API（65分） / 6.4 サーバーを起動する

package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"example.com/go-lab/internal/httpapi"
	"example.com/go-lab/internal/task"
)

func main() {
	store := &task.Memory{}

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           httpapi.New(store),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on http://127.0.0.1:8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
