package logger

import "go.uber.org/zap"

func NewZapLogger() (Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.DisableStacktrace = true

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
