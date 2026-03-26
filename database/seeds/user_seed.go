package seeds

import (
	"workzen-be/internal/core/domain/model"
	"workzen-be/lib/conv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate password")
	}

	admin := model.User{
		TenantID: nil,
		Name:     "Admin",
		Email:    "admin@mail.com",
		Status:   "ACTIVE",
		Role:     "SUPER_ADMIN",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "admin@mail.com"}).Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to create user")
	} else {
		log.Info().Msg("Admin role seeded successfully")
	}
}
