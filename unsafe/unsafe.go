package main

import (
	"log"
	"unsafe"
)

type A0 struct {
}

type A1 struct {
	A11 int32
}

// go:orderedfields
type A2 struct {
	A21 int16
	A22 int32
	A23 bool
	A24 bool
	A25 int64
}

type A3 struct {
	A23 bool `json:"id"`
	A24 bool
	A21 int16
	A22 int32
	A25 int64
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	log.Printf("A0 -> %d", unsafe.Sizeof(A0{}))
	log.Printf("A1 -> %d", unsafe.Sizeof(A1{}))
	log.Printf("A2 -> %d", unsafe.Sizeof(A2{}))
	log.Printf("A3 -> %d", unsafe.Sizeof(A3{}))
}
