package services

import (
	"authservice/src/config"
	"authservice/src/domain"
	"authservice/src/repositories"

	"go.uber.org/zap"
)

type ClaimService struct {
	claimRepo *repositories.ClaimRepository
	logger    *zap.SugaredLogger
}

func InitClaimService(serviceCfg *config.ServiceConfig) *ClaimService {
	return &ClaimService{
		claimRepo: repositories.InitClaimRepository(serviceCfg),
		logger:    serviceCfg.Logger,
	}
}

func (s *ClaimService) Add(claim domain.Claim) (domain.Claim, error) {
	if len(claim.Claim) <= 0 {
		s.logger.Warn("claim name cannot be empty")
	}

	return s.claimRepo.Add(claim)
}

func (s *ClaimService) Update(claim domain.Claim) error {
	if len(claim.Claim) <= 0 {
		s.logger.Warn("claim name cannot be empty")
	}

	return s.claimRepo.Update(claim)
}

func (s *ClaimService) Delete(claim domain.Claim) error {
	return s.claimRepo.Delete(claim)
}

func (s *ClaimService) GetAll() ([]domain.Claim, error) {
	return s.claimRepo.GetAll()
}
