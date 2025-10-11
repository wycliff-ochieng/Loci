package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/wycliff-ochieng/internal/models"
	"github.com/wycliff-ochieng/internal/service"
	"github.com/wycliff-ochieng/pkg/middleware"
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

type location struct {
	Lat  float64
	Long float64
}

type CreateLociReq struct {
	UserID   uuid.UUID
	Message  string
	Location location
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

	if RegisterUser.Username == "" || RegisterUser.Firstname == "" || RegisterUser.Lastname == "" || RegisterUser.Email == "" {
		http.Error(w, "some requireed fields are missing values", http.StatusExpectationFailed)
		return
	}

	user, err := h.us.Register(ctx, RegisterUser.Username, RegisterUser.Firstname, RegisterUser.Lastname, RegisterUser.Email, RegisterUser.Password)
	if err != nil {
		http.Error(w, "service layer failure: ", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
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

	var loci *models.LociResponse

	ctx := r.Context()

	var req CreateLociReq

	userID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		http.Error(w, "Could get userID from context", http.StatusExpectationFailed)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&loci)
	if err != nil {
		http.Error(w, "error decoding message from json", http.StatusInternalServerError)
		return
	}

	//validate message

}
