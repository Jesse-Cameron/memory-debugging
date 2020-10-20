package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
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

func leakyFunction(wg sync.WaitGroup) {
	defer wg.Done()
	s := make([]string, 3)
	for i := 0; i < 10000000; i++ {
		s = append(s, "Surely this won't be expensive 🤔")
		if (i % 100000) == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	fmt.Printf("Starting Server on port :%s\n", "8080")

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requests.doWork()
		r.Body.Close()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Gee whiz. I hope my app isn't leaky 🙏\n"))
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// leaky function
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go leakyFunction(wg)
	// wg.Wait()
}
