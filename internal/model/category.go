package model

type Category struct {
	Base
	Name string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug string    `gorm:"type:varchar(100);unique" json:"slug"`
	Articles []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}