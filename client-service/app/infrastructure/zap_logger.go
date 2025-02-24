package infrastructure

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/yakob-abada/delfare/client-service/domain"
)

type ZapLogger struct {
	logger *zap.Logger
	closer func()
}

func NewZapLogger(production bool) (*ZapLogger, error) {
	var config zap.Config

	if production {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	closer := func() {
		_ = logger.Sync()
	}

	return &ZapLogger{logger: logger, closer: closer}, nil
}

func (z *ZapLogger) Debug(ctx domain.LogContext, msg string, fields ...interface{}) {
	z.logger.Debug(msg, z.buildFields(ctx, fields)...)
}

func (z *ZapLogger) Info(ctx domain.LogContext, msg string, fields ...interface{}) {
	z.logger.Info(msg, z.buildFields(ctx, fields)...)
}

func (z *ZapLogger) Warn(ctx domain.LogContext, msg string, fields ...interface{}) {
	z.logger.Warn(msg, z.buildFields(ctx, fields)...)
}

func (z *ZapLogger) Error(ctx domain.LogContext, msg string, fields ...interface{}) {
	z.logger.Error(msg, z.buildFields(ctx, fields)...)
}

func (z *ZapLogger) Fatal(ctx domain.LogContext, msg string, fields ...interface{}) {
	z.logger.Fatal(msg, z.buildFields(ctx, fields)...)
}

func (z *ZapLogger) Close() {
	if z.closer != nil {
		z.closer()
	}
}

// buildFields combines the log context and additional fields into a slice of zap.Field
func (z *ZapLogger) buildFields(ctx domain.LogContext, fields []interface{}) []zap.Field {
	var zapFields []zap.Field

	// Add context fields
	if ctx.RequestID != "" {
		zapFields = append(zapFields, zap.String("request_id", ctx.RequestID))
	}
	if ctx.CorrelationID != "" {
		zapFields = append(zapFields, zap.String("correlation_id", ctx.CorrelationID))
	}
	if ctx.UserID != "" {
		zapFields = append(zapFields, zap.String("user_id", ctx.UserID))
	}

	// Add additional fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 >= len(fields) {
			break
		}
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}

	return zapFields
}
