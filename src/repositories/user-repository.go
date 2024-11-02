package repositories

import (
	"authservice/src/config"
	"authservice/src/domain"
	"authservice/src/dtos"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func InitUserRepositoy(serviceCfg *config.ServiceConfig) *UserRepository {
	return &UserRepository{
		db:     serviceCfg.Db,
		logger: serviceCfg.Logger,
	}
}

func (r *UserRepository) Add(user domain.User) (domain.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		r.logger.Errorf("error creating user: %s with error: %v", user.Username, err)
		return user, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserPassword(updatePasswordDto dtos.UserUpdatePasswordDto) error {
	err := r.db.Model(&domain.User{}).
		Where("user_id = ?", updatePasswordDto.UserId).
		Update("password", updatePasswordDto.NewPassword).
		Error

	if err != nil {
		r.logger.Errorf("error updating userpassword with error: %v", err)
		return err
	}

	return nil
}

func (r *UserRepository) Delete(user domain.User) error {
	if err := r.db.Delete(&user).Error; err != nil {
		r.logger.Errorf("error deleting user: %s with error: %v", user.Username, err)
		return err
	}

	return nil
}

func (r *UserRepository) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("username %s not found as an active user", username)
		} else {
			r.logger.Errorf("error finging user with username %s with error %v", username, err)
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetById(userId uint) (domain.User, error) {
	var user domain.User

	if err := r.db.First(&user, user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) IncrementFailedLoginAttempt(userId uint, attempts int) error {
	if err := r.db.Model(&domain.User{}).Where("id = ?", userId).Update("failed_login_attempts", attempts).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ResetFailedLoginAttempt(userId uint) error {
	if err := r.db.Model(&domain.User{}).Where("id = ?", userId).Update("failed_login_attempts", 0).Error; err != nil {
		return err
	}

	return nil
}
