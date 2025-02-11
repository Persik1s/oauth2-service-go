package domain

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	UserId       uuid.UUID
	SessionToken string
	TTL          time.Duration
}
