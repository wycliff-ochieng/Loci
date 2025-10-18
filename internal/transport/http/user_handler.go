package http

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wycliff-ochieng/internal/models"
	"github.com/wycliff-ochieng/internal/service"
	"github.com/wycliff-ochieng/pkg/middleware"
	"github.com/wycliff-ochieng/sqlc"
)

type UserHandler struct {
	logger *slog.Logger
	us     *service.UserService
}

type RegisterReq struct {
	Username  string
	Firstname string
	Lastname  string
	Email     string
	Password  string
}

type LoginReq struct {
	Email    string
	Username string
	Password string
}

type location struct {
	Lat  float64
	Long float64
}

type CreateLociReq struct {
	UserID   uuid.UUID `json:"userID"`
	Message  string    `json:"message"`
	Location location  `json:"loaction"`
}

type AuthenticationResponse struct {
	User         interface{}
	AccessToken  string
	RefreshToken string
}

func NewUserHandler(l *slog.Logger, us *service.UserService) *UserHandler {
	return &UserHandler{
		logger: l,
		us:     us,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Register user handler touched ,sending data to service layer")

	ctx := r.Context()

	var RegisterUser *RegisterReq

	err := json.NewDecoder(r.Body).Decode(&RegisterUser)
	if err != nil {
		http.Error(w, "failed to register user for some reason", http.StatusExpectationFailed)
		return
	}

	if RegisterUser.Username == "" || RegisterUser.Firstname == "" || RegisterUser.Lastname == "" || RegisterUser.Email == "" || RegisterUser.Password == "" {
		http.Error(w, "some requireed fields are missing values", http.StatusExpectationFailed)
		return
	}

	user, err := h.us.Register(ctx, RegisterUser.Username, RegisterUser.Firstname, RegisterUser.Lastname, RegisterUser.Email, RegisterUser.Password)
	if err != nil {
		log.Printf("Error due to: %s", err)
		http.Error(w, "service layer failure: ", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error picking up json response", http.StatusInternalServerError)
		return
	}

	//validate user input
	if req.Email == " " || req.Password == " " {
		http.Error(w, "email or password required", http.StatusExpectationFailed)
		return
	}

	//authenticate user (user service transactions)
	//token, user, err :=

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	token, user, err := h.us.LoginUser(ctx, req.Username, req.Email, req.Password)
	if err == service.ErrNotFound || err == service.ErrInvalidPassword {
		log.Printf("error due to: %s", err)
		http.Error(w, "USER NOT FOUND,INVALID PASSWORD", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "FAILED TO SIGN IN", http.StatusInternalServerError)
		//h.logger.Info("reason: due to")
		log.Printf("due to: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthenticationResponse{
		User:         user,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

// api :: GET -> api/loci?{SouthWestlat=}&{}&{}&{}
func (h *UserHandler) GetLociInGeoFencedLocation(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Messages within geofenced location being viewed")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//geeting the radius/ geolocation points from the client
	vars := mux.Vars(r)
	SWLatStr := vars["SouthWestlat"]
	SWLongStr := vars["SouthWestLong"]
	NELatStr := vars["NorthEastLat"]
	NELongStr := vars["NorthEastLong"]

	//parse the string to float
	SWLat, err1 := strconv.ParseFloat(SWLatStr, 64)
	SWLong, err2 := strconv.ParseFloat(SWLongStr, 64)
	NELat, err3 := strconv.ParseFloat(NELatStr, 64)
	NELong, err4 := strconv.ParseFloat(NELongStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Error(w, "Error converting the coordinates to float", http.StatusFailedDependency)
		return
	}

	box := models.BoundBox{
		NorthEastLat:  NELat,
		NorthEastLong: NELong,
		SouthWestLat:  SWLat,
		SouthWestLong: SWLong,
	}

	//call the service layer
	AllLoci, err := h.us.GetLociWithinBounds(ctx, box)
	if err != nil {
		http.Error(w, "some error while fetching loci within the geolocation from the service layer", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&AllLoci)

}

func (h *UserHandler) CreateLoci(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("create message handler in action")

	//var loci *models.LociResponse

	var req CreateLociReq

	ctx := r.Context()

	userID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		http.Error(w, "Could get userID from context", http.StatusExpectationFailed)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error decoding message from json", http.StatusInternalServerError)
		return
	}
	//validation
	if req.Message == "" || req.Location.Lat == 0.0 || req.Location.Long == 0.0 || len(req.Message) > 250 {
		http.Error(w, "invalid data input types,check location points/message length", http.StatusExpectationFailed)
		return

	}

	serviceParam := sqlc.CreateLociParams{
		UserID:   userID,
		Message:  req.Message,
		Location: req.Location,
	}

	//call user service
	locus, err := h.us.CreateLoci(ctx, userID, serviceParam)
	if err != nil {
		http.Error(w, "something happened in the service layer", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Tpe", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&locus)

}
