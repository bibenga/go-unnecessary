package main

// https://go.dev/doc/tutorial/generics

import (
	"log"
)

func SumIntsOrFloats[K comparable, V int64 | float32 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	ints := map[string]int64{
		"first":  1,
		"second": 2,
	}
	ints["thurd"] = 3
	log.Printf("ints -> %v", SumIntsOrFloats(ints))

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  1.1,
		"second": 2.2,
	}
	log.Printf("floats -> %v", SumIntsOrFloats(floats))
}
