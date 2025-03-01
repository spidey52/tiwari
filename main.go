package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"rmt/handler"
	"rmt/utils"
	"runtime"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/shlex"
)

const (
	decrypt_visibility = false
)

//go:embed ui/dist/*
//go:embed ui/dist/**/*
var reactFiles embed.FS

type RmtCommandDto struct {
	Filepath       string
	Command        string
	EncodedCommand string
	ArrayCommand   []string
}

// schtasks /create /sc onstart /tn "MyStartupTask" /tr "C:\Path\To\rmt.exe" /ru "SYSTEM"
type Settings struct {
	AllowedOrigins []string `json:"allowed_origins"`
}

// read json files
func getAllowedOrigins() (Settings, error) {
	default_allowed_origins := []string{"http://localhost:3000"}

	file, err := os.Open("settings.json")

	if err != nil {
		return Settings{AllowedOrigins: default_allowed_origins}, nil
	}

	defer file.Close()

	var settings Settings

	err = json.NewDecoder(file).Decode(&settings)

	if err != nil {
		return Settings{AllowedOrigins: default_allowed_origins}, nil
	}

	return settings, nil
}

func main() {

	allowedOrigins, err := getAllowedOrigins()

	if err != nil {
		fmt.Println(err)
		return
	}

	server := gin.Default()

	fmt.Println(allowedOrigins.AllowedOrigins)

	server.Use(cors.New(cors.Config{
		// AllowOrigins:     allowedOrigins.AllowedOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
	}))

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/drives", func(c *gin.Context) {

		veraDetails, err := utils.GetVeraDetails()

		if err != nil {
			c.IndentedJSON(500, gin.H{
				"error":   err.Error(),
				"message": "failed to get vera details",
			})
			return
		}

		c.IndentedJSON(200, gin.H{
			"drives":        veraDetails.Volumes,
			"backup_drives": veraDetails.BackupVolume,
		})
	})

	server.GET("/mounted-drives", func(c *gin.Context) {
		cmd := exec.Command("cmd", "/C", "wmic", "logicaldisk", "get", "name")

		output, err := cmd.CombinedOutput()

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		csvData := string(output)

		csvData = csvData[44:]

		csvData = csvData[:len(csvData)-2]

		csvData = csvData + "\n"

		c.JSON(200, gin.H{
			"drives": csvData,
		})

	})

	server.POST("/mount", handler.MountDriveHandler)
	server.POST("/unmount", handler.UnmountDriveHandler)
	server.POST("/shutdown", handler.ShutdownHandler)
	server.POST("/reboot", handler.RebootHandler)
	server.POST("/encrypt-password", handler.EncryptPasswordHandler)
	server.POST("/decrypt-password", func(c *gin.Context) {
		if !decrypt_visibility {
			c.JSON(400, gin.H{"error": "endpoint not found"})
			c.Abort()
			return
		}
		c.Next()
	}, handler.DecryptPasswordHandler)

	server.POST("/backup-vera", handler.BackupVeraHandler)

	server.POST("/rmt", func(c *gin.Context) {

		var rmtCommandDto RmtCommandDto

		if err := c.ShouldBindJSON(&rmtCommandDto); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if rmtCommandDto.Filepath != "" {
			_, err := os.Stat(rmtCommandDto.Filepath)

			if os.IsNotExist(err) {
				c.JSON(400, gin.H{"error": "file not found"})
				return
			}

			if os.IsPermission(err) {
				c.JSON(400, gin.H{"error": "permission denied"})
				return
			}

			var cmd *exec.Cmd

			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/C", rmtCommandDto.Filepath)
			} else {
				cmd = exec.Command(rmtCommandDto.Filepath)
			}

			output, err := cmd.CombinedOutput()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
				return
			}

			c.JSON(200, gin.H{
				"message": "file executed",
				"output":  string(output),
			})

			return
		}

		if rmtCommandDto.Command != "" {

			splittedCommand, err := shlex.Split(rmtCommandDto.Command)

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if len(splittedCommand) == 0 {
				c.JSON(400, gin.H{"error": "command is invalid"})
				return
			}

			if runtime.GOOS == "windows" {
				splittedCommand = append([]string{"cmd", "/C"}, splittedCommand...)
			}

			cmd := exec.Command(splittedCommand[0], splittedCommand[1:]...)

			output, err := cmd.CombinedOutput()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
				return
			}

			fmt.Println(string(output))

			c.JSON(200, gin.H{
				"message": "command executed",
				"output":  string(output),
			})

			return
		}

		if rmtCommandDto.EncodedCommand != "" {
			// decode base64
			decodedCommand, err := base64.StdEncoding.DecodeString(rmtCommandDto.EncodedCommand)

			fmt.Println(string(decodedCommand))

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			splittedCommand, err := shlex.Split(string(decodedCommand))

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if len(splittedCommand) == 0 {
				c.JSON(400, gin.H{"error": "command is invalid"})
				return
			}

			if runtime.GOOS == "windows" {
				splittedCommand = append([]string{"cmd", "/C"}, splittedCommand...)
			}

			cmd := exec.Command(splittedCommand[0], splittedCommand[1:]...)

			output, err := cmd.CombinedOutput()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
				return
			}

			fmt.Println(string(output))

			c.JSON(200, gin.H{
				"message": "command executed",
				"output":  string(output),
			})

		}

		if rmtCommandDto.ArrayCommand != nil {

			if len(rmtCommandDto.ArrayCommand) == 0 {
				c.JSON(400, gin.H{"error": "command is invalid"})
				return
			}

			if runtime.GOOS == "windows" {
				rmtCommandDto.ArrayCommand = append([]string{"cmd", "/C"}, rmtCommandDto.ArrayCommand...)
			}

			cmd := exec.Command(rmtCommandDto.ArrayCommand[0], rmtCommandDto.ArrayCommand[1:]...)

			output, err := cmd.CombinedOutput()

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
				return
			}

			fmt.Println(string(output))

			c.JSON(200, gin.H{
				"message": "command executed",
				"output":  string(output),
			})
			return
		}

		c.JSON(400, gin.H{"error": "file path or command is required"})

	})

	reactFS, err := fs.Sub(reactFiles, "ui/dist")

	if err != nil {
		fmt.Println(err)
		return
	}

	server.NoRoute(func(c *gin.Context) {
		fmt.Println("Request:", c.Request.Method, c.Request.URL.Path)

		filePath := c.Request.URL.Path

		// Default to index.html for root path
		if filePath == "/" {
			filePath = "index.html"
		} else {
			filePath = filePath[1:]
		}

		data, err := fs.ReadFile(reactFS, filePath)
		if err != nil {
			data, err = fs.ReadFile(reactFS, "index.html")
			if err != nil {
				c.JSON(404, gin.H{
					"error":   "file not found",
					"message": fmt.Sprintf("file %s not found", filePath),
				})
				return
			}

			filePath = "index.html"
		}

		// Determine content type based on file extension
		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = http.DetectContentType(data)
		}

		// Serve the file with the correct content type
		c.Data(200, contentType, data)
	})

	server.Run(":8089")
}
