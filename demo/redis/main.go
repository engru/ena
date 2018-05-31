package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var _ = redis.Dial

func main() {
	fmt.Println("main")
}
