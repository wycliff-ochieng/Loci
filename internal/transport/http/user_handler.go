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

func NewUserHandler(l *slog.Logger, us *service.UserService) *UserHandler {
	return &UserHandler{
		logger: l,
		us:     us,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Register user handler touched ,sending data to service layer")

	err := json.NewDecoder(r.Body).Decode()
}
