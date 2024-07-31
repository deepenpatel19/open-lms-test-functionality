package models

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/open-lms-test-functionality/logger"
	"go.uber.org/zap"
)

type TestCreateSchema struct {
	Title string `json:"title" form:"title"`
}

type TestResponseSchema struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

func (data TestCreateSchema) Insert(uuidString string) (int64, error) {
	query := `INSERT INTO
				tests
					(title)
				VALUES
					($1)
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.Title)
	return id, err
}

func (data TestCreateSchema) Update(uuidString string, testId int64) (int64, error) {
	query := `UPDATE
				tests
					set title=$1
				WHERE id=$2
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.Title, testId)
	return id, err
}

func DeleteTest(uuidString string, testId int64) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM tests WHERE id=%d`, testId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchTest(uuidString string, testId int64) (TestResponseSchema, error) {
	logger.Logger.Info("MODELS :: Will fetch test details ", zap.Int64("testId", testId), zap.String("requestId", uuidString))

	var testData TestResponseSchema
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return testData, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							t.id,
							t.title
							FROM tests t
							WHERE t.id=%d LIMIT 1`, testId)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))
	err = tx.QueryRow(ctx, query).Scan(
		&testData.Id,
		&testData.Title,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("query", query))
			return testData, nil
		}
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.Error(err),
		)
		return testData, err
	}
	return testData, nil
}

func FetchTests(uuidString string, limit int, offset int) ([]TestResponseSchema, int, error) {
	logger.Logger.Info("MODELS :: Will fetch tests ", zap.String("requestId", uuidString))

	var testData []TestResponseSchema
	var count int
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return testData, count, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							t.id,
							t.title,
							COUNT(*) OVER() AS total 
							FROM tests t
							ORDER BY id DESC LIMIT %d OFFSET %d`, limit, offset)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("requestId", uuidString), zap.String("query", query))
			return testData, count, nil
		}
		logger.Logger.Error("MODELS :: Error while fetching tests", zap.String("requestId", uuidString), zap.String("query", query), zap.Error(err))
		return testData, count, err
	}
	defer rows.Close()

	for rows.Next() {
		var singleTestData TestResponseSchema
		err := rows.Scan(
			&singleTestData.Id,
			&singleTestData.Title,
		)
		if err != nil {
			logger.Logger.Error("MODELS :: Error while iterating rows", zap.String("requestId", uuidString), zap.Error(err))
			return testData, count, err
		}

		testData = append(testData, singleTestData)
	}

	err = rows.Err()
	if err != nil {
		logger.Logger.Error("MODELS :: Error while at rows level", zap.String("requestId", uuidString), zap.Error(err))
		return testData, count, err
	}

	return testData, count, nil
}
