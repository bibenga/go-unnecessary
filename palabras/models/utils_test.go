package models

import (
	"reflect"
	"testing"
)

func TestFilterEmptyString(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		src := []string{"", "a", "", "b"}
		res := FilterEmptyString(src)
		if !reflect.DeepEqual(res, []string{"a", "b"}) {
			t.Errorf("FilterEmptyString failed")
		}
	})
}

func TestSplitLines(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		res := SplitLines("a\n\nb\n")
		if !reflect.DeepEqual(res, []string{"a", "b"}) {
			t.Errorf("SplitLines failed: %+v", res)
		}
	})
}

func TestSplitWords(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		res := SplitWords("más  ¿o?   menos")
		if !reflect.DeepEqual(res, []string{"más", "o", "menos"}) {
			t.Errorf("SplitWords failed: %+v", res)
		}
	})
}

func TestUnidecode(t *testing.T) {
	t.Run("es", func(t *testing.T) {
		res := Unidecode("más")
		if res != "mas" {
			t.Errorf("Unidecode failed: %+v", res)
		}
	})

	t.Run("ru", func(t *testing.T) {
		res := Unidecode("привет")
		if res != "privet" {
			t.Errorf("Unidecode failed: %+v", res)
		}
	})

	t.Run("fr", func(t *testing.T) {
		res := Unidecode("réveillez-vous")
		if res != "reveillez-vous" {
			t.Errorf("Unidecode failed: %+v", res)
		}
	})
}
