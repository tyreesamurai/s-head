package db

import (
	"shithead/internal/models"

	"gorm.io/gorm"
)

func CreateGame(db *gorm.DB, game models.Game) error {
	return db.Create(&game).Error
}

func GetGameByID(db *gorm.DB, id string) (*models.Game, error) {
	return &models.Game{}, db.First(&models.Game{}, "id = ?", id).Error
}
