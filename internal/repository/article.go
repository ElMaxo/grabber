package repository

import "time"

type Article struct {
	CreatedAt   time.Time `json:"-"`
	ID          string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Link        string
	Title       string
	Description string
	RowsCount   int `gorm:"-"`
}
