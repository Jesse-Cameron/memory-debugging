package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

var (
	requests RequestTracker
)

type RequestTracker struct {
	sync.Mutex
	requests [][]byte
}

func (rt *RequestTracker) doWork() {
	rt.Lock()
	defer rt.Unlock()
	// alloc 1MB for each request
	rt.requests = append(rt.requests, bytes.Repeat([]byte("a"), 1000000))
	PrintMemUsage()
}

func main() {
	fmt.Printf("Starting Server on port :%s 🎉🎉🎉\n", "8080")

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requests.doWork()
		r.Body.Close()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Gee whiz. I hope my app isn't leaky 🙏\n"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
