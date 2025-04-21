package models

import "time"

type Weather struct {
	City        string    `gorm:"primaryKey" json:"city"        binding:"required"       validate:"required"`
	Temperature float64   `json:"temperature"                  validate:"required"`
	Humidity    int       `json:"humidity"                     binding:"required"       validate:"required,gte=0,lte=100"`
	Condition   string    `json:"condition"                    binding:"required"       validate:"required"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}
