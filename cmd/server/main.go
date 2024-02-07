package main

import (
	"net/http"

	"github.com/madcarpet/metrics/internal/handlers"
	"github.com/madcarpet/metrics/internal/storage"
)

func run() {
	s := storage.NewMemStorage()
	rs := http.NewServeMux()
	rs.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Update(w, r, s)
	})
	http.ListenAndServe(`localhost:8080`, rs)
}

func main() {
	run()
}
