package main

// https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/viper"
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

func GetDsn() string {
	// return "host=host.docker.internal port=5432 user=rds password=sqlsql dbname=go TimeZone=UTC sslmode=disable"
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s TimeZone=%s sslmode=%s",
		viper.GetString("DB_POSTGRES_HOST"),
		viper.GetString("DB_POSTGRES_PORT"),
		viper.GetString("DB_POSTGRES_USERNAME"),
		viper.GetString("DB_POSTGRES_PASSWORD"),
		viper.GetString("DB_POSTGRES_DB"),
		viper.GetString("DB_POSTGRES_TIMEZONE"),
		viper.GetString("DB_POSTGRES_SSLMODE"),
	)
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("[migrate] ")

	log.Print("Load config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// log.Printf(" - %+v", viper.Get("DB_POSTGRES_HOST"))

	log.Print("connect to DB")
	dsn := GetDsn()
	log.Printf("[unsecure] dsn is '%s'", dsn)

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
