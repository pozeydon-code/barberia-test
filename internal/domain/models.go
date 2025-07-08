package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound       = errors.New("record not found")
	ErrorInvalidInput = errors.New("invalid input provided")
)

type Appointment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ClientName string    `gorm:"size:255;not null" json:"client_name"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
	Products   []Product `gorm:"many2many:appointment_products;" json:"products"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Price       float64   `gorm:"not null" json:"price"`
	Description string    `gorm:"size:500" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
