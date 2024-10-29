package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Item struct {
	Id   int
	Name string
	Age  float64
}

func main() {
	f, err := os.Open("csv/csv.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	bom := [3]byte{0xEF, 0xBB, 0xBF}
	utf8sigHeader := make([]byte, 3)
	n, err := f.Read(utf8sigHeader)
	if err != nil {
		panic(err)
	}
	if n != 3 {
		panic("n != 3")
	}
	if utf8sigHeader[0] == bom[0] && utf8sigHeader[1] == bom[1] && utf8sigHeader[2] == bom[2] {
		fmt.Println("hay utf8sigHeader")
	} else {
		fmt.Println("no hay utf8sigHeader")
		f.Seek(0, 0)
	}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	headers, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", headers)
	headersMap := map[string]int{}
	for index, name := range headers {
		headersMap[name] = index
	}
	fmt.Printf("%#v\n", headersMap)

	var items []*Item
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		id, err := strconv.Atoi(rec[headersMap["Id"]])
		if err != nil {
			panic(err)
		}

		name := rec[headersMap["Name"]]

		ageS := rec[headersMap["Age"]]
		age, err := strconv.ParseFloat(ageS, 64)
		if err != nil {
			age, err = strconv.ParseFloat(strings.Replace(ageS, ",", ".", 1), 64)
			if err != nil {
				panic(err)
			}
		}

		item := &Item{Id: id, Name: name, Age: age}
		fmt.Printf("%#v\n", item)
		items = append(items, item)
	}
	fmt.Printf("%#v\n", items)
	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

	f2, err := os.Create("csv/output.csv")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	// writer := bufio.NewWriter(f2)
	if _, err := f2.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(f2)
	csvWriter.Comma = ';'
	csvWriter.UseCRLF = true
	defer csvWriter.Flush()

	if err := csvWriter.Write(headers); err != nil {
		panic(err)
	}

	for _, item := range items {
		var row []string = make([]string, len(headers))
		for i, name := range headers {
			switch name {
			case "Id":
				row[i] = strconv.Itoa(item.Id)
			case "Name":
				row[i] = item.Name
			case "Age":
				row[i] = strconv.FormatFloat(item.Age, 'f', -1, 64)
			default:
				panic(fmt.Sprintf("Unknown name: %s", name))
			}
		}
		if err := csvWriter.Write(row); err != nil {
			panic(err)
		}
	}

}
