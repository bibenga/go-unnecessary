package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.NullUUID  `gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Email     string         `gorm:"not null;uniqueIndex;comment:user email"`
	Vesion    uint           `gorm:"default:1"`
}

func (User) TableName() string {
	return "users"
}

type ProductMeta struct {
	Mode *string
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
	Meta  *datatypes.JSONType[ProductMeta]
}

// DictPlatform
type DictPlatform struct {
	ID          int16 `gorm:"primaryKey"`
	Name        string
	DisplayName *string
	CreatedTs   *time.Time `gorm:"autoCreateTime;type:timestamp with time zone"`
	ModifiedTs  *time.Time `gorm:"autoUpdateTime;type:timestamp with time zone"`
	DeletedTs   *time.Time
}

func (*DictPlatform) TableName() string {
	return "dict_platform"
}

// Tag
type Tag struct {
	ID         int64 `gorm:"primaryKey;autoIncrement:true"`
	Name       string
	ModifiedTs *time.Time `gorm:"autoUpdateTime;type:timestamp with time zone"`
}

func (*Tag) TableName() string {
	return "tag"
}

// Application
type Application struct {
	ID             int64 `gorm:"primaryKey;autoIncrement:true"`
	DictPlatformID int16
	Name           string
	DisplayName    *string
	CreatedTs      *time.Time `gorm:"autoCreateTime;type:timestamp with time zone"`
	ModifiedTs     *time.Time `gorm:"autoUpdateTime;type:timestamp with time zone"`
	DeletedTs      *time.Time `gorm:"type:timestamp with time zone"`
	SomeFlg        bool
	SomeDouble1    float64
	SomeDouble2    *float64
	DictPlatform   *DictPlatform `gorm:"foreignKey:DictPlatformID"`
	Tags           []*Tag        `gorm:"many2many:application_tag"`
}

func (*Application) TableName() string {
	return "application"
}

// ApplicationTag
type ApplicationTag struct {
	ID            int64 `gorm:"primaryKey;autoIncrement:true"`
	ApplicationID int64
	TagID         int64
	ModifiedTs    *time.Time   `gorm:"autoUpdateTime;type:timestamp with time zone"`
	Application   *Application `gorm:"foreignKey:ApplicationID"`
	Tag           *Tag         `gorm:"foreignKey:TagID"`
}

func (*ApplicationTag) TableName() string {
	return "application_tag"
}
