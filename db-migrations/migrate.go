package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	// _ "github.com/golang-migrate/migrate/v4/database/postgres"
	// _ "github.com/golang-migrate/migrate/v4/database/sqlite"
)

//go:embed sql
var embededDbMigrationsFiles embed.FS

type MigratorLoggerImpl struct {
}

func (l *MigratorLoggerImpl) Verbose() bool {
	return true
}

func (l *MigratorLoggerImpl) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func newMigrator(sqlDb *sql.DB) *migrate.Migrate {
	migratorDriver, err := pgx.WithInstance(sqlDb, &pgx.Config{
		// SchemaName: "go",
	})
	if err != nil {
		panic(err)
	}
	// migrator, err := migrate.NewWithDatabaseInstance(
	// 	"file://migrations/sql",
	// 	"postgres",
	// 	migratorDriver,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	source, err := iofs.New(embededDbMigrationsFiles, "sql")
	if err != nil {
		panic(err)
	}
	migrator, err := migrate.NewWithInstance(
		"iofs", source,
		"postgres", migratorDriver,
	)
	if err != nil {
		panic(err)
	}
	migrator.Log = &MigratorLoggerImpl{}
	return migrator
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[migrate] ")

	log.Print("connect to DB")
	dsn := fmt.Sprintf(
		"host=%s user=rds password=sqlsql dbname=go port=5432 TimeZone=UTC sslmode=disable",
		"host.docker.internal",
	)

	// sqlDb, err := sql.Open("postgres", dsn)
	sqlDb, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	// Drop delete all tables, but MigrationTable created inside migratorPostgres.WithInstance
	log.Println("drop all")
	migrator := newMigrator(sqlDb)
	err = migrator.Drop()
	if err != nil {
		panic(err)
	}

	log.Println("migrate")
	migrator = newMigrator(sqlDb)
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	log.Println("done")
}
