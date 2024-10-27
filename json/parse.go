package main

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	l := slog.Default()

	f, err := os.Open("data.json")
	if err != nil {
		l.Error("can't open file", "error", err)
		return
	}
	l.Info("file opened")
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		l.Error("can't read file", "error", err)
		return
	}
	l.Info("json loaded")

	var a = make(map[string]interface{})
	err = json.Unmarshal(bytes, &a)
	if err != nil {
		l.Error("invalid json: %+v", err)
		return
	}

	l.Info("data", "json", a)

	am := a["meta"].(map[string]interface{})
	l.Info("raw", "meta.source", am["source"])
}
