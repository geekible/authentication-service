package config

import (
	"authservice/src/domain"

	"gorm.io/gorm"
)

type DatabaseMigration struct {
	db *gorm.DB
}

func InitDatabaseMigration(db *gorm.DB) *DatabaseMigration {
	return &DatabaseMigration{
		db: db,
	}
}

func (m *DatabaseMigration) DoMigration() error {
	if err := m.db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}

	if err := m.db.AutoMigrate(&domain.Claim{}); err != nil {
		return err
	}

	m.db.FirstOrCreate(&domain.Claim{
		Claim: "Administrator",
	}, 1)
	m.db.FirstOrCreate(&domain.Claim{
		Claim: "User",
	}, 2)

	if err := m.db.AutoMigrate(&domain.UserClaim{}); err != nil {
		return err
	}

	return nil
}
