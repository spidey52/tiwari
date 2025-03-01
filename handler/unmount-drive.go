package handler

import (
	"fmt"
	"os/exec"
	"rmt/utils"
	"runtime"

	"github.com/gin-gonic/gin"
)

type UnmountDriveDto struct {
	DriveLetter string `json:"drive_letter"`
}

func unMountDrive(driveLetter string) error {
	veraDetails, err := utils.GetVeraDetails()

	if err != nil {
		return err
	}

	for _, volume := range veraDetails.Volumes {
		if volume.DriveLetter == driveLetter {

			// unmount drive
			commands := []string{
				veraDetails.VeraCryptPath,
				"/u",
				volume.DriveLetter,
				"/q",
			}

			if runtime.GOOS == "windows" {
				commands = append([]string{"cmd", "/C"}, commands...)
			}

			cmd := exec.Command(commands[0], commands[1:]...)

			cmd.CombinedOutput()

			return nil

		}
	}

	return fmt.Errorf("drive not found")

}

func UnmountDriveHandler(c *gin.Context) {

	var unmountDriveDto UnmountDriveDto

	if err := c.ShouldBindJSON(&unmountDriveDto); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if unmountDriveDto.DriveLetter == "" {
		c.JSON(400, gin.H{
			"error": "drive_letter is required",
		})
		return
	}

	// unmount drive
	veraDetails, err := utils.GetVeraDetails()

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, volume := range veraDetails.Volumes {
		if volume.DriveLetter == unmountDriveDto.DriveLetter {

			// unmount drive
			commands := []string{
				veraDetails.VeraCryptPath,
				"/u",
				volume.DriveLetter,
				"/q",
			}

			if runtime.GOOS == "windows" {
				commands = append([]string{"cmd", "/C"}, commands...)
			}

			cmd := exec.Command(commands[0], commands[1:]...)

			output, err := cmd.CombinedOutput()

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "drive unmounted",
				"output":  string(output),
			})

			return

		}
	}

	c.JSON(200, gin.H{
		"message": "drive not found",
	})
}
