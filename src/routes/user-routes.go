package routes

import (
	"authservice/src/helpers"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserRoutes struct {
	baseEndpoint string
	mux          *chi.Mux
	jsonHelpers  *helpers.JsonHelpers
	logger       *zap.SugaredLogger
}
