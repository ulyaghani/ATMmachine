package models

import "time"

type Account struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Pin       string    `gorm:"type:varchar(10);not null" json:"pin"`
	Balance   float64   `gorm:"type:decimal(15,2);default:0" json:"balance"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Account) TableName() string {
	return "accounts"
}
