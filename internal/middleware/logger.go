package middleware

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (m *Middleware) StyledLogger(appEnv string) fiber.Handler {
	var logFile *os.File
	if appEnv == "dev" {
		file, err := os.OpenFile("./bin/requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		logFile = file
	} else {
		file, err := os.OpenFile("requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		logFile = file
	}

	logConfig := logger.Config{
		Done: func(c *fiber.Ctx, logString []byte) {
			fmt.Print(string(logString))
		},
		Output: logFile,
	}

	return logger.New(logConfig)
}
