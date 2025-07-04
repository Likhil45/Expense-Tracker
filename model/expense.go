package model

import "time"

type Expense struct {
	Id          string  `gorm:"primaryKey"`
	User_id     string  `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
	Currency    string  `gorm:"default:USD;not null"`
	Category    string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	TimeStamp   time.Time
}
