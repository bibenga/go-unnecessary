package main

// https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"unnecessary/db-gorm/models"

	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initGorm() *gorm.DB {
	gormLogger := logger.New(
		log.New(os.Stdout, "[sql] -", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// dsn := "host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "host=db user=postgres password=postgres dbname=postgres port=5432 TimeZone=UTC"
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: dsn,
		}),
		&gorm.Config{
			Logger: gormLogger,
			// Logger: log,
		},
	)
	if err != nil {
		panic(err)
	}
	return db
}

func playWithSomeModels(db *gorm.DB) {
	log.Printf("------------------------------")
	log.Printf("playWithSomeModels")

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Product{})

	db.Transaction(func(tx *gorm.DB) error {
		return nil
	}, &sql.TxOptions{ReadOnly: true})

	u1 := models.User{Email: fmt.Sprintf("email-%d", rand.Int31()), Vesion: 1}
	log.Printf("Before user created -> %+v\n", &u1)
	db.Create(&u1)
	log.Printf("User created -> %+v\n", &u1)

	var count int64
	db.Model(&models.Product{}).Count(&count)
	log.Printf("Product count -> %+v\n", count)

	// Create
	// pMeta := datatypes.JSON([]byte(`{"a":1}`))
	pMetaMode := "Unnecessary"
	// pMeta := datatypes.JSONType[models.ProductMeta]{
	// 	// Data: models.ProductMeta{Mode: &pMetaMode},
	// 	Mode: &pMetaMode,
	// }
	pMeta := datatypes.NewJSONType(models.ProductMeta{
		Mode: &pMetaMode,
	})
	p := models.Product{
		Code:  "D42",
		Price: 100,
		Meta:  &pMeta,
	}
	db.Create(&p)
	log.Printf("Product created -> %+v\n", &p)

	// Read
	var product models.Product
	db.First(&product, p.ID)              // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(models.Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, p.ID)
}

func playWithSqlModels(db *gorm.DB) {
	log.Printf("------------------------------")
	log.Printf("playWithSqlModels")

	err := db.Transaction(func(tx *gorm.DB) error {
		log.Printf("inside transaction")

		// load DictPlatform
		var android models.DictPlatform
		result := tx.Model(&models.DictPlatform{}).Where(
			"LOWER(name) = LOWER(?)", "ANDROID",
		).Take(&android)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		log.Printf("loaded DictPlatform: %+v", android)
		if result.RowsAffected == 1 {
			log.Printf("delete Applications")
			result := tx.Where(
				"dict_platform_id = ?", android.ID,
			).Delete(&models.Application{})
			if result.Error != nil {
				return result.Error
			}
			log.Printf("deleted %v Applications", result.RowsAffected)

			log.Printf("delete loaded DictPlatform")
			result = tx.Delete(&android)
			if result.Error != nil {
				return result.Error
			}
		} else {
			log.Printf("DictPlatform not found (RowsAffected=%v)", result.RowsAffected)
		}

		// create DictPlatform
		displayName := "Android"
		android = models.DictPlatform{
			ID:          1,
			Name:        "ANDROID",
			DisplayName: &displayName,
		}
		result = tx.Create(&android)
		if result.Error != nil {
			return result.Error
		}
		log.Printf("created DictPlatform: %+v", android)

		// load Tag
		var tag1 models.Tag
		result = tx.Model(&models.Tag{}).Where(
			"LOWER(name) = LOWER(?)", "TAG1",
		).Take(&tag1)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		log.Printf("loaded Tag: %+v", tag1)
		if result.RowsAffected == 1 {
			log.Printf("delete loaded Tag")
			result := tx.Delete(&tag1)
			if result.Error != nil {
				return result.Error
			}
		} else {
			log.Printf("Tag not found (RowsAffected=%v)", result.RowsAffected)
		}
		// create Tag
		tag1 = models.Tag{
			ID:   1,
			Name: "TAG1",
		}
		result = tx.Create(&tag1)
		if result.Error != nil {
			return result.Error
		}
		log.Printf("created Tag: %+v", tag1)

		// load Application
		var safari models.Application
		result = tx.Model(&models.Application{}).Where(
			// "DictPlatformID = ? AND Name ILIKE ?",
			"dict_platform_id = ? AND lower(name) = lower(?)",
			android.ID, "CHROME",
		).Take(&safari)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		log.Printf("loaded Application: %+v", safari)
		if result.RowsAffected == 1 {
			result = tx.Delete(&safari)
			if result.Error != nil {
				return result.Error
			}
		}
		safari = models.Application{
			DictPlatform:   &android,
			DictPlatformID: android.ID,
			Name:           "CHROME",
			Tags:           []*models.Tag{&tag1},
		}
		log.Printf("create Application: %+v", safari)
		result = tx.Omit(
			"DictPlatform",
			"Tags.*",
		).Create(&safari)
		if result.Error != nil {
			return result.Error
		}
		log.Printf("created Application: %+v", tag1)

		log.Printf("replace Tags")
		err := tx.Model(&safari).Omit("Tags.*").Association("Tags").Replace([]*models.Tag{&tag1})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("jet - ")

	db := initGorm()
	// playWithSomeModels(db)
	playWithSqlModels(db)
}
