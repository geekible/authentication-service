package routes

import (
	"authservice/src/config"
	"authservice/src/domain"
	"authservice/src/dtos"
	"authservice/src/helpers"
	"authservice/src/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	userErrSrc = "UserRoutes"
)

type UserRoutes struct {
	baseEndpoint string
	mux          *chi.Mux
	userService  *services.UserService
	jsonHelpers  *helpers.JsonHelpers
	logger       *zap.SugaredLogger
}

func InitUserRoutes(serviceCfg *config.ServiceConfig) *UserRoutes {
	return &UserRoutes{
		baseEndpoint: "/user",
		mux:          serviceCfg.Mux,
		userService:  services.InitUserService(serviceCfg),
		jsonHelpers:  helpers.InitJsonHelpers(serviceCfg.Logger),
		logger:       serviceCfg.Logger,
	}
}

func (a *UserRoutes) Register() {
	a.mux.Post(a.baseEndpoint, a.addUser)

	// protected routes
	a.mux.Post(fmt.Sprintf("%s/add-admin-user", a.baseEndpoint), a.addAdminUser)
	a.mux.Put(fmt.Sprintf("%s/update-password", a.baseEndpoint), a.updateUserPassword)
	a.mux.Delete(a.baseEndpoint, a.deleteUser)

	a.mux.Post(fmt.Sprintf("%s/login", a.baseEndpoint), a.login)
}

func (a *UserRoutes) addUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := a.jsonHelpers.ReadJSON(w, r, &user); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, userErrSrc)
		return
	}

	user, err := a.userService.AddUser(user, false)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, userErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusCreated, nil, nil)
}

func (a *UserRoutes) addAdminUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := a.jsonHelpers.ReadJSON(w, r, &user); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, userErrSrc)
		return
	}

	user, err := a.userService.AddUser(user, true)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, userErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusCreated, nil, nil)
}

func (a *UserRoutes) updateUserPassword(w http.ResponseWriter, r *http.Request) {
	var updatePasswordDto dtos.UserUpdatePasswordDto
	if err := a.jsonHelpers.ReadJSON(w, r, &updatePasswordDto); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, userErrSrc)
		return
	}

	if err := a.userService.UpdateUserPassword(updatePasswordDto); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, userErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusAccepted, nil, nil)
}

func (a *UserRoutes) deleteUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := a.jsonHelpers.ReadJSON(w, r, &user); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, userErrSrc)
		return
	}

	if err := a.userService.DeleteUser(user); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, userErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusAccepted, nil, nil)
}

func (a *UserRoutes) login(w http.ResponseWriter, r *http.Request) {
	var loginDto dtos.LoginDto
	if err := a.jsonHelpers.ReadJSON(w, r, &loginDto); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, userErrSrc)
		return
	}

	user, err := a.userService.GetByUsernameAndPassword(loginDto.Username, loginDto.Password)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, errors.New("invalid login attempt"), http.StatusUnauthorized, userErrSrc)
		return
	}

	token, err := a.userService.GenerateUserToken(user)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, errors.New("invalid login attempt"), http.StatusUnauthorized, userErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusOK, token, nil)
}
