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
        // Production config but with debug level (keeps JSON format)
        config := zap.NewProductionConfig()
        config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
        Logger, err = config.Build()
    } else {
        Logger, err = zap.NewProduction()
    }
    
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer Logger.Sync()
}