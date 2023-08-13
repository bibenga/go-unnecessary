package main

//go:generate jet -source=sqlite -dsn=../palabras.db -schema=palabras -path=./

import (
	"flag"
	"log"
	"os"
	"time"
	"unnecessary/palabras/models"

	"github.com/glebarez/sqlite"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initGorm(filename string) *gorm.DB {
	gormLogger := logger.New(
		log.New(os.Stdout, "[sql] -", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(
		sqlite.Open(filename),
		&gorm.Config{
			Logger: gormLogger,
		},
	)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("")

	isMigrationDisabled := flag.Bool("no-migrate", false, "disable db migration with gorm")
	flag.Parse()
	log.Printf("should we migrate a db: %v", !*isMigrationDisabled)

	db := initGorm("palabras.db")
	log.Printf("db: %v", db)

	if !*isMigrationDisabled {
		db.AutoMigrate(&models.User{}, &models.TextPair{}, &models.StudyState{})
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		result := tx.
			Where("lower(email) = lower(?)", "user1").
			Attrs(models.User{Email: "user1"}).
			FirstOrCreate(&user)
		if result.Error != nil {
			return result.Error
		}
		created := result.RowsAffected == 1
		log.Printf("user was created or fetched: created=%+v, user=%+v", created, user)
		return nil
	})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
