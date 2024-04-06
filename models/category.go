package models

import "gorm.io/gorm"

type Category struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `json:"name"`
	Foods      []Food `gorm:"foreignKey:CategoryID"` // Foods relationship
	gorm.Model `json:"-" swaggerignore:"true"`
}
