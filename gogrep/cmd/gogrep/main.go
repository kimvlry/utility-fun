package main

import (
	"fmt"
	"unsafe"
)

type cringe struct {
	b  byte
	i  int64
	bb byte
	ii int64
}

type cool struct {
	i  int64
	ii int64
	b  byte
	bb byte
}

func main() {
	fmt.Println(unsafe.Sizeof(cringe{}))
	fmt.Println(unsafe.Sizeof(cool{}))
}
