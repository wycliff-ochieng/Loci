package models

import (
	"fmt"
	"math"
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
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Message      string    `json:"message"`
	Location     GeoPoint  `json:"location"`
	Createdat    time.Time `json:"created_at"`
	Viewscount   int64     `json:"view_count"`
	Repliescount int64     `json:"replies_count"`
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

type GeoPoint struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type View struct {
	UserID   uuid.UUID
	LocusID  uuid.UUID
	ViewedAT time.Time
}

type Reply struct {
	ReplyID   uuid.UUID `json:"replyid"`
	LocusID   uuid.UUID `json:"locusid"`
	UserID    uuid.UUID `json:"userid"`
	Content   string    `json:"content"`
	CreatedAT time.Time `json:"createdat"`
}

type ReplyEvent struct {
	ReplyID       uuid.UUID `json:"replyid"`
	LocusID       uuid.UUID `json:"locusid"`
	UserName      string    `json:"username"`
	LocusLocation GeoPoint  `json:"locus_location"`
	CreatedAt     time.Time `json:"createdat"`
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

func ComparePassword(harshedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(harshedPassword), []byte(password))
}

const (
	earthRadiusKM = 6471
)

func CalculateDistance(pt1, pt2 *GeoPoint) float64 {
	//implement Haverside formula ( converting degrees to radians)
	//apply the haverside formula to the distance
	// return the clculated distance
	lat1Rad := pt1.Lat * math.Pi / 180
	lon1Rad := pt1.Long * math.Pi / 180
	lat2Rad := pt2.Lat * math.Pi / 180
	long2Rad := pt2.Long * math.Pi / 180

	diffLat := lat1Rad - lat2Rad
	diffLong := lon1Rad - long2Rad

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(diffLong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKM * c

	return distance
}

type ViewEvent struct {
	UserID        uuid.UUID `json:"userid"`
	LocusID       uuid.UUID `json:"locusid"`
	LocusLocation GeoPoint  `json:"locus_location"`
	ViewedAt      time.Time `json:"viewedat"`
}
