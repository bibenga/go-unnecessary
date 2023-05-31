package main

// https://go.dev/doc/tutorial/generics

import (
	"log"
)

type S1 struct {
	V1 string
}

func (m *S1) Debug() {
	log.Printf("Debug: s1 -> %+v", m)
}

type S2 struct {
	S1
	V2 string
}

type JopaV1 struct {
	V1 string
}

// func (m *S2) Debug() {
// 	log.Printf("Debug: s2 -> %+v", m)
// }

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	s1 := S1{"Unnecessary1"}
	log.Printf("s1 -> %+v", s1)
	s1.Debug()

	s2 := S2{S1{"Unnecessary1"}, "Unnecessary2"}
	log.Printf("s1 -> %+v", s2)
	s2.Debug()
}
