package routes

import (
	"authservice/src/config"
	"authservice/src/domain"
	"authservice/src/helpers"
	"authservice/src/services"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ClaimsRoutes struct {
	baseEndpoint string
	mux          *chi.Mux
	claimService *services.ClaimService
	jsonHelpers  *helpers.JsonHelpers
	logger       *zap.SugaredLogger
}

const (
	claimErrSrc = "ClaimRoutes"
)

func InitClaimRoutes(serviceCfg *config.ServiceConfig) *ClaimsRoutes {
	return &ClaimsRoutes{
		baseEndpoint: "/claims",
		mux:          serviceCfg.Mux,
		claimService: services.InitClaimService(serviceCfg),
		jsonHelpers:  helpers.InitJsonHelpers(serviceCfg.Logger),
		logger:       serviceCfg.Logger,
	}
}

func (a *ClaimsRoutes) Register() {
	a.mux.Post(a.baseEndpoint, a.addClaim)
	a.mux.Put(a.baseEndpoint, a.updateClaim)
	a.mux.Delete(a.baseEndpoint, a.deleteClaim)

	a.mux.Get(fmt.Sprintf("%s/get-all", a.baseEndpoint), a.getAll)
}

func (a *ClaimsRoutes) addClaim(w http.ResponseWriter, r *http.Request) {
	var claim domain.Claim
	if err := a.jsonHelpers.ReadJSON(w, r, &claim); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, claimErrSrc)
		return
	}

	claim, err := a.claimService.Add(claim)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, claimErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusCreated, claim)
}

func (a *ClaimsRoutes) updateClaim(w http.ResponseWriter, r *http.Request) {
	var claim domain.Claim
	if err := a.jsonHelpers.ReadJSON(w, r, &claim); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, claimErrSrc)
		return
	}

	err := a.claimService.Update(claim)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, claimErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusAccepted, claim)
}

func (a *ClaimsRoutes) deleteClaim(w http.ResponseWriter, r *http.Request) {
	var claim domain.Claim
	if err := a.jsonHelpers.ReadJSON(w, r, &claim); err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusBadRequest, claimErrSrc)
		return
	}

	err := a.claimService.Delete(claim)
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, claimErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusAccepted, claim)
}

func (a *ClaimsRoutes) getAll(w http.ResponseWriter, r *http.Request) {
	claims, err := a.claimService.GetAll()
	if err != nil {
		a.jsonHelpers.ErrorJSON(w, err, http.StatusInternalServerError, claimErrSrc)
		return
	}

	a.jsonHelpers.WriteJSON(w, http.StatusAccepted, claims)
}
