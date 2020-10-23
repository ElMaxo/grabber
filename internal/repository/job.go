package repository

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Job struct {
	CreatedAt time.Time `json:"-"`
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Url       string    `gorm:"not null"`
	Query     postgres.Jsonb
	Period    int64 `gorm:"not null"`
}
