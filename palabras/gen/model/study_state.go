//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type StudyState struct {
	ID         *int32 `sql:"primary_key"`
	Question   string
	Answer     string
	TextPairID int32
	IsDone     float64
	IsSkiped   float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
