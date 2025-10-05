package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	ErrNotFound    = errors.New("No user with that email/ username")
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

func (us *UserService) LoginUser(ctx context.Context, username string, email string, password string) (*models.UserResponse, error) {

	//sqlc.LoginParams
	var User models.User

	user, err := us.query.Login(ctx, sqlc.LoginParams{Email: email, Username: username})
	if err != nil || err == sql.ErrNoRows {
		return nil, ErrNotFound
	}

	if err := User.ComparePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password!!Try again")
	}

	return &models.UserResponse{
		UserID:    user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Email:     user.Email,
	}, nil

}

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
