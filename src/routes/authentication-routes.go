package routes

import (
	"authservice/src/config"
	"authservice/src/dtos"
	"authservice/src/helpers"
	"authservice/src/services"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AuthenticationRoutes struct {
	baseEndpoint string
	mux          *chi.Mux
	authService  *services.AuthenticationService
	jsonHelpers  *helpers.JsonHelpers
	logger       *zap.SugaredLogger
}

const (
	authErrSrc = "AuthenticationRoutes"
)

func InitAuthenticationRoutes(serviceConfig *config.ServiceConfig) *AuthenticationRoutes {
	return &AuthenticationRoutes{
		mux:          serviceConfig.Mux,
		baseEndpoint: "/authentication",
		authService:  services.InitAuthenticationService(serviceConfig),
		jsonHelpers:  helpers.InitJsonHelpers(serviceConfig.Logger),
		logger:       serviceConfig.Logger,
	}
}

func (a *AuthenticationRoutes) Register() {
	a.mux.Post(fmt.Sprintf("%s/login", a.baseEndpoint), a.login)
}

func (a *AuthenticationRoutes) login(w http.ResponseWriter, r *http.Request) {
	var loginDto dtos.LoginDto
	if err := a.jsonHelpers.ReadJSON(w, r, &loginDto); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, authErrSrc)
		return
	}

	_, tokenString, err := a.authService.Login(loginDto.Username, loginDto.Password)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, authErrSrc)
		return
	}

	a.logger.Infof("user %s granted access at ", loginDto.Username, time.Now())
	a.jsonHelpers.WriteJSON(w, http.StatusOK, tokenString)
}
