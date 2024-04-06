package models

import (
	"gorm.io/gorm"
)

type Food struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `json:"name"`
	Quanity      string        `json:"quanity"`
	CategoryID   uint          // Foreign key for Category
	Category     Category      `json:"category"` // Category relationship
	MealTypeID   uint          // Foreign key for MealType
	MealType     MealType      `json:"meal_type"` // MealType relationship
	gorm.Model   `json:"-" swaggerignore:"true"`
}

type FoodHandler struct {
	db *gorm.DB
}

func NewFoodHandler(db *gorm.DB) *FoodHandler {
	return &FoodHandler{db}
}

func (f *FoodHandler) CreateFood(food *Food) error {
	return f.db.Create(food).Error
}

func (f *FoodHandler) GetFood(id uint) (*Food, error) {
	var food Food
	result := f.db.First(&food, id)
	return &food, result.Error
}

func (f *FoodHandler) GetFoods() ([]Food, error) {
	var foods []Food
	result := f.db.Find(&foods)
	return foods, result.Error
}

func (f *FoodHandler) UpdateFood(id uint, food *Food) error {
	result := f.db.Model(&Food{}).Where("id = ?", id).Updates(food)
	return result.Error
}

func (f *FoodHandler) DeleteFood(id uint) error {
	result := f.db.Delete(&Food{}, id)
	return result.Error
}
