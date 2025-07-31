package models

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}
