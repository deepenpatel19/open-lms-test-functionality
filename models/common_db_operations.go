package models

import (
	"context"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/open-lms-test-functionality/logger"
	"go.uber.org/zap"
)

type QueryStructToExecute struct {
	Query string
}

func (query QueryStructToExecute) InsertOrUpdateOperations(uuidString string, args ...any) (int64, error) {
	logger.Logger.Info("MODELS :: Will do insert operations", zap.String("requestId", uuidString), zap.String("query", query.Query), zap.Any("args", args))

	var id int64
	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err), zap.String("requestId", uuidString))
		return id, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = tx.QueryRow(ctx, query.Query, args...).Scan(&id)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return id, err
	}

	return id, nil
}

func (query QueryStructToExecute) DeleteOperation(uuidString string) (bool, error) {
	logger.Logger.Info("MODELS :: Will do delete operation.", zap.String("requestId", uuidString), zap.String("query", query.Query))
	dbConnection := DbPool()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadWrite})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err), zap.String("requestId", uuidString))
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	tx.Exec(ctx, query.Query)
	return true, nil
}
