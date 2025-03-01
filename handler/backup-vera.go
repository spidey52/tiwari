package handler

import (
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"rmt/utils"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Preserve file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// copyDirectory recursively copies a source directory to a destination
func copyDirectory(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Construct the new path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// Create the directory if it doesn’t exist
			return os.MkdirAll(dstPath, info.Mode())
		} else {
			// Copy the file
			return copyFile(path, dstPath)
		}
	})
}

type BackupDriveDto struct {
	DriveLetter string `json:"drive_letter"`
	Password    string `json:"password"`
}

func BackupVeraHandler(c *gin.Context) {

	var backupDriveDto BackupDriveDto

	if err := c.ShouldBindJSON(&backupDriveDto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	veraDetails, err := utils.GetVeraDetails()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "message": "failed to get vera details"})
		return
	}

	idx := -1

	for i, volume := range veraDetails.BackupVolume {
		if volume.DriveLetter == backupDriveDto.DriveLetter {
			idx = i
			break
		}
	}

	if idx == -1 {
		c.JSON(400, gin.H{"error": "drive letter not found"})
		return
	}

	details := veraDetails.BackupVolume[idx]

	fileName := path.Base(details.VolumePath)

	if fileName == "." || fileName == "/" {
		c.JSON(400, gin.H{"error": "invalid file name"})
		return
	}

	// mount
	mountDrive(backupDriveDto.DriveLetter, backupDriveDto.Password)

	//  vera drive path
	veraDrivePath := details.DriveLetter + ":\\" + "backup"

	// remove existing backup
	os.RemoveAll(veraDrivePath)

	// create backup folder
	os.MkdirAll(veraDrivePath, os.ModePerm)

	// copy files

	copyDirectory(details.BackupFolder, veraDrivePath)

	// unmounting drive
	unMountDrive(backupDriveDto.DriveLetter)

	// rename file with timestamp for backing up, this filename will be on the cloud
	fileName = time.Now().Format(time.DateTime) + "-" + fileName

	backup_path := path.Join("mega:/vera-backup", fileName)

	commands := []string{
		"rclone",
		"copy",
		details.VolumePath,
		backup_path,
	}

	if runtime.GOOS == "windows" {
		commands = append([]string{"cmd", "/C"}, commands...)
	}

	cmd := exec.Command(commands[0], commands[1:]...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
		return
	}

	c.JSON(200, gin.H{
		"message": "backup completed",
	})

}
