package models

import "gorm.io/gorm"

type Sides struct {
	ID           uint          `gorm:"primaryKey"`
	Name         string        `json:"name"`
	Quantity     string        `json:"quantity"`
	gorm.Model   `json:"-" swaggerignore:"true"`
}

type SidesHandler struct {
	db *gorm.DB
}

func NewSidesHandler(db *gorm.DB) *SidesHandler {
	return &SidesHandler{db}
}

func (h *SidesHandler) CreateSides(sides *Sides) error {
	return h.db.Create(sides).Error
}

func (h *SidesHandler) GetSide(id uint) (*Sides, error) {
	var sides Sides
	result := h.db.First(&sides, id)
	return &sides, result.Error
}

func (h *SidesHandler) GetSides() ([]Sides, error) {
	var sides []Sides
	result := h.db.Find(&sides)
	return sides, result.Error
}

func (h *SidesHandler) UpdateSides(id uint, sides *Sides) error {
	result := h.db.Model(&Sides{}).Where("id = ?", id).Updates(sides)
	return result.Error
}

func (h *SidesHandler) DeleteSides(id uint) error {
	result := h.db.Delete(&Sides{}, id)
	return result.Error
}
