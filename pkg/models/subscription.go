package models

import "time"

type Subscription struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Email             string     `gorm:"size:100;not null;uniqueIndex:idx_email_city" json:"email" binding:"required,email"`
	City              string     `gorm:"size:100;not null;uniqueIndex:idx_email_city" json:"city"  binding:"required"`
	Condition         string     `gorm:"size:255;not null"                json:"condition" binding:"required,condition"`
	Verified          bool       `gorm:"default:false" json:"verified"`
	VerificationToken string     `gorm:"size:64;index" json:"-"`
	TokenExpiresAt    *time.Time `json:"-"`
	LastSent          *time.Time `json:"last_sent"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
