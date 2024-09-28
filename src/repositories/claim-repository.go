package repositories

import (
	"authservice/src/config"
	"authservice/src/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ClaimRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func InitClaimRepository(serviceCfg *config.ServiceConfig) *ClaimRepository {
	return &ClaimRepository{
		db:     serviceCfg.Db,
		logger: serviceCfg.Logger,
	}
}

func (r *ClaimRepository) Add(claim domain.Claim) (domain.Claim, error) {
	if err := r.db.Create(&claim).Error; err != nil {
		r.logger.Errorf("error creating claim %+v with error %v", claim, err)
		return claim, err
	}

	return claim, nil
}

func (r *ClaimRepository) Update(claim domain.Claim) error {
	if err := r.db.Save(&claim).Error; err != nil {
		r.logger.Errorf("error updating claim %+v with error %v", claim, err)
		return err
	}

	return nil
}

func (r *ClaimRepository) Delete(claim domain.Claim) error {
	if err := r.db.Delete(&claim).Error; err != nil {
		r.logger.Errorf("error deleting claim %+v with error %v", claim, err)
		return err
	}

	return nil
}

func (r *ClaimRepository) GetAll() ([]domain.Claim, error) {
	var claims []domain.Claim

	if err := r.db.Find(&claims).Error; err != nil {
		r.logger.Errorf("error getting all claims %v", err)
		return []domain.Claim{}, err
	}

	return claims, nil
}
