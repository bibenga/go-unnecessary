package main

import (
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

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'

	headers, err := csvReader.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", headers)
	headersMap := map[string]int{}
	for index, name := range headers {
		headersMap[name] = index
	}
	fmt.Printf("%+v\n", headersMap)

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

	f2, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	writer := csv.NewWriter(f2)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, item := range items {
		var row []string = make([]string, len(headers))
		for i, name := range headers {
			if name == "Id" {
				row[i] = strconv.Itoa(item.Id)
			} else if name == "Name" {
				row[i] = item.Name
			} else if name == "Age" {
				row[i] = strconv.FormatFloat(item.Age, 'f', -1, 64)
			}
		}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}

}
