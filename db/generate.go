package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	log.SetPrefix("gorm - ")

	dsn := "host=host.docker.internal port=5432 user=rds password=sqlsql dbname=go TimeZone=UTC sslmode=disable"
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: dsn,
		}),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		// db:                db,
		OutPath:           "db/dao",
		ModelPkgPath:      "db/model",
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		// WithUnitTest:      true,
	})
	g.UseDB(db)

	tag := g.GenerateModel("tag")

	dict_platform := g.GenerateModel(
		"dict_platform",
		gen.FieldGORMTag(
			"create_ts",
			// "column:create_ts;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP;autoCreateTime",
			// "column:create_ts;type:timestamp with time zone;not null;autoCreateTime",
			func(tag field.GormTag) field.GormTag {
				tag.Remove(field.TagKeyGormDefault)
				tag.Set(field.TagKeyGormColumn, "create_ts")
				tag.Set(field.TagKeyGormType, "timestamp with time zone")
				tag.Set(field.TagKeyGormNotNull, "")
				tag.Set("autoCreateTime", "")
				return tag
			},
		),
		gen.FieldGORMTag(
			"modified_ts",
			// "column:modified_ts;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP;autoUpdateTime",
			// "column:modified_ts;type:timestamp with time zone;not null;autoUpdateTime",
			func(tag field.GormTag) field.GormTag {
				tag.Remove(field.TagKeyGormDefault)
				tag.Set(field.TagKeyGormColumn, "modified_ts")
				tag.Set(field.TagKeyGormType, "timestamp with time zone")
				tag.Set(field.TagKeyGormNotNull, "")
				tag.Set("autoUpdateTime", "")
				return tag
			},
		),
	)

	application := g.GenerateModel(
		"application",
		// gen.FieldGenType("id", "uuid.NullUUID"),
		// gen.FieldGenTypeReg("deleted_ts", "DeletedAt"),
		// gen.FieldGenType("deleted_ts", "gorm.DeletedAt"),
		gen.FieldRelate(field.BelongsTo, "DictPlatform", dict_platform,
			&field.RelateConfig{
				GORMTag:       field.GormTag{"foreignKey": "DictPlatformID"},
				RelatePointer: true,
			}),
		gen.FieldRelate(field.Many2Many, "Tags", tag,
			&field.RelateConfig{
				GORMTag:            field.GormTag{"many2many": "application_tag"},
				RelateSlicePointer: true,
			}),
	)

	application_tag := g.GenerateModel(
		"application_tag",
		gen.FieldRelate(field.BelongsTo, "Application", application,
			&field.RelateConfig{
				// GORMTag:       "foreignKey:ApplicationID",
				GORMTag:       field.GormTag{"foreignKey": "ApplicationID"},
				RelatePointer: true,
			}),
		gen.FieldRelate(field.BelongsTo, "Tag", tag,
			&field.RelateConfig{
				// GORMTag:       "foreignKey:TagID",
				GORMTag:       field.GormTag{"foreignKey": "TagID"},
				RelatePointer: true,
			}),
	)

	audit := g.GenerateModel("audit")

	g.ApplyBasic(tag, dict_platform, application, application_tag, audit)

	g.Execute()
}
