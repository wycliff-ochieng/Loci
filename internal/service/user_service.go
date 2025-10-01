package service

import (
	"context"

	"github.com/wycliff-ochieng/internal/models"
	"github.com/wycliff-ochieng/internal/store"
	//"github.com/wycliff-ochieng/internal/transport/http"
)

type UserService struct {
	db store.Postgis
	//uh http.UserHandler
}

func NewUserService(db store.Postgis) *UserService {
	return &UserService{
		db: db,
	}
}

func (us *UserService) Register(ctx context.Context, firstname string, lastname string, email string, password string) (*models.Users, error) {
	return nil,nil
}

func(us *UserService)Login()
