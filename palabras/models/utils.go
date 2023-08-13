package models

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

func FilterEmptyString(items []string) []string {
	var res []string
	for _, item := range items {
		item = strings.Trim(item, " ")
		if item != "" {
			res = append(res, item)
		}
	}
	return res
}

var linesRe = regexp.MustCompile("[\n\r]+")

func SplitLines(text string) []string {
	lines := linesRe.Split(text, -1)
	return FilterEmptyString(lines)
}

var wordsRe = regexp.MustCompile("[ \r\n¡!¿?.,:;'\"]+")

func SplitWords(line string) []string {
	words := wordsRe.Split(line, -1)
	return FilterEmptyString(words)
}

func Unidecode(src string) string {
	return unidecode.Unidecode(src)
}
