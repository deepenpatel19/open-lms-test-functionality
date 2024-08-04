package models

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/open-lms-test-functionality/logger"
	"go.uber.org/zap"
)

type TestQuestionsSchema struct {
	Id           int64                  `json:"id"`
	TestId       int64                  `json:"test_id"`
	QuestionData QuestionResponseSchema `json:"question_data"`
}

type TestQuestionSchemaForTakeTest struct {
	Id           int64                             `json:"id"`
	TestId       int64                             `json:"test_id"`
	QuestionData QuestionResponseSchemaForTakeTest `json:"question_data"`
}

func CreateTestQuestionary(uuidString string, testId int64, questionIds []int64) (int64, error) {
	var listOfQuestions []string

	if len(questionIds) == 0 {
		return 0, nil
	}
	for _, questionId := range questionIds {
		query := fmt.Sprintf(`INSERT INTO
				test_questions
					(test_id, question_id)
				VALUES
					(%d, %d) RETURNING id`, testId, questionId)
		listOfQuestions = append(listOfQuestions, query)
	}
	queryToExecute := QueryStructToExecute{QueryList: listOfQuestions}
	id, err := queryToExecute.InsertOrUpdateMultipleQueries(uuidString)
	return id, err

}

func DeleteTestQuestionary(uuidString string, testId int64, questionId int64) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM test_questions WHERE test_id = %d AND question_id = %d`, testId, questionId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchTestQuestionaryForTeacher(uuidString string, testId int64, limit int, offset int) ([]TestQuestionsSchema, int, error) {
	logger.Logger.Info("MODELS :: Will fetch questions for teacher ", zap.String("requestId", uuidString))

	var data []TestQuestionsSchema
	var count int
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	// ctx := context.Background()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return data, count, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							tq.id,
							tq.test_id,
							q.id,
							q.type,
							q.question_data,
							q.answer_data,
							COUNT(*) OVER() AS total
							FROM test_questions tq
							JOIN tests t on t.id = tq.test_id
							JOIN questions q on q.id = tq.question_id
							WHERE tq.test_id = %d
							ORDER BY tq.id DESC LIMIT %d OFFSET %d`, testId, limit, offset)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("requestId", uuidString), zap.String("query", query))
			return data, count, nil
		}
		logger.Logger.Error("MODELS :: Error while fetching tests", zap.String("requestId", uuidString), zap.String("query", query), zap.Error(err))
		return data, count, err
	}
	defer rows.Close()

	for rows.Next() {
		var singleData TestQuestionsSchema
		err := rows.Scan(
			&singleData.Id,
			&singleData.TestId,
			&singleData.QuestionData.Id,
			&singleData.QuestionData.Type,
			&singleData.QuestionData.QuestionData,
			&singleData.QuestionData.AnswerData,
			&count,
		)
		if err != nil {
			logger.Logger.Error("MODELS :: Error while iterating rows", zap.String("requestId", uuidString), zap.Error(err))
			return data, count, err
		}

		data = append(data, singleData)
	}

	err = rows.Err()
	if err != nil {
		logger.Logger.Error("MODELS :: Error while at rows level", zap.String("requestId", uuidString), zap.Error(err))
		return data, count, err
	}

	return data, count, nil
}

func FetchTestQuestionaryForStrudent(uuidString string, testId int64, limit int, offset int) ([]TestQuestionSchemaForTakeTest, int, error) {
	logger.Logger.Info("MODELS :: Will fetch questions for student ", zap.String("requestId", uuidString))

	var data []TestQuestionSchemaForTakeTest
	var count int
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	// ctx := context.Background()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return data, count, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							tq.id,
							tq.test_id,
							q.id,
							q.type,
							q.question_data,
							COUNT(*) OVER() AS total
							FROM test_questions tq
							JOIN tests t on t.id = tq.test_id
							JOIN questions q on q.id = tq.question_id
							WHERE tq.test_id = %d
							ORDER BY tq.id DESC LIMIT %d OFFSET %d`, testId, limit, offset)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("requestId", uuidString), zap.String("query", query))
			return data, count, nil
		}
		logger.Logger.Error("MODELS :: Error while fetching tests", zap.String("requestId", uuidString), zap.String("query", query), zap.Error(err))
		return data, count, err
	}
	defer rows.Close()

	for rows.Next() {
		var singleData TestQuestionSchemaForTakeTest
		err := rows.Scan(
			&singleData.Id,
			&singleData.TestId,
			&singleData.QuestionData.Id,
			&singleData.QuestionData.Type,
			&singleData.QuestionData.QuestionData,
			&count,
		)
		if err != nil {
			logger.Logger.Error("MODELS :: Error while iterating rows", zap.String("requestId", uuidString), zap.Error(err))
			return data, count, err
		}

		data = append(data, singleData)
	}

	err = rows.Err()
	if err != nil {
		logger.Logger.Error("MODELS :: Error while at rows level", zap.String("requestId", uuidString), zap.Error(err))
		return data, count, err
	}

	return data, count, nil
}
