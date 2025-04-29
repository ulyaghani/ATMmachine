package models

import "time"

type Transaction struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID int       `gorm:"not null" json:"account_id"`
	Type      string    `gorm:"type:enum('deposit', 'withdraw', 'transfer_in', 'transfer_out');not null" json:"type"`
	Amount    float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	TargetID  *int      `gorm:"type:int;default:null" json:"target_id"` // Nullable
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
