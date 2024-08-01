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

const (
	TRUEORFALSE       string = "true_or_false"
	MULTIPLECHOICE    string = "multiple_choice"
	TRUEORFALSEINT    int    = 1
	MULTIPLECHOICEINT int    = 2
)

func ValidateQuestionType(questionType string) string {
	if questionType == MULTIPLECHOICE {
		return "multiple_choice"
	} else if questionType == TRUEORFALSE {
		return "true_or_false"
	} else {
		return ""
	}
}

func GetQuestionType(questionType string) int {
	if questionType == MULTIPLECHOICE {
		return MULTIPLECHOICEINT
	} else if questionType == TRUEORFALSE {
		return TRUEORFALSEINT
	} else {
		return 0
	}
}

type QuestionCreateSchema struct {
	Type         string                 `json:"type"`
	QuestionData map[string]interface{} `json:"question_data"`
	AnswerData   map[string]interface{} `json:"answer_data"`
}

type QuestionResponseSchemaForTakeTest struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`
	QuestionData string `json:"question_data"`
}

type QuestionResponseSchema struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`
	QuestionData string `json:"question_data"`
	AnswerData   string `json:"answer_data"`
}

func (data QuestionCreateSchema) Insert(uuidString string) (int64, error) {
	questionData, err := json.Marshal(data.QuestionData)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while json marshalling question data ", zap.String("requestId", uuidString), zap.Error(err))
		return 0, err
	}
	answerData, err := json.Marshal(data.AnswerData)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while json marshalling answer data ", zap.String("requestId", uuidString), zap.Error(err))
		return 0, err
	}
	questionType := GetQuestionType(data.Type)
	query := `INSERT INTO
				questions
					(type, question_data, answer_data)
				VALUES
					($1, $2, $3)
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, questionType, string(questionData), string(answerData))
	return id, err
}

func (data QuestionCreateSchema) Update(uuidString string, questionId int64) (int64, error) {
	questionData, err := json.Marshal(data.QuestionData)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while json marshalling question data ", zap.String("requestId", uuidString), zap.Error(err))
		return 0, err
	}
	answerData, err := json.Marshal(data.AnswerData)
	if err != nil {
		logger.Logger.Error("MODELS :: Error while json marshalling answer data ", zap.String("requestId", uuidString), zap.Error(err))
		return 0, err
	}
	questionType := GetQuestionType(data.Type)
	query := `UPDATE
				questions
					set type=$1, question_data=$2, answer_data=$3
				WHERE id= $4
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, questionType, questionData, answerData, questionId)
	return id, err
}

func DeleteQuestion(uuidString string, questionId int64) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM questions WHERE id=%d`, questionId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchQuestion(uuidString string, questionId int64) (QuestionResponseSchema, error) {
	logger.Logger.Info("MODELS :: Will fetch test details ", zap.Int64("questionId", questionId), zap.String("requestId", uuidString))

	var questionData QuestionResponseSchema
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return questionData, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							q.id,
							q.type,
							q.question_data,
							q.answer_data
							FROM questions q
							WHERE q.id=%d LIMIT 1`, questionId)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))
	err = tx.QueryRow(ctx, query).Scan(
		&questionData.Id,
		&questionData.Type,
		&questionData.QuestionData,
		&questionData.AnswerData,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("query", query))
			return questionData, nil
		}
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.Error(err),
		)
		return questionData, err
	}
	return questionData, nil
}

func FetchQuestions(uuidString string, limit int, offset int) ([]QuestionResponseSchema, int, error) {
	logger.Logger.Info("MODELS :: Will fetch tests ", zap.String("requestId", uuidString))

	var questionsData []QuestionResponseSchema
	var count int
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return questionsData, count, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							q.id,
							q.type,
							q.question_data,
							q.answer_data,
							COUNT(*) OVER() AS total
							FROM questions q
							ORDER BY id DESC LIMIT %d OFFSET %d`, limit, offset)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query), zap.String("requestId", uuidString))

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("requestId", uuidString), zap.String("query", query))
			return questionsData, count, nil
		}
		logger.Logger.Error("MODELS :: Error while fetching tests", zap.String("requestId", uuidString), zap.String("query", query), zap.Error(err))
		return questionsData, count, err
	}
	defer rows.Close()

	for rows.Next() {
		var singleQuestionData QuestionResponseSchema
		err := rows.Scan(
			&singleQuestionData.Id,
			&singleQuestionData.Type,
			&singleQuestionData.QuestionData,
			&singleQuestionData.AnswerData,
			&count,
		)
		if err != nil {
			logger.Logger.Error("MODELS :: Error while iterating rows", zap.String("requestId", uuidString), zap.Error(err))
			return questionsData, count, err
		}

		questionsData = append(questionsData, singleQuestionData)
	}

	err = rows.Err()
	if err != nil {
		logger.Logger.Error("MODELS :: Error while at rows level", zap.String("requestId", uuidString), zap.Error(err))
		return questionsData, count, err
	}

	return questionsData, count, nil
}
