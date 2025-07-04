package model

type User struct {
	UserId   string `gorm:"primaryKey"`
	UserName string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}
