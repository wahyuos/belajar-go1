package model

type File struct {
	Base
	Filename  string    `gorm:"type:varchar(255);not null" json:"filename"`
	Filetype  string    `gorm:"type:varchar(50);not null" json:"filetype"`
	Filepath  string    `gorm:"type:varchar(255);not null" json:"filepath"`
}