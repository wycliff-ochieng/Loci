package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	UserID    uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	About     string    `json:"about"`
	Password  string    `json:"password"`
}

type LociResponse struct {
	LociID  uuid.UUID
	UserID  uuid.UUID
	Message string
	//Location       point.Point
	CreatedAT      time.Time
	AuthorUserName string
}

type Locus struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	message    string
	location   float64
	createdat  time.Time
	viewscount int64
}

type BoundBox struct {
	NorthEastLat  float64
	NorthEastLong float64
	SouthWestLat  float64
	SouthWestLong float64
}

type Metadata struct {
	Location string `json:"location"`
}

type User struct {
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
}

func NewUser(id int, username string, firstname string, lastname string, email string, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("HASH_ERROR:Failed to harsh password:%v", err)
	}
	return &User{
		ID:        id,
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Password:  string(hashedPassword),
	}, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
