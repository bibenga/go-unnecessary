package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primarykey;autoIncrement"`
	Email     string         `gorm:"not null;uniqueIndex:u_user_email,expression:lower(email)"`
	IsActive  bool           `gorm:"not null;default:True"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime;<-:create"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index:u_user_deleted_at"`
}

func (User) TableName() string {
	return "user"
}

type TextPair struct {
	ID        uint64         `gorm:"primarykey;autoIncrement"`
	UserID    uint64         `gorm:"not null"`
	User      *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Text1     string         `gorm:"not null"`
	Text2     string         `gorm:"not null"`
	IsLearned bool           `gorm:"not null;default:False"`
	CreatedAt time.Time      `gorm:"not null;autoCreateTime;<-:create"`
	UpdatedAt time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index:u_text_pair_deleted_at"`
}

func (TextPair) TableName() string {
	return "text_pair"
}

type StudyState struct {
	ID         uint64         `gorm:"primarykey;autoIncrement"`
	Question   string         `gorm:"not null"`
	Answer     string         `gorm:"not null"`
	TextPairID uint64         `gorm:"not null"`
	TextPair   *TextPair      `gorm:"foreignKey:TextPairID;constraint:OnDelete:CASCADE"`
	IsDone     bool           `gorm:"not null;default:False"`
	IsSkiped   bool           `gorm:"not null;default:False"`
	CreatedAt  time.Time      `gorm:"not null;autoCreateTime;<-:create"`
	UpdatedAt  time.Time      `gorm:"not null;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index:u_study_state_deleted_at"`
}

func (StudyState) TableName() string {
	return "study_state"
}
