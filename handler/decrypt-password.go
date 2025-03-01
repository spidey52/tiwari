package handler

import (
	"rmt/utils"

	"github.com/gin-gonic/gin"
)

type DecryptPasswordDto struct {
	Password string `json:"password"`
}

func DecryptPasswordHandler(c *gin.Context) {
	// decrypt password

	var decryptPasswordDto DecryptPasswordDto

	if err := c.ShouldBindJSON(&decryptPasswordDto); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	if decryptPasswordDto.Password == "" {
		c.JSON(400, gin.H{
			"error": "password is required",
		})
	}

	decrypted, err := utils.DecryptPassword(decryptPasswordDto.Password)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message":  "command executed",
		"password": decrypted,
	})

}
