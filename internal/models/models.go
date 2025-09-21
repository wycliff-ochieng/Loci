package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	UserID    uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	About     string    `json:"about"`
	Password  string    `json:"password"`
}

type Loci struct {
	LociID         uuid.UUID
	UserID         uuid.UUID
	Message        string
	Location       point.Point
	CreatedAT      time.Time
	AuthorUserName string
}
