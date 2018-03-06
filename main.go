package main

import (
	"fmt"

	"github.com/lsytj0413/deustgo/atexit"
)

func main() {
	atexit.RegisterHandler(func() {
		fmt.Println("handler 0")
	})

	atexit.HandleInterrupts()
	atexit.Exit(0)
}
