package log

import (
	"database/sql"
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Adapter struct {
	Config         config.LogConfig
	DB             *gorm.DB
	names          map[string]string     //srv-captain--jsonbeautifier.1.qa9gcu6usinw06lqcfu286wsc -> srv-captain--jsonbeautifier
	lastTimestamps map[string]*time.Time //srv-captain--jsonbeautifier->"timestamp..."
}

func New(cfg config.LogConfig) (*Adapter, error) {
	dbLogger := logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(sqlite.Open(cfg.DB), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return nil, fmt.Errorf("open local database error: %w", err)
	}

	rawDB, _ := db.DB()

	if err := migrate(rawDB, cfg.DB); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}
	return &Adapter{
		Config:         cfg,
		DB:             db,
		names:          make(map[string]string, 10),
		lastTimestamps: make(map[string]*time.Time, 10),
	}, nil
}

func migrate(db *sql.DB, dbPath string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("error on set goose dialect: %w", err)
	}

	if err := goose.Up(db, "./configs/migrations/local"); err != nil {
		return fmt.Errorf("error on applying migrations: %w", err)
	}
	return nil
}
