package database

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes the SQLite database
func Init(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(User{}, Review{}, Question{}, QuestionCorrespondence{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
