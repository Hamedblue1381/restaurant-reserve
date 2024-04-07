package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID         uint      `gorm:"primaryKey"`
	FoodID     uint      // Foreign key for Food
	Food       Food      `json:"food"` // Food relationship
	UserID     uint      // Foreign key for User
	User       User      `json:"user"` // User relationship
	SideID     uint      // Foreign key for Sides
	Side       Sides     `json:"side"` // Sides relationship
	Date       time.Time `json:"date"`
	IsPaid     bool      `json:"-"`
	gorm.Model `json:"-" swaggerignore:"true"`
}

type ReservationHandler struct {
	db *gorm.DB
}

func NewReservationHandler(db *gorm.DB) *ReservationHandler {
	return &ReservationHandler{db}
}

func (r *ReservationHandler) Reserve(reservation *Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *ReservationHandler) DeleteReservation(id uint) error {
	result := r.db.Delete(&Reservation{}, id)
	return result.Error
}

func (r *ReservationHandler) UpdateReservation(id uint, reservation *Reservation) error {
	result := r.db.Model(&Reservation{}).Where("id = ?", id).Updates(reservation)
	return result.Error
}

func (r *ReservationHandler) IsBlackListed(id uint) error {
	var count int64
	result := r.db.Model(&Reservation{}).Where("user_id = ? AND is_paid = ?", id, false).Count(&count)

	if count > 3 {
		return errors.New("User is blacklisted due to having more than 3 unpaid reservations")
	}

	return result.Error
}

func (r *ReservationHandler) ListReservations(startDate, endDate time.Time) ([]Reservation, error) {
	var reservations []Reservation
	query := r.db

	if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	}

	if !endDate.IsZero() {
		query = query.Where("date <= ?", endDate)
	}

	result := query.Preload("User").Preload("Food").Preload("Side").Find(&reservations)
	return reservations, result.Error
}

func (r *ReservationHandler) GetReservation(id uint) (*Reservation, error) {
	var reservation Reservation
	result := r.db.First(&reservation, id)
	return &reservation, result.Error
}

func (r *ReservationHandler) GetReservationsByUserID(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	result := r.db.Preload("User").Preload("Restaurant").Where("user_id = ?", userID).Find(&reservations)

	if result.Error != nil {
		return nil, result.Error
	}
	return reservations, nil
}
