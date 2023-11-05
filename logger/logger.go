package logger

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger

func InitLogger() {
    var err error
    Logger, err = zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer Logger.Sync()
}