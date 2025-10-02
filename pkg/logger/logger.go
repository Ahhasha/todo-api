package logger

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func Info(message string) {
	logger.Println(message)
}
