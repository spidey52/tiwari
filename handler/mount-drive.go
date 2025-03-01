package handler

import (
	"fmt"
	"os/exec"
	"rmt/utils"
	"runtime"

	"github.com/gin-gonic/gin"
)

type MountDriveDto struct {
	DriveLetter string `json:"drive_letter"`
	Password    string `json:"password"`
}

func mountDrive(driveLetter, password string) error {

	veraDetails, err := utils.GetVeraDetails()

	if err != nil {
		return err
	}

	decryptedPassword, err := utils.DecryptPassword(password)

	if err != nil {
		return err
	}

	for _, volume := range veraDetails.Volumes {
		if volume.DriveLetter == driveLetter {


			// fmt.Println(volume.VolumeName)
			// fmt.Println(volume.VolumePath)
			// fmt.Println(volume.DriveLetter)

			commands := []string{
				veraDetails.VeraCryptPath,
				"/v",
				volume.VolumePath,
				"/l",
				volume.DriveLetter,
				"/q",
				"/p",
				decryptedPassword,
				"/m",
				"rm",
			}

			if runtime.GOOS == "windows" {
				commands = append([]string{"cmd", "/C"}, commands...)
			}

			fmt.Println("mounting drive", commands, runtime.GOOS, volume.VolumePath, volume.DriveLetter)

			cmd := exec.Command(commands[0], commands[1:]...)

			output, err := cmd.CombinedOutput()

			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil

		}
	}

	return fmt.Errorf("drive letter not found")
}

func MountDriveHandler(c *gin.Context) {

	var mountDriveDto MountDriveDto

	if err := c.ShouldBindJSON(&mountDriveDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if mountDriveDto.DriveLetter == "" {
		c.JSON(400, gin.H{"error": "drive_letter is required"})
		return
	}

	if mountDriveDto.Password == "" {
		c.JSON(400, gin.H{"error": "password is required"})
		return
	}

	// mount drive
	err := mountDrive(mountDriveDto.DriveLetter, mountDriveDto.Password)
	fmt.Println(err)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "drive mounted",
		"drive":   mountDriveDto.DriveLetter,
	})

}
