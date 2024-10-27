package main

import (
	"fmt"
	"os"

	"github.com/tobgu/qframe"
	"github.com/tobgu/qframe/config/csv"
)

func main() {
	f, err := os.Open("csv/csv.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	df := qframe.ReadCSV(f, csv.Delimiter(';'))
	fmt.Println(df)
}
