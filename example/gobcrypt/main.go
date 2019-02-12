package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var password = []byte("X123D433yuchao")
var hashedPassword = []byte("$2a$08$yKdOMV5XGnKecSY1P733VOuEEe1bTrbMD3tvXE2C.UBJor77vf9w.")

func bcryptDemo() {
	hashedValue, err := bcrypt.GenerateFromPassword(password, 8)
	if err != nil {
		fmt.Println("GenerateFromPassword error: ", err)
		return
	}
	fmt.Println("hashedValue: ", string(hashedValue))

	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		fmt.Println("CompareHashAndPassword failed, ", err)
		return
	}
	fmt.Println("CompareHashAndPassword ok")
}

func timef(f func(), prefix string) {
	start := time.Now()
	defer func() {
		end := time.Now()
		fmt.Printf("%s: Start[%v], End[%v], Cost[%v]\n", prefix, start, end, end.Sub(start).String())
	}()

	f()
}

func main() {
	timef(bcryptDemo, "bcrypt Demo: ")
}
