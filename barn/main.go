package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/bibenga/barn-go"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	pgxslog "github.com/mcosta74/pgx-slog"
)

// const driver string = "sqlite3"
// const dsn string = "file:barn/_barn.db?cache=shared&mode=rwc&_journal_mode=WAL&_loc=UTC"

const driver string = "pgx"
const dsn string = "host=host.docker.internal port=5432 user=rds password=sqlsql dbname=barn TimeZone=UTC sslmode=disable"

func main() {
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.SetPrefix("")

	// os.Remove("./barn/_barn.db")

	// db, err := sql.Open("sqlite3", "./barn/_barn.db")
	// db, err := sql.Open("sqlite3", "file:barn/_barn.db?cache=shared&mode=rwc&_journal_mode=WAL&_loc=UTC")

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		slog.Error("db config errorz", "error", err)
		panic(err)
	}
	connConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxslog.NewLogger(slog.Default()),
		LogLevel: tracelog.LogLevelDebug,
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sql.Open(driver, connStr)

	// db, err := sql.Open(driver, dsn)
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

	manager := barn.NewLockManager(db)
	err = manager.InitializeDB()
	if err != nil {
		slog.Error("db error", "error", err)
		panic(err)
	}
	// go manager.Run()

	s := <-osSignal
	slog.Info("os signal received", "signal", s)

	// manager.Stop()
	// scheduler.Stop()
}
