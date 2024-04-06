package models

import "gorm.io/gorm"

type MealType struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `json:"name"`
	Foods      []Food `gorm:"foreignKey:MealTypeID"`
	gorm.Model `json:"-" swaggerignore:"true"`
}

type MealTypeHandler struct {
	db *gorm.DB
}

func NewMealTypeHandler(db *gorm.DB) *MealTypeHandler {
	return &MealTypeHandler{db}
}

func (h *MealTypeHandler) CreateMealType(mealType *MealType) error {
	return h.db.Create(mealType).Error
}

func (h *MealTypeHandler) GetMealType(id uint) (*MealType, error) {
	var mealType MealType
	result := h.db.First(&mealType, id)
	return &mealType, result.Error
}

func (h *MealTypeHandler) GetMealTypes() ([]MealType, error) {
	var mealTypes []MealType
	result := h.db.Find(&mealTypes)
	return mealTypes, result.Error
}

func (h *MealTypeHandler) UpdateMealType(id uint, mealType *MealType) error {
	result := h.db.Model(&MealType{}).Where("id = ?", id).Updates(mealType)
	return result.Error
}

func (h *MealTypeHandler) DeleteMealType(id uint) error {
	result := h.db.Delete(&MealType{}, id)
	return result.Error
}
