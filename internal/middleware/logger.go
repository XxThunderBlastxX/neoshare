package middleware

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (m *Middleware) StyledLogger(appEnv string) fiber.Handler {
	const filePerm uint32 = 0o644

	var logFile *os.File
	if appEnv == "dev" {
		file, err := os.OpenFile("./bin/requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(filePerm))
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		logFile = file
	} else {
		file, err := os.OpenFile("requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(filePerm))
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		logFile = file
	}

	// Sets the log output to the file and stdout
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	logConfig := logger.Config{
		Done: func(_ *fiber.Ctx, logString []byte) {
			fmt.Print(string(logString))
		},
		TimeFormat: "2006-01-02 15:04:05-0700",
		Output:     logFile,
	}

	return logger.New(logConfig)
}
