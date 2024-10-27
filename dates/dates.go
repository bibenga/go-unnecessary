package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	f1 := 0.3333333333333 * 60 * 60
	d := time.Duration(f1) * time.Second
	fmt.Println("d =", d)
	t = t.Add(d)
	fmt.Println("Now? =", t)
	fmt.Println("Weekday =", t.Weekday())
}
