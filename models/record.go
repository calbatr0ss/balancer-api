package models

type Record struct {
	ID      uint64 `gorm:"primaryKey"`
	Name    string
	Balance float64
	Type    string // enum?
}
