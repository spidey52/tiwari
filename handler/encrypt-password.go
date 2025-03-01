package handler

import (
	"rmt/utils"

	"github.com/gin-gonic/gin"
)

type EncryptPasswordDto struct {
	Password string `json:"password"`
}

func EncryptPasswordHandler(c *gin.Context) {
	// encrypt password

	var encryptPasswordDto EncryptPasswordDto

	if err := c.ShouldBindJSON(&encryptPasswordDto); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	if encryptPasswordDto.Password == "" {
		c.JSON(400, gin.H{
			"error": "password is required",
		})

		return
	}

	encrypted, err := utils.EncryptPassword(encryptPasswordDto.Password)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":  "command executed",
		"password": encrypted,
	})
}
