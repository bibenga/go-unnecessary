package main

import (
	"testing"
)

func Test1and2(t *testing.T) {
	t.Run("check 1", func(t *testing.T) {
		got := 1
		if got != 1 {
			t.Errorf("Abs(-1) = %d; want 1", got)
		}
	})
	t.Run("check 2", func(t *testing.T) {
	})
}
