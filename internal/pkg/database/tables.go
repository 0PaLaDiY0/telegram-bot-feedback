package database

import (
	"gorm.io/gorm"
)

// User table
type User struct {
	gorm.Model
	ChatID     int
	State      int
	Nickname   string
	IsEmployee bool       `gorm:"default:false"`
	IsReceiver bool       `gorm:"default:false"`
	Review     []Review   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Question   []Question `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Review table
type Review struct {
	gorm.Model
	Rating int
	Text   string
	UserID int
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Question table
type Question struct {
	gorm.Model
	Header                 string
	UserID                 int
	User                   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AnswererID             int
	Answerer               User                     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	QuestionCorrespondence []QuestionCorrespondence `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HaveAnswer             bool                     `gorm:"default:false"`
	IsClosed               bool                     `gorm:"default:false"`
}

// QuestionCorrespondence table
type QuestionCorrespondence struct {
	gorm.Model
	QuestionID int
	MessageID  int
	UserID     int
	User       User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsEmployee bool
}
