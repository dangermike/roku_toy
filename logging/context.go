package logging

import (
	"context"

	"go.uber.org/zap"
)

type keyType struct{}

var key keyType

func Configure(debug bool) *zap.Logger {
	var l zap.Config
	if debug {
		l = zap.NewDevelopmentConfig()
	} else {
		l = zap.NewProductionConfig()
	}
	l.DisableCaller = true
	logger, err := l.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	return logger
}

func NewContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, key, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger := ctx.Value(key).(*zap.Logger)
	if logger == nil {
		return zap.L()
	}
	return logger
}
