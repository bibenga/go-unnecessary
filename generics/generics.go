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

type Key struct {
	p1 uint32
	p2 uint32
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

	//
	vector := map[[2]uint16]int64{{1, 1}: 1, {2, 2}: 2}
	vector[[2]uint16{3, 3}] = 3
	log.Printf("vector -> %v", vector)
	vector[[2]uint16{3, 3}] = 4
	log.Printf("vector -> %#v", vector)

	//
	vector2 := map[Key]int64{{1, 1}: 1, {2, 2}: 2}
	vector2[Key{3, 3}] = 3
	log.Printf("vector2 -> %v", vector2)
	vector2[Key{3, 3}] = 4
	log.Printf("vector2 -> %#v", vector2)
	delete(vector2, Key{1, 1})
	log.Printf("vector2 -> %#v", vector2)
}
