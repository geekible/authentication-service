package services

import (
	"authservice/src/config"
	"authservice/src/domain"
	"authservice/src/dtos"
	"authservice/src/helpers"
	"authservice/src/repositories"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo      *repositories.UserRepository
	userClaimRepo *repositories.UserClaimRepository
	emailService  *EmailService
	tokenAuth     *jwtauth.JWTAuth
	logger        *zap.SugaredLogger
	clientSecret  string
}

func InitUserService(serviceCfg *config.ServiceConfig) *UserService {
	return &UserService{
		userRepo:      repositories.InitUserRepositoy(serviceCfg),
		userClaimRepo: repositories.InitUserClaimRepository(serviceCfg),
		emailService:  InitEmailService(),
		tokenAuth:     jwtauth.New("HS256", []byte(serviceCfg.ClientSecret), nil),
		logger:        serviceCfg.Logger,
		clientSecret:  serviceCfg.ClientSecret,
	}
}

func (s *UserService) validateUser(user domain.User) error {
	if len(user.Username) <= 0 {
		return errors.New("username must be supplied")
	}

	if !s.emailService.ValidateEmail(user.EmailAddress) {
		return fmt.Errorf("the email address %s is not in a valid format", user.EmailAddress)
	}

	pwdService := helpers.InitPasswordHelper(user.Password)
	if err := pwdService.ValidateComplexity(); err != nil {
		return err
	}

	return nil
}

func (s *UserService) AddUser(user domain.User, isAdminUser bool) (domain.User, error) {
	if err := s.validateUser(user); err != nil {
		return user, err
	}

	pwd, err := helpers.InitCryptoHelper().Encrypt(user.Password)
	if err != nil {
		s.logger.Errorf("error encrypting password with error %v", user.Username, err)
		return user, err
	}
	user.Password = pwd

	user, err = s.userRepo.Add(user)
	if err != nil {
		s.logger.Errorf("error adding user %s with error %v", user.Username, err)
		return user, err
	}

	userClaim := domain.UserClaim{
		UserId: user.ID,
	}

	if isAdminUser {
		userClaim.ClaimId = 1
	} else {
		userClaim.ClaimId = 2
	}

	if err := s.userClaimRepo.Add(userClaim); err != nil {
		s.logger.Errorf("error adding user claim for user %s with error %v", user.Username, err)
		return user, err
	}

	return user, nil
}

func (s *UserService) UpdateUserDetails(user dtos.UserDto) error {
	usr, err := s.userRepo.GetById(user.UserId)
	if err != nil {
		return errors.New("details do not match")
	}

	usr.FirstName = user.FirstName
	usr.Surname = user.Surname
	usr.EmailAddress = user.EmailAddress
	usr.Username = user.Username

	if err := s.userRepo.UpdateUser(usr); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserPassword(updateUserPassword dtos.UserUpdatePasswordDto) error {
	user, err := s.userRepo.GetById(updateUserPassword.UserId)
	if err != nil {
		return errors.New("details do not match")
	}

	if !helpers.InitCryptoHelper().IsHashMatched(user.Password, updateUserPassword.OldPassword) {
		return errors.New("details do not match")
	}

	err = s.userRepo.UpdateUserPassword(updateUserPassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(user domain.User) error {
	if err := s.userRepo.Delete(user); err != nil {
		s.logger.Errorf("error deleting user %s with error %v", user.Username, err)
		return err
	}

	return nil
}

func (s *UserService) GetByUsername(username string) (dtos.UserDto, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		s.logger.Warnf("invalid user %s", username)
		return dtos.UserDto{}, errors.New("user not found")
	}

	return dtos.UserDto{
		UserId:       user.ID,
		Username:     user.Username,
		EmailAddress: user.EmailAddress,
		FirstName:    user.FirstName,
		Surname:      user.Surname,
	}, nil
}

func (s *UserService) GetByUsernameAndPassword(username, password string) (dtos.UserLoginResponseDto, error) {
	loginErrMsg := "username or password does not match"
	if len(username) <= 0 || len(password) <= 0 {
		s.logger.Warnf("invalid login attempt for user %s", username)
		return dtos.UserLoginResponseDto{}, errors.New(loginErrMsg)
	}

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		s.logger.Warnf("invalid login attempt for user %s", username)
		return dtos.UserLoginResponseDto{}, errors.New(loginErrMsg)
	}

	if !helpers.InitCryptoHelper().IsHashMatched(user.Password, password) {
		s.logger.Warnf("invalid login attempt for user %s", username)
		return dtos.UserLoginResponseDto{}, errors.New(loginErrMsg)
	}

	resp := dtos.UserLoginResponseDto{
		Username:     user.Username,
		EmailAddress: user.EmailAddress,
	}

	claims, err := s.userClaimRepo.GetClaimsByUserId(user.ID)
	if err != nil {
		return dtos.UserLoginResponseDto{}, err
	}
	resp.UserClaims = claims

	return resp, nil
}

func (s *UserService) GenerateUserToken(loginResponse dtos.UserLoginResponseDto) (string, error) {
	permissions := []string{}
	permissions = append(permissions, loginResponse.UserClaims...)

	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{
		"username":      loginResponse.Username,
		"email_address": loginResponse.EmailAddress,
		"exp":           time.Now().Add(1 * time.Hour),
		"issueed_at":    time.Now(),
		"permissions":   permissions,
	})

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) CustomJWTAuthVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("access_token")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.clientSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		expChecked := false

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			for key, value := range claims {
				//chaeck exp time
				if key == "exp" {
					expChecked = true
					exp := value.(float64)
					now := time.Now().Unix()

					if int64(exp) < now {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
				}
			}
		}

		if !expChecked {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token.Valid {
			next.ServeHTTP(w, r)
			return
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
