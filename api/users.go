package api

import (
	"github.com/gin-gonic/gin"
	"github.com/open-lms-test-functionality/core"
	"github.com/open-lms-test-functionality/logger"
	"github.com/open-lms-test-functionality/models"
	"github.com/open-lms-test-functionality/schemas"
	"github.com/open-lms-test-functionality/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var userData models.UserCreateSchema
	if err := c.Bind(&userData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with user create schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}

	if userType := models.GetUserType(userData.Type); userType == 0 {
		c.JSON(400, gin.H{
			"message": "please check request body",
		})
		return

	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(userData.Password), core.Config.PasswordHashCost)
	if err != nil {
		logger.Logger.Error("API :: Error while generating password hash", zap.String("requestId", uuidString), zap.Error(err))
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}
	userData.Password = string(passwordHashBytes)

	id, err := userData.Insert(uuidString)
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

func UpdateUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var userData models.UserCreateSchema
	if err := c.Bind(&userData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with user update schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}

	id, err := userData.Update(uuidString, uri.UserId)
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

func DeleteUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	status, err := models.DeleteUserFromDB(uuidString, uri.UserId)
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
