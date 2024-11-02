package services

import (
	"authservice/src/config"
	"authservice/src/helpers"
	"authservice/src/repositories"
	"errors"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

type AuthenticationService struct {
	clientId     string
	clientSecret string
	tokenAuth    *jwtauth.JWTAuth
	userRepo     *repositories.UserRepository
	cryptoHelper *helpers.CryptoHelper
	logger       *zap.SugaredLogger
}

func InitAuthenticationService(serviceConfig *config.ServiceConfig) *AuthenticationService {
	return &AuthenticationService{
		clientId:     serviceConfig.ClientId,
		tokenAuth:    jwtauth.New("HS256", []byte(serviceConfig.ClientSecret), nil),
		userRepo:     repositories.InitUserRepositoy(serviceConfig),
		cryptoHelper: helpers.InitCryptoHelper(),
		logger:       serviceConfig.Logger,
	}
}

func (s *AuthenticationService) Login(username, password string) (jwt.Token, string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if !s.cryptoHelper.IsHashMatched(user.Password, password) {
		s.logger.Warnf("invalid login attempt for user %s", username)
		return nil, "", errors.New("username or password does not not")
	}

	return s.tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
}
