package api

import (
	"github.com/gin-gonic/gin"
	"github.com/open-lms-test-functionality/utils"
)

func CreateUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)
}
