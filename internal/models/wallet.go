package models

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Balance   int64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
