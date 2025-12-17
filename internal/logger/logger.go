package logger

import (
    "go.uber.org/zap"
)

var Log *zap.Logger

func Init() error {
    var err error
    // Production configuration
    config := zap.NewProductionConfig()
    config.DisableStacktrace = true
    
    Log, err = config.Build()
    if err != nil {
        return err
    }
    
    zap.ReplaceGlobals(Log)
    return nil
}

func Sync() {
    if Log != nil {
        _ = Log.Sync()
    }
}