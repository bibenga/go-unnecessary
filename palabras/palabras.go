package main

//go:generate jet -source=sqlite -dsn=../palabras.db -schema=palabras -path=./

import (
	"errors"
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

	db := initGorm("palabras.db")
	log.Printf("db: %v", db)

	db.AutoMigrate(&models.User{})

	err := db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		result := tx.Model(&models.User{}).Where(
			// "DictPlatformID = ? AND Name ILIKE ?",
			"lower(email) = lower(?)",
			"user1",
		).Take(&user)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		if result.RowsAffected == 1 {
			log.Printf("user was founded: %v", user)
		} else {
			user = models.User{
				Email: "user1",
			}
			result = tx.Create(&user)
			if result.Error != nil {
				return result.Error
			}
			log.Printf("user was created: %v", user)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}
