package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Telephone  string `json:"telephone"`
	Role       string `json:"role"`
	Password   string `json:"password"`
	gorm.Model `json:"-" swaggerignore:"true"`
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

func (h *UserHandler) CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return h.db.Create(user).Error
}

func (h *UserHandler) CheckPassword(email, password string) bool {
	var user User
	if err := h.db.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (h *UserHandler) GetUser(id uint) (*User, error) {
	var user User
	result := h.db.First(&user, id)
	return &user, result.Error
}

func (h *UserHandler) GetUsers() ([]User, error) {
	var users []User
	result := h.db.Find(&users)
	return users, result.Error
}

func (h *UserHandler) UpdateUser(id uint, user *User) error {
	result := h.db.Model(&User{}).Where("id = ?", id).Updates(user)
	return result.Error
}

func (h *UserHandler) DeleteUser(id uint) error {
	result := h.db.Delete(&User{}, id)
	return result.Error
}

func (h *UserHandler) GetUserByEmail(email string) (*User, error) {
	var user User
	result := h.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, result.Error
}
