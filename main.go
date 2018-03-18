package main

import (
	"fmt"
	"sync"

	"github.com/lsytj0413/ena/atexit"
	"github.com/lsytj0413/ena/conc"
)

func main() {
	atexit.RegisterHandler(func() {
		fmt.Println("handler 0")
	})

	atexit.HandleInterrupts()
	// atexit.Exit(0)

	arr := conc.NewConcurrentArray(10)
	l := 10
	var wg sync.WaitGroup
	wg.Add(l)
	for i := 0; i < l; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < 1000; j++ {
				err := arr.Set(uint32(i), j)
				if err != nil {
					fmt.Println(err)
				}
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < l; i++ {
		fmt.Println(arr.Get(uint32(i)))
	}
}
