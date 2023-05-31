package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"unnecessary/db/dao"
	"unnecessary/db/model"

	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// go get -u gorm.io/gen
// go install gorm.io/gen/tools/gentool@latest

func initGorm() *gorm.DB {
	log.Printf("--------------------------------------------")

	gormLogger := logger.New(
		log.Default(),
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
	log.Printf("connected")
	return db
}

func playWithDao(db *gorm.DB) {
	log.Printf("--------------------------------------------")
	log.Printf("playWithDbConGormAndGen")
	dao.SetDefault(db)

	log.Printf("run transaction")
	err := dao.Q.Transaction(func(t *dao.Query) error {
		log.Printf("inside transaction")

		ctx := context.Background()
		q := t.WithContext(ctx)

		android, err := q.DictPlatform.Where(
			t.DictPlatform.Name.Eq("ANDROID"),
		).Take()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				displayName := "Android"
				android = &model.DictPlatform{
					ID:          1,
					Name:        "ANDROID",
					DisplayName: &displayName,
				}
				err = q.DictPlatform.Create(android)
				log.Printf("create DictPlatform")
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		log.Printf("Update DictPlatform")
		dn := fmt.Sprint(rand.Int63())
		ts := time.Now()
		android.DisplayName = &dn
		android.ModifiedTs = &ts
		err = q.DictPlatform.Save(android)
		if err != nil {
			panic(err)
		}
		log.Printf("Update DictPlatform 2")
		updated, err := q.DictPlatform.Where(
			t.DictPlatform.ID.Eq(android.ID),
		).Updates(model.DictPlatform{
			DisplayName: android.DisplayName,
			ModifiedTs:  android.ModifiedTs,
		})
		if err != nil {
			panic(err)
		}
		log.Printf("RowsAffected: %d", updated.RowsAffected)

		log.Printf("Create tag")
		tag1 := &model.Tag{
			Name:       "TAG1",
			ModifiedTs: &ts,
		}
		tag1, err = q.Tag.Where(
			t.Tag.Name.Eq(tag1.Name),
		).Attrs(
			t.Tag.Name,
		).FirstOrCreate()
		if err != nil {
			return err
		}

		log.Printf("Create tag 2")
		err = q.Tag.Clauses(
			clause.OnConflict{DoNothing: true},
		).Create(tag1)
		if err != nil {
			return err
		}

		// log.Printf("Create tag 3")
		// err = q.Tag.Clauses(clause.OnConflict{
		// 	Columns: []clause.Column{
		// 		{Name: tx.Tag.Name.ColumnName().String()},
		// 	},
		// 	DoUpdates: clause.Assignments(map[string]interface{}{
		// 		// tx.Tag.Name.ColumnName().String():       tag1.Name,
		// 		tx.Tag.ModifiedTs.ColumnName().String(): tag1.ModifiedTs,
		// 	}),
		// }).Create(tag1)
		// if err != nil {
		// 	return err
		// }

		log.Printf("Load tag")
		tag1, err = q.Tag.Where(
			t.Tag.Name.Like(tag1.Name),
		).First()
		if err != nil {
			return err
		}

		log.Printf("try to find Application \"CHROME\"")
		chrome, err := q.Application.Where(
			t.Application.DictPlatformID.Eq(android.ID),
			t.Application.Name.Eq("CHROME"),
		).Joins(
			dao.Application.DictPlatform,
		).Preload(
			dao.Application.Tags,
		).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if chrome != nil {
			log.Printf("delete Application")
			_, err = q.Application.Where(t.Application.ID.Eq(chrome.ID)).Delete()
			if err != nil {
				return err
			}
		}
		displayName := "Chrome"
		// cascade insert with update DictPlatform
		// chrome = &model.Application{DictPlatform: android, Name: "CHROME", DisplayName: &displayName}
		chrome = &model.Application{
			DictPlatformID: android.ID,
			DictPlatform:   android,
			Name:           "CHROME",
			DisplayName:    &displayName,
			Tags:           []*model.Tag{tag1},
		}
		log.Printf("before Create Application")
		err = q.Application.Omit(
			t.Application.DictPlatform.Field(),
			t.Application.Tags.Field("*"),
		).Create(chrome)
		if err != nil {
			return err
		}
		// err = tx.Application.Tags.Model(chrome).Append(tag1)
		// if err != nil {
		// return err
		// }

		return nil
	})
	if err != nil {
		panic(err)
	}
	log.Printf("done")
}

func playWithModelOnly(db *gorm.DB) {
	log.Printf("--------------------------------------------")
	log.Printf("playWithDbConGormAndGenModel")

	err := db.Transaction(func(tx *gorm.DB) error {
		log.Printf("inside transaction")

		var android model.DictPlatform
		result := tx.Model(&model.DictPlatform{}).Where(
			"Name ILIKE ?", "ANDROID",
		).Take(&android)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		log.Printf("loaded DictPlatform: %+v", android)
		if result.RowsAffected == 0 {
			displayName := "Android"
			android = model.DictPlatform{
				ID:          1,
				Name:        "ANDROID",
				DisplayName: &displayName,
			}
			result := tx.Create(&android)
			if result.Error != nil {
				return result.Error
			}
			log.Printf("created DictPlatform: %+v", android)
		}

		var tag1 model.Tag
		result = tx.Model(&model.Tag{}).Where(
			"LOWER(Name) = LOWER(?)", "TAG1",
		).Take(&tag1)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		log.Printf("loaded Tag: %+v", tag1)
		if result.RowsAffected == 0 {
			tag1 = model.Tag{
				Name: "TAG1",
			}
			result := tx.Create(&tag1)
			if result.Error != nil {
				return result.Error
			}
			log.Printf("created Tag: %+v", tag1)
		}

		var safari model.Application
		result = tx.Model(&model.Application{}).Where(
			// "DictPlatformID = ? AND Name ILIKE ?",
			"dict_platform_id = ? AND name ILIKE ?",
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
			log.Printf("loaded Application deleted")
		}
		safari = model.Application{
			DictPlatform:   &android,
			DictPlatformID: android.ID,
			Name:           "CHROME",
			Tags:           []*model.Tag{&tag1},
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
		err := tx.Model(&safari).Omit("Tags.*").Association("Tags").Replace([]*model.Tag{&tag1})
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
	playWithDao(db)
	// playWithModelOnly(db)
}
