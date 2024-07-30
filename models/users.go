package models

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/open-lms-test-functionality/logger"
	"go.uber.org/zap"
)

const (
	STUDENT int = 1
	TEACHER int = 2
)

func ValidateUserType(userType int) string {
	if userType == TEACHER {
		return "teacher"
	} else if userType == STUDENT {
		return "student"
	} else {
		return ""
	}
}

func GetUserType(userType string) int {
	if userType == "teacher" {
		return TEACHER
	} else if userType == "student" {
		return STUDENT
	} else {
		return 0
	}
}

type UserCreateSchema struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
	Type      string `json:"type" form:"type"`
}

type UserSchema struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Type      string `json:"type"`
}

type UserResponseSchema struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Type      string `json:"type"`
}

func (data UserCreateSchema) Insert(uuidString string) (int64, error) {
	query := `INSERT INTO 
				users
					(first_name, last_name, email, password, type)
				VALUES
					($1, $2, $3, $4, $5)
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.FirstName, data.LastName, data.Email, data.Password, GetUserType(data.Type))
	return id, err
}

func (data UserCreateSchema) Update(uuidString string, id int64) (int64, error) {
	query := `UPDATE 
				users
					SET first_name=$1, last_name=$2
				WHERE id=$3
				RETURNING id`
	queryToExecute := QueryStructToExecute{Query: query}
	id, err := queryToExecute.InsertOrUpdateOperations(uuidString, data.FirstName, data.LastName, id)
	return id, err

}

func DeleteUserFromDB(uuidString string, userId int64) (bool, error) {
	query := fmt.Sprintf(`DELETE FROM users WHERE id=%d`, userId)
	queryToExecute := QueryStructToExecute{Query: query}
	status, err := queryToExecute.DeleteOperation(uuidString)
	return status, err
}

func FetchUserForAuth(email string) UserSchema {
	logger.Logger.Info("MODELS :: Will fetch user details for auth", zap.String("email", email))

	var userData UserSchema
	dbConnection := DbPool()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := dbConnection.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		logger.Logger.Error("MODELS :: Error while begin transaction", zap.Error(err))
		return userData
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := fmt.Sprintf(`SELECT
							u.id,
							u.first_name, 
							u.last_name, 
							u.email, 
							u.password,
							u.type
							FROM users u
							WHERE u.email='%s' LIMIT 1`, email)
	logger.Logger.Info("MODELS :: Query", zap.String("query", query))
	err = tx.QueryRow(ctx, query).Scan(
		&userData.Id,
		&userData.FirstName,
		&userData.LastName,
		&userData.Email,
		&userData.Password,
		&userData.Type,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Info("MODELS :: Query - No rows found. ", zap.String("query", query))
			return userData
		}
		logger.Logger.Error("MODELS :: Error while executing query.",
			zap.Error(err),
		)
		return userData
	}
	return userData
}
