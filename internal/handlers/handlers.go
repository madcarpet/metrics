package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/madcarpet/metrics/internal/storage"
)

const (
	gauge storage.MetricType = 1 + iota
	counter
)

type Repositories interface {
	Update(k string, v float64, t storage.MetricType)
}

// Handler for /update/ path
func Update(w http.ResponseWriter, r *http.Request, s Repositories) {
	if method := r.Method; method != http.MethodPost {
		http.Error(w, "Неверный метод запроса", http.StatusBadRequest)
		return
	}
	// get slice of resuest parameters
	params := strings.Split(r.URL.Path, "/")[1:]
	// check if all parameters sent
	if len(params) != 4 {
		if len(params) == 3 {
			http.Error(w, "Не найдено имя метрики", http.StatusNotFound)
			return
		}
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}
	// check metrics types
	switch params[1] {
	case "gauge":
		v, err := strconv.ParseFloat(params[3], 64)
		if err != nil {
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}
		s.Update(params[2], v, gauge)
	case "counter":
		_, err := strconv.Atoi(params[3])
		if err != nil {
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}
		v, err := strconv.ParseFloat(params[3], 64)
		if err != nil {
			http.Error(w, "Некорректный запрос", http.StatusBadRequest)
			return
		}
		s.Update(params[2], v, counter)
	default:
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
