package services

import (
	"authservice/src/config"
	"authservice/src/repositories"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type AuthenticationService struct {
	clientId     string
	clientSecret string
	tokenAuth    *jwtauth.JWTAuth
	userRepo     *repositories.UserRepository
}

func InitAuthenticationService(serviceConfig *config.ServiceConfig) *AuthenticationService {
	return &AuthenticationService{
		clientId:  serviceConfig.ClientId,
		tokenAuth: jwtauth.New("HS256", []byte(serviceConfig.ClientSecret), nil),
		userRepo:  repositories.InitUserRepositoy(serviceConfig.Db, serviceConfig.Logger),
	}
}

func (s *AuthenticationService) Login(username, password string) (jwt.Token, string, error) {
	user, err := s.userRepo.GetByUsernameAndPassword(username, password)
	if err != nil {
		return nil, "", err
	}

	return s.tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
}
