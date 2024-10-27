package main

import (
	"log"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type Row struct {
	Id     int64   `parquet:"name=id, type=INT64"`
	Name   string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Age    int32   `parquet:"name=age, type=INT32, encoding=PLAIN"`
	Weight float32 `parquet:"name=weight, type=FLOAT"`
}

func main() {
	var err error
	fw, err := local.NewLocalFileWriter("0-olala.parquet")
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}
	pw, err := writer.NewParquetWriter(fw, new(Row), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}
	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.PageSize = 8 * 1024              //8K
	pw.CompressionType = parquet.CompressionCodec_SNAPPY
	num := 100
	for i := 0; i < num; i++ {
		stu := Row{
			Name:   "Olala",
			Age:    int32(20 + i%5),
			Id:     int64(i),
			Weight: float32(50.0 + float32(i)*0.1),
		}
		if err = pw.Write(stu); err != nil {
			log.Println("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
		return
	}
	log.Println("Write Finished")
	fw.Close()
}
