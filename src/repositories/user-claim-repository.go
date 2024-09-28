package repositories

import (
	"authservice/src/config"
	"authservice/src/domain"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserClaimRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func InitUserClaimRepository(servicecfg *config.ServiceConfig) *UserClaimRepository {
	return &UserClaimRepository{
		db:     servicecfg.Db,
		logger: servicecfg.Logger,
	}
}

func (r *UserClaimRepository) Add(userClaim domain.UserClaim) error {
	if err := r.db.Create(&userClaim).Error; err != nil {
		r.logger.Errorf("error creating claim for user id %d", userClaim.UserId)
		return err
	}

	return nil
}

func (r *UserClaimRepository) Delete(userClaim domain.UserClaim) error {
	if err := r.db.Delete(&userClaim).Error; err != nil {
		r.logger.Errorf("error creating claim for user id %d", userClaim.UserId)
		return err
	}

	return nil
}

func (r *UserClaimRepository) GetClaimsByUserId(userId uint) ([]string, error) {
	var userClaims []string

	if err := r.db.Where(&userClaims, "user_id = ?", userId).Error; err != nil {
		return []string{}, err
	}

	rows, err := r.db.
		Table("claims").
		Select("claims.claim").
		Joins("inner join user_claims on claims.id = user_claims.claim_id").
		Rows()
	if err != nil {
		r.logger.Errorf("error locating user claims with error %v", err)
		return []string{}, errors.New("error locating user claims")
	}

	for rows.Next() {
		var claim string
		rows.Scan(&claim)
		userClaims = append(userClaims, claim)
	}

	return userClaims, nil
}
