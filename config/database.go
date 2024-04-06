package config

import (
	"log"
	"os"

	"github.com/Hamedblue1381/restaurant-reserve/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDBConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	db.AutoMigrate(&models.User{}, &models.Food{}, &models.Sides{}, &models.Category{}, &models.MealType{})

	return db
}
