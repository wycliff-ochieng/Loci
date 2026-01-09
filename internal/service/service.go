package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/wycliff-ochieng/internal/limitter"
	"github.com/wycliff-ochieng/internal/models"
	"github.com/wycliff-ochieng/internal/socket"
	"github.com/wycliff-ochieng/internal/store"
	"github.com/wycliff-ochieng/pkg/auth"
	"golang.org/x/crypto/bcrypt"

	//"github.com/wycliff-ochieng/pkg/middleware/auth"
	"github.com/wycliff-ochieng/sqlc"
)

type UserService struct {
	db    *store.Postgis
	query sqlc.Queries
	rtl   *limitter.RedisLimitter
	hub   *socket.Hub
	//lc sqlc.CreateLociParams
	//uh http.UserHandler
}

type RateLimitAction string

const (
	ActionPostLocus = "post_locus"
)

var (
	ErrForbidden       = errors.New("Cannot perform this operation")
	ErrOutOfBounds     = errors.New("Error setting bounds")
	ErrNotFound        = errors.New("No user with that email/ username")
	ErrRateLimited     = errors.New("user not allowed to create new locus, maximum retires exceeded")
	ErrInvalidPassword = errors.New("wrong password input")
)

func NewUserService(db *store.Postgis, query sqlc.Queries, rtl *limitter.RedisLimitter, hub *socket.Hub) *UserService {
	return &UserService{
		db:    db,
		query: query,
		rtl:   rtl,
		hub:   hub,
	}
}

func (us *UserService) Register(ctx context.Context, username string, firstname string, lastname string, email string, password string) (*models.UserResponse, error) {
	//var exists bool

	exist, err := us.query.UserExists(ctx, email)
	if err != nil {
		log.Printf("User Exists Error: %s", err)
		return nil, err
	}

	if exist {
		return nil, ErrForbidden
	}

	harshedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Generate password error: %s", err)
		return nil, err
	}

	newUserID := uuid.New()

	err = us.query.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           newUserID,
		Username:     username,
		Firstname:    firstname,
		Lastname:     lastname,
		Email:        email,
		PasswordHash: string(harshedPassword),
		//CreatedAt:    createdat,
	})
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		UserID:    newUserID,
		Firstname: firstname,
		Lastname:  lastname,
		Username:  username,
		Email:     email,
	}, nil
}

func (us *UserService) LoginUser(ctx context.Context, username string, email string, password string) (*auth.TokenPair, *models.UserResponse, error) {

	//sqlc.LoginParams
	//var User models.User
	log.Println("user service login function touched")

	user, err := us.query.Login(ctx, sqlc.LoginParams{Email: email, Username: username})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, err
		}
		log.Printf("error due to: %s", err)
		return nil, nil, err
	}

	//if err := User.ComparePassword(password); err != nil {
	//	return nil, nil, fmt.Errorf("invalid password!!Try again")
	//}

	if err := models.ComparePassword(user.PasswordHash, password); err != nil {
		log.Printf("comparing password error: %s", err)
		return nil, nil, fmt.Errorf("invalid password")
	}

	token, err := auth.GenerateTokenPair(
		user.ID,
		//user.UserID,
		//role,
		user.Email,
		string(auth.JWTsecret),     //jwtsecret
		string(auth.Refreshsecret), //refreshsecret
		time.Hour*24,
		time.Hour*24*7,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return token, &models.UserResponse{
		UserID:    user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Email:     user.Email,
	}, nil

}

func (us *UserService) GetLociWithinBounds(ctx context.Context, bounds models.BoundBox) ([]sqlc.Loci, error) {
	log.Println("Database operation in motion")

	log.Printf("[SERVICE] Received BoundingBox: %+v", bounds)

	messages, err := us.query.GetLociInBounds(ctx, sqlc.GetLociInBoundsParams{
		StMakeenvelope:   bounds.SouthWestLong,
		StMakeenvelope_2: bounds.SouthWestLat,
		StMakeenvelope_3: bounds.NorthEastLong,
		StMakeenvelope_4: bounds.NorthEastLat,
	})
	if err != nil {
		return nil, err
	}

	return messages, nil

}

func (us *UserService) CreateLoci(ctx context.Context, userID uuid.UUID, params sqlc.CreateLociParams) ([]sqlc.CreateLociRow, error) {

	//message := "test message"
	//location :=

	identifier := fmt.Sprintf("%d", userID)

	redisKey := generateRedisKey(ActionPostLocus, identifier)

	allowed, err := us.rtl.AllowPost(ctx, redisKey)
	if err != nil {
		log.Printf("Error in the rate limit check due to %s", err)
		return nil, err
	}

	if !allowed {
		return nil, ErrRateLimited
	}

	//content moderation logic -> check if message contains only allowed words
	newUserID := uuid.New()

	dbparams := sqlc.CreateLociParams{
		ID:            newUserID,
		UserID:        params.UserID,
		Message:       params.Message,
		StMakepoint:   params.StMakepoint,
		StMakepoint_2: params.StMakepoint_2,
		//Location: params.Location,
	}
	//calling the db
	loci, err := us.query.CreateLoci(ctx, dbparams)
	if err != nil {
		return nil, err
	}

	//delegeation to websocket
	//trigger real time broadcast
	if len(loci) > 0 {
		row := loci[0]

		//convert sqlc row to shared schema(models)
		newLoci := &models.Locus{
			ID:         row.ID,
			UserID:     row.UserID,
			Message:    row.Message,
			Location:   models.GeoPoint{},
			Createdat:  row.CreatedAt,
			Viewscount: int64(row.ViewCount),
		}

		//push loci to hub
		us.hub.BroadcastLocus <- newLoci
	}

	return loci, nil
}

func generateRedisKey(action RateLimitAction, identifier string) string {
	return fmt.Sprintf("ratelimit:%s:%s", action, identifier)
}

func (us *UserService) RecordView(ctx context.Context, userID uuid.UUID, locusID uuid.UUID) (*models.View, error) {

	//start transaction
	txs, err := us.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Printf("failed to set up transactions with pgx: %s", err)
	}

	defer txs.Rollback(ctx)

	qtx := us.query.WithTx(txs)

	viewValues := sqlc.CreateViewParams{
		UserID:  userID,
		LocusID: locusID,
	}

	/*locusView, err := us.query.CreateView(ctx, viewValues)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		return err
	}*/

	locusView, err := qtx.CreateView(ctx, viewValues)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.View{}, nil
		}
		return nil, err
	}

	err = qtx.IncrementViewCount(ctx, locusView.LocusID)
	if err != nil {
		return nil, err
	}

	return &models.View{
		UserID:   locusView.UserID,
		LocusID:  locusView.LocusID,
		ViewedAT: time.Now().Local(),
	}, err
}
