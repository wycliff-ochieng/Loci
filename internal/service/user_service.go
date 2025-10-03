package service

import (
	"context"
	"errors"
	"log"

	"github.com/wycliff-ochieng/internal/models"
	"github.com/wycliff-ochieng/internal/store"
	"github.com/wycliff-ochieng/sqlc"
	//"github.com/wycliff-ochieng/internal/tra"
)

type UserService struct {
	db    *store.Postgis
	query sqlc.Queries
	//uh http.UserHandler
}

var (
	ErrForbidden   = errors.New("Cannot perform this operation")
	ErrOutOfBounds = errors.New("Error setting bounds")
)

func NewUserService(db *store.Postgis, query sqlc.Queries) *UserService {
	return &UserService{
		db:    db,
		query: query,
	}
}

func (us *UserService) Register(ctx context.Context, username string, firstname string, lastname string, email string, password string) (*models.UserResponse, error) {
	//var exists bool

	exist, err := us.query.UserExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, ErrForbidden
	}

	err = us.query.CreateUser(ctx, sqlc.CreateUserParams{
		//ID:           ID,
		Username:     username,
		Firstname:    firstname,
		Lastname:     lastname,
		Email:        email,
		PasswordHash: password,
		//CreatedAt:    createdat,
	})
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		//UserID:    ID,
		Firstname: firstname,
		Lastname:  lastname,
		Username:  username,
		Email:     email,
	}, nil
}

/*func (us *UserService) Login(ctx context.Context, email string, password string) (*models.UserResponse, error) {

	user, err := us.query.Login(ctx, email)
	if err != nil {
		return nil, err
	}

}
*/

func (us *UserService) GetLociWithinBounds(ctx context.Context, bounds models.BoundBox) ([]sqlc.GetLociInBoundsRow, error) {
	log.Println("Database operation in motion")

	messages, err := us.query.GetLociInBounds(ctx, sqlc.GetLociInBoundsParams{
		StMakeenvelope:   bounds.NorthEastLat,
		StMakeenvelope_2: bounds.NorthEastLong,
		StMakeenvelope_3: bounds.SouthWestLat,
		StMakeenvelope_4: bounds.SouthWestLong,
	})
	if err != nil {
		return nil, err
	}

	return messages, nil

}
