package domain

type Logger interface {
	Debug(ctx LogContext, msg string, fields ...interface{})
	Info(ctx LogContext, msg string, fields ...interface{})
	Warn(ctx LogContext, msg string, fields ...interface{})
	Error(ctx LogContext, msg string, fields ...interface{})
	Fatal(ctx LogContext, msg string, fields ...interface{})
}
