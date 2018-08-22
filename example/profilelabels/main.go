package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
	"time"
)

func main() {
	ctx := context.Background()
	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 1000000000000; i++ {
			labels := pprof.Labels("handler", "hello")
			pprof.Do(ctx, labels, func(ctx context.Context) {
				generate(1, 10)
			})
		}
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 1000000000000; i++ {
			labels := pprof.Labels("handler2", "hello")
			pprof.Do(ctx, labels, func(ctx context.Context) {
				generate(1, 10)
			})
		}
	}()

	log.Fatal(http.ListenAndServe("localhost:5555", nil))
}

func generate(duration int, usage int) { /* string concats in loops */ }
