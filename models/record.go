package models

type Record struct {
	ID      uint64  `gorm:"primaryKey" json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Type    string  `json:"type"`
}
