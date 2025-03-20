package handler

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"rmt/utils"
	"runtime"
	"strings"
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

func GetBasePath(filepath string) string {
	if strings.Contains(filepath, "\\") {
		filepath = strings.ReplaceAll(filepath, "\\", "/")
	}

	basePath := path.Base(filepath)
	return basePath
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
		fmt.Println("Volume", volume.DriveLetter, "Backup Drive", backupDriveDto.DriveLetter)
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

	// mount
	err = mountDrive(backupDriveDto.DriveLetter, backupDriveDto.Password)

	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "failed to mount drive",
		})
		return
	}

	//  vera drive path
	veraDrivePath := details.DriveLetter + ":\\" + "backup"

	// remove existing backup
	os.RemoveAll(veraDrivePath)
	fmt.Println("removing existing backup")

	// create backup folder
	os.MkdirAll(veraDrivePath, os.ModePerm)
	fmt.Println("creating backup folder")

	// copy files

	fmt.Println("Backup folder", details.BackupFolder, "Vera Drive Path", veraDrivePath)
	copyDirectory(details.BackupFolder, veraDrivePath)
	fmt.Println("copying files")

	unMountDrive(backupDriveDto.DriveLetter)

	drive_filename := time.Now().Format(time.DateTime) + "-" + strings.ToLower(details.VolumeName)
	backup_path := path.Join("mega:/vera-backup", drive_filename)

	commands := []string{
		"rclone",
		"copy",
		details.VolumePath,
		backup_path,
	}

	if runtime.GOOS == "windows" {
		commands = append([]string{"cmd", "/C"}, commands...)
	}

	fmt.Println("backup command", commands)
	cmd := exec.Command(commands[0], commands[1:]...)

	fmt.Println("backup started")
	output, err := cmd.CombinedOutput()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
		return
	}

	fmt.Println("backup completed")
	fmt.Println(string(output))

	c.JSON(200, gin.H{
		"message": "backup completed",
	})

}
