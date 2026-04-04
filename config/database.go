package config

import (
	"fmt"
	"time" // Tambahkan ini
	"workzen-be/database/seeds"
	"workzen-be/internal/core/domain/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName)

	var db *gorm.DB
	var err error

	// --- LOGIKA RETRY DIMULAI ---
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
		if err == nil {
			// Jika koneksi berhasil, keluar dari loop
			log.Info("Successfully connected to database")
			break
		}

		log.Warnf("[Attempt %d/%d] Database not ready, retrying in 3 seconds... (Host: %s)", i, maxRetries, cfg.Psql.Host)

		time.Sleep(3 * time.Second) // Tunggu 3 detik sebelum coba lagi
	}
	// --- LOGIKA RETRY SELESAI ---

	if err != nil {
		log.Errorf("[ConnectionPostgres-Final] Failed to connect to database after retries: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("[ConnectionPostgres-2] Error getting database connection: %v", err)
		return nil, err
	}

	// Hybrid Migration Strategy:
	// Only run AutoMigrate in development/staging for faster iteration.
	// In production, we rely on manual SQL migrations for better control and safety.
	if cfg.App.AppEnv != "production" {
		log.Info("Running database auto-migration (Development/Staging)...")
		err = db.AutoMigrate(
			&model.User{},
			&model.Tenant{},
			&model.Client{},
			&model.Category{},
			&model.Content{},
			&model.Employee{},
			&model.ManpowerRequest{},
			&model.Candidate{},
			&model.CandidateApplication{},
			&model.Interview{},
			&model.Offer{},
			&model.EmployeeAssignment{},
			&model.SubscriptionPlan{},
			&model.TenantSubscription{},
		)
		if err != nil {
			log.Error("[ConnectionPostgres-Migration] Failed to auto-migrate database: ", err)
			return nil, err
		}
	} else {
		log.Info("Skipping database auto-migration (Production)...")
	}

	seeds.SeedRoles(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
