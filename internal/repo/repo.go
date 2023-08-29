package repo

import (
	"database/sql"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

type Repo struct {
	DB *gorm.DB
}

func New(cfg config.DBConfig) (*Repo, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("current directory", path)

	path, err = os.Executable()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("current executable", path)

	dbLogger := logger.Default.LogMode(logger.Info)
	if !cfg.Log {
		dbLogger = logger.Discard
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBFile), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return nil, fmt.Errorf("open local database error: %w", err)
	}

	rawDB, _ := db.DB()

	if err := migrate(rawDB, "./configs/migrations/local"); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	return &Repo{DB: db}, nil
}

func migrate(db *sql.DB, dbPath string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("error on set goose dialect: %w", err)
	}

	if err := goose.Up(db, dbPath); err != nil {
		return fmt.Errorf("error on applying migrations: %w", err)
	}
	return nil
}
