package main

import (
	"net/http"
	"strconv"
	"strings"
)

// type for gauge metric
type gauge float64

// type for counter metric
type counter int64

// type for metrics storage
type MemStorage struct {
	gauge   map[string]gauge
	counter map[string]counter
}

func (m MemStorage) updateGauge(n string, p gauge) {
	m.gauge[n] = p
}

func (m MemStorage) updateCounter(n string, p counter) {
	m.counter[n] += p
}

var base MemStorage

// func to validate request parameters
func checkRequest(w http.ResponseWriter, r *http.Request) (bool, []string) {
	params := []string{}
	if method := r.Method; method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return false, params
	}
	// get slice of parameters
	params = strings.Split(r.URL.Path, "/")[1:]
	// check if all parameters sent
	if len(params) != 4 {
		if len(params) == 3 {
			w.WriteHeader(http.StatusNotFound)
			return false, params
		}
		w.WriteHeader(http.StatusBadRequest)
		return false, params
	}
	// check metrics types
	switch params[1] {
	case "gauge":
		_, err := strconv.ParseFloat(params[3], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return false, params
		}
	case "counter":
		_, err := strconv.Atoi(params[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return false, params
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return false, params
	}
	return true, params[1:]
}

// update path handler
func updateHandler(w http.ResponseWriter, r *http.Request) {
	ok, params := checkRequest(w, r)
	if !ok {
		return
	}
	switch params[0] {
	case "gauge":
		floatGauge, _ := strconv.ParseFloat(params[2], 64)
		base.updateGauge(params[1], gauge(floatGauge))
	case "counter":
		intCounter, _ := strconv.Atoi(params[2])
		base.updateCounter(params[1], counter(intCounter))
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

// function to start server
func run() {
	base = MemStorage{
		gauge:   make(map[string]gauge),
		counter: make(map[string]counter),
	}
	rs := http.NewServeMux()
	rs.HandleFunc("/update/", updateHandler)
	err := http.ListenAndServe(`localhost:8080`, rs)
	if err != nil {
		panic(err)
	}

}

// main function
func main() {
	run()
}
