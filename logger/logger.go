package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
    var err error
    
    if os.Getenv("LOG_LEVEL") == "debug" {
        Logger, err = zap.NewDevelopment()
    } else {
        Logger, err = zap.NewProduction()
    }
    
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer Logger.Sync()
}