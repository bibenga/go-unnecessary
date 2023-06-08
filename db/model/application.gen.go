// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameApplication = "application"

// Application mapped from table <application>
type Application struct {
	ID             int64         `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:Application ID" json:"id"`
	DictPlatformID int16         `gorm:"column:dict_platform_id;type:smallint;not null;uniqueIndex:application_dict_platform_id_name_key,priority:1" json:"dict_platform_id"`
	Name           string        `gorm:"column:name;type:character varying;not null" json:"name"`
	DisplayName    *string       `gorm:"column:display_name;type:character varying" json:"display_name"`
	CreatedTs      *time.Time    `gorm:"column:created_ts;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_ts"`
	ModifiedTs     *time.Time    `gorm:"column:modified_ts;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"modified_ts"`
	DeletedTs      *time.Time    `gorm:"column:deleted_ts;type:timestamp with time zone" json:"deleted_ts"`
	SomeFlg        bool          `gorm:"column:some_flg;type:boolean;not null" json:"some_flg"`
	SomeDouble1    float64       `gorm:"column:some_double1;type:double precision;not null" json:"some_double1"`
	SomeDouble2    *float64      `gorm:"column:some_double2;type:double precision" json:"some_double2"`
	SomeJSON       *string       `gorm:"column:some_json;type:jsonb;not null;default:'{}'::jsonb" json:"some_json"`
	DictPlatform   *DictPlatform `gorm:"foreignKey:DictPlatformID" json:"dict_platform"`
	Tags           []*Tag        `gorm:"many2many:application_tag" json:"tags"`
}

// TableName Application's table name
func (*Application) TableName() string {
	return TableNameApplication
}