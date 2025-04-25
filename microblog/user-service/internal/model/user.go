package model

type User struct {
	ID           string `gorm:"primaryKey;type:uuid"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
}
