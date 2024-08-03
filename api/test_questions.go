package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/open-lms-test-functionality/logger"
	"github.com/open-lms-test-functionality/models"
	"github.com/open-lms-test-functionality/schemas"
	"github.com/open-lms-test-functionality/utils"
	"go.uber.org/zap"
)

func CreateTestQuestionary(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	user, _ := c.Get("id")
	userEmail := user.(*models.UserSchema).Email

	userDataFromDb := models.FetchUserForAuth(userEmail)

	userType, err := strconv.Atoi(userDataFromDb.Type)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	userTypeStr := models.ValidateUserType(userType)
	if userTypeStr == "" || userTypeStr == "student" {
		c.JSON(400, gin.H{
			"message": "you're not allowed for this operation",
		})
		return
	}

	allQuestions, count, err := models.FetchQuestions(uuidString, 50, 0, true)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	if count == 0 {
		c.JSON(400, gin.H{
			"message": "no questions to create questionary",
		})
		return
	}

	var questionIds []int64
	for _, questionData := range allQuestions {
		questionIds = append(questionIds, questionData.Id)
	}

	_, err = models.CreateTestQuestionary(uuidString, uri.TestId, questionIds)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "created",
	})
}

func AddTestQuestion(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	user, _ := c.Get("id")
	userEmail := user.(*models.UserSchema).Email

	userDataFromDb := models.FetchUserForAuth(userEmail)

	userType, err := strconv.Atoi(userDataFromDb.Type)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	userTypeStr := models.ValidateUserType(userType)
	if userTypeStr == "" || userTypeStr == "student" {
		c.JSON(400, gin.H{
			"message": "you're not allowed for this operation",
		})
		return
	}

	_, err = models.CreateTestQuestionary(uuidString, uri.TestId, []int64{uri.QuestionId})
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "created",
	})

}

func DeleteTestQuestion(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	user, _ := c.Get("id")
	userEmail := user.(*models.UserSchema).Email

	userDataFromDb := models.FetchUserForAuth(userEmail)

	userType, err := strconv.Atoi(userDataFromDb.Type)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	userTypeStr := models.ValidateUserType(userType)
	if userTypeStr == "" || userTypeStr == "student" {
		c.JSON(400, gin.H{
			"message": "you're not allowed for this operation",
		})
		return
	}

	status, err := models.DeleteTestQuestionary(uuidString, uri.TestId, uri.QuestionId)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": status,
	})

}

func FetchTestQuestionary(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	user, _ := c.Get("id")
	userEmail := user.(*models.UserSchema).Email

	userDataFromDb := models.FetchUserForAuth(userEmail)

	userType, err := strconv.Atoi(userDataFromDb.Type)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	limitQuery := c.DefaultQuery("limit", "0")
	offsetQuery := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitQuery)
	offset, _ := strconv.Atoi(offsetQuery)

	if limit > 50 {
		c.JSON(400, gin.H{
			"message": "please check query params - param should not greater than 50",
		})
		return
	}
	if limit == 0 {
		limit = 10
	}

	userTypeStr := models.ValidateUserType(userType)
	if userTypeStr == "student" {
		data, count, err := models.FetchTestQuestionaryForStrudent(uuidString, uri.TestId, limit, offset)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "something went wrong",
			})
			return
		}

		if count == 0 {
			emptyArray := make([]string, 0)
			c.JSON(200, gin.H{
				"message": emptyArray,
				"count":   count,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": data,
			"count":   count,
		})
		return
	} else if userTypeStr == "teacher" {
		data, count, err := models.FetchTestQuestionaryForTeacher(uuidString, uri.TestId, limit, offset)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "something went wrong",
			})
			return
		}

		if count == 0 {
			emptyArray := make([]string, 0)
			c.JSON(200, gin.H{
				"message": emptyArray,
				"count":   count,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": data,
			"count":   count,
		})
		return
	} else {
		c.JSON(400, gin.H{
			"message": "you're not allowed for this operation",
		})
		return
	}

}
