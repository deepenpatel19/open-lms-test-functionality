package models

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"github.com/open-lms-test-functionality/core"
	"github.com/open-lms-test-functionality/logger"
	"go.uber.org/zap"
)

var DBPoolConnection *pgxpool.Pool

func CreateConnection() {
	databaseString := core.Config.DBString
	connConf, err := pgxpool.ParseConfig(databaseString)
	if err != nil {
		logger.Logger.Error("DATABASE :: Could not able to parse db config.", zap.Error(err))
	}
	connConf.MaxConnIdleTime = 600 * time.Second
	connConf.HealthCheckPeriod = 15 * time.Second
	connConf.MaxConnLifetime = 1800 * time.Second
	connConf.MinConns = 7
	connConf.MaxConns = 40
	connConf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol // For pgbouncer
	connConf.ConnConfig.DescriptionCacheCapacity = 1024                        // For pgbouncer

	dbpool, err := pgxpool.NewWithConfig(context.Background(), connConf)
	if err != nil {
		logger.Logger.Error("DATABASE :: Could not create pool.", zap.Error(err))
	}
	DBPoolConnection = dbpool
}

func DbPool() *pgxpool.Pool {
	return DBPoolConnection
}

func RunMigrations() {
	environment := core.Config.Environment
	databaseString := core.Config.DBString

	// Add connection timeout
	queryParamExists := strings.Contains(databaseString, "?")
	if queryParamExists {
		databaseString = databaseString + "&connect_timeout=30"
	} else {
		databaseString = databaseString + "?connect_timeout=30"
	}

	logger.Logger.Info("DATABASE :: run migrations ", zap.Any("environment", environment))
	var migrationFilesLocation string
	if environment == "local" || environment == "" {

		cwd, err := os.Getwd()
		if err != nil {
			logger.Logger.Error("DATABASE :: directory not found", zap.Error(err))
		}
		migrationFilesLocation = cwd + "/migrations"
	} else {
		migrationFilesLocation = "/usr/bin/migrations"
	}
	logger.Logger.Info("DATABASE :: SQL will open.", zap.Any("migration location", migrationFilesLocation))
	db1, err := sql.Open("pgx", databaseString)
	if err != nil {
		logger.Logger.Error("DATABASE :: DB is not initialized", zap.Error(err))
	}
	logger.Logger.Info("DATABASE :: SQL opened", zap.Any("db1", db1), zap.Any("db connector", db1.Ping()))

	driver, err := postgres.WithInstance(db1, &postgres.Config{})
	if err != nil {
		logger.Logger.Error("DATABASE :: Error while setup db driver", zap.Error(err))
		return
	}
	logger.Logger.Info("DATABASE :: driver setup", zap.Any("driver", driver))
	migrationDir := flag.String("migration.files", "./migrations", "./migrations")
	logger.Logger.Info("DATABASE :: Migration dir", zap.String("dir", *migrationDir), zap.String("path", migrationFilesLocation))

	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+migrationFilesLocation,
		"postgres", driver)
	if err != nil {
		logger.Logger.Error("DATABASE :: Error for new with db instance", zap.Error(err))
		return
	}

	m.Up()
	logger.Logger.Info("DATABASE :: Up command executed.")

	err = driver.Close()
	if err != nil {
		logger.Logger.Error("DATABASE :: Error while closing connection ", zap.Error(err))
	}
	logger.Logger.Info("DATABASE :: Driver connection closed successfully.")
}
