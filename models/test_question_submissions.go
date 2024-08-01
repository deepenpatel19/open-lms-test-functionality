package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/open-lms-test-functionality/logger"
	"go.uber.org/zap"
)

type TestQuestionSubmissionSchema struct {
	Id            int64                  `json:"id"`
	UserId        int64                  `json:"user_id"`
	TestId        int64                  `json:"test_id"`
	QuestionData  QuestionResponseSchema `json:"question"`
	SubmittedData string                 `json:"submitted_data"`
	AnswerStatus  bool                   `json:"answer_status"`
}

func FetchTestQuestionSubmissions(uuidString string, testId int64, userId int64, limit int, offset int) ([]TestQuestionSubmissionSchema, int, error) {
	logger.Logger.Info("MODELS :: Will fetch test question submissions ", zap.String("requestId", uuidString), zap.Int64("testId", testId), zap.Int64("userId", userId))

	var data []TestQuestionSubmissionSchema
	var count int
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
							tq.user_id,
							tq.submitted_data,
							tq.answer_status,
							q.id,
							q.type,
							q.question_data,
							q.answer_data,
							COUNT(*) OVER() AS total
							FROM test_question_submissions tq
							JOIN questions q on q.id = tq.question_id
							WHERE tq.test_id = %d AND tq.user_id = %d
							ORDER BY id DESC LIMIT %d OFFSET %d`, testId, userId, limit, offset)
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
		var singleData TestQuestionSubmissionSchema
		err := rows.Scan(
			&singleData.Id,
			&singleData.TestId,
			&singleData.UserId,
			&singleData.SubmittedData,
			&singleData.AnswerStatus,
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

func CreateOrUpdateTestQuestionSubmission(uuidString string, testId int64, userId int64, questionId int64, answerData map[string][]string, questionData QuestionResponseSchema) (int64, error) {
	logger.Logger.Info("MODELS :: Will create or update test question submission data for student", zap.String("requestId", uuidString), zap.Int64("testId", testId), zap.Int64("userId", userId), zap.Int64("questionId", questionId), zap.Any("answerData", answerData), zap.Any("questionanswer data", questionData.AnswerData))

	var id int64
	var testQuestionSubmissionQuery string
	var answerStatus bool
	var questionAnswerData string
	var answerDataChoices []string
	var questionAnswerChoices []string
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

	questionAnswerDataQuery := fmt.Sprintf(`SELECT q.answer_data FROM questions q WHERE q.id=%d`, questionId)
	err = tx.QueryRow(ctx, questionAnswerDataQuery).Scan(&questionAnswerData)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing fetch question answer data query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return id, err
	}

	var testQuestionUnmarshal map[string][]string
	_ = json.Unmarshal([]byte(questionAnswerData), &testQuestionUnmarshal)
	logger.Logger.Debug("MODELS :: question answer data ", zap.Any("data", questionAnswerData), zap.Any("unmarshal", testQuestionUnmarshal))

	answerDataChoices = answerData["answer_data"]
	questionAnswerChoices = testQuestionUnmarshal["choices"]
	logger.Logger.Debug("MODELS :: question choices", zap.Any("choices ", answerDataChoices), zap.Any("raw", answerData["answer_data"]))
	var correctAnswer int
	for _, givenAnswer := range answerDataChoices {
		for _, actualAnswer := range questionAnswerChoices {
			logger.Logger.Debug("MODELS :: checking ", zap.Any("given", givenAnswer), zap.Any("actual", actualAnswer))
			if givenAnswer == actualAnswer {
				correctAnswer += 1
			}
		}
	}
	if correctAnswer == len(questionAnswerChoices) {
		answerStatus = true
	}

	logger.Logger.Debug("MODELS :: question answer", zap.Any("answer ", answerStatus), zap.Any("count", correctAnswer), zap.Any("a", len(questionAnswerChoices)))

	selectTestQuestionSubmissionQuery := fmt.Sprintf(`SELECT
														id
													FROM
														test_question_submissions
													WHERE test_id=%d AND user_id = %d AND question_id = %d`, testId, userId, questionId)

	err = tx.QueryRow(ctx, selectTestQuestionSubmissionQuery).Scan(&id)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return id, err
	}

	if id > 0 {
		testQuestionSubmissionQuery = fmt.Sprintf(`UPDATE
													test_question_submissions
												SET submitted_data='%s', answer_status=%t
												WHERE test_id = %d AND user_id = %d AND question_id = %d
												`, answerData, answerStatus, testId, userId, questionId)

	} else {
		testQuestionSubmissionQuery = fmt.Sprintf(`INSERT INTO
													test_question_submissions
														(test_id, user_id, question_id, submitted_data, answer_status)
													VALUES
														(%d, %d, %d, '%s', %t)`,
			testId, userId, questionId, answerData, answerStatus)
	}

	err = tx.QueryRow(ctx, testQuestionSubmissionQuery).Scan(&id)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		return id, err
	}
	return id, nil
}
