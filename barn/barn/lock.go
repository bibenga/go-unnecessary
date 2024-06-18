package barn

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Lock struct {
	Name     string
	LockedAt time.Time
	LockedBy string
}

type LockManager struct {
	hostname   string
	lockName   string
	interval   time.Duration
	expiration time.Duration
	isLocked   bool
	db         *sql.DB
	tableName  string
	stop       chan struct{}
	stopped    chan struct{}
}

func NewLockManager(db *sql.DB) *LockManager {
	uuid, _ := uuid.NewRandom()
	manager := LockManager{
		hostname:   uuid.String(),
		lockName:   "barn",
		interval:   1 * time.Second,
		expiration: 10 * time.Second,
		isLocked:   false,
		db:         db,
		tableName:  "barn_lock",
		stop:       make(chan struct{}),
		stopped:    make(chan struct{}),
	}
	return &manager
}

func (manager *LockManager) InitializeDB() error {
	db := manager.db
	slog.Info("create table", "table", manager.tableName)
	sqlStmt := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS "%s"  (
	    name VARCHAR NOT NULL,
        locked_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP) NOT NULL,
        locked_by VARCHAR NOT NULL DEFAULT '',
        PRIMARY KEY (name)
	)`, manager.tableName)
	_, err := db.Exec(sqlStmt)
	return err
}

func (manager *LockManager) Stop() {
	slog.Info("stopping", "lock", manager.lockName)
	manager.stop <- struct{}{}
	<-manager.stopped
	close(manager.stop)
	close(manager.stopped)
	slog.Info("stopped", "lock", manager.lockName)
}

func (manager *LockManager) Run() {
	if manager.isExist() {
		slog.Info("the lock is exists", "lock", manager.lockName)
	} else {
		manager.create()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	slog.Info("started")
	for {
		select {
		case <-manager.stop:
			slog.Info("terminate")
			manager.stopped <- struct{}{}
			return
		case <-ticker.C:
			manager.run()
		}
	}
}

func (manager *LockManager) run() {
	dbLock := manager.getDbLock()
	if manager.isLocked {
		if dbLock.LockedBy == manager.hostname {
			if manager.tryUpdate(dbLock) {
				slog.Info("the lock is still captured", "lock", manager.lockName)
			} else {
				slog.Warn("the lock was captured unexpectedly by someone", "lock", manager.lockName)
				manager.isLocked = false
				manager.onReleased()
			}
		} else {
			slog.Warn("the lock was captured by someone", "lock", manager.lockName)
			manager.isLocked = false
			manager.onReleased()
		}
	} else if time.Since(dbLock.LockedAt) > manager.expiration {
		slog.Info("the lock is rotten", "lock", manager.lockName)
		if manager.tryUpdate(dbLock) {
			manager.isLocked = true
			manager.onCaptured()
		}
	}
}

func (manager *LockManager) isExist() bool {
	db := manager.db
	stmt, err := db.Prepare(
		`select 1 
		from barn_lock 
		where name = $1 
		limit 1`)
	if err != nil {
		slog.Error("cannot prepare query", "error", err)
		panic(err)
	}
	defer stmt.Close()
	var count int
	row := stmt.QueryRow(manager.lockName)
	switch err := row.Scan(&count); err {
	case sql.ErrNoRows:
		slog.Info("the lock is not exist", "lock", manager.lockName)
		return false
	case nil:
		slog.Info("the lock is exist", "lock", manager.lockName)
		return true
	default:
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}
}

func (manager *LockManager) create() {
	db := manager.db
	slog.Info("create the lock", "name", manager.lockName)

	tx, err := db.Begin()
	if err != nil {
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}

	res, err := tx.Exec(
		`insert into barn_lock(name, locked_at, locked_by) 
		values ($1, $2, $3) 
		on conflict (name) do nothing`,
		manager.lockName, time.Now().Add((-300*24)*time.Hour), "",
	)
	if err != nil {
		slog.Error("cannot create lock", "name", manager.lockName, "error", err)
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}
	if rowsAffected == 1 {
		slog.Info("the lock is created", "name", manager.lockName)
	} else {
		slog.Info("the lock was created by someone", "name", manager.lockName)
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}
}

func (manager *LockManager) getDbLock() *Lock {
	db := manager.db
	stmt, err := db.Prepare(
		`select locked_at, locked_by 
		from barn_lock 
		where name = $1`,
	)
	if err != nil {
		slog.Error("cannot prepare query", "error", err)
		panic(err)
	}
	defer stmt.Close()
	var dbLock Lock = Lock{Name: manager.lockName}
	row := stmt.QueryRow(manager.lockName)
	switch err := row.Scan(&dbLock.LockedAt, &dbLock.LockedBy); err {
	case sql.ErrNoRows:
		slog.Info("the lock is not found", "lock", manager.lockName)
		return nil
	case nil:
		slog.Info("the lock is found", "lock", manager.lockName, "state", dbLock)
		return &dbLock
	default:
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}
}

func (manager *LockManager) tryUpdate(dbLock *Lock) bool {
	db := manager.db
	res, err := db.Exec(
		`update barn_lock 
		set locked_at=$1, locked_by=$2 
		where name=$3 and locked_by=$4 and locked_at=$5`,
		time.Now(), manager.hostname,
		manager.lockName, dbLock.LockedBy, dbLock.LockedAt,
	)
	if err != nil {
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("db error", "lock", manager.lockName, "error", err)
		panic(err)
	}
	return rowsAffected == 1
}

func (manager *LockManager) onError(err error) {
	slog.Error("unexpected error", "lock", manager.lockName, "error", err)
	panic(err)
}

func (manager *LockManager) onCaptured() {
	slog.Info("lock is captured", "lock", manager.lockName)
}

func (manager *LockManager) onReleased() {
	slog.Warn("lock is released", "lock", manager.lockName)
}
