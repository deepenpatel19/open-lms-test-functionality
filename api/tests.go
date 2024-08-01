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

func CreateTest(c *gin.Context) {
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

	var testData models.TestCreateSchema
	if err := c.Bind(&testData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with test create schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}

	id, err := testData.Insert(uuidString)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": id,
	})

}

func UpdateTest(c *gin.Context) {
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

	var testData models.TestCreateSchema
	if err := c.Bind(&testData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with test create schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}

	id, err := testData.Update(uuidString, uri.TestId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": id,
	})

}

func DeleteTest(c *gin.Context) {
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

	status, err := models.DeleteTest(uuidString, uri.TestId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": status,
	})

}

func FetchTest(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	testData, err := models.FetchTest(uuidString, uri.TestId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": testData,
	})
}

func FetchTests(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
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

	testData, count, err := models.FetchTests(uuidString, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{
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
		"message": testData,
		"count":   count,
	})

}
