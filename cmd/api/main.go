package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"re-partners/internal"
	"re-partners/internal/config"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func main() {
	log := initLogger()
	defer log.Sync() // flushes buffer, if any

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalw("loading config failed", "error", err)
	}

	if err := runMigrations(log, cfg.DatabaseURL); err != nil {
		log.Fatalw("migration failed", "error", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	defer pool.Close()

	e := internal.NewServer(log, pool).SetupRoutes()
	if err := e.Start(":" + cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("server error", "error", err)
	}
}

func runMigrations(log *zap.SugaredLogger, databaseURL string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsPath := "file://" + filepath.Join(wd, "migrations")

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return err
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Infof("migration source close error: %v", srcErr)
		}

		if dbErr != nil {
			log.Infof("migration db close error: %v", dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
