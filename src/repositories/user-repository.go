package repositories

import (
	"authservice/src/domain"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func InitUserRepositoy(db *gorm.DB, logger *zap.SugaredLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) Add(user domain.User) (domain.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		r.logger.Errorf("error creating user: %s with error: %v", user.Username, err)
		return user, err
	}

	return user, nil
}

func (r *UserRepository) Update(user domain.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		r.logger.Errorf("error updating user: %s with error: %v", user.Username, err)
		return err
	}

	return nil
}

func (r *UserRepository) GetByUsernameAndPassword(username, password string) (domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "username = ? and password = ?").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warnf("username %s not found as an active user", username)
		} else {
			r.logger.Errorf("error finging user with username %s with error %v", username, err)
		}

		return domain.User{}, err
	}

	return user, nil
}
