package model

import "github.com/google/uuid"

type Article struct {
	Base
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	CategoryID uuid.UUID `gorm:"type:uuid" json:"category_id"`
	Category   *Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
	// Tambahkan relasi ke User (Many-to-One)
	UserID      uuid.UUID `gorm:"type:char(36);not null" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user"`
}