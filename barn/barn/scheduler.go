package barn

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/adhocore/gronx"
)

type Entry struct {
	Id       int32
	Name     string
	IsActive bool
	Cron     *string
	NextTs   *time.Time
	LastTs   *time.Time
	Message  *string
}

func (e Entry) LogValue() slog.Value {
	// return slog.AnyValue(computeExpensiveValue(e.arg))
	var args []slog.Attr
	args = append(args, slog.Int("Id", int(e.Id)))
	args = append(args, slog.Bool("IsActive", e.IsActive))
	if e.Cron != nil {
		args = append(args, slog.String("Cron", *e.Cron))
	}
	if e.NextTs != nil {
		args = append(args, slog.String("NextTs", e.NextTs.String()))
	}
	return slog.GroupValue(
		// slog.Int("Id", int(e.Id)),
		// // slog.String("Name", e.Name),
		// slog.Bool("IsActive", e.IsActive),
		// // slog.Any("Cron", e.Cron),
		// // slog.Any("NextTs", e.NextTs),
		// // slog.Any("Message", e.Message),
		args...,
	)
}

type EntryMap map[int32]*Entry

type Scheduler struct {
	entries EntryMap
	db      *sql.DB
	stop    chan struct{}
	stopped chan struct{}
	timer   *time.Timer
	entry   *Entry
}

func NewScheduler(db *sql.DB) *Scheduler {
	logger := slog.Default().With("a", 1)
	logger.Info("created")
	manager := Scheduler{
		entries: make(EntryMap),
		db:      db,
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
	// db := scheduler.db
	// stmt, err := db.Prepare(
	// 	`select id, name, is_active, cron, next_ts, last_ts, message
	// 	from barn_entry`,
	// )
	// if err != nil {
	// 	slog.Error("db", "error", err)
	// 	panic(err)
	// }
	// defer stmt.Close()

	slog.Info("started")

	err := scheduler.reload()
	if err != nil {
		slog.Error("db", "error", err)
		panic(err)
	}

	reloader := time.NewTicker(5 * time.Second)
	defer reloader.Stop()

	// I don't know how to be
	scheduler.timer = time.NewTimer(1 * time.Second)
	defer scheduler.timer.Stop()
	select {
	case <-scheduler.timer.C:
	default:
	}
	scheduler.scheduleNext()

	for {
		select {
		case <-scheduler.stop:
			slog.Info("terminate")
			scheduler.stopped <- struct{}{}
			return
		case <-scheduler.timer.C:
			entry := scheduler.entry
			if entry != nil {
				// process
				slog.Info("tik ", "entry", entry.Id, "nextTs", entry.NextTs)
				// calculate next time
				if entry.Cron != nil {
					nextTs, err := gronx.NextTick(*entry.Cron, false)
					if err != nil {
						panic(err)
					}
					entry.LastTs = entry.NextTs
					entry.NextTs = &nextTs
				} else {
					entry.IsActive = false
				}
				err = scheduler.update(entry)
				if err != nil {
					slog.Error("db", "error", err)
					panic(err)
				}
			}
			scheduler.scheduleNext()
		case <-reloader.C:
			err = scheduler.reload()
			if err != nil {
				slog.Error("db", "error", err)
				panic(err)
			}
		}
	}
}

func (scheduler *Scheduler) reload() error {
	entries, err := scheduler.getEntries()
	if err != nil {
		return err
	}

	for id, newEntry := range entries {
		if oldEntry, ok := scheduler.entries[id]; ok {
			// exists
			if oldEntry.NextTs != newEntry.NextTs || oldEntry.Cron != newEntry.Cron {
				// changed
				slog.Info("changed entry", "entry", newEntry.Id)
				oldEntry.Cron = newEntry.Cron
				oldEntry.Message = newEntry.Message
				if newEntry.NextTs == nil {
					nextTs2, err := gronx.NextTick(*newEntry.Cron, true)
					if err != nil {
						return err
					}
					oldEntry.NextTs = &nextTs2
				} else {
					oldEntry.NextTs = newEntry.NextTs
				}
				scheduler.update(oldEntry)
			}
			oldEntry.Name = newEntry.Name
		} else {
			// added
			slog.Info("new entry", "entry", newEntry.Id)
			scheduler.entries[id] = newEntry
		}
	}

	for id, oldEntry := range scheduler.entries {
		if _, ok := entries[id]; !ok {
			slog.Info("deleted entry", "entry", oldEntry.Id)
			delete(scheduler.entries, oldEntry.Id)
		}
	}

	// scheduler.entries = entries

	if scheduler.entry != nil {
		entry2 := scheduler.entries[scheduler.entry.Id]
		if entry2.NextTs.Equal(*scheduler.entry.NextTs) {
			scheduler.entry = entry2
		} else {
			slog.Info("RESCHEDULE", "id", scheduler.entry.Id, "t1", entry2.NextTs, "t2", scheduler.entry.NextTs)
			scheduler.scheduleNext()
		}
	}
	return nil
}

func (scheduler *Scheduler) scheduleNext() {
	var next *Entry = scheduler.getNext()
	// if next != nil && scheduler.entry != nil && next.Id == scheduler.entry.Id {
	// 	return
	// }
	scheduler.entry = next

	var d time.Duration
	if next != nil {
		d = time.Until(*next.NextTs)
		slog.Info("next", "entry", next.Id, "nextTs", next.NextTs)
	} else {
		d = 1 * time.Second
		slog.Info("next", "entry", nil)
	}

	// scheduler.timer.Reset(time.Since(*next.NextTs))
	scheduler.timer.Stop()
	select {
	case <-scheduler.timer.C:
	default:
	}
	scheduler.timer.Reset(d)
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

func (scheduler *Scheduler) getEntries() (EntryMap, error) {
	db := scheduler.db
	stmt, err := db.Prepare(
		`select id, name, is_active, cron, next_ts, last_ts, message 
		from barn_entry`,
	)
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
		if e.IsActive {
			if e.Message == nil {
				// we don't know what to do...
				slog.Warn("invalid entry", "entry", e)
				scheduler.deactivate(&e)
			} else if e.Cron == nil && e.NextTs == nil {
				// we don't know when to do...
				slog.Warn("invalid entry", "entry", e)
				scheduler.deactivate(&e)
			} else {
				if e.NextTs == nil {
					nextTs2, err := gronx.NextTick(*e.Cron, true)
					if err != nil {
						slog.Info("invalid cron string", "entry", e)
						continue
					}
					e.NextTs = &nextTs2
					scheduler.update(&e)
				}
				slog.Info("the entry is active", "entry", e)
				entries[e.Id] = &e
			}
		} else {
			slog.Info("the entry is inactive", "entry", e)
		}
	}
	return entries, nil
}

func (scheduler *Scheduler) Add(name string, cron *string, nextTs *time.Time, message string) error {
	// fake 1
	// cron := "*/5 * * * * *"

	if cron == nil && nextTs == nil {
		return fmt.Errorf("invalid args")
	}
	// if cron != nil && nextTs == nil {
	// 	nextTs2, err := gronx.NextTick(*cron, true)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	nextTs = &nextTs2
	// }

	slog.Info("create the entry", "name", name, "cron", cron, "message", message)

	db := scheduler.db
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
		slog.Info("the entry was already deleted", "entry", id)
	}
	return nil
}

func (scheduler *Scheduler) update(entry *Entry) error {
	db := scheduler.db

	slog.Info("update the entry", "entry", entry)
	res, err := db.Exec(
		`update barn_entry 
		set is_active=?, next_ts=?, last_ts=? 
		where id=?`,
		entry.IsActive, entry.NextTs, entry.LastTs, entry.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		// not an erros, need to reload entries...
		return fmt.Errorf("an object deleted")
	}
	return nil
}

func (scheduler *Scheduler) deactivate(entry *Entry) error {
	db := scheduler.db

	entry.IsActive = false
	slog.Info("deactivate the entry", "entry", entry)
	res, err := db.Exec(
		`update barn_entry 
		set is_active=? 
		where id=?`,
		false, entry.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		slog.Info("the entry was deleted somewhen", "entry", entry)
	}
	return nil
}

func (scheduler *Scheduler) DeleteByName(name string) error {
	return nil
}
