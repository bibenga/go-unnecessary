package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"unnecessary/barn/barn"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

// const driver string = "sqlite3"
// const dsn string = "file:barn/_barn.db?cache=shared&mode=rwc&_journal_mode=WAL&_loc=UTC"

const driver string = "pgx"
const dsn string = "host=host.docker.internal port=5432 user=rds password=sqlsql dbname=barn TimeZone=UTC sslmode=disable"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	slog.Info("Hello")

	// os.Remove("./barn/_barn.db")

	// db, err := sql.Open("sqlite3", "./barn/_barn.db")
	// db, err := sql.Open("sqlite3", "file:barn/_barn.db?cache=shared&mode=rwc&_journal_mode=WAL&_loc=UTC")
	db, err := sql.Open(driver, dsn)
	if err != nil {
		slog.Error("db error", "error", err)
		panic(err)
	}
	defer db.Close()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	scheduler := barn.NewScheduler(db)
	err = scheduler.InitializeDB()
	if err != nil {
		slog.Error("db error", "error", err)
		panic(err)
	}

	err = scheduler.DeleteAll()
	if err != nil {
		slog.Error("db error", "error", err)
		panic(err)
	}

	cron1 := "*/5 * * * * *"
	err = scheduler.Add("olala1", &cron1, nil, "{\"type\":\"olala1\"}")
	if err != nil {
		slog.Error("db error", "error", err)
		panic(err)
	}

	nextTs2 := time.Now().UTC().Add(-20 * time.Second)
	// err = scheduler.Add("olala2", nil, &nextTs2)
	cron2 := "*/10 * * * * *"
	err = scheduler.Add("olala2", &cron2, &nextTs2, "{\"type\":\"olala2\"}")
	if err != nil {
		slog.Error("db error", "error", err)
		panic(err)
	}
	go scheduler.Run()

	// manager := barn.NewLockManager(db)
	// err = manager.InitializeDB()
	// if err != nil {
	// 	slog.Error("db error", "error", err)
	// 	panic(err)
	// }
	// go manager.Run()

	s := <-osSignal
	slog.Info("os signal received", "signal", s)

	// manager.Stop()
	scheduler.Stop()
}
