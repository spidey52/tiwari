package handler

import (
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
)

func ShutdownHandler(c *gin.Context) {
	// shutdown
	commands := []string{}

	if runtime.GOOS == "windows" {
		commands = append(commands, "cmd", "/C", "shutdown", "/s", "/t", "0")

	} else if runtime.GOOS == "linux" {
		commands = append(commands, "shutdown", "now")
	} else if runtime.GOOS == "darwin" {
		commands = append(commands, "shutdown", "-h", "now")
	}

	if len(commands) == 0 {
		c.JSON(400, gin.H{"error": "command is invalid"})
		return
	}

	cmd := exec.Command(commands[0], commands[1:]...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "output": string(output)})
		return
	}

	c.JSON(200, gin.H{
		"message": "command executed",
		"output":  string(output),
	})

}
