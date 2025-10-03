package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wycliff-ochieng/internal/service"
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

	//err := json.NewDecoder(r.Body).Decode()
}
