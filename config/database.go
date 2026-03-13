package config

import (
	"bwanews/database/seeds"
	"fmt"
	"time" // Tambahkan ini

	"github.com/rs/zerolog/log"
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
			log.Info().Msg("Successfully connected to database")
			break
		}

		log.Warn().
			Msgf("[Attempt %d/%d] Database not ready, retrying in 3 seconds... (Host: %s)", i, maxRetries, cfg.Psql.Host)

		time.Sleep(3 * time.Second) // Tunggu 3 detik sebelum coba lagi
	}
	// --- LOGIKA RETRY SELESAI ---

	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-Final] Failed to connect to database after retries")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-2] Error getting database connection")
		return nil, err
	}

	//log.Info().Msg("Running database migration...")
	//err = db.AutoMigrate(
	//	&model.User{},
	//	&model.Content{},
	//	&model.Client{},
	//	&model.Category{},
	//	&model.Employee{},
	//	&model.ManpowerRequest{},
	//	&model.Tenant{},
	//)
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to migrate database")
	//	return nil, err
	//}

	seeds.SeedRoles(db)

	sqlDB.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Psql.DBMaxIdle)

	return &Postgres{DB: db}, nil
}
