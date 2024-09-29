package services

import (
	"authservice/src/config"
	"authservice/src/domain"
	"authservice/src/dtos"
	"authservice/src/helpers"
	"authservice/src/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo      *repositories.UserRepository
	userClaimRepo *repositories.UserClaimRepository
	emailService  *EmailService
	tokenAuth     *jwtauth.JWTAuth
	logger        *zap.SugaredLogger
}

func InitUserService(serviceCfg *config.ServiceConfig) *UserService {
	return &UserService{
		userRepo:      repositories.InitUserRepositoy(serviceCfg),
		userClaimRepo: repositories.InitUserClaimRepository(serviceCfg),
		emailService:  InitEmailService(),
		tokenAuth:     jwtauth.New("HS256", []byte(serviceCfg.ClientSecret), nil),
		logger:        serviceCfg.Logger,
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
	loginResponse.Exp = time.Now().Add(8 * time.Hour)
	_, tokenString, err := s.tokenAuth.Encode(map[string]interface{}{
		"profile": loginResponse,
	})

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
