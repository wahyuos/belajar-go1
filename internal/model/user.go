package model

// User merepresentasikan model pengguna di database
type User struct {
	Base
	Username string `gorm:"type:varchar(100);uniqueIndex" json:"username"`
	Password string `gorm:"type:varchar(255)" json:"-"`

	// Tambahkan relasi One-to-Many ke Article (GORM)
	Articles []Article `gorm:"foreignKey:UserID"`
}