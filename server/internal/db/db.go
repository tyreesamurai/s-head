package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	ID         string  `gorm:"primaryKey"`
	CreatorID  string  `gorm:"not null"`
	Status     string  `gorm:"not null"`
	MaxPlayers int     `gorm:"not null"`
	Private    bool    `gorm:"not null"`
	Password   *string `gorm:"default:null"`
	WinnderID  *string `gorm:"default:null"`
}

type Store struct {
	db *gorm.DB
}

func NewStore() (*Store, error) {
	db, err := gorm.Open(sqlite.Open("games.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Game{}); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}
