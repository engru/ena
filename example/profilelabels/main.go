package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"runtime/pprof"
	"sync/atomic"

	_ "net/http/pprof"
)

var i uint32

func update(ctx context.Context) {
	for j := 0; j < 100; j++ {
		atomic.AddUint32(&i, 10)
	}
}

func main() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		labels := pprof.Labels("httppath", "/user")
		pprof.Do(context.TODO(), labels, func(ctx context.Context) {
			i := rand.Int31()
			data, _ := json.Marshal(struct {
				ID uint64
			}{
				ID: uint64(i),
			})
			w.Write(data)
		})
	})
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		labels := pprof.Labels("worker", "purge")
		pprof.Do(context.TODO(), labels, func(ctx context.Context) {
			go update(ctx)
		})

		w.Write([]byte("/message"))
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
