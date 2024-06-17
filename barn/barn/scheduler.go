package barn

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/adhocore/gronx"
)

type Entry struct {
	Id       int32
	Name     string
	Cron     *string
	IsActive bool
	NextTs   *time.Time
	LastTs   *time.Time
	Message  *string
}

type EntryMap map[int32]*Entry

type Scheduler struct {
	db      *sql.DB
	entries EntryMap
	stop    chan struct{}
	stopped chan struct{}
	timer   *time.Timer
	entry   *Entry
}

func NewScheduler(db *sql.DB) *Scheduler {
	manager := Scheduler{
		db:      db,
		entries: nil,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	return &manager
}

func (scheduler *Scheduler) InitializeDB() error {
	db := scheduler.db
	slog.Info("create table", "table", "barn_entry")
	sqlStmt := `
	CREATE TABLE  IF NOT EXISTS "barn_entry" (
        id INTEGER NOT NULL, 
        name VARCHAR NOT NULL, 
        is_active BOOLEAN DEFAULT True NOT NULL, 
        cron VARCHAR, 
        next_ts TIMESTAMP, 
        last_ts TIMESTAMP, 
        message JSON, 
        PRIMARY KEY (id), 
        UNIQUE (name)
	)`
	_, err := db.Exec(sqlStmt)
	return err
}

func (scheduler *Scheduler) Stop() {
	slog.Info("stopping")
	scheduler.stop <- struct{}{}
	<-scheduler.stopped
	close(scheduler.stop)
	close(scheduler.stopped)
	slog.Info("stopped")
}

func (scheduler *Scheduler) Run() {
	db := scheduler.db
	stmt, err := db.Prepare(`select id, name, is_active, cron, next_ts, last_ts, message from barn_entry`)
	if err != nil {
		slog.Error("db", "error", err)
		panic(err)
	}
	defer stmt.Close()

	slog.Info("started")

	err = scheduler.reload()
	if err != nil {
		slog.Error("db", "error", err)
		panic(err)
	}

	reloader := time.NewTicker(5 * time.Second)
	defer reloader.Stop()

	// scheduler.timer = time.NewTimer(1 * time.Second)
	// scheduler.scheduleNext()

	for {
		// timer := scheduler.timer
		entry := scheduler.getNext()
		var timer *time.Timer
		if entry == nil {
			timer = time.NewTimer(1 * time.Second)
		} else {
			timer = time.NewTimer(time.Until(*entry.NextTs))
		}
		defer timer.Stop()

		select {
		case <-scheduler.stop:
			slog.Info("terminate")
			scheduler.stopped <- struct{}{}
			return
		case <-timer.C:
			id := -1
			// entry := scheduler.entry
			if entry != nil {
				id = int(entry.Id)
				if entry.Cron != nil {
					nextTs, err := gronx.NextTick(*entry.Cron, false)
					if err != nil {
						panic(err)
					}
					entry.LastTs = entry.NextTs
					entry.NextTs = &nextTs
				}
			}
			slog.Info("tik", "entry", id)
		case <-reloader.C:
			timer.Stop()
			// err = scheduler.reload()
			// if err != nil {
			// 	slog.Error("db", "error", err)
			// 	panic(err)
			// }
		}
	}
}

func (scheduler *Scheduler) reload() error {
	entries, err := scheduler.getEntries()
	if err != nil {
		return err
	}
	scheduler.entries = entries
	return nil
}

func (scheduler *Scheduler) getNext() *Entry {
	var next *Entry = nil
	for _, entry := range scheduler.entries {
		if next == nil {
			next = entry
			// slog.Info("=> ", "next", next.NextTs)
		} else {
			// slog.Info("=> ", "next", next.NextTs, "entry", entry.NextTs)
			if entry.NextTs.Before(*next.NextTs) {
				next = entry
			}
		}
	}
	return next
}

func (scheduler *Scheduler) scheduleNext() {
	var next *Entry = nil
	for _, entry := range scheduler.entries {
		if next == nil {
			next = entry
			// slog.Info("=> ", "next", next.NextTs)
		} else {
			// slog.Info("=> ", "next", next.NextTs, "entry", entry.NextTs)
			if entry.NextTs.Before(*next.NextTs) {
				next = entry
			}
		}
	}

	if next != nil {
		// scheduler.timer.Reset(time.Since(*next.NextTs))
		// scheduler.entry = next
		// slog.Info("next", "entry", next.Id)
	} else {
		// scheduler.timer.Reset(1 * time.Second)
		// scheduler.entry = nil
		// slog.Info("next", "entry", nil)
	}
}

func (scheduler *Scheduler) getEntries() (EntryMap, error) {
	db := scheduler.db
	stmt, err := db.Prepare(`select id, name, is_active, cron, next_ts, last_ts, message from barn_entry`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	entries := make(EntryMap)

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e Entry = Entry{}
		err := rows.Scan(&e.Id, &e.Name, &e.IsActive, &e.Cron, &e.NextTs, &e.LastTs, &e.Message)
		if err != nil {
			return nil, err
		}
		slog.Info("loaded", "entry", e)
		if e.IsActive {
			entries[e.Id] = &e
		}
	}
	return entries, nil
}

func (scheduler *Scheduler) Add(name string, cron *string, nextTs *time.Time) error {
	// fake 1
	// cron := "*/5 * * * * *"

	if cron == nil && nextTs == nil {
		return fmt.Errorf("invalid args")
	}
	if cron != nil && nextTs == nil {
		nextTs2, err := gronx.NextTick(*cron, true)
		if err != nil {
			return err
		}
		nextTs = &nextTs2
	}

	// fake 2
	var message *string = nil
	if name == "olala1" {
		var m = make(map[string]interface{})
		m["extra"] = 1
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}
		m2 := string(b)
		message = &m2
	}
	db := scheduler.db
	if message != nil {
		slog.Info("create the entry", "name", name, "cron", cron, "message", *message)
	} else {
		slog.Info("create the entry", "name", name, "cron", cron, "message", message)
	}

	stmt, err := db.Prepare(
		`insert into barn_entry(name, cron, next_ts, message) 
		values (?, ?, ?, ?) 
		returning id`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, cron, nextTs, message).Scan(&id)
	if err != nil {
		return err
	}
	slog.Info("the entry is created", "name", name, "id", id)
	return nil
}

func (scheduler *Scheduler) Delete(id int) error {
	db := scheduler.db

	slog.Info("delete the entry", "id", id)
	res, err := db.Exec(
		`delete from barn_entry where id=?`,
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("an object does not exist")
	}
	return nil
}

func (scheduler *Scheduler) update(entry *Entry) error {
	db := scheduler.db

	slog.Info("update the entry", "entry", entry)
	res, err := db.Exec(
		`update barn_entry set next_ts=?, last_ts=? where id=?`,
		entry.NextTs, entry.LastTs, entry.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("an object deleted")
	}
	return nil
}

func (scheduler *Scheduler) DeleteByName(name string) error {
	return nil
}
